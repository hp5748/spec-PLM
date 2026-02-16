package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:plm123456@tcp(localhost:3306)/plm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to generate hash:", err)
	}

	result := db.Exec("UPDATE users SET password = ? WHERE username = ?", string(hash), "admin")
	if result.Error != nil {
		log.Fatal("Failed to update password:", result.Error)
	}

	fmt.Println("Password updated successfully for admin user")
	fmt.Println("New hash:", string(hash))
}
