// Special Thanks: https://github.com/dtgreene/ivy2
package poisonIvy

import (
	"bytes"
	"encoding/binary"
)

const StartCode uint16 = 17167

func GetBaseMessage(command uint16, flag1, flag2 bool) []byte {
	buf := new(bytes.Buffer)

	var b1 int16 = 1
	var b2 int8 = 32
	if flag1 {
		b1, b2 = -1, -1
	}

	binary.Write(buf, binary.BigEndian, StartCode)
	binary.Write(buf, binary.BigEndian, b1)
	binary.Write(buf, binary.BigEndian, b2)
	binary.Write(buf, binary.BigEndian, command)

	var f2 byte = 0
	if flag2 {
		f2 = 1
	}
	binary.Write(buf, binary.BigEndian, f2)

	// pad to 34 bytes
	padding := make([]byte, 34-buf.Len())
	buf.Write(padding)

	return buf.Bytes()
}

// GetPrintReadyMessage creates the specific buffer for starting a print
func GetPrintReadyMessage(length int) []byte {
	msg := GetBaseMessage(769, false, false)

	// manually pack the length into the specific offsets (8-13)
	msg[8] = byte((int32(-16777216) & int32(length)) >> 24)
	msg[9] = byte((16711680 & length) >> 16)
	msg[10] = byte((65280 & length) >> 8)
	msg[11] = byte(length & 255)
	msg[12] = 1 // b4
	msg[13] = 1 // b5 (flag false)

	return msg
}
