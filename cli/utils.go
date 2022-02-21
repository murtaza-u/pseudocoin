package cli

import (
	"crypto/rand"
	"io"
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

var chars = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e',
	'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

func randomString() (string, error) {
	b := make([]byte, 6)
	_, err := io.ReadAtLeast(rand.Reader, b, 6)

	if err != nil {
		return "", err
	}

	for i := 0; i < len(b); i++ {
		b[i] = chars[int(b[i])%len(chars)]
	}

	return string(b), nil
}
