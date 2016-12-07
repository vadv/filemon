package server

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func init() {
	registerCommand("AvgField", &AvgField{})
}

// просуммированое поле
type AvgField struct {
	lock     *sync.Mutex
	counter  int64
	value    float64
	fieldNum int
	regExpr  *regexp.Regexp
}

func (c *AvgField) New(expr string, args []string) (command, error) {
	regExpr, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	result := &AvgField{}
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
	return result, nil
}

func (c *AvgField) Process(line string) {
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

func (c *AvgField) Result() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	result := 0.0
	if c.counter > 0 {
		result = c.value / float64(c.counter)
	}
	c.value, c.counter = 0.0, 0
	return result
}
