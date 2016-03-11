一个通用的解决方案是将一个函数分离为多个函数，比如我们把Deposit分离成两个：一个不导出的函数deposit，
这个函数假设锁总是会被保持并去做实际的操作，另一个是导出的函数Deposit，这个函数会调用deposit，但在调用前会先去获取锁。
同理我们可以将Withdraw也表示成这种形式：

```
func Withdraw(amount int) bool {
    mu.Lock()
    defer mu.Unlock()
    deposit(-amount)
    if balance < 0 {
        deposit(amount)
        return false // insufficient funds
    }
    return true
}

func Deposit(amount int) {
    mu.Lock()
    defer mu.Unlock()
    deposit(amount)
}

func Balance() int {
    mu.Lock()
    defer mu.Unlock()
    return balance
}

// This function requires that the lock be held.
func deposit(amount int) { balance += amount }

```

当然，这里的存款deposit函数很小实际上取款withdraw函数不需要理会对它的调用，尽管如此，这里的表达还是表明了规则
用限制一个程序中的意外交互的方式，可以使我们获得数据结构的不变性。因为某种原因，封装还帮我们获得了并发的不变性。
当你使用mutex时，确保mutex和其保护的变量没有被导出(在go里也就是小写，且不要被大写字母开头的函数访问啦)，无论这些变量是包级的变量还是
一个struct的字段。

