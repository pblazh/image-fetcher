package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

type Config struct {
	Bucket      string
	Out         string
	Mask        string
	IdWorkers   int
	FileWorkers int
}

func (cfg *Config) Parse() error {
	flag.StringVar(&cfg.Out, "out", "out", "out folder")
	flag.StringVar(&cfg.Bucket, "bucket", "", "bucket name")
	flag.StringVar(&cfg.Mask, "mask", "_*_*", "file mask")
	flag.IntVar(&cfg.IdWorkers, "idWorkers", 5, "ids workers")
	flag.IntVar(&cfg.FileWorkers, "fileWorkers", 3, "file workers")

	envProd := false
	flag.BoolVar(&envProd, "prod", false, "Revolut vision prod bucket")

	envDev := false
	flag.BoolVar(&envDev, "dev", false, "Revolut vision dev bucket")

	flag.Parse()

	if cfg.Out == "" {
		return fmt.Errorf("output is empty")
	}
	absRoot, err := filepath.Abs(cfg.Out)
	if err != nil {
		return fmt.Errorf("resolve output: %w", err)
	}
	cfg.Out = absRoot

	if cfg.IdWorkers <= 0 {
		return fmt.Errorf("too few id workers")
	}

	if cfg.FileWorkers <= 0 {
		return fmt.Errorf("too few file workers")
	}

	if envProd {
		cfg.Bucket = "revolut-prod-apps_vision-scans"
	}

	if envDev {
		cfg.Bucket = "revolut-dev-apps_vision-scans"
	}

	if cfg.Bucket == "" {
		return fmt.Errorf("no bucket provided")
	}

	out, err := filepath.Abs(cfg.Out)
	if err != nil {
		return fmt.Errorf("can not resolve output path: %w", err)
	}
	cfg.Out = out

	return nil
}
