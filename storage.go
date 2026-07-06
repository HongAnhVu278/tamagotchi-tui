package main

import (
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

		initialData := saveData{Hunger: defaultStat, Happiness: defaultStat}

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
