#并发的非阻塞缓存

本节中我们会做一个无阻塞的缓存，这种工具可以帮助我们来解决现实世界中并发程序出现但没有现成的库可以解决的问题。这个问题叫作缓存(memoizing)函数(译注：Memoization的定义： memoization 一词是Donald Michie 根据拉丁语memorandum杜撰的一个词。相应的动词、过去分词、ing形式有memoiz、memoized、memoizing.)，也就是说，我们需要缓存函数的返回结果，这样在对函数进行调用的时候，我们就只需要一次计算，之后只要返回计算的结果就可以了。我们的解决方案会是并发安全且会避免对整个缓存加锁而导致所有操作都去争一个锁的设计。

我们将使用下面的httpGetBody函数作为我们需要缓存的函数的一个样例。这个函数会去进行HTTP GET请求并且获取http响应body。对这个函数的调用本身开销是比较大的，所以我们
尽量尽量避免在不必要的时候反复调用。

```
func httpGetBody(url string) (interface{}, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

```
最后一行稍微隐藏了一些细节。ReadAll会返回两个结果，一个[]byte数组和一个错误，不过这两个对象可以被赋值给httpGetBody的返回声明里的interface{}和error类型，所以我们也就可以这样返回结果并且不需要额外的工作了。我们在httpGetBody中选用这种返回类型是为了使其可以与缓存匹配。”

下面是我们要设计的cache的第一个“草稿”：”

```
// Package memo provides a concurrency-unsafe
// memoization of a function of type Func.

package memo

// A Memo caches the results of calling a Func.
type Memo struct {
    f     Func
    cache map[string]result
}

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

type result struct {
    value interface{}
    err   error
}

func New(f Func) *Memo {
    return &Memo{f: f, cache: make(map[string]result)}
}

// NOTE: not concurrency-safe!
func (memo *Memo) Get(key string) (interface{}, error) {
    res, ok := memo.cache[key]
    if !ok {
        res.value, res.err = memo.f(key)
        memo.cache[key] = res
    }
    return res.value, res.err
}

```
Memo实例会记录需要缓存的函数f(类型为Func)，以及缓存内容(里面是一个string到result映射的map)。每一个result都是都是简单的函数返回的值对儿--一个值和一个错误值。继续下去我们会展示一些Memo的变种，不过所有的例子都会遵循这些上面的这些方面。

下面是一个使用Memo的例子。对于流入的URL的每一个元素我们都会调用Get，并打印调用延时以及其返回的数据大小的log：

```
m := memo.New(httpGetBody)
for url := range incomingURLs() {
    start := time.Now()
    value, err := m.Get(url)
    if err != nil {
        log.Print(err)
    }
    fmt.Printf("%s, %s, %d bytes\n",
    url, time.Since(start), len(value.([]byte)))
}

```
从下面的测试输出，我们可以看到URL流包含了一些重复的情况，尽管我们第一次对每一个URL的(*Memo).Get的调用都会花上几百毫秒，但第二次就只需要花1毫秒就可以返回完整的数据了。

$ go test -v gopl.io/ch9/memo1
=== RUN   Test
https://golang.org, 175.026418ms, 7537 bytes
https://godoc.org, 172.686825ms, 6878 bytes
https://play.golang.org, 115.762377ms, 5767 bytes
http://gopl.io, 749.887242ms, 2856 bytes
https://golang.org, 721ns, 7537 bytes
https://godoc.org, 152ns, 6878 bytes
https://play.golang.org, 205ns, 5767 bytes
http://gopl.io, 326ns, 2856 bytes
--- PASS: Test (1.21s)
PASS
ok  gopl.io/ch9/memo1   1.257s
这个测试是顺序地去做所有的调用的。

由于这种彼此独立的HTTP请求可以很好地并发，我们可以把这个测试改成并发形式。可以使用sync.WaitGroup来等待所有的请求都完成之后再返回。

```
m := memo.New(httpGetBody)
var n sync.WaitGroup
for url := range incomingURLs() {
    n.Add(1)
    go func(url string) {
       start := time.Now()
        value, err := m.Get(url)
        if err != nil {
            log.Print(err)
        }
        fmt.Printf("%s, %s, %d bytes\n",
        url, time.Since(start), len(value.([]byte)))
        n.Done()
    }(url)
}
n.Wait()

```
这次测试跑起来更快了，然而不幸的是貌似这个测试不是每次都能够正常工作。我们注意到有一些意料之外的cache miss(缓存未命中)，
或者命中了缓存但却返回了错误的值，或者甚至会直接崩溃。

