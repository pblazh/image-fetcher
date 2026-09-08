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
}

func ListObjects(ctx context.Context, bkt *storage.BucketHandle, out string, id string, requests chan Request) {
	objs := bkt.Objects(ctx, &storage.Query{MatchGlob: id + "_*_*"})

	for {
		obj, err := objs.Next()
		if err == iterator.Done || err == context.Canceled {
			break
		}

		if err != nil {
			log.Println(fmt.Errorf("failed to list objects, %w", err))
			log.Println(err)
		}

		requests <- Request{Id: id, Name: obj.Name, Out: out}
	}

	close(requests)
}

func fetchForIds(ctx context.Context, cfg Config, ids <-chan string) {
	client, err := storage.NewClient(ctx)
	if err == context.Canceled {
		return
	}

	if err != nil {
		log.Println(err)
	}

	bkt := client.Bucket(cfg.Bucket)

	for id := range ids {
		fetchForId(ctx, bkt, cfg, id)
	}
}

func fetchForId(ctx context.Context, bkt *storage.BucketHandle, config Config, id string) {
	err := makeDir(path.Join(config.Out, id))
	if err == context.Canceled {
		return
	}
	if err != nil {
		log.Println(err)
	}

	requests := make(chan Request)
	go ListObjects(ctx, bkt, config.Out, id, requests)

	var wg sync.WaitGroup

	for range config.FileWorkers {
		wg.Go(func() {
			FetchObjects(ctx, bkt, requests)
		})
	}

	wg.Wait()
}
