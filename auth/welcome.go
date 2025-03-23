package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var loaderProgram *tea.Program

type stopMsg struct{} // Custom message to stop the loader

type model struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case stopMsg: // Stop the loader when this message is received
		m.quitting = true
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.quitting {
		return "" // Clear the spinner when stopped
	}
	return fmt.Sprintf("\n\n   %s Loading ...\n\n", m.spinner.View())
}

// **StartLoader runs the spinner in a separate goroutine**
func StartLoader() {
	loaderProgram = tea.NewProgram(initialModel())

	// Run the loader in a goroutine
	go func() {
		if _, err := loaderProgram.Run(); err != nil {
			fmt.Println("Loader error:", err)
		}
	}()
}

// **StopLoader sends a stop message to end the loader**
func StopLoader() {
	if loaderProgram != nil {
		loaderProgram.Send(stopMsg{}) // Stop the loader using a custom message
		loaderProgram = nil
	}

	// Clear terminal output
	fmt.Print("\033[H\033[2J")
}

// **Welcome function to display text and start the loader**
func Welcome() {
	cellStyle := lipgloss.NewStyle().
		Width(50).
		Height(5).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("205"))

	table := cellStyle.Render("🎉 Welcome to ProctiNet! 🚀\nYour cutting-edge Anti-Botnet Solution.\n\nPlease login to secure your experience 🔒")
	detailedInfo := "Main Features:\n" +
		"1. Real-time Monitoring: Instantly track network threats as they occur.\n" +
		"2. Live Dashboard Preview: View a dynamic and interactive display of all logged threats.\n" +
		"3. Automated Botnet Protection: Proactively detect and defend against botnet activities."

	infoBox := lipgloss.NewStyle().
		Width(70).
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("82")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("82")).
		Render(detailedInfo)

	fmt.Println(infoBox)
	fmt.Println(table)

	// Start loader
	StartLoader()

	// Automatically stop loader after 3 seconds
	
}

