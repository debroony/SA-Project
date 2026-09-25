package main

import (
	"fmt"
	"net/http"
)

func main() {
	// ชี้เป้าหมายไปที่โฟลเดอร์ frontend ที่เราเก็บไฟล์ HTML ไว้ (ถอยหลัง 1 โฟลเดอร์เพื่อเข้า frontend)
	fs := http.FileServer(http.Dir("../frontend"))
	
	// ตั้งค่าให้เมื่อเข้า URL หลัก (/) ให้เสิร์ฟไฟล์จากโฟลเดอร์ frontend
	http.Handle("/", fs)

	fmt.Println("Server is running on http://localhost:8080...")
	// สั่งรันเซิร์ฟเวอร์
	http.ListenAndServe(":8080", nil)
}