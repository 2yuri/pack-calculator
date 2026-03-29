package db

import (
	"fmt"

	"github.com/2yuri/pack-calculator/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New() (*sqlx.DB, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Instance().DB.User,
		config.Instance().DB.Password,
		config.Instance().DB.Host,
		config.Instance().DB.Port,
		config.Instance().DB.Name,
	)

	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
