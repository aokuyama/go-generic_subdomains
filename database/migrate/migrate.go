package migrate

import (
	"database/sql"
	"errors"
	"flag"
	"github.com/aokuyama/go-generic_subdomains/database"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrate struct {
	migrate *migrate.Migrate
}

func NewMigrate(d *database.Dsn, dir string) (*Migrate, error) {
	db, err := sql.Open("postgres", d.Dsn())
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(dir, "postgres", driver)
	if err != nil {
		return nil, err
	}

	mm := Migrate{migrate: m}
	return &mm, nil
}

func (m *Migrate) Up() error {
	return m.migrate.Up()
}

func (m *Migrate) Down() error {
	return m.migrate.Down()
}

func (m *Migrate) Force(v int) error {
	return m.migrate.Force(v)
}

func (m *Migrate) DoFlag() {
	err := m.Flag()
	if err != nil {
		panic(err)
	}
}
func (m *Migrate) Flag() error {
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		if args[0] == "up" {
			return m.Up()
		}
		if args[0] == "down" {
			return m.Down()
		}
	} else if len(args) == 2 {
		if args[0] == "force" {
			ver, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			return m.Force(ver)
		}
	}
	return errors.New("up/down/force [int]")
}
