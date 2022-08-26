package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initQuestions()
	m := questions[currentQ]
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
