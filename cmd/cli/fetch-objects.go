package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

func FetchObjects(ctx context.Context, bkt *storage.BucketHandle, requests <-chan Request) {
	for req := range requests {
		obj := bkt.Object(req.Name)

		r, err := obj.NewReader(ctx)
		if errors.Is(err, context.Canceled) {
			break
		}

		if err != nil {
			log.Println(fmt.Errorf("failed to create a reader for %s/%s, %w", req.Id, req.Name, err))
			break
		}

		defer func() {
			err := r.Close()
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Println(err)
			}
		}()

		file, err := os.Create(req.FileName)

		if errors.Is(err, context.Canceled) {
			break
		}

		if err != nil {
			log.Println(fmt.Errorf("failed to create %s, %w", req.FileName, err))
		}

		if _, err := io.Copy(file, r); err != nil {
			if err != context.Canceled {
				log.Println(fmt.Errorf("failed to write %s, %w", req.FileName, err))
			}
		}
	}
}
