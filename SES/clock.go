package SES

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	INT_SIZE    = 4
	PORT_OFFSET = 60000
	MAX_MESSAGE = 150
)

type LogicClock struct {
	NInstance  int
	InstanceID int
	Clock      []int32
}

func NewLogicClock(nInstance, instanceID int, zeroFill bool) *LogicClock {
	clock := make([]int32, nInstance)
	if zeroFill {
		for i := range clock {
			clock[i] = 0
		}
	} else {
		for i := range clock {
			clock[i] = -1
		}
	}
	return &LogicClock{
		NInstance:  nInstance,
		InstanceID: instanceID,
		Clock:      clock,
	}
}

func (lc *LogicClock) String() string {
	return fmt.Sprintf("%v", lc.GetTime())
}

func (lc *LogicClock) Equal(other *LogicClock) bool {
	for i := 0; i < lc.NInstance; i++ {
		if lc.Clock[i] != other.GetTime()[i] {
			return false
		}
	}
	return true
}

func (lc *LogicClock) LessThan(other *LogicClock) bool {
	return lc.LessThanOrEqual(other) && !lc.Equal(other)
}

func (lc *LogicClock) LessThanOrEqual(other *LogicClock) bool {
	for i := 0; i < lc.NInstance; i++ {
		if lc.Clock[i] > other.GetTime()[i] {
			return false
		}
	}
	return true
}

func (lc *LogicClock) Serialize() []byte {
	/*data := make([]byte, INT_SIZE)
	for i := 0; i < lc.NInstance; i++ {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, uint32(lc.Clock[i]))
		if err != nil {
			// Xử lý lỗi nếu cần thiết
		}
		data = append(data, buf.Bytes()...)
	}

	return data*/
	data := bytes.Buffer{}
	for _, value := range lc.Clock {
		binary.Write(&data, binary.BigEndian, int32(value))
	}
	return data.Bytes()
}

func (lc *LogicClock) Deserialize(data []byte) *LogicClock {
	newClock := NewLogicClock(lc.NInstance, lc.InstanceID, false)
	for i := 0; i < lc.NInstance; i++ {
		newClock.Clock[i] = int32(binary.BigEndian.Uint32(data[INT_SIZE*i : INT_SIZE*(i+1)]))
	}
	return newClock
}

func (lc *LogicClock) GetTime() []int32 {
	return lc.Clock
}

func (lc *LogicClock) Increase() {
	lc.Clock[lc.InstanceID]++
}

func (lc *LogicClock) IsNull() bool {
	for _, value := range lc.Clock {
		if value == -1 {
			return true
		}
	}
	return false
}

func (lc *LogicClock) Merge(other *LogicClock) {
	if lc.IsNull() {
		for i := range lc.Clock {
			lc.Clock[i] = other.GetTime()[i]
		}
	} else {
		for i := range lc.Clock {
			lc.Clock[i] = max(lc.Clock[i], other.GetTime()[i])
		}
	}
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
