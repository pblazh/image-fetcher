package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

func FetchObjects(ctx context.Context, bkt *storage.BucketHandle, requests <-chan Request) {
	for req := range requests {
		err := fetchObject(ctx, bkt, req)
		if err != nil {
			log.Println(err)
		}
	}
}

func fetchObject(ctx context.Context, bkt *storage.BucketHandle, req Request) error {
	obj := bkt.Object(req.Name)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a reader for %s/%s, %w", req.Id, req.Name, err)
	}

	stat, err := os.Stat(req.FileName)
	if err == nil && stat.Size() == attrs.Size && !stat.IsDir() {
		return nil
	}

	r, err := obj.NewReader(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to create a reader for %s/%s, %w", req.Id, req.Name, err)
	}

	defer func() { _ = closeF(r) }()

	file, err := os.CreateTemp(filepath.Dir(req.FileName), filepath.Base(req.FileName)+"_")
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to create a tmp file %s, %w", req.Name, err)
	}

	defer func() { _ = closeF(file) }()
	defer func() { _ = rmF(file) }()
	fileName := file.Name()

	err = os.Chmod(fileName, 0o600)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to create %s, %w", req.FileName, err)
	}

	if _, err := io.Copy(file, r); err != nil {
		if err != context.Canceled {
			return fmt.Errorf("failed to write %s, %w", req.FileName, err)
		}
		return err
	}

	err = closeF(file)
	if err != nil {
		return fmt.Errorf("failed to write a file %s, %w", req.FileName, err)
	}

	err = os.Rename(fileName, req.FileName)
	if err != nil {
		return fmt.Errorf("failed to rename a tmp file %s, %w", req.FileName, err)
	}
	return nil
}

func closeF(c io.Closer) error {
	err := c.Close()
	if err != nil && errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func rmF(f *os.File) error {
	return os.Remove(f.Name())
}
