package participle

import (
	"reflect"
	"testing"
	"text/scanner"

	"github.com/google/go-cmp/cmp"
	"github.com/gotestyourself/gotestyourself/assert"

	"github.com/alecthomas/participle/lexer"
)

func TestStructLexerTokens(t *testing.T) {
	type testScanner struct {
		A string `12`
		B string `34`
	}

	scan := lexStruct(reflect.TypeOf(testScanner{}))
	t12 := lexer.Token{Type: scanner.Int, Value: "12", Pos: lexer.Position{Line: 1, Column: 1}}
	t34 := lexer.Token{Type: scanner.Int, Value: "34", Pos: lexer.Position{Line: 2, Column: 1}}
	assert.DeepEqual(t, t12, scan.Peek())
	assert.DeepEqual(t, 0, scan.field)
	assert.DeepEqual(t, t12, scan.Next())

	assert.DeepEqual(t, t34, scan.Peek())
	assert.DeepEqual(t, 0, scan.field)
	assert.DeepEqual(t, t34, scan.Next())
	assert.DeepEqual(t, 1, scan.field)

	assert.DeepEqual(t, lexer.EOFToken, scan.Next())
}

func TestStructLexer(t *testing.T) {
	g := struct {
		A string `"a"|`
		B string `"b"`
	}{}

	gt := reflect.TypeOf(g)
	r := lexStruct(gt)
	f := []reflect.StructField{}
	s := ""
	for {
		r.Peek()
		rn := r.Next()
		if rn.EOF() {
			break
		}
		f = append(f, r.Field())
		s += string(rn.String())
	}
	assert.DeepEqual(t, `a|b`, s)
	f0 := gt.Field(0)
	f1 := gt.Field(1)
	assert.DeepEqual(t, []reflect.StructField{f0, f0, f1}, f, cmp.Comparer(func(x, y reflect.Type) bool {
		return x == y
	}))
}
