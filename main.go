package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"installer/auth"
	"installer/db"
	"os"
	"os/exec"
	"strings"
	"time"
)

//go:embed suricata.yaml
var data string

// Command represents a shell command with a description and emoji
type Command struct {
	cmd         string
	description string
	emoji       string
}

func main() {
	// Check if running with sudo
	if os.Geteuid() != 0 {
		fmt.Println("❌ This script must be run with sudo privileges")
		fmt.Println("Please run: sudo go run main.go")
		os.Exit(1)
	}	

	dsn := db.EnvLoader()

	auth.Authentication(dsn)

	commands := []Command{
		{"apt update", "Updating package lists", "📦"},
		{"apt upgrade -y", "Upgrading packages", "⬆️"},
		{"apt -y install libnetfilter-queue-dev libnetfilter-queue1 libnfnetlink-dev libnfnetlink0 jq", "Installing dependencies", "🔧"},
		{"add-apt-repository ppa:oisf/suricata-stable -y", "Adding Suricata repository", "📚"},
		{"apt install suricata -y", "Installing Suricata", "🛡️"},
		{"systemctl stop suricata.service", "Stopping Suricata service", "🛑"},
	}

	for _, cmd := range commands {
		executeCommand(cmd)
	}

	// Inject custom Suricata configuration
	fmt.Println("\n📝 Updating Suricata configuration file...")
	updateSuricataConfig()

	// Update rules
	fmt.Println("\n📜 Listing available rule sources...")
	listRuleSources()

	// Restart Suricata with new configuration
	finalCommands := []Command{
		{"suricata-update", "Updating Suricata rules", "🔄"},
		{"suricata -T -c /etc/suricata/suricata.yaml -v", "Testing configuration", "🧪"},
		{"systemctl restart suricata.service", "Restarting Suricata service", "♻️"},
		{"curl http://testmynids.org/uid/index.html", "Testing IDS functionality", "🌐"},
		{"cat /var/log/suricata/fast.log", "Checking logs", "📋"},
	}

	for _, cmd := range finalCommands {
		executeCommand(cmd)
	}

	fmt.Println("\n✅ Suricata installation and configuration complete! 🚀")

	// Clear logs
	clearSuricataLogs()

}



// Execute a shell command and log output
func executeCommand(cmd Command) {
	fmt.Printf("\n%s %s...\n", cmd.emoji, cmd.description)
	command := exec.Command("bash", "-c", cmd.cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		fmt.Printf("❌ Error executing command: %v\n", err)
		fmt.Println("Would you like to continue anyway? (y/n)")
		if !confirmAction() {
			os.Exit(1)
		}
	}
	time.Sleep(1 * time.Second) // Small delay for readability
}

// Inject the embedded Suricata config into /etc/suricata/suricata.yaml
func updateSuricataConfig() {
	configPath := "/etc/suricata/suricata.yaml"
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to update Suricata config: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var ifn = gtInterfaceDetails("ethernet")
	if ifn == nil {
		panic("No wan network interface found")
	}

	ifnStr, ok := ifn.(string)
	if !ok {
		fmt.Printf("❌ Failed to convert interface name to string\n")
		os.Exit(1)
	}

	newData := strings.Replace(data, "_IFACE_", ifnStr, 1)
	_, err = file.WriteString(newData)
	if err != nil {
		fmt.Printf("❌ Error writing to Suricata config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Suricata configuration updated successfully! 🎉")
}

// List available rule sources
func listRuleSources() {
	cmd := exec.Command("suricata-update", "list-sources")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Enable selected rule sources
func enableRuleSources() {
	fmt.Println("Enter the names of the sources you want to enable (space-separated):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	sources := strings.TrimSpace(input)

	if sources != "" {
		cmd := exec.Command("bash", "-c", "suricata-update enable-source "+sources)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

// Confirm user action

func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

// Clear the testing Logs
func clearSuricataLogs() {
	fmt.Println("\n🧹 Clearing Suricata log files...")

	commands := []string{
		`sudo su -c 'echo "" > /var/log/suricata/eve.json'`,
		`sudo su -c 'echo "" > /var/log/suricata/fast.log'`,
	}

	for _, cmd := range commands {
		command := exec.Command("bash", "-c", cmd)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err := command.Run()
		if err != nil {
			fmt.Printf("❌ Error clearing logs: %v\n", err)
		} else {
			fmt.Printf("✅ Successfully cleared: %s\n", cmd)
		}
	}
}
