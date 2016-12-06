package server

import (
	"regexp"
	"sync"
)

func init() {
	registerCommand("CountLines", &CountLines{})
}

// кол-во линий после последнего запроса
type CountLines struct {
	lock    *sync.Mutex
	counter int64
	regExpr *regexp.Regexp
}

func (c *CountLines) New(expr string, args []string) (command, error) {
	regExpr, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		return nil, ArgsError
	}
	result := &CountLines{}
	result.regExpr = regExpr
	result.lock = &sync.Mutex{}
	return result, nil
}

func (c *CountLines) Process(line string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.regExpr.MatchString(line) {
		c.counter++
	}
}

func (c *CountLines) Result() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	result := float64(c.counter)
	c.counter = 0
	return result
}
