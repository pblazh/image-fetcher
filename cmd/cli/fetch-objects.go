package main

import (
	"context"
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
		if err != nil {
			log.Println(fmt.Errorf("failed to create a reader for %s/%s, %w", req.Id, req.Name, err))
			break
		}
		defer r.Close()

		fileName := makeFileName(req.Out, req.Id, req.Name)
		file, err := os.Create(fileName)
		if err != nil {
			log.Println(fmt.Errorf("failed to create %s, %w", fileName, err))
		}

		if _, err := io.Copy(file, r); err != nil {
			log.Println(fmt.Errorf("failed to write %s, %w", fileName, err))
		}
	}
}
