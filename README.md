Personal Web Store - E-Commerce with Midtrans & MySQL.

Personal Web Store is a Go (Golang)-based eCommerce platform that integrates with Midtrans for secure payment processing. It supports product management, shopping cart functionality, and online payments using various methods like Virtual Accounts, E-Wallets, and Credit/Debit Cards.

🚀 Upcoming Update: User authentication (Register & Login) will be added soon!

✨ Features
✅ Product Catalog (List, View, Manage)
✅ Shopping Cart & Checkout
✅ Order Management
✅ Midtrans Payment Gateway Integration
✅ Admin Panel for Managing Orders & Products
✅ Responsive UI with HTML & CSS
🔜 User Authentication (Register & Login) - Coming Soon!

🛠 Tech Stack
Backend: Go (Fiber/Gin Framework)
Database: MySQL
Frontend: HTML, CSS, JavaScript
ORM: GORM
Payment Gateway: Midtrans
🚀 Installation
1️⃣ Clone the Repository
bash
Copy
Edit
git clone https://github.com/agungferdi/golang-web-payment-itegration.git
cd personal-web-store
2️⃣ Install Dependencies
bash
Copy
Edit
go mod tidy
3️⃣ Configure .env File
Create a .env file in the root directory with the following variables:

ini
Copy
Edit
# Database Configuration
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=personal_store
DB_HOST=localhost
DB_PORT=3306

# Midtrans Payment Gateway
MIDTRANS_SERVER_KEY=your-midtrans-server-key
MIDTRANS_CLIENT_KEY=your-midtrans-client-key
MIDTRANS_ENV=sandbox   # Use "sandbox" for testing, "production" for live transactions
🔹 DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT → MySQL Database connection settings.
🔹 MIDTRANS_SERVER_KEY → Your Midtrans private server key (for backend transactions).
🔹 MIDTRANS_CLIENT_KEY → Your Midtrans public client key (for frontend integration).
🔹 MIDTRANS_ENV → Set to "sandbox" for testing and "production" for live payments.

4️⃣ Run Database Migrations
bash
Copy
Edit
go run migrate.go
5️⃣ Start the Server
bash
Copy
Edit
go run main.go
Now, visit: http://localhost:8080 🎉

🔗 API Endpoints
🛍 Products
GET /products → Get all products
GET /products/:id → Get product details
POST /products → Add new product (Admin only)
🛒 Orders & Checkout
POST /checkout → Create an order & generate Midtrans payment link
GET /orders/:id → Retrieve order details
POST /orders/webhook → Handle Midtrans payment status updates
🔐 Upcoming Authentication API
🚀 In the next update, user authentication (Register & Login) will be added, allowing customers to create accounts, log in, and track orders.

💳 Midtrans Payment Integration
This project integrates with Midtrans to support multiple payment methods:

✔ Virtual Accounts (BCA, Mandiri, BNI, etc.)
✔ E-Wallets (GoPay, ShopeePay, OVO)
✔ Credit/Debit Cards (Visa, MasterCard, JCB)

How It Works:
1️⃣ User checks out and selects a payment method.
2️⃣ The backend calls Midtrans API to create a payment transaction.
3️⃣ Midtrans returns a payment URL where the user completes the payment.
4️⃣ Once paid, Midtrans webhook updates the order status in the database.

📸 Screenshots

<img width="1440" alt="Screenshot 2025-02-25 at 14 28 12" src="https://github.com/user-attachments/assets/e8912930-b3ac-4e66-b82a-8d304f8218a3" />
<img width="1440" alt="Screenshot 2025-02-25 at 14 28 33" src="https://github.com/user-attachments/assets/01b190eb-4057-4a32-a66a-04ab4e652b61" />
<img width="1440" alt="Screenshot 2025-02-25 at 14 28 37" src="https://github.com/user-attachments/assets/76fd789d-9945-4662-a3e0-dd2d47429acc" />
<img width="1440" alt="Screenshot 2025-02-25 at 14 28 54" src="https://github.com/user-attachments/assets/f37cee4f-aaf0-4bdf-a31a-90c5ef3c0036" />

