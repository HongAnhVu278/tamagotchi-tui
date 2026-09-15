package main

import (
	"errors"
	"fmt"
	"io"
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
	feedAmount  = 20 // hunger gained per commit
	playAmount  = 10 // happiness gained per gratitude entry

	hungerDecayPerDay    = 20.0
	happinessDecayPerDay = 10.0

	sadBelow      = 50
	hungryBelow   = 50
	thrivingAbove = 80

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

// usage writes bitly's own help text. Takes an io.Writer so "help" can print
// to stdout (it's the answer) while an unknown command prints to stderr.
func usage(w io.Writer) {
	fmt.Fprint(w, `bitly - a terminal pet fed by your git commits

usage:
  bitly              launch the pet (same as "bitly run")
  bitly run          launch the pet
  bitly commit ...   run "git commit ..." and feed the pet on success
  bitly speak        print a random line from your gratitude log
  bitly help         show this message

everything after "bitly commit" is passed straight through to git:
  bitly commit -m "fix the thing"
  bitly commit --amend --no-edit
`)
}

func main() {
	cmd := "run"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	var err error

	switch cmd {
	case "speak":
		err = runSpeak()
	case "run":
		err = runTui()
	case "commit":
		err = runGit()
	case "help", "--help", "-h":
		usage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "bitly: unknown command %q\n\n", cmd)
		usage(os.Stderr)
		os.Exit(2)
	}

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "bitly:", err)
		os.Exit(1)
	}
}
