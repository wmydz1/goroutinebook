#golang 并发总结：
并发两种方式：sync.WaitGroup，该方法最大优点是Wait()可以阻塞到队列中的所有任务都执行完才解除阻塞，但是它的缺点是不能够指定并发协程数量．
channel优点：能够利用带缓存的channel指定并发协程goroutine，比较灵活．但是它的缺点是如果使用不当容易造成死锁；并且他还需要自己判定并发goroutine是否执行完．

但是相对而言，channel更加灵活，使用更加方便，同时通过超时处理机制可以很好的避免channel造成的程序死锁，因此利用channel实现程序并发，更加方便，更加易用．