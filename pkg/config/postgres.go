package config

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pperaltaisern/financing/internal/esrc/esrcpg"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/spf13/viper"
)

type PostgresConfig struct {
	ConnectionString string
}

func LoadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		ConnectionString: viper.GetString("DB_CONNECTION_STRING"),
	}
}

func (c PostgresConfig) Build() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), c.ConnectionString)
}

func (c PostgresConfig) BuildRepositories() (Repositories, error) {
	pool, err := c.Build()
	if err != nil {
		return Repositories{}, err
	}

	es := esrcpg.NewEventStore(pool)
	repos := Repositories{
		Issuers:   financing.NewIssuerRepository(es),
		Investors: financing.NewInvestorRepository(es),
		Invoices:  financing.NewInvoiceRepository(es),
	}
	return repos, nil
}
