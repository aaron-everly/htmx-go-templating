package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type Poll struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Published   bool
	Description string
	Votes       int
}

func main() {
	var port = ":8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var polls []Poll
	rows, err := conn.Query(context.Background(), `SELECT id, title, published, "createdAt", "updatedAt", votes, description FROM "Poll"`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var poll Poll
		if err := rows.Scan(&poll.ID, &poll.Title, &poll.Published, &poll.CreatedAt, &poll.UpdatedAt, &poll.Votes, &poll.Description); err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			os.Exit(1)
		}
		polls = append(polls, poll)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error iterating through rows: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected!")
	fmt.Println("Polls:", polls)

	pageData := struct {
		Polls []Poll
	}{
		Polls: polls,
	}

	modalHandler := func(w http.ResponseWriter, r *http.Request) {
		// Assuming modal.html is the template that contains your modal's HTML
		templ := template.Must(template.ParseFiles("modal.html"))
		if err := templ.Execute(w, nil); // Pass any necessary data your modal might need
		err != nil {
			log.Printf("Error executing modal template: %v", err)
			http.Error(w, "Error generating modal content", http.StatusInternalServerError)
		}
	}

	serveHandler := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))
		if err := templ.Execute(w, pageData); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Error generating page", http.StatusInternalServerError)
		}
	}

	http.HandleFunc("/", serveHandler)
	http.HandleFunc("/modal", modalHandler)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
