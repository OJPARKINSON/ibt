package ibt

import (
	"fmt"
	"reflect"
)

// Tick is a single instance of telemetry data
type Tick map[string]interface{}

// TickValueType is an interface containing all possible types for the value of a telemetry variable
type TickValueType interface {
	uint8 | []uint8 | bool | []bool | int | []int | string | []string | float32 | []float32 | float64 | []float64
}

// Filter the tick for only the given whitelisted fields
func (t Tick) Filter(whitelist ...string) Tick {
	// For now, use simple allocation until we can properly implement pooling
	// The pool approach was causing shared state issues in tests
	partialTick := make(Tick, len(whitelist))

	for _, field := range whitelist {
		if val, exists := t[field]; exists {
			partialTick[field] = val
		}
	}

	return partialTick
}

// GetTickValue will retrieve and type assert the given variable.
func GetTickValue[T TickValueType](tick Tick, key string) (T, error) {
	var def T

	rawValue, ok := tick[key]
	if !ok {
		return def, fmt.Errorf("key %s not found in tick", key)
	}

	value, ok := rawValue.(T)
	if !ok {
		return def, fmt.Errorf("value of %s was %s not %s", key, reflect.TypeOf(rawValue).String(), reflect.TypeOf(def).String())
	}

	return value, nil
}
