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

var (
	IDS_WORKERS  = 5
	FILE_WORKERS = 3
)

type Request struct {
	Id   string
	Name string
	Out  string
}

func ListObjects(ctx context.Context, bkt *storage.BucketHandle, out string, id string, requests chan Request) {
	objs := bkt.Objects(ctx, &storage.Query{MatchGlob: id + "_*_*"})
	// objs := bkt.Objects(ctx, &storage.Query{MatchGlob: id + "_*_SELFIE"})

	for {
		obj, err := objs.Next()
		if err == iterator.Done {
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

func fetchForIds(ctx context.Context, bucket string, out string, ids <-chan string) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println(err)
	}

	bkt := client.Bucket(bucket)

	for id := range ids {
		fetchForId(ctx, bkt, out, id)
	}
}

func fetchForId(ctx context.Context, bkt *storage.BucketHandle, out string, id string) {
	err := makeDir(path.Join(out, id))
	if err != nil {
		log.Println(err)
	}

	requests := make(chan Request)
	go ListObjects(ctx, bkt, out, id, requests)

	var wg sync.WaitGroup

	for range FILE_WORKERS {
		wg.Go(func() {
			FetchObjects(ctx, bkt, requests)
		})
	}

	wg.Wait()
}
