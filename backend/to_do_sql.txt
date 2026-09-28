ช่วยเขียนโค้ดสำหรับเริ่มต้นทำ Backend (Golang) จากโครงสร้าง DBML นี้ให้หน่อยครับ:

1. ขอคำสั่ง SQL (CREATE TABLE) สำหรับสร้างฐานข้อมูล (อิงตาม MySQL)
2. ขอโค้ด Go Struct สำหรับทุกตาราง พร้อมใส่ Tags สำหรับ JSON และถ้าใช้ ORM (เช่น GORM) ให้ใส่ Tags มาให้ด้วย

โครงสร้าง Database:
Table CUSTOMER {
  CustomerID int [pk]
  Username varchar
  Password varchar
  Title varchar
  Name varchar
  Surname varchar
  Phone varchar
  Email varchar
}

Table EMPLOYEE {
  EmployeeID int [pk]
  Username varchar
  Password varchar
  Role varchar
}

Table ORDERS {
  OrderID int [pk]
  OrderDate datetime
  ReceivedDate datetime
  TotalAmount decimal
  DepositAmount decimal
  Status varchar
  CustomerID int
  EmployeeID int
}

Table PAYMENT {
  PaymentID int [pk]
  Amount decimal
  PaymentType varchar
  SlipImage varchar
  OrderID int
}

Table PRODUCT {
  ProductID int [pk]
  ProductName varchar
  Price decimal
}

Table INGREDIENT {
  IngredientID int [pk]
  IngredientName varchar
  StockQty decimal
}

Table ORDER_PRODUCT {
  OrderID int [pk]
  ProductID int [pk]
  Quantity int
  SubTotal decimal
}

Table PRODUCT_INGREDIENT {
  ProductID int [pk]
  IngredientID int [pk]
}