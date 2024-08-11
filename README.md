# squeal

## What

Squeal is going to be an opinionated postgresql linter.

## Why

There is currently no good opinionated SQL linter.
I also fancy a golang project, and I want to try using codecov on something I'm making.

## How

1. Parse SQL using [PG Query](github.com/pganalyze/pg_query_go)
2. Apply formatting rules
   - Hard-code for now. Maybe introduce the notion of rules later on.
3. Output formatted SQL to stdout (write to file later).

## How to run

```bash
go run squeal.go select.go create.go
```

## How to test

Not yet...

## Progress

I may want to completely rewrite what I have done for now. It feels stupidly sloppy and short sighted. I'm really not pleased with it at all right now.
Development should be clear so that anyone can add to it. I'm not sure I can say that about this project right now.

- [x] Parse SQL
- [x] Format `create table` statements
- [ ] Format `select` statements
- [ ] Format `insert` statements
- [ ] Format `update` statements
- [ ] Format `delete` statements
- [ ] Format `alter table` statements
- [ ] Format `drop table` statements
- [ ] Format `drop database` statements
- [ ] Format `create database` statements
- [ ] Format `create index` statements
- [ ] Format `drop index` statements
- [ ] Format `create view` statements
- [ ] Format `drop view` statements
- [ ] Format `create function` statements
- [ ] Format `drop function` statements
- [ ] Format `create trigger` statements
- [ ] Format `drop trigger` statements
- [ ] Format `create schema` statements
- [ ] Format `drop schema` statements
- [ ] Format `create type` statements
- [ ] Format `drop type` statements
- [ ] Format `create extension` statements
- [ ] Format `drop extension` statements
