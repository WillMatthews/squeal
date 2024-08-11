package main

import (
	"fmt"
	"os"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	in, err := os.ReadFile("./test-sql/unformatted.sql")
	check(err)

	s, err := parse(string(in))
	check(err)

	// fmt.Println((s))
	deparsePretty(s)

	// out, err := deparse(s)
	// check(err)

	// fmt.Println(out)
}

func parse(s string) (*pg_query.ParseResult, error) {
	result, err := pg_query.Parse(s)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deparsePretty(sql *pg_query.ParseResult) (string, error) {
	for _, statement := range sql.Stmts {
		stmt := statement.GetStmt()
		fmt.Println("\n\n\n=================STATEMENT=================")

		d, err := decideDeparse(stmt)
		check(err)

		_ = d

	}

	return "", nil
}

func decideDeparse(stmt *pg_query.Node) (string, error) {
	switch stmt.Node.(type) {
	case *pg_query.Node_CreateStmt:
		cs, err := deparseCreate(stmt.GetCreateStmt())
		check(err)
		return printCreate(cs), nil
	case *pg_query.Node_SelectStmt:
		return deparseSelect(stmt.GetSelectStmt())
	default:
		panic("Unknown statement type" + fmt.Sprintf("%T", stmt.Node))
	}
}

type CreateStatement struct {
	Relation string
	Columns  []Column
}

type Column struct {
	Name        string
	Type        string
	Constraints []string
}

func deparseCreate(stmt *pg_query.CreateStmt) (CreateStatement, error) {
	fmt.Println(stmt)

	tbl := CreateStatement{
		Relation: stmt.Relation.Relname,
	}

	for _, elt := range stmt.GetTableElts() {
		switch elt.Node.(type) {
		case *pg_query.Node_ColumnDef:
			col := Column{
				Name:        elt.GetColumnDef().Colname,
				Type:        elt.GetColumnDef().TypeName.String(),
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

func printCreate(cs CreateStatement) string {
	builder := new(strings.Builder)
	builder.WriteString("create table " + cs.Relation + " (\n")
	for _, col := range cs.Columns {
		builder.WriteString(col.Name + " " + col.Type + ",\n")
	}
	builder.WriteString(")")

	fmt.Println(builder.String())

	return builder.String()
}

func deparseSelect(stmt *pg_query.SelectStmt) (string, error) {
	fmt.Println(stmt)

	return "", nil
}

// func deparseDefault(sql *pg_query.ParseResult) (string, error) {
// 	stmt, err := pg_query.Deparse(sql)
// 	if err != nil {
// 		return "", err
// 	}

// 	return stmt, nil
// }
