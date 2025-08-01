package data

import (
	"database/sql"
	"todo-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewTodoRepo)

// Data .
type Data struct {
	db *sql.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.NewHelper(logger).Errorf("failed to close database: %v", err)
		}
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{db: db},
		cleanup,
		nil
}
