package source

import (
	"context"
	"fetcher/models"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/samber/lo"
	"sync"
)

type Sources interface {
	Fetch(ctx context.Context) ([]models.Item, error)
}

type Source struct {
	Urls []string
}

func NewSource(urls []string) *Source {
	return &Source{Urls: urls}
}

type syncItem struct {
	mu   sync.Mutex
	item *[]models.Item
}

func (s *Source) Fetch(ctx context.Context) ([]models.Item, error) {
	var lst []models.Item
	si := &syncItem{item: &lst}
	var wg sync.WaitGroup
	errCh := make(chan error)
	for _, i := range s.Urls {
		wg.Add(1)
		i := i
		go func() {
			feed, err := s.loadFeed(ctx, i)
			if err != nil {
				errCh <- err
				return
			}
			fmt.Println(feed)
			g := lo.Map(feed.Items, func(item *rss.Item, _ int) models.Item {
				return models.Item{
					Title:       item.Title,
					Link:        item.Link,
					PubDate:     item.Date,
					Description: item.Content,
				}
			})
			si.mu.Lock()
			*si.item = append(*si.item, g[0])
			si.mu.Unlock()
			defer wg.Done()
		}()
	}

	wg.Wait()
	err := <-errCh
	if err != nil {
		return nil, err
	}
	return *si.item, err
}
func (s *Source) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)
	go func() {
		feed, err := rss.Fetch(url)

		if err != nil {
			errCh <- err
			return
		}
		feedCh <- feed
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}

}
