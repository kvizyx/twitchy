package helix

import "context"

const (
	defaultPageLimit = 100
	maxPageLimit     = 10000
)

type pagerConfig struct {
	pageLimit int
}

type PagerOption func(*pagerConfig)

func WithPageLimit(limit int) PagerOption {
	return func(config *pagerConfig) {
		config.pageLimit = limit
	}
}

type pageFetcher[T any] func(context.Context, string) (*Response[T], error)

type Pager[T any] struct {
	fetch      pageFetcher[T]
	pageLimit  int
	pageCount  int
	nextCursor string
	seen       map[string]struct{}
	page       *Response[T]
	err        error
	exhausted  bool
	capPending bool
}

func newPager[T any](fetch pageFetcher[T], options ...PagerOption) (*Pager[T], error) {
	return newPagerAt(fetch, "", options...)
}

func newPagerAt[T any](fetch pageFetcher[T], initialCursor string, options ...PagerOption) (*Pager[T], error) {
	config := pagerConfig{pageLimit: defaultPageLimit}
	for _, option := range options {
		if option == nil {
			return nil, ErrInvalidOption
		}
		option(&config)
	}
	if config.pageLimit < 1 || config.pageLimit > maxPageLimit {
		return nil, ErrInvalidOption
	}
	pager := &Pager[T]{
		fetch:      fetch,
		pageLimit:  config.pageLimit,
		nextCursor: initialCursor,
		seen:       make(map[string]struct{}),
	}
	if initialCursor != "" {
		pager.seen[initialCursor] = struct{}{}
	}
	return pager, nil
}

func (p *Pager[T]) Next(ctx context.Context) bool {
	if p.exhausted {
		return false
	}
	if p.capPending {
		p.capPending = false
		p.exhausted = true
		p.err = ErrPageLimit
		return false
	}
	if err := ctx.Err(); err != nil {
		p.exhausted = true
		p.err = err
		return false
	}

	page, err := p.fetch(ctx, p.nextCursor)
	if err != nil {
		p.exhausted = true
		if contextErr := ctx.Err(); contextErr != nil {
			p.err = contextErr
		} else {
			p.err = err
		}
		return false
	}
	cursor := page.Pagination.Cursor()
	if cursor != "" {
		if _, exists := p.seen[cursor]; exists {
			p.exhausted = true
			p.err = ErrCursorCycle
			return false
		}
	}

	p.page = page
	p.pageCount++
	if cursor == "" {
		p.exhausted = true
		return true
	}
	p.seen[cursor] = struct{}{}
	p.nextCursor = cursor
	if p.pageCount >= p.pageLimit {
		p.capPending = true
	}
	return true
}

func (p *Pager[T]) Page() *Response[T] {
	return p.page
}

func (p *Pager[T]) Err() error {
	return p.err
}
