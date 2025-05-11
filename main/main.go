package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/knblkxx/golangcalcV3/calculator"
	"github.com/knblkxx/golangcalcV3/handlers"
	"github.com/knblkxx/golangcalcV3/auth"
	"log"
	"net/http"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"context"
	"time"
)

func main() {
	calculator.Initial()
	handlers.Initial()
	
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	conn, err := sql.Open("sqlite3", "database/database.sql")
	if err != nil {
		os.Create("database/database.sql")
		conn, err = sql.Open("sqlite3", "database/database.sql")
	}
	defer conn.Close()
	

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	schemaSql := `
	CREATE TABLE IF NOT EXISTS auth (
    login    TEXT UNIQUE
                  PRIMARY KEY
                  NOT NULL,
    password TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS expressions (
    id     INTEGER PRIMARY KEY
                   UNIQUE
                   NOT NULL,
    status TEXT    NOT NULL
                   DEFAULT PENDING,
    result NUMERIC NOT NULL
);`
	_, err = conn.Exec(schemaSql)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if _, err := conn.ExecContext(ctx, schemaSql); err != nil {
        log.Fatalf("failed to init schema: %v", err)
    }

	auth.InitAuth(conn)

	log.Println("Starting Server")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", auth.LoginHandler).Methods("POST")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(auth.JwtMiddleware)
	api.HandleFunc("/calculate", handlers.CalculatorHandler)
	api.HandleFunc("/expressions/{id}", handlers.ExpressionAnswer)
	api.HandleFunc("/expressions", handlers.ExpressionsList)

	r.HandleFunc("/internal/task", handlers.OrkestratorHandler)

	http.ListenAndServe(":8080", r)
}
