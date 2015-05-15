#Channel机制：

相对sync.WaitGroup而言，golang中利用channel实习同步则简单的多．channel自身可以实现阻塞，其通过<-进行数据传递，channel是golang中一种内置基本类型，对于channel操作只有４种方式：

创建channel(通过make()函数实现，包括无缓存channel和有缓存channel);

向channel中添加数据（channel<-data）;

从channel中读取数据（data<-channel）;

关闭channel(通过close()函数实现，关闭之后无法再向channel中存数据，但是可以继续从channel中读取数据）


channel分为有缓冲channel和无缓冲channel,两种channel的创建方法如下:

var ch = make(chan int) //无缓冲channel,等同于make(chan int ,0)

var ch = make(chan int,10) //有缓冲channel,缓冲大小是５

其中无缓冲channel在读和写是都会阻塞，而有缓冲channel在向channel中存入数据没有达到channel缓存总数时，可以一直向里面存，直到缓存已满才阻塞．由于阻塞的存在，所以使用channel时特别注意使用方法，防止死锁的产生．例子如下：

无缓存channel:

```
package main

import "fmt"

func Afuntion(ch chan int) {
	fmt.Println("finish")
	<-ch
}

func main() {
	ch := make(chan int) //无缓冲的channel
	go Afuntion(ch)
	ch <- 1
	
	// 输出结果：
	// finish
}

```

代码分析：首先创建一个无缓冲channel　ch,　然后执行　go Afuntion(ch),此时执行＜-ch，则Afuntion这个函数便会阻塞，不再继续往下执行，直到主进程中ch<-1向channel　ch 中注入数据才解除Afuntion该协程的阻塞．

```
package main

import "fmt"

func Afuntion(ch chan int) {
	fmt.Println("finish")
	<-ch
}

func main() {
	ch := make(chan int) //无缓冲的channel
	//只是把这两行的代码顺序对调一下
	ch <- 1
	go Afuntion(ch)

	// 输出结果：
	// 死锁，无结果
}

```

代码分析：首先创建一个无缓冲的channel,　然后在主协程里面向channel　ch 中通过ch<-1命令写入数据，则此时主协程阻塞，就无法执行下面的go Afuntions(ch),自然也就无法解除主协程的阻塞状态，则系统死锁

总结：
对于无缓存的channel,放入channel和从channel中向外面取数据这两个操作不能放在同一个协程中，防止死锁的发生；同时应该先利用go 开一个协程对channel进行操作，此时阻塞该go 协程，然后再在主协程中进行channel的相反操作（与go 协程对channel进行相反的操作），实现go 协程解锁．即必须go协程在前，解锁协程在后．

带缓存channel:
对于带缓存channel，只要channel中缓存不满，则可以一直向 channel中存入数据，直到缓存已满；同理只要channel中缓存不为０，便可以一直从channel中向外取数据，直到channel缓存变为０才会阻塞．

由此可见，相对于不带缓存channel，带缓存channel不易造成死锁，可以同时在一个goroutine中放心使用，