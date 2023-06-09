package goroutinechan

import (
	"fmt"
	"testing"
	"time"
)

func say(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // 把 sum 发送到通道 c
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func Test_main(t *testing.T) {
	/*
		Go 语言支持并发，我们只需要通过 go 关键字来开启 goroutine 即可。
		goroutine 是轻量级线程，goroutine 的调度是由 Golang 运行时进行管理的。
		goroutine 语法格式：
			go funcname(params)
		或者
			go func(params){
				todo
			}()
	*/
	go say("world")
	say("hello")
	go func() {
		fmt.Println("hello world")
	}()

	/*
		通道（channel）是用来传递数据的一个数据结构。
		通道可用于两个 goroutine 之间通过传递一个指定类型的值来同步运行和通讯。操作符 <- 用于指定通道的方向，发送或接收。如果未指定方向，则为双向通道。
		ch <- v    // 把 v 发送到通道 ch
		v := <-ch  // 从 ch 接收数据并把值赋给 v
	*/

	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // 从通道 c 中接收

	fmt.Println(x, y, x+y)

	/*
		通道可以设置缓冲区，通过 make 的第二个参数指定缓冲区大小
		不带缓冲区的chan必须同时收发（类似于在线传输）
		带缓冲区的通道允许发送端的数据发送和接收端的数据获取处于异步状态（类似于离线传输）
		就是说发送端发送的数据可以放在缓冲区里面，可以等待接收端去获取数据，而不是立刻需要接收端去获取数据。
		不过由于缓冲区的大小是有限的，所以还是必须有接收端来接收数据的，否则缓冲区一满，数据发送端就无法再发送数据了。
	*/

	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	// 获取这两个数据
	fmt.Println(<-ch)
	fmt.Println(<-ch)

	// channel 可以关闭close
	c2 := make(chan int, 10)
	go fibonacci(cap(c2), c2)
	/*
		range 函数遍历每个从通道接收到的数据，因为 c 在发送完 10 个
		数据之后就关闭了通道，所以这里我们 range 函数在接收到 10 个数据
		之后就结束了。如果上面的 c 通道不关闭，那么 range 函数就不
		会结束，从而在接收第 11 个数据的时候就阻塞了。
	*/
	for i := range c2 {
		fmt.Println(i)
	}
}
