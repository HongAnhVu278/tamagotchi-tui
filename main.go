package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	intro string
}

// define initial state
func initialModel() model {
	return model{
		intro: "hello, this is your pet!",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// “something happened” comes in the form of a Msg, which can be any type
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// is it a key pressed?
	case tea.KeyPressMsg:

		// what's the key pressed
		switch msg.String() {

		// key to exit program
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m model) View() tea.View {
	s := m.intro

	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
