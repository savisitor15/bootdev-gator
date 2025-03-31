package config

import (
	"fmt"
	"os"
	"io"
	"encoding/json"
)

const configName = ".gatorconfig.json";

func getConfigFilePath() (string, error){
	pth, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home path! err:\n %w", err)
	}
	pth = pth + string(os.PathSeparator) + configName
	return pth, nil
}

func backup(sourceFile string, destinationFile string) error{
		source, err := os.Open(sourceFile)  //open the source file 
		if err != nil {
		   return err
		}
		defer source.Close()
	 
		destination, err := os.Create(destinationFile)  //create the destination file
		if err != nil {
		   return err
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)  //copy the contents of source to destination file
		if err != nil {
		   return err
		}
		return nil
}

func write(cfg Config) error {
	pth, err := getConfigFilePath()
	if err != nil {
		return err
	}
	rawJson, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("unable to marshal cfg:\n %w", err)
	}
	// backup the file
	err = backup(pth, pth + ".bak")
	if err != nil{
		return fmt.Errorf("error taking snapshot of cfg:\n %w", err)
	}
	err = os.WriteFile(pth, rawJson, 0644)
	if err != nil {
		return fmt.Errorf("error writing back cfg:\n %w", err)
	}
	os.Remove(pth + ".bak") // we don't need the backup anymore
	return nil
}

func Read()(Config, error){
	var outConfig Config
	pth, err := getConfigFilePath()
	if err != nil{
		return outConfig, err
	}
	jsonByte, err := os.ReadFile(pth)
	if err != nil {
		return outConfig, fmt.Errorf("unable to open config json:\n %w", err)
	}
	err = json.Unmarshal(jsonByte, &outConfig)
	if err != nil {
		return outConfig, fmt.Errorf("unable to parse config json:\n %w", err)
	}
	return outConfig, nil
}