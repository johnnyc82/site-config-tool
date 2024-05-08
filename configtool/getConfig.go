package configtool

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const CONFIG_PATH = ".folder"

func ConfigRootDir() (string, error) {
	usrRoot, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// fmt.Println(usrRoot)

	configPath := fmt.Sprintf("%s/%s", usrRoot, CONFIG_PATH)
	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("cannot read expected config folder location")
	}

	return configPath, nil
}

func PullProjectFiles() error {
	path, err := ConfigRootDir()
	if err != nil {
		return err
	}

	gitPath := fmt.Sprintf("%s/project-config/config", path)
	contentsPath := fmt.Sprintf("%s/project-config/config", path)

	if _, err := os.Stat(contentsPath); os.IsNotExist(err) {
		log.Println("Cloning the project config repo...")
		cmd := exec.Command("git", "clone", "git@bitbucket.org:ivystreetptyltd/project-config.git")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = path
		err := cmd.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Println(os.Getwd())

		log.Println("Checking for updates to project config repo...")
		cmd := exec.Command("git", "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = gitPath
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
