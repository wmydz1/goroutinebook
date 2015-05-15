package main
import (

    "sync"
    "fmt"
    "runtime"
)
var waitgroup sync.WaitGroup
func work(i int) {
    fmt.Println("I am working...",i)
    waitgroup.Done()
}
func main() {
    //使用多核进行计算
    runtime.GOMAXPROCS(runtime.NumCPU())
    for i := 0; i<1000; i++ {
        waitgroup.Add(1)
        go work(i)
    }
    waitgroup.Wait()
}

