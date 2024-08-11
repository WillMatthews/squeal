package main

import (
	"fmt"
	"slices"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type CreateStmt struct {
	Relation string
	Columns  []CreateStmtColumn
}

type CreateStmtColumn struct {
	Name        string
	Type        string
	Constraints []string
}

func printCreate(cs CreateStmt) string {
	builder := new(strings.Builder)
	builder.WriteString("create table " + cs.Relation + " (\n")

	typePosition, constraintPosition := calcPositions(cs.Columns)

	writeName := func(b *strings.Builder, name string, position int) {
		b.WriteString(name)
		spaces := position - len(name)
		for i := 0; i < spaces; i++ {
			b.WriteString(" ")
		}
	}

	writeType := func(b *strings.Builder, typ string, position int) {
		b.WriteString(typ)
		spaces := position - len(typ)
		for i := 0; i < spaces; i++ {
			b.WriteString(" ")
		}
	}

	writeConstraints := func(b *strings.Builder, cons []string) {
		for i, con := range cons {
			builder.WriteString(con)
			if i != len(cons)-1 {
				builder.WriteString(" ")
			}
		}
	}

	for iCol, col := range cs.Columns {
		builder.WriteString("    ")
		writeName(builder, col.Name, typePosition)

		// If no constraints, do not add a space
		position := constraintPosition
		if len(col.Constraints) == 0 {
			position = 0
		}
		writeType(builder, col.Type, position)
		writeConstraints(builder, col.Constraints)

		if iCol != len(cs.Columns)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	builder.WriteString(")")

	fmt.Println(builder.String())
	return builder.String()
}

func calcPositions(columns []CreateStmtColumn) (int, int) {
	typePosition, constraintPosition := 0, 0
	for _, col := range columns {
		if len(col.Name) > typePosition {
			typePosition = len(col.Name) + 1
		}
		if len(col.Type) > constraintPosition {
			constraintPosition = len(col.Type) + 1
		}
	}
	return typePosition, constraintPosition
}

func parseCreate(stmt *pg_query.CreateStmt) (CreateStmt, error) {
	fmt.Println(stmt)

	tbl := CreateStmt{
		Relation: stmt.Relation.Relname,
	}

	for _, elt := range stmt.GetTableElts() {
		switch elt.Node.(type) {
		case *pg_query.Node_ColumnDef:
			col := CreateStmtColumn{
				Name:        elt.GetColumnDef().Colname,
				Type:        getColType(elt.GetColumnDef()),
				Constraints: parseConstraints(elt.GetColumnDef().GetConstraints()),
			}

			tbl.Columns = append(tbl.Columns, col)
		default:
			panic("Unknown table element type" + fmt.Sprintf("%T", elt.Node))
		}
	}
	return tbl, nil
}

func parseConstraints(constraints []*pg_query.Node) []string {
	var cons []string
	for _, c := range constraints {
		switch c.Node.(type) {
		case *pg_query.Node_Constraint:
			cons = append(cons, c.GetConstraint().Contype.String())
		default:
			panic("Unknown constraint type" + fmt.Sprintf("%T", c.Node))
		}
	}
	return cons
}

func parseSval(s string) string {
	switch s {
	case "int4":
		return "integer"
	case "bool":
		return "boolean"
	default:
		return s
	}
}

func getColType(cdef *pg_query.ColumnDef) string {
	names := cdef.TypeName.Names
	typmods := cdef.TypeName.Typmods

	mod := ""
	if len(typmods) > 0 {
		ival := typmods[0].GetAConst().GetIval().Ival
		mod = fmt.Sprintf("(%d)", ival)
	}

	endSval := names[len(names)-1].GetString_().GetSval()
	startSval := names[0].GetString_().GetSval()

	avoid := []string{"pg_catalog", "public", "serial", "date", "uuid"}
	if slices.Contains(avoid, startSval) {
		return parseSval(endSval) + mod
	}

	return startSval + "." + parseSval(endSval) + mod
}
