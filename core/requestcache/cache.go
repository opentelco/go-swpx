package requestcache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Cache is a Cache for requests between SWPX
// and the DNC dispatcher
type Cache interface {
	Put(ctx context.Context, id uuid.UUID) chan interface{}
	Pop(id uuid.UUID) (chan interface{}, error)
	Delete(id uuid.UUID) error
	GetSize() int
}

func SuccessFn(p *payload) error {
	return nil
}
func TimeoutFn() error {
	return nil
}

// New creates a new cache that can be used to manage API-Request channels.
func New() Cache {
	o := &requestCache{}
	return Cache(o)
}

type payload struct {
	ID           uuid.UUID
	ResponseChan chan interface{}
	ctx          context.Context
}

// Create a payload from ID and context. returns a channel to listen on.
func createPayload(ctx context.Context, id uuid.UUID) *payload {
	p := &payload{
		ID:           id,
		ResponseChan: make(chan interface{}, 1),
		ctx:          ctx,
	}

	return p
}

// requestCache is the impl
type requestCache struct {
	payloads []*payload

	sync.Mutex
}

func (o *requestCache) Put(ctx context.Context, id uuid.UUID) chan interface{} {
	p := createPayload(ctx, id)
	o.payloads = append(o.payloads, p)
	return p.ResponseChan
}

// Pop a channel if it exists
func (o *requestCache) Pop(id uuid.UUID) (chan interface{}, error) {
	o.Lock()
	defer o.Unlock()

	for ix, p := range o.payloads {
		o.payloads = append(o.payloads[:ix], o.payloads[ix+1:]...)
		if p.ID == id {
			return p.ResponseChan, nil
		}
	}
	return nil, fmt.Errorf("failed to pop request, id: %s does not exist", id.String())
}

func (o *requestCache) Delete(id uuid.UUID) error {
	o.Lock()
	defer o.Unlock()
	for ix, p := range o.payloads {
		if p.ID == id {
			log.Println("delete the selected id: ", id)
			o.payloads = append(o.payloads[:ix], o.payloads[ix+1:]...)
			return nil
		}
	}
	return fmt.Errorf("failed to delete, id: %s does not exist", id.String())
}

// GetSize returns the size of the request cache
func (o requestCache) GetSize() int {
	o.Lock()
	defer o.Unlock()
	return len(o.payloads)
}
