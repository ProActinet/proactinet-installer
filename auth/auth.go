package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)




func Authentication(dsn string) {
	Welcome()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	StopLoader()
	
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to connect to database: %v\033[0m\n", err)
		os.Exit(1)
	}

	// Fetch all usernames and emails from the database
	var users []struct {
		Username string
		Email    string
	}

	if err := db.Table("users_user").Select("username, email").Scan(&users).Error; err != nil {
		fmt.Printf("\033[31m❌ Error fetching user data: %v\033[0m\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\033[36mEnter username (or type 'exit' to quit): \033[0m")
		inputUsername, _ := reader.ReadString('\n')
		inputUsername = strings.TrimSpace(inputUsername)

		if strings.ToLower(inputUsername) == "exit" {
			fmt.Println("\033[33mExiting program.\033[0m")
			os.Exit(0)
		}

		fmt.Print("\033[36mEnter email: \033[0m")
		inputEmail, _ := reader.ReadString('\n')
		inputEmail = strings.TrimSpace(inputEmail)

		// Check if username and email match any record
		validUser := false
		for _, user := range users {
			if user.Username == inputUsername && user.Email == inputEmail {
				validUser = true
				break
			}
		}

		if validUser {
			fmt.Printf("\033[32m✅ Welcome, %s!\033[0m\n", inputUsername)
			SaveCreds(inputUsername, inputEmail, "licence ABC")
			break
		} else {
			fmt.Println("\033[31m❌ Invalid username or email. Please try again.\033[0m")
		}
	}
}
