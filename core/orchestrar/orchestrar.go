package orchestrar

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Orchestrer is the interface used to
type Orchestrer interface {
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

// New creates a new Orchestrar that can be used to manage API-Request channels.
func New() Orchestrer {
	o := &orchestrar{}
	return Orchestrer(o)
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

// Orchestrar is main struct of this package
type orchestrar struct {
	payloads []*payload

	sync.Mutex
}

func (o *orchestrar) Put(ctx context.Context, id uuid.UUID) chan interface{} {
	p := createPayload(ctx, id)
	o.payloads = append(o.payloads, p)
	return p.ResponseChan
}

// Pop a channel if it exists
func (o *orchestrar) Pop(id uuid.UUID) (chan interface{}, error) {
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

func (o *orchestrar) Delete(id uuid.UUID) error {
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
func (o orchestrar) GetSize() int {
	o.Lock()
	defer o.Unlock()
	return len(o.payloads)
}
