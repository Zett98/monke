package eval

import (
	"bytes"
	"encoding/json"
	"monke/ast"
	"monke/lexer"
	"monke/parser"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func parse(t *testing.T, s *string) ast.Node {
	lex := lexer.New(s)
	parser := parser.New(lex)
	prog := parser.ParseProgram()
	if err := parser.Errors(); len(err) != 0 {
		t.Logf("Unexpected error %v\n", err)
		t.FailNow()
	}
	return prog
}
func TestEval(t *testing.T) {
	type args struct {
		node ast.Node
	}
	tests := []struct {
		name string
		args string
	}{
		{
			"eval int",
			"5",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"truthy `!true`",
			"!true",
		},
		{
			"truthy `!false`",
			"!false",
		},
		{
			"truthy `!5`",
			"!5",
		},
		{
			"truthy `!!true`",
			"!!true",
		},
		{
			"truthy `!!false`",
			"!!false",
		},
		{
			"truthy `!!5`",
			"!!5",
		},
		{
			"truthy `!!0`",
			"!!0",
		},
		{
			"prefix minus",
			"-5",
		},
		{
			"infix",
			"true == true",
		},
		{
			"infix",
			"false == false",
		},
		{
			"infix",
			"true == false",
		},
		{
			"infix",
			"true != false",
		},
		{
			"infix",
			"false != true",
		},
		{
			"infix",
			"(1 < 2) == true",
		},
		{
			"infix",
			"(1 < 2) == false",
		},
		{
			"infix",
			"(1 > 2) == true",
		},
		{
			"infix",
			"(1 > 2) == false",
		},
		{
			"cond",
			"if (1 < 2) { 10 }",
		},
		{
			"cond",
			"if (1 > 2) { 10 }",
		},
		{
			"cond",
			"if (1 > 2) { 10 } else { 20 }",
		},
		{
			"cond",
			"if (1 < 2) { 10 } else { 20 }",
		},
		{
			"return",
			"return 10;",
		},
		{
			"sample program",
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`,
		},
		{
			"error",
			"5 + true",
		},
		{
			"environment",
			"let x = 5; return x + 5;",
		},
		{
			"environment",
			`let newAdder = fn(x) { fn(y) { x + y }; };
				let addTwo = newAdder(2);
				addTwo(2);`,
		},
		{
			"environment",
			"let f = fn(x) { let x = 10; x + 2; }; f(5);",
		},
		{
			"environment",
			"let x = 1; let f = fn(x) { x + 2; };  f(5);",
		},
	}
	for _, tt := range tests {
		env := NewEnvironment()
		t.Run(tt.name, func(t *testing.T) {
			parsed := parse(t, &tt.args)

			result := struct {
				Input string `json:"input"`
				Out   string `json:"out"`
			}{
				Input: tt.args,
				Out:   Eval(parsed, env).Inspect(),
			}
			buf := bytes.NewBufferString("")
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			enc.SetIndent("", "  ")
			err := enc.Encode(result)
			if err != nil {
				t.Errorf("%v", err)
				t.FailNow()
			}
			// NB: object.Bool.false are omitted during printing
			// because `false` is a default value
			snaps.MatchSnapshot(t,
				buf.String())
		})
	}
}

func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}
