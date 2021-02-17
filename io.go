// +build !js

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func DataDir() (string, error) {
	var xdg string
	if runtime.GOOS == "windows" {
		xdg = os.Getenv("LOCALAPPDATA")
	} else {
		xdg = os.Getenv("XDG_DATA_HOME")
	}
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	dataDir := filepath.Join(xdg, "harmonist")
	_, err := os.Stat(dataDir)
	if err != nil {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return "", fmt.Errorf("building data directory: %v\n", err)
		}
	}
	return dataDir, nil
}

func (g *game) Save() error {
	dataDir, err := DataDir()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	saveFile := filepath.Join(dataDir, "save")
	data, err := g.GameSave()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		g.Print(err.Error())
		return err
	}
	return nil
}

func RemoveSaveFile() error {
	return RemoveDataFile("save")
}

func RemoveReplay() {
	RemoveDataFile("replay.part")
}

func (g *game) Load() (bool, error) {
	dataDir, err := DataDir()
	if err != nil {
		return false, err
	}
	saveFile := filepath.Join(dataDir, "save")
	_, err = os.Stat(saveFile)
	if err != nil {
		// no save file, new state
		return false, nil
	}
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return false, err
	}
	lg, err := g.DecodeGameSave(data)
	if err != nil {
		return false, err
	}
	if lg.Version != Version {
		return false, nil
	}
	*g = *lg
	return true, nil
}

func SaveConfig() error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}
	saveFile := filepath.Join(dataDir, "config.gob")
	data, err := GameConfig.ConfigSave()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (bool, error) {
	dataDir, err := DataDir()
	if err != nil {
		return false, err
	}
	saveFile := filepath.Join(dataDir, "config.gob")
	_, err = os.Stat(saveFile)
	if err != nil {
		// no save file, new state
		return false, nil
	}
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return false, err
	}
	c, err := DecodeConfigSave(data)
	if err != nil {
		return false, err
	}
	if c.Version != GameConfig.Version {
		return false, nil
	}
	GameConfig = *c
	return true, nil
}

func RemoveDataFile(file string) error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}
	dataFile := filepath.Join(dataDir, file)
	_, err = os.Stat(dataFile)
	if err == nil {
		err := os.Remove(dataFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *game) WriteDump() error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dataDir, "dump"), []byte(g.Dump()), 0644)
	if err != nil {
		return fmt.Errorf("writing dump statistics: %v", err)
	}
	return nil
}
