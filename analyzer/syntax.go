package Compiler

import "fmt"

var tokenIndex = 0
var currentToken = ""

// Constructor
func Syntax() {
	tokens := []string{"program", "teste", ";", "var", "teste", "{", "TDC"}
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
	nextToken(tokens) // program
	if currentToken == "program" {
		nextToken(tokens) // identifies
		nextToken(tokens) // ;
		if currentToken == ";" {
			nextToken(tokens) // var
			globalStatement(tokens)
		}
	} else {
		// error
	}
}

// <GlobalStatement> ::= <VarStatement> <ConstStatement> <RegisterStatement><ProcedureStatement><FunctionStatement> <Main>
func globalStatement(tokens []string) {
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
		nextToken(tokens) // {
		if currentToken == "{" {
			nextToken(tokens) // integer | string | real | boolean | char | identifier
			varList(tokens)
		}
	} else {
		// error
	}
}

// <VarList>::= <VarDeclaration> <VarList> | '}'
func varList(tokens []string) {
	if currentToken == "}" {
		nextToken(tokens)
	} else {
		varDeclaration(tokens)
		varList(tokens)
	}
}

// <VarDeclaration>::= <VarType> Identifier <VarDeclaration1>
func varDeclaration(tokens []string) {
	varType(tokens)
	nextToken(tokens) // Identifier

	// todo: Identifier

	nextToken(tokens) // ,
	varDeclaration1(tokens)
}

// <VarDeclaration1>::= ',' Identifier <VarDeclaration1> | ';'
func varDeclaration1(tokens []string) {
	if currentToken == "," {
		nextToken(tokens)

		// todo: Identifier

		varDeclaration1(tokens)
	} else if currentToken == ";" {
		nextToken(tokens)
	} else {
		// error
	}
}

// <VarType>::= 'integer' | 'string' | 'real' | 'boolean' | 'char' | Identifier
func varType(tokens []string) {
	if currentToken == "integer" {
		fmt.Print("integer")
		// todo
	} else if currentToken == "string" {
		fmt.Print("string")
		// todo
	} else if currentToken == "real" {
		fmt.Print("real")
		// todo
	} else if currentToken == "boolean" {
		fmt.Print("boolean")
		// todo
	} else if currentToken == "char" {
		fmt.Print("char")
		// todo
	} else {
		// todo Identifier
	}
}
