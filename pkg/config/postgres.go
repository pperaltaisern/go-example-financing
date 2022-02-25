package config

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/esrcpg"
	"github.com/pperaltaisern/financing/pkg/financing"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func (c PostgresConfig) BuildGORM() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.ConnectionString), &gorm.Config{})
}

func (c PostgresConfig) BuildRepositories() (Repositories, esrc.EventStore, error) {
	pool, err := c.Build()
	if err != nil {
		return Repositories{}, nil, err
	}

	es := esrcpg.NewEventStore(pool)
	repos := Repositories{
		Issuers:   financing.NewIssuerRepository(es),
		Investors: financing.NewInvestorRepository(es),
		Invoices:  financing.NewInvoiceRepository(es),
	}
	return repos, es, nil
}
