package main

import (
	"fmt"
	"slices"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

const (
	defaultTypeSpace       = 1
	defaultConstraintSpace = 1

	useTabs    = false
	tabSpacing = 4
)

type TableCreationStatement struct {
	Relation string
	Columns  []TableColumn
}

type TableColumn struct {
	Name        string
	Type        string
	Constraints []string
}

func (cs TableCreationStatement) String() string {
	builder := new(strings.Builder)
	builder.WriteString("create table " + cs.Relation + " (\n")

	typePosition, constraintPosition := cs.calcPositions()
	for iCol, col := range cs.Columns {
		writeColumn(builder, col, typePosition, constraintPosition)
		if iCol < len(cs.Columns)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}
	builder.WriteString(");")

	return builder.String()
}

func pad(b *strings.Builder, position int) {
	for i := 0; i < position; i++ {
		b.WriteString(" ")
	}
}

func writeName(b *strings.Builder, name string, position int) {
	b.WriteString(name)
	pad(b, position-len(name))
}

func writeType(b *strings.Builder, typ string, position int) {
	b.WriteString(typ)
	pad(b, position-len(typ))
}

func writeConstraints(b *strings.Builder, cons []string) {
	for i, con := range cons {
		b.WriteString(con)
		if i != len(cons)-1 {
			b.WriteString(" ")
		}
	}
}

func writeTab(b *strings.Builder) {
	if useTabs {
		b.WriteString("\t")
	} else {
		pad(b, tabSpacing)
	}
}

func writeColumn(b *strings.Builder, col TableColumn, typePosition, constraintPosition int) {
	writeTab(b)
	writeName(b, col.Name, typePosition)

	// If no constraints, do not add a space
	position := constraintPosition
	if len(col.Constraints) == 0 {
		position = 0
	}
	writeType(b, col.Type, position)
	writeConstraints(b, col.Constraints)
}

func (cs TableCreationStatement) calcPositions() (int, int) {
	typePosition, constraintPosition := 0, 0
	for _, col := range cs.Columns {
		if len(col.Name) > typePosition {
			typePosition = len(col.Name) + defaultTypeSpace
		}
		if len(col.Type) > constraintPosition {
			constraintPosition = len(col.Type) + defaultConstraintSpace
		}
	}
	return typePosition, constraintPosition
}

func parseCreate(stmt *pg_query.CreateStmt) (TableCreationStatement, error) {
	// fmt.Println(stmt)

	tbl := TableCreationStatement{
		Relation: stmt.Relation.Relname,
	}

	for _, elt := range stmt.GetTableElts() {
		switch elt.Node.(type) {
		case *pg_query.Node_ColumnDef:
			col := TableColumn{
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

func parseDefault(c *pg_query.Constraint) string {
	b := new(strings.Builder)
	b.WriteString("default")
	if c.GetRawExpr() != nil {
		aconst := c.GetRawExpr().GetAConst()
		switch {
		case aconst.Val == nil:
			b.WriteString(" null")
		case aconst.GetBoolval() != nil:
			bval := aconst.GetBoolval().Boolval
			b.WriteString(fmt.Sprintf(" %t", bval))
		case aconst.GetFval() != nil:
			fval := aconst.GetFval().Fval
			b.WriteString(fmt.Sprintf(" %s", fval))
		case aconst.GetIval() != nil:
			ival := aconst.GetIval().Ival
			b.WriteString(fmt.Sprintf(" %d", ival))
		case aconst.GetSval() != nil:
			sval := aconst.GetSval()
			b.WriteString(fmt.Sprintf(" '%s'", sval))
		default:
			panic("Unknown constraint default type" + fmt.Sprintf("%T", aconst.Val))
		}
	}
	return b.String()
}

func parseFK(c *pg_query.Constraint) string {
	b := new(strings.Builder)
	b.WriteString(c.Pktable.Relname)
	b.WriteString("(")
	b.WriteString(c.PkAttrs[0].GetString_().GetSval())
	b.WriteString(")")
	return b.String()
}

type ConstraintNode = *pg_query.Node

func parseConstraints(constraints []ConstraintNode) []string {
	var cons []string
	for _, c := range constraints {
		switch c.Node.(type) {
		case *pg_query.Node_Constraint:
			switch c.GetConstraint().Contype {
			case pg_query.ConstrType_CONSTR_PRIMARY:
				cons = append(cons, "primary key")
			case pg_query.ConstrType_CONSTR_UNIQUE:
				cons = append(cons, "unique")
			case pg_query.ConstrType_CONSTR_FOREIGN:
				ref := "references " + parseFK(c.GetConstraint())
				cons = append(cons, ref)
			case pg_query.ConstrType_CONSTR_NOTNULL:
				cons = append(cons, "not null")
			case pg_query.ConstrType_CONSTR_DEFAULT:
				cons = append(cons, parseDefault(c.GetConstraint()))
			case pg_query.ConstrType_CONSTR_CHECK:
				cons = append(cons, "check")
			default:

				// Feel free to open a PR to add the missing constraint types
				panic("Unknown constraint type" + fmt.Sprintf("%T", c.Node))
			}

		default:
			panic("Not a constraint" + fmt.Sprintf("%T", c.Node))
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
