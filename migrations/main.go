package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:example@127.0.0.1:6543/postgres?sslmode=disable")
	if err != nil {
		logrus.Fatal("open connection: ", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.Fatal("create instance: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://",
		"postgres", driver)
	if err != nil {
		logrus.Fatal("create migration object: ", err)
	}

	isDrop := flag.Bool("drop", false, "use flag -drop to drop DB")
	flag.Parse()

	if *isDrop {
		if err = m.Drop(); err != nil {
			logrus.Fatal(fmt.Errorf("migration down: %w", err))
		}
	} else {
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logrus.Println("no change")
				os.Exit(0)
			}
			logrus.Fatal(fmt.Errorf("migration up: %w", err))
		}
	}
	logrus.Println("migration successful")
}
