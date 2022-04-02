package lexical_analysis

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	itemError      itemType = iota // whether an error has ocurred.
	itemKeyword                    // No.1
	itemLetter                     // letter
	itemIdentifier                 // identifier
	itemDigit                      // single digit
	itemNumber                     // number

	// itemDot
	itemEOF
	itemElse  // else keyword
	itemEnd   // end keyword
	itemField // identifier, starting with '.'

	itemIf              // if keyword
	itemLeftMeta        // left meta-string
	itemPipe            // pipe symbol
	itemRange           // range keyword
	itemRawString       // raw quoted string (includes quotes)
	itemRightMeta       // right meta-string
	itemString          // quoted string (includes quotes)
	itemText            // plain text
	itemMalformedNumber // Error when lexing number
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

// item Defines a Token (item) structure.
type item struct {
	typ itemType
	val string
}

type itemType int

// lexer Holds the state of the scanner.
// start is where the next token sent out begins.
// pos is where we are in the scanning.
type lexer struct {
	name  string    // used for error reports
	input string    // string being scanned
	start int       // start position of this item
	pos   int       // current position in the input
	width int       // width of last rune read
	items chan item // channel of the scanned items
}

// stateFn Represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// Lex a constructor.
func Lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
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
	close(l.items) // No more tokens will be delivered.
}

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.isEOF() {
		l.width = 0
		return rune(itemEOF)
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
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
		case strings.IndexRune("+-", r) >= 0, unicode.IsNumber(r): // todo: verify whether is plus or minus, based on that do the following possible verifications.
			return lexNumber
		case r == rune(itemEOF):
			l.emit(itemEOF)
			return nil
		}
	}
}

// lexLetter
// If next rune is a whitespace, then emit a letter (token) and go back to initial state.
// Otherwise it could be other stuff (check-it in lexInsideAction).
func lexLetter(l *lexer) stateFn {
	switch r := l.next(); {
	case r == rune(itemEOF):
		l.emit(itemLetter)
		l.emit(itemEOF)
		return lexText
	case l.isIdentifier(r):
		switch nextRune := l.peek(); {
		case l.isIdentifier(nextRune):
			return lexLetter
		default:
			return lexIdentifier
		}
	default:
		l.backup()
		l.emit(itemLetter)
		return lexText
	}
}

// lexIdentifier At this point, the word is already formed,
// thus we need to verify whether is either an identifier or keyword.
func lexIdentifier(l *lexer) stateFn {
	if !strings.Contains("_", l.input[l.start:l.pos]) { // Verify whether is a keyword, only if it does not contain underline character
		if !l.emitIfKeyword() { // If not a keyword, emit an identifier then go back to initial state.
			l.emit(itemIdentifier)
			return lexText
		}
		return lexText
	}
	l.emit(itemIdentifier)
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
		l.emit(itemNumber)
		return lexText
	}
	l.emit(itemDigit)
	return lexText
	// backup codes  todo: review before delete.
	/*if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	} */

	/*	if strings.IndexRune(".", r) >= 0 {
			rStart, _ := utf8.DecodeRuneInString(l.input[l.start:])
			if unicode.IsNumber(rStart) {
				return true
			}
		}
		return false
	*/
}

func (l *lexer) isIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || strings.IndexRune("_", r) >= 0
	/*	if unicode.IsNumber(r) {
			rStart, _ := utf8.DecodeRuneInString(l.input[l.start:])
			if unicode.IsLetter(rStart) {
				return true
			}
		}
		return false*/ //  Backup Code (for now) todo: review before delete.
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
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
		l.emit(itemKeyword)
		return true
	case word == varKeyword:
		l.emit(itemKeyword)
		return true
	case word == constKeyword:
		l.emit(itemKeyword)
		return true
	case word == registerKeyword:
		l.emit(itemKeyword)
		return true
	case word == functionKeyword:
		l.emit(itemKeyword)
		return true
	case word == procedureKeyword:
		l.emit(itemKeyword)
		return true
	case word == returnKeyword:
		l.emit(itemKeyword)
		return true
	case word == mainKeyword:
		l.emit(itemKeyword)
		return true
	case word == ifKeyword:
		l.emit(itemKeyword)
		return true
	case word == elseKeyword:
		l.emit(itemKeyword)
		return true
	case word == whileKeyword:
		l.emit(itemKeyword)
		return true
	case word == readKeyword:
		l.emit(itemKeyword)
		return true
	case word == writeKeyword:
		l.emit(itemKeyword)
		return true
	case word == integerKeyword:
		l.emit(itemKeyword)
		return true
	case word == realKeyword:
		l.emit(itemKeyword)
		return true
	case word == booleanKeyword:
		l.emit(itemKeyword)
		return true
	case word == charKeyword:
		l.emit(itemKeyword)
		return true
	case word == stringKeyword:
		l.emit(itemKeyword)
		return true
	case word == trueKeyword:
		l.emit(itemKeyword)
		return true
	case word == falseKeyword:
		l.emit(itemKeyword)
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
	fmt.Println("Item list: ")
	for i := range l.items {
		fmt.Println("value: ", i.val)
		fmt.Println("type: ", i.typ)
		fmt.Println("length: ", len(string((i.val))))
		fmt.Println("------------------")
	}
}
