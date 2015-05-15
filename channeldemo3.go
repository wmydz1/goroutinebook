package main
import (

    "fmt"
)
/*
对于无缓存的channel,放入channel和从channel中向外面取数据这两个操作不能放在同一个协程中，
防止死锁的发生；同时应该先利用go 开一个协程对channel进行操作，此时阻塞该go 协程，然后再在主协程中进行channel的相反操作
（与go 协程对channel进行相反的操作），实现go 协程解锁．即必须go协程在前，解锁协程在后．
带缓存channel:
对于带缓存channel，只要channel中缓存不满，则可以一直向 channel中存入数据，直到缓存已满；
同理只要channel中缓存不为０，便可以一直从channel中向外取数据，直到channel缓存变为０才会阻塞．
由此可见，相对于不带缓存channel，带缓存channel不易造成死锁，可以同时在一个goroutine中放心使用，
close主要用来关闭channel通道其用法为close(channel)，并且实在生产者的地方关闭channel，而不是在
消费者的地方关闭．并且关闭channel后，便不可再想channel中继续存入数据，但是可以继续从channel中读取数据
*/
func main() {
    ch := make(chan int, 20)
    for i := 0; i<10; i++ {
        ch <- i
    }
    close(ch)
    for i := range ch {
        fmt.Println(i)
    }



}
