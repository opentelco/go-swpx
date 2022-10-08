package lexer

import "fmt"

type Tree struct {
	lex       *lexer
	peekCount int
	token     [3]item
}

func NewTree(lexer *lexer) *Tree {
	return &Tree{lex: lexer}

}

func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

// peek returns the next token but does not consume it
func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

// backup2 backs the input stream up two tokens.
// The zeroth token is already there.
func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

// nextNonSpace returns the next non-space item.
func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}

func (t *Tree) parse() {
	var s string
	for t.peek().typ != itemEOF {

		if t.peek().typ == itemSectionMark {
			fmt.Println("------ section ------ ")
		}
		if t.peek().typ == itemCommandLineTermination {
			fmt.Println("------ ----- ------ ")
		}
		token := t.nextNonSpace()
		if token.typ == itemSectionMark {
			continue
		}

		s += fmt.Sprintf("%s ", token.val)
		if t.peek().typ == itemCommandLineTermination {
			fmt.Println(s)
			s = ""
			t.next()
			continue
		}

	}
}
