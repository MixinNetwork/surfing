package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MixinNetwork/surfing/config"
	"github.com/MixinNetwork/surfing/durable"
	"github.com/MixinNetwork/surfing/services"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	service := flag.String("service", "http", "run a service")
	flag.Parse()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DatebaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)
	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	database, err := durable.NewDatabase(context.Background(), dbpool)
	if err != nil {
		log.Panicln(err)
	}

	switch *service {
	case "http":
		err := StartHTTP(database)
		if err != nil {
			log.Println(err)
		}
	default:
		hub := services.NewHub(database)
		err := hub.StartService(*service)
		if err != nil {
			log.Println(err)
		}
	}
}
