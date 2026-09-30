-- MySQL CREATE TABLE Statements based on DBML

CREATE TABLE CUSTOMER (
  CustomerID INT AUTO_INCREMENT PRIMARY KEY,
  Username VARCHAR(255) NOT NULL,
  Password VARCHAR(255) NOT NULL,
  Title VARCHAR(50),
  Name VARCHAR(255) NOT NULL,
  Surname VARCHAR(255) NOT NULL,
  Phone VARCHAR(50),
  Email VARCHAR(255)
);

CREATE TABLE EMPLOYEE (
  EmployeeID INT AUTO_INCREMENT PRIMARY KEY,
  Username VARCHAR(255) NOT NULL,
  Password VARCHAR(255) NOT NULL,
  Role VARCHAR(50) NOT NULL
);

CREATE TABLE ORDERS (
  OrderID INT AUTO_INCREMENT PRIMARY KEY,
  OrderDate DATETIME,
  ReceivedDate DATETIME,
  TotalAmount DECIMAL(10, 2),
  DepositAmount DECIMAL(10, 2),
  Status VARCHAR(50),
  CustomerID INT,
  EmployeeID INT,
  FOREIGN KEY (CustomerID) REFERENCES CUSTOMER(CustomerID),
  FOREIGN KEY (EmployeeID) REFERENCES EMPLOYEE(EmployeeID)
);

CREATE TABLE PAYMENT (
  PaymentID INT AUTO_INCREMENT PRIMARY KEY,
  Amount DECIMAL(10, 2),
  PaymentType VARCHAR(50),
  SlipImage VARCHAR(255),
  OrderID INT,
  FOREIGN KEY (OrderID) REFERENCES ORDERS(OrderID)
);

CREATE TABLE PRODUCT (
  ProductID INT AUTO_INCREMENT PRIMARY KEY,
  ProductName VARCHAR(255) NOT NULL,
  Price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE INGREDIENT (
  IngredientID INT AUTO_INCREMENT PRIMARY KEY,
  IngredientName VARCHAR(255) NOT NULL,
  StockQty DECIMAL(10, 2) NOT NULL
);

CREATE TABLE ORDER_PRODUCT (
  OrderID INT,
  ProductID INT,
  Quantity INT,
  SubTotal DECIMAL(10, 2),
  PRIMARY KEY (OrderID, ProductID),
  FOREIGN KEY (OrderID) REFERENCES ORDERS(OrderID),
  FOREIGN KEY (ProductID) REFERENCES PRODUCT(ProductID)
);

CREATE TABLE PRODUCT_INGREDIENT (
  ProductID INT,
  IngredientID INT,
  PRIMARY KEY (ProductID, IngredientID),
  FOREIGN KEY (ProductID) REFERENCES PRODUCT(ProductID),
  FOREIGN KEY (IngredientID) REFERENCES INGREDIENT(IngredientID)
);

-- ==========================================
-- SEED DATA (Initial Data for Testing)
-- ==========================================

-- Default Admin Employee
INSERT INTO EMPLOYEE (Username, Password, Role) VALUES 
('admin', 'admin123', 'admin');

-- Traditional Thai Desserts (matching order.html menu items)
INSERT INTO PRODUCT (ProductID, ProductName, Price) VALUES
(1, 'ทองหยิบ', 50.00),
(2, 'ทองหยอด', 45.00),
(3, 'ฝอยทอง', 60.00),
(4, 'เม็ดขนุน', 45.00),
(5, 'ขนมชั้น', 55.00),
(6, 'บัวลอยเผือกไข่หวาน', 50.00),
(7, 'ขนมเปียกปูนกะทิสด', 45.00),
(8, 'ขนมหม้อแกงเผือก', 70.00),
(9, 'ตะโก้แห้วเผือก', 50.00),
(10, 'สังขยาฟักทอง', 80.00);

-- Default Test Customer
INSERT INTO CUSTOMER (Username, Password, Title, Name, Surname, Phone, Email) VALUES
('somchai', '123456', 'นาย', 'สมชาย', 'ใจดี', '081-2345678', 'somchai@email.com');

-- Seed Ingredients for Inventory & Purchasing System
INSERT INTO INGREDIENT (IngredientID, IngredientName, StockQty) VALUES
(1, 'แป้งข้าวเจ้า (g)', 3000.00),
(2, 'น้ำตาลปี๊บ (g)', 2500.00),
(3, 'กะทิสด (ml)', 4000.00),
(4, 'ไข่ไก่ (ฟอง)', 250.00),
(5, 'เผือกนึ่ง (g)', 1200.00),
(6, 'ถั่วเขียวเลาะเปลือก (g)', 350.00); -- ต่ำกว่า 500 (Trigger Purchasing Alert)

-- Seed Product-Ingredient Mapping
INSERT INTO PRODUCT_INGREDIENT (ProductID, IngredientID) VALUES
(1, 1), (1, 2), (1, 4),
(2, 1), (2, 2), (2, 4),
(3, 1), (3, 2), (3, 4),
(4, 1), (4, 2), (4, 6),
(5, 1), (5, 3),
(6, 3), (6, 4), (6, 5),
(7, 1), (7, 3),
(8, 3), (8, 4), (8, 5),
(9, 1), (9, 3), (9, 5),
(10, 3), (10, 4), (10, 5);

