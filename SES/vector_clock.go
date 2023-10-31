package SES

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//const (
//	INT_SIZE = 4
//)

type VectorClock struct {
	InstanceID int
	NInstance  int
	Vectors    []*LogicClock
}

func NewVectorClock(nInstance, instanceID int) *VectorClock {
	vectors := make([]*LogicClock, nInstance)
	for i := 0; i < nInstance; i++ {
		vectors[i] = NewLogicClock(nInstance, i, i == instanceID)
	}
	return &VectorClock{
		InstanceID: instanceID,
		NInstance:  nInstance,
		Vectors:    vectors,
	}
}

func (vc *VectorClock) String() string {
	result := fmt.Sprintf("(%d,%d)", vc.NInstance, vc.InstanceID)
	for _, vector := range vc.Vectors {
		result += fmt.Sprintf("\n%s", vector)
	}
	return result
}

func (vc *VectorClock) Serialize(packet []byte) []byte {
	data := bytes.Buffer{}
	err := binary.Write(&data, binary.BigEndian, int32(vc.InstanceID))
	if err != nil {
		fmt.Println("Lỗi:", err)

	}
	for i := 0; i < vc.NInstance; i++ {
		data.Write(vc.Vectors[i].Serialize())
	}

	return append(data.Bytes(), packet...)
}

func (vc *VectorClock) Deserialize(packet []byte) (*VectorClock, []byte) {
	dataSize := INT_SIZE * (vc.NInstance*vc.NInstance + 1)
	data, packet := packet[0:dataSize], packet[dataSize:]

	newInstanceID := int(binary.BigEndian.Uint32(data[:INT_SIZE]))
	newVectorClock := NewVectorClock(vc.NInstance, newInstanceID)
	data = data[INT_SIZE:]

	for i := 0; i < vc.NInstance; i++ {
		newVectorClock.Vectors[i] = newVectorClock.Vectors[i].Deserialize(data[i*INT_SIZE*vc.NInstance : (i+1)*INT_SIZE*vc.NInstance])
	}

	return newVectorClock, packet
}

func (vc *VectorClock) Increase() {
	vc.Vectors[vc.InstanceID].Increase()
}

func (vc *VectorClock) SelfMerge(sourceID, destinationID int) {
	vc.Vectors[destinationID].Merge(vc.Vectors[sourceID])
}

func (vc *VectorClock) Merge(sourceVectorClock *VectorClock, sourceID, destinationID int) {
	vc.Vectors[destinationID].Merge(sourceVectorClock.Vectors[sourceID])
}

func (vc *VectorClock) GetClock(index int) *LogicClock {
	return vc.Vectors[index]
}

//func main() {
// Testing
// for n := 2; n <= 100; n++ {
// 	k := rand.Intn(n)
// 	vc1 := NewVectorClock(n, k)
// 	for i := 0; i < n; i++ {
// 		for j := 0; j < n; j++ {
// 			vc1.vectors[i].clock[j] = rand.Intn(1000) - 500
// 		}
// 	}
// 	temp := vc1.Serialize([]byte{})
// 	A, _ := vc1.Deserialize(temp)
// 	B := vc1.vectors
// 	fmt.Println(A)
// 	fmt.Println(B)
// 	fmt.Println(A.String() == B.String())
// }

// vc1 := NewVectorClock(2, 0)
// vc1.vectors[0].clock[1] = 1000
// fmt.Println("INPUT")
// fmt.Println(vc1)
// temp := vc1.Serialize([]byte{})
// A, _ := vc1.Deserialize(temp)
// B := vc1.vectors
// fmt.Println(A.vectors)
// fmt.Println(B)
// fmt.Println("OUTPUT")
// fmt.Println(A.String() == B.String())
//}
