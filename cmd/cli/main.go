package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"path"
	"sync"
)

var BUCKET = "revolut-prod-apps_vision-scans"

func main() {
	var out string
	flag.StringVar(&out, "out", ".", "out folder")

	var bucket string
	flag.StringVar(&bucket, "bucket", BUCKET, "bucket name")

	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	out = path.Join(pwd, out)

	ctx := context.Background()
	ids := make(chan string)

	go ListIds(ctx, ids)

	var wg sync.WaitGroup
	for range IDS_WORKERS {
		wg.Go(func() {
			fetchForIds(ctx, bucket, out, ids)
		})
	}

	wg.Wait()
}

func ListIds(ctx context.Context, ch chan string) {
	scaner := bufio.NewScanner(os.Stdin)
	for scaner.Scan() {
		id := scaner.Text()
		ch <- id
		log.Println(id)
	}

	close(ch)
}
