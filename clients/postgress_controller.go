package clients

import (
	//"context"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cmartinc22/devhub-backend-core/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pedidosya/peya-go/logs"
)

const (
	defauiltConnString = "postgres://%s:%s@%s:%v/%s?sslmode=disable"
)
type PostgressControllerSpec interface {
	Ping(context.Context) error
	GetConnectionString(context.Context) string
	Db() *sqlx.DB
}

type PostgressClient struct {
	Config *models.DbConfiguration
	Database     *sqlx.DB
}

// NewClient creates an STS client that contains the introspection service.
func NewPostgressClient() (PostgressControllerSpec, error) {
	cfg := readDBConfig()

	db, err := sqlx.Open("postgres", formatConnectinoString(*cfg))
	if err != nil {
		return nil, err
	}
	db = db.Unsafe()

	r := &PostgressClient{Config: cfg, Database: db}

	r.Database.SetMaxOpenConns(10)   // The default is 0 (unlimited)
	r.Database.SetMaxIdleConns(2)    // defaultMaxIdleConns = 2
	r.Database.SetConnMaxLifetime(0) // 0, connections are reused forever.

	if err = r.Database.Ping(); err != nil {
		return nil, err
	}

	return r, nil
}

func readDBConfig() *models.DbConfiguration {
	// Get Configuration FROM ENV
	cfg := &models.DbConfiguration{}

	// POSTGRES DB configurations
	if v, ok := os.LookupEnv("POSTGRES_HOST"); ok {
		cfg.Server = v
	} else {
		logs.Fatal("[main] missing POSTGRES_HOST environment")
	}

	if v, ok := os.LookupEnv("POSTGRES_PORT"); ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			cfg.Port = i
		} else {
			logs.Debug("[main] using default port for postgres")
			cfg.Port = 5432
		}
	} else {
		logs.Debug("[main] missing POSTGRES_PORT environment. Using default port for postgres")
		cfg.Port = 5432
	}

	if v, ok := os.LookupEnv("POSTGRES_DB"); ok {
		if ok {
			cfg.Database = v
		} else {
			logs.Fatal("[main] missing POSTGRES_DB environment")
		}
	}

	if v, ok := os.LookupEnv("POSTGRES_USER"); ok {
		cfg.User = v
	} else {
		logs.Fatal("[main] missing POSTGRES_USER environment")
	}

	if v, ok := os.LookupEnv("POSTGRES_PASSWORD"); ok {
		cfg.Password = v
	} else {
		logs.Fatal("[main] missing POSTGRES_PASSWORD environment or db.pwd setting")
	}

	if v, ok := os.LookupEnv("POSTGRES_CONN_STR"); ok {
		cfg.ConnectionStr = v
	} else {
		cfg.ConnectionStr = defauiltConnString
		logs.Debug("[main] missing POSTGRES_CONN_STR environment. Using Default")
	}


	return cfg
}

func (db *PostgressClient) Ping(context.Context) error {
	if db.Db != nil {
		return db.Database.Ping()
	}
	return errors.New("missing db conection")
}

func formatConnectinoString(cfg models.DbConfiguration) string {
	return fmt.Sprintf(cfg.ConnectionStr, cfg.User, cfg.Password, cfg.Server, cfg.Port, cfg.Database)
}

func (db *PostgressClient) GetConnectionString(context.Context) string {
	return formatConnectinoString(*db.Config)
}

func (db *PostgressClient) Db() *sqlx.DB {
	return db.Database
}