package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func speak() (string, error) {
	path, err := bitlyPath("gratitude.log")

	if err != nil {
		return "", fmt.Errorf("get log path to feed pet: %w", err)
	}

	lines, err := readLines(path)

	if err != nil {
		return "", fmt.Errorf("get output from log file: %w", err)
	}

	if len(lines) == 0 {
		return "", fmt.Errorf("empty log: %w", err)
	}

	randomIdx := rand.IntN(len(lines))

	speakLine := lines[randomIdx]

	return speakLine, nil
}

func currentHappiness(m model) float64 {
	return decayed(m.happiness, m.lastGratitude, time.Now().Unix(), happinessDecayPerDay)
}
