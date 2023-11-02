package test

import (
	"p2p-ses-clocksync/vectorclock"
	"testing"
)

func TestVectorClock(t *testing.T) {
	// Create two vector clocks
	vcA := vectorclock.NewVectorClock()
	vcB := vectorclock.NewVectorClock()

	// Increment vector clock A for peer X
	vcA.Increment("X")
	if vcA.GetClock()["X"] != 1 {
		t.Errorf("Expected clock value for peer X to be 1, got %d", vcA.GetClock()["X"])
	}

	// Increment vector clock B for peer Y
	vcB.Increment("Y")
	if vcB.GetClock()["Y"] != 1 {
		t.Errorf("Expected clock value for peer Y to be 1, got %d", vcB.GetClock()["Y"])
	}

	// Compare A and B, A should be before B
	result := vcA.Compare(vcB)
	if result != 0 {
		t.Errorf("Expected A to be concurrent to B, got %d", result)
	}

	// Serialize vector clock A
	serializedA, err := vcA.Serialize()
	if err != nil {
		t.Errorf("Serialization error: %v", err)
	}

	// Deserialize vector clock A
	newVCA := vectorclock.NewVectorClock()
	err = newVCA.Deserialize(serializedA)
	if err != nil {
		t.Errorf("Deserialization error: %v", err)
	}

	// Check if the deserialized vector clock is the same as the original
	if !isEqualVectorClock(newVCA, vcA) {
		t.Errorf("Deserialized vector clock is not equal to the original")
	}
}

// Helper function to compare two vector clocks
func isEqualVectorClock(vc1, vc2 *vectorclock.VectorClock) bool {
	if len(vc1.GetClock()) != len(vc2.GetClock()) {
		return false
	}
	for k, v := range vc1.GetClock() {
		if v != vc2.GetClock()[k] {
			return false
		}
	}
	return true
}
