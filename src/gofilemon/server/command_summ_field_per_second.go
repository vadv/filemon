package server

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	registerCommand("SummFieldPerSecond", &SummFieldPerSecond{})
}

type SummFieldPerSecond struct {
	lock         *sync.Mutex
	value        float64
	fieldNum     int
	regExpr      *regexp.Regexp
	lastResultTs int64 // unixts in milliseconds
}

func (c *SummFieldPerSecond) New(expr string, args []string) (command, error) {
	regExpr, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	result := &SummFieldPerSecond{}
	// parse fieldNum
	if len(args) != 1 {
		return nil, ArgsError
	} else {
		if fieldNum, err := strconv.ParseInt(args[0], 10, 32); err != nil {
			return nil, err
		} else {
			result.fieldNum = int(fieldNum)
		}
	}
	result.regExpr = regExpr
	result.lock = &sync.Mutex{}
	result.lastResultTs = time.Now().UnixNano()
	return result, nil
}

func (c *SummFieldPerSecond) Process(line string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.regExpr.MatchString(line) {
		return
	}

	data, value := strings.Split(line, " "), ""

	if c.fieldNum > 0 {
		if len(data) < c.fieldNum {
			return
		}
		value = data[c.fieldNum]
	} else {
		fieldNum := len(data) + c.fieldNum
		if fieldNum < 0 {
			return
		}
		value = data[c.fieldNum]
	}

	// parse value
	if result, err := strconv.ParseFloat(value, 64); err != nil {
		return
	} else {
		// increment
		c.value += result
	}
}

func (c *SummFieldPerSecond) Result() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now().UnixNano()
	result := float64(time.Millisecond) * float64(c.value) / float64(now-c.lastResultTs)
	c.value = 0.0
	c.lastResultTs = now

	return result
}
