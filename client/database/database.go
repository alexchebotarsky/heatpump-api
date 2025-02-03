package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Database struct {
	file *os.File
	data map[string]string
}

func New(filename string, defaults map[string]string) (*Database, error) {
	var d Database
	var err error

	d.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error opening/creating database file %q: %v", filename, err)
	}
	d.data = make(map[string]string, len(defaults))

	err = json.NewDecoder(d.file).Decode(&d.data)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error decoding json from database file: %v", err)
	}

	for key, value := range defaults {
		_, ok := d.data[key]
		if !ok {
			d.data[key] = value
		}
	}

	err = d.updateFile()
	if err != nil {
		return nil, fmt.Errorf("error updating database file: %v", err)
	}

	return &d, nil
}

func (d *Database) updateFile() error {
	// Go to the beginning of the file
	_, err := d.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error going to the beginning of database file: %v", err)
	}

	// Clear the file
	err = d.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error clearing database file: %v", err)
	}

	// Encode the data to the file
	err = json.NewEncoder(d.file).Encode(d.data)
	if err != nil {
		return fmt.Errorf("error encoding json to database file: %v", err)
	}

	return nil
}

func (d *Database) Close() error {
	err := d.file.Close()
	if err != nil {
		return fmt.Errorf("error closing database file: %v", err)
	}

	return nil
}

func (d *Database) GetStr(key string) (string, error) {
	value, ok := d.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in database", key)
	}

	return value, nil
}

func (d *Database) GetInt(key string) (int, error) {
	value, ok := d.data[key]
	if !ok {
		return 0, fmt.Errorf("key %q not found in database", key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("error converting value %q to int: %v", value, err)
	}

	return intValue, nil
}

func (d *Database) Set(key, value string) error {
	d.data[key] = value

	err := d.updateFile()
	if err != nil {
		return fmt.Errorf("error updating database file: %v", err)
	}

	return nil
}

func (d *Database) Delete(key string) error {
	delete(d.data, key)

	err := d.updateFile()
	if err != nil {
		return fmt.Errorf("error updating database file: %v", err)
	}

	return nil
}
