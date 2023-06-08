package config

import (
	"bufio"
	"os"
	"path"
)

var (
	tokenFile = pathFromHome(".diarier_token.txt")
	_token    string
	_loaded   bool
)

// pathFromHome appends the home dir, in front of the given relativePath
func pathFromHome(relativePath string) string {
	home := os.Getenv("HOME")
	return path.Join(home, relativePath)
}

func ReadToken() (string, error) {
	if _loaded {
		return _token, nil
	}

	file, err := os.Open(tokenFile)
	_loaded = true
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	_token = scanner.Text()
	return _token, scanner.Err()
}

func UpdateToken(token string) error {
	file, err := os.Create(tokenFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return err
	}

	_token = token
	_loaded = true

	return nil
}
