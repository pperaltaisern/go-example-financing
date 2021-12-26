package main

import (
	"context"
	"fmt"
	"ledger/pkg/financing"
	"ledger/pkg/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	connectionString := "user=postgres password=postgres host=localhost port=5432 dbname=postgres pool_max_conns=10"
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		panic(err)
	}

	repo := postgres.NewInvestorRepository(pool)

	investor := financing.NewInvestor("test_investor")
	investor.AddFunds(100)

	err = repo.Add(context.Background(), investor)
	if err != nil {
		panic(err)
	}

	ret, err := repo.ByID(context.Background(), "test_investor")
	if err != nil {
		panic(err)
	}
	fmt.Println("obtained: ", ret)
}
