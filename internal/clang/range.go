package clang

import (
	"fmt"
	"go/ast"
	"go/types"
)

// emitIntRange emits a range loop over an integer.
func (g *Generator) emitIntRange(stmt *ast.RangeStmt) {
	w := g.state.writer
	if stmt.Key == nil {
		// Basic form: `for range n { ... }`
		fmt.Fprintf(w, "%sfor (so_int _i = 0; _i < ", g.indent())
		g.emitExpr(stmt.X)
		fmt.Fprintf(w, "; _i++) {\n")
		g.emitBlock(stmt.Body)
		fmt.Fprintf(w, "%s}\n", g.indent())
		return
	}

	key := stmt.Key.(*ast.Ident)
	cType := g.mapType(stmt, g.types.Defs[key].Type())
	fmt.Fprintf(w, "%sfor (%s %s = 0; %s < ", g.indent(), cType, key.Name, key.Name)
	g.emitExpr(stmt.X)
	fmt.Fprintf(w, "; %s++) {\n", key.Name)
	g.emitBlock(stmt.Body)
	fmt.Fprintf(w, "%s}\n", g.indent())
}

// emitArrayRange emits a range loop over a fixed-size array.
func (g *Generator) emitArrayRange(stmt *ast.RangeStmt) {
	if _, ok := stmt.X.(*ast.CompositeLit); ok {
		g.fail(stmt.X, "for-range over literal not supported")
	}
	w := g.state.writer
	if stmt.Key == nil {
		// Basic form: `for range arr { ... }`
		arrType := g.types.TypeOf(stmt.X).Underlying().(*types.Array)
		fmt.Fprintf(w, "%sfor (so_int _i = 0; _i < %d; _i++) {\n", g.indent(), arrType.Len())
		g.emitBlock(stmt.Body)
		fmt.Fprintf(w, "%s}\n", g.indent())
		return
	}

	key := stmt.Key.(*ast.Ident)
	arrType := g.types.TypeOf(stmt.X).Underlying().(*types.Array)
	elemType := g.mapType(stmt, arrType.Elem())
	cType := g.mapType(stmt, g.types.Defs[key].Type())

	fmt.Fprintf(w, "%sfor (%s %s = 0; %s < %d; %s++) {\n",
		g.indent(), cType, key.Name, key.Name, arrType.Len(), key.Name)

	// Emit value variable if present (e.g. `for i, v := range nums`).
	if stmt.Value != nil {
		if valIdent, ok := stmt.Value.(*ast.Ident); ok && valIdent.Name != "_" {
			g.state.indent++
			fmt.Fprintf(w, "%s%s %s = ", g.indent(), elemType, valIdent.Name)
			g.emitExpr(stmt.X)
			fmt.Fprintf(w, "[%s];\n", key.Name)
			g.state.indent--
		}
	}

	g.emitBlock(stmt.Body)
	fmt.Fprintf(w, "%s}\n", g.indent())
}

// emitSliceRange emits a range loop over a slice.
func (g *Generator) emitSliceRange(stmt *ast.RangeStmt) {
	if _, ok := stmt.X.(*ast.CompositeLit); ok {
		g.fail(stmt.X, "for-range over literal not supported")
	}
	w := g.state.writer
	if stmt.Key == nil {
		// Basic form: `for range slice { ... }`
		fmt.Fprintf(w, "%sfor (so_int _i = 0; _i < so_len(", g.indent())
		g.emitExpr(stmt.X)
		fmt.Fprintf(w, "); _i++) {\n")
		g.emitBlock(stmt.Body)
		fmt.Fprintf(w, "%s}\n", g.indent())
		return
	}

	key := stmt.Key.(*ast.Ident)
	sliceType := g.types.TypeOf(stmt.X).Underlying().(*types.Slice)
	elemType := g.mapType(stmt, sliceType.Elem())
	cType := g.mapType(stmt, g.types.Defs[key].Type())

	fmt.Fprintf(w, "%sfor (%s %s = 0; %s < so_len(", g.indent(), cType, key.Name, key.Name)
	g.emitExpr(stmt.X)
	fmt.Fprintf(w, "); %s++) {\n", key.Name)

	// Emit value variable if present (e.g. `for i, v := range nums`).
	if stmt.Value != nil {
		if valIdent, ok := stmt.Value.(*ast.Ident); ok && valIdent.Name != "_" {
			g.state.indent++
			fmt.Fprintf(w, "%s%s %s = so_at(%s, ", g.indent(), elemType, valIdent.Name, elemType)
			g.emitExpr(stmt.X)
			fmt.Fprintf(w, ", %s);\n", key.Name)
			g.state.indent--
		}
	}

	g.emitBlock(stmt.Body)
	fmt.Fprintf(w, "%s}\n", g.indent())
}

// emitStringRange emits a range loop over a string (rune iteration).
func (g *Generator) emitStringRange(stmt *ast.RangeStmt) {
	w := g.state.writer
	if stmt.Key == nil {
		// Basic form: `for range str { ... }`
		fmt.Fprintf(w, "%sfor (so_int _i = 0; _i < so_len(", g.indent())
		g.emitExpr(stmt.X)
		fmt.Fprintf(w, ");) {\n")
		g.state.indent++
		fmt.Fprintf(w, "%sint _iw = 0;\n", g.indent())
		fmt.Fprintf(w, "%sso_utf8_decode(", g.indent())
		g.emitExpr(stmt.X)
		fmt.Fprintf(w, ", _i, &_iw);\n")
		g.state.indent--
		g.emitBlock(stmt.Body)
		g.state.indent++
		fmt.Fprintf(w, "%s_i += _iw;\n", g.indent())
		g.state.indent--
		fmt.Fprintf(w, "%s}\n", g.indent())
		return
	}

	key := stmt.Key.(*ast.Ident)
	cType := g.mapType(stmt, g.types.Defs[key].Type())

	fmt.Fprintf(w, "%sfor (%s %s = 0; %s < so_len(", g.indent(), cType, key.Name, key.Name)
	g.emitExpr(stmt.X)
	fmt.Fprintf(w, ");) {\n")

	// Decode rune and width once per iteration.
	g.state.indent++
	widthVar := "_" + key.Name + "w"
	fmt.Fprintf(w, "%sint %s = 0;\n", g.indent(), widthVar)
	if stmt.Value != nil {
		if valIdent, ok := stmt.Value.(*ast.Ident); ok && valIdent.Name != "_" {
			fmt.Fprintf(w, "%sso_rune %s = so_utf8_decode(", g.indent(), valIdent.Name)
		} else {
			fmt.Fprintf(w, "%sso_utf8_decode(", g.indent())
		}
	} else {
		fmt.Fprintf(w, "%sso_utf8_decode(", g.indent())
	}
	g.emitExpr(stmt.X)
	fmt.Fprintf(w, ", %s, &%s);\n", key.Name, widthVar)
	g.state.indent--

	g.emitBlock(stmt.Body)

	// Advance index by rune width.
	g.state.indent++
	fmt.Fprintf(w, "%s%s += %s;\n", g.indent(), key.Name, widthVar)
	g.state.indent--

	fmt.Fprintf(w, "%s}\n", g.indent())
}
