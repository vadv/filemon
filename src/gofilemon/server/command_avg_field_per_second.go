package server

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	registerCommand("AvgFieldPerSecond", &AvgFieldPerSecond{})
}

type AvgFieldPerSecond struct {
	lock         *sync.Mutex
	counter      int64
	value        float64
	fieldNum     int
	regExpr      *regexp.Regexp
	lastResultTs int64 // unixts in milliseconds
}

func (c *AvgFieldPerSecond) New(expr string, args []string) (command, error) {
	regExpr, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	result := &AvgFieldPerSecond{}
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

func (c *AvgFieldPerSecond) Process(line string) {
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
		value = data[fieldNum]
	}
	value = strings.Trim(value, `"`)
	value = strings.Trim(value, ` `)

	// parse value
	if result, err := strconv.ParseFloat(value, 64); err != nil {
		return
	} else {
		// increment
		c.value += result
		c.counter++
	}
}

func (c *AvgFieldPerSecond) Result() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	now, result := time.Now().UnixNano(), 0.0
	if c.counter > 0 {
		result = float64(time.Second) * float64(c.value) / (float64(now-c.lastResultTs) * float64(c.counter))
	}
	c.value, c.counter = 0.0, 0
	c.lastResultTs = now

	return result
}
