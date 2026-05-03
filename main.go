package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	intro     string
	hunger    int // higher hunger = more full
	happiness int
}

// define initial state
func initialModel() model {
	return model{
		intro:     "hello, this is your pet!",
		hunger:    50,
		happiness: 50,
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

		// key to feed
		case "f":
			m.hunger += 10
			if m.hunger > 100 {
				m.hunger = 100
			}

		// key to play:
		case "p":
			m.happiness += 10
			if m.happiness > 100 {
				m.happiness = 100
			}

		// key to exit program
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m model) View() tea.View {
	s := fmt.Sprintf("%s\nhunger:%d\nhappiness:%d", m.intro, m.hunger, m.happiness)
	s += "\n----------------"
	s += "\n[f]eed [p]lay [q]uit"

	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}
}
