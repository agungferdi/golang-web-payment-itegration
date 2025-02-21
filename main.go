package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"coba/config"
	"coba/routes"
)

func main() {
	// Load Midtrans Server Key dari environment variable
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY is not set or empty")
	}
	fmt.Println("MIDTRANS_SERVER_KEY is set")

	// Inisialisasi Midtrans
	config.InitMidtrans(serverKey)

	// Setup Routes
	http.HandleFunc("/", routes.IndexHandler)
	http.HandleFunc("/buy_page.html", routes.BuyPageHandler)
	http.HandleFunc("/checkout", routes.CheckoutHandler)
	http.HandleFunc("/success", routes.SuccessHandler)
	http.HandleFunc("/pending", routes.PendingHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
