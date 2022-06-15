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
	constStatement(tokens)
	registerStatement(tokens)
	// procedureStatement(tokens)
	// functionStatement(tokens)
	// theMain(tokens)
}

// ====================================== VAR ======================================

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

// ====================================== CONST ======================================

// <ConstStatement> ::= 'const' '{' <ConstList>
func constStatement(tokens []string) {
	if currentToken == "const" {
		nextToken(tokens) // {
		if currentToken == "{" {
			nextToken(tokens) // integer | string | real | boolean | char | identifier
			constList(tokens)
		}
	} else {
		// error
	}
}

// <ConstList>::= <ConstDeclaration> <ConstList>
func constList(tokens []string) {
	if currentToken == "}" {
		nextToken(tokens)
	} else {
		constDeclaration(tokens)
		constList(tokens)
	}
}

// <ConstDeclaration> ::= <ConstType> Identifier '=' <Value> <ConstDeclaration1>
func constDeclaration(tokens []string) {
	constType(tokens)
	nextToken(tokens) // Identifier

	// todo: Identifier

	nextToken(tokens) // ,
	constDeclaration1(tokens)
}

// <ConstDeclaration1> ::= ',' Identifier  '=' <Value> <ConstDeclaration1> | ';'
func constDeclaration1(tokens []string) {
	if currentToken == "," {
		nextToken(tokens)

		// todo: Identifier

		constDeclaration1(tokens)
	} else if currentToken == ";" {
		nextToken(tokens)
	} else {
		// error
	}
}

// <ConstType> ::= 'integer' | 'string' | 'real' | 'boolean' | 'char'
func constType(tokens []string) {
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
		// erro
	}
}

// ====================================== REGISTER ======================================

// <RegisterStatementMultiple> ::= <RegisterStatement> |
func registerStatementMultiple(tokens []string) {
	registerStatement(tokens)
}

// <RegisterStatement> ::= 'register' Identifier '{' <RegisterList>
func registerStatement(tokens []string) {
	if currentToken == "register" {
		nextToken(tokens)
		// todo Identifier
		nextToken(tokens)
		if currentToken == "{" {
			nextToken(tokens)
			registerList(tokens)
		}
	}
}

// <RegisterList> ::= <RegisterDeclaration> <RegisterList1>
func registerList(tokens []string) {
	registerDeclaration(tokens)
	registerList1(tokens)
}

// <RegisterList1> ::= <RegisterDeclaration> <RegisterList1> | '}' <RegisterStatementMultiple>
func registerList1(tokens []string) {
	if currentToken == "}" {
		nextToken(tokens)
		registerStatementMultiple(tokens)
	} else {
		registerDeclaration(tokens)
		registerList1(tokens)
	}
}

// <RegisterDeclaration> ::= <ConstType> Identifier <RegisterDeclaration1>
func registerDeclaration(tokens []string) {
	constType(tokens)
	nextToken(tokens)
	// todo Identifier
	registerDeclaration1(tokens)
}

// <RegisterDeclaration1> ::= ',' Identifier <RegisterDeclaration1> | ';'
func registerDeclaration1(tokens []string) {

}

// ====================================== PROCEDURE ======================================

// <Value>  ::= Decimal | RealNumber | StringLiteral | Identifier <ValueRegister> | Char | Boolean
func value(tokens []string) {
	// todo: values
}
