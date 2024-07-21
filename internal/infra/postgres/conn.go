package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type Config struct {
	Host string
	Port string
	User string
	Pswd string
	DBNm string
}

type PgInstance struct {
	config *Config
	ctx    context.Context
	conn   *pgx.Conn
}

var singleInstance *PgInstance
var once = &sync.Once{}

// initMigration initialize migrations, connect with table if tbName is different than empty
// passing the folder path containing the migrations
func (i *PgInstance) initMigration(tbName, filePath string) error {
	path := ""
	if tbName != "" {
		path = fmt.Sprintf("/%s", i.config.DBNm)
	}
	dburl := url.URL{
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", i.config.Host, i.config.Port),
		User:     url.UserPassword(i.config.User, i.config.Pswd),
		Path:     path,
		RawQuery: "sslmode=disable",
	}
	// open coon
	db, err := sql.Open("postgres", dburl.String())
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return err
	}
	// migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	wd, _ := os.Getwd()
	// parse to work on windows
	u, _ := url.Parse(wd)
	// create db path
	dbPath := fmt.Sprintf("file:%s/%s", u, filePath)
	// create instance
	m, err := migrate.NewWithDatabaseInstance(dbPath, "postgres", driver)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	// TODO: Melhorar
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	// success
	return nil
}

// SetConfiguration set configuration
func SetConfiguration(config Config) {
	if singleInstance == nil {
		once.Do(func() {
			singleInstance = &PgInstance{
				config: &config,
				ctx:    context.Background(),
			}
		})
	}
}

func (i *PgInstance) Init(withMigration bool) error {
	// create db url, example := "postgres://username:password@localhost:5432/database_name"
	dburl := url.URL{
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", i.config.Host, i.config.Port),
		User:     url.UserPassword(i.config.User, i.config.Pswd),
		Path:     fmt.Sprintf("/%s", i.config.DBNm),
		RawQuery: "sslmode=disable",
	}
	// // create database
	// err := i.initMigration("", "migrations/database")
	// if err != nil {
	// 	return err
	// }
	// do migrations
	if withMigration {
		err := i.initMigration(i.config.DBNm, "migrations")
		if err != nil {
			return err
		}
	}
	// connect
	conn, err := pgx.Connect(i.ctx, dburl.String())
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	// ping to test
	err = conn.Ping(i.ctx)
	if err != nil {
		return err
	}
	i.conn = conn
	return nil
}

// GetConn recover conn from instance
func (i *PgInstance) GetConn() (pgx.Tx, error) {
	// case not exist, initialize
	if i.conn == nil {
		err := i.Init(false)
		if err != nil {
			return nil, err
		}
	}
	tx, err := i.conn.BeginTx(i.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	// return conn
	return tx, nil
}

// GetInstance get a singleton instance of conn
func GetInstance() *PgInstance {
	return singleInstance
}
