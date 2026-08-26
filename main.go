package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"
)

const (
	frameCount  = 6
	minStat     = 0
	maxStat     = 100
	defaultStat = 40
	feedAmount  = 10 // hunger gained per commit
	playAmount  = 10 // happiness gained per gratitude entry

	hungerDecayPerDay    = 20.0
	happinessDecayPerDay = 10.0

	secondsPerDay = 86400
	threshold     = 5 * secondsPerDay // 5 days with no commit = death
)

func runTui() error {
	data, err := loadDataFromFile()

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	p := tea.NewProgram(initialModel(data))

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("there's been an error: %w", err)
	}
	return nil
}

func runSpeak() error {
	speakLine, err := speak()

	if err != nil {
		return fmt.Errorf("speak: %w\n", err)
	}

	fmt.Println(speakLine)
	return nil

}

func runGit() error {
	cmd := exec.Command("git", os.Args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// only a commit feeds the pet
	if os.Args[1] != "commit" {
		return nil
	}

	err := appendLine("commits.log", fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("record commit: %w", err)
	}

	return nil

}

func main() {
	var err error
	if len(os.Args) < 2 {
		err = runTui()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	switch os.Args[1] {
	case "speak":
		err = runSpeak()

	default:
		err = runGit()
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Println(err)
		os.Exit(1)
	}

}
