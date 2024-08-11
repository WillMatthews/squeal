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

	s, err := makeSqlTree(string(in))
	check(err)

	// fmt.Println((s))
	out, err := parsePretty(s)
	check(err)
	fmt.Println(out)
	_ = out
}

func makeSqlTree(s string) (*pg_query.ParseResult, error) {
	result, err := pg_query.Parse(s)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parsePretty(sql *pg_query.ParseResult) (string, error) {
	sb := new(strings.Builder)
	for _, statement := range sql.Stmts {
		stmt := statement.GetStmt()
		fmt.Printf("\n---------------%T---------------\n", stmt.Node)
		d, err := decideDeparse(stmt)
		check(err)

		sb.WriteString(d)
	}

	return sb.String(), nil
}

func decideDeparse(stmt *pg_query.Node) (string, error) {
	switch stmt.Node.(type) {
	case *pg_query.Node_CreateStmt:
		cs, err := parseCreate(stmt.GetCreateStmt())
		check(err)
		return printCreate(cs), nil
	case *pg_query.Node_SelectStmt:
		sel, err := parseSelect(stmt.GetSelectStmt())
		check(err)

		return printSelect(sel), nil
	default:
		panic("Unknown statement type" + fmt.Sprintf("%T", stmt.Node))
	}
}
