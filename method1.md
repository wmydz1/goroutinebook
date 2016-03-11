> 第一种方法是不要去写变量。考虑一下下面的map，会被“懒”填充，也就是说在每个key被第一次请求到的时候才会去填值
。如果Icon是被顺序调用的话，这个程序会工作很正常，但
如果Icon被并发调用，那么对于这个map来说就会存在数据竞争。

```
var icons = make(map[string]image.Image)
func loadIcon(name string) image.Image

// NOTE: not concurrency-safe!
func Icon(name string) image.Image {
    icon, ok := icons[name]
    if !ok {
        icon = loadIcon(name)
        icons[name] = icon
    }
    return icon
}

```
- 反之，如果我们在创建goroutine之前的初始化阶段，就初始化了map中的所有条目并且再也不去修改它们，
那么任意数量的goroutine并发访问Icon都是安全的，因为每一个goroutine都只是去读取而已。

```
var icons = map[string]image.Image{
    "spades.png":   loadIcon("spades.png"),
    "hearts.png":   loadIcon("hearts.png"),
    "diamonds.png": loadIcon("diamonds.png"),
    "clubs.png":    loadIcon("clubs.png"),
}

// Concurrency-safe.
func Icon(name string) image.Image { 
return icons[name] 
}

```
上面的例子里icons变量在包初始化阶段就已经被赋值了，包的初始化是在程序main函数开始执行之前就完成了的。  
只要初始化完成了，icons就再也不会修改的或者不变量是本来就并发安全的，这种变量不需要进行同步。不过显然我们没法用[…]”

