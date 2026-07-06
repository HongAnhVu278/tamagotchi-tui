package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
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

	path, err := bitlyPath("commits.log")

	if err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)

	if err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Printf("there's been an error: %v", err)
		os.Exit(1)
	}

	defer f.Close()

	fmt.Fprintf(f, "%d commit\n", time.Now().Unix())
}
