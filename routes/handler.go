package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"coba/config"
	"coba/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// IndexHandler untuk index.html
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading index.html", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// BuyPageHandler untuk buy_page.html
func BuyPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving buy_page.html...")

	productName := r.URL.Query().Get("name")
	productPrice := r.URL.Query().Get("price")

	tmpl, err := template.ParseFiles("templates/buy_page.html")
	if err != nil {
		log.Println("Error loading buy_page.html:", err)
		http.Error(w, "Error loading buy_page.html", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"ProductName":  productName,
		"ProductPrice": productPrice,
	}
	tmpl.Execute(w, data)
}

// CheckoutHandler untuk menangani pembayaran
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received order: %+v", order)

	if order.TotalPrice <= 0 {
		http.Error(w, "Total price must be greater than zero", http.StatusBadRequest)
		return
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  generateOrderID(),
			GrossAmt: int64(order.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.FullName,
			Phone: order.Phone,
			BillAddr: &midtrans.CustomerAddress{
				Address:  order.Address,
				Postcode: order.PostalCode,
			},
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "1",
				Name:  order.Product,
				Price: int64(order.UnitPrice),
				Qty:   int32(order.Quantity),
			},
			{
				ID:    "shipping",
				Name:  "Shipping Cost",
				Price: int64(order.ShippingCost),
				Qty:   1,
			},
			{
				ID:    "insurance",
				Name:  "Insurance",
				Price: int64(order.InsuranceCost),
				Qty:   1,
			},
		},
	}

	snapResp, err := config.MidtransClient.CreateTransaction(snapReq)
	if err != nil {
		log.Println("Error creating transaction:", err)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	log.Printf("Snap Response: %+v", snapResp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": snapResp.Token,
	})
}

// generateOrderID membuat Order ID unik
func generateOrderID() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

// SuccessHandler menangani halaman sukses
func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	fmt.Fprintf(w, "Payment successful! Order ID: %s", orderID)
}

// PendingHandler menangani halaman pending
func PendingHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	fmt.Fprintf(w, "Payment pending! Order ID: %s", orderID)
}
