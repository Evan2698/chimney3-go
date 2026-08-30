package proxy

import (
	"chimney3-go/mem"
	"chimney3-go/utils"
	"log"
	"net"
	"sync"
)

func startBidirectionalRelay(left, right net.Conn) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(2)
	go relay(left, right, &wg)
	go relay(right, left, &wg)
	return &wg
}

func relay(src, dst net.Conn, wg *sync.WaitGroup) {
	defer func() {
		utils.Recover("proxy.relay")
	}()

	defer wg.Done()

	buf := mem.GetLarge()
	defer func() {
		mem.PutLarge(buf)
	}()
	for {
		n, err := src.Read(buf)
		if err != nil {
			log.Println("read failed ", err)
			break
		}

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Println("write failed ", err)
			break
		}
	}
}
