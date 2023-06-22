package clients

import (
	//"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cmartinc22/devhub-backend-core/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pedidosya/peya-go/logs"
)

type PostgressControllerSpec interface {
	Ping() error
}

type PostgressClient struct {
	Config               *models.DbConfiguration
	Db *sqlx.DB
}

// NewClient creates an STS client that contains the introspection service.
func NewPostgressClient() (PostgressControllerSpec, error) {
	cfg := readDBConfig()
	
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Server, cfg.Port, cfg.Database)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db = db.Unsafe()

	r := &PostgressClient{Config: cfg, Db: db}

	r.Db.SetMaxOpenConns(10)   // The default is 0 (unlimited)
	r.Db.SetMaxIdleConns(2)    // defaultMaxIdleConns = 2
	r.Db.SetConnMaxLifetime(0) // 0, connections are reused forever.

	if err = r.Db.Ping(); err != nil {
		return nil, err
	}

	return r, nil
}

func readDBConfig() (*models.DbConfiguration) {
	// Get Configuration FROM ENV
	cfg := &models.DbConfiguration{}

	//POSTGRES DB configurations
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
	return cfg
}

func (db *PostgressClient) Ping() error {
	if db.Db != nil {
		return db.Db.Ping()
	}
	return errors.New("missing db conection")
}