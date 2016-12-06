package server

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func init() {
	registerCommand("SummField", &SummField{})
}

// просуммированое поле
type SummField struct {
	lock     *sync.Mutex
	value    float64
	fieldNum int
	regExpr  *regexp.Regexp
}

func (c *SummField) New(expr string, args []string) (command, error) {
	regExpr, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	result := &SummField{}
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

func (c *SummField) Process(line string) {
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

func (c *SummField) Result() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	result := c.value
	c.value = 0.0
	return result
}
