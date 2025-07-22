package main

import (
	"fmt"
	"log"
	"os"
	
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run debug_password.go <password> <hash>")
		os.Exit(1)
	}
	
	password := os.Args[1]
	hash := os.Args[2]
	
	fmt.Printf("Testing password: %s\n", password)
	fmt.Printf("Against hash: %s\n", hash)
	
	if checkPasswordHash(password, hash) {
		fmt.Println("✅ Password matches!")
	} else {
		fmt.Println("❌ Password does not match!")
	}
}
