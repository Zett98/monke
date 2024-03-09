package repl

import (
	"bytes"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestReplLoop(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	let add = fn(x, y) { x + y; }
	let result = add(five, ten);
	5 / 5 / -5;
	5 < 10 > 5;

	if (5 < 10) { return true; } else { return false; }

	10 == 10;
    10 != 9;
	let x 12 * 3
	`
	in := bytes.NewBufferString(input)
	out := bytes.NewBuffer(make([]byte, 0))
	t.Run("repl tests", func(t *testing.T) {
		Start(in, out)
		result := string(out.Bytes())
		snaps.MatchSnapshot(t, result)
	})
}
func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}
