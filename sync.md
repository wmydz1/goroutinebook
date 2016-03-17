#同步
通道并非是用来取代锁,它们有各自不同的使用场景。channel 倾向于解决逻辑层次的并 发处理架构,而 mutex 则用来保护局部范围内的数据安全。

```
BenchmarkTest-4            1    4299047783 ns/op        3296 B/op          8 allocs/op
BenchmarkTestBlock-4      10     122825583 ns/op      401516 B/op          2 allocs/op
￼￼!156

```

标准库 sync 所提供的 Mutex、RWMutex 使用并不复杂,只几个地方需要注意。
将 Mutex 作为匿名字段时,相关方法必须实现为 pointer-receiver。否则会因复制的关系, 导致锁机制失效。

```Go
type data struct {
    sync.Mutex
}
func (d data) test(s string) {
    d.Lock()
    defer d.Unlock()
    for i := 0; i < 5; i++ {
        println(s, i)
        time.Sleep(time.Second)
    }
}
func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    var d data
    go func() {
        defer wg.Done()
        d.test("read")
    }()
    go func() {
        defer wg.Done()
        d.test("write")
    }()
    wg.Wait() 
}

```

输出:

```
write 0
read 0
read 1
write 1
write 2
read 2
read 3
write 3
write 4
read 4

```
锁失效,改为 *data 方法后正常。

应将 Mutex 粒度控制在最小范围内,及早释放。

```Go
// 
func doSomething() {
    m.Lock()
    url := cache["key"]
    http.Get(url)
    m.Unlock()
}
// 
func doSomething() {
    m.Lock()
    url := cache["key"]
    m.Unlock()
    http.Get(url)
}
// 
//  defer Get 

```

Mutex 不支持递归锁,即便在同一 goroutine 下也会导致死锁。

```
func main() {
    var m sync.Mutex
    m.Lock()
{
    m.Lock()
    m.Unlock()
}
    m.Unlock() 
}

```
输出:

fatal error: all goroutines are asleep - deadlock!

在设计并发安全类型时,千万注意此类问题。

```Go
￼type cache struct {
    sync.Mutex
    data []int 
}
func (c *cache) count() int {
    c.Lock()
    n := len(c.data)
    c.Unlock()
    return n 
}
func (c *cache) get() int {
    c.Lock()
    defer c.Unlock()
    var d int
    if n := c.count(); n > 0 { // count 重复锁定，造成死锁
    d = c.data[0]
        c.data = c.data[1:]
    }
    return d
}
func main() {
    c := cache{
        data: []int{1, 2, 3, 4},
    }
    println(c.get())
}
// count 

```
fatal error: all goroutines are asleep - deadlock!

1. 对性能要求较高时,应避免使用 defer Unlock。  
2. 读写并发时,用 RWMutex 性能会更好一些。 
3. 执行严格测试,尽可能打开数据竞争检查。

