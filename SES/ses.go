package SES

import (
	"context"
	"fmt"
	"io"
	"log"

	"golang.org/x/sync/semaphore"
)

var loggerSend = log.New(io.Writer(nil), "__sender_log__", log.LstdFlags)
var loggerReceive = log.New(io.Writer(nil), "__receiver_log__", log.LstdFlags)

type SES struct {
	VectorClock *VectorClock
	Queue       []interface{}
	Lock        *semaphore.Weighted
}

func NewSES(nInstance int, instanceID int) *SES {
	vectorClock := NewVectorClock(nInstance, instanceID)
	queue := make([]interface{}, 0)
	var sem = semaphore.NewWeighted(int64(1))
	return &SES{
		VectorClock: vectorClock,
		Queue:       queue,
		Lock:        sem,
	}
}

func (s *SES) String() string {
	return fmt.Sprintf("%v\n%v", s.VectorClock, s.Queue)
}

func (s *SES) GetSenderLog(destinationID int, packet []byte) string {

	log := fmt.Sprintf("Send Packet Info:\n")
	log += fmt.Sprintf("\tSender ID: %d\n", s.VectorClock.InstanceID)
	log += fmt.Sprintf("\tReceiver ID: %d\n", destinationID)
	log += fmt.Sprintf("\tPacket Content: %s\n", string(packet))
	log += fmt.Sprintf("\tSender Clock:\n")
	log += fmt.Sprintf("\t\tLocal logical clock: %d\n", s.VectorClock.GetClock(s.VectorClock.InstanceID))
	log += fmt.Sprintf("\t\tLocal process vectors:\n")
	for i := 0; i < s.VectorClock.NInstance; i++ {
		if i != s.VectorClock.InstanceID && s.VectorClock.GetClock(i).IsNull() {
			log += fmt.Sprintf("\t\t\t<P_%d: %d>\n", i, s.VectorClock.GetClock(i))
		}
	}
	return log
}

func (s *SES) GetDeliverLog(tM *LogicClock, sourceVectorClock *VectorClock, packet []byte, status, header string, printCompare bool) string {
	//s.Lock.Lock()
	//defer s.Lock.Unlock()

	log := fmt.Sprintf("Received Packet Info %s:\n", header)
	log += fmt.Sprintf("\tSender ID: %d\n", sourceVectorClock.InstanceID)
	log += fmt.Sprintf("\tReceiver ID: %d\n", s.VectorClock.InstanceID)
	log += fmt.Sprintf("\tPacket Content: %s\n", string(packet))
	log += fmt.Sprintf("\tPacket Clock:\n")
	log += fmt.Sprintf("\t\tt_m: %d\n", tM)
	log += fmt.Sprintf("\t\ttP_snd: %d\n", sourceVectorClock.GetClock(s.VectorClock.InstanceID))
	log += fmt.Sprintf("\tReceiver Logical Clock (tP_rcv):\n")
	log += fmt.Sprintf("\t\t%d\n", s.VectorClock.GetClock(s.VectorClock.InstanceID))
	log += fmt.Sprintf("\tStatus: %s\n", status)
	if printCompare {
		log += fmt.Sprintf("\tDelivery Condition: %d > %d\n", s.VectorClock.GetClock(s.VectorClock.InstanceID), tM)
	}
	return log
}

func (s *SES) Serialize(packet []byte) []byte {

	return s.VectorClock.Serialize(packet)
}

func (s *SES) Deserialize(packet []byte) (*VectorClock, []byte) {

	vector_clock, packet := s.VectorClock.Deserialize(packet)
	return vector_clock, packet
}

func (ses *SES) Merge(sourceVectorClock *VectorClock) {
	for i := 0; i < ses.VectorClock.NInstance; i++ {
		if i != ses.VectorClock.InstanceID && i != sourceVectorClock.InstanceID {
			ses.VectorClock.Merge(sourceVectorClock, i, i)
		}
	}
	ses.VectorClock.Merge(sourceVectorClock, sourceVectorClock.InstanceID, ses.VectorClock.InstanceID)
	ses.VectorClock.Increase()
}

func (ses *SES) Deliver(packet []byte) {
	ctx := context.TODO()
	ses.Lock.Acquire(ctx, 1)
	sourceVectorClock, packet := ses.Deserialize(packet)
	fmt.Println("source_vector_clock\n", sourceVectorClock)
	fmt.Println("packet\n", packet)
	tP := ses.VectorClock.GetClock(ses.VectorClock.InstanceID)
	tM := sourceVectorClock.GetClock(ses.VectorClock.InstanceID)
	fmt.Println("Cai t_p\n", tP)
	fmt.Println("Cai t_m\n", tM)
	if tM.LessThanOrEqual(tP) {
		// Deliver
		//fmt.Printf(ses.GetDeliverLog(tM, sourceVectorClock, packet, "delivering", "BEFORE DELIVERED", true))
		ses.Merge(sourceVectorClock)
	} else {
		// Queue
		ses.Queue = append(ses.Queue, []interface{}{tM, sourceVectorClock, packet})
		//fmt.Printf(ses.GetDeliverLog(tM, sourceVectorClock, packet, "buffered", "BEFORE DELIVERED", true))
		breakFlag := false
		for !breakFlag {
			breakFlag = true
			for index, item := range ses.Queue {
				elements := item.([]interface{})
				tM := elements[0].(*LogicClock)
				sourceVectorClock := elements[1].(*VectorClock)
				//packet := item[2].([]byte)
				if tM.LessThanOrEqual(tM) {
					//fmt.Println(ses.GetDeliverLog(tM, sourceVectorClock, packet, "delivering from buffer", "BEFORE DELIVERED FROM BUFFERED", true))
					ses.Merge(sourceVectorClock)
					ses.Queue = append(ses.Queue[:index], ses.Queue[index+1:]...)
					breakFlag = false
					break
				}
			}
		}
	}
	ses.Lock.Release(1)
}

func (s *SES) Send(destinationID int, packet []byte) []byte {
	ctx := context.TODO()
	s.Lock.Acquire(ctx, 1)
	s.VectorClock.Increase()
	//fmt.Println(s.GetSenderLog(destinationID, packet))
	result := s.Serialize(packet)
	s.VectorClock.SelfMerge(s.VectorClock.InstanceID, destinationID)
	s.Lock.Release(1)
	return result
}
