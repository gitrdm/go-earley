# Enhancing go-earley with Grammar-IR, miniKanren, and Z3

This document outlines a design to augment the existing `patrickhuber/go-earley` library by:

1. Introducing a **generic Grammar-IR** layer  
2. Integrating a **miniKanren** engine for relational AST extraction and generation  
3. Hooking in **Z3** for heavy semantic constraints  
4. Dog-fooding via the existing PDL importer to bootstrap additional front-ends

---

## 1. Goals & Scope

- **Unify** multiple grammar dialects (PDL, EBNF, YACC, ANTLR, JSON) into one in-memory IR  
- **Decouple** parsing (Earley chart) from semantic work (AST assembly, type checks, constraint solving)  
- **Leverage** logic programming (miniKanren) to walk shared-packed parse forests, prune by semantics, and even *generate* code fragments  
- **Offload** complex, global checks (e.g. name-binding, non-linear arithmetic, user‐defined invariants) to Z3  
- **Preserve** the existing performance and robustness benefits (Leo, SPPF, memoization)

---

## 2. High-Level Architecture

```text
[ Grammar Files ]  ──> Importers ──> [ Grammar-IR ] ──┐
                                                     │
[ go-earley Scanner ] ──> Earley Parser ──> Chart/SPPF │
                                                     ├─> Kanren AST Extraction ──> AST+
                                                     │      (backtracking + SMT)
                                                     │
                                                     └─> Datalog/CI Passes, Codegen, REPL
```

### 2.1 Layers

1. **Importers**  
   - Existing `pdl/pdl.pdl` remains a PDL importer.  
   - New importers (EBNF, YACC, JSON) produce the same Go structs of Grammar-IR.

2. **Grammar-IR**  
   - Go package `grammarir/` defines:
     ```go
     type Terminal    struct { Name, Pattern string; Skip bool }
     type NonTerminal struct { Name string }
     type Production  struct {
       LHS  string
       RHS  []string
       Meta map[string]interface{}
     }
     type Precedence struct { Level int; Symbols []string; Assoc string }
     type Grammar struct {
       Terminals    []Terminal
       NonTerminals []NonTerminal
       Productions  []Production
       Precedences  []Precedence
       Start        string
     }
     ```

3. **Core Parser (go-earley)**  
   - Accepts a `Grammar` instance to build its `Rules`, `DottedRules`, `Chart`, and SPPF.  
   - Retains nullability checking and Leo/Aycock–Horspool optimizations.

4. **Relational AST Extraction (miniKanren)**  
   - Go port (e.g. `github.com/awalterschulze/gominikanren`) embedded as `kanren/`.  
   - For each nonterminal `N`, generate a goal:
     ```go
     func GoalN(i, j int, out *ASTNode) Goal { … }
     ```
   - A master goal does:
     ```go
     runAll(func(tree ASTNode) {
       GoalStart(0, len(tokens), &tree)
       // optionally: smtCheck(tree)
       emit(tree)
     })
     ```

5. **Heavy Constraint Solving (Z3)**  
   - Use Go bindings (e.g. `github.com/mitchellh/go-z3`).  
   - Expose a Kanren primitive:
     ```go
     func SMTCheck(cs …Constraint) Goal { … }
     ```
   - Invoke inline during AST extraction to prune invalid parses early.

6. **Downstream Passes**  
   - (Optional) Datalog or custom Go visitors for name resolution, type checking, codegen.  
   - Integration with CI to validate grammar changes, re-gen code, run tests.

---

## 3. Detailed Workflow

1. **Bootstrapping Importers**  
   - Start by using the existing PDL grammar (`pdl/pdl.pdl`) and go-earley to parse new importer grammars, e.g. `ebnf/ebnf.pdl`.  
   - Map their parse trees into `grammarir.Grammar` via a visitor in Go.

2. **Grammar-IR Construction**  
   - In CLI or library entrypoint:
     ```go
     raw, _ := os.ReadFile("grammar.pdl")
     syntaxTree, _ := dsl.Parse(raw)
     ir := grammarir.FromPDL(syntaxTree)
     ```

3. **Table Generation & Parsing**  
   - Build `Rules`, compute nullable set and right-recursive symbols.  
   - Stream tokens from a lexer into `parser.Parse(ir, tokens)` to produce an SPPF.

4. **AST Extraction via Kanren**  
   - Code-generate Go stubs for each nonterminal:
     ```sh
     go run cmd/codegen gen_go_ir grammar.json > ir_gen.go
     ```
   - In `parser/ast.go`:
     ```go
     func extractAST(chart Chart) []ASTNode {
       var results []ASTNode
       runAll(func(root ASTNode) {
         GoalStart(0, chart.Length(), &root)(Env{}, func(_ Env) {
           results = append(results, root)
         })
       })
       return results
     }
     ```

5. **SMT-Backed Semantic Filters**  
   - Decorate IR productions with metadata:
     ```pdl
     PROD: Expr [ Expr "+" Expr ] { z3: "(> left right)" }
     ```
   - During AST extraction, read `Meta["z3"]`, `SMTCheck` the constraint.

6. **Testing & CI**  
   - Unit tests for each importer: grammar → IR correctness.  
   - Fuzz tests: run parser+Kanren on generated strings.  
   - End-to-end tests: parse real PDL files and validate expected ASTs.

---

## 4. Implementation Plan

1. Create a new Go module `grammarir/` with the core IR types and helper builders.  
2. Refactor existing PDL importer (`dsl.Parse`) to output `grammarir.Grammar`.  
3. Wire `go-earley` core to accept `grammarir.Grammar` instead of internal `dsl.Grammar`.  
4. Add MiniKanren (`kanren/`) and implement codegen stubs from `grammarir.Grammar`.  
5. Integrate Z3 binding and provide an API to embed constraints in IR/productions.  
6. Update CLI (`cmd/earley/`) with new flags (`--ir-json`, `--generate-kanren`, `--smt`).  
7. Extend tests: importer specs, AST comparisons, SMT‐pruning unit tests.  
8. Document dog-fooding process: show how to write an `ebnf.pdl` and bootstrap its importer.

---

## 5. Risks & Mitigations

- **Performance overhead** of Kanren backtracking & SMT calls  
  - Mitigate by caching Z3 contexts and pruning as early as possible.  
  - Limit Kanren goals to ambiguous/non-nullable spans only.

- **Complexity of dogfooding importers**  
  - Start with one additional importer (EBNF) before extending to YACC/ANTLR.  
  - Provide clear examples and ship a “meta-grammar.pdl” that describes the IR itself.

- **User onboarding**  
  - Ship templates and a guided tutorial in `docs/`.  
  - Offer default IR→Go codegen so users don’t write boilerplate.

---

## 6. Next Steps

1. **Spike**: prototype `grammarir/` and refactor PDL importer.  
2. **PoC**: codegen a simple two-rule grammar to Go Kanren stubs and solve AST.  
3. **Benchmark**: compare native go-earley parse vs. parse+Kanren for small test cases.  
4. **Review**: iterate on design with maintainers and early adopters.  
5. **Merge**: roll out incremental PRs—first IR refactor, then Kanren, then SMT.

---

*Prepared by @gitrdm on 2025-10-22*  