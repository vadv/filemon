package server

import (
	"regexp"
	"sync"
	"time"
)

func init() {
	registerCommand("CountLinesPerSecond", &CountLinesPerSecond{})
}

// количество линий в секунду
type CountLinesPerSecond struct {
	CountLines         // тот же счетчик, только делим на время
	lastResultTs int64 // unixts in milliseconds
}

func (c *CountLinesPerSecond) New(expr string, args []string) (command, error) {
	regExpr, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		return nil, ArgsError
	}
	result := &CountLinesPerSecond{}
	result.regExpr = regExpr
	result.lock = &sync.Mutex{}
	result.lastResultTs = time.Now().UnixNano()
	return result, nil
}

func (c *CountLinesPerSecond) Process(line string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.regExpr.MatchString(line) {
		c.counter++
	}
}

func (c *CountLinesPerSecond) Result() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now().UnixNano()
	result := float64(time.Millisecond) * float64(c.counter) / float64(now-c.lastResultTs)
	c.counter = 0
	c.lastResultTs = now

	return result
}
