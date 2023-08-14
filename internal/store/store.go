package store

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/atotto/clipboard"
)

// Returns the contents of $HOME/.hitman
// Returns placeholder text when an error is encountered
func LoadText() string {
	defaultText := "GET www.example.com"

	var home string

	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	if bytes, err := ioutil.ReadFile(path.Join(home, ".hitman")); err != nil {
		return defaultText
	} else {
		return string(bytes)
	}
}

// Save the provided string into $HOME/.hitman
// Fails silently
func SaveText(text string) {
	var home string

	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	_ = ioutil.WriteFile(path.Join(home, ".hitman"), []byte(text), 0644)
}

func CopyText(text string) error {
	return clipboard.WriteAll(text)
}
