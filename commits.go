package main

import (
	"fmt"
	"strconv"
	"time"
)

func applyCommits(hunger float64, cursor int64, lines []string, hatched bool) (float64, int64) {
	for _, line := range lines {
		ts, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			continue
		}
		if ts <= cursor {
			continue
		}
		if hatched {
			hunger = decayed(hunger, cursor, ts, hungerDecayPerDay)
		}
		hunger = clamp(hunger + feedAmount)
		cursor = ts
	}
	return hunger, cursor
}

func feed(m model) (model, error) {

	path, err := bitlyPath("commits.log")

	if err != nil {
		return m, fmt.Errorf("get log path to feed pet: %w", err)
	}

	lines, err := readLines(path)

	if err != nil {
		return m, fmt.Errorf("get output from log file: %w", err)
	}

	m.hunger, m.lastRead = applyCommits(m.hunger, m.lastRead, lines, isHatched(m))

	return m, nil
}

func currentHunger(m model) float64 {
	if isHatched(m) {
		return decayed(m.hunger, m.lastRead, time.Now().Unix(), hungerDecayPerDay)
	}
	return m.hunger

}

func isDead(m model) bool {
	if !isHatched(m) {
		return false // unborn/egg: never fed, can't die yet
	}

	now := time.Now().Unix()
	return (now - m.lastRead) > threshold

}
