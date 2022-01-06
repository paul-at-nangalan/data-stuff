package cache

import (
	"container/ring"
	"log"
	"sync"
)

///A simple, thread safe, size limit cache
/// The size must be greater than the max number of concurrent threads that
/// can set data in the cache
type FifoCache struct{
	maplock sync.RWMutex
	keys    *ring.Ring  ///store the keys in order, when we reach capacity, remove the first key added
	size int
	store   map[string]interface{}
}

func NewFifoCache(size int)*FifoCache{
	return &FifoCache{
		keys:  ring.New(size),
		size: size,
		store: make(map[string]interface{}),
	}
}

func (p *FifoCache)replaceFirst(newkey string){

	if len(p.store) > p.size {
		log.Panic("Store size has exceeded limit", p.size, len(p.store))
	}
	last := p.keys
	if len(p.store) == p.size {
		first := last
		///this should always be valid, cos we fill the cache before we start
		// removing items from it
		firstoutkey := first.Value.(string)
		//fmt.Println("Removing key ", firstoutkey)
		delete(p.store, firstoutkey)
	}
	last.Value = newkey
	p.keys = last.Next()
}

func (p *FifoCache)Set(key string, value interface{}){
	p.maplock.Lock()
	defer p.maplock.Unlock()

	////Don't fill the cache with the same item
	_, found := p.store[key]
	/// if it's not found, this is a new item, so replace the oldest item in the cache
	if !found {
		p.replaceFirst(key)
	}
	//fmt.Println("Putting ", key, " into the cache, with val ", value)
	p.store[key] = value
}

func (p *FifoCache)Find(key string)(value interface{}, found bool){
	p.maplock.RLock()
	defer p.maplock.RUnlock()

	value, found = p.store[key]
	return value, found
}
