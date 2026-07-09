package main

import (
	"fmt"

	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	intro     string
	hunger    int // higher hunger = more full
	happiness int
	lastRead  int64
	isDead    bool
	frame     int
	textInput textinput.Model
}

type saveData struct {
	Hunger    int   `json:"hunger"`
	Happiness int   `json:"happiness"`
	LastRead  int64 `json:"last_read"`
}

// set up tick for time-based animation
type tickMsg time.Time

// tickCmd() return tea.Cmd
// which is tea.Tick(...)
// tea.Tick waits 300ms then runs func(t time.Time) which returns TickMsg
func tickCmd() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// define initial state
func initialModel(data saveData) model {

	return model{
		intro:     "hi! i'm bitly",
		hunger:    data.Hunger,
		happiness: data.Happiness,
		lastRead:  data.LastRead,
		frame:     0,
		textInput: newTextInput(),
	}
}

// initialize tickCmd so that it sends a TickMsg to Update
func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), textinput.Blink)
}

// “something happened” comes in the form of a Msg, which can be any type
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		m, err := feed(m)
		_ = err //TODO: stash on model for view (?)
		m.isDead = isDead(m)

		m.frame = (m.frame + 1) % frameCount

		// must return a fresh tea.Tick command on every tick so that further ticks are scheduled
		return m, tickCmd()

	// is it a key pressed?
	case tea.KeyPressMsg:

		// what's the key pressed
		switch msg.String() {

		// key to feed
		// case "f":
		// 	m.hunger += feedAmount
		// 	if m.hunger > maxStat {
		// 		m.hunger = maxStat
		// 	}

		// // key to play:
		// case "p":
		// 	m.happiness += playAmount
		// 	if m.happiness > maxStat {
		// 		m.happiness = maxStat
		// 	}
		case "enter":
			//record whatever the user enter and log it into gratitude.log
			/*TODO:
			+) check if the text is empty
			+) prompt the answer again after the user enter
			*/
			if m.textInput.Value() == "" {
				break
			}
			dailyGrat := m.textInput.Value()

			_ = appendLine("gratitude.log", dailyGrat) // swallow err

			m.happiness += feedAmount
			if m.happiness > maxStat {
				m.happiness = maxStat
			}
			m.textInput = newTextInput()

		// key to exit program
		case "ctrl+c", "esc":
			_ = saveDataToFile(m)
			return m, tea.Quit
		}

	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() tea.View {
	//var cursor *tea.Cursor

	// if !m.textInput.VirtualCursor() {
	// 	cursor = m.textInput.Cursor()
	// 	cursor.Y += lipgloss.Height(m.headerView())
	// }

	s := m.intro
	s += "\n----------------\n"

	// state-based rendering
	if m.isDead {
		s += deadFrames[m.frame]
		return tea.NewView(s)
	}

	if m.happiness < 50 {
		s += sadFrames[m.frame]
	} else if m.hunger < 50 {
		s += hungryFrames[m.frame]
	} else if m.happiness > 80 && m.hunger > 80 {
		s += happyFrames[m.frame]
	} else {
		s += neutralFrames[m.frame]
	}

	s += "\n----------------"
	s += fmt.Sprintf("\nhunger:%d\nhappiness:%d", m.hunger, m.happiness)
	s += "\n----------------"

	if m.lastRead == 0 {
		s += "\ncommit your code to start the challenge!"
	}

	s += lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())

	view := tea.NewView(s)
	//view.Cursor = cursor

	return view
}

func (m model) headerView() string {
	return "\nCheer me up!! What made today a little better? (Or tell me a joke!)\n"
}
func (m model) footerView() string { return "\n(esc to quit)\n" }

func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Words go here, human"
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return ti
}
