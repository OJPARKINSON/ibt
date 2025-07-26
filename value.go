package ibt

import (
	"sync"
	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
)

// Slice pools to reduce allocations for array variables
var (
	uint8Pool = sync.Pool{New: func() interface{} { return make([]uint8, 0, 32) }}
	boolPool = sync.Pool{New: func() interface{} { return make([]bool, 0, 32) }}
	intPool = sync.Pool{New: func() interface{} { return make([]int, 0, 32) }}
	stringPool = sync.Pool{New: func() interface{} { return make([]string, 0, 32) }}
	float32Pool = sync.Pool{New: func() interface{} { return make([]float32, 0, 32) }}
	float64Pool = sync.Pool{New: func() interface{} { return make([]float64, 0, 32) }}
)

// readVarValue extracts the telemetry variable value from the given buffer based on the provided metadata.
//
// This function will ensure that the underlying type of the value is correct.
func readVarValue(buf []byte, vh headers.VarHeader) interface{} {
	var rbuf []byte

	offset := vh.Offset
	var value interface{}

	if vh.Count > 1 {
		switch vh.Rtype {
		case 0:
			rbuf = buf[offset:vh.Count]
			res := uint8Pool.Get().([]uint8)[:0] // Reuse slice from pool
			for _, x := range rbuf[offset : offset+vh.Count] {
				res = append(res, uint8(x))
			}
			// Make a copy to return, put slice back in pool for reuse
			result := make([]uint8, len(res))
			copy(result, res)
			uint8Pool.Put(res)
			value = result
		case 1:
			rbuf = buf[offset : offset+vh.Count]
			res := boolPool.Get().([]bool)[:0] // Reuse slice from pool
			for _, x := range rbuf {
				res = append(res, x > 0)
			}
			// Make a copy to return, put slice back in pool for reuse
			result := make([]bool, len(res))
			copy(result, res)
			boolPool.Put(res)
			value = result
		case 2:
			rbuf = buf[offset : offset+vh.Count*4]
			res := intPool.Get().([]int)[:0] // Reuse slice from pool
			for i := 0; i < len(rbuf); i += 4 {
				res = append(res, utilities.Byte4ToInt(rbuf[i:i+4]))
			}
			// Make a copy to return, put slice back in pool for reuse
			result := make([]int, len(res))
			copy(result, res)
			intPool.Put(res)
			value = result
		case 3:
			rbuf = buf[offset : offset+vh.Count*4]
			res := stringPool.Get().([]string)[:0] // Reuse slice from pool
			for i := 0; i < len(rbuf); i += 4 {
				res = append(res, utilities.Byte4toBitField(rbuf[i:i+4]))
			}
			// Make a copy to return, put slice back in pool for reuse
			result := make([]string, len(res))
			copy(result, res)
			stringPool.Put(res)
			value = result
		case 4:
			rbuf = buf[offset : offset+vh.Count*4]
			res := float32Pool.Get().([]float32)[:0] // Reuse slice from pool
			for i := 0; i < len(rbuf); i += 4 {
				res = append(res, utilities.Byte4ToFloat(rbuf[i:i+4]))
			}
			// Make a copy to return, put slice back in pool for reuse
			result := make([]float32, len(res))
			copy(result, res)
			float32Pool.Put(res)
			value = result
		case 5:
			rbuf = buf[offset : offset+vh.Count*8]
			res := float64Pool.Get().([]float64)[:0] // Reuse slice from pool
			for i := 0; i < len(rbuf); i += 8 {
				res = append(res, utilities.Byte8ToFloat(rbuf[i:i+8]))
			}
			// Make a copy to return, put slice back in pool for reuse
			result := make([]float64, len(res))
			copy(result, res)
			float64Pool.Put(res)
			value = result
		}
	} else {
		switch vh.Rtype {
		case 0:
			rbuf = buf[offset : offset+1]
			value = uint8(rbuf[0])
		case 1:
			rbuf = buf[offset : offset+1]
			value = int(rbuf[0]) > 0
		case 2:
			rbuf = buf[offset : offset+4]
			value = utilities.Byte4ToInt(rbuf)
		case 3:
			rbuf = buf[offset : offset+4]
			value = utilities.Byte4toBitField(rbuf)
		case 4:
			rbuf = buf[offset : offset+4]
			value = utilities.Byte4ToFloat(rbuf)
		case 5:
			rbuf = buf[offset : offset+8]
			value = utilities.Byte8ToFloat(rbuf)
		}
	}

	return value
}
