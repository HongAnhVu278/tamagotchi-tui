package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func readLines(path string) ([]string, error) {
	var output []string
	file, err := os.Open(path)

	// if no file exist => no commit/no food yet => completely ok
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	// file exists => err opening the file
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineStr := scanner.Text()
		output = append(output, lineStr)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file: %w", err)
	}

	return output, nil

}

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
