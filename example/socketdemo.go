package main
import (

    "net"
    "fmt"
    "os"
    "time"
)
const (
    MAX_CONN_NUM = 5
)
func EchoFunc(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    for {
        _, err := conn.Read(buf)
        if err!=nil {

            return
        }
        //send messsage
        _, err=conn.Write(buf)
        if err!=nil {
//            panic(err)
            return
        }
    }
}

func main() {

    listener, err := net.Listen("tcp", "0.0.0.0:8088")
    if err!=nil {
        fmt.Println("error listening", err.Error())
        os.Exit(1)
    }
    defer listener.Close()
    fmt.Println("running.....")


    var cur_conn_num int = 0
    conn_chan := make(chan net.Conn)
    ch_conn_change := make(chan int)

    go func() {
        for conn_change := range ch_conn_change {
            cur_conn_num+=conn_change
        }
    }()

    go func() {
        for _ = range time.Tick(1e8) {
            fmt.Printf("cur conn num %f\n", cur_conn_num)
        }
    }()

    for i := 0; i<MAX_CONN_NUM; i++ {
        go func() {
            for conn := range conn_chan {
                ch_conn_change <- 1
                EchoFunc(conn)
                ch_conn_change <- 1
            }
        }()
    }

    for {
        conn, err := listener.Accept()
        if err !=nil {
            fmt.Println("Error Accept", err.Error())
            return
        }
        conn_chan <- conn
    }



}
