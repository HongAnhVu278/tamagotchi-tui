package main

import (
	"fmt"
	"strconv"
	"time"
)

func applyCommits(hunger int, cursor int64, lines []string) (int, int64) {
	for _, line := range lines {
		ts, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			continue
		}
		if ts > cursor {
			hunger += feedAmount
			if hunger > maxStat {
				hunger = maxStat
			}
			cursor = ts
		}
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

	m.hunger, m.lastRead = applyCommits(m.hunger, m.lastRead, lines)

	return m, nil
}

func isDead(m model) bool {
	if m.lastRead == 0 {
		return false // unborn/egg: never fed, can't die yet
	}

	now := time.Now().Unix()
	return (now - m.lastRead) > threshold

}
