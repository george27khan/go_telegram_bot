package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func connect(ctx context.Context) *pgx.Conn {
	// loads DB settings from .env into the system
	if err := godotenv.Load("db.env"); err != nil {
		log.Print("No .env file found")
	}
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")

	//"postgres://username:password@localhost:5432/database_name"
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbUser, dbPwd, dbName)

	conn, err := pgx.Connect(, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Connect to DB %s", dbName)
	}
	return conn
}
func Init() {
	ctx := context.Background()
	conn := connect(ctx)
	defer conn.Close(ctx)
	//conn.
	//
	//var name string
	//var weight int64
	//err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}
}
