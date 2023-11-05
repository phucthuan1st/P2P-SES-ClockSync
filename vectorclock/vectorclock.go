package vectorclock

import (
	"encoding/json"
	"sort"
	"sync"
)

type ClockEntry struct {
	NodeID string
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

func (vc *VectorClock) SetClock(clock []ClockEntry) {
	vc.clock = clock
}

func (vc *VectorClock) Increment(peerID string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	for i, entry := range vc.clock {
		if entry.NodeID == peerID {
			vc.clock[i].Value++
			return
		}
	}

	// If the peerID doesn't exist in the clock, add it with value 1
	vc.clock = append(vc.clock, ClockEntry{NodeID: peerID, Value: 1})
}

func (vc *VectorClock) Compare(other *VectorClock) int {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	vc.normalize(other)
	other.normalize(vc)

	// A == B
	if vc.areEqual(other) {
		return 0
	}

	// A < B
	if vc.isLessThan(other) {
		return -1
	}

	// A > B
	if other.isLessThan(vc) {
		return 1
	}

	// concurrently
	return -101
}

// Rule 1: Equal (ta = tb iff ta[i] = tb[i])
func (vc *VectorClock) areEqual(other *VectorClock) bool {
	for i := 0; i < len(vc.clock); i++ {
		if vc.clock[i].Value != other.clock[i].Value {
			return false
		}
	}
	return true
}

// Rule 3: Less than (ta < tb iff ta[i] <= tb[i] and ta[i] != tb[i], for all i)
func (vc *VectorClock) isLessThan(other *VectorClock) bool {
	for i := 0; i < len(vc.clock); i++ {
		if vc.clock[i].Value >= other.clock[i].Value {
			return false
		}
	}
	return true
}

func MergeClock(c1, c2 *VectorClock) *VectorClock {
	c1.mu.Lock()
	defer c1.mu.Unlock()

	c2.mu.Lock()
	defer c2.mu.Unlock()

	mergedClock := NewVectorClock()

	for _, entry := range c1.clock {
		mergedClock.Increment(entry.NodeID)
	}
	for _, entry := range c2.clock {
		mergedClock.Increment(entry.NodeID)
	}

	return mergedClock
}

// normalize the vector clock by ensuring it has all keys from another vector
func (vc *VectorClock) normalize(other *VectorClock) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	// Create a map to store entries from the current vector clock
	ownEntries := make(map[string]int64)
	for _, entry := range vc.clock {
		ownEntries[entry.NodeID] = entry.Value
	}

	for _, entry := range other.clock {
		if _, found := ownEntries[entry.NodeID]; !found {
			vc.clock = append(vc.clock, ClockEntry{
				NodeID: entry.NodeID,
				Value:  0,
			})
		}
	}

	// Sort the clock entries to ensure they are in ascending order by PeerID
	vc.sortClockEntries()
}

// Add a new function to sort the clock entries by PeerID
func (vc *VectorClock) sortClockEntries() {
	sort.Slice(vc.clock, func(i, j int) bool {
		return vc.clock[i].NodeID < vc.clock[j].NodeID
	})
}

// serialize current clock to json format
func (vc *VectorClock) Serialize() (string, error) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	data, err := json.Marshal(vc.clock)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// deserialize json string to clock object
func (vc *VectorClock) Deserialize(serialized string) error {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	var data []ClockEntry
	err := json.Unmarshal([]byte(serialized), &data)
	if err != nil {
		return err
	}
	vc.clock = data
	vc.sortClockEntries() // Ensure the clock entries are sorted after deserialization
	return nil
}

// Clone current clock to a new one
func (vc *VectorClock) Clone() *VectorClock {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	// Create a new VectorClock
	clone := NewVectorClock()

	// Copy the clock entries
	for _, entry := range vc.clock {
		clone.clock = append(clone.clock, ClockEntry{
			NodeID: entry.NodeID,
			Value:  entry.Value,
		})
	}

	return clone
}

// Merge timestamp to a vectorclokc vc
func (vc *VectorClock) Merge(Timestamp []ClockEntry) *VectorClock {
	this := vc.Clone()

	other := NewVectorClock()
	other.SetClock(Timestamp)

	this.normalize(other)
	other.normalize(this)

	// TODO: choose max value between vc and other for each entry
	for i := range this.clock {
		if other.clock[i].Value > this.clock[i].Value {
			this.clock[i].Value = other.clock[i].Value
		}
	}

	return this
}
