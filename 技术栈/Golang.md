# Golang
## 1  Golang面试题
### 1.1  Go语言的优势（为什么选择Go语言）
Go语言语法规则清晰，接近我对自然语言的理解。而且Go语言拥有着几乎跟C语言一样的性能，却比C语言的开发效率更高。且有着强大的标准库，自主的内存管理方式和高效的垃圾回收机制，当然最重要的一点是天生支持并发，只需一个go语句就可以开启一个协程，也就是goroutine。goroutine类似线程却比线程更加轻量，所占内存更小，管理起来也更加方便，对程序的性能有很大的提高。而且我对新的技术也比较感兴趣，学一门很新的语言让我觉得比较有成就感。

### 1.2  内存分配
Go语言是自主管理内存分配的，优势就是不用每次分配内存都进行系统调用，可以提升性能，缺点就是实现起来比较复杂，全靠go语言内置运行时也就是runtime进行管理。
Go 分配内存的过程，主要由四大组件管理，级别从上到下分别是：mheap、mcentral、mcache和mspan。Go程序在启动时，首先会向操作系统申请一大块内存，并交由mheap结构全局管理。mheap 会将这一大块内存，切分成不同规格的小内存块，我们称之为 mspan。然后将某一种规格交给对应的某种mcentral来管理。而mcentral是全局可见的，我们不可能在每个goroutine来申请内存的时候都加锁约束。解决问题的办法是在本地的mcache上分配。就是在一个Go 程序里，工作线程M和处理器P关联之后才能运行goroutine，每个P都会绑定一个叫 mcache 的本地缓存。然后goroutine可以在这个缓存里申请内存，不需要加锁。一般的内存分配情况就是这样的，如果是申请很小的内存比如小于16个字节会使用mcache上的微型分配器分配。如果超过 32KB 的内存申请，会直接从mheap（堆）上分配对应数量的内存页给程序，每页的大小是8kb。

### 1.3  什么是堆内存和栈内存（内存逃逸）
* 堆内存由内存分配器分配由垃圾收集器负责回收。一个程序在运行过程中，只会有一个堆内存（mheap），所以多个goroutine在堆中申请内存需要加锁避免冲突。并且堆内存需要进行垃圾回收，如果有大量的GC操作，将会使程序性能下降的很厉害。所以为了提高程序的性能，应当尽量减少内存在堆上分配。
* 而栈内存在函数结束后会自动回收，不需要GC操作，性能相对堆内存要好得多。所以我们在编写程序时要尽量避免一些栈内变量逃逸到堆上，而且每一次变量逃逸都需要一次额外的内存分配过程，所以避免栈内变量的逃逸对性能的提高是很有效的。比如避免返回指针类型、引用类型的变量、避免在闭包函数中使用外部变量、避免传入不能确定类型的参数如空接口、避免变量的大小大于32k等等。当不可避免地需要使用到这些变量时，我们可以提前在堆上分配，避免了逃逸后进行二次分配，对性能的提高也是有帮助的。
* 编码的时候应该注意：全局变量、指针、引用类型都是在堆上分配的，所以对于小变量，传值效率比传指针效率更高（值在栈中分配）。当我们要在函数中使用到大变量时，先将其在堆上分配然后只需在函数内使用到的时候声明一个变量指向它就可以了，避免了内存逃逸，当已经没有任何变量指向它了，那么它就会被GC回收 。                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       

### 1.4  GMP模型（是抢占式的两级线程模型（m:n））
* G：goroutine是用户态的轻量级的线程，一个goroutine大概需要4k的内存（运行中再动态扩容），一个Go程序创建几百万个goroutine应该是没有问题的。
* M：是工作线程，也就是操作系统内核线程，goruntime默认最多允许创建10000个，超出了就会抛出异常（可以调大）
* P：包含Go代码等的必要资源和可以调度goroutine的处理器，默认等于CPU的核心数，可以通过GOMAXPROCS调小
* Sched：调度器的结构，维护M和G的全局队列，以及调度器的一些状态信息
* 主要的的调度过程就是某个工作线程M先和某个处理器P进行绑定关联形成一个有效的运行环境，然后再选择某个goroutine来运行。
* 详细的调度过程就是每个处理器P都会维护一个本地的队列用来存放goroutine，全局中还有一个全局队列，当我们使用go语句时就会产生一个goroutine优先放到当前P的本地队列，当本地队列满了之后才会放到全局队列。运行的时候P也会是优先在本地队列的队头取出一个goroutine来运行，当本地队列运行完了之后再取全局队列里的goroutine来运行，当全局队列也运行完了操作系统还会在其它P的队列里取出一半（负载均衡的思想）来给这个P运行，也就是work stealing机制。假设工作线程M1正在绑定处理器P1运行goroutine G1的工作，当这个G1发生io调用进行阻塞时，M1会释放绑定的P1（也就是hand off机制）。这个P1会将G1从本地队列中移除，在线程池中取出一个M2来运行G2，如果线程池里面没有工作线程，他就会自己创建一个。

（如何查看运行时的调度信息）

### 1.5  内存对齐

### 1.6  垃圾回收
Go的垃圾回收采用三色标记清除法结合混合写屏障技术。有手动触发和自动触发两种方式；手动触发需要调用runtime包的GC函数，自动触发可以配置多久触发一次或者内存分配达到多少触发一次。一次完整的垃圾回收会分为四个阶段：
（1）标记准备：需STW然后打开写屏障（STW：Stop The World）
（2）标记开始：使用三色标记法标记，和用户程序并发执行
（3）标记终止：STW然后重新扫描触发写屏障的对象，然后在关闭写屏障
（4）清理：把回收的内存归还到堆中，把过多的内存归还给操作系统，清理过程与用户程序并发执行
然后三色标记法可以分为6个步骤：
（1）创建白、灰、黑三个集合
（2）将所有对象放入白色集合中
（3）遍历所有root对象，把遍历到的对象从白色集合放入灰色集合
（4）遍历灰色集合，将灰色对象引用的对象从白色集合放入灰色集合，它自身会被放入黑色集合
（5） 重复步骤4，直到灰色集合中无任何对象（用到2个机制：混合写屏障和辅助GC）
（6） 收集所有白色对象也就是内存垃圾
为什么要STW
STW 的目的是防止 GC 扫描时内存变化引起的混乱，1.18版本采用混合写屏障之后的GC操作几乎不需要STW，只是在标记准备和结束的时候停止大约几十微秒的时间主要是告诉程序开始和结束GC，和一些GC之前的初始化操作。
混合写屏障：是让垃圾清理和用户程序可以并发执行的技术，大大减少了 STW 的时间。开启后指针传递时会把指针标记，即本轮不回收，下次 GC 时再确定。
辅助 GC：为了防止内存分配过快，协助 GC 做一部分工作。

### 1.7  GC调优
(1)控制内存分配的速度，避免变量重复扩容和内存逃逸（如少量使用+连接string，slice提前分配足够的内存来降低扩容带来的拷贝）
(2)限制 Goroutine 的数量
(3)避免map的key过多，导致扫描时间增加
(4)变量复用，减少对象分配，例如使用 sync.Pool 来复用需要频繁创建的临时对象、和使用全局变量等
(5)增大 GOGC 的值，降低 GC 的运行频率
(6)提高赋值器 mutator 的 CPU 利用率（降低GC的CPU利用率）
### 1.8  查看GC信息
(1)GODEBUG='gctrace=1'; go run main.go
(2)go run main.go; go tool trace trace.out
(3)debug.ReadGCStats
(4)runtime.ReadMemStats：

### 1.9  defer（v.推迟）
return有两步：为返回值赋值和返回调用处，defer在其中间执行
defer一般用于成对的操作（如开关文件）、函数收尾和异常捕获
defer函数参数会先在入栈之前求值再入栈，后入的defer先执行
panic之后的defer不会被执行

```go
func F() {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println("捕获异常:", err)
        }
        fmt.Println("b")
    }()
    panic("a")
    defer fmt.Println("不会被执行")
}
```

### 1.10  make和new的区别
（1）首先是应用场景的不一样
make() 函数专门用来为 slice、map、channel 这三种类型来分配内存并进行初始化的；而 new一般用来给这三类之外的其它类型分配内存，实际上new可以给任意类型分配内存，但是用来给这三类分配内存有点不合适。因为new只会返回一个指针指向分配好的内存，而这三类还有其他不同的字段来管理底层数据结构。
以 slice 类型为例，当创建一个 slice 类型的值时，会创建一个轻量级的数据结构，这个数据结构包含三个字段：一个是指向底层数组的指针、另外两个是这个数组的长度和容量。而用new来创建切片时，切片的长度和容量默认被初始化为0值而不能指定其它值，这样在后面赋值的时候会重新分配内存，降低性能。
（2）传入的参数不同
new只能传入一个类型参数；make可以传入一个、两个、或者三个参数
（3）返回值的类型不一样
new() 函数返回一个指向该参数类型的指针。
make() 函数返回值的类型和它接收的第一个参数类型一样。

### 1.11  数组和切片的区别
* 数组

数组是值类型长度是固定的，必须在编译时就确定数组长度的大小，而且数组的长度是数组类型的一部分，在运行的过程中也不能动态扩容。且长度无记录，求数组长度时需要遍历，O(n)。而且数组是值类型，大数组传给函数做参数的时候还会生成副本，导致性能降低。所以除非已经知道了固定长度，否则我们很少使用数组。

* 切片

切片是引用类型长度是可变的，可以在程序的运行过程中动态扩容。切片是一种轻量的数据结构，包含三个字段：指针、长度和容量。指针指到底层数组可以从这个切片访问的第一个元素，长度是当前切片的元素个数，容量是指从当前切片的第一个元素开始到底层数组最后一个元素的元素个数

### 1.12  切片使用的注意事项
（1）切片是一个轻量数据结构包括指向底层数组的指针和切片的长度和容量
（2）切片在使用之前需要用make进行初始化，在make之前是nil切片，给它赋值会报错。所有的nil切片地址都是一样的，输出为0x0。make初始化之后切片内的所有元素都是该类型的零值。当make指定切片的长度为0时，该切片就是一个空切片，空切片指向一个内存地址，但是没有分配内存空间。不管是使用 nil 切片还是空切片，对其调用内置函数 append，len 和 cap 的效果都是一样的。
（3）扩容：扩容后，新切片指向的数组是一个全新的数组。如果新申请容量大于等于原来的两倍，扩容后的容量等于申请的容量（双数）。如果切片的容量小于 1024 个元素，扩容的时候会翻倍增加容量，如果元素个数超过1024个，扩容后容量为原来的1.25倍。
（4）创建切片的时候建议使用字面量创建，而不是在原数组上创建。在数组上创建切片会共享当前数组，也就是说切片在修改的时候，数组也会被修改，导致一些意料之外的情况。
（5）切片不是线程安全的（在底层结构没有加锁的字段），但在并发执行中不会报错，只是值会被随意更改，不能控制。
（6）深拷贝（修改不影响原切片值）：copy和遍历赋值
浅拷贝（修改会影响原切片值）：slice2=slice1

### 1.13  map底层实现原理
首先是使用，map是引用类型，它的零值是nil，在添加元素之前必须要使用make进行初始化。它的键的类型必须是可比较的。
Go Map在源码中是一个叫hmap的结构体。其中的buckets字段是一个切片，里面存放着哈希桶。也是结构体，有四个字段，tophash，keys，elements，overflow。Go的Map是将key和value分开存放的，可以使内存更紧凑。tophash存放哈希值的高字节（key的hash值的高8位），keys存放了键，elems存放数据。每个桶只能存放最多8个键值对，超出时就会使用溢出桶，overflow指向下一个溢出桶的位置。当溢出桶的数量大于原桶或平均每个桶里的键值对>=6.5（测试后选择6.5）时就会导致扩容。扩容分两种情况：如果溢出桶过多但是数据不多只会等量扩容，也就是桶的数量不变，只是进行整理。如果是数据太多了就会翻倍扩容，扩容后普通桶的数量为原来的两倍。

### 1.14  map使用的注意事项以及如何保证线程安全
(1)Go map遍历是无序的，因为每次遍历是从随机的bucket开始的，而且key的位置也在不时改变（不可寻址）。如果想按某顺序，应先对key进行排序，再用key遍历map。
(2)Go map默认是不支持并发安全的。Go官方认为 Go map 更多在典型场景使用，而不应为了小部分情况去加锁，导致更大性能的消耗。同时Go官方提供了并发安全的sync.Map让我们在写并发的时候使用。或者我们也可以在操作map之前先加锁。在“读”多时候适合使用sync.Map，在“写”多的时候适合使用锁和map
(3)map冲突：Go采用链地址法解决冲突，具体就是插入key到map中时，当key定位的桶填满8个元素后，将会创建一个溢出桶，并且将溢出桶插入当前桶所在链表尾部。

### 1.15  channel
通道本质上是一个环形队列，在源码中是一个结构体，其中的buf字段是一个指针，指向在堆内存上的环形队列，里面存放着通道中的数据。lock字段用来进行每次访问的加锁和解锁，保证并发访问的安全性。sendx和recvx表示下次发送数据的位置和接收的位置。sendq和recvq表示在两端阻塞等待的goroutine队列。
### 1.16  channel注意事项
(1)通道属于引用类型，是多个gorotine之间传递数据和同步的重要工具，在使用之前需要用make进行初始化，make之后通道的缓冲容量就被固定了，如果还没有初始化就向通道发送或者接收元素的话会导致当前goroutine永久阻塞。
(2)在make的时候可以用第二个参数指定通道的缓冲容量，如果不指定那么就是一个非缓冲通道。缓冲通道有以下机制可以保证并发安全：1.每次只能有一个goroutine可以向通道中发送数据，有多个goroutine需排队发送，如果通道已经满了，试图发送数据的gorotine就会阻塞，直到可以发送数据。2.每次也只能有一个goroutine可以从通道中取数据（随机一个取），取的时候通道为空的话则当前gorotine会阻塞直到可以成功接收数据。非缓冲通道有一些区别：向非缓冲通道中发送或接收数据，当前goroutine都会阻塞，直到有别的goroutine在通道的另一端和它进行握手交换数据。
(3)channel会发生死锁的场景：
①非缓冲通道只写不读
②非缓冲通道先读后写
③缓存channel写入超过缓冲区数量
④没初始化就读或者写
⑤多个协程互相等待
(4)注意避免goroutine泄露：
多个gorotine排队向无缓冲通道发送数据，后面发送的goroutine因为没有其他gorotine来接收而被永久阻塞，这种情况称为goroutine泄露。和垃圾变量不同，泄露的goroutine并不会被自动回收，所以我们应该保证每一个goroutine都能正常退出。还有一些别的情况比如死锁等导致goroutine资源一直无法释放。

> 排查Goroutine泄露
>
> ①单个函数：调用 runtime.NumGoroutine 方法来打印 执行代码前后Goroutine 的运行数量，进行比较
>
> ②生产/测试环境：使用PProf实时监测Goroutine的数量：
>
> -import _ "net/http/pprof"
>
> -程序中开启HTTP监听服务 http.ListenAndServe("localhost:6060", nil)
>
> -执行命令：go tool pprof -http=:1248 http://127.0.0.1:6060/debug/pprof/goroutine

### 1.17  channel和select


### 1.17  Goroutine的底层原理
goroutine是用户态的轻量级的线程（协程），由Go runtime管理而不是操作系统。在源码中是一个结构体，有一个goid字段用来标识（类似进程的唯一Pid）、一个schedule字段是结构体类型用来在切换的时候保存当前goroutine的上下文、一个stack字段也是结构体，存放栈的上下界内存地址。
### 1.18  Golang控制同步的方式
主要有六种：
（1）channel
channel的存取数据的机制可以用来控制goroutine的同步，且channel是协程安全的不会有数据冲突，比锁好用多了，但是不太适应比较复杂的场景，比如控制多个goroutine的同步问题，而且channel存在死锁问题。（比如存取位置不对，一直阻塞。）
（2）sync.Mutex
（3）sync.WaitGroup
有三个函数Add( ), Done( ), Wait( )	（先统一Add再并发Done最后Wait）
（4）sync.Context主要配合select使用
（5）sync.Pool和sync.Once用的比较少

### 1.19  context
context主要用于控制goroutine的同步以及超时控制、取消方法的调用等。特别是在子goroutine之下还有子goroutine的情况下使用起来就很方便。
context的底层结构就是一个叫Context的接口，其中有四个函数，Deadline，Done，Err和Value。Deadline返回一个time包的Time类型，表示当前Context应该结束的时间；Done主要是告诉context相关的函数要停止工作了；Err表示context被取消的原因；Value是context实现共享数据存储的地方，是协程安全的。
context包还有5个函数，Background和TODO用来返回一个Context；另外三个分别是返回带有取消函数和带有截止日期以及带有超时时间段的上下文和一个取消函数。

> 具体的使用在标准库的context包中


### 1.20  Go Module
Go modules 是 Go 语言的依赖解决方案，是GOPATH的使用模式的进化。主要是解决了GOPATH没有版本控制的概念。现在的gomudle默认都是开启的，只需要在程序运行之前使用go mod init 项目名 生成go.mod文件进行初始化就可以用了。
go mod init 生成 go.mod 文件
go mod download 下载 go.mod 文件中指明的所有依赖
go mod tidy 整理现有的依赖
go mod graph 查看现有的依赖结构
go mod edit 编辑 go.mod 文件
go mod vendor 导出项目所有的依赖到vendor目录
go mod verify 校验一个模块是否被篡改过
go mod why 查看为什么需要依赖某模块

### 1.21  Go Workspace
在之前是先将某个模块上传到github然后再需要用到的时候再import下来。而使用Workspace之后能够在本地项目的go.work文件中，通过设置一系列的模块本地路径，再将路径下的模块组成一个当前的工作区。

### 1.22  值类型和引用类型
值类型通常会在栈上分配：
包括所有的整型int和浮点型float、布尔型、字符串、数组和结构体
值类型存储的就是数据
引用类型通常会在堆上分配：
包括指针、切片、map、通道channel、接口
引用类型存储的是指向底层数据结构的地址，切片、map、通道使用var声明之后必须使用make函数对其进行分配。make之前值都为nil，给它赋值时会导致panic

### 1.23  for和for range的区别

for range 中可以用两个变量来接受值，假设是i和v，它们只会声明一次然后循环接收不同的值，它们相当于只是一个用来接受副本的容器，对它们进行操作并不能改变原始数组的值。如果for range中只使用一个变量来接受值这个for range语句就相当于普通的for循环语句，用得到的i也就是数据组的下标操作数组能改变数组的值。

> for range和for循环中的i和v的地址是不会改变的，&nums[i]是一直在改变的，大小比上次循环的大小加数组元素所占字节数。

### 1.24  内存泄露有哪些场景？

内存泄漏指的是程序在运行过程中分配的内存空间无法被及时释放和回收，导致系统中存在大量不再使用的内存占用，最终导致系统性能下降甚至崩溃的现象。

内存泄漏的场景：

* 通过new()和make()等方式动态分配的内存得不到合理释放
* 多个对象循环引用，却不是程序需要的
* 打开文件、链接等资源后没有及时关闭

### atomic包（原子操作）

### 

### 设计一个携程池

### Golang的Pool的结构体

### Golang的内存模型

### golang里实现栈和队列

### 垃圾回收的根对象有哪些

### 调试（GDB 和 pprof）

### golang用代码实现共享内存



## 2  Golang的安装和环境配置

到官网下载msi文件 --> go version 查看是否安装成果 --> go env -w GOPROXY=https://goproxy.cn,direct 设置Go的代理网站 --> 在环境变量里设置GOPATH路径 --> go env 查看go语言的环境详情-->打开VSCode 新建main.go文件，右下角会提示是否下载依赖包，选择Istall all --> 安装完了之后go语言就可以自动格式化和导包之类的了（安装的文件会放在GOPATH的pkg和bin文件夹下）

## 3  Golang编译器
### 3.1  GoLand

#### 3.1.1  GoLand创建并运行.go文件

文件 --> New --> Project --> 选择目录 --> Create --> Add Configuration（添加配置） --> "+" --> Go Build --> 选择要运行的目录 --> 新建.go文件夹 --> 点绿色三角形运行代码->

#### 3.1.2 GoLand主题设置

File --> Settings --> 插件 --> 搜主题插件 --> 应用

#### 3.1.3 快捷键

| 快捷键    | 功能               |
| --------- | ------------------ |
| Alt+5     | 打开、关闭调试窗口 |
| Alt+Enter | 显示快速修复       |
|           |                    |
|           |                    |
|           |                    |
|           |                    |

### 3.2  VSCode

下载Go扩展 --> 新建文件main.go --> go mod init *Gotest* --> go run main.go

#### 3.2.1  VScode的快捷键

> VSCode可以点击左下角设置符号->键盘快捷方式设置快捷键

| **快捷键**           | **作用**                 | **备注**           |
| -------------------- | ------------------------ | ------------------ |
| Ctrl+A；Ctrl+shift+A | 全选                     | 原生；个人习惯设置 |
| Ctrl + Shift + P；F1 | 显示命令面板             | 原生               |
| Ctrl+R               | 打开最近的文件           | 原生               |
| Shift+Alt + ↓ / ↑    | 向上/向下复制行          | 原生               |
| Ctrl+Enter           | 在下方插入新行           | 原生               |
| Ctrl+Shift+Enter     | 在上方插入新行           | 原生               |
| ctrl+Tab             | 切换文档                 | 原生               |
| ctrl+g               | 切换主侧栏可见性         | 个人习惯设置       |
| ctrl+b               | 转到资源管理器           | 个人习惯设置       |
| ctrl+n               | 新建文件夹（new）        | 个人习惯设置       |
| ctrl+m               | 新建文件                 | 个人习惯设置       |
| Ctrl+Alt+ S          | 全部保存 Save All        | 个人习惯设置       |
| F12                  | 转到定义                 | 原生               |
| F8                   | 转到文件下一个错误或警告 | 原生               |
| Ctrl +  Shift + N    | 打开新的VSCode窗口       | 原生               |
| Ctrl+X               | 剪切行（空选定）         | 原生               |
| Ctrl+C               | 复制行（空选定）         | 原生               |
| Alt+ ↑ /  ↓          | 向上/向下移动行          | 原生               |
| Ctrl+Shift+K         | 删除行                   | 原生               |
| Home                 | 转到行首                 | 原生               |
| End                  | 转到行尾                 | 原生               |
| Ctrl+Home            | 转到文件开头             | 原生               |
| Ctrl+End             | 转到文件末尾             | 原生               |
| Ctrl+↑ /  ↓          | 滚动当前文档             | 原生               |
| Ctrl + F             | 查找                     | 原生               |
| Ctrl + H             | 替换                     | 原生               |
| Tab                  | Emmet 展开缩写           | 原生               |
| Ctrl+`               | 显示集成终端             | 原生               |

#### 3.2.2  VScode的扩展

| 扩展                | 功能           | 备注                   |
| ------------------- | -------------- | ---------------------- |
| 简体中文            | 汉化           |                        |
| go                  |                |                        |
| Live Server         | 同步浏览器开发 | 前端开发               |
| Markdown All in One |                |                        |
| Night Owl           | 主题           | vs里的这个主题也很好看 |



## 4  Golang数据类型
### 4.1  所有基本类型（18个）
`bool`  `byte`  `rune`  `int/uint`  `int8/uint8`  `int16/uint16`  `int32/uint32`  `int64/uint64`  `float32`  `float64`  <u>**`complex64`**</u>  <u>**`complex128`**</u>  `string`

> 这些基本类型都是可比较的

> rune:4字节专门用于存储Unicode字符（如中文）

> complex64：64位复数类型，float32的实部和虚部联合表示

> complex128：128位复数类型，float64的实部和虚部联合表示

> 
>

### 4.2  所有关键字（25个）
* 程序声明（10个）

  `import`  `package` `type` `interface` `const` `var ` `struct` `func` `chan` `map`   

* 流程控制（15个）

  `defer` `for` `range` `continue` `break` `go` `select` `if` `else` `switch` `case` `fallthrough` `default` `goto` `return`   

### 4.3  所有内嵌函数（15个）
`append`  `cap`    `close`   **`complex`**  <u>**`copy `**</u> 

`delete`   <u>**`imag`**</u>   `len`     `make`     `new`  

`panic`   ` print`   `println`  <u>**`real `**</u>     `recover`

> 内嵌函数不需要引入任何包就可以使用它们

* copy()：用于将一个切片或字符串中的元素复制到另一个切片或字符串中，返回复制的元素个数。
* delete()：用于删除一个映射中的键值对，如果映射是nil或没有该键，不做任何操作。
* complex()：用于创建一个复数，返回一个复数类型的值。
* imag()：用于返回一个复数的虚部，返回一个浮点数类型的值。
* real()：用于返回一个复数的实部，返回一个浮点数类型的值。

###1.4  Go常用操作符（21个）

| \|\|   | &&       | ==  | !=  | <   | <=   | >    | \>=    | +   | -   |     |
| ------ | -------- | --- | --- | --- | ---- | ---- | ------ | --- | --- | --- |
|        |          |     |     |     |      |      |        |     |     |     |
| \|     | ^        | *   | /   | %   | <<   | \>>  | &      | &^  | !   | <-  |
| 按位或 | 按位异或 |     |     |     | 左移 | 右移 | 按位与 |     |     |     |

### 4.5 值类型和引用类型

**1.5.1  值类型**

18个基本类型、数组、结构体

**1.5.2  引用类型 **

指针、切片、map、chan、接口

###1.6  可比较类型和不可比较的类型

在 Go 1.18 引入泛型之后，可以用comparable作类型约束，限制参数的类型必须是可比较的，comparable表示的是所有可以比较类型的集合。

```go
func max[T comparable](x, y T) T {
if x > y {
return x
}
return y
}
fmt.Println(max(3, 5)) // 5
fmt.Println(max(time.Now(), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))) // 2024-01-01 00:00:00 +0000 UT
```



**1.6.1  可比较类型**

1. 18个基本类型都是可比较的（包括字符串）
2. 数组（两个数组的元素类型相同且对应位置上的元素值也相等时，它们才被认为是相等的，比较的是数组的值而不是地址。如果数组的元素类型不同或者数组长度不同，那么它们是不可比较的。）
3. 结构体（两个结构体的字段类型和名称都相同且对应字段的值也相等时，它们才被认为是相等的。如果结构体的字段类型或者字段名称不同，那么它们是不可比较的。）
4. 接口（nil接口可比较；动态类型相同且可比较，那么这个接口可比较，否则是不可比较的。）

**1.6.2  不可比较类型**

切片、map、chan  、函数、包含不可比较类型的结构体

###  1.7  Go语言的数据类型

#### byte

byte是uint8的别名

```go

```

#### rune

rune是int32的别名，可以表示任何Unicode编码的字符

```go

```

#### 字符

字符一般用byte，rune和int类型声明，赋值时用单引号包裹。直接输出字符得到的是数字（Unicode编码），用%c输出才是字符。

#### 字符串  string

```go
var name string = "Tom"
```

1. 字符串赋值后是不可变的。

2. 对字符串使用for range得到的v是一个int32类型的字符。直接输出则是数字，用%c输出才是字符。

3. 对字符串使用下标name[0]得到的是一个uint8类型的字符。直接输出则是数字，用%c输出才是字符。（中文字符的输出会有写奇怪的问题）

4. 字符串有两种表现形式：原生字符串用反引号包裹且不会转义；解释型字符串用双引号包裹且可以转义。原生字符串和解释型字符串之间都可以“+”也都可以互相加：a:=”abc”+\`efg`。

5. 对字符串使用`+`拼接生成的是新的字符串，需分配新的地址，性能较低。而使用`strings.Builder`和`bytes.Buffer`能更高效地拼接字符串，避免多次分配内存。

   ```go
   var builder strings.Builder
   builder.Grow (1000) 					// 提前分配1000字节的容量，避免扩容
   builder.WriteString("hello ")
   builder.WriteString("World!")
   var str = builder.String()
   fmt.Println(str)						//Hello World!

   var buffer bytes.Buffer
   buffer.Grow (1000) 						// 提前分配1000字节的容量，避免扩容
   buffer.WriteString("你 ")
   buffer.WriteString("好!")
   var str2 = buffer.String()
   fmt.Println(str2)						//你 好！
   ```

   |              | strings.Builder                                              | bytes.Buffer                                                 |
   | ------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
   | 底层数据结构 | 字节切片[]byte                                               | 字节切片[]byte                                               |
   | 默认大小     | 0字节                                                        | 64字节                                                       |
   | 扩容         | 小于1024字节时扩容会翻倍，大于1024字节时扩容会以125%扩容。   | 小于1024字节时扩容会翻倍，大于1024字节时扩容会以150%扩容。   |
   | 接受类型     | 只能接受字节切片或字符串类型的参数                           | 可以接受任何实现了io.Writer接口的类型的参数                  |
   | string的生成 | 它的String方法将底层的字节切片直接转换为string返回，不会产生额外的内存分配和拷贝 | 它的String方法会重新申请一块空间，存放生成的string变量，然后返回，会产生额外的内存分配和拷贝 |
   | 性能         | 较高                                                         | 较低                                                         |
   | 使用场景     | 拼接字符串                                                   | 处理任何字节流（当然也包括拼接字符串）                       |
   | 提供的方法   |                                                              |                                                              |

   **strings包里的Reader**

   > Reader类型通过从一个字符串读取数据，实现了io.Reader、io.Seeker、io.ReaderAt、io.WriterTo、io.ByteScanner、io.RuneScanner接口

   ```go
   type Reader struct{
       //...
   }
   func NewReader(s string) *Reader//使用此方法可以生成一个指向新生成的Reader的地址

   func (r *Reader) Len() int
   func (r *Reader) Read(b []byte) (n int, err error)
   func (r *Reader) ReadByte() (b byte, err error)
   func (r *Reader) UnreadByte() error
   func (r *Reader) ReadRune() (ch rune, size int, err error)
   func (r *Reader) UnreadRune() error
   func (r *Reader) Seek(offset int64, whence int) (int64, error)
   func (r *Reader) ReadAt(b []byte, off int64) (n int, err error)
   func (r *Reader) WriteTo(w io.Writer) (n int64, err error)
   ```

   **strings包里的Replacer**

   > Replacer类型可以用来进行一系列字符串的替换。

   ```go
   type Replacer struct {
       // 内含隐藏或非导出字段
   }
   func NewReplacer(oldnew ...string) *Replacer

   //Replace返回s的所有替换进行完后的拷贝
   func (r *Replacer) Replace(s string) string 
   //WriteString向w中写入s的所有替换进行完后的拷贝
   func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)
   ```

   **strings包里操作字符串的方法**

   > str1 = "Hello"; str2 = " World!"

   | 方法                                                       | 作用                                                         |
   | ---------------------------------------------------------- | ------------------------------------------------------------ |
   | func **Join**(a []string, sep string) string               | 将一系列字符串**连接为一个字符串**，之间用sep来分隔          |
   | func **EqualFold**(s, t string) bool                       | 判断两个utf-8编码字符串**是否相同**。（将unicode大写、小写、标题三种格式字符视为相同） |
   | func **HasPrefix**(s, prefix string) bool                  | 判断s是否有**前缀**字符串 prefix                             |
   | func **HasSuffix**(s, suffix string) bool                  | 判断s是否有**后缀**字符串suffix                              |
   | func **Contains**(s, substr string) bool                   | 判断字符串s是否**包含子串**substr                            |
   | func **ContainsRune**(s string, r rune) bool               | 判断字符串s是否**包含utf-8码值r**                            |
   | func **ContainsAny**(s, chars string) bool                 | 判断字符串s是否**包含字符串chars中的任一字符**               |
   | func **Count**(s, sep string) int                          | 返回字符串s中有几个**不重复的sep子串**                       |
   | func **Index**(s, sep string) int                          | 子串sep在字符串s中**第一次出现的位置**，不存在则返回-1       |
   | func **LastIndex**(s, sep string) int                      | 子串sep在字符串s中**最后一次出现的位置**，不存在则返回-1     |
   | func **IndexByte**(s string, c byte) int                   | **字符c在s中第一次出现的位置**，不存在则返回-1               |
   | func **IndexRune**(s string, r rune) int                   | **unicode码值r在s中第一次出现的位置**，不存在则返回-1        |
   | func **IndexAny**(s, chars string) int                     | **字符串chars中的任一utf-8码值在s中第一次出现的位置**，如果不存在或者chars为空字符串则返回-1 |
   | func **LastIndexAny**(s, chars string) int                 | **字符串chars中的任一utf-8码值在s中最后一次出现的位置**，如不存在或者chars为空字符串则返回-1 |
   | func **IndexFunc**(s string, f func(rune) bool) int        | **s中第一个满足函数f的位置i**（该处的utf-8码值r满足f(r)==true），不存在则返回-1 |
   | func **LastIndexFunc**(s string, f func(rune) bool) int    | **s中最后一个满足函数f的unicode码值的位置i**，不存在则返回-1 |
   | func **Title**(s string) string                            | 单词**首字母**改为Title格式                                  |
   | func **ToTitle**(s string) string                          | **所有字母**改为Title格式                                    |
   | func **ToUpper**(s string) string                          | 返回将所有字母都转为对应的**大写**版本的拷贝                 |
   | func **ToLower**(s string) string                          | 返回将所有字母都转为对应的**小写**版本的拷贝                 |
   | func **Repeat**(s string, count int) string                | 返回count个s**串联**的字符串                                 |
   | func **Replace**(s, old, new string, n int) string         | 返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串 |
   | func **Map**(mapping func(rune) rune, s string) string     | 将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝。如果mapping返回一个负值，将会丢弃该码值而不会被替换。（返回值中对应位置将没有码值） |
   | func **Trim**(s string, cutset string) string              | 返回将s前后端所有cutset包含的utf-8码值都**去掉**的字符串     |
   | func **TrimSpace**(s string) string                        | 返回将s前后端**所有空格都去掉**的字符串                      |
   | func **TrimFunc**(s string, f func(rune) bool) string      | 返回将s前后端所有**满足f的unicode码值都去掉**的字符串        |
   | func **TrimLeft**(s string, cutset string) string          | 返回将s前端所有cutset包含的utf-8码值都去掉的字符串           |
   | func **TrimRight**(s string, cutset string) string         | 返回将s后端所有cutset包含的utf-8码值都去掉的字符串           |
   | func **TrimLeftFunc**(s string, f func(rune) bool) string  | 返回将s前端所有满足f的unicode码值都去掉的字符串。            |
   | func **TrimRightFunc**(s string, f func(rune) bool) string | 返回将s后端所有满足f的unicode码值都去掉的字符串              |
   | func **TrimPrefix**(s, prefix string) string               | 返回**去除s可能的前缀prefix**的字符串                        |
   | func **TrimSuffix**(s, suffix string) string               | 返回**去除s可能的后缀suffix**的字符串                        |
   | func **Fields**(s string) []string                         | 返回将字符串**按照空格分割**的多个字符串。如果字符串全部是空白或者是空字符串的话，会返回空切片 |
   | func **Split(**s, sep string) []string                     | 用去掉s中出现的sep的方式进行**分割**，会分割到结尾，并返回生成的所有片段组成的切片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如果sep为空字符，Split会将s切分成每一个unicode码值一个字符串。 |
   | func **SplitN**(s, sep string, n int) []string             | 类似函数Split， 参数n决定返回的切片的数目：n > 0 : 返回的切片最多n个子字符串；最后一个子字符串包含未进行切割的部分。 n == 0: 返回nil n < 0 : 返回所有的子字符串组成的切片 |
   | func **SplitAfter**(s, sep string) []string                | 用从s中出现的sep**后面切断的方式进行分割**，会分割到结尾，并返回生成的所有片段组成的切片（每一个sep都会进行一次切割，即使两个sep相邻，也会进行两次切割）。如果sep为空字符，Split会将s切分成每一个unicode码值一个字符串。 |
   | func **SplitAfterN**(s, sep string, n int) []string        | 类似函数Split， 参数n决定返回的切片的数目：n > 0 : 返回的切片最多n个子字符串；最后一个子字符串包含未进行切割的部分。 n == 0: 返回nil n < 0 : 返回所有的子字符串组成的切片 |

   **strconv 包里转换字符串的方法**

   | 方法                                                         | 作用                                     |
   | ------------------------------------------------------------ | ---------------------------------------- |
   | func **ParseInt**(s string, base int, bitSize int) (i int64, err error) | 将字符串转换为整数，可以指定转换的进制   |
   | func **FormatInt**(i int64, base int) string                 | 返回i的base进制的字符串表示              |
   | func **AppendInt**(dst []byte, i int64, base int) []byte     | 等价于append(dst, FormatInt(I, base)...) |

   ​



#### 数组

```go
var family [2]string = [2]string{"Tom","Jerry"}
```

1. 数组在内存中的布局
2. 数组的for-range
3. 数组长度不可变，创建后没有赋值会有默认值。
4. ​

#### 结构体

```go
type Student struct{
    name string
    age int
}
var tom Student= Student{
	name:"tom",
	age:18}
var jerry = struct{}{}					//此jerry输出的值为{}，这是结构体的零值
```

结构体中可以有嵌入字段
结构体的所有字段在内存中是连续的

#### 函数

```go
func FuncName(x int, y int)int
```

* 调用一个函数时，程序会分配给这个函数一个新的栈空间。当函数调用完毕后，程序会自动回收这个栈空间。

* Go的函数支持多个返回值。当不需要接收某个返回值时可以用占位符“_”忽略掉。

* Go的函数不支持重载，多个同名函数会报错。Go的函数也是一种数据类型，可以赋值给一个变量，通过该变量可以对函数调用。函数也可以作为另一个函数的形参。

* **参数**

  * 参数必须有名称，结果的名称则可有可无，要么全有，要么全无

  * 参数是值传递，传入时会在函数内生成一个副本。传入引用类型时的注意事项：

    ```go
    //首先要知道：引用类型是一个轻量的数据类型。切片则是一个结构体，包含容量大小和指针
    func main() {
    	var nums []int = make([]int, 2, 2)
    	nums[0] = 2
    	nums[1] = 4
    	fmt.Printf("%p,%p\n", &nums, nums)
    	test(nums)
    	fmt.Printf("%v,%p,%p\n", nums, &nums, nums)
    }

    func test(nums []int) {
    	fmt.Printf("%p,%p\n", &nums, nums)	//当传入切片时，会生成一个切片结构体的副本，这两个结构体的指向底层数组的指针是一样的。&nums和函数外的不一样，而nums和函数外的一样。
    	nums[0] = 8							//引用类型在函数内修改会影响函数外的值
    	nums = append(nums, 6)	//而当切片扩容时，函数内的切片结构体的指针指向了新的底层数组
    	nums[1] = 3				//所以再次修改切片时，不会影响到函数外的切片
    	fmt.Printf("%v,%p,%p\n", nums, &nums, nums)
    }
    //当传入的是切片地地址时则不需要考虑这些问题
    ```

  * Go语言支持可变参数。

    * 可变参数必须是函数签名中的最后一个参数。
    * 可变参数可以直接传入多个值，也可以传入切片的展开形式slice...

    ```go
    //...int表示可以输入0个或多个参数，ags在函数内为一个切片，里面包含传入的参数
    func sum(ags ...int)int{}					//支持0到多个参数
    func sum(i int,ags ...int)int{}				//支持1到多个参数

    sum(2,4,6)				//正确
    sum([]int{2,4,6}...)	//正确
    sum([]int{2,4,6})		//错误
    ```

  * 1

* **匿名函数**

  > 匿名函数就是没有名字的函数，一般是声明时直接调用，且只能使用一次。也可以将匿名函数赋值给某个变量，那么就可以多次调用了。

  ```go
  func (x,y int){
  	fmt.Println(x+y)
  }(1,2)//带括号，直接执行

  fun3:=func(){
  	fmt.Println(x+y)
  }//还未执行
  fun3()
  ```

  ​


* **回调函数**

  > 将函数func2 作为函数func1的一个参数，那么func1叫做高阶函数，func2叫做回调函数

* **闭包函数**

  > 闭包函数是指在一个函数内部定义的函数，可以访问并操作其外部函数的局部变量

  * golang在编译时会把闭包函数和它引用的对象一起打包成闭包对象，存放在堆内存中，由golang的垃圾回收器管理，不会随外部函数的生命周期结束而被销毁，同时也能有效的避免内存泄露。


  * 闭包就是一个函数和与其相关的引用环境组合的一个整体(实体)。
  * 当匿名函数引用了外部的变量，那么这个匿名函数就变成了闭包函数

* **特殊函数main和init**

  * main函数和init函数在定义时不能有任何参数和返回值，且只能被程序自动调用，不可以被引用。
  * init 可以应用于任意包中，且可以重复定义多个。 main 函数只能用于 main 包中，且只能定义一个。
  * init 可以应用于任意包中，且可以重复定义多个。 main 函数只能用于 main 包中，且只能定义一个。
  * 即使一个包被其它多个包导入，它的init函数也只会执行一次。
  * 执行顺序：import->const->var->init()->main()



​

​

#### 方法

* 方法可以定义在任何类型上，但是不能定义在接口类型或者内置的基本类型（int float等）上。
* 同一个包内不能有相同名字的方法，即使它们的接收者的类型不同。
* ​

值类型和指针类型都可以调用对应的方法，是否能对值做出改变取决于方法定义的时候使用的是值还是指针类型来定义的

```go
type Tom struct {
    name string
}

//值类型的方法
func (t Tom) ChangeName() {
    t.name = "Jerry"
}
t1 := Tom{width: "Tom"}
t1.ChangeName()										//此时t1还是{"Tom"}
    

//指针类型的方法
func (t *Tom) ChangeName2() {
    t.name = "Jerry"					
}
t2 := &Tom{name: "Tom"} 
t2.ChangeName2()										//此时t2变成{"Jerry"}
```



#### 切片 slice

```go
var name []string = make([]string,2,4)
var name []string = []string{"Tom","Jerry"}	//此时长度和容量默认为2


```

* 初始化切片时，make函数需要有2个或3个参数，当只有两个参数时，容量默认和长度相等。

（1）切片是一个轻量数据结构包括指向底层数组的指针和切片的长度和容量
（2）切片在使用之前需要用make进行初始化，在make之前是nil切片，给它赋值会报错。所有的nil切片地址都是一样的，输出为0x0。make初始化之后切片内的所有元素都是该类型的零值。当make指定切片的长度为0时，该切片就是一个空切片，空切片指向一个内存地址，但是没有分配内存空间。不管是使用 nil 切片还是空切片，对其调用内置函数 append，len 和 cap 的效果都是一样的。
（3）扩容：扩容后，新切片指向的数组是一个全新的数组。如果新申请容量大于等于原来的两倍，扩容后的容量等于申请的容量（双数）。如果切片的容量小于 1024 个元素，扩容的时候会翻倍增加容量，如果元素个数超过1024个，扩容后容量为原来的1.25倍。
（4）创建切片的时候建议使用字面量创建，而不是在原数组上创建。在数组上创建切片会共享当前数组，也就是说切片在修改的时候，数组也会被修改，导致一些意料之外的情况。
（5）切片不是线程安全的（在底层结构没有加锁的字段），但在并发执行中不会报错，只是值会被随意更改，不能控制。
（6）深拷贝（修改不影响原切片值）：copy和遍历赋值
浅拷贝（修改会影响原切片值）：slice2=slice1

#### 映射 map

map可以用len

map可以用for range，但是遍历顺序不一样，map可以边遍历边删除

在map中访问一个不存在的键时，它会返回该值类型的零值。

map的key为什么是无序的？





#### 通道 channel

```go
var ch chan int
ch = make(chan int)			//无缓冲channel
ch = make(chan int, 2)		//有缓冲channel，此时len（ch）为0，cap（ch）为2
```

* make函数在初始化channel时只能接受一个或两个参数，channel只能指定容量不能指定长度

#### 接口 interface

```go
type Student interface{
	Study()
}
```

* 如果结构体Tom实现了接口Student，那么可以推断出*Tom也实现了Student；然而反过来推断并不成立。

* 实现了接口的变量都可以赋值给接口变量，但是接口对象不能调用实现对象的属性。

* 接口应该尽量小而专注于一个特定的任务，提高代码可读性降低耦合度。例如io.Reader和io.Writer，它们只定义了一个方法，但是可以适用于很多类型。

* 类型断言会破坏接口的抽象性，增加代码复杂度，应尽量避免，然后使用更小更具体的接口代替。

* 应该使用接口组合来代替接口的嵌套，可以让代码更加灵活，降低耦合度。（例如io.ReadWriter就是通过组合io.Reader和io.Writer来定义的。）

  > golang的接口没有继承的定义，但是可以嵌入任意其它接口类型来实现继承。



**空接口**

* 空接口空接口可以存储任意类型的值，但是在使用这些值之前，需要进行类型断言。

  ```go
   s, ok := i.(string)					//如果ok是true，则s是string类型的值
  ```

*  使用空接口可以提高灵活性。

* 使用空接口会降低程序的性能。因为空接口中存的是两个指针，一个指向值，一个指向值的类型，取值的时候根据类型来调用相应的方法，比直接读取要慢。

```go
package main

import "fmt"

func main() {
	any := make([]interface{}, 5)
	any[0] = 11
	any[1] = "hello world"
	any[2] = []int{11, 22, 33, 44}
	for _, value := range any {
		fmt.Println(value)
	}
}


func main() {
	testSlice := []int{11,22,33,44}

	// 成功拷贝
	var newSlice []int
	newSlice = testSlice
	fmt.Println(newSlice)

	// 拷贝失败
	var any []interface{}
	any = testSlice
	fmt.Println(any)
}
```



##### 类型断言

类型断言v1.（I1）：判断一个接口的实际类型是否为某个类型
注意事项：
（1）v1必须是一个接口值（不是得转换）
（2）结果为否会产生一个恐慌，解决办法：i1,ok:=interface{}(v1).(I1)

#### any

any是interface{}的别名

### 5  Go语言的流程语句

#### 5.1  if else

```go
if age < 18 {
fmt.Println("You are  a child.")
} else if a {
fmt.Println("You are not an adult.")
}
```

#### 5.2  switch

Go语言的case语句执行后会自动返回，不需要break语句。

```go
switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	default:
		fmt.Println("Invalid day")
	}
```



（1）case后的语句都会被求值
（2）switch语句也可以包含一条子语句来初始化局部变量
（3）case后可以跟多个表达式，逗号隔开，满足一个即可
（4）break可以退出当前switch，也可以break+标签
类型sswitch语句：
（5）switch v.(type){}，case后只跟类型
（6）类型switch不允许有fallthrough

#### for

#### for range语句：

for range 中可以用两个变量来接收值，i是下标v是值的副本。它们只会声明一次然后循环接收不同的值，它们相当于只是一个用来接受副本的容器，对它们进行操作并不能改变原始数组的值。

* 等号右边只有一个接收，得到的是下标。此时和普通的for循环一样。
* for range和for循环中的i和v的地址是不会改变的，&nums[i]是一直在改变的，大小比上次循环的大小加数组元素所占字节数。
* 可以用for range迭代的数据结构：字符串、数组、切片、map、通道
* 迭代无元素的数组，nil切片、nil map 一开始就直接结束
* 迭代nil通道，当前goroutine会永远阻塞在此

#### goto语句

#### defer语句

defer后面必须是函数或者方法的调用。所以defer的数据结构跟一般函数类似。\_defer结构体中包含了被调用的函数的函数指针、参数、返回值、堆栈指针、程序计数器、关联的panic等信息，以及一个link字段，用于链接下一个_defer结构体，形成一个单链表。新声明的defer会插入链表的表头。

```go
//defer在源码中的结构：
type _defer struct {
siz     int32 // includes both arguments and results
started bool
sp      uintptr // sp at time of defer
pc      uintptr
fn      *funcval
_panic  *_panic // panic that is running defer
link    *_defer
}
```



* return有两步：为返回值赋值和返回调用处，defer在其中间执行。

* defer一般用于成对的操作（如开关文件）、函数收尾和异常捕获。

  > 在申请资源后应立即使用defer关闭资源。

* 要在defer函数中使用外部变量应通过参数传入。defer函数参数会先在入栈之前求值再入栈，后入的defer先执行。

  > * defer 在入栈时会生成值的副本。需要注意的是如果defer内生成的是变量地址的副本，那么defer外对变量的值的修改也会影响defer内的变量的值。
  > * 对defer外的变量i求地址和对defer内的变量i求的地址值是一样的。
  > * defer被设置为后进先出的原因：资源的申请也是有顺序的，比如先申请A资源再申请B资源，此时也需要关闭B资源再关闭A资源。

* defer可以改变返回值

```go
//当函数的返回值有名称时，defer可以改变返回值的值
func deferFuncReturn() (i int) {		//调用这个函数得到的是2，defer改变了i的值
    i = 1
    defer func() {
       i++
    }()
    return i
}
func deferFuncReturn() int {			//调用这个函数得到的是，defer没有改变了i的值
	i := 1
	defer func() {
		i++
	}()
	return i
}
```

* 如果`defer`后面跟的是多级函数的调用，只有最后一个函数会被延迟执行。例如 `defer NewFoo(p).Bar(p)`中的NewFoo(p)会直接执行，只有.Bar()函数被defer了。

#### panic

`panic()`是Go语言的内嵌函数，用于抛出一个异常。当某个函数调用了panic，这个函数的执行就会停止，在panic之前定义的defer操作都会被执行，然后函数返回。

#### recover

`recover()`是一个让函数从panic状态恢复的内嵌函数，需要在defer 后面的函数里使用。函数正常时返回nil，函数异常时可以捕获异常，然后恢复函数的正常运行。

```go
func F() {
    defer func() {							//recover一般放在函数体的开始处
        if err := recover(); err != nil {
            fmt.Println("捕获异常:", err)
        }
        fmt.Println("b")
    }()
    panic("a")
    defer fmt.Println("不会被执行")				//panic之后的defer不会被执行
}
```

## 6  Go语言特性和机制

### 6.1  反射

> 反射是指程序在运行时动态地获取和操作变量和类型的能力。主要通过refleg包实现。

```go
//常用的reflect包函数如下：
//reflect.TypeOf：获取变量的类型信息
//reflect.ValueOf：获取变量的值
//Value.Interface：将Value转换为interface{}类型
//Value.CanSet：判断Value是否可被修改
//Value.Set：设置变量的值
//Value.MethodByName：根据方法名调用变量的方法

package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{Name: "Alice", Age: 25}

	// 获取变量的类型信息
	pType := reflect.TypeOf(p)
	fmt.Println("Type:", pType.Name())

	// 获取变量的值
	pValue := reflect.ValueOf(p)

	// 遍历结构体字段并获取字段的值
	for i := 0; i < pType.NumField(); i++ {
		field := pType.Field(i)
		value := pValue.Field(i).Interface()
		fmt.Printf("%s: %v\n", field.Name, value)
	}
}
//输出结果：
//Type: Person
//Name: Alice
//Age: 25
```



## 7  Go语言的IO

### 7.1 fmt包

***

| Input     | Output | 备注 | 使用                          |
| --------- | ------ | :--: | ----------------------------- |
| Scan()    |        |      | fmt.Scan(&name, &age)         |
| Scanf()   |        |      | fmt.Scanf("%s %d",&name,&age) |
| Scanln()  |        |      | fmt.Scanln(&name, &age)       |
| Sscan()   |        |      |                               |
| Sscanf()  |        |      |                               |
| Sscanln   |        |      |                               |
| Fscan()   |        |      |                               |
| Fscanf()  |        |      |                               |
| Fscanln() |        |      |                               |

### 7.2  os包

* os包中的权限：文件类型|自己|组|其它<--> `-`是文件；d是目录|0是无权限；1是可执行；2是可写；4是可读

```go
//FileInfo接口中定义了File信息相关的方法。
type FileInfo interface {
    Name() string       // base name of the file	 文件名.扩展名 如a.txt
    Size() int64        // 文件大小，字节数 		  12540
    Mode() FileMode     // 文件权限 	  			 -rw-rw-rw-
    ModTime() time.Time // 修改时间 				 2018-04-13 16:30:53 +0800 CST
    IsDir() bool        // 是否文件夹
    Sys() interface{}   // 基础数据源接口(can return nil)
}

//实现了FileInfo接口就可以使用它的方法
func Stat(name string) (fi FileInfo, err error)					//尝试跳转链接
func Lstat(name string) (fi FileInfo, err error)				//不会跳转链接

//File代表一个打开的文件对象
type File struct {//...}
    
//创建文件，返回的file默认是0666（任何人可读写不可执行）
func Create(name string) (file *File, err error)
 
//创建文件，可以指定文件描述符
func NewFile(fd uintptr, name string) *File
    
//返回两个文件对象，从r读取然后写入w
func Pipe() (r *File, w *File, err error)
  
//打开一个文件用于读取，对应O_RDONLY
func Open(name string) (file *File, err error)

//打开一个文件，可以指定打开方式例如0755和O_RDONLY等
func OpenFile(name string, flag int, perm FileMode) (file *File, err error)
 

func (f *File) Name() string
func (f *File) Stat() (fi FileInfo, err error)
func (f *File) Fd() uintptr
func (f *File) Chdir() error
func (f *File) Chmod(mode FileMode) error
func (f *File) Chown(uid, gid int) error
func (f *File) Readdir(n int) (fi []FileInfo, err error)
func (f *File) Readdirnames(n int) (names []string, err error)
func (f *File) Truncate(size int64) error
func (f *File) Read(b []byte) (n int, err error)
func (f *File) ReadAt(b []byte, off int64) (n int, err error)
func (f *File) Write(b []byte) (n int, err error)
func (f *File) WriteString(s string) (ret int, err error)
func (f *File) WriteAt(b []byte, off int64) (n int, err error)
func (f *File) Seek(offset int64, whence int) (ret int64, err error)
func (f *File) Sync() (err error)
func (f *File) Close() error
//创建目录
directoryPath := "./example_directory"
func createDirectory(directoryPath string) {
	err := os.Mkdir(directoryPath, 0755)
	if err == nil {
		fmt.Printf("创建目录成功：%s\n", directoryPath)
	} else if os.IsExist(err) {
		fmt.Printf("目录已存在：%s\n", directoryPath)
	} else {
		fmt.Printf("创建目录失败：%s\n", err.Error())
	}
}

//创建文件
filePath := "./example_directory/example.txt"
func createFile(filePath string) {
	data := []byte("这是一个示例文件\n")
	err := ioutil.WriteFile(filePath, data, 0644)
	if err == nil {
		fmt.Printf("创建文件成功：%s\n", filePath)
	} else {
		fmt.Printf("创建文件失败：%s\n", err.Error())
	}
}

//读取文件
func readFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err == nil {
		fmt.Printf("读取文件内容成功：\n%s\n", content)
	} else if os.IsNotExist(err) {
		fmt.Printf("文件不存在：%s\n", filePath)
	} else {
		fmt.Printf("读取文件内容失败：%s\n", err.Error())
	}
}

//删除文件
func deleteFile(filePath string) {
	err := os.Remove(filePath)
	if err == nil {
		fmt.Printf("删除文件成功：%s\n", filePath)
	} else if os.IsNotExist(err) {
		fmt.Printf("文件不存在：%s\n", filePath)
	} else {
		fmt.Printf("删除文件失败：%s\n", err.Error())
	}
}

//删除目录
func deleteDirectory(directoryPath string) {
	err := os.Remove(directoryPath)
	if err == nil {
		fmt.Printf("删除目录成功：%s\n", directoryPath)
	} else if os.IsNotExist(err) {
		fmt.Printf("目录不存在：%s\n", directoryPath)
	} else {
		fmt.Printf("删除目录失败：%s\n", err.Error())
	}
}
```

### 7.2  io包

* io包的Reader：

  ```go
  type Reader interface {
          Read(p []byte) (n int, err error)
  }
  ```

  > Read 将 len(p) 个字节读取到 p 中。它返回读取的字节数 n（0 <= n <= len(p)）以及任何遇到的错误。即使 Read 返回的 n < len(p)，它也会在调用过程中使用 p的全部作为暂存空间。若一些数据可用但不到 len(p) 个字节，Read 会照例返回可用的东西，而不是等待更多。
  >
  > 当 Read 在成功读取 n > 0 个字节后遇到一个错误或 EOF 情况，它就会返回读取的字节数。它会从相同的调用中返回（非nil的）错误或从随后的调用中返回错误（和 n == 0）。这种一般情况的一个例子就是 Reader 在输入流结束时会返回一个非零的字节数，可能的返回不是 err == EOF 就是 err == nil。无论如何，下一个 Read 都应当返回 0, EOF。
  >
  > 调用者应当总在考虑到错误 err 前处理 n > 0 的字节。这样做可以在读取一些字节，以及允许的 EOF 行为后正确地处理I/O错误。
  >
  > Read 的实现会阻止返回零字节的计数和一个 nil 错误，调用者应将这种情况视作空操作。

* io包的Writer：

  ```g
  type Writer interface {
      Write(p []byte) (n int, err error)
  }
  ```

  > Writer接口用于包装基本的写入方法。
  >
  > Write方法len(p) 字节数据从p写入底层的数据流。它会返回写入的字节数(0 <= n <= len(p))和遇到的任何导致写入提取结束的错误。Write必须返回非nil的错误，如果它返回的 n < len(p)。Write不能修改切片p中的数据，即使临时修改也不行。







## 8 Go并发

保证并发安全性的方法：锁、原子操作、通道、同步机制（WaitGroup、Cond、Context等）

### 8.1  通道 channel



```go
type Product struct {
	id     int
	product_id int
}

//通道也拥有长度（len）和容量（cap），但是初始化的时候只能初始化容量（下面的make中加入长度编译就会报错）
var ch chan Product = make(chan Product, 1)			

```

- chan后跟类型才是通道真正的类型，例如`chan int`。

- 通道有三种操作：发送、读取和关闭。

- 通道关闭后不能发送但是可以读取，读到真值时第二个参数为true，读到零值时为false。

- 通道在初始化之前是nil通道，往里面发送或接受值都会使当前goroutine永久阻塞。

- 可以在函数的参数中传入单向通道。

  > 在声明时也可以声明单向通道，但是一般会声明双向通道，然后在传入某个函数时，在参数限制它在这个函数作用域内仅可以作为某个方向的通道使用。发送通道只能向它发送数据，接受通道只能用来接受且不能被关闭。

- 遍历通道

  ```go
  //一般都会用一个无限循环来读取通道，避免goroutine泄露
  for {
  		i, ok := <-ch
  		if !ok {
  			break
  		}
  		fmt.Println(i)
  	}

  //使用for和select组合遍历通道会更加灵活
  for {
  		select {
  		case i := <-ch1:
  			fmt.Println("From ch1:", i)
  		case j := <-ch2:
  			fmt.Println("From ch2:", j)
  		default:							//当此刻i，j中都没有数据时会走default分支
  			fmt.Println("No data")
  			time.Sleep(time.Second / 2)		//睡眠半秒
  		}
  	}
  ```


  ```

- 通道也可以用来传递通道

  ```go
  var chch1 chan chan int
  ```

### 8.2  channel与select

select的case后面跟着对通道的收发操作，select会随机选择一个不阻塞的case运行，如果没有则执行default。如果没有default则会阻塞当前goroutine直到有某个case可以运行。

```go
	loop:											//标签，代表跟在它后边的循环结构
	for {											//select一般配合无限循环使用
		select {
		case x := <-ch1:
			fmt.Println("Received from ch1:", x)
		case ch2 <- y:
			fmt.Println("Sent to ch2:", y)
		case <-quit:
			fmt.Println("Quit")
			break loop								//结束这个loop代表的循环体
		default:
			fmt.Println("No communication")
		}
	}
```

### 8.3为channel设置超时时间

### 原子操作

### 锁

####互斥锁

互斥锁Mutex，可以创建为其他结构体的字段，零值为解锁状态。

```go
type Mutex struct {/* 包含隐藏或非导出字段*/ }		//在源码sync.go中的定义

```

#### 读写锁

```go
type RWMutex struct {
    // 包含隐藏或非导出字段
}
```

### Context

Context主要在异步场景中用于实现并发协调以及对 goroutine 的生命周期控制。

 ```go
type Context interface {							   //Context在源码context.go里的声明
	Deadline() (deadline time.Time, ok bool)			//返回 context 的过期时间；
	Done() <-chan struct{}							   //返回 context 中的 channel；
	Err() error										  //返回错误
	Value(key any) any								   //返回 context 中的对应 key 的值.
}
 ```



### 9  Go面向对象编程

1、Golang是基于结构体实现面向对象编程的，通过接口进行关联。

2、面向对象三大特性：封装、继承、多态

(1)  封装（结构体封装事物属性）

(2)  继承（内嵌共有字段的结构体）

(3)  多态（通过接口实现）

3、工厂模式

4、MVC分别表示什么，有什么作用（结合项目经历来讲）

Model：模型，用来封装对象，查询数据库等完成具体的业务操作。

View：视图，对数据进行表示。

Controller：控制器，获取View的请求，调用模型进行操作，然后将结果进行展示

在我做过的一个前后端分离的清单小项目中，对MVC有一个比较深的体会。假设分为多个小组来实现它，一个小组负责前端，一个小组负责后端，他们的使用的语言甚至不一样代码可以完全分离开，这实现了代码低耦合的思想。在后端，我会设置一个model文件夹，存放和数据库进行数据迁移的模型以及可以操作这些数据实现具体业务的方法，将它们封装起来作为一个模型。然后再设置一个控制器controller文件夹，里面只存放着业务调用的逻辑。这样一来前端View通过选择控制器controller提供服务就可以让控制器对模型进行操作，实现具体的业务功能，之后将结果返回给前端Vie让他展现出来。这样就实现了高内聚低耦合的思想。这样去开发可以将软件用户界面和业务逻辑分离，使代码可扩展性，可复用性和可维护性都加强，真正实现高内聚，低耦合，而这也正是MVC的宗旨。而且MVC更容易去理解，所以现在的很多项目都会选择MVC这个设计模式。

 