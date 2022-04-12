package Compiler

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tokenError      tokenType = iota // whether an error has ocurred.
	tokenKeyword                     // No.1
	tokenLetter                      // letter
	tokenIdentifier                  // identifier
	tokenDigit                       // single digit
	tokenNumber                      // number
	tokenArithmeticOp
	tokenBlockComment
	tokenEOF
	tokenRelationalOp
	tokenLogicalOp
	tokenDelimiter
	tokenChar
	tokenMalformedChar

	tokenIf              // if keyword
	tokenLeftMeta        // left meta-string
	tokenPipe            // pipe symbol
	tokenRange           // range keyword
	tokenRawString       // raw quoted string (includes quotes)
	tokenRightMeta       // right meta-string
	tokenString          // quoted string (includes quotes)
	tokenText            // plain text
	tokenMalformedNumber // Error when lexing number
	tokenMalformedComment
	// Keywords
	programKeyword   = "program"
	varKeyword       = "var"
	constKeyword     = "const"
	registerKeyword  = "register"
	functionKeyword  = "function"
	procedureKeyword = "procedure"
	returnKeyword    = "return"
	mainKeyword      = "main"
	ifKeyword        = "if"
	elseKeyword      = "else"
	whileKeyword     = "while"
	readKeyword      = "read"
	writeKeyword     = "write"
	integerKeyword   = "integer"
	realKeyword      = "real"
	booleanKeyword   = "boolean"
	charKeyword      = "char"
	stringKeyword    = "string"
	trueKeyword      = "true"
	falseKeyword     = "false"
)

// token Defines a Token (token) structure.
type token struct {
	typ tokenType
	val string
}

type tokenType int

// lexer Holds the state of the scanner.
// start is where the next token sent out begins.
// pos is where we are in the scanning.
type lexer struct {
	name   string     // used for error reports
	input  string     // string being scanned
	start  int        // start position of this token
	pos    int        // current position in the input
	width  int        // width of last rune read
	tokens chan token // channel of the scanned items
}

// stateFn Represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// Lex a constructor.
func Lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
	}
	go l.run() // Concurrently runs the state machine.
	//time.Sleep(time.Second)
	return l
}

// run Lexes the input by executing state functions
// until the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens) // No more tokens will be delivered.
}

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.isEOF() {
		l.width = 0
		return rune(tokenEOF)
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// lexText A partir de um estado inicial procura por lexemas tratando-os individualmente.
func lexText(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsLetter(r):
			return lexLetter
		case unicode.IsSpace(r):
			l.ignore()
		case strings.IndexRune("+-", r) >= 0, unicode.IsNumber(r):
			if strings.IndexRune("+", r) >= 0 || strings.IndexRune("-", r) >= 0 {
				if unicode.IsNumber(l.peek()) {
					return lexNumber
				}
			}
			if unicode.IsNumber(r) {
				return lexNumber
			}
			l.backup()
			return lexArithmeticOperator
		case strings.IndexRune("&|!", r) >= 0:
			if strings.IndexRune("=", l.peek()) >= 0 {
				return lexRelationalOperator
			}
			l.backup()
			return lexLogigcalOperator
		case strings.IndexRune("=!><", r) >= 0:
			if strings.IndexRune("=", l.peek()) >= 0 {
				return lexRelationalOperator
			}
			l.backup()
			return lexRelationalOperator
		case strings.IndexRune(";,(){}[].:", r) >= 0:
			l.backup()
			return lexDelimiter

		// TODO: Working for rights cases. Not working for: 'a''dfds'
		case strings.IndexRune("'", r) >= 0:
			if unicode.IsLetter(l.peek()) {
				if strings.IndexRune("'", l.next()) >= 0 {
					return lexChar
				}
			} else {
				l.emit(tokenMalformedChar)
				return lexText
			}
			return lexChar
		case strings.IndexRune("/", r) >= 0 || strings.IndexRune("*", r) >= 0:
			if strings.IndexRune("/", r) >= 0 { // If (r == / or *) check whether it could possibly be a comment block.
				if !(strings.IndexRune("#", l.peek()) >= 0) { // If not, emit an arithmetic operator (/)
					l.emit(tokenArithmeticOp)
				} else { // Otherwise, lex a comment block.
					return lexCommentBlock
				}
			}
			if strings.IndexRune("*", r) >= 0 { // Verification not necessary but left intentionally for legibility.
				l.emit(tokenArithmeticOp)
			}
		case r == rune(tokenEOF):
			l.emit(tokenEOF)
			return nil
		}
	}
}

// lexLetter If next rune is a whitespace, then emit a letter (token) and go back to initial state.
// Otherwise it could be a keyword or identifier, then handle it as so.
func lexLetter(l *lexer) stateFn {
	switch r := l.next(); {
	case r == rune(tokenEOF):
		l.emit(tokenLetter)
		l.emit(tokenEOF)
		return nil
	case l.isIdentifier(r):
		switch nextRune := l.peek(); {
		case l.isIdentifier(nextRune):
			return lexLetter
		default:
			return lexIdentifier
		}
	default:
		l.backup()
		l.emit(tokenLetter)
		return lexText
	}
}

// lexIdentifier At this point, the word is already formed,
// thus we need to verify whether is either an identifier or keyword.
func lexIdentifier(l *lexer) stateFn {
	if !strings.Contains("_", l.input[l.start:l.pos]) { // Verify whether is a keyword, only if it does not contain underline character
		if !l.emitIfKeyword() { // If not a keyword, emit an identifier then go back to initial state.
			l.emit(tokenIdentifier)
			return lexText
		}
		return lexText
	}
	l.emit(tokenIdentifier)
	return lexText
}

// lexNumber lexes a signed number (digit or multiple digits) including floating point.
func lexNumber(l *lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") { // todo: check error possibilities to implement here
		l.acceptRun(digits)
	}
	if len(l.input[l.start:l.pos]) > 1 {
		l.emit(tokenNumber)
		return lexText
	}
	l.emit(tokenDigit)
	return lexText
}

func lexArithmeticOperator(l *lexer) stateFn {
	switch r := l.next(); {
	case strings.IndexRune("+", r) >= 0 && strings.IndexRune("+", l.peek()) >= 0:
		l.next()
		l.emit(tokenArithmeticOp)
	case strings.IndexRune("-", r) >= 0 && strings.IndexRune("-", l.peek()) >= 0:
		l.next()
		l.emit(tokenArithmeticOp)
	default:
		l.emit(tokenArithmeticOp)
	}
	return lexText
}

// lexCommentBlock Lexes a comment block & identify possible errors.
func lexCommentBlock(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case strings.IndexRune("#", r) >= 0:
			if strings.IndexRune("/", l.peek()) >= 0 {
				l.next()
				l.emit(tokenBlockComment)
				return lexText
			}
		case r == rune(tokenEOF):
			l.emit(tokenMalformedComment)
			l.emit(tokenEOF)
			return nil
		}
	}
}

func lexRelationalOperator(l *lexer) stateFn {
	l.next()
	l.emit(tokenRelationalOp)
	return lexText
}

func lexDelimiter(l *lexer) stateFn {
	l.next()
	l.emit(tokenDelimiter)
	return lexText
}

func lexChar(l *lexer) stateFn {
	l.next()
	l.emit(tokenChar)
	return lexText
}

func lexLogigcalOperator(l *lexer) stateFn {
	switch r := l.next(); {
	case strings.IndexRune("&", r) >= 0 && strings.IndexRune("&", l.peek()) >= 0:
		l.next()
		l.emit(tokenLogicalOp)
	case strings.IndexRune("|", r) >= 0 && strings.IndexRune("|", l.peek()) >= 0:
		l.next()
		l.emit(tokenLogicalOp)
	default:
		l.emit(tokenLogicalOp)
	}
	return lexText
}

// isIdentifier return whether it could be identifier or not.
func (l *lexer) isIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || strings.IndexRune("_", r) >= 0
}

func (i token) String() string {
	switch i.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

/* emitIfKeyword Identify whether is a keyword if so, emit.
We could have worked with an array instead of a switch case,
although, intentionally using a switch may be useful for distinction
(if needed) in the future.
*/
func (l *lexer) emitIfKeyword() bool {
	switch word := l.input[l.start:l.pos]; {
	case word == programKeyword:
		l.emit(tokenKeyword)
		return true
	case word == varKeyword:
		l.emit(tokenKeyword)
		return true
	case word == constKeyword:
		l.emit(tokenKeyword)
		return true
	case word == registerKeyword:
		l.emit(tokenKeyword)
		return true
	case word == functionKeyword:
		l.emit(tokenKeyword)
		return true
	case word == procedureKeyword:
		l.emit(tokenKeyword)
		return true
	case word == returnKeyword:
		l.emit(tokenKeyword)
		return true
	case word == mainKeyword:
		l.emit(tokenKeyword)
		return true
	case word == ifKeyword:
		l.emit(tokenKeyword)
		return true
	case word == elseKeyword:
		l.emit(tokenKeyword)
		return true
	case word == whileKeyword:
		l.emit(tokenKeyword)
		return true
	case word == readKeyword:
		l.emit(tokenKeyword)
		return true
	case word == writeKeyword:
		l.emit(tokenKeyword)
		return true
	case word == integerKeyword:
		l.emit(tokenKeyword)
		return true
	case word == realKeyword:
		l.emit(tokenKeyword)
		return true
	case word == booleanKeyword:
		l.emit(tokenKeyword)
		return true
	case word == charKeyword:
		l.emit(tokenKeyword)
		return true
	case word == stringKeyword:
		l.emit(tokenKeyword)
		return true
	case word == trueKeyword:
		l.emit(tokenKeyword)
		return true
	case word == falseKeyword:
		l.emit(tokenKeyword)
		return true
	default:
		return false
	}
}

/*
	Helper Functions Below:
		ignore Skips over the pending input before this point.
		backup Steps back one rune.
		peek   Returns but not consume the next rune in the input.
	Aceptors:
		accept Consumes the next rune if it is from the valid set.
		acceptRun Consumes a run of runes from the valid set.
*/
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup Can only be called once per call of next !!
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) isEOF() bool {
	if l.pos >= len(l.input) {
		return true
	}
	return false
}

// Display tokens.
func (l *lexer) Debug() {
	fmt.Println("Lexer state info: ")
	fmt.Println("Start index: ")
	fmt.Println(l.start)
	fmt.Println("Current index: ")
	fmt.Println(l.pos)
	fmt.Println("Last token width: ")
	fmt.Println(l.width)
	fmt.Println("token list: ")
	for i := range l.tokens {
		fmt.Println("value: ", i.val)
		fmt.Println("type index: ", i.typ)
		fmt.Println("type name: ", [...]string{"tokenError",
			"tokenKeyword",
			"tokenLetter",
			"tokenIdentifier",
			"tokenDigit",
			"tokenNumber",
			"tokenArithmeticOp",
			"tokenBlockComment",
			"tokenEOF",
			"tokenRelationalOp",
			"tokenLogicalOp",
			"tokenDelimiter",
			"tokenChar",
			"tokenMalformedChar",

			"tokenIf",
			"tokenLeftMeta",
			"tokenPipe",
			"tokenRange",
			"tokenRawString",
			"tokenRightMeta",
			"tokenString",
			"tokenText",
			"tokenMalformedNumber",
			"tokenMalformedComment",

			"programKeyword",
			"varKeyword",
			"constKeyword",
			"registerKeyword",
			"functionKeyword",
			"procedureKeyword",
			"returnKeyword",
			"mainKeyword",
			"ifKeyword",
			"elseKeyword",
			"whileKeyword",
			"readKeyword",
			"writeKeyword",
			"integerKeyword",
			"realKeyword",
			"booleanKeyword",
			"charKeyword",
			"stringKeyword",
			"trueKeyword",
			"falseKeyword"}[i.typ])
		fmt.Println("length: ", len(string((i.val))))
		fmt.Println("--------------------------------")
	}
}
