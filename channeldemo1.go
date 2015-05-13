package main
import (

    "fmt"
    "runtime"
)

var over chan bool
//开启的goroutine的数目
var gcNum int


func worker(i int) {

    fmt.Println("I am serving ....", i)
    over <- true
}

func main() {
    //使用多核进行计算
    runtime.GOMAXPROCS(runtime.NumCPU())
    over =make(chan bool)
    gcNum=10000
    for i := 0; i<gcNum; i++ {
        go worker(i)
    }
    for i := 0; i<gcNum; i++ {
        <-over
    }

}
