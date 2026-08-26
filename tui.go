package main

import (
	"fmt"

	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// hunger/happiness are stored as of their anchor (lastRead / lastGratitude),
// not as of now. Read them through currentHunger / currentHappiness.
type model struct {
	intro         string
	hunger        float64 // higher hunger = more full
	happiness     float64
	lastRead      int64 // commit ts: log cursor + hunger anchor
	lastGratitude int64 // happiness anchor
	isDead        bool
	frame         int
	textInput     textinput.Model
}

type saveData struct {
	Hunger        float64 `json:"hunger"`
	Happiness     float64 `json:"happiness"`
	LastRead      int64   `json:"last_read"`
	LastGratitude int64   `json:"last_gratitude"`
}

// check if the user commit code and add gratitude
func isHatched(m model) bool {
	return m.lastRead != 0 && m.lastGratitude != 0
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
		intro:         "hi! i'm bitly",
		hunger:        data.Hunger,
		happiness:     data.Happiness,
		lastRead:      data.LastRead,
		lastGratitude: data.LastGratitude,
		frame:         0,
		textInput:     newTextInput(),
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

			now := time.Now().Unix()
			m.happiness = decayed(m.happiness, m.lastGratitude, now, happinessDecayPerDay)
			m.happiness = clamp(m.happiness + playAmount)
			m.lastGratitude = now

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

	hunger := currentHunger(m)
	happiness := currentHappiness(m)

	s := m.intro
	s += "\n----------------\n"

	// state-based rendering
	if m.isDead {
		s += deadFrames[m.frame]
		return tea.NewView(s)
	}

	// rendering based on hatch:

	if !isHatched(m) {
		if m.lastRead == 0 {
			s += "Commit your code to hatch your egg!\n"
		}
		if m.lastGratitude == 0 {
			s += "Write down something you're grateful for to hatch your egg!\n"
		}
		s += "----------------\n"
		s += eggFrames[m.frame]

	} else {
		s += currentFrames(hunger, happiness)[m.frame]

		s += "\n----------------"
		s += fmt.Sprintf("\nhunger:%.0f\nhappiness:%.0f", hunger, happiness)

	}
	s += "\n----------------"
	s += lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())

	view := tea.NewView(s)
	//view.Cursor = cursor

	return view
}

func currentFrames(hunger, happiness float64) []string {
	switch {
	case hunger < hungryBelow:
		return hungryFrames
	case happiness < sadBelow:
		return sadFrames
	case hunger > thrivingAbove && happiness > thrivingAbove:
		return happyFrames
	default:
		return neutralFrames
	}
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
