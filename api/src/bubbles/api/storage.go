package api

import (
	"container/list"
	"errors"
	"log"
	"sync"
	"time"

	// "github.com/google/btree"

	"bubbles/config"
)

var (
	storage sync.Map
	gc      = make(chan *lifetime, 10)
	life    = list.New()
)

type lifetime struct {
	id     string
	expire time.Time
}

func store(id string, body []byte, ttl int) error {
	if _, loaded := storage.LoadOrStore(id, body); loaded {
		return errors.New("already exist")
	}
	if config.Debug {
		log.Printf("Store %s of %d bytes", id, len(body))
	}
	gc <- &lifetime{id, time.Now().Add(time.Duration(ttl) * time.Second)}
	return nil
}

func retrieve(id string) ([]byte, error) {
	if maybe, exist := storage.Load(id); !exist {
		return nil, errors.New("not found")
	} else {
		body := maybe.([]byte)
		if config.Debug {
			log.Printf("Retrieve %s of %d bytes", id, len(body))
		}
		return body, nil
	}
}

func Init() {
	go dispose()
}

func dispose() {
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-tick:
			now := time.Now()
			if config.Trace {
				log.Printf("Expiring < %d", now.Unix())
			}
			for elem := life.Front(); elem != nil && elem.Value.(*lifetime).expire.Before(now); {
				next := elem.Next()
				id := elem.Value.(*lifetime).id
				if config.Trace {
					log.Printf("Expiring %s", id)
				}
				storage.Delete(id)
				life.Remove(elem)
				elem = next
			}

		case l := <-gc:
			elem := life.Back()
			for elem != nil && elem.Value.(*lifetime).expire.After(l.expire) {
				elem = elem.Prev()
			}
			if config.Trace {
				log.Printf("Expire %s after %d", l.id, l.expire.Unix())
			}
			if elem == nil {
				life.PushFront(l)
			} else {
				life.InsertAfter(l, elem)
			}
		}
	}
}
