package main

import (
	"database/sql"
	"fmt"
	"time"
)

// OrderItem represents a product item in an order
type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	SubTotal  float64 `json:"sub_total"`
}

// normalizeDate converts DD-MM-YYYY or YYYY-MM-DD to YYYY-MM-DD for SQL
func normalizeDate(dateStr string) string {
	if t, err := time.Parse("02-01-2006", dateStr); err == nil {
		return t.Format("2006-01-02")
	}
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t.Format("2006-01-02")
	}
	return dateStr
}

// CheckDailyQuota checks if adding newOrderQty exceeds the 999 pieces quota for the given receiveDate
func CheckDailyQuota(db *sql.DB, receiveDate string, newOrderQty int) (bool, int, error) {
	normDate := normalizeDate(receiveDate)
	// Query to sum the quantity of products ordered for a specific ReceivedDate
	query := `
		SELECT COALESCE(SUM(op.Quantity), 0) 
		FROM ORDER_PRODUCT op
		JOIN ORDERS o ON op.OrderID = o.OrderID
		WHERE DATE(o.ReceivedDate) = DATE(?) AND o.Status != 'Cancelled'
	`

	var currentTotal int
	err := db.QueryRow(query, normDate).Scan(&currentTotal)
	if err != nil && err != sql.ErrNoRows {
		return false, 0, fmt.Errorf("failed to query quota: %v", err)
	}

	maxQuota := 999
	if currentTotal+newOrderQty > maxQuota {
		// Quota exceeded
		return false, currentTotal, nil
	}

	// Within quota
	return true, currentTotal, nil
}

// CreateOrder inserts a new order into ORDERS and its items into ORDER_PRODUCT using a transaction
func CreateOrder(db *sql.DB, customerID int, employeeID int, receiveDate string, totalAmount float64, depositAmount float64, items []OrderItem) (int64, error) {
	normDate := normalizeDate(receiveDate)

	// 1. Calculate total quantity of new order
	totalNewQty := 0
	for _, item := range items {
		totalNewQty += item.Quantity
	}

	// 2. Check quota before creating order
	allowed, currentTotal, err := CheckDailyQuota(db, normDate, totalNewQty)
	if err != nil {
		return 0, err
	}
	if !allowed {
		return 0, fmt.Errorf("quota exceeded: current reserved is %d, requested %d, max is 999", currentTotal, totalNewQty)
	}

	// 3. Start Transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Handle nullable EmployeeID
	var empIDVal interface{}
	if employeeID > 0 {
		empIDVal = employeeID
	} else {
		empIDVal = nil
	}

	// 4. Insert into ORDERS
	orderQuery := `
		INSERT INTO ORDERS (OrderDate, ReceivedDate, TotalAmount, DepositAmount, Status, CustomerID, EmployeeID)
		VALUES (NOW(), ?, ?, ?, 'รอตรวจสอบ', ?, ?)
	`
	result, err := tx.Exec(orderQuery, normDate, totalAmount, depositAmount, customerID, empIDVal)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %v", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	// 5. Insert each item into ORDER_PRODUCT
	itemQuery := `
		INSERT INTO ORDER_PRODUCT (OrderID, ProductID, Quantity, SubTotal)
		VALUES (?, ?, ?, ?)
	`
	for _, item := range items {
		_, err := tx.Exec(itemQuery, orderID, item.ProductID, item.Quantity, item.SubTotal)
		if err != nil {
			return 0, fmt.Errorf("failed to insert order product item: %v", err)
		}
	}

	// 6. Commit Transaction
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return orderID, nil
}
