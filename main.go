package main

import (
	"SES_Algorithm/SES"
	"fmt"
)

func main() {
	p0 := SES.NewSES(3, 0)
	p1 := SES.NewSES(3, 1)
	p2 := SES.NewSES(3, 2)

	fmt.Println(p0)
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println("--------------")

	pk1 := p1.Send(0, []byte("pk1"))
	fmt.Println(p1)
	fmt.Println("--------------")

	pk2 := p1.Send(0, []byte("pk2"))
	fmt.Println(p1)
	fmt.Println("--------------")

	p0.Deliver(pk2)
	fmt.Println("**************")
	fmt.Println(p0)
	fmt.Println("--------------")
	p0.Deliver(pk1)
	fmt.Println("**************")
	fmt.Println(p0)
	fmt.Println("--------------")

	pk3 := p2.Send(0, []byte("pk3"))
	fmt.Println(p2)
	fmt.Println("--------------")

	pk4 := p2.Send(1, []byte("pk4"))
	fmt.Println(p2)
	fmt.Println("--------------")

	p1.Deliver(pk4)
	fmt.Println("**************")
	fmt.Println(p1)
	fmt.Println("--------------")

	pk5 := p1.Send(0, []byte("pk5"))
	fmt.Println(p1)
	fmt.Println("--------------")

	p0.Deliver(pk5)
	fmt.Println("**************")
	fmt.Println(p0)
	fmt.Println("--------------")

	p0.Deliver(pk3)
	fmt.Println("**************")
	fmt.Println(p0)
	fmt.Println("--------------")

	sendPacket := p1.Send(0, []byte{})
	fmt.Println(p1.Deserialize(sendPacket))
	fmt.Println(p1)
}
