package main

import (
	"fmt"
	"os"
	"time"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	intro     string
	hunger    int // higher hunger = more full
	happiness int
	frame     int
}

// set up tick for time-based animation
type tickMsg time.Time

// tickCmd() return tea.Cmd
// which is tea.Tick(...)
// tea.Tick waits 500ms then runs func(t time.Time) which returns TickMsg
func tickCmd() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}


// define initial state
func initialModel() model {
	return model{
		intro:     "hello, i'm bitly!!",
		hunger:    40,
		happiness: 40,
		frame:     0,

	}
}

// initialize tickCmd so that it sends a TickMsg to Update
func (m model) Init() tea.Cmd {
	return tickCmd()
}

// “something happened” comes in the form of a Msg, which can be any type
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.frame = (m.frame + 1) % len(neutralFrames)

		// must return a fresh tea.Tick command on every tick so that further ticks are scheduled
		return m, tickCmd()

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
	s := m.intro
	s += "\n----------------\n"

	// state-based rendering

	if m.happiness < 50 {
		s += sadFrames[m.frame]
	} else if m.hunger < 50 {
		s += hungryFrames[m.frame]
	} else if m.happiness > 80 && m.hunger > 80 {
		s += happyFrames[m.frame]
	} else {
		s += neutralFrames[m.frame]
	}
	// s += fmt.Sprintf("\n[debug] frame=%d", m.frame)

	s += "\n----------------"
	s += fmt.Sprintf("\nhunger:%d\nhappiness:%d", m.hunger, m.happiness)
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
