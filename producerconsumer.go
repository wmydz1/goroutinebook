package main
import (

    "fmt"
    "time"
)
/*
写一个Server所使用的代码更少，也更简单。写一个Server除了网络，另外就是并发，相对python等其它语言，Go对并发支持使得它有更好的性能。
Goroutine和channel是Go在“并发”方面两个核心feature。
Channel是goroutine之间进行通信的一种方式，它与Unix中的管道类似。
Channel声明：
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
例如：
var ch chan int
var ch1 chan<- int  //ch1只能写
var ch2 <-chan int  //ch2只能读
channel是类型相关的，也就是一个channel只能传递一种类型。例如，上面的ch只能传递int。
*/
func Producter(queue  chan <-int) {
    for i := 0; i<1000; i++ {
        queue <- i
    }
}
func Consumer(queue <-chan int) {
    for i := 0; i<1000; i++ {
        fmt.Println("recive: ",<-queue)
    }
}
func main() {
    queue := make(chan int, 1)
    go Producter(queue)
    go Consumer(queue)

    time.Sleep(time.Duration(19)*time.Second)
}