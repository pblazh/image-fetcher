package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Request struct {
	Id   string
	Name string

	FileName string
}

func ListObjects(ctx context.Context, config Config, bkt *storage.BucketHandle, id string, requests chan<- Request) {
	objects := bkt.Objects(ctx, &storage.Query{MatchGlob: id + config.Mask})
	defer close(requests)

	for {
		obj, err := objects.Next()
		if err == iterator.Done || err == context.Canceled {
			break
		}

		if err != nil {
			log.Println(fmt.Errorf("failed to list objects, %w", err))
			break
		}

		fileName, err := makeFileName(config.Out, id, obj.Name)
		if err != nil {
			log.Println(fmt.Errorf("rejected unsafe output path, %w", err))
			continue
		}

		requests <- Request{
			Id: id, Name: obj.Name,
			FileName: fileName,
		}
	}
}

func fetchForIds(ctx context.Context, config Config, ids <-chan string) error {
	client, err := storage.NewClient(ctx)
	if err == context.Canceled {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to create a client, %w", err)
	}
	defer func() { _ = client.Close() }()

	bucket := client.Bucket(config.Bucket)

	for id := range ids {
		fetchForId(ctx, config, bucket, id)
	}
	return nil
}

func fetchForId(ctx context.Context, config Config, bkt *storage.BucketHandle, id string) {
	err := validatePathComponent("object id", id)
	if err != nil {
		log.Println(fmt.Errorf("rejected unsafe output directory, %w", err))
		return
	}

	err = makeDir(config.Out)
	if err == context.Canceled {
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("failed to create a dir, %w", err))
		return
	}

	requests := make(chan Request)
	go ListObjects(ctx, config, bkt, id, requests)

	var wg sync.WaitGroup

	for range config.FileWorkers {
		wg.Go(func() {
			FetchObjects(ctx, bkt, requests)
		})
	}

	wg.Wait()
}
