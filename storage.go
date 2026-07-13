package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func bitlyPath(name string) (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("get user dir: %w", err)
	}

	path := filepath.Join(home, ".bitly", name)

	return path, nil
}

func loadDataFromFile() (saveData, error) {
	// check if file exist

	path, err := bitlyPath("save.json")
	if err != nil {
		return saveData{}, fmt.Errorf("get save path: %w", err)
	}

	_, err = os.Stat(path)

	// if not, create new file with default val + return the default model
	if errors.Is(err, os.ErrNotExist) {

		initialData := saveData{Hunger: defaultStat, Happiness: defaultStat, LastRead: 0}

		//marshall data into json
		data, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return saveData{}, fmt.Errorf("marshall data: %w", err)
		}

		// create the dir
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return saveData{}, fmt.Errorf("create the dir: %w", err)
		}

		//write json data to a file
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			return saveData{}, fmt.Errorf("write json data to a file: %w", err)
		}

		return initialData, nil
	}
	// if yes, return model from data file
	// read the json file
	data, err := os.ReadFile(path)

	if err != nil {
		return saveData{}, fmt.Errorf("read the json file: %w", err)
	}

	// initialize var to hold data
	var initialData saveData

	// unmarshall json bytes into the struct
	err = json.Unmarshal(data, &initialData)
	if err != nil {
		return saveData{}, fmt.Errorf("unmarshall json: %w", err)
	}

	return initialData, nil

}

func saveDataToFile(m model) error {
	// get the current data from model
	// saved to saved Model
	// marshall to json and overwritee

	var savedData saveData

	savedData.Happiness = m.happiness
	savedData.Hunger = m.hunger
	savedData.LastRead = m.lastRead

	path, err := bitlyPath("save.json")
	if err != nil {
		return fmt.Errorf("get save path: %w", err)
	}

	//marshall data into json
	data, err := json.MarshalIndent(savedData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshall data into json: %w", err)
	}

	//write data to file
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("write data to file: %w", err)
	}

	return nil
}

func readLines(name string) ([]string, error) {
	var output []string
	file, err := os.Open(name)

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

func appendLine(name string, line string) error {
	path, err := bitlyPath(name)

	if err != nil {
		return fmt.Errorf("resolve log path: %v", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create %s: %v", dir, err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %v", f, err)
	}

	defer f.Close()
	if _, err := fmt.Fprintf(f, "%s\n", line); err != nil {
		return fmt.Errorf("write %s: %w", path, err)

	}
	return nil
}
