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
	feedAmount  = 5
	playAmount  = 10
	threshold   = 86400 //5 days
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

	err := appendLine("commits.log", fmt.Sprintf("%d\n", time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("record commit: %w\n", err)
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
