package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/urfave/cli/v2"
)

type connectionString struct {
	hostname string
	database string
	port     uint
	username string
	password string
}

type Order struct {
	Id      uint      `db:"id"`
	Content string    `db:"content"`
	Created time.Time `db:"created"`
	AgentId uint      `db:"agent_id"`
}

func (o Order) String() string {
	return fmt.Sprintf("id: %d, content: %s, created: %s, agentId: %d", o.Id, o.Content, o.Created, o.AgentId)
}

func newConnectionString(hostname string, database string, port uint, username string, password string) *connectionString {
	conn := connectionString{
		hostname: hostname,
		database: database,
		port:     port,
		username: username,
		password: password,
	}

	return &conn
}

var ordersCount uint = 0

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hostname",
				Value: "localhost",
				Usage: "message before counter",
			},
			&cli.StringFlag{
				Name:  "database",
				Value: "postgres",
				Usage: "name of the database",
			},
			&cli.StringFlag{
				Name:  "username",
				Value: "admin",
				Usage: "name of user",
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "poney",
				Usage: "password (use env PG_PASSWORD)",
			},
			&cli.UintFlag{
				Name:  "port",
				Value: 5432,
				Usage: "instance port",
			},
			// &cli.DurationFlag{
			// 	Name:  "frequency",
			// 	Value: time.Duration(1 * time.Second),
			// 	Usage: "elasped time between message",
			// },
		},
		Name:  "load-test-rds",
		Usage: "Connect to postgres",
		Action: func(ctx *cli.Context) error {
			hostname := ctx.String("hostname")
			port := ctx.Uint("port")
			database := ctx.String("database")
			username := ctx.String("username")
			password := ctx.String("password")

			connStr := newConnectionString(hostname, database, port, username, password)
			db, err := open(connStr.build())
			if err != nil {
				return err
			}
			defer db.Close(context.Background())

			if err = nextSQLCommand(db, 1, countOrders, createOrder, listOrders); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func (connection *connectionString) build() string {
	// "postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]"

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		connection.username,
		connection.password,
		connection.hostname,
		connection.port,
		connection.database,
	)
}

func open(connection string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return db, nil
}

func listOrders(db *pgx.Conn, agentId uint) error {
	rows, err := db.Query(context.Background(), "SELECT id, content, created, agent_id FROM orders WHERE agent_id =$1 ORDER BY id", agentId)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[Order])

	if err != nil {
		return fmt.Errorf("collectRows failed: %w", err)
	}

	for _, o := range orders {
		fmt.Fprintf(os.Stdout, "%s\n", o.String())
	}

	return nil
}

func createOrder(db *pgx.Conn, agentId uint) error {
	content := fmt.Sprintf("Order %d from agent %d", ordersCount+1, agentId)

	tag, err := db.Exec(context.Background(), "INSERT INTO Orders (content, agent_id) VALUES($1, $2)", content, agentId)
	if err != nil {
		return fmt.Errorf("insert into orders failed: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return errors.New("insert into orders failed: no new order have been inserted")
	}
	ordersCount++

	return nil
}

func countOrders(db *pgx.Conn, agentId uint) error {
	var count uint

	err := db.QueryRow(context.Background(), "select COUNT(id) FROM orders WHERE agent_id =$1", agentId).Scan(&count)
	if err != nil {
		return fmt.Errorf("countOrders failed: %w", err)
	}

	ordersCount = count

	return nil
}

func nextSQLCommand(db *pgx.Conn, agentId uint, commands ...func(*pgx.Conn, uint) error) error {
	for _, f := range commands {
		var err error
		if err = f(db, agentId); err != nil {
			return err
		}
	}
	return nil
}
