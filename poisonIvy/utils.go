// Special Thanks: https://github.com/dtgreene/ivy2
package poisonIvy

import "math"

func ParseBitRange(input int, size uint) int {
	var result int
	for i := uint(0); i < size; i++ {
		if (input>>i)&1 == 1 {
			result += int(math.Pow(2, float64(size-1-i)))
		}
	}
	return result
}

func ParseIncomingMessage(data []byte) ([]byte, []byte, uint16, byte) {
	if len(data) < 8 {
		return data, nil, 0, 0
	}
	payload := data[8:]
	ack := uint16(data[6]&255) | (uint16(data[5]&255) << 8)
	error := data[7] & 255
	return data, payload, ack, error
}
