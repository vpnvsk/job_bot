package rss

import (
	"context"
	"encoding/xml"
	"fetcher/models"
	"fmt"
	"io"
	"net/http"
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
			si.mu.Lock()
			*si.item = append(*si.item, feed...)
			si.mu.Unlock()
			defer wg.Done()
			return
		}()
	}

	wg.Wait()
	//close(errCh)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	default:
		return *si.item, nil
	}
}

func (s *Source) loadFeed(ctx context.Context, url string) ([]models.Item, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request:", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching RSS feed: %s", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d\n", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body:", err)
	}
	fmt.Println("third")

	var rss models.RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling XML:", err)
	}

	//feedCh <- rss.Channel.Items
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return rss.Channel.Items, err
	}
}

//	func (s *Source) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
//		var (
//			feedCh = make(chan *rss.Feed)
//			errCh  = make(chan error)
//		)
//		go func() {
//			ht := &http.Client{}
//			feed, err := rss.FetchByClient(url, ht)
//
//			if err != nil {
//				errCh <- err
//				return
//			}
//			feedCh <- feed
//		}()
//		select {
//		case <-ctx.Done():
//			return nil, ctx.Err()
//		case err := <-errCh:
//			return nil, err
//		case feed := <-feedCh:
//			return feed, nil
//		}
//
// }
