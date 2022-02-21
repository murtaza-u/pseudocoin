package cli

import (
	"io/ioutil"
	"os"
	"path"
)

func getHome() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(userHome, ".pseudocoin"), nil
}

func createHome() error {
	home, err := getHome()
	if err != nil {
		return err
	}

	return os.MkdirAll(home, 0700)
}

func write(file string, data []byte) error {
	err := createHome()
	if err != nil {
		return err
	}

	home, err := getHome()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(home, file), data, 0600)
}
