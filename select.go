package main

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type SelectStmt struct {
	From       string
	TargetList []string
	Where      string
	GroupBy    []string
	Having     string
	OrderBy    []string
	Sort       []string
	Limit      string
}

func parseSelect(stmt *pg_query.SelectStmt) (SelectStmt, error) {
	fmt.Println(stmt)

	sel := SelectStmt{}

	fmt.Println("TargetList")
	for _, a := range stmt.TargetList {
		fmt.Println(a)
	}

	fmt.Println("FromClause")
	for _, a := range stmt.FromClause {
		fmt.Println(a)
	}

	fmt.Println("ValuesLists")
	for _, a := range stmt.ValuesLists {
		fmt.Println(a)
	}

	fmt.Println("WhereClause")
	fmt.Println(stmt.WhereClause)

	fmt.Println("GroupClause")
	for _, a := range stmt.GroupClause {
		fmt.Println(a)
	}

	fmt.Println("HavingClause")
	fmt.Println(stmt.HavingClause)

	fmt.Println("SortClause")
	for _, a := range stmt.SortClause {
		fmt.Println(a)
	}

	fmt.Println("LimitOption")
	fmt.Println(stmt.LimitOption)

	return sel, nil
}

func printSelect(ss SelectStmt) string {
	return ""
}
