# Enhancing go-earley with Grammar-IR, miniKanren, Z3, and Prolog

This document outlines a design to augment the existing `patrickhuber/go-earley` library by:

1. Introducing a **generic Grammar-IR** layer (Go structs with multiple input formats)
2. Integrating a **miniKanren** engine for relational AST extraction and generation  
3. Hooking in **Z3** for heavy semantic constraints
4. Adding **Prolog** (via `ichiban/prolog`) for global semantic analysis (name resolution, control-flow, data-flow)
5. Using **HCL** as the primary authoring format (parsed via external `hashicorp/hcl/v2` library)

---

## 1. Goals & Scope

- **Unify** multiple grammar dialects (PDL, EBNF, YACC, ANTLR, HCL, JSON) into one in-memory IR  
- **Decouple** parsing (Earley chart) from semantic work (AST assembly, type checks, constraint solving)  
- **Grammar-IR defined as Go structs** - external formats (HCL, PDL, etc.) parse into these structs
- **Use HCL as primary authoring format** - leverages external `hashicorp/hcl/v2` library (no bootstrapping needed)
- **Support legacy PDL format** - existing go-earley grammars can be imported alongside HCL
- **Leverage** logic programming (miniKanren) to walk shared-packed parse forests, prune by semantics, and even *generate* code fragments  
- **Offload** local constraints (type checks, arithmetic) to Z3 SMT solver
- **Perform** global semantic analysis (name resolution, control-flow, data-flow) with Prolog (tabled evaluation via `ichiban/prolog`)
- **Preserve** the existing performance and robustness benefits (Leo, SPPF, memoization)

---

## 2. High-Level Architecture

```text
┌──────────────────────────────────────────────────────────────────┐
│ External Grammar Source Files                                    │
├──────────────────────────────────────────────────────────────────┤
│  calc.grammar.hcl   │  legacy.pdl   │  grammar.ebnf  │ lang.y    │
│   (PRIMARY)         │  (LEGACY)     │  (IMPORT)      │ (IMPORT)  │
└─────────┬───────────┴──────┬────────┴───────┬────────┴───────┬───┘
          │                  │                │                │
          ↓                  ↓                ↓                ↓
    ┌─────────┐       ┌──────────┐    ┌───────────┐   ┌──────────┐
    │   HCL   │       │   PDL    │    │   EBNF    │   │   YACC   │
    │ Parser  │       │ Importer │    │ Importer  │   │ Importer │
    └────┬────┘       └─────┬────┘    └─────┬─────┘   └─────┬────┘
         │                  │               │               │
         └──────────────────┴───────────────┴───────────────┘
                            ↓
      ┌──────────────────────────────────────────────────────┐
      │         Grammar-IR (Canonical In-Memory)             │
      │  - Go structs: Grammar, Symbol, Production, etc.     │
      │  - All external formats parse into these structs     │
      │  - Format-agnostic internal representation           │
      └─────────────────────┬────────────────────────────────┘
                            ↓
      ┌─────────────────────────────────────────────────┐
      │      Earley Parser (EXISTING, Enhanced)         │
      │  Input:  Grammar-IR + Token Stream              │
      │  Output: SPPF (Shared-Packed Parse Forest)      │
      │  Note:   SPPF contains ALL valid parses         │
      └─────────────────────┬───────────────────────────┘
                            ↓
      ┌─────────────────────────────────────────────────┐
      │    AST Extractor (NEW - miniKanren-based)       │
      │  Input:  SPPF + Grammar-IR metadata             │
      │  Process: Logic programming traversal           │
      │  Output: Candidate AST(s) - may be multiple     │
      └─────────────────────┬───────────────────────────┘
                            ↓
      ┌─────────────────────────────────────────────────┐
      │   Semantic Validator (NEW - Z3-based)           │
      │  Input:  Candidate ASTs + Constraints           │
      │  Process: SMT constraint solving                │
      │  Output: Valid AST(s) satisfying all constraints│
      └─────────────────────┬───────────────────────────┘
                            ↓
      ┌─────────────────────────────────────────────────┐
      │  Semantic Analyzer (STANDARD - Prolog-based)    │
      │  Input:  Valid AST(s) + Prolog Rules            │
      │  Process: Tabled evaluation (fixed-point)       │
      │  Output: Semantic facts (symbol tables, CFG)    │
      │  Note:   Empty ruleset = no-op (zero cost)      │
      └─────────────────────┬───────────────────────────┘
                            ↓
                    [ Application Code ]
         (codegen, optimization, interpretation, etc.)
```

### 2.1 Layers (New Capabilities Highlighted)

1. **External Grammar Sources (All are INPUT formats)**
   - **HCL format** (`.grammar.hcl`): PRIMARY authoring format - human-friendly, clean syntax
   - **PDL format** (`.pdl`): LEGACY - Existing format in go-earley
   - **EBNF format** (`.ebnf`): IMPORT - Standard grammar notation
   - **YACC format** (`.y`): IMPORT - Parser generator input files
   - **JSON format**: EXPORT/IMPORT - Machine interchange (Grammar-IR can serialize to JSON)

2. **Format-Specific Converters (All produce Grammar-IR Go structs)**  
   - **HCL Parser**: Uses external `hashicorp/hcl/v2` library → parses `.grammar.hcl` → `grammarir.Grammar` structs
   - **PDL Importer**: Uses existing go-earley PDL parser → `grammarir.Grammar` structs (legacy compatibility)
   - **EBNF Importer**: Parses `.ebnf` → `grammarir.Grammar` structs (future)
   - **YACC Importer**: Parses `.y` → `grammarir.Grammar` structs (future)
   - **JSON marshaler**: `grammarir.Grammar` ↔ JSON (bidirectional for tooling)
   - **All formats converge to the same Go struct representation**
   - **No bootstrapping needed**: HCL library is external, mature, and already built

3. **Grammar-IR (Canonical Representation = Go Structs)**  
   - **Definition**: Grammar-IR IS the Go type system in `grammarir/` package
   - **Core Types**:
     - `type Grammar struct` - Complete grammar with start symbol
     - `type Symbol struct` - Terminal or NonTerminal
     - `type Production struct` - Rule mapping LHS → RHS sequence
     - `type Metadata map[string]any` - Extensible annotations
   - **Properties**:
     - **Format-agnostic**: Doesn't care if loaded from HCL, PDL, EBNF, or JSON
     - **Validated**: Well-formedness checked after construction (see grammar_ir_spec.md)
     - **Serializable**: Can export to JSON, HCL, or legacy PDL as needed
     - **Extensible**: Metadata carries constraints, AST hints, custom annotations
   - **Why HCL is "primary"**: Just easier for humans to write than JSON or EBNF
   
   - **Extensibility Pattern** (Critical Design Principle):
     - **Stable Core Fields**: Each struct has minimal, rarely-changing typed fields
       - Grammar: `StartSymbol`, `Symbols`, `Productions` (essentials only)
       - Production: `LHS`, `RHS` (structural minimum)
       - Symbol: `Name`, `Kind` (identity only)
     - **Metadata Maps**: Each struct includes `Metadata map[string]any`
       - Carries all extensions: constraints, AST hints, tooling annotations
       - HCL files populate metadata freely without Go code changes
       - Namespaced keys prevent collisions: `z3.constraint`, `ast.node_type`, `kanren.goal`
     - **Helper Accessors**: Convenience methods extract typed values from metadata
       - `func (p *Production) GetZ3Constraint() (string, bool)`
       - `func (p *Production) GetASTNodeType() (string, bool)`
       - `func (s *Symbol) GetPrecedence() (int, bool)`
       - Adding helpers does NOT change core structs or break serialization
     - **Philosophy**: 
       - Go structs change only for fundamental grammar theory updates
       - Feature additions live in metadata + helper functions
       - HCL carries the expressive weight, Go provides stable foundation
       - New tools (future constraint solver? new AST extractor?) add metadata namespaces, not struct fields

4. **Core Parser (go-earley) - EXISTING, Enhanced**  
   - **Current capability**: Builds SPPF from token stream
   - **Enhancement**: Accepts Grammar-IR (from HCL/PDL/etc.) instead of internal grammar format
   - **Retains**: Nullability checking, Leo/Aycock–Horspool optimizations
   - **Agnostic**: Works with any Grammar-IR source (HCL, PDL, EBNF, etc.)
   - **Output**: SPPF containing ALL valid parses (ambiguity preserved)

5. **AST Extractor (miniKanren) - NEW CAPABILITY**  
   - **Purpose**: Convert SPPF to concrete AST(s) - fills gap in existing library
   - **Current state**: Users must write custom visitors; no standard AST extraction
   - **Approach**: Logic programming (miniKanren) for flexible traversal
   - **Input**: SPPF forest + Grammar-IR metadata (AST structure hints)
   - **Process**: 
     - Generate Kanren goals from Grammar-IR productions
     - Explore SPPF paths via backtracking
     - Unify with AST structure constraints
   - **Output**: One or more candidate ASTs (ambiguity may still exist)
   - **Properties**: 
     - Handles ambiguity explicitly (can return multiple ASTs)
     - Compositional (small goals combine into larger goals)
     - Lazy (can stop after finding first valid AST if desired)

6. **Semantic Validator (Z3) - NEW CAPABILITY**  
   - **Purpose**: Filter ASTs by semantic constraints (type compatibility, etc.)
   - **Why after miniKanren**: Structural extraction first, then semantic validation
   - **Input**: Candidate AST(s) from miniKanren + constraint metadata from Grammar-IR
   - **Process**:
     - Read constraints from Production.Metadata (e.g., `z3_constraint`)
     - Translate to Z3 SMT-LIB format
     - Check satisfiability for each candidate AST
   - **Output**: Valid AST(s) that satisfy all constraints
   - **Properties**:
     - Prunes semantically invalid parses
     - Can reason about complex properties (type equality, arithmetic, etc.)
     - May reduce ambiguity (some parses fail constraints)

7. **Semantic Analyzer (Prolog) - STANDARD CAPABILITY**  
   - **Purpose**: Global semantic analysis via logic programming with tabled evaluation
   - **Library**: Uses `github.com/ichiban/prolog` (maintained, pure Go, ISO Prolog)
   - **Why standard (not optional)**: 
     - Zero cost when unused (empty ruleset = no-op)
     - Essential for real languages (scoping, control-flow, name resolution)
     - Uniform pipeline (no "with/without Prolog" modes)
     - Completes the semantic parsing framework
   - **Why after Z3**: Local constraints first, then global reachability/flow analysis
   - **Input**: Valid AST(s) from Z3 + Prolog rules (opt-in)
   - **Process**:
     - Convert AST to Prolog facts (nodes, edges, attributes)
     - Load analysis rules from:
       - Grammar-IR metadata (`metadata["prolog.rules"]`)
       - External `.pl` files (user-provided or standard library)
       - Empty ruleset → immediate pass-through (no analysis)
     - Execute queries with tabled evaluation (`:- table` directive for fixed-point)
   - **Output**: Semantic facts (symbol tables, CFG, type environment, warnings)
   - **Standard Library Rules** (shipped with go-earley):
     - `rules/name_resolution.pl` - Lexical scoping, symbol lookup
     - `rules/control_flow.pl` - CFG construction, reachability
     - `rules/data_flow.pl` - Def-use chains, liveness analysis
     - `rules/type_inference.pl` - Simple type propagation
     - `rules/lint.pl` - Common linting (unused vars, dead code)
   - **Use cases**:
     - Name binding and scope resolution
     - Control-flow and data-flow analysis
     - Type inference (complementing Z3's type checking)
     - Linting rules (unused variables, dead code)
     - Grammar-IR meta-queries (production usage, symbol references)
   - **Properties**:
     - Declarative logic programming
     - Tabled evaluation for fixed-point queries (transitive closure)
     - Backtracking for multiple solutions
     - Opt-in usage: users provide rules only when needed

8. **Downstream Application Logic**  
   - Receives valid AST(s) and semantic facts from Datalog
   - Traditional compiler passes: codegen, optimization, interpretation
   - Grammar-IR metadata can guide these passes (e.g., AST node types)
   - Semantic facts available even if Datalog rules were empty (just no derived facts)

---

## 3. Detailed Workflow

### 3.1 No Bootstrap Required

**Why no bootstrapping?**
- HCL parsing handled by external `hashicorp/hcl/v2` library (mature, maintained, already built)
- Grammar-IR is just Go structs - no code generation needed
- HCL → Grammar-IR mapping is straightforward Go code (read HCL AST, populate structs)
- PDL support uses existing go-earley PDL parser (no new parser needed)

**Implementation approach:**
1. Import `github.com/hashicorp/hcl/v2` and `github.com/hashicorp/hcl/v2/hclsimple`
2. Define HCL schema matching Grammar-IR structs
3. Write mapper: HCL attributes → Grammar-IR fields
4. For PDL: Existing parser → Grammar-IR structs (refactor existing code)

**No circular dependency**: We're using HCL as a data format, not defining HCL's grammar

### 3.2 Production Workflow (End-to-End)

#### Step 1: Grammar-IR Loading
- **Input**: HCL file (primary) or legacy format (PDL, EBNF, YACC)
- **Process**: Parse external format → construct Grammar-IR in memory
- **Output**: Validated Grammar-IR instance
- **Invariants**: All symbols referenced in productions must be declared; start symbol exists

#### Step 2: Table Generation & Parsing (EXISTING)
- **Input**: Grammar-IR instance + token stream
- **Process**: 
  - Compute nullable set, right-recursive productions, dotted rules
  - Execute Earley algorithm with Leo optimization
- **Output**: SPPF (Shared-Packed Parse Forest)
- **Properties**: 
  - SPPF preserves ALL valid parses (no disambiguation yet)
  - Common subtrees shared (space-efficient)
  - Ambiguity explicitly represented

#### Step 3: AST Extraction via miniKanren (NEW)
- **Why this step**: Existing library stops at SPPF; users must write custom extraction
- **Input**: 
  - SPPF from parser
  - Grammar-IR metadata (AST structure hints)
- **Process**: 
  - Generate miniKanren goals from productions
  - Traverse SPPF using logic programming
  - Unify with AST structural constraints
- **Output**: Candidate AST(s) - one per valid parse path
- **Properties**: 
  - Multiple ASTs if ambiguous
  - Backtracking explores alternatives
  - Can limit to first N solutions
- **Decision point**: Take all ASTs or first valid one?

#### Step 4: Semantic Validation via Z3 (NEW)
- **Why after miniKanren**: Structure first, semantics second
- **Input**: 
  - Candidate AST(s) from miniKanren
  - Constraint expressions from Grammar-IR metadata
- **Process**: 
  - For each candidate AST:
    - Extract relevant values (e.g., left/right operand types)
    - Translate constraint to Z3 SMT-LIB
    - Check satisfiability
  - Filter out unsatisfiable ASTs
- **Output**: Valid AST(s) satisfying all constraints
- **Properties**: 
  - May reduce ambiguity (some parses semantically invalid)
  - Can handle complex constraints (arithmetic, type systems, etc.)
  - Failure modes: no valid ASTs → parse error

#### Step 5: Semantic Analysis via Prolog (STANDARD)
- **Why after Z3**: Local constraints validated, now compute global facts
- **Input**: 
  - Valid AST(s) from Z3
  - Analysis rules from Grammar-IR metadata or external `.pl` files
- **Process**: 
  - Convert AST to Prolog facts (AST nodes → predicates)
  - Load rules:
    - Check `metadata["prolog.rules"]` in Grammar-IR
    - Load external files specified by user
    - If no rules provided → empty ruleset (pass-through)
  - Execute queries with tabled evaluation (`:- table` for fixed-point)
- **Output**: Semantic facts database
  - Symbol tables (name → declaration mapping)
  - Control-flow graphs (statement → successor edges)
  - Type environment (expression → inferred type)
  - Warnings (unused variables, dead code, etc.)
- **Properties**: 
  - Zero cost if no rules provided (no queries executed)
  - Standard library rules available for common tasks
  - Users can extend with custom Prolog predicates

#### Step 6: Application Logic
- **Input**: Valid AST(s) + semantic facts
- **Handling ambiguity**:
  - If single AST: proceed with compilation
  - If multiple ASTs: error (truly ambiguous) or user choice
- **Downstream passes**: Codegen, optimization, interpretation
  - Can query semantic facts: "what's the type of this expression?"
  - Can use symbol tables for name resolution
  - Can traverse CFG for optimization passes

#### Export & Interchange (Orthogonal)
- **Purpose**: Grammar-IR can round-trip through formats
- **Formats**: HCL ↔ JSON ↔ PDL (semantic preservation)
- **Use cases**: Tool interchange, version control, debugging

---

## 4. Extensibility Design Pattern

### 4.1 Core Principle: Stable Schema + Rich Metadata

Grammar-IR structs are designed for **long-term stability** while allowing **unlimited extensibility**:

- **The Problem**: Every new feature (constraint type, AST hint, tool annotation) could require adding struct fields
- **The Cost**: Each field addition breaks serialization, requires migration, forces recompilation
- **The Solution**: Minimal core fields + metadata maps for everything else

### 4.2 Struct Design Guidelines

#### Minimal Core Fields
Each struct includes ONLY fields that:
1. **Define identity**: What makes this thing what it is (e.g., Symbol.Name)
2. **Define structure**: Required for grammar well-formedness (e.g., Production.RHS)
3. **Used by all tools**: Parser, validators, all importers need it (e.g., Grammar.StartSymbol)

**Anti-pattern**: Adding `Z3Constraint string` directly to Production
- Breaks users who don't use Z3
- Requires migration when constraint format changes
- Clutters struct with tool-specific concerns

**Correct pattern**: Adding `Metadata map[string]any` once, forever
- Z3 uses `metadata["z3.constraint"]`
- Future tools use `metadata["tool.field"]`
- HCL populates freely: `metadata { z3.constraint = "..." }`

#### Metadata Namespace Conventions
```
<tool>.<category>.<field>
```
Examples:
- `z3.constraint.type_equality` - Z3 type constraint
- `z3.constraint.arithmetic` - Z3 arithmetic constraint
- `ast.node_type` - AST node hint for miniKanren
- `ast.node_name` - Custom AST node name
- `kanren.goal` - Custom Kanren goal expression
- `prolog.rules` - Prolog analysis rules (inline or file paths)
- `prolog.fact_type` - Prolog fact schema for AST node
- `prolog.table` - Predicates to table for fixed-point evaluation
- `precedence.level` - Operator precedence (future parser enhancement)
- `associativity.direction` - Left/right associativity
- `codegen.emit_as` - Codegen hint (future)
- `lint.suppress` - Linting suppression (future)

**Collision prevention**: Tool name prefix + dot notation creates natural namespaces

#### Helper Accessor Pattern
```go
// In production.go
func (p *Production) GetZ3Constraint() (string, bool) {
    if val, ok := p.Metadata["z3.constraint"]; ok {
        if str, ok := val.(string); ok {
            return str, true
        }
    }
    return "", false
}

func (p *Production) GetASTNodeType() (string, bool) {
    if val, ok := p.Metadata["ast.node_type"]; ok {
        if str, ok := val.(string); ok {
            return str, true
        }
    }
    return "", false
}

// In symbol.go
func (s *Symbol) GetPrecedence() (int, bool) {
    if val, ok := s.Metadata["precedence.level"]; ok {
        if i, ok := val.(int); ok {
            return i, true
        }
    }
    return 0, false
}
```

**Benefits**:
- Type safety at use site (callers get `string, bool` not `any`)
- Centralized key naming (avoids typos like `"z3_constraint"` vs `"z3.constraint"`)
- Default values easy to implement
- Documentation lives with code
- Adding helpers never breaks existing code

### 4.3 HCL Expressiveness

HCL metadata blocks carry arbitrary structured data:

```hcl
production "add_expr" {
  lhs = "expr"
  rhs = ["expr", "+", "expr"]
  
  metadata = {
    ast.node_type = "BinaryOp"
    ast.op_field   = "Add"
    
    z3.constraint.type_equality = "typeof(expr[0]) == typeof(expr[2])"
    z3.constraint.result_type   = "typeof(expr[0])"
    
    prolog.fact_type = "binary_op"
    prolog.rules = """
      :- table type_of/2.
      type_of(BinOp, T) :- binary_op(BinOp, _, Left, Right), type_of(Left, T), type_of(Right, T).
    """
    
    precedence.level = 10
    precedence.associativity = "left"
    
    codegen.emit_as = "IR_ADD"
    
    custom.my_tool.my_data = {
      nested = "structures"
      also   = "work"
    }
  }
}
```

**When Grammar-IR changes**:
- Adding `z3.constraint.new_feature` → HCL only (no Go changes)
- Adding new tool support → Add helper accessors (non-breaking)
- Changing constraint format → HCL migration (not Go struct migration)

**When Grammar-IR stays stable**:
- Fundamental grammar theory changes (e.g., adding "priority" to Symbol for parser disambiguation)
- Structural invariants change (e.g., allowing epsilon in RHS, requiring annotations)

### 4.4 Evolution Examples

#### Example 1: Adding Precedence Support
**Wrong approach**:
```go
type Symbol struct {
    Name       string
    Kind       SymbolKind
    Precedence int  // BAD: breaks all existing grammars
}
```

**Correct approach**:
```go
// No struct change!

// Add helper (non-breaking):
func (s *Symbol) GetPrecedence() (int, bool) {
    if val, ok := s.Metadata["precedence.level"]; ok {
        if i, ok := val.(int); ok {
            return i, true
        }
    }
    return 0, false  // default for symbols without precedence
}

// HCL users add metadata:
symbol "+" {
  kind = "terminal"
  metadata = {
    precedence.level = 10
  }
}
```

#### Example 2: Adding Custom AST Nodes
**Wrong approach**:
```go
type Production struct {
    LHS        Symbol
    RHS        []Symbol
    ASTNodeType string  // BAD: most productions don't need this
}
```

**Correct approach**:
```go
// No struct change!

// Add helper:
func (p *Production) GetASTNodeType() (string, bool) {
    if val, ok := p.Metadata["ast.node_type"]; ok {
        if str, ok := val.(string); ok {
            return str, true
        }
    }
    return "", false  // no custom node type
}

// HCL users opt-in:
production "func_def" {
  lhs = "function"
  rhs = ["def", "name", "(", "params", ")", "block"]
  
  metadata = {
    ast.node_type = "FunctionDef"
    ast.name_field = "params[1]"  // which RHS element is the name
  }
}
```

#### Example 3: Future Constraint Solver Integration
Suppose a new constraint solver "CVC5" is added in 2027:

**No Go struct changes needed**:
```hcl
production "safe_div" {
  lhs = "expr"
  rhs = ["expr", "/", "expr"]
  
  metadata = {
    # Z3 constraint (existing)
    z3.constraint = "expr[2] != 0"
    
    # New CVC5 constraint (no Go changes!)
    cvc5.constraint = "(assert (not (= divisor 0)))"
    cvc5.solver_options = {
      logic = "QF_NIA"
      timeout = 5000
    }
  }
}
```

**Only additions**:
- `cvc5/` package with solver integration
- Helper accessors: `GetCVC5Constraint()`, `GetCVC5Options()`
- Zero impact on existing code

### 4.5 Serialization & Round-Tripping

**Critical property**: Metadata must survive all format conversions

```
HCL → Grammar-IR → JSON → Grammar-IR → HCL
```

All metadata preserved because:
- `map[string]any` serializes naturally to JSON
- HCL metadata blocks accept arbitrary values
- PDL can embed JSON in comments if needed (legacy compatibility)

**Validation**: Metadata is opaque to core Grammar-IR
- Well-formedness checks (symbol references, start symbol exists) ignore metadata
- Tool-specific validation happens at use time (e.g., Z3 package validates `z3.constraint` syntax)
- Invalid metadata = tool failure, not grammar invalidity

---

## 5. Implementation Plan

### Phase 1: Grammar-IR Core (Weeks 1-2)
1. **Define stable core structs** in `grammarir/` package
   - Grammar: minimal fields (StartSymbol, Symbols, Productions, Metadata)
   - Symbol: minimal fields (Name, Kind, Metadata)
   - Production: minimal fields (LHS, RHS, Metadata)
   - Metadata type: `map[string]any`
   - **Design constraint**: No tool-specific fields in core structs
2. **Define HCL schema for Grammar-IR**
   - Grammar block with metadata support
   - Symbol blocks with metadata support
   - Production blocks with metadata support
   - **Validation**: Metadata blocks accept arbitrary keys (no schema restrictions)
3. **Implement HCL mapper** using `hashicorp/hcl/v2`
   - Import external HCL library: `go get github.com/hashicorp/hcl/v2`
   - Parse HCL files using `hclsimple.DecodeFile()`
   - Map HCL attributes → Grammar-IR struct fields
   - Preserve all metadata (no filtering)
4. **Implement JSON export/import** for interchange
   - Round-trip test: HCL → IR → JSON → IR → equality
   - Metadata preservation test

### Phase 2: PDL Importer Refactor (Week 3)
1. **Refactor existing PDL parser** to output `grammarir.Grammar`
   - Update `dsl.Parse()` to return `grammarir.Grammar` instead of internal format
   - Map PDL constructs → Grammar-IR structs
   - Preserve PDL-specific metadata (comments, annotations)
2. **Create HCL versions of example grammars**
   - Migrate calculator grammar: `calc.pdl` → `calc.grammar.hcl`
   - Migrate expression grammar: `expr.pdl` → `expr.grammar.hcl`
   - Keep PDL originals for compatibility testing
3. **Validate both input paths**
   - HCL → Grammar-IR → Parse pipeline
   - PDL → Grammar-IR → Parse pipeline
   - Ensure identical Grammar-IR output from equivalent grammars

### Phase 3: Parser Integration (Week 4)
1. Wire `go-earley` core to accept `grammarir.Grammar` instead of internal `dsl.Grammar`
2. Ensure nullability computation, Leo optimization work with Grammar-IR
3. Test with grammars loaded from both HCL and PDL

### Phase 4: Kanren & Z3 (Weeks 5-6)
1. **Add MiniKanren (`kanren/`)** for AST extraction
   - Implement goal generation from Production metadata
   - Read `metadata["ast.node_type"]`, `metadata["kanren.goal"]`
   - **No Grammar-IR struct changes** - only new helpers like `GetASTNodeType()`
2. **Integrate Z3 (`z3/`)** for constraint validation
   - Implement constraint extraction from Production metadata
   - Read `metadata["z3.constraint"]`, `metadata["z3.constraint.type_equality"]`, etc.
   - **No Grammar-IR struct changes** - only helpers like `GetZ3Constraint()`
3. **Document metadata conventions**
   - Create `docs/metadata_conventions.md` listing namespaces
   - `z3.*` namespace for Z3 constraints
   - `ast.*` namespace for AST hints
   - `kanren.*` namespace for Kanren goals
   - `prolog.*` namespace for Prolog rules/facts
4. **Update HCL examples** with metadata usage
   - Calculator grammar with `z3.constraint` for type checking
   - Expression grammar with `ast.node_type` for AST structure

### Phase 5: Prolog Integration (Week 7)
1. **Add Prolog engine (`prolog/`)** for semantic analysis
   - Import `github.com/ichiban/prolog` (maintained, pure Go)
   - Implement AST → Prolog facts conversion
   - Read `metadata["prolog.rules"]`, `metadata["prolog.fact_type"]`
   - **No Grammar-IR struct changes** - only helpers like `GetPrologRules()`
2. **Create standard library rulesets** in `rules/`
   - `rules/name_resolution.pl` - Lexical scoping, symbol tables
   - `rules/control_flow.pl` - CFG construction (with tabling)
   - `rules/data_flow.pl` - Def-use chains, liveness
   - `rules/type_inference.pl` - Type propagation
   - `rules/lint.pl` - Common linting rules
3. **Implement empty ruleset optimization**
   - Detect when no rules provided → skip evaluation (zero cost)
   - Pass-through mode: AST in, no queries executed
4. **Test with real semantic analyses**
   - Name resolution in calculator (variable scoping)
   - Control-flow in expression language with tabled predicates
   - Benchmark empty ruleset vs full analysis

### Phase 6: Tooling & Migration (Week 8)
1. **Update CLI** (`cmd/earley/`) with new flags
   - `--format=hcl|pdl|ebnf|json` - Specify input format
   - `--export-json` - Export Grammar-IR as JSON
   - `--prolog-rules=file.pl` - Load external Prolog rules
   - `--validate` - Grammar-IR well-formedness check only
2. **Create PDL → HCL migration tool**
   - Parse PDL → Grammar-IR → Export as HCL
   - Preserve metadata, comments as HCL comments
   - Automated batch migration for existing grammars
3. **Documentation updates**
   - Clarify: No bootstrapping needed (HCL is external library)
   - Document: HCL is primary, PDL is legacy compatibility
   - Migration guide: PDL users → HCL syntax

### Phase 7: Testing & Documentation (Week 9)
1. Extend tests: HCL parsing, importer compatibility, AST comparisons
2. Document HCL schema and usage examples
3. Create migration guide for PDL users
4. Benchmark HCL parsing vs. PDL parsing

---

## 5. Risks & Mitigations

- **HCL adoption barrier**  
  - Risk: Users unfamiliar with HCL syntax
  - Mitigation: Provide migration tool from PDL; comprehensive examples; HCL is similar to Terraform (widely known)

- **PDL transition friction**  
  - Risk: Existing PDL grammars need conversion
  - Mitigation: Keep PDL importer working; automated PDL→HCL converter; gradual migration path

- **Performance of HCL parsing**  
  - Risk: HCL parser might be slower than PDL
  - Mitigation: Parse once, cache Grammar-IR as JSON; HCL used at build time, not runtime

- **Performance overhead** of Kanren backtracking, SMT calls, and Prolog evaluation
  - Mitigate by caching Z3 contexts and pruning as early as possible
  - Limit Kanren goals to ambiguous/non-nullable spans only
  - Empty Prolog ruleset → zero-cost pass-through (optimization critical)
  - Benchmark: ensure trivial grammars see <1ms overhead for full pipeline

- **Complexity of additional importers**  
  - Start with HCL and PDL; add EBNF/YACC only when needed
  - Each importer maps to same Grammar-IR, so adding is incremental

- **Prolog library maturity**  
  - Choice: `github.com/ichiban/prolog` (maintained, pure Go, ISO Prolog)
  - Pros: Active development (2024 commits), no CGO, good test coverage
  - Cons: May not have all SWI-Prolog features (but sufficient for our needs)
  - Mitigation: Stick to ISO Prolog subset; document supported features

- **User onboarding**  
  - Ship HCL templates and guided tutorial in `docs/`
  - Offer default IR→Go codegen so users don't write boilerplate
  - Provide VSCode extension for `.grammar.hcl` syntax highlighting
  - Include standard Prolog library for common analyses (most users won't write rules)

---

## 6. Key Design Questions for TLA+ Modeling

Before implementation, these aspects need formal modeling to ensure correctness:

### 6.1 Grammar-IR Invariants
- **Well-formedness**: All RHS symbols must be declared; start symbol must exist
- **Consistency**: Symbol kinds (terminal/nonterminal) must be unambiguous
- **Completeness**: Nullable computation must terminate and be correct
- **Transformation preservation**: Importer conversions preserve language semantics
- **Metadata preservation** (NEW):
  - Round-trip property: `HCL → IR → JSON → IR → equality` preserves all metadata
  - Metadata independence: Grammar well-formedness independent of metadata content
  - Namespace isolation: Metadata key collisions cannot break grammar validity
  - Tool isolation: Invalid metadata for tool X does not affect tool Y

### 6.2 SPPF Construction (Enhanced existing behavior)
- **Correctness**: SPPF represents exactly the set of valid parse trees
- **Sharing**: Common subtrees are represented once (not duplicated)
- **Ambiguity handling**: Multiple parses coexist without interference
- **Traversal**: Any valid path through SPPF corresponds to a valid parse

### 6.3 miniKanren AST Extraction (NEW - Critical questions)
- **Soundness**: Every extracted AST corresponds to a valid SPPF path
- **Completeness**: Every valid SPPF path has at least one corresponding AST
- **Termination**: Kanren search must not loop infinitely (even with cycles in grammar)
- **Exploration strategy**: 
  - Depth-first vs breadth-first?
  - How to bound search space?
  - When to stop (first solution vs all solutions)?
- **Goal generation**: How to map Grammar-IR productions to Kanren goals?
- **Metadata interpretation**: How do AST hints in Grammar-IR guide extraction?

### 6.4 Z3 Constraint Integration (NEW - Critical questions)
- **Evaluation order**: miniKanren first, then Z3 (as designed) - why not interleave?
- **Constraint attachment**: Stored in Production.Metadata in Grammar-IR
- **Translation**: Grammar-IR constraint string → Z3 SMT-LIB - correctness?
- **Failure handling**: 
  - No valid ASTs → parse error (what error message?)
  - All ASTs fail constraints → which constraint violated?
- **Context passing**: How do child node values propagate to parent constraints?
- **Performance**: Check all ASTs or fail-fast?

### 6.5 Pipeline Composition
- **Question**: Why miniKanren before Z3, not interleaved?
  - **Answer candidate**: Structural extraction is grammar-local; semantics may be global
- **Question**: Can Z3 constraints prune SPPF before miniKanren?
  - **Trade-off**: Earlier pruning vs complexity of SPPF annotation
- **Question**: What if constraint checking is expensive?
  - **Strategy**: Lazy evaluation, memoization, incremental Z3?

### 6.6 Ambiguity Resolution
- **Multiple ASTs after Z3**: How to handle?
  - Return all (user chooses)?
  - Use heuristics (prefer simpler trees)?
  - Make it an error?
- **Constraint-based disambiguation**: Can constraints fully resolve ambiguity?
- **User control**: Grammar-IR metadata for disambiguation strategy?

### 6.7 Format Round-Tripping
- **Semantic preservation**: HCL → IR → JSON → IR → HCL preserves meaning
- **Metadata preservation**: All metadata keys/values survive round-trips exactly
- **Normalization**: When/if to normalize grammars (deduplication, sorting)
- **Versioning**: How grammar format versions interact
- **Extensibility invariants** (NEW - Critical for long-term stability):
  - Adding helper accessors never breaks serialization
  - New metadata namespaces never conflict with existing ones
  - Grammar-IR schema changes only for fundamental grammar theory updates
  - Tool-specific data lives in metadata, not core struct fields
  - **Question for TLA+**: Can we model "no struct changes needed for N new tools"?
  - **Question for TLA+**: What properties guarantee metadata extensibility without breaking changes?

### 6.8 Prolog Integration (STANDARD COMPONENT)
**Decision**: Prolog is a standard pipeline component with opt-in usage via rulesets

#### Rationale for Standard (Not Optional)
1. **Zero cost when unused**: Empty ruleset = immediate pass-through (no computation)
2. **Essential for real languages**: Any language with scoping, control-flow, or name resolution needs it
3. **Uniform pipeline**: No "with/without Prolog" modes to maintain
4. **Library completeness**: Completes semantic parsing framework (parse → AST → local constraints → global analysis)
5. **Negligible overhead**: Small grammars see <1ms; large grammars benefit significantly
6. **Standard library value**: Ship common analyses (name resolution, CFG, data flow) ready to use
7. **Mature library**: `ichiban/prolog` is actively maintained, pure Go, ISO Prolog compliant

#### Prolog's Strengths
- **Tabled evaluation**: Fixed-point computation via `:- table` directive (same semantics as Datalog)
- **Logic programming**: Declarative queries with backtracking and unification
- **Proven use cases**: 
  - Program analysis (points-to, def-use chains, call graphs)
  - Name resolution (scope chains, symbol tables)
  - Type inference (constraint generation and solving)
  - Semantic queries over ASTs
- **Advantage over Datalog**: More expressive (function symbols, negation-as-failure, DCGs)
- **ISO standard**: Widely known syntax, extensive literature

#### Standard Library Rulesets (Shipped with go-earley)
- `rules/name_resolution.pl` - Lexical scoping, symbol table construction
- `rules/control_flow.pl` - CFG construction, reachability analysis (tabled)
- `rules/data_flow.pl` - Def-use chains, live variable analysis (tabled)
- `rules/type_inference.pl` - Type propagation, constraint solving
- `rules/lint.pl` - Common linting (unused variables, dead code, unreachable statements)

#### Usage Pattern
```hcl
grammar "my_lang" {
  start = "program"
  
  metadata = {
    # Opt-in to standard analyses
    prolog.rules = ["name_resolution", "control_flow", "lint"]
    
    # Or provide custom rules
    prolog.custom_rules = file("my_analysis.pl")
  }
}
```

**Example Prolog rule (with tabling for fixed-point):**
```prolog
:- table reachable/2.

reachable(X, Y) :- edge(X, Y).
reachable(X, Z) :- reachable(X, Y), edge(Y, Z).

type_of(Expr, int) :- integer_literal(Expr).
type_of(BinOp, T) :- 
    binary_op(BinOp, '+', Left, Right),
    type_of(Left, T), 
    type_of(Right, T).
```

**No rules provided** → Prolog engine passes AST through (no queries executed)

#### Overlap Analysis: Prolog vs Existing Tools

| **Task**                     | **miniKanren** | **Z3**     | **Prolog** | **Verdict**              |
|------------------------------|----------------|------------|------------|--------------------------|
| AST extraction from SPPF     | ✅ Ideal       | ❌ No      | ⚠️ Possible | miniKanren wins (relational) |
| Constraint satisfaction      | ⚠️ Slow        | ✅ Ideal   | ❌ No      | Z3 wins (SMT solver)     |
| Transitive closure           | ⚠️ Manual      | ❌ No      | ✅ Ideal   | Prolog wins (tabled)     |
| Control-flow analysis        | ❌ Awkward     | ⚠️ Slow    | ✅ Ideal   | Prolog wins              |
| Name resolution              | ⚠️ Possible    | ❌ No      | ✅ Ideal   | Prolog wins              |
| Type inference               | ⚠️ Possible    | ✅ Good    | ✅ Good    | Prolog or Z3 (depends)   |
| Grammar meta-queries         | ❌ No          | ❌ No      | ✅ Ideal   | Prolog wins              |

#### Architecture Decision

**Pipeline position**: 
- **After Z3**: `Parser → SPPF → miniKanren → Z3 → Prolog → Application`
- **Scope**: Post-parse semantic analysis (name resolution, control-flow, data-flow, linting)

**Why this is NOT bloat**:
- miniKanren: Structure extraction (tree traversal with backtracking)
- Z3: Local constraint satisfaction (types, arithmetic, invariants)
- Prolog: Global reachability and fixed-point queries (scopes, flows, graphs)
- **Complementary strengths**: Each tool excels at different problem classes
- **Zero cost if unused**: Empty ruleset → no queries executed
- **High value when needed**: Standard library provides 80% of common semantic analyses

**Why Prolog over Datalog**:
- ✅ Mature Go library (`ichiban/prolog` - actively maintained, 2024 commits)
- ✅ Supports tabling (same fixed-point semantics as Datalog)
- ✅ More expressive when needed (function symbols, DCGs, cuts)
- ✅ ISO standard with extensive literature
- ❌ No mature Go Datalog libraries (abandoned or unmaintained)

**Implementation strategy**:
- Prolog engine always present (standard component via `ichiban/prolog`)
- Usage is opt-in (provide rules to enable analysis)
- Standard library makes it useful without custom rule writing
- Clear documentation: when to use Prolog vs Z3 vs miniKanren vs custom visitor

**Use case decision tree**:
- Need to query "all variables reachable from entry point"? → **Prolog** (tabled transitive closure)
- Need to check "this expression has type Int"? → **Z3** (local constraint)
- Need to extract "list of statements in function body"? → **miniKanren** (tree traversal)
- Need custom transformation logic? → **Traditional visitor pattern** (full control)

## 7. Next Steps

1. **Define HCL Schema**: Complete HCL schema specification for Grammar-IR structs
2. **Spike Grammar-IR Core**: Prototype `grammarir/` package with minimal structs
3. **Integrate HCL Library**: Import `hashicorp/hcl/v2`, write HCL → Grammar-IR mapper
4. **Refactor PDL Importer**: Update existing PDL parser to produce Grammar-IR
5. **PoC**: Load HCL grammar, generate parser tables, parse calculator expressions
6. **Benchmark**: Compare HCL vs PDL loading time; optimize if needed
7. **Review**: Iterate on HCL schema with early feedback
8. **Merge**: Roll out incremental PRs—Grammar-IR core, HCL mapper, PDL refactor, Kanren, Z3, Prolog

---

## 8. TLA+ Modeling Roadmap

### Goals
- Validate high-risk design decisions before implementation
- Prove core invariants hold across all importers and format conversions
- Define error handling for ambiguous/invalid cases
- Build confidence in novel integration points (miniKanren, Z3, Prolog)
- Document verified properties for future maintainers

### Scope (What to Model)

**In Scope:**
- ✅ Grammar-IR well-formedness and metadata preservation
- ✅ Integration with existing go-earley parser API
- ✅ Pipeline composition and error handling
- ✅ Extensibility pattern (metadata namespace isolation)

**Out of Scope:**
- ❌ External library behavior (HCL parsing, Z3 solving, Prolog evaluation)
- ❌ Performance characteristics (TLA+ is for correctness)
- ❌ Algorithmic correctness of Earley parser (already proven in literature)
- ❌ Complete miniKanren implementation (too large, implementation detail)

### Phase 1: Grammar-IR Core (Week 1, ~40 hours)

**Spec**: `irexpress_design/specs/GrammarIR.tla`

**Models**:
- Core types: `Symbol`, `Production`, `Grammar`
- Metadata as function: `Symbol -> [String -> Value]`
- Well-formedness predicate: `WellFormed(grammar)`
- Importer interface: `Import(externalFormat) -> Grammar`

**Properties to Prove**:
1. **Symbol Reference Validity**: All RHS symbols exist in `grammar.symbols`
   ```tla
   SymbolsValid(g) == 
     \A p \in g.productions: 
       \A s \in p.rhs: s \in g.symbols
   ```

2. **Start Symbol Exists**: Start symbol is declared and non-terminal
   ```tla
   StartSymbolValid(g) ==
     /\ g.start \in g.symbols
     /\ g.start \in NonTerminals
   ```

3. **Metadata Preservation**: Round-trip through formats preserves all metadata
   ```tla
   RoundTrip(g) ==
     LET json == ToJSON(g)
         g2 == FromJSON(json)
     IN g2.metadata = g.metadata
   ```

4. **Invalid Imports Fail**: Malformed external formats cannot create well-formed Grammar-IR
   ```tla
   InvalidImportsFail ==
     \A ext \in InvalidExternalFormats:
       ~WellFormed(Import(ext))
   ```

**Exit Criteria**: 
- TLC model-checks all properties without violations (state space < 10^8 states)
- All invariants pass with realistic grammar examples (calculator, expression)

### Phase 2: Integration Constraints (Week 2, ~40 hours)

**Spec**: `irexpress_design/specs/GrammarIRIntegration.tla`

**Models**:
- Existing `grammar.Grammar` API surface from go-earley
- Refinement mapping: `GrammarIRToExisting(grammarIR) -> grammar.Grammar`
- Computed properties: nullable set, right-recursive set, rule registry

**Properties to Prove**:
1. **Nullable Computation**: Nullable set is computable from productions alone
   ```tla
   NullableComputable(gir) ==
     LET nullable == ComputeNullable(gir.productions)
         existing == ExistingNullable(GrammarIRToExisting(gir))
     IN nullable = existing
   ```

2. **Right-Recursive Computation**: Can identify right-recursive productions
   ```tla
   RightRecursiveComputable(gir) ==
     LET rr == ComputeRightRecursive(gir.productions)
         existing == ExistingRightRecursive(GrammarIRToExisting(gir))
     IN rr = existing
   ```

3. **API Completeness**: Grammar-IR supports all existing parser operations
   ```tla
   SupportsExistingAPI(gir) ==
     /\ CanComputeRulesFor(gir)       \* RulesFor(NonTerminal)
     /\ CanComputeStartProductions(gir) \* StartProductions()
     /\ CanComputeNullable(gir)       \* IsTransitiveNullable(Symbol)
     /\ CanComputeRightRecursive(gir) \* IsRightRecursive(Production)
   ```

4. **Test Case Compatibility**: Existing test cases pass with Grammar-IR
   ```tla
   ExistingTestsPass ==
     /\ TransitiveNullTest     \* from grammar_test.go
     /\ RightRecursiveTest     \* from grammar_test.go
   ```

**Exit Criteria**: 
- Refinement mapping verified correct
- All existing test patterns provably supported

### Phase 3: Pipeline Composition (Week 3, ~40 hours)

**Spec**: `irexpress_design/specs/Pipeline.tla`

**Models**:
- Pipeline stages: `Parser -> SPPF -> miniKanren -> Z3 -> Prolog -> App`
- Data flow: `SPPF -> Set(AST) -> Set(ValidAST) -> Facts -> Result`
- Error states: empty sets at each stage
- Ambiguity: multiple ASTs persisting through stages

**Properties to Prove**:
1. **Extraction Completeness**: Non-empty SPPF produces at least one AST
   ```tla
   ExtractionCompleteness ==
     \A sppf \in ValidSPPF:
       sppf # {} => MiniKanren(sppf) # {}
   ```

2. **Constraint Filtering**: Z3 never creates ASTs, only filters
   ```tla
   Z3OnlyFilters ==
     \A asts \in Set(AST):
       Z3Validate(asts) \subseteq asts
   ```

3. **Error Reporting**: All failure modes have defined error messages
   ```tla
   ErrorsWellDefined ==
     /\ Z3Validate(asts) = {} => ReportConstraintViolation
     /\ MiniKanren(sppf) = {} => ReportExtractionFailure
     /\ ParseError => ReportSyntaxError
   ```

4. **Zero-Cost Empty Rules**: Empty Prolog rules don't execute queries
   ```tla
   EmptyRulesNoOp ==
     \A ast \in AST:
       PrologRules = {} => PrologAnalyze(ast, {}) = EmptyFacts
   ```

5. **Pipeline Never Stuck**: Every state transitions to output or error
   ```tla
   PipelineTerminates ==
     []<>(Output \/ Error)  \* Eventually reaches terminal state
   ```

**Exit Criteria**: 
- All error paths defined and tested
- Ambiguity resolution strategy proven correct
- Empty/no-op optimizations verified

### Decision Point (End of Week 3)

**If no issues found:**
- Proceed to implementation with confidence
- Use TLA+ specs as living documentation
- Reference proven properties in code comments

**If issues found:**
- Document design flaws discovered
- Iterate on design (estimate 1-2 weeks)
- Re-model affected areas
- Do NOT proceed to implementation until properties pass

### Phase 4: Extensibility Properties (Optional - Week 4, if needed)

**Spec**: `irexpress_design/specs/Extensibility.tla`

**Models**:
- Metadata namespace system
- Tool isolation via prefixes
- Schema stability over tool additions

**Properties to Prove**:
1. **Namespace Isolation**: Tool metadata doesn't affect other tools
   ```tla
   NamespaceIsolation ==
     \A t1, t2 \in Tools:
       t1 # t2 => metadata[t1.*] \cap metadata[t2.*] = {}
   ```

2. **N Tools Without Struct Changes**: Adding tools doesn't require Go changes
   ```tla
   NoStructChanges ==
     \A n \in 1..10:  \* Prove for 10 new tools
       AddNTools(n) => GrammarIRSchema = OriginalSchema
   ```

3. **Serialization Preserves Metadata**: All namespaces survive round-trips
   ```tla
   AllNamespacesPreserved ==
     \A ns \in Namespaces:
       RoundTrip(g).metadata[ns] = g.metadata[ns]
   ```

**Exit Criteria**: Mathematical proof that extensibility pattern scales

### Success Metrics

- **Coverage**: 3-4 TLA+ specs covering all high-risk areas
- **Properties**: 12-15 key properties formally verified
- **Bugs Found**: Design flaws caught before writing Go code (target: 2-5 issues)
- **Confidence**: Team comfortable proceeding to implementation
- **Documentation**: Specs serve as formal documentation for future maintainers

### Tools & Resources

- **TLA+ Toolbox**: Version 2.19 (already installed at `~/tla2tools.jar`)
- **VS Code Extension**: TLA+ (alygin.vscode-tlaplus) for syntax highlighting
- **Model Checker**: TLC with multi-threading enabled
- **Reference Book**: "Practical TLA+" by Hillel Wayne (available online)
- **Community**: TLA+ users group (groups.google.com/g/tlaplus) for questions
- **Examples**: TLA+ examples repository (github.com/tlaplus/Examples)

### Integration with Implementation Plan

TLA+ modeling runs **in parallel** with early implementation phases:

```
Weeks 1-2: Grammar-IR Core Implementation + TLA+ Phase 1
Week 3:    Parser Integration + TLA+ Phase 2
Week 4:    Continue Implementation + TLA+ Phase 3
Week 5:    Kanren/Z3/Prolog (only if TLA+ passed)
```

**Critical path**: TLA+ Phase 3 must complete before Week 5 implementation starts.

---

*Prepared by @gitrdm on 2025-10-22*  