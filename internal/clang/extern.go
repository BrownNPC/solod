package clang

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

// externInfo holds metadata parsed from a so:extern directive.
type externInfo struct {
	name    string // C name override (empty = use default)
	nodecay bool   // skip decay for call args
}

// collectFileExterns collects extern symbols from a single file's declarations.
func (g *Generator) collectFileExterns(typesInfo *types.Info, file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			found, info := parseExternDirective(d.Doc)
			if !found {
				continue
			}
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					g.markExtern(typesInfo.Defs[s.Name], info)
					g.markExternFields(typesInfo, s, info)
				case *ast.ValueSpec:
					for _, name := range s.Names {
						g.markExtern(typesInfo.Defs[name], info)
					}
				}
			}
		case *ast.FuncDecl:
			found, info := parseExternDirective(d.Doc)
			if d.Body == nil || found {
				if d.Recv != nil {
					g.markExtern(typesInfo.Defs[d.Name], info)
				} else {
					g.markExtern(typesInfo.Defs[d.Name], info)
				}
			}
		}
	}
}

// callExtern returns the extern metadata for a call expression, if it
// targets an extern C function.
func (g *Generator) callExtern(call *ast.CallExpr) (externInfo, bool) {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Local package call.
		return g.getExtern(g.types.Uses[fun])
	case *ast.SelectorExpr:
		// Package-qualified call (e.g. stdio.Printf).
		if ident, ok := fun.X.(*ast.Ident); ok {
			if _, ok := g.types.Uses[ident].(*types.PkgName); ok {
				return g.getExtern(g.types.Uses[fun.Sel])
			}
		}
		// Function pointer field on an extern struct (e.g. acc.write(...)).
		return g.callExternField(fun)
	}
	return externInfo{}, false
}

// callExternField checks whether a selector targets a function pointer field
// on an extern struct (e.g. acc.write).
func (g *Generator) callExternField(sel *ast.SelectorExpr) (externInfo, bool) {
	selection, ok := g.types.Selections[sel]
	if !ok || selection.Kind() != types.FieldVal {
		return externInfo{}, false
	}
	info, ok := g.externs[selection.Obj()]
	return info, ok
}

// markExternFields registers function pointer fields of an extern struct type,
// so that calls like acc.write(...) can be resolved via a map lookup.
func (g *Generator) markExternFields(typesInfo *types.Info, spec *ast.TypeSpec, info externInfo) {
	obj := typesInfo.Defs[spec.Name]
	if obj == nil {
		return
	}
	st, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return
	}
	fieldInfo := externInfo{nodecay: info.nodecay}
	for field := range st.Fields() {
		if _, ok := field.Type().Underlying().(*types.Signature); ok {
			g.externs[field] = fieldInfo
		}
	}
}

// markExtern marks a types.Object as extern.
func (g *Generator) markExtern(obj types.Object, info externInfo) {
	g.externs[obj] = info
}

// hasExtern reports whether a types.Object is marked as extern.
func (g *Generator) hasExtern(obj types.Object) bool {
	_, ok := g.externs[obj]
	return ok
}

// getExtern returns the extern metadata for a types.Object.
func (g *Generator) getExtern(obj types.Object) (externInfo, bool) {
	info, ok := g.externs[obj]
	return info, ok
}

// cIntrinsic checks whether an expression is a c.Raw or c.Val call
// and returns the raw string content. The argument must be a string literal.
func (g *Generator) cIntrinsic(expr ast.Expr) (string, bool) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return "", false
	}
	// Unwrap IndexExpr for generic calls like c.Val[T]("...").
	fun := call.Fun
	if idx, ok := fun.(*ast.IndexExpr); ok {
		fun = idx.X
	}
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok || (sel.Sel.Name != "Raw" && sel.Sel.Name != "Val") {
		return "", false
	}
	obj := g.types.Uses[sel.Sel]
	if obj == nil || obj.Pkg() == nil || obj.Pkg().Path() != "solod.dev/so/c" {
		return "", false
	}
	if len(call.Args) != 1 {
		g.fail(call, "Raw C call requires exactly one argument")
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		g.fail(call, "Raw C call argument must be a string literal")
	}

	// Extract raw content: strip backticks for raw strings, unquote interpreted strings.
	var s string
	if strings.HasPrefix(lit.Value, "`") {
		s = lit.Value[1 : len(lit.Value)-1]
	} else {
		var err error
		s, err = strconv.Unquote(lit.Value)
		if err != nil {
			g.fail(lit, "Raw C call: invalid string literal")
		}
	}
	return dedent(s), true
}

// hasInlineDirective reports whether a comment group
// contains the so:inline directive.
func hasInlineDirective(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, c := range doc.List {
		if strings.TrimSpace(c.Text) == "//so:inline" {
			return true
		}
	}
	return false
}

// parseExternDirective checks if a comment group contains the so:extern
// directive and parses its options (name override and nodecay flag).
func parseExternDirective(doc *ast.CommentGroup) (bool, externInfo) {
	if doc == nil {
		return false, externInfo{}
	}
	for _, c := range doc.List {
		text := strings.TrimSpace(c.Text)
		rest, ok := strings.CutPrefix(text, "//so:extern")
		if !ok {
			continue
		}
		var info externInfo
		for tok := range strings.FieldsSeq(rest) {
			if tok == "nodecay" {
				info.nodecay = true
			} else {
				info.name = tok
			}
		}
		return true, info
	}
	return false, externInfo{}
}

// dedent removes common leading whitespace from a multi-line string.
// Also trims leading and trailing blank lines.
func dedent(s string) string {
	lines := strings.Split(s, "\n")

	// Trim leading and trailing blank lines.
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return ""
	}

	// Find the shortest whitespace prefix among non-empty lines.
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent <= 0 {
		return strings.Join(lines, "\n")
	}

	// Strip the common prefix.
	for i, line := range lines {
		if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}
	return strings.Join(lines, "\n")
}
