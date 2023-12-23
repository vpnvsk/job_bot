package source

import (
	"context"
	"fmt"
	"time"
)

type Fetcher interface {
	Sources
}
type Fetch struct {
	fetchInterval time.Duration
	Sources       Sources
}

func New(fetchInterval time.Duration, sources *Source) *Fetch {
	return &Fetch{
		fetchInterval: fetchInterval,
		Sources:       sources,
	}
}

func (f Fetch) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()
	if err := f.Fetch(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f Fetch) Fetch(ctx context.Context) error {
	items, err := f.Sources.Fetch(ctx)
	if err != nil {
		return err
	}
	fmt.Println(items)
	return err
}
