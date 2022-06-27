package Compiler

import (
	"fmt"
	"os"
)

var tokens = make([]token, 0)
var tokenIndex = -1

// https://www.youtube.com/watch?v=tfIQzjUMKXA - 25:46

// Constructor
func Syntax(l *lexer) {
	// todo: remove it later
	tempProgramToken := token{27, "\"program\"", 0}
	tokens = append(tokens, tempProgramToken)
	// ================================================

	for t := range l.tokens {
		tokens = append(tokens, t)
	}

	fmt.Println(tokens)
	fmt.Println("")

	start()
}

func nextToken() {
	if tokenIndex < len(tokens) {
		tokenIndex++
	}
}

func lookAhead(index int) token {
	if len(tokens) == 0 {
		tempProgramToken := token{8, "EOF", 0}
		return tempProgramToken
	} else if tokenIndex+index >= len(tokens) {
		return tokens[len(tokens)-1]
	} else if tokenIndex+index < 0 {
		return tokens[0]
	}
	return tokens[tokenIndex+index]
}

func compareType(t token, tp2 string) bool {
	return t.typ == parseTokenTypeByString(tp2)
}

func verifyToken(tokenType string, tokenValue string, next bool) bool {
	if compareType(lookAhead(1), tokenType) && (tokenValue == "{any token value here}" || lookAhead(1).val == tokenValue || lookAhead(1).typ == 8) {
		if next {
			nextToken()
		}
		return true
	}
	return false
}

func verifyDelimiter(tokenValue string, next bool) bool {
	return verifyToken("tokenDelimiter", tokenValue, next)
}

func match(tokenTpy string, lookAheadSync int) bool {
	if compareType(lookAhead(1), tokenTpy) {
		fmt.Println("Match", lookAhead(1))
		nextToken()
		return true
	} else {
		syntaxError(tokenTpy, lookAheadSync)
		return false
	}
}

func syntaxError(tokenTpy string, lookAheadSync int) {
	// todo: implement a better way or not
	fmt.Println("Error na linha", lookAhead(0).line+1, "Esperando", lookAhead(0).val, "(", lookAhead(0).typ, ")", "porém foi recebido", "\""+lookAhead(1).val+"\"", "(", lookAhead(1).typ, ")")
	for i := 0; i < lookAheadSync; i++ {
		nextToken()
	}
}

func debug(die bool, message ...string) {
	fmt.Println()
	fmt.Println("======================= [[ DEBUG ]] =======================")
	fmt.Println()
	if message != nil {
		fmt.Println("Message", message)
	}
	fmt.Println("Previus token is  ", "\""+lookAhead(-1).val+"\"", "=", types[lookAhead(1).typ], "(", lookAhead(1).typ, ")")
	fmt.Println("Current token is  ", "\""+tokens[tokenIndex].val+"\"", "=", types[tokens[tokenIndex].typ], "(", tokens[tokenIndex].typ, ")")
	fmt.Println("Next token is     ", "\""+lookAhead(1).val+"\"", "=", types[lookAhead(1).typ], "(", lookAhead(1).typ, ")")
	fmt.Println()
	fmt.Println("===================== [[ END DEBUG ]] =====================")
	fmt.Println()
	if die {
		os.Exit(0)
	}
}

// <Start> ::= 'program' Identifier ';' <GlobalStatement>
func start() {
	match("programKeyword", 3)
	match("tokenIdentifier", 2)
	match("tokenDelimiter", 1)
	globalStatement()
}

// <GlobalStatement> ::= <VarStatement> <ConstStatement> <RegisterStatement><ProcedureStatement><FunctionStatement> <Main>
func globalStatement() {
	varStatement()
	constStatement()
	registerStatement()
	// procedureStatement()
	// functionStatement()
	// theMain()
}

// ====================================== VAR ======================================

// <VarStatement>::= 'var' '{' <VarList>
func varStatement() {
	match("tokenKeyword", 2)
	match("tokenDelimiter", 1)
	varList()
}

// <VarList>::= <VarDeclaration> <VarList> | '}'
func varList() {
	if verifyDelimiter("}", true) { // }
		return
	} else {
		varDeclaration()
		varList()
	}
}

// <VarDeclaration>::= <VarType> Identifier <VarDeclaration1>
func varDeclaration() {
	varType()
	match("tokenIdentifier", 1)
	varDeclaration1()
}

// <VarDeclaration1>::= ',' Identifier <VarDeclaration1> | ';'
func varDeclaration1() {
	if verifyDelimiter(";", true) { // ;
		return
	} else {
		match("tokenDelimiter", 2)
		match("tokenIdentifier", 1)
		varDeclaration1()
	}
}

// <VarType>::= 'integer' | 'string' | 'real' | 'boolean' | 'char' | Identifier
func varType() {
	if compareType(lookAhead(1), "tokenKeyword") {
		match("tokenKeyword", 1)
	} else if compareType(lookAhead(1), "tokenIdentifier") {
		match("tokenIdentifier", 1)
	}
}

// ====================================== CONST ======================================

// <ConstStatement> ::= 'const' '{' <ConstList>
func constStatement() {
	match("tokenKeyword", 2)
	match("tokenDelimiter", 1)
	constList()
}

// <ConstList>::= <ConstDeclaration> <ConstList> | '}'
func constList() {
	if verifyDelimiter("}", true) { // }
		return
	} else {
		constDeclaration()
		constList()
	}
}

// <ConstDeclaration> ::= <ConstType> Identifier '=' <Value> <ConstDeclaration1>
func constDeclaration() {
	varType()
	match("tokenIdentifier", 1)
	match("tokenArithmeticOp", 1)

	//match("tokenNumber", 1) // value
	value()

	constDeclaration1()
}

// <ConstDeclaration1> ::= ',' Identifier  '=' <Value> <ConstDeclaration1> | ';'
func constDeclaration1() {
	if verifyDelimiter(";", true) { // ;
		return
	} else {
		match("tokenDelimiter", 2)
		match("tokenIdentifier", 1)
		constDeclaration1()
	}
}

// <Value>  ::= Decimal | RealNumber | StringLiteral | Identifier <ValueRegister> | Char | Boolean
func value() {
	if verifyToken("tokenNumber", "{any token value here}", true) {
		return
	} else if verifyToken("tokenString", "{any token value here}", true) {
		return
	} else if verifyToken("tokenIdentifier", "{any token value here}", true) {
		valueRegister()
	} else if verifyToken("tokenChar", "{any token value here}", true) {
		return
	} else if verifyToken("tokenKeyword", "{any token value here}", true) {
		return
	}
}

// <ValueRegister> ::= '.' Identifier |
func valueRegister() {
	if verifyDelimiter(".", false) {
		match("tokenDelimiter", 2)
		match("tokenIdentifier", 1)
	}
}

// ====================================== REGISTER ======================================

// <RegisterStatementMultiple> ::= <RegisterStatement> |
func registerStatementMultiple() {
	if lookAhead(1).val == "register" {
		registerStatement()
	}
}

// <RegisterStatement> ::= 'register' Identifier '{' <RegisterList>
func registerStatement() {
	if lookAhead(1).val == "register" && match("tokenKeyword", 1) {
		match("tokenIdentifier", 2)
		match("tokenDelimiter", 1)
		registerList()
	}
}

// <RegisterList> ::= <RegisterDeclaration> <RegisterList1>  | '}'
func registerList() {
	if compareType(lookAhead(1), "tokenDelimiter") && lookAhead(1).val == "}" { // }
		nextToken()
	} else {
		registerDeclaration()
		registerList1()
	}
}

// <RegisterList1> ::= <RegisterDeclaration> <RegisterList1> | '}' <RegisterStatementMultiple>
func registerList1() {
	if verifyDelimiter("}", false) { // ;
		registerStatementMultiple()
	} else {
		registerDeclaration()
		registerList1()
	}
}

// <RegisterDeclaration> ::= <ConstType> Identifier <RegisterDeclaration1>
func registerDeclaration() {
	varType()
	match("tokenIdentifier", 1)
	registerDeclaration1()
}

// <RegisterDeclaration1> ::= ',' Identifier <RegisterDeclaration1> | ';'
func registerDeclaration1() {
	if verifyDelimiter(";", true) { // ;
		return
	} else {
		match("tokenDelimiter", 2)
		match("tokenIdentifier", 1)
		registerDeclaration1()
	}
}

// ====================================== PROCEDURE ======================================
