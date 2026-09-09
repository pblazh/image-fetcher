package main

import (
	"flag"
	"log"
	"os"
	"path"
)

type Config struct {
	Bucket      string
	Out         string
	Mask        string
	IdWorkers   int
	FileWorkers int
}

func (cfg *Config) Parse() {
	flag.StringVar(&cfg.Out, "out", "out", "out folder")
	flag.StringVar(&cfg.Bucket, "bucket", "revolut-prod-apps_vision-scans", "bucket name")
	flag.StringVar(&cfg.Mask, "mask", "_*_*", "file mask")
	flag.IntVar(&cfg.IdWorkers, "idWorkers", 5, "ids workers")
	flag.IntVar(&cfg.FileWorkers, "fileWorkers", 3, "file workers")

	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Out = path.Join(pwd, cfg.Out)
}
