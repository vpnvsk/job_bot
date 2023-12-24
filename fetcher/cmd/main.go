package main

import (
	"context"
	"errors"
	source "fetcher/internal/rss"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sources := []string{"https://jobs.dou.ua/vacancies/feeds/?exp=1-3&category=Golang",
		"https://jobs.dou.ua/vacancies/feeds/?exp=1-3&category=Python"}

	sour := source.NewSource(sources)
	fet := source.New(1*time.Minute, sour)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := fet.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to run fetcher: %v", err)
			return
		}

		log.Printf("[INFO] fetcher stopped")
	}

}

//import (
//	"context"
//	"fmt"
//	"github.com/SlyMarbo/rss"
//	"log"
//	"os"
//	"os/signal"
//	"sync"
//	"syscall"
//)
//
//func main() {
//	t := Fetcher{1}
//	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
//	defer cancel()
//	t.Fetch(ctx)
//}
//
//type Fetcher struct {
//	iii int
//}
//
//func (f *Fetcher) Fetch(ctx context.Context) error {
//
//	sources := []string{"https://justjoin.it/all-locations/python/experience-level_junior",
//		"https://justjoin.it/krakow/javascript/experience-level_junior",
//		"https://www.glassdoor.com/Job/poland-python-junior-jobs-SRCH_IL.0,6_IN193_KO7,20.htm",
//		"https://nofluffjobs.com/pl/praca-zdalna/backend?criteria=city%3Dkrakow%20requirement%3DPython%20%20seniority%3Dtrainee,junior&page=1"}
//
//	var wg sync.WaitGroup
//
//	for _, rss := range sources {
//		wg.Add(1)
//
//		go func(rss string) {
//			defer wg.Done()
//
//			err := RFetch(ctx, rss)
//			if err != nil {
//				log.Printf("[ERROR] failed to fetch items from rss %q", err)
//				return
//			}
//
//			//if err := f.processItems(ctx, rss, items); err != nil {
//			//	log.Printf("[ERROR] failed to process items from rss %q: %v", rss.Name(), err)
//			//	return
//			//}
//		}(rss)
//	}
//
//	wg.Wait()
//
//	return nil
//}
//
//func RFetch(ctx context.Context, url string) error {
//	feed, err := loadFeed(ctx, url)
//	if err != nil {
//		return err
//	}
//	fmt.Print(feed.Items)
//	return nil
//}
//func loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
//	const op = "loadFeed"
//	var (
//		feedCh = make(chan *rss.Feed)
//		errCh  = make(chan error)
//	)
//
//	go func() {
//		feed, err := rss.Fetch(url)
//		if err != nil {
//
//			errCh <- fmt.Errorf("%s: something :::%s", err, op)
//			return
//		}
//		feedCh <- feed
//	}()
//
//	select {
//	case <-ctx.Done():
//		return nil, ctx.Err()
//	case err := <-errCh:
//		return nil, err
//	case feed := <-feedCh:
//		return feed, nil
//	}
//}
