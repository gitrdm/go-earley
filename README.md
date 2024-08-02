# go-earley 

A go earley parsing library

# Install

```bash
go get -u github.com/patrickhuber/go-earley
```

# Usage

## Create a grammar

> calculator.pdl

```
Calculator 
	= Expression;
		
Expression 
	= Expression '+' Term
	| Term;
		
Term 
	= Term '*' Factor
	| Factor;
		
Factor 
	= Number ;
	
Number 
	= Digits;
		
Digits ~ /[0-9]+/ ;
Whitespace ~ /[\s]+/ ;
	
:start = Calculator;
:ignore = Whitespace;
```

## Create a Grammar Instance

```golang
grammarString, err := os.ReadFile("calculator.pdl")
if err != nil{
    fmt.Fatal(err)
    os.Exit(1)
}
g, err := dsl.Parse(grammarString)
if err != niL{
    fmt.Fatal(err)
    os.Exit(1)
}
```

## Parse Some Expressions

```golang
```