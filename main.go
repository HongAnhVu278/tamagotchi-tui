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

func main() {
	if len(os.Args) < 2 {
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
		return

	}
	switch os.Args[1] {
	case "speak":
		speakLine, err := speak()

		if err != nil {
			fmt.Fprintf(os.Stderr, "speak: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(speakLine)
		return
	default:
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

		err := appendLine("commits.log", fmt.Sprintf("%d\n", time.Now().Unix()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "record commit: %v\n", err)
			os.Exit(1)
		}

	}

}
