package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"coba/config"
	"coba/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var DB *sql.DB

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println("IndexHandler error:", err)
		return
	}
	tmpl.Execute(w, nil)
}

func BuyPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/buy_page.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println("BuyPageHandler error:", err)
		return
	}

	data := map[string]string{
		"ProductName":  r.URL.Query().Get("name"),
		"ProductPrice": r.URL.Query().Get("price"),
	}
	tmpl.Execute(w, data)
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("CheckoutHandler decode error:", err)
		return
	}

	order.ID = generateOrderID()
	order.MidtransOrderID = order.ID
	order.CreatedAt = time.Now()
	order.PaymentStatus = "pending"

	_, err := DB.Exec(`
		INSERT INTO orders (
			id, product, quantity, unit_price, 
			shipping_cost, insurance_cost, total_price, 
			full_name, address, postal_code, phone, 
			payment_status, created_at, midtrans_order_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		order.ID,
		order.Product,
		order.Quantity,
		order.UnitPrice,
		order.ShippingCost,
		order.InsuranceCost,
		order.TotalPrice,
		order.FullName,
		order.Address,
		order.PostalCode,
		order.Phone,
		order.PaymentStatus,
		order.CreatedAt,
		order.MidtransOrderID,
	)

	if err != nil {
		log.Printf("DB insert error: %v", err)
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	// Create Midtrans request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.MidtransOrderID,
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
			{ID: "1", Name: order.Product, Price: int64(order.UnitPrice), Qty: int32(order.Quantity)},
			{ID: "shipping", Name: "Shipping Cost", Price: int64(order.ShippingCost), Qty: 1},
			{ID: "insurance", Name: "Insurance", Price: int64(order.InsuranceCost), Qty: 1},
		},
	}

	snapResp, err := config.MidtransClient.CreateTransaction(snapReq)

	if err != nil {

		if snapResp != nil && snapResp.Token != "" {
			log.Printf("Midtrans transaction succeeded with warnings: %v", err)
		} else {
			log.Printf("Midtrans transaction failed: %v", err)
			http.Error(w, "Payment gateway error", http.StatusInternalServerError)
			return
		}
	}

	if snapResp == nil || snapResp.Token == "" {
		log.Printf("Midtrans transaction failed: empty response")
		http.Error(w, "Payment gateway error", http.StatusInternalServerError)
		return
	}

	log.Printf("Payment initiated: OrderID=%s, Token=%s", order.ID, snapResp.Token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":    snapResp.Token,
		"order_id": order.ID,
	})
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var notification map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid notification", http.StatusBadRequest)
		return
	}

	orderID, ok := notification["order_id"].(string)
	if !ok {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	status, ok := notification["transaction_status"].(string)
	if !ok {
		http.Error(w, "Invalid transaction status", http.StatusBadRequest)
		return
	}

	var paymentStatus string
	switch status {
	case "capture", "settlement":
		paymentStatus = "settlement"
	case "pending":
		paymentStatus = "pending"
	case "expire":
		paymentStatus = "expire"
	case "deny", "cancel":
		paymentStatus = "failed"
	default:
		paymentStatus = "pending"
	}

	_, err := DB.Exec(
		"UPDATE orders SET payment_status = ? WHERE midtrans_order_id = ?",
		paymentStatus,
		orderID,
	)

	if err != nil {
		log.Printf("Update error: %v", err)
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	_, err := DB.Exec("UPDATE orders SET payment_status = 'settlement' WHERE id = ?", orderID)
	if err != nil {
		log.Printf("DB update error: %v", err)
	}
	fmt.Fprintf(w, "Payment successful! Order ID: %s", orderID)
}

func PendingHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	fmt.Fprintf(w, "Payment pending! Order ID: %s", orderID)
}

func generateOrderID() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

// TransactionListHandler shows the transaction list
func TransactionListHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT 
			created_at,
			midtrans_order_id,
			'Payment' AS transaction_type,
			'Bank Transfer' AS channel,
			payment_status,
			total_price
		FROM orders
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []map[string]interface{}
	for rows.Next() {
		var (
			createdAt  time.Time
			orderID    string
			txType     string
			channel    string
			status     string
			totalPrice int
		)

		err := rows.Scan(&createdAt, &orderID, &txType, &channel, &status, &totalPrice)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		transactions = append(transactions, map[string]interface{}{
			"date":           createdAt.Format("02 Jan 2006, 15:04"),
			"order_id":       orderID,
			"type":           txType,
			"channel":        channel,
			"status":         status,
			"amount":         formatRupiah(totalPrice),
			"customer_email": "-",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func formatRupiah(amount int) string {
	return fmt.Sprintf("Rp%s", formatNumber(amount))
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)
	startOffset := 0
	if n < 0 {
		startOffset = 1
	}

	var result []rune
	for i, c := range s[startOffset:] {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, c)
	}

	if n < 0 {
		return "-" + string(result)
	}
	return string(result)
}
