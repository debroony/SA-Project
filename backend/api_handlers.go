package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RegisterRequest represents payload for registration
type RegisterRequest struct {
	Fullname string `json:"fullname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents payload for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// StatusUpdateRequest represents payload for updating order status by admin
type StatusUpdateRequest struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

// RegisterHandler handles customer registration
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Split fullname into Name and Surname roughly if possible, or store in Name
		nameParts := strings.SplitN(req.Fullname, " ", 2)
		firstName := nameParts[0]
		lastName := ""
		if len(nameParts) > 1 {
			lastName = nameParts[1]
		}

		query := `
			INSERT INTO CUSTOMER (Username, Password, Title, Name, Surname, Phone, Email)
			VALUES (?, ?, 'คุณ', ?, ?, ?, ?)
		`
		_, err := db.Exec(query, req.Username, req.Password, firstName, lastName, req.Phone, req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to register: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "สมัครสมาชิกสำเร็จ",
		})
	}
}

// LoginHandler handles customer and employee login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// 1. Check EMPLOYEE table first (Admin / Staff)
		var empID int
		var empRole string
		empQuery := `SELECT EmployeeID, Role FROM EMPLOYEE WHERE Username = ? AND Password = ?`
		err := db.QueryRow(empQuery, req.Username, req.Password).Scan(&empID, &empRole)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"role":    "admin",
				"id":      empID,
				"message": "Login as Admin successful",
			})
			return
		}

		// 2. Check CUSTOMER table
		var custID int
		var custUsername string
		custQuery := `SELECT CustomerID, Username FROM CUSTOMER WHERE (Username = ? OR Email = ?) AND Password = ?`
		err = db.QueryRow(custQuery, req.Username, req.Username, req.Password).Scan(&custID, &custUsername)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success":  true,
				"role":     "customer",
				"id":       custID,
				"username": custUsername,
				"message":  "Login successful",
			})
			return
		}

		// Invalid credentials
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "ชื่อผู้ใช้งานหรือรหัสผ่านไม่ถูกต้อง",
		})
	}
}

// PaymentHandler handles slip upload and payment confirmation
func PaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form (max 10MB)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		orderID := r.FormValue("order_id")
		amount := r.FormValue("amount")
		paymentType := r.FormValue("payment_type")

		file, handler, err := r.FormFile("slip_image")
		var slipFilename string
		if err == nil {
			defer file.Close()
			// Ensure uploads dir exists
			os.MkdirAll("../frontend/uploads", os.ModePerm)
			slipFilename = fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
			dstPath := filepath.Join("../frontend/uploads", slipFilename)

			dst, err := os.Create(dstPath)
			if err != nil {
				http.Error(w, "Failed to save slip image", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, file)
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Insert PAYMENT
		paymentQuery := `INSERT INTO PAYMENT (Amount, PaymentType, SlipImage, OrderID) VALUES (?, ?, ?, ?)`
		_, err = tx.Exec(paymentQuery, amount, paymentType, slipFilename, orderID)
		if err != nil {
			http.Error(w, "Failed to insert payment", http.StatusInternalServerError)
			return
		}

		// Update ORDERS Status to 'รอตรวจสอบ'
		updateOrderQuery := `UPDATE ORDERS SET Status = 'รอตรวจสอบ' WHERE OrderID = ?`
		_, err = tx.Exec(updateOrderQuery, orderID)
		if err != nil {
			http.Error(w, "Failed to update order status", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "อัปโหลดสลิปและบันทึกการชำระเงินสำเร็จ",
		})
	}
}

// AdminOrdersHandler returns all orders for admin dashboard
func AdminOrdersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		query := `
			SELECT o.OrderID, COALESCE(c.Name, 'ลูกค้าทั่วไป'), o.ReceivedDate, o.TotalAmount, o.Status,
			       COALESCE(GROUP_CONCAT(CONCAT(p.ProductName, ' (', op.Quantity, ')') SEPARATOR ', '), '-') AS items_summary
			FROM ORDERS o
			LEFT JOIN CUSTOMER c ON o.CustomerID = c.CustomerID
			LEFT JOIN ORDER_PRODUCT op ON o.OrderID = op.OrderID
			LEFT JOIN PRODUCT p ON op.ProductID = p.ProductID
			GROUP BY o.OrderID, c.Name, o.ReceivedDate, o.TotalAmount, o.Status
			ORDER BY o.OrderID DESC
		`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type AdminOrder struct {
			OrderID      int     `json:"order_id"`
			CustomerName string  `json:"customer_name"`
			ReceivedDate string  `json:"received_date"`
			TotalAmount  float64 `json:"total_amount"`
			Status       string  `json:"status"`
			ItemsSummary string  `json:"items_summary"`
		}

		var orders []AdminOrder
		for rows.Next() {
			var o AdminOrder
			var recDate time.Time
			if err := rows.Scan(&o.OrderID, &o.CustomerName, &recDate, &o.TotalAmount, &o.Status, &o.ItemsSummary); err == nil {
				o.ReceivedDate = recDate.Format("02-01-2006")
				orders = append(orders, o)
			}
		}

		json.NewEncoder(w).Encode(orders)
	}
}

// CustomerOrdersHandler returns orders for a specific customer
func CustomerOrdersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		customerID := r.URL.Query().Get("customer_id")
		if customerID == "" {
			http.Error(w, "Missing customer_id", http.StatusBadRequest)
			return
		}

		query := `
			SELECT o.OrderID, o.ReceivedDate, o.TotalAmount, o.Status,
			       COALESCE(GROUP_CONCAT(CONCAT(p.ProductName, ' (', op.Quantity, ')') SEPARATOR ', '), '-') AS items_summary
			FROM ORDERS o
			LEFT JOIN ORDER_PRODUCT op ON o.OrderID = op.OrderID
			LEFT JOIN PRODUCT p ON op.ProductID = p.ProductID
			WHERE o.CustomerID = ?
			GROUP BY o.OrderID, o.ReceivedDate, o.TotalAmount, o.Status
			ORDER BY o.OrderID DESC
		`
		rows, err := db.Query(query, customerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type CustomerOrder struct {
			OrderID      int     `json:"order_id"`
			ReceivedDate string  `json:"received_date"`
			TotalAmount  float64 `json:"total_amount"`
			Status       string  `json:"status"`
			ItemsSummary string  `json:"items_summary"`
		}

		var orders []CustomerOrder
		for rows.Next() {
			var o CustomerOrder
			var recDate time.Time
			if err := rows.Scan(&o.OrderID, &o.ReceivedDate, &o.TotalAmount, &o.Status, &o.ItemsSummary); err == nil {
				o.ReceivedDate = recDate.Format("02-01-2006")
				orders = append(orders, o)
			}
		}

		json.NewEncoder(w).Encode(orders)
	}
}

// AdminUpdateStatusHandler updates order status
func AdminUpdateStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req StatusUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		query := `UPDATE ORDERS SET Status = ? WHERE OrderID = ?`
		_, err := db.Exec(query, req.Status, req.OrderID)
		if err != nil {
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "อัปเดตสถานะสำเร็จ",
		})
	}
}
