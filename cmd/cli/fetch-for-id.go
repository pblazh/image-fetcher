package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Request struct {
	Id   string
	Name string
	Out  string

	FileName string // := makeFileName(req.Out, req.Id, req.Name)
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

		requests <- Request{
			Id: id, Name: obj.Name, Out: config.Out,
			FileName: makeFileName(config.Out, id, obj.Name),
		}
	}
}

func fetchForIds(ctx context.Context, config Config, ids <-chan string) {
	client, err := storage.NewClient(ctx)
	if err == context.Canceled {
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("failed to create a client, %w", err))
		return
	}

	bucket := client.Bucket(config.Bucket)

	for id := range ids {
		fetchForId(ctx, config, bucket, id)
	}
}

func fetchForId(ctx context.Context, config Config, bkt *storage.BucketHandle, id string) {
	err := makeDir(path.Join(config.Out, id))
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
