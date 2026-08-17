package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	LastCompanyID int    `json:"last_company_id"`
	Theme         string `json:"theme"`
}

func settingsPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "retui-accountant-settings.json"
	}

	return filepath.Join(dir, "retui-accountant", "settings.json")
}

func LoadSettings() (*Settings, error) {
	data, err := os.ReadFile(settingsPath())
	if os.IsNotExist(err) {
		return &Settings{}, nil
	}
	if err != nil {
		return nil, err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func SaveSettings(s *Settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	path := settingsPath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func SaveLastCompany(id int) error {
	s, err := LoadSettings()
	if err != nil {
		return err
	}

	s.LastCompanyID = id

	return SaveSettings(s)
}

func ClearLastCompany() error {
	s, err := LoadSettings()
	if err != nil {
		return err
	}

	s.LastCompanyID = 0

	return SaveSettings(s)
}
