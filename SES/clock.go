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
	Clock      []int
}

func NewLogicClock(nInstance, instanceID int, zeroFill bool) *LogicClock {
	clock := make([]int, nInstance)
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
	data := bytes.Buffer{}
	for _, value := range lc.Clock {
		binary.Write(&data, binary.BigEndian, value)
	}
	return data.Bytes()
}

func (lc *LogicClock) Deserialize(data []byte) *LogicClock {
	newClock := NewLogicClock(lc.NInstance, lc.InstanceID, false)
	for i := 0; i < lc.NInstance; i++ {
		value := 0
		binary.Read(bytes.NewReader(data[i*INT_SIZE:(i+1)*INT_SIZE]), binary.BigEndian, &value)
		newClock.Clock[i] = value
	}
	return newClock
}

func (lc *LogicClock) GetTime() []int {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
