package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"coba/config"
	"coba/routes"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv" // Import godotenv
)

func main() {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY is not set")
	}
	fmt.Println("MIDTRANS_SERVER_KEY is set:", serverKey)

	config.InitMidtrans(serverKey)

	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/coba")
	if err != nil {
		log.Fatal("MySQL connection failed:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("MySQL ping failed:", err)
	}
	fmt.Println("Connected to MySQL")

	routes.DB = db

	http.HandleFunc("/", routes.IndexHandler)
	http.HandleFunc("/buy_page.html", routes.BuyPageHandler)
	http.HandleFunc("/checkout", routes.CheckoutHandler)
	http.HandleFunc("/success", routes.SuccessHandler)
	http.HandleFunc("/pending", routes.PendingHandler)
	http.HandleFunc("/webhook", routes.WebhookHandler)
	http.HandleFunc("/transactions", routes.TransactionListHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
