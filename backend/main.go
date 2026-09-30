package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // โหลด Driver สำหรับ MySQL
)

// OrderRequest represents the incoming JSON payload for creating an order
type OrderRequest struct {
	CustomerID    int         `json:"customer_id"`
	EmployeeID    int         `json:"employee_id"`
	ReceivedDate  string      `json:"received_date"`
	TotalAmount   float64     `json:"total_amount"`
	DepositAmount float64     `json:"deposit_amount"`
	Items         []OrderItem `json:"items"`
}

func main() {
	// 1. ตั้งค่า DSN (Data Source Name) สำหรับเชื่อมต่อ MySQL
	dsn := "root:@tcp(127.0.0.1:3306)/bakery_db?parseTime=true"

	// 2. เปิดการเชื่อมต่อ
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("เกิดข้อผิดพลาดในการเปิดฐานข้อมูล: %v", err)
	}
	defer db.Close()

	// 3. ทดสอบการเชื่อมต่อ (Ping) ว่าต่อติดจริงๆ ไหม
	err = db.Ping()
	if err != nil {
		log.Fatalf("ไม่สามารถเชื่อมต่อฐานข้อมูลได้: %v", err)
	}

	fmt.Println("🎉 เชื่อมต่อฐานข้อมูลสำเร็จแล้ว!")

	// 4. เสิร์ฟไฟล์ Static จากโฟลเดอร์ frontend
	frontendDir := "../frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "frontend"
	}
	fs := http.FileServer(http.Dir(frontendDir))
	http.Handle("/", fs)

	// 5. ลงทะเบียน API Endpoints ทั้งหมด
	http.HandleFunc("/api/register", RegisterHandler(db))
	http.HandleFunc("/api/login", LoginHandler(db))
	http.HandleFunc("/api/payment", PaymentHandler(db))
	http.HandleFunc("/api/payment/balance", BalancePaymentHandler(db))
	http.HandleFunc("/api/admin/orders", AdminOrdersHandler(db))
	http.HandleFunc("/api/admin/update-status", AdminUpdateStatusHandler(db))
	http.HandleFunc("/api/admin/ingredients", IngredientsHandler(db))
	http.HandleFunc("/api/admin/production-plan", ProductionPlanHandler(db))
	http.HandleFunc("/api/admin/daily-report", DailyReportHandler(db))
	http.HandleFunc("/api/customer/orders", CustomerOrdersHandler(db))

	// API Endpoint สำหรับสร้างออเดอร์ (/api/order)
	http.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
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

		var req OrderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		// เรียกใช้ฟังก์ชัน CreateOrder จาก order_service.go
		orderID, err := CreateOrder(db, req.CustomerID, req.EmployeeID, req.ReceivedDate, req.TotalAmount, req.DepositAmount, req.Items)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// ส่งผลลัพธ์สำเร็จกลับไป
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"order_id": orderID,
			"message":  "สร้างคำสั่งซื้อและตรวจสอบโควต้าสำเร็จ",
		})
	})

	fmt.Println("🚀 Server is running on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
