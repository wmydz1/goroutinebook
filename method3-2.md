 > 一系列的导出函数封装了一个或多个变量，那么访问这些变量唯一的方式就是通过这些函数来做(或者方法，对于一个对象的变量来说)。
每一个函数在一开始就获取互斥锁并在最后释放锁，从而保证共享变量不会被并发访问。这种函数、互斥锁和变量的编排叫作监控monitor
(这种老式单词的monitor是受"monitor goroutine"的术语启发而来的。两种用法都是一个代理人保证变量被顺序访问)。
由于在存款和查询余额函数中的临界区代码这么短--只有一行，没有分支调用--在代码最后去调用Unlock就显得更为直截了当。
在更复杂的临界区的应用中，尤其是必须要尽早处理错误并返回的情况下，就很难去(靠人)判断对Lock和Unlock的调用是在所有路径中都能够严格
配对的了。Go语言里的defer简直就是这种情况下的救星：

我们用defer来调用Unlock，临界区会隐式地延伸到函数作用域的最后，这样我们就从“总要记得在函数返回之后或者发生错误返回时要记得调用一次
Unlock”这种状态中获得了解放。Go会自动帮我们完成这些事情。

```
func Balance() int {
    mu.Lock()
    defer mu.Unlock()
    return balance
}
```
上面的例子里Unlock会在return语句读取完balance的值之后执行，所以Balance函数是并发安全的。这带来的另一点好处是，
我们再也不需要一个本地变量b了。
此外，一个deferred Unlock即使在临界区发生panic时依然会执行，这对于用recover (§5.10)来恢复的程序来说是很重要的。
defer调用只会比显式地调用Unlock成本高那么一点点，不过却在很大程度上保证了代码的整洁性。大多数情况下对于并发程序来说，
代码的整洁性比过度的优化更重要。如果可能的话尽量使用defer来将临界区扩展到函数的结束。

考虑一下下面的Withdraw函数。成功的时候，它会正确地减掉余额并返回true。但如果银行记录资金对交易来说不足，
那么取款就会恢复余额，并返回false。

```
// NOTE: not atomic!
func Withdraw(amount int) bool {
    Deposit(-amount)
    if Balance() < 0 {
        Deposit(amount)
        return false // insufficient funds
    }
    return true
}
```

函数终于给出了正确的结果，但是还有一点讨厌的副作用。当过多的取款操作同时执行时，balance可能会瞬时被减到0以下。
这可能会引起一个并发的取款被不合逻辑地拒绝。所以如果Bob尝试买一辆sports car时，Alice可能就没办法为她的早咖啡付款了。
这里的问题是取款不是一个原子操作：它包含了三个步骤，每一步都需要去获取并释放互斥锁，但任何一次锁都不会锁上整个取款流程。
理想情况下，取款应该只在整个操作中获得一次互斥锁。下面这样的尝试是错误的：

```
// NOTE: incorrect!
func Withdraw(amount int) bool {
    mu.Lock()
    defer mu.Unlock()
    Deposit(-amount)
    if Balance() < 0 {
        Deposit(amount)
        return false // insufficient funds
    }
    return true
}

```
上面这个例子中，Deposit会调用mu.Lock()第二次去获取互斥锁，但因为mutex已经锁上了，而无法被重入
(译注：go里没有重入锁，关于重入锁的概念，请参考java)--也就是说没法对一个已经锁上的mutex来再次上锁--这会导致程序死锁，
没法继续执行下去，Withdraw会永远阻塞下去。

关于Go的互斥量不能重入这一点我们有很充分的理由。互斥量的目的是为了确保共享变量在程序执行时的关键点上能够保证不变性。
不变性的其中之一是“没有goroutine访问共享变量”。但实际上对于mutex保护的变量来说，不变性还包括其它方面。当一个goroutine获得了一个互斥锁时，
它会断定这种不变性能够被保持。其获取并保持锁期间，可能会去更新共享变量，这样不变性只是短暂地被破坏。然而当其释放锁之后，
它必须保证不变性已经恢复原样。尽管一个可以重入的mutex也可以保证没有其它的goroutine在访问共享变量，但这种方式没法保证这些变量额外的不变性。





