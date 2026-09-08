package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"path"
	"regexp"
)

var (
	reSelfie = regexp.MustCompile("^*.SELFIE$")
	reVideo  = regexp.MustCompile("^*.VIDEO$")
)

func makeFileName(out string, id string, name string) string {
	fullName := path.Join(id, name)
	var ext string

	switch {
	case reSelfie.MatchString(name):
		ext = ".jpg"
	case reVideo.MatchString(name):
		ext = ".mp4"
	default:
		ext = ".bin"
	}

	return path.Join(out, fullName+ext)
}

func makeDir(dir string) error {
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dir, 0o777)
	}
	return err
}

func ListIds(ctx context.Context, ch chan string) {
	scaner := bufio.NewScanner(os.Stdin)
loop:
	for scaner.Scan() {
		id := scaner.Text()
		select {
		case ch <- id:
			log.Println(id)
		case <-ctx.Done():
			break loop
		}
	}
	close(ch)
}
