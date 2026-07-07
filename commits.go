package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func feed(m model) (model, error) {

	cursor := m.lastRead
	hunger := m.hunger

	curr := cursor

	path, err := bitlyPath("commits.log")

	if err != nil {
		return m, fmt.Errorf("get log path to feed pet: %w", err)
	}

	file, err := os.Open(path)

	// if no file exist => no commit/no food yet => completely ok
	if errors.Is(err, os.ErrNotExist) {
		return m, nil
	}

	// commits.log exists => err opening the file
	if err != nil {
		return m, fmt.Errorf("open commits.log: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineStr := scanner.Text()

		ts, err := strconv.ParseInt(lineStr, 10, 64)

		if err != nil {
			// TODO: count skipped lines and render instead
			// log.Printf("not a valid timestamp: %v", err)
			continue
		}

		if ts > cursor {
			hunger += feedAmount
			if hunger > maxStat {
				hunger = maxStat
			}

			curr = ts
		}
	}

	if err := scanner.Err(); err != nil {
		return m, fmt.Errorf("scan commits.log: %w", err)
	}

	m.lastRead = curr
	m.hunger = hunger
	return m, nil
}

func isDead(m model) bool {
	if m.lastRead == 0 {
		return false // unborn/egg: never fed, can't die yet
	}

	now := time.Now().Unix()
	return (now - m.lastRead) > threshold

}
