package lexical_analysis

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TODO: (1) Consertar cast de itemEOF como runa/ fixar tipo padrão como runa .

const (
	itemError itemType = iota // whether an error has ocurred.

	itemDot
	itemEOF
	itemElse          // else keyword
	itemEnd           // end keyword
	itemField         // identifier, starting with '.'
	itemIdentifier    // identifier
	itemIf            // if keyword
	itemLeftMeta      // left meta-string
	itemNumber        // number
	itemLetter        // letter
	itemPipe          // pipe symbol
	itemRange         // range keyword
	itemRawString     // raw quoted string (includes quotes)
	itemRightMeta     // right meta-string
	itemString        // quoted string (includes quotes)
	itemText          // plain text
	leftMeta          = ""
	programIdentifier = "program"
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

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// lexText A partir de um estado inicial, eg., leftMeta, procura por lexemas,
// tratando-os individualmente.
func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], programIdentifier) {
			if l.pos > l.start { // TODO: Verificar o que isso faz.
				fmt.Println("Entrou condicao 0")
				l.emit(itemText)
			}
			return lexProgramIdentifier
		}
		// Se letra, então identificador, ou letra.
		r := l.next()
		if unicode.IsLetter(r) {
			return lexLetter
		}
		if unicode.IsSpace(r) {
			l.ignore()
		}
		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta // Next state.
		}
		if l.next() == rune(itemEOF) {
			break
		}
	}
	// If correctly reached EOF do:
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexLetter(l *lexer) stateFn {
	switch r := l.peek(); {
	case unicode.IsSpace(r):
		l.emit(itemLetter)
		return lexText
	}
	return lexInsideAction
}

func lexProgramIdentifier(l *lexer) stateFn {
	l.pos += len(programIdentifier)
	l.emit(itemIdentifier)
	return lexInsideAction // Now inside {{ }}.
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(leftMeta)
	l.emit(itemLeftMeta)
	return lexInsideAction // Now inside {{ }}.
}

func lexInsideAction(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == rune(itemEOF):
			return nil
		case unicode.IsLetter(r):
			if unicode.IsSpace(l.peek()) {
				l.emit(itemLetter)
				return lexText
			}
		case unicode.IsSpace(r):
			l.ignore()
		}
	}
}

func lexNumber(l *lexer) stateFn {
	l.accept("+-")
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	/*if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q",
			l.input[l.start:l.pos])
	} */
	l.emit(itemNumber)
	return lexInsideAction
}

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return rune(itemEOF)
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// Display display tokens.
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
		fmt.Println("length: ", len(string((i.val))))
		fmt.Println("------------------")
	}
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
