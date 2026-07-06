package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	intro     string
	hunger    int // higher hunger = more full
	happiness int
	frame     int
}

type saveData struct {
	Hunger    int `json:"hunger"`
	Happiness int `json:"happiness"`
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

func saveFilePath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	path := filepath.Join(home, ".bitly", "save.json")

	return path, nil

}

func loadDataFromFile() (saveData, error) {
	// check if file exist

	path, err := saveFilePath()
	if err != nil {
		return saveData{}, err
	}

	_, err = os.Stat(path)

	// if not, create new file with default val + return the default model
	if errors.Is(err, os.ErrNotExist) {

		initialData := saveData{40, 40}

		//marshall data into json
		data, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return saveData{}, err
		}

		// create the dir
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return saveData{}, err
		}

		//write json data to a file
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			return saveData{}, err
		}

		return initialData, nil
	} else {
		// if yes, return model from data file
		// read the json file
		data, err := os.ReadFile(path)

		if err != nil {
			return saveData{}, err
		}

		// initialize var to hold data
		var initialData saveData

		// unmarshall json bytes into the struct
		err = json.Unmarshal(data, &initialData)
		if err != nil {
			return saveData{}, err
		}

		return initialData, nil
	}

}

func saveDataToFile(m model) error {
	// get the current data from model
	// saved to saved Model
	// marshall to json and overwritee

	var savedData saveData

	savedData.Happiness = m.happiness
	savedData.Hunger = m.hunger

	path, err := saveFilePath()
	if err != nil {
		return err
	}

	//marshall data into json
	data, err := json.MarshalIndent(savedData, "", "  ")
	if err != nil {
		return err
	}

	//write data to file
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// define initial state
func initialModel(data saveData) model {

	return model{
		intro:     "welcome to bitly",
		hunger:    data.Hunger,
		happiness: data.Happiness,
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
			_ = saveDataToFile(m)
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
	/*
		TODO:
		whenever ./bitly commit => append the timestampt line to /.bitly/commits.log
	*/
	switch arg := os.Args[0]; arg {
	case "./bitly":
		cmd := exec.Command("git", os.Args[1:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			os.Exit(1)
		}
		home, _ := os.UserHomeDir()
		f, _ := os.OpenFile(filepath.Join(home, ".bitly", "commits.log"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		defer f.Close()
		//

		fmt.Fprintf(f, "%d commit\n", time.Now().Unix())

	default:
		data, err := loadDataFromFile()

		if err != nil {
			fmt.Println("error loading model", err)
			os.Exit(1)
		}

		p := tea.NewProgram(initialModel(data))

		if _, err := p.Run(); err != nil {
			fmt.Printf("there's been an error: %v", err)
			os.Exit(1)
		}
	}
}
