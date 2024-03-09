package parser

import (
	"fmt"
	"monke/lexer"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func checkParserErrors(p *Parser) []error {
	errors := p.Errors()
	var allErrors []error
	if len(errors) == 0 {
		return nil
	}
	for _, msg := range errors {
		allErrors = append(allErrors, fmt.Errorf("parser error: %q", msg))
	}
	return allErrors
}

func TestParse(t *testing.T) {
	okInputs := []string{
		`
		let z = 5 * 5;
		let y = 5 < 10;
		let foobar = (2 + 2) / 4;
		`,
		`
		return 5;
		return 10;
		return 993322;
		`,
		"foobar;",
		"5;",
		"-15",
		"!5",
		"5 + 5;",
		"5 * 5;",
		"5 < 5;",
		"5 > 5;",
		"5 / 5;",
		"5 == 5;",
		"5 != 5;",
		"true;",
		"false;",
		`
		if (x < y) {
			 x 
		}
		`,
		`
		if (x < y) {
			 x 
		} else {
			y
		}
		`,
		`
		fn(x, y) { 
			x + y; 
		};
		`,
		`
		fn() { 
			
		};
		`,
		"add(1, 2 * 3, 4 + 5);",
	}

	failingInputs := []string{
		`
		let x 5;
		let = 10;
		let 838383;
		`,
	}

	t.Run("succesfull passes", func(t *testing.T) {
		for _, v := range okInputs {
			lexer := lexer.New(&v)
			parser := New(lexer)

			program := parser.ParseProgram()
			errors := checkParserErrors(parser)
			if errors != nil {
				t.Fatalf("Expected no errors %v", errors)
				t.FailNow()
			}
			snaps.MatchSnapshot(t, v, program.Statements)
		}
	})
	t.Run("failing passes", func(t *testing.T) {
		for _, v := range failingInputs {
			lexer := lexer.New(&v)
			parser := New(lexer)

			parser.ParseProgram()
			errors := checkParserErrors(parser)
			if errors == nil {
				t.Fatalf("Expected to have errors")
				t.FailNow()
			}
			snaps.MatchSnapshot(t, v, errors)
		}
	})
}

func TestPrec(t *testing.T) {
	okInputs := []string{

		"-a * b",
		"!-a",
		"a + b + c",
		"a + b - c",
		"a * b * c",
		"a * b / c",
		"a + b / c",
		"a + b * c + d / e - f",
		"3 + 4; -5 * 5",
		"5 > 4 == 3 < 4",
		"5 < 4 != 3 > 4",
		"3 + 4 * 5 == 3 * 1 + 4 * 5",
		"3 + 4 * 5 == 3 * 1 + 4 * 5",
		"2 / (5 + 5)",
		"!(true == true)",
		"a + add(b * c) + d",
		"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		"add(a + b + c * d / f + g)",
	}

	t.Run("Precedence parsing", func(t *testing.T) {
		for _, v := range okInputs {
			lexer := lexer.New(&v)
			parser := New(lexer)

			program := parser.ParseProgram()
			errors := checkParserErrors(parser)
			if errors != nil {
				t.Fatalf("Expected no errors %v", errors)
				t.FailNow()
			}
			snaps.MatchSnapshot(t, v, program.String())
		}
	})
}
func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run `go-snaps` will sort snapshots
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}
