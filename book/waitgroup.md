#WaitGroup
golang中的同步是通过sync.WaitGroup来实现的．WaitGroup的功能：它实现了一个类似队列的结构，可以一直向队列中添加任务，当任务完成后便从队列中删除，如果队列中的任务没有完全完成，可以通过Wait()函数来出发阻塞，防止程序继续进行，直到所有的队列任务都完成为止．

WaitGroup总共有三个方法：Add(delta int),　Done(),　Wait()。

Add:添加或者减少等待goroutine的数量

Done:相当于Add(-1)

Wait:执行阻塞，直到所有的WaitGroup数量变成0

```
package main

import (
	"fmt"
	"sync"
)

var waitgroup sync.WaitGroup

func Afunction(shownum int) {
	fmt.Println(shownum)
	waitgroup.Done() //任务完成，将任务队列中的任务数量-1，其实.Done就是.Add(-1)
}

func main() {
	for i := 0; i < 10; i++ {
		waitgroup.Add(1) //每创建一个goroutine，就把任务队列中任务的数量+1
		go Afunction(i)
	}
	waitgroup.Wait() //.Wait()这里会发生阻塞，直到队列中所有的任务结束就会解除阻塞
}

```

使用场景：  

程序中需要并发，需要创建多个goroutine，并且一定要等这些并发全部完成后才继续接下来的程序执行．WaitGroup的特点是Wait()可以用来阻塞直到队列中的所有任务都完成时才解除阻塞，而不需要sleep一个固定的时间来等待．但是其缺点是无法指定固定的goroutine数目．

