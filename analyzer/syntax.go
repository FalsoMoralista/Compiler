package Compiler

import "fmt"

var tokenIndex = 0
var currentToken = ""

// Constructor
func Syntax() {
	tokens := []string{"program", "teste", ";"}
	start(tokens)
}

func nextToken(tokens []string) {
	if tokenIndex < len(tokens) {
		currentToken = tokens[tokenIndex]
		tokenIndex++
	}
}

// <Start> ::= 'program' Identifier ';' <GlobalStatement>
func start(tokens []string) {
	nextToken(tokens)
	if currentToken == "program" {
		nextToken(tokens) // identifies
		nextToken(tokens)

		if currentToken == ";" {
			globalStatement(tokens)
		}
	} else {
		// error
	}
}

// <GlobalStatement> ::= <VarStatement> <ConstStatement> <RegisterStatement><ProcedureStatement><FunctionStatement> <Main>
func globalStatement(tokens []string) {
	nextToken(tokens)
	varStatement(tokens)
	// constStatement(tokens)
	// registerStatement(tokens)
	// procedureStatement(tokens)
	// functionStatement(tokens)
	// theMain(tokens)
}

// <VarStatement>::= 'var' '{' <VarList>
func varStatement(tokens []string) {
	if currentToken == "var" {
		nextToken(tokens)
		if currentToken == "{" {
			fmt.Print("Var")
			// varList()
		}
	}
}
