package prices

import (
	"github.com/paul-at-nangalan/errorhandler/handlers"
	"sort"
	"strconv"
	"sync"
)

type Positions struct {
	/////map price to volume
	lock      sync.RWMutex
	position  map[float64]float64
	timestamp uint64
	///so we can report any positions that have been removed to who ever is using this class
	removed []float64
}

func NewPositions() *Positions {
	return &Positions{}
}

type PriceVol struct {
	Price float64
	Vol   float64
}

func (p *Positions) GetAllUnordered() []PriceVol {
	p.lock.RLock()
	defer p.lock.RUnlock()

	pricevol := make([]PriceVol, len(p.position)+len(p.removed))
	i := 0
	for price, vol := range p.position {
		pricevol[i].Price = price
		pricevol[i].Vol = vol
		i++
	}
	for _, removed := range p.removed {
		pricevol[i].Price = removed
		pricevol[i].Vol = 0
		i++
	}
	p.removed = p.removed[:0]
	return pricevol
}

func (p *Positions) GetAllOrderedByPrice() []PriceVol {
	pricevol := p.GetAllUnordered()
	sort.Slice(pricevol, func(i int, j int) bool {
		if pricevol[i].Price < pricevol[j].Price {
			return true
		}
		return false
	})
	return pricevol
}

func (p *Positions) Fill(pos map[string]interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.position = make(map[float64]float64)
	for pricestr, volstr := range pos {
		price := toFloat(pricestr)
		vol := toFloat(volstr.(string))
		//fmt.Println("Adding vol for ", price, vol)
		p.position[price] = vol
	}
}

func (p *Positions) Update(pos []interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for i := 0; i < len(pos); i += 3 { ///tripplets <price>, <vol>, <time>
		price := toFloat(pos[i].(string))
		vol := toFloat(pos[i+1].(string))
		if vol > 0 {
			p.position[price] = vol
		} else {
			delete(p.position, price)
			p.removed = append(p.removed, price)
		}
		ts, err := strconv.ParseUint(pos[i+2].(string), 10, 64)
		handlers.PanicOnError(err)
		p.timestamp = ts
	}
}

func toFloat(str string)float64{
	fl, err := strconv.ParseFloat(str, 64)
	handlers.PanicOnError(err)
	return fl
}
