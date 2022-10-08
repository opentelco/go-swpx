package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ   itemType // The type of this item.
	pos   Pos      // The starting position, in bytes, of this item in the input string.
	val   string   // The value of this item.
	line  int      // The line number at the start of this item.
	block bool     // is the item part of a block
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemCommandLineTermination:
		return "command termination"
	case i.typ == itemSpace:
		return "space"
	case i.typ == itemError:
		return i.val
	}
	return fmt.Sprintf("block:%v, %s:%s", i.block, i.typ, strings.TrimSpace(i.val))
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemKeyword
	itemText
	itemChar
	itemComment
	itemString
	itemBreak
	itemSpace
	itemIdentifier
	itemContinue
	itemSectionMark // section header
	itemCommandLineTermination
	itemField
	itemBool
)

var key = map[string]itemType{
	"#": itemSectionMark,
}

func (i itemType) String() string {
	switch i {
	case itemError:
		return "error"
	case itemEOF:
		return "EOF"
	case itemKeyword:
		return "keyword"
	case itemText:
		return "text"
	case itemString:
		return "string"
	case itemSpace:
		return "space"
	case itemChar:
		return "char"
	case itemComment:
		return "comment"
	case itemSectionMark:
		return "######### section mark"
	case itemIdentifier:
		return "identifier"

	case itemCommandLineTermination:
		return "command termination"
	}

	return fmt.Sprintf("%d: unknown", i)

}

const eof = -1

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

const (
	spaceChars    = " \t\r\n"  // These are the space characters defined by Go itself.
	trimMarker    = '-'        // Attached to left/right delimiter, trims trailing spaces from preceding/following text.
	trimMarkerLen = Pos(1 + 1) // marker plus space before or after
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name             string // the name of the input; used only for error reports
	input            string // the string being scanned
	sectionDelimiter string // start of section marker
	pos              Pos    // current position in the input
	start            Pos    // start position of this item
	atEOF            bool   // we have hit the end of input and returned eof
	parenDepth       int    // nesting depth of ( ) exprs
	line             int    // 1+number of newlines seen
	startLine        int    // start line of this item
	item             item   // item to return to parser
	insideSection    bool   // are we inside an config section
	insideBlock      bool   // are we inside a config block?
	options          lexOptions
}

// lexOptions control behavior of the lexer. All default to false.
type lexOptions struct {
	emitComment bool // emit itemComment tokens.
	breakOK     bool // break keyword allowed
	continueOK  bool // continue keyword allowed
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.atEOF = true
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += Pos(w)
	if r == '\n' {
		l.line++
	}
	return r
}

const (
	sectionDelimiter    = "\n#\n"
	sectionVersionStamp = "\n!\n"
	commandDelimter     = "\n"
	comment             = "//"
)

// lex creates a new scanner for the input string.
func lex(name, input, secDel string) *lexer {
	if secDel == "" {
		secDel = sectionDelimiter
	}
	l := &lexer{
		name:             name,
		input:            input,
		sectionDelimiter: secDel,
		line:             1,
		startLine:        1,
		insideSection:    false,
	}
	return l
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	l.item = item{itemEOF, l.pos, "EOF", l.startLine, l.insideBlock}
	state := lexText

	//If we are inside an action, we need to return the next item
	if l.insideSection {
		state = lexInsideSection
	}

	for {
		state = state(l)
		if state == nil {
			return l.item
		}
	}
}

// emit passes the trailing text as an item back to the parser.
func (l *lexer) emit(t itemType) stateFn {
	return l.emitItem(l.thisItem(t))
}

// emitItem passes the specified item to the parser.
func (l *lexer) emitItem(i item) stateFn {
	l.item = i
	return nil
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	if !l.atEOF && l.pos > 0 {
		r, w := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= Pos(w)
		// Correct newline count.
		if r == '\n' {
			l.line--
		}
	}
}

// thisItem returns the item at the current input point with the specified type
// and advances the input.
func (l *lexer) thisItem(t itemType) item {
	i := item{t, l.start, l.input[l.start:l.pos], l.startLine, l.insideBlock}
	l.start = l.pos
	l.startLine = l.line
	return i
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || r == '/' || r == '.' || r == '-' || r == ':' || r == '@' || r == '#'
}

func hasLeftTrimMarker(s string) bool {
	return len(s) >= 2 && s[0] == trimMarker && isSpace(rune(s[1]))
}

func hasRightTrimMarker(s string) bool {
	return len(s) >= 2 && isSpace(rune(s[0])) && s[1] == trimMarker
}

// rightTrimLength returns the length of the spaces at the end of the string.
func rightTrimLength(s string) Pos {
	return Pos(len(s) - len(strings.TrimRight(s, spaceChars)))
}

// leftTrimLength returns the length of the spaces at the beginning of the string.
func leftTrimLength(s string) Pos {
	return Pos(len(s) - len(strings.TrimLeft(s, spaceChars)))
}

// ignore skips over the pending input before this point.
// It tracks newlines in the ignored text, so use it only
// for text that is skipped without calling l.next.
func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLine = l.line
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...any) stateFn {
	l.item = item{itemError, l.start, fmt.Sprintf(format, args...), l.startLine, l.insideBlock}
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil
}

func lexText(l *lexer) stateFn {
	if x := strings.Index(l.input[l.pos:], l.sectionDelimiter); x >= 0 {
		if x > 0 {
			l.pos += Pos(x)
			// Do we trim any trailing space?
			trimLength := Pos(0)
			delimEnd := l.pos + Pos(len(l.sectionDelimiter))
			if hasLeftTrimMarker(l.input[delimEnd:]) {
				trimLength = rightTrimLength(l.input[l.start:l.pos])
			}
			l.pos -= trimLength
			l.line += strings.Count(l.input[l.start:l.pos], "\n")
			i := l.thisItem(itemText)
			l.pos += trimLength
			l.ignore()
			if len(i.val) > 0 {
				return l.emitItem(i)
			}
		}

		return lexSectionStart
	}
	l.pos = Pos(len(l.input))

	// Correctly reached EOF.
	if l.pos > l.start {
		l.line += strings.Count(l.input[l.start:l.pos], "\n")
		return l.emit(itemText)
	}
	return l.emit(itemEOF)
}

func lexSectionStart(l *lexer) stateFn {
	l.pos += Pos(len(l.sectionDelimiter))
	trimSpace := hasLeftTrimMarker(l.input[l.pos:])
	afterMarker := Pos(0)
	if trimSpace {
		afterMarker = trimMarkerLen
	}
	if strings.HasPrefix(l.input[l.pos+afterMarker:], comment) {
		l.pos += afterMarker
		l.ignore()
		return lexComment
	}
	i := l.thisItem(itemSectionMark)
	l.insideBlock = false
	l.insideSection = true
	l.pos += afterMarker
	l.ignore()
	l.parenDepth = 0
	return l.emitItem(i)
}

func lexInsideSection(l *lexer) stateFn {

	// if l.pos == Pos(1) {
	// 	fmt.Println("lexInsideSection", l.pos, l.input[l.pos:])
	// }

	switch r := l.next(); {
	case r == '\n':
		return l.emit(itemCommandLineTermination)

	case r == eof:
		return l.emit(itemEOF)

	case r == '"':
		return lexQuote

	case isSpace(r):
		if l.pos > l.start {
			l.insideBlock = true
		}
		return lexSpace

	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case r <= unicode.MaxASCII && unicode.IsPrint(r):
		return l.emit(itemChar)

	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
}

func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	return l.emit(itemString)
}
func lexIdentifier(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}
			switch {
			case key[word] > itemKeyword:
				item := key[word]
				if item == itemBreak && !l.options.breakOK || item == itemContinue && !l.options.continueOK {
					return l.emit(itemIdentifier)
				}
				return l.emit(item)
			case word[0] == '.':
				return l.emit(itemField)
			case word == "true", word == "false":
				return l.emit(itemBool)
			default:
				return l.emit(itemIdentifier)
			}
		}
	}
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(':
		return true
	}

	return strings.HasPrefix(l.input[l.pos:], l.sectionDelimiter)
}

func lexLineTermination(lexer *lexer) stateFn {
	lexer.pos += Pos(len(commandDelimter))
	lexer.ignore()
	return lexText
}

func lexComment(l *lexer) stateFn {
	l.pos += Pos(len(comment))
	x := strings.Index(l.input[l.pos:], commandDelimter)
	if x < 0 {
		return l.errorf("unclosed comment")
	}
	l.pos += Pos(x + len(commandDelimter))

	i := l.thisItem(itemComment)

	l.ignore()
	if l.options.emitComment {
		return l.emitItem(i)
	}
	return lexText
}

// lexSpace scans a run of space characters.
// We have not consumed the first space, which is known to be present.
// Take care if there is a trim-marked right delimiter, which starts with a space.
func lexSpace(l *lexer) stateFn {
	var r rune
	var numSpaces int
	for {
		r = l.peek()
		if !isSpace(r) {
			break
		}
		l.next()
		numSpaces++
	}

	return l.emit(itemSpace)
}
