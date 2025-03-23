package auth

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define a struct for the data
type LicenseData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	License string `json:"license"`
}

func SaveCreds(name string, email string, licence string ) {
	data := LicenseData{
		License: licence,
		Name:    name,
		Email:   email,
	}

	// Convert struct to JSON (with indentation for readability)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Define file path
	filePath := "./license.json"

	// Write JSON data to file
	err = os.WriteFile(filePath, jsonData, 0644) // 0644 = Read/Write for owner, Read for others
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("✅ JSON file saved successfully at", filePath)
}