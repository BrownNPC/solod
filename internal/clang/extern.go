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
func (g *Generator) collectFileExterns(pkgName string, file *ast.File) {
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
					g.markExtern(pkgName, s.Name.Name, info)
					g.markExternFields(pkgName, s, info)
				case *ast.ValueSpec:
					for _, name := range s.Names {
						g.markExtern(pkgName, name.Name, info)
					}
				}
			}
		case *ast.FuncDecl:
			found, info := parseExternDirective(d.Doc)
			if d.Body == nil || found {
				g.markExtern(pkgName, externFuncKey(d), info)
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
		return g.getExtern("", fun.Name)
	case *ast.SelectorExpr:
		// Package-qualified call (e.g. stdio.Printf).
		if ident, ok := fun.X.(*ast.Ident); ok {
			if pkgName, ok := g.types.Uses[ident].(*types.PkgName); ok {
				return g.getExtern(pkgName.Name(), fun.Sel.Name)
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
	recv := selection.Recv()
	if ptr, ok := recv.(*types.Pointer); ok {
		recv = ptr.Elem()
	}
	named, ok := recv.(*types.Named)
	if !ok {
		return externInfo{}, false
	}
	pkgName := ""
	if typePkg := named.Obj().Pkg(); typePkg != nil && typePkg.Path() != g.pkg.PkgPath {
		// The type is imported from another package.
		pkgName = typePkg.Name()
	}
	return g.getExtern(pkgName, named.Obj().Name()+"."+sel.Sel.Name)
}

// markExternFields registers function pointer fields of an extern struct type,
// so that calls like acc.write(...) can be resolved via a map lookup.
func (g *Generator) markExternFields(pkgName string, spec *ast.TypeSpec, info externInfo) {
	obj := g.types.Defs[spec.Name]
	if obj == nil {
		return
	}
	st, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return
	}
	fieldInfo := externInfo{nodecay: info.nodecay}
	for field := range st.Fields() {
		if _, ok := field.Type().(*types.Signature); ok {
			g.markExtern(pkgName, spec.Name.Name+"."+field.Name(), fieldInfo)
		}
	}
}

// markExtern marks a symbol in a package as extern.
func (g *Generator) markExtern(pkgName, name string, info externInfo) {
	if pkgName != "" {
		name = pkgName + "." + name
	}
	g.externs[name] = info
}

// hasExtern reports whether a symbol in a package is marked as extern.
func (g *Generator) hasExtern(pkgName, name string) bool {
	if pkgName != "" {
		name = pkgName + "." + name
	}
	_, ok := g.externs[name]
	return ok
}

// getExtern returns the extern metadata for a symbol.
func (g *Generator) getExtern(pkgName, name string) (externInfo, bool) {
	if pkgName != "" {
		name = pkgName + "." + name
	}
	info, ok := g.externs[name]
	return info, ok
}

// externFuncKey returns a map key for a function or method declaration.
// Functions use their bare name (e.g. "Foo"), while methods use
// "ReceiverType.Name" (e.g. "T.Foo") to avoid collisions.
func externFuncKey(decl *ast.FuncDecl) string {
	if decl.Recv != nil {
		return recvTypeName(decl.Recv.List[0]) + "." + decl.Name.Name
	}
	return decl.Name.Name
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
