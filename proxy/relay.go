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
	var closeOnce sync.Once
	closeBoth := func() {
		closeOnce.Do(func() {
			utils.CloseQuietly(left)
			utils.CloseQuietly(right)
		})
	}
	wg.Add(2)
	go relay(left, right, &wg, closeBoth)
	go relay(right, left, &wg, closeBoth)
	return &wg
}

func relay(src, dst net.Conn, wg *sync.WaitGroup, closeBoth func()) {
	defer func() {
		utils.Recover("proxy.relay")
	}()

	defer wg.Done()
	defer closeBoth()

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
