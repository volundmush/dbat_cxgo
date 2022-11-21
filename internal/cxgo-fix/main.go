package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	path := "."
	if len(args) >= 1 {
		path = args[0]
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var last error
	for _, p := range pkgs {
		for name, f := range p.Files {
			fpath := filepath.Join(path, name)
			if err := processFile(fset, fpath, f); err != nil {
				fmt.Fprintln(os.Stderr, name, err)
				last = err
			}
		}
	}
	return last
}

var unwrapCStringFuncs = map[string][]int{
	"send_to_char":     {1},
	"write_to_output":  {1},
	"vwrite_to_output": {1},
	"send_to_all":      {0},
	"send_to_room":     {1},
	"send_to_range":    {2},
	"send_to_outdoor":  {0},
	"send_to_planet":   {2},
	"send_to_moon":     {0},
	"send_to_imm":      {0},
	"act":              {0},
	"act_to_room":      {0},
	"basic_mud_log":    {0},
	"basic_mud_vlog":   {0},
	"mudlog":           {3},
}
var replaceTrueFalseFuncs = map[string][]int{
	"act": {1},
}

func processFile(fset *token.FileSet, path string, f *ast.File) error {
	changed := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CallExpr:
			fnc, _ := identFullName(n.Fun)
			if argi, ok := unwrapCStringFuncs[fnc]; ok {
				for _, i := range argi {
					if arg, ok := unwrapCString(n.Args[i]); ok {
						n.Args[i] = arg
						changed = true
					}
				}
			}
			if argi, ok := replaceTrueFalseFuncs[fnc]; ok {
				for _, i := range argi {
					if arg, ok := replaceTrueFalse(n.Args[i]); ok {
						n.Args[i] = arg
						changed = true
					}
				}
			}
		}
		return true
	})
	if !changed {
		return nil
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func identFullName(n ast.Node) (string, bool) {
	switch n := n.(type) {
	case *ast.Ident:
		return n.Name, true
	case *ast.SelectorExpr:
		x, ok := identFullName(n.X)
		if !ok {
			return "", false
		}
		return x + "." + n.Sel.Name, true
	}
	return "", false
}

func isZero(n ast.Node) bool {
	x, ok := n.(*ast.BasicLit)
	if !ok || x.Kind != token.INT {
		return false
	}
	v, err := strconv.ParseInt(x.Value, 0, 8)
	return err == nil && v == 0
}

func unwrapCString(n ast.Expr) (ast.Expr, bool) {
	switch n := n.(type) {
	default:
		return n, false
	case *ast.CallExpr: // libc.CString("...") -> "..."
		if len(n.Args) != 1 {
			return n, false
		}
		fnc, ok := identFullName(n.Fun)
		if !ok {
			return n, false
		}
		switch fnc {
		default:
			return n, false
		case "libc.CString":
			return n.Args[0], true
		}
	case *ast.UnaryExpr: // &x[0] -> libc.GoStringS(x[0:])
		if n.Op != token.AND {
			return n, false
		}
		idx, ok := n.X.(*ast.IndexExpr)
		if !ok {
			return n, false
		}
		low := idx.Index
		if isZero(low) {
			low = nil
		}
		return &ast.CallExpr{
			Fun: ast.NewIdent("libc.GoStringS"),
			Args: []ast.Expr{&ast.SliceExpr{
				X:   idx.X,
				Low: low,
			}},
		}, true
	case *ast.SelectorExpr,
		//*ast.Ident,
		*ast.IndexExpr: // x -> libc.GoString(x)
		return &ast.CallExpr{
			Fun:  ast.NewIdent("libc.GoString"),
			Args: []ast.Expr{n},
		}, true
	}
}

var (
	astTrue  = ast.NewIdent("true")
	astFalse = ast.NewIdent("false")
)

func replaceTrueFalse(n ast.Expr) (ast.Expr, bool) {
	switch n := n.(type) {
	case *ast.Ident:
		switch n.Name {
		case "TRUE":
			return astTrue, true
		case "FALSE":
			return astFalse, true
		case "true", "false":
			return n, false
		}
		return &ast.BinaryExpr{X: n, Op: token.NEQ, Y: &ast.BasicLit{Kind: token.INT, Value: "0"}}, true
	case *ast.SelectorExpr:
		return &ast.BinaryExpr{X: n, Op: token.NEQ, Y: &ast.BasicLit{Kind: token.INT, Value: "0"}}, true
	case *ast.BasicLit:
		if n.Kind == token.INT {
			switch n.Value {
			case "1":
				return astTrue, true
			case "0":
				return astFalse, true
			}
		}
	}
	return n, false
}
