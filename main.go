package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	dsn := os.Getenv("POSTGRES_DSN")
	slog.Info("connectiong to postgres server", slog.String("url", dsn))
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		panic(err)
	}

	const dropQuery = `DROP TABLE IF EXISTS t CASCADE;`
	slog.Info("executing table creation query", slog.String("query", dropQuery))
	_, err = conn.Exec(ctx, dropQuery)
	if err != nil {
		panic(err)
	}

	const createQuery = `CREATE TABLE public.t (id serial primary key, field1 text null, field2 int not null, field3 jsonb not null)`
	slog.Info("executing table creation query", slog.String("query", createQuery))
	_, err = conn.Exec(ctx, createQuery)
	if err != nil {
		panic(err)
	}

	slog.Info("starting `copy from`")
	insertedRows, err := conn.CopyFrom(ctx, pgx.Identifier{"public", "t"}, []string{"field1", "field2", "field3"}, pgx.CopyFromRows([][]any{
		{
			"test text",
			42,
			[]byte(`{"some": "json"}`),
		},
		{
			nil,
			7,
			[]byte(`null`),
		},
	}))
	if err != nil {
		panic(err)
	}
	slog.Info("finished `copy from`", slog.Int64("insertedRows", insertedRows))
}
