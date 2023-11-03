package vectorclock

import (
	"encoding/json"
	"sync"
)

type ClockEntry struct {
	PeerID string
	Value  int64
}

type VectorClock struct {
	clock []ClockEntry
	mu    sync.Mutex
}

func NewVectorClock() *VectorClock {
	return &VectorClock{
		clock: make([]ClockEntry, 0),
	}
}

func (vc *VectorClock) GetClock() []ClockEntry {
	return vc.clock
}

func (vc *VectorClock) Increment(peerID string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	for i, entry := range vc.clock {
		if entry.PeerID == peerID {
			vc.clock[i].Value++
			return
		}
	}

	// If the peerID doesn't exist in the clock, add it with value 1
	vc.clock = append(vc.clock, ClockEntry{PeerID: peerID, Value: 1})
}

func (vc *VectorClock) Compare(other *VectorClock) int {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	for _, entry := range other.clock {
		found := false
		for _, ownEntry := range vc.clock {
			if ownEntry.PeerID == entry.PeerID {
				found = true
				if ownEntry.Value < entry.Value {
					return -1
				} else if ownEntry.Value > entry.Value {
					return 1
				}
			}
		}
		if !found {
			// Missing entry in the other clock
			return -1
		}
	}

	for _, ownEntry := range vc.clock {
		found := false
		for _, entry := range other.clock {
			if ownEntry.PeerID == entry.PeerID {
				found = true
			}
		}
		if !found {
			// Missing entry in own clock
			return 1
		}
	}

	return 0
}

func (vc *VectorClock) Serialize() (string, error) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	data, err := json.Marshal(vc.clock)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (vc *VectorClock) Deserialize(serialized string) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	var data []ClockEntry
	err := json.Unmarshal([]byte(serialized), &data)
	if err != nil {
		return err
	}
	vc.clock = data
	return nil
}

func MergeClock(c1, c2 *VectorClock) *VectorClock {
	c1.mu.Lock()
	defer c1.mu.Unlock()

	c2.mu.Lock()
	defer c2.mu.Unlock()

	mergedClock := NewVectorClock()

	for _, entry := range c1.clock {
		mergedClock.Increment(entry.PeerID)
	}
	for _, entry := range c2.clock {
		mergedClock.Increment(entry.PeerID)
	}

	return mergedClock
}
