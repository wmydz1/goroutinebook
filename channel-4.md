### 8.4.3. 单方向的Channel

随着程序的增长，人们习惯于将大的函数拆分为小的函数。我们前面的例子中使用了三个goroutine，然后用两个channels连链接它们，它们都是main函数的局部变量。将三个goroutine拆分为以下三个函数是自然的想法：

```Go
func counter(out chan int)
func squarer(out, in chan int)
func printer(in chan int)
```

其中squarer计算平方的函数在两个串联Channels的中间，因此拥有两个channels类型的参数，一个用于输入一个用于输出。每个channels都用有相同的类型，但是它们的使用方式想反：一个只用于接收，另一个只用于发送。参数的名字in和out已经明确表示了这个意图，但是并无法保证squarer函数向一个in参数对应的channels发送数据或者从一个out参数对应的channels接收数据。

这种场景是典型的。当一个channel作为一个函数参数是，它一般总是被专门用于只发送或者只接收。


这是改进的版本，这一次参数使用了单方向channel类型：

<u><i>gopl.io/ch8/pipeline3</i></u>
```Go
func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)
	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}
```

调用counter(naturals)将导致将`chan int`类型的naturals隐式地转换为`chan<- int`类型只发送型的channel。
调用printer(squares)也会导致相似的隐式转换，这一次是转换为`<-chan int`类型只接收型的channel。任何双向channel向单向channel变量的赋值操作
都将导致该隐式转换。这里并没有反向转换的语法：也就是不能一个将类似`chan<- int`类型的单向型的channel转换为`chan int`类型的双向型的channel。


