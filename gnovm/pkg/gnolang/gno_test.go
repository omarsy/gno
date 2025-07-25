package gnolang

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"unsafe"

	// "github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMachine(b *testing.B, numValues, numStmts, numExprs, numBlocks, numFrames, numExceptions int) *Machine {
	b.Helper()

	m := &Machine{
		Ops:       make([]Op, 100),
		Values:    make([]TypedValue, numValues),
		Exprs:     make([]Expr, numExprs),
		Stmts:     make([]Stmt, numStmts),
		Blocks:    make([]*Block, numBlocks),
		Frames:    make([]Frame, numFrames),
		Exception: nil,
	}
	return m
}

func BenchmarkStringLargeData(b *testing.B) {
	m := setupMachine(b, 10000, 5000, 5000, 2000, 3000, 1000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = m.String()
	}
}

func TestBuiltinIdentifiersShadowing(t *testing.T) {
	t.Parallel()
	tests := map[string]string{}

	uverseNames := []string{
		"iota",
		"append",
		"cap",
		"copy",
		"delete",
		"len",
		"make",
		"new",
		"panic",
		"print",
		"println",
		"recover",
		"nil",
		"bool",
		"byte",
		"float32",
		"float64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"rune",
		"string",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"error",
		"true",
		"false",
		"any",
	}

	for _, name := range uverseNames {
		tests[("struct builtin " + name)] = fmt.Sprintf(`
			package test

			type %v struct {}

			func main() {}
		`, name)

		tests[("var builtin " + name)] = fmt.Sprintf(`
			package test

			func main() {
				%v := 1
			}
		`, name)

		tests[("var declr builtin " + name)] = fmt.Sprintf(`
			package test

			func main() {
				var %v int
			}
		`, name)
	}

	for n, s := range tests {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("shadowing test: `%s` should have failed but didn't\n", n)
				}
			}()

			m := NewMachine("test", nil)
			nn := MustParseFile("main.go", s)
			m.RunFiles(nn)
			m.RunMain()
		})
	}
}

func TestConvertTo(t *testing.T) {
	t.Parallel()

	testFunc := func(source, msg string) {
		defer func() {
			if len(msg) == 0 {
				return
			}

			r := recover()

			if r == nil {
				t.Fail()
			}

			err := r.(*PreprocessError)
			c := strings.Contains(err.Error(), msg)
			if !c {
				t.Fatalf(`expected "%s", got "%s"`, msg, r)
			}
		}()

		m := NewMachine("test", nil)

		n := MustParseFile("main.go", source)
		m.RunFiles(n)
		m.RunMain()
	}

	type cases struct {
		source string
		msg    string
	}

	tests := []cases{
		{
			`package test
		func main() {
			var t interface{}
			t = 2
			var g = float32(t)
			println(g)
		}
		`, `test/main.go:5:12-22: cannot convert interface{} to float32: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = int(t)
		   println(g)
		}
		`, `test/main.go:5:14-20: cannot convert interface{} to int: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = int8(t)
		   println(g)
		}
		`, `test/main.go:5:14-21: cannot convert interface{} to int8: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = int16(t)
		   println(g)
		}
		`, `test/main.go:5:14-22: cannot convert interface{} to int16: need type assertion`,
		},
		{
			`package test
				func main() {
				   var t interface{}
				   t = 2
				   var g = int32(t)
				   println(g)
				}
				`, `test/main.go:5:16-24: cannot convert interface{} to int32: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = int64(t)
		   println(g)
		}
		`, `test/main.go:5:14-22: cannot convert interface{} to int64: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = uint(t)
		   println(g)
		}
		`, `test/main.go:5:14-21: cannot convert interface{} to uint: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = uint8(t)
		   println(g)
		}
		`, `test/main.go:5:14-22: cannot convert interface{} to uint8: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = uint16(t)
		   println(g)
		}
		`, `test/main.go:5:14-23: cannot convert interface{} to uint16: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = uint32(t)
		   println(g)
		}
		`, `test/main.go:5:14-23: cannot convert interface{} to uint32: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 2
		   var g = uint64(t)
		   println(g)
		}
		`, `test/main.go:5:14-23: cannot convert interface{} to uint64: need type assertion`,
		},

		// Built-in non-numeric types
		{
			`package test
		func main() {
		   var t interface{}
		   t = "hello"
		   var g = string(t)
		   println(g)
		}
		`, `test/main.go:5:14-23: cannot convert interface{} to string: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = true
		   var g = bool(t)
		   println(g)
		}
		`, `test/main.go:5:14-21: cannot convert interface{} to bool: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = 'a'
		   var g = rune(t)
		   println(g)
		}
		`, `test/main.go:5:14-21: cannot convert interface{} to int32: need type assertion`,
		},
		{
			`package test
		func main() {
		   var t interface{}
		   t = byte(65)
		   var g = byte(t)
		   println(g)
		}
		`, `test/main.go:5:14-21: cannot convert interface{} to uint8: need type assertion`,
		},

		{
			`package test
		type MyInt int
		func main() {
		   var t interface{}
		   t = MyInt(2)
		   var g = MyInt(t)
		   println(g)
		}
		`, `test/main.go:6:14-22: cannot convert interface{} to test.MyInt: need type assertion`,
		},
		{
			`package test

		func main() {
			const a int = -1
		   println(uint(a))
		}`,
			`test/main.go:5:14-21: cannot convert constant of type IntKind to UintKind`,
		},
		{
			`package test

		func main() {
			const a int = -1
		   println(uint8(a))
		}`,
			`test/main.go:5:14-22: cannot convert constant of type IntKind to Uint8Kind`,
		},
		{
			`package test

		func main() {
			const a int = -1
		   println(uint16(a))
		}`,
			`test/main.go:5:14-23: cannot convert constant of type IntKind to Uint16Kind`,
		},
		{
			`package test

		func main() {
			const a int = -1
		   println(uint32(a))
		}`,
			`test/main.go:5:14-23: cannot convert constant of type IntKind to Uint32Kind`,
		},
		{
			`package test

		func main() {
			const a int = -1
		   println(uint64(a))
		}`,
			`test/main.go:5:14-23: cannot convert constant of type IntKind to Uint64Kind`,
		},
		{
			`package test

		func main() {
			const a float32 = 1.5
		   println(int32(a))
		}`,
			`test/main.go:5:14-22: cannot convert constant of type Float32Kind to Int32Kind`,
		},
		{
			`package test

		func main() {
		   println(int32(1.5))
		}`,
			`test/main.go:4:14-24: cannot convert (const (1.5 <untyped> bigdec)) to integer type`,
		},
		{
			`package test

		func main() {
			const a float64 = 1.5
		   println(int64(a))
		}`,
			`test/main.go:5:14-22: cannot convert constant of type Float64Kind to Int64Kind`,
		},
		{
			`package test

		func main() {
		   println(int64(1.5))
		}`,
			`test/main.go:4:14-24: cannot convert (const (1.5 <untyped> bigdec)) to integer type`,
		},
		{
			`package test

				func main() {
					const f = float64(1.0)
				   println(int64(f))
				}`,
			``,
		},
	}

	for _, tc := range tests {
		testFunc(tc.source, tc.msg)
	}
}

// run empty main().
func TestRunEmptyMain(t *testing.T) {
	t.Parallel()

	m := NewMachine("test", nil)
	// []Stmt{} != nil, as nil means that in the source code not even the
	// brackets are present and is reserved for external (ie. native) functions.
	main := FuncD("main", nil, nil, []Stmt{})
	m.RunDeclaration(main)
	m.RunMain()
}

// run main() with a for loop.
func TestRunLoopyMain(t *testing.T) {
	t.Parallel()

	m := NewMachine("test", nil)
	c := `package test
func main() {
	for i:=0; i<1000; i++ {
		if i == -1 {
			return
		}
	}
}`
	n := MustParseFile("main.go", c)
	m.RunFiles(n)
	m.RunMain()
}

func BenchmarkPreprocessForLoop(b *testing.B) {
	m := NewMachine("test", nil)
	c := `package test
func main() {
	for i:=0; i<10000; i++ {}
}`
	n := MustParseFile("main.go", c)
	m.RunFiles(n)

	for i := 0; i < b.N; i++ {
		m.RunMain()
	}
}

func TestOptimizeConversion(t *testing.T) {
	t.Parallel()

	m := NewMachine("test", nil)
	c := `package test
func main() {}

func foo(a int) {
    b := int(a)
    println(b)
}`
	n := MustParseFile("main.go", c)
	m.RunFiles(n)
	fn := n.Decls[1].(*FuncDecl)
	as := fn.Body[0].(*AssignStmt)
	ne := as.Rhs[0].(*NameExpr)
	if ne.Name != "a" {
		t.Fatalf("expecting optimized 'a', got %v", ne.String())
	}
}

func BenchmarkIfStatement(b *testing.B) {
	m := NewMachine("test", nil)
	c := `package test
func main() {
	for i:=0; i<10000; i++ {
		if i > 10 {

		}
	}
}`
	n := MustParseFile("main.go", c)
	m.RunFiles(n)

	for i := 0; i < b.N; i++ {
		m.RunMain()
	}
}

func TestDoOpEvalBaseConversion(t *testing.T) {
	m := NewMachine("test", nil)

	type testCase struct {
		input     string
		expect    string
		expectErr bool
	}

	testCases := []testCase{
		// binary
		{input: "0b101010", expect: "42", expectErr: false},
		{input: "0B101010", expect: "42", expectErr: false},
		{input: "0b111111111111111111111111111111111111111111111111111111111111111", expect: "9223372036854775807", expectErr: false},
		{input: "0b0", expect: "0", expectErr: false},
		{input: "0b000000101010", expect: "42", expectErr: false},
		{input: " 0b101010", expectErr: true},
		{input: "0b", expectErr: true},
		{input: "0bXXXX", expectErr: true},
		{input: "42b0", expectErr: true},
		// octal
		{input: "0o42", expect: "34", expectErr: false},
		{input: "0o0", expect: "0", expectErr: false},
		{input: "042", expect: "34", expectErr: false},
		{input: "0777", expect: "511", expectErr: false},
		{input: "0O0000042", expect: "34", expectErr: false},
		{input: "0777777777777777777777", expect: "9223372036854775807", expectErr: false},
		{input: "0o777777777777777777777", expect: "9223372036854775807", expectErr: false},
		{input: "048", expectErr: true},
		{input: "0o", expectErr: true},
		{input: "0oXXXX", expectErr: true},
		{input: "0OXXXX", expectErr: true},
		{input: "0o42x42", expectErr: true},
		{input: "0O42x42", expectErr: true},
		{input: "0420x42", expectErr: true},
		{input: "0o420o42", expectErr: true},
		// hex
		{input: "0x2a", expect: "42", expectErr: false},
		{input: "0X2A", expect: "42", expectErr: false},
		{input: "0x7FFFFFFFFFFFFFFF", expect: "9223372036854775807", expectErr: false},
		{input: "0x2a ", expectErr: true},
		{input: "0x", expectErr: true},
		{input: "0xXXXX", expectErr: true},
		{input: "0xGHIJ", expectErr: true},
		{input: "0x42o42", expectErr: true},
		{input: "0x2ax42", expectErr: true},
		// decimal
		{input: "42", expect: "42", expectErr: false},
		{input: "0", expect: "0", expectErr: false},
		{input: "0000000000", expect: "0", expectErr: false},
		{input: "9223372036854775807", expect: "9223372036854775807", expectErr: false},
	}

	for _, tc := range testCases {
		m.PushExpr(&BasicLitExpr{
			Kind:  INT,
			Value: tc.input,
		})

		if tc.expectErr {
			assert.Panics(t, func() { m.doOpEval() })
		} else {
			m.doOpEval()
			v := m.PopValue()
			assert.Equal(t, tc.expect, v.V.String())
		}
	}
}

func TestEval(t *testing.T) {
	t.Parallel()

	m := NewMachine("test", nil)
	c := `package test
func next(i int) int {
	return i+1
}`
	n := MustParseFile("main.go", c)
	m.RunFiles(n)
	res := m.Eval(Call("next", "1"))
	fmt.Println(res)
}

func assertOutput(t *testing.T, input string, output string) {
	t.Helper()

	buf := new(bytes.Buffer)
	m := NewMachineWithOptions(MachineOptions{
		PkgPath: "test",
		Output:  buf,
	})
	n := MustParseFile("main.go", input)
	m.RunFiles(n)
	m.RunMain()
	assert.Equal(t, output, string(buf.Bytes()))
	err := m.CheckEmpty()
	assert.Nil(t, err)
}

// ----------------------------------------
// Benchmarks

func BenchmarkPreprocess(b *testing.B) {
	pkg := &PackageNode{
		PkgName: "main",
		PkgPath: ".main",
		FileSet: nil,
	}
	pkg.InitStaticBlock(pkg, nil)
	main := FuncD("main", nil, nil, Ss(
		A("mx", ":=", "1000000"),
		For(
			A("i", ":=", "0"),
			X("i < mx"),
			Inc("i"),
		),
	))
	pn := NewPackageNode("hey", "gno.land/p/hey", nil)
	copies := make([]*FuncDecl, b.N)
	for i := 0; i < b.N; i++ {
		copies[i] = main.Copy().(*FuncDecl)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// initStaticBlocks is always performed before a Preprocess
		initStaticBlocks(nil, pn, copies[i])
		main = Preprocess(nil, pkg, copies[i]).(*FuncDecl)
	}
}

type bdataParams struct {
	N     int
	Param string
}

func BenchmarkBenchdata(b *testing.B) {
	const bdDir = "./benchdata"
	files, err := os.ReadDir(bdDir)
	require.NoError(b, err)
	for _, file := range files {
		// Read file and parse template.
		bcont, err := os.ReadFile(filepath.Join(bdDir, file.Name()))
		cont := string(bcont)
		require.NoError(b, err)
		tpl, err := template.New("").Parse(cont)
		require.NoError(b, err)

		// Determine parameters.
		const paramString = "// param: "
		var params []string
		pos := strings.Index(cont, paramString)
		if pos >= 0 {
			paramsRaw := strings.SplitN(cont[pos+len(paramString):], "\n", 2)[0]
			params = strings.Fields(paramsRaw)
		} else {
			params = []string{""}
		}

		for _, param := range params {
			name := file.Name()
			if param != "" {
				name += "_param:" + param
			}
			b.Run(name, func(b *testing.B) {
				// Gen template with N and param.
				var buf bytes.Buffer
				require.NoError(b, tpl.Execute(&buf, bdataParams{
					N:     b.N,
					Param: param,
				}))

				// Set up machine.
				m := NewMachineWithOptions(MachineOptions{
					PkgPath: "main",
					Output:  io.Discard,
				})
				n := MustParseFile("main.go", buf.String())
				m.RunFiles(n)

				b.ResetTimer()
				m.RunMain()
			})
		}
	}
}

// ----------------------------------------
// Unsorted

type Struct1 struct {
	A int
	B int
}

func TestModifyTypeAsserted(t *testing.T) {
	t.Parallel()

	x := Struct1{1, 1}
	var v any = x
	x2 := v.(Struct1)
	x2.A = 2

	// only x2 is changed.
	assert.Equal(t, 1, x.A)
	assert.Equal(t, 1, v.(Struct1).A)
	assert.Equal(t, 2, x2.A)
}

type Interface1 interface {
	Foo()
}

func TestTypeConversion(t *testing.T) {
	t.Parallel()

	x := 1
	var v any = x
	if _, ok := v.(Interface1); ok {
		panic("should not happen")
	}
	v = nil
	if _, ok := v.(Interface1); ok {
		panic("should not happen")
	}
	assert.Panics(t, func() {
		// this would panic.
		z := v.(Interface1)
		fmt.Println(z)
	})
}

func TestSomething(t *testing.T) {
	t.Parallel()

	type Foo struct {
		X any
	}

	type Bar struct {
		X any
		Y bool
	}

	fmt.Println(unsafe.Sizeof(Foo{}))                        // 16
	fmt.Println(unsafe.Sizeof(Foo{X: reflect.ValueOf(0)}.X)) // still 16? weird.
	fmt.Println(unsafe.Sizeof(reflect.ValueOf(0)))           // 24
	fmt.Println(unsafe.Sizeof(Bar{}))
}

// XXX is there a way to test in Go as well as Gno?
func TestDeferOrder(t *testing.T) {
	t.Parallel()

	a := func() func(int, int) int {
		fmt.Println("a constructed")
		return func(x int, y int) int {
			fmt.Println("a called")
			return x + y
		}
	}
	b := func() int {
		fmt.Println("b constructed")
		return 1
	}
	c := func() int {
		fmt.Println("c constructed")
		return 2
	}
	defer a()(b(), c())
	fmt.Println("post defer")

	// should print
	// a constructed
	// b constructed
	// c constructed
	// post defer
	// a called
}

// XXX is there a way to test in Go as well as Gno?
func TestCallOrder(t *testing.T) {
	t.Parallel()

	a := func() func(int, int) int {
		fmt.Println("a constructed")
		return func(x int, y int) int {
			fmt.Println("a called")
			return x + y
		}
	}
	b := func() int {
		fmt.Println("b constructed")
		return 1
	}
	c := func() int {
		fmt.Println("c constructed")
		return 2
	}
	a()(b(), c())

	// should print
	// a constructed
	// b constructed
	// c constructed
	// a called
}

// XXX is there a way to test in Go as well as Gno?
func TestBinaryShortCircuit(t *testing.T) {
	t.Parallel()

	tr := func() bool {
		fmt.Println("t called")
		return true
	}
	fa := func() bool {
		fmt.Println("f called")
		return false
	}
	if fa() && tr() {
		fmt.Println("error")
	} else {
		fmt.Println("done")
	}
}

// XXX is there a way to test in Go as well as Gno?
func TestSwitchDefine(t *testing.T) {
	t.Parallel()

	var x any = 1
	switch y := x.(type) {
	case int:
		fmt.Println("int", y) // XXX
	default:
		fmt.Println("not int")
	}
}

// XXX is there a way to test in Go as well as Gno?
func TestBinaryCircuit(t *testing.T) {
	t.Parallel()

	tr := func() bool {
		fmt.Println("tr() called")
		return true
	}
	fa := func() bool {
		fmt.Println("fa() called")
		return false
	}

	fmt.Println("case 1")
	fmt.Println(fa() && tr())
	fmt.Println("case 1")
	fmt.Println(tr() || fa())

	// should print
	// case 1
	// fa() called
	// false
	// case 1
	// tr() called
	// true
}

func TestMultiAssignment(t *testing.T) {
	t.Parallel()

	buf := make([]int, 4)
	ref := func(i int) *int {
		fmt.Printf("ref(%v) called\n", i)
		return &buf[i]
	}
	val := func(i int) int {
		fmt.Printf("val(%v) called\n", i)
		return i
	}

	*ref(0), *ref(1), *ref(2), *ref(3) = val(11), val(22), val(33), val(44)

	/*
		ref(0) called
		ref(1) called
		ref(2) called
		ref(3) called
		val(11) called
		val(22) called
		val(33) called
		val(44) called
	*/
}

// XXX is there a way to test in Go as well as Gno?
func TestCallLHS(t *testing.T) {
	t.Parallel()

	x := 1
	xptr := func() *int {
		return &x
	}
	*xptr() = 2
	assert.Equal(t, 2, x)
}

// XXX is there a way to test in Go as well as Gno?
func TestCallFieldLHS(t *testing.T) {
	t.Parallel()

	type str struct {
		X int
	}
	x := str{}
	xptr := func() *str {
		return &x
	}
	y := 0
	xptr().X, y = 2, 3
	assert.Equal(t, 2, x.X)
	assert.Equal(t, 3, y)
}
