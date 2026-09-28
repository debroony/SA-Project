package main

import (
	"time"
)

// Customer struct for CUSTOMER table
type Customer struct {
	CustomerID int    `gorm:"primaryKey;column:CustomerID" json:"customer_id"`
	Username   string `gorm:"column:Username" json:"username"`
	Password   string `gorm:"column:Password" json:"password"`
	Title      string `gorm:"column:Title" json:"title"`
	Name       string `gorm:"column:Name" json:"name"`
	Surname    string `gorm:"column:Surname" json:"surname"`
	Phone      string `gorm:"column:Phone" json:"phone"`
	Email      string `gorm:"column:Email" json:"email"`
}

func (Customer) TableName() string {
	return "CUSTOMER"
}

// Employee struct for EMPLOYEE table
type Employee struct {
	EmployeeID int    `gorm:"primaryKey;column:EmployeeID" json:"employee_id"`
	Username   string `gorm:"column:Username" json:"username"`
	Password   string `gorm:"column:Password" json:"password"`
	Role       string `gorm:"column:Role" json:"role"`
}

func (Employee) TableName() string {
	return "EMPLOYEE"
}

// Order struct for ORDERS table
type Order struct {
	OrderID       int       `gorm:"primaryKey;column:OrderID" json:"order_id"`
	OrderDate     time.Time `gorm:"column:OrderDate" json:"order_date"`
	ReceivedDate  time.Time `gorm:"column:ReceivedDate" json:"received_date"`
	TotalAmount   float64   `gorm:"column:TotalAmount" json:"total_amount"`
	DepositAmount float64   `gorm:"column:DepositAmount" json:"deposit_amount"`
	Status        string    `gorm:"column:Status" json:"status"`
	CustomerID    int       `gorm:"column:CustomerID" json:"customer_id"`
	EmployeeID    int       `gorm:"column:EmployeeID" json:"employee_id"`

	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Employee *Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

func (Order) TableName() string {
	return "ORDERS"
}

// Payment struct for PAYMENT table
type Payment struct {
	PaymentID   int     `gorm:"primaryKey;column:PaymentID" json:"payment_id"`
	Amount      float64 `gorm:"column:Amount" json:"amount"`
	PaymentType string  `gorm:"column:PaymentType" json:"payment_type"`
	SlipImage   string  `gorm:"column:SlipImage" json:"slip_image"`
	OrderID     int     `gorm:"column:OrderID" json:"order_id"`

	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (Payment) TableName() string {
	return "PAYMENT"
}

// Product struct for PRODUCT table
type Product struct {
	ProductID   int     `gorm:"primaryKey;column:ProductID" json:"product_id"`
	ProductName string  `gorm:"column:ProductName" json:"product_name"`
	Price       float64 `gorm:"column:Price" json:"price"`
}

func (Product) TableName() string {
	return "PRODUCT"
}

// Ingredient struct for INGREDIENT table
type Ingredient struct {
	IngredientID   int     `gorm:"primaryKey;column:IngredientID" json:"ingredient_id"`
	IngredientName string  `gorm:"column:IngredientName" json:"ingredient_name"`
	StockQty       float64 `gorm:"column:StockQty" json:"stock_qty"`
}

func (Ingredient) TableName() string {
	return "INGREDIENT"
}

// OrderProduct struct for ORDER_PRODUCT table (Many-to-Many relationship table)
type OrderProduct struct {
	OrderID   int     `gorm:"primaryKey;column:OrderID" json:"order_id"`
	ProductID int     `gorm:"primaryKey;column:ProductID" json:"product_id"`
	Quantity  int     `gorm:"column:Quantity" json:"quantity"`
	SubTotal  float64 `gorm:"column:SubTotal" json:"sub_total"`

	Order   *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (OrderProduct) TableName() string {
	return "ORDER_PRODUCT"
}

// ProductIngredient struct for PRODUCT_INGREDIENT table (Many-to-Many relationship table)
type ProductIngredient struct {
	ProductID    int `gorm:"primaryKey;column:ProductID" json:"product_id"`
	IngredientID int `gorm:"primaryKey;column:IngredientID" json:"ingredient_id"`

	Product    *Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient,omitempty"`
}

func (ProductIngredient) TableName() string {
	return "PRODUCT_INGREDIENT"
}
