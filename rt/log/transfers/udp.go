package transfers

import (
	"fmt"
	"net"
)

// CreateUDPTransfer 创建 UDP 转发
func CreateUDPTransfer(cfg *TransferConfigure) TransferFn {
	server := fmt.Sprintf("%s:%d", cfg.Server, cfg.Port)

	return func(ch chan []byte) {
		conn, err := net.ListenUDP("udp", nil)
		if err != nil {
			panic(err)
		}

		to, err := net.ResolveUDPAddr("udp", server)
		if err != nil {
			panic(err)
		}

		for {
			data := <-ch
			_, err = conn.WriteTo(data, to)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("udp", string(data))
		}
	}
}
