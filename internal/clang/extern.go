package clang

import (
	"go/ast"
	"go/types"
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
