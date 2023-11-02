package vectorclock

import (
	"encoding/json"
)

type VectorClock struct {
	clock map[string]int64
}

func (c *VectorClock) GetClock() map[string]int64 {
	return c.clock
}

func NewVectorClock() *VectorClock {
	return &VectorClock{
		clock: make(map[string]int64),
	}
}

// TODO: Increment the vector clock for a specific peer
func (vc *VectorClock) Increment(peerID string) {
	if _, ok := vc.clock[peerID]; !ok {
		vc.clock[peerID] = 1
	} else {
		vc.clock[peerID]++
	}
}

func (vc *VectorClock) isLessThan(other *VectorClock) bool {
	for k, v1 := range vc.clock {
		v2, ok := other.clock[k]

		if !ok {
			continue
		}

		if v1 > v2 {
			return false
		}
	}

	return true
}

func (vc *VectorClock) isMoreThan(other *VectorClock) bool {
	for k, v1 := range vc.clock {
		v2, ok := other.clock[k]

		if !ok {
			continue
		}

		if v1 < v2 {
			return false
		}
	}

	return true
}

// Normalize the vector clock by ensuring it has all keys from another vector
func (vc *VectorClock) Normalize(other *VectorClock) {
	for k := range other.clock {
		if _, ok := vc.clock[k]; !ok {
			vc.clock[k] = 0
		}
	}
}

func (vc *VectorClock) Compare(other *VectorClock) int {
	// Normalize both vector clocks
	vc.Normalize(other)
	other.Normalize(vc)

	// A is before B
	if vc.isLessThan(other) {
		return -1
	}

	// A is after B
	if vc.isMoreThan(other) {
		return 1
	}

	// A and B are concurrent
	return 0
}

// Serialize the vector clock to a JSON string
func (vc *VectorClock) Serialize() (string, error) {
	data, err := json.Marshal(vc.clock)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Deserialize a JSON string to a vector clock
func (vc *VectorClock) Deserialize(serialized string) error {
	var data map[string]int64
	err := json.Unmarshal([]byte(serialized), &data)
	if err != nil {
		return err
	}
	vc.clock = data
	return nil
}

// Merge 2 vectors clock
func MergeClock(c1, c2 *VectorClock) *VectorClock {
	c1.Normalize(c2)
	c2.Normalize(c1)

	mergeClock := NewVectorClock()

	for k, v := range c1.clock {
		mergeClock.clock[k] = v
	}

	for k, v := range c2.clock {
		if v > mergeClock.clock[k] {
			mergeClock.clock[k] = v
		}
	}

	return mergeClock
}
