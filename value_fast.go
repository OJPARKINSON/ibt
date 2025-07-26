package ibt

import (
	"unsafe"
	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
)

// readVarValueFast is an optimized version of readVarValue using unsafe operations
// and pre-allocated slices for better performance
func readVarValueFast(buf []byte, vh headers.VarHeader) interface{} {
	offset := vh.Offset
	var value interface{}

	if vh.Count > 1 {
		switch vh.Rtype {
		case 0: // uint8 array
			if offset+vh.Count <= len(buf) {
				// Direct slice without allocation for single bytes
				result := make([]uint8, vh.Count)
				copy(result, buf[offset:offset+vh.Count])
				value = result
			} else {
				value = []uint8{}
			}
		case 1: // bool array
			res := boolPool.Get().([]bool)[:0]
			for i := 0; i < vh.Count && offset+i < len(buf); i++ {
				res = append(res, buf[offset+i] > 0)
			}
			result := make([]bool, len(res))
			copy(result, res)
			boolPool.Put(res)
			value = result
		case 2: // int array
			if offset+vh.Count*4 <= len(buf) {
				res := intPool.Get().([]int)[:0]
				for i := 0; i < vh.Count*4; i += 4 {
					intVal := fastByte4ToInt(buf[offset+i:])
					res = append(res, intVal)
				}
				result := make([]int, len(res))
				copy(result, res)
				intPool.Put(res)
				value = result
			} else {
				value = []int{}
			}
		case 3: // string/bitfield array
			res := stringPool.Get().([]string)[:0]
			for i := 0; i < vh.Count*4; i += 4 {
				if offset+i+4 <= len(buf) {
					res = append(res, utilities.Byte4toBitField(buf[offset+i:offset+i+4]))
				}
			}
			result := make([]string, len(res))
			copy(result, res)
			stringPool.Put(res)
			value = result
		case 4: // float32 array
			if offset+vh.Count*4 <= len(buf) {
				res := float32Pool.Get().([]float32)[:0]
				for i := 0; i < vh.Count*4; i += 4 {
					floatVal := fastByte4ToFloat(buf[offset+i:])
					res = append(res, floatVal)
				}
				result := make([]float32, len(res))
				copy(result, res)
				float32Pool.Put(res)
				value = result
			} else {
				value = []float32{}
			}
		case 5: // float64 array
			if offset+vh.Count*8 <= len(buf) {
				res := float64Pool.Get().([]float64)[:0]
				for i := 0; i < vh.Count*8; i += 8 {
					floatVal := fastByte8ToFloat(buf[offset+i:])
					res = append(res, floatVal)
				}
				result := make([]float64, len(res))
				copy(result, res)
				float64Pool.Put(res)
				value = result
			} else {
				value = []float64{}
			}
		}
	} else {
		// Single values - use unsafe for direct conversion
		switch vh.Rtype {
		case 0: // uint8
			if offset < len(buf) {
				value = uint8(buf[offset])
			} else {
				value = uint8(0)
			}
		case 1: // bool
			if offset < len(buf) {
				value = buf[offset] > 0
			} else {
				value = false
			}
		case 2: // int
			if offset+4 <= len(buf) {
				value = fastByte4ToInt(buf[offset:])
			} else {
				value = 0
			}
		case 3: // string/bitfield
			if offset+4 <= len(buf) {
				value = utilities.Byte4toBitField(buf[offset : offset+4])
			} else {
				value = "0x0"
			}
		case 4: // float32
			if offset+4 <= len(buf) {
				value = fastByte4ToFloat(buf[offset:])
			} else {
				value = float32(0)
			}
		case 5: // float64
			if offset+8 <= len(buf) {
				value = fastByte8ToFloat(buf[offset:])
			} else {
				value = float64(0)
			}
		}
	}

	return value
}

// Fast unsafe byte conversion functions
func fastByte4ToFloat(buf []byte) float32 {
	if len(buf) < 4 {
		return 0
	}
	return *(*float32)(unsafe.Pointer(&buf[0]))
}

func fastByte4ToInt(buf []byte) int {
	if len(buf) < 4 {
		return 0
	}
	converted := *(*uint32)(unsafe.Pointer(&buf[0]))
	if converted == 4294967295 {
		return -1
	}
	return int(converted)
}

func fastByte8ToFloat(buf []byte) float64 {
	if len(buf) < 8 {
		return 0
	}
	return *(*float64)(unsafe.Pointer(&buf[0]))
}