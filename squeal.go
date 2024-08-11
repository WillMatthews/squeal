package main

import (
	"fmt"
	"os"

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
		cs, err := parseCreate(stmt.GetCreateStmt())
		check(err)
		return printCreate(cs), nil
	case *pg_query.Node_SelectStmt:
		return deparseSelect(stmt.GetSelectStmt())
	default:
		panic("Unknown statement type" + fmt.Sprintf("%T", stmt.Node))
	}
}

func deparseSelect(stmt *pg_query.SelectStmt) (string, error) {
	fmt.Println(stmt)

	return "", nil
}
