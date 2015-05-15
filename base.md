在go语言中，有4种引用类型：slice，map，channel，interface。

Slice，map，channel一般都通过make进行初始化：

ci := make(chan int)            // unbuffered channel of integers

cj := make(chan int, 0)         // unbuffered channel of integers

cs := make(chan *os.File, 100)  // buffered channel of pointers to Files

创建channel时可以提供一个可选的整型参数，用于设置该channel的缓冲区大小。该值缺省为0，用来构建默认的“无缓冲channel”，也称为“同步channel”。
Channel作为goroutine间的一种通信机制，与操作系统的其它通信机制类似，一般有两个目的：同步，或者传递消息。