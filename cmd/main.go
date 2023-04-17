package main

import (
	"fmt"
	"github.com/nnqq/spider/pkg/config"
	"github.com/nnqq/spider/pkg/console"
	"github.com/nnqq/spider/pkg/list"
	"github.com/nnqq/spider/pkg/spider"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("config.NewConfig: %w", err))
	}

	err = spider.NewSpider(
		list.NewReader(cfg.URLList),
		console.NewWriter(),
		cfg.Concurrency,
	).Run()
	if err != nil {
		log.Fatal(fmt.Errorf("spider.NewSpider.Run: %w", err))
	}
	log.Println("Spider finished")
}
