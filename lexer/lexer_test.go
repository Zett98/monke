package lexer

import (
	"monke/token"
	"os"

	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
	x + y;
	};
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	if (5 < 10) {
	return true;
	} else {
	return false;
	}
	10 == 10;
    10 != 9;
	`

	l := New(&input)
	t.Run("lexer tests", func(t *testing.T) {
		var tok token.Token
		for tok.Type != token.EOF {
			tok = l.NextToken()
			snaps.MatchSnapshot(t, tok)
		}
	})
}

func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}
