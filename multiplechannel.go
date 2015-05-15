package main
import (

    "time"
    "fmt"
)
/*
在实际编程中，经常会遇到在一个goroutine中处理多个channel的情况。我们不可能阻塞在两个channel，
这时就该select场了。与C语言中的select可以监控多个fd一样，go语言中select可以等待多个channel。
*/
func main() {
    c1 := make(chan string)
    c2 := make(chan string)
    go func() {
        time.Sleep(time.Second*1)
        c1 <- "one"

    }()
    go func() {
        time.Sleep(time.Second*2)
        c2 <- "two"

    }()
    for i := 0; i<2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("receved ", msg1)
        case msg2 := <-c2:
            fmt.Println("recived", msg2)
        }
    }
}