package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reData  = regexp.MustCompile("^*.((DATA|FINGERPRINT))$")
	reVideo = regexp.MustCompile("^*.VIDEO$")
)

func makeFileName(out string, id string, name string) (string, error) {
	if err := validatePathComponent("object id", id); err != nil {
		return "", err
	}
	if err := validatePathComponent("object name", name); err != nil {
		return "", err
	}

	var ext string

	switch {
	case reData.MatchString(name):
		ext = ".json"
	case reVideo.MatchString(name):
		ext = ".mp4"
	default:
		ext = ".jpg"
	}

	return path.Join(out, id, name+ext), nil
}

func validatePathComponent(label string, component string) error {
	if component == "" || component == "." || component == ".." {
		return fmt.Errorf("%s is empty or reserved", label)
	}
	if filepath.IsAbs(component) {
		return fmt.Errorf("%s must be relative", label)
	}
	if strings.ContainsRune(component, '\x00') {
		return fmt.Errorf("%s contains a NUL byte", label)
	}
	if strings.ContainsAny(component, `/\`) {
		return fmt.Errorf("%s contains a path separator", label)
	}
	if strings.Contains(component, "..") {
		return fmt.Errorf("%s contains a parent path", label)
	}
	if filepath.Clean(component) != component {
		return fmt.Errorf("%s is not a clean path component", label)
	}

	return nil
}

func makeDir(dir string) error {
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0o700)
	}
	return err
}

func ListIds(ctx context.Context, ch chan string) {
	scaner := bufio.NewScanner(os.Stdin)
loop:
	for scaner.Scan() {
		id := scaner.Text()
		id = strings.Trim(id, " \t")

		if id == "" {
			continue
		}

		select {
		case ch <- id:
			log.Println(id)
		case <-ctx.Done():
			break loop
		}
	}

	err := scaner.Err()
	if err != nil {
		log.Println(err)
	}

	close(ch)
}
