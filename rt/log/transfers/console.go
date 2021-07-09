package transfers

import "fmt"

func CreateConsoleTransfer(cfg *TransferConfigure) TransferFn {
	return func(ch chan []byte) {
		for {
			data := <-ch
			fmt.Println(string(data))
		}
	}
}
