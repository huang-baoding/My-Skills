# GoWeb

# 1  Web基础

### 1.1  URL

URL是“统一资源定位符”（Uniform Resource Locator）的缩写，用于描述一个网络上的资源。

```text
URL：scheme://host[:port#]/path/.../[?query-string][#anchor]
scheme         指定底层使用的协议(例如：http, https, ftp)
host           HTTP 服务器的 IP 地址或者域名
port#          HTTP 服务器的默认端口是 80，这种情况下端口号可以省略。如果使用了别的端口，则必须指明。
path           访问资源的路径
query-string   发送给 http 服务器的数据
anchor         锚

```

### 1.2  net/http 包

#### 1.2.1  用http包搭建一个Web服务器

> 在main.go中保存然后go run，然后在浏览器中输入http://localhost:9090

```go
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()								// 解析参数，默认是不会解析的
	fmt.Println(r.Form)							 // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") 				// 这个写入到 w 的是输出到客户端的
}

func main() {
	http.HandleFunc("/", sayhelloName)				//设置访问的路由
	err := http.ListenAndServe(":9090", nil)		// 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
```

#### 1.2.2  用http包发送GET、POST、UPDATE、DELETE请求

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func get() {
	r, err := http.Get("http://httpbin.org/get")//这个网站是用于简单测试HTTP请求的网站
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", content)
}

func post() {
	r, err := http.Post("http://httpbin.org/post", "", nil)
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", content)
}

func put() {
	request, err := http.NewRequest(http.MethodPut, "http://httpbin.org/put", nil)
	if err != nil {
		panic(err)
	}
	r, err := http.DefaultClient.Do(request) // enter 键
	if err != nil {
		panic(err)
	}

	defer func() { _ = r.Body.Close() }()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", content)
}

func del() {
	request, err := http.NewRequest(http.MethodDelete, "http://httpbin.org/delete", nil)
	if err != nil {
		panic(err)
	}
	r, err := http.DefaultClient.Do(request) // enter 键
	if err != nil {
		panic(err)
	}

	defer func() { _ = r.Body.Close() }()
	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", content)
}

func main() {
	post()
}
```

#### 1.2.3  对请求的设置

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func printBody(r *http.Response) {			//读取并打印HTTP响应体的内容
	defer func() { _ = r.Body.Close() }()
	content, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", content)
}

func requestByParams() {			//设置请求的查询参数
	request, err := http.NewRequest(http.MethodGet, "http://httpbin.org/get", nil)
	if err != nil {
		panic(err)
	}
	params := make(url.Values)
	params.Add("name", "poloxue")
	params.Add("age", "18")
	request.URL.RawQuery = params.Encode()
	r, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	printBody(r)
}

func requestByHead() {			//定制请求头
	request, err := http.NewRequest(http.MethodGet, "http://httpbin.org/get", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("user-agent", "chrome")
	r, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	printBody(r)
}

func main() {
	// 如何设置请求的查询参数，http://httpbin.org/get?name=poloxue&age=18
	// 如何定制请求头，比如修改 user-agent
	requestByParams()
	requestByHead()
}
```

#### 1.2.4  http的响应和编码信息

读取响应信息

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	r, err := http.Get("https://baidu.com")
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	//读取响应内容
	content, _ := io.ReadAll(r.Body)
	fmt.Printf("%s", content)

	//读取响应状态码和状态描述
	fmt.Println(r.StatusCode)
	fmt.Println(r.Status)

	//读取响应头
	fmt.Println(r.Header.Get("content-type"))
}
```

获取编码信息

> 一般网页的content-type 或者 html的head meta 会提供编码信息。如果都没有可以通过网页的头部猜网页的编码信息

```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func encoding(r *http.Response) {
	bufReader := bufio.NewReader(r.Body)
	bytes, _ := bufReader.Peek(1024) // 不会移动 reader 的读取位置

	e, _, _ := charset.DetermineEncoding(bytes, r.Header.Get("content-type")) //获取网页的编码信息
	fmt.Println(e)
	bodyReader := transform.NewReader(bufReader, e.NewDecoder())//通过编码信息去解码网页
	content, _ := io.ReadAll(bodyReader)
	fmt.Printf("%s", content)
}

func main() {
	r, err := http.Get("https://baidu.com")
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()
	encoding(r)
}
```

#### 1.2.5  使用http包进行文件的下载和进度显示

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Reader struct {
	io.Reader
	Total   int64 //总大小
	Current int64 //当前已下载大小
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	r.Current += int64(n)
	fmt.Printf("\r进度 %.2f%%", float64(r.Current*10000/r.Total)/100)

	return
}

func DownloadFileProgress(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	reader := &Reader{
		Reader: r.Body,
		Total:  r.ContentLength,
	}

	_, _ = io.Copy(f, reader)		//在读取这个reader的内容时，会触发它的Read方法
}

func main() {
	// 自动文件下载，比如自动下载图片、压缩包
	url := "https://img-home.csdnimg.cn/images/20240129062750.png"
	filename := "poloxue.png"
	DownloadFileProgress(url, filename)
}
```

#### 1.2.6  Post请求提交Form和Json

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func postForm() {
	data := make(url.Values)
	data.Add("name", "poloxue")
	data.Add("age", "18")
	payload := data.Encode()

	r, _ := http.Post(
		"http://httpbin.org/post",
		"application/x-www-form-urlencoded",
		strings.NewReader(payload),
	)
	defer func() { _ = r.Body.Close() }()

	content, _ := io.ReadAll(r.Body)
	fmt.Printf("%s", content)
}

func postJson() {
	u := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "poloxue",
		Age:  18,
	}
	payload, _ := json.Marshal(u)
	r, _ := http.Post(
		"http://httpbin.org/post",
		"application/json",
		bytes.NewReader(payload),
	)
	defer func() { _ = r.Body.Close() }()

	content, _ := io.ReadAll(r.Body)
	fmt.Printf("%s", content)
}

func main() {
	postForm()
	postJson()
}
```

#### 1.2.7  Post请求提交文件

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile() {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	_ = writer.WriteField("words", "123")

	// 一个是输入表单的 name，一个上传的文件名称
	upload1Writer, _ := writer.CreateFormFile("uploadfile1", "uploadfile1")

	uploadFile1, _ := os.Open("uploadfile1")
	defer func() { _ = uploadFile1.Close() }()

	_, _ = io.Copy(upload1Writer, uploadFile1)

	// 一个是输入表单的 name，一个上传的文件名称
	upload2Writer, _ := writer.CreateFormFile("uploadfile2", "uploadfile2")

	uploadFile2, _ := os.Open("uploadfile2")
	defer func() { _ = uploadFile2.Close() }()

	_, _ = io.Copy(upload2Writer, uploadFile2)

	_ = writer.Close()

	fmt.Println(writer.FormDataContentType())
	fmt.Println(body.String())
	r, _ := http.Post("http://httpbin.org/post",
		writer.FormDataContentType(),
		body,
	)
	defer func() { _ = r.Body.Close() }()

	content, _ := io.ReadAll(r.Body)

	fmt.Printf("%s", content)
}

func main() {
	postFile()
}
```

#### 1.2.8  http处理重定向

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
)

func redirectLimitTimes() {
	// 限制重定向的次数
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 {
				return errors.New("redirect too times")
			}
			return nil
		},
	}

	request, _ := http.NewRequest(
		http.MethodGet,
		"http://httpbin.org/redirect/20",
		nil,
	)
	_, err := client.Do(request)
	if err != nil {
		panic(err)
	}
}

func redirectForbidden() {
	// 禁止重定向; 登录请求，防止重定向到首页
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, _ := http.NewRequest(
		http.MethodGet,
		"http://httpbin.org/cookies/set?name=poloxue",
		nil,
	)
	r, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer func() {_ = r.Body.Close()}()
	fmt.Println(r.Request.URL)
}

func main() {
	redirectForbidden()
}
```

#### 1.2.9  http使用cookie

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

// 使用过程比较繁琐，主要是为了用来理解cookie的使用
func rrCookie() {
	// 模拟完成一个登录
	// 请求一个页面，传递基本的登录信息，将响应的 cookie 设置到下一次之上重新请求
	// 请求 http://httpbin.org/cookies/set?name=poloxue&password=123456
	// 返回 set-cookie:
	// 再一次请求呢携带上 cookie，
	// 首页 http://httpbin.org/cookies 就会通过 body 打印出已经设置 cookie
	// http://httpbin.org/cookies/set? => response
	// request => http://httpbin.org/cookies
	client := &http.Client{ //禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	firstRequest, _ := http.NewRequest( //第一次提交请求
		http.MethodGet,
		"http://httpbin.org/cookies/set?name=poloxue&password=123456",
		nil,
	)
	firstResponse, _ := client.Do(firstRequest)
	defer func() { _ = firstResponse.Body.Close() }()

	secondRequest, _ := http.NewRequest( //第二次提交请求
		http.MethodGet,
		"http://httpbin.org/cookies",
		nil,
	)

	for _, cookie := range firstResponse.Cookies() {
		secondRequest.AddCookie(cookie) //添加过程
	}

	secondResponse, _ := client.Do(secondRequest)
	defer func() { _ = secondResponse.Body.Close() }()

	content, _ := io.ReadAll(secondResponse.Body)
	fmt.Printf("%s\n", content)
}

func main() {
	rrCookie()
}

```

#### 1.2.10 http cookie持久化

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	cookiejar2 "github.com/juju/persistent-cookiejar"
)

func jarCookie() {
	jar, _ := cookiejar.New(nil) //创建一个系统自带的cookie
	client := &http.Client{
		Jar: jar,
	}
	r, _ := client.Get("http://httpbin.org/cookies/set?username=poloxue&password=123456")
	defer func() { _ = r.Body.Close() }()

	_, _ = io.Copy(os.Stdout, r.Body) //将请求copy到标准输出
}

func login(jar http.CookieJar) {

	client := &http.Client{
		Jar: jar,
	}
	r, _ := client.PostForm(
		"http://localhost:8080/login",
		url.Values{"username": {"poloxue"}, "password": {"poloxue123"}},
	)
	defer func() { _ = r.Body.Close() }()
	fmt.Println(r.Cookies())

	_, _ = io.Copy(os.Stdout, r.Body)
}

func center(jar http.CookieJar) {
	client := &http.Client{
		Jar: jar,
	}
	r, _ := client.Get("http://localhost:8080/center")
	defer func() { _ = r.Body.Close() }()

	_, _ = io.Copy(os.Stdout, r.Body)
}

func main() {
	// cookie 的分类有两种 一种是会话期 cookie 一种是持久性 cookie
	// jar, _ := cookiejar.New(nil)//标准库提供的cookie是会话期cookie
	jar, _ := cookiejar2.New(nil) //持久性cookie
	// login(jar)
	center(jar)
	_ = jar.Save() //使用持久性cookie需要对cookie进行保存
}

```

1.2.11  http设置超时时间和代理

```go
package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

func main() {
	// https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return net.DialTimeout(network, addr, 2*time.Second)
			},
			ResponseHeaderTimeout: 5 * time.Second,
			TLSHandshakeTimeout:   2 * time.Second,
			IdleConnTimeout:       60 * time.Second,
		},
	}
	_, _ = client.Get("http://httpbin.org/delay/10")
}

```

### 1.3  Web访问数据库

#### 1.3.1  database/sql 接口

```go
  import (
  	"database/sql"
   	_ "github.com/mattn/go-sqlite3"					//调用其中的init函数使用数据库驱动
  )
```

#### 1.3.2 使用MySQL数据库

数据准备

```mysql
CREATE TABLE `userinfo` (
	`uid` INT(10) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NULL DEFAULT NULL,
	`department` VARCHAR(64) NULL DEFAULT NULL,
	`created` DATE NULL DEFAULT NULL,
	PRIMARY KEY (`uid`)
);

CREATE TABLE `userdetail` (
	`uid` INT(10) NOT NULL DEFAULT '0',
	`intro` TEXT NULL,
	`profile` TEXT NULL,
	PRIMARY KEY (`uid`)
)
```

在go中对数据进行增删改查

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:5555)/dbname?charset=utf8")	//打开一个注册过的数据库驱动
	checkErr(err)

	// 插入数据
	stmt, err := db.Prepare("INSERT userinfo SET username=?,department=?,created=?")//得到将要执行的SQL；=？可以一定程度防止SQL注入
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")			//执行 stmt 准备好的SQL语句
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	// 更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// 查询数据
	rows, err := db.Query("SELECT * FROM userinfo")  	//直接执行查询SQL并返回查询结果
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	// 删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

```

#### 1.3.3  使用Redis数据库

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

var ( //声明全局变量Pool，存储连接池
	Pool *redis.Pool
)

func init() { //初始化连接池
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,                 //最大空闲连接数
		IdleTimeout: 240 * time.Second, //空闲连接的超时时间为240秒

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func close() {
	c := make(chan os.Signal, 1)   //接收操作系统信号的通道
	signal.Notify(c, os.Interrupt) //将通道与这些信号关联起来
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() { //启动一个通道等待信号
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

func main() {
	test, err := Get("test")
	fmt.Println(test, err)
}

```

### 1.4  cookie和session

> 使用HTTP协议的每次连接都是无状态的新连接，可以使用cookie和session来记录用户的历史信息以保持连接状态。

#### 1.4.1 session和cookie的区别

session：服务端机制，服务器检查客户端携带的session id去检索session。没有则会生成一个返回给客户端。

cookie：客户端机制，由服务器创建然后由浏览器保存，有会话期cookie和持久cookie

#### 1.4.2  Go使用session和cookie

```go
//设置cookie
expiration := time.Now()
expiration = expiration.AddDate(1, 0, 0)
cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
http.SetCookie(w, &cookie)

//读取 cookie
for _, cookie := range r.Cookies() {
	fmt.Fprint(w, cookie.Name)
}
```

### 1.5 文本处理

#### 1.5.1  Go和XML

> encoding/xml 是go专门用来处理XML文本的标准库。

```go
package main

import (
	"encoding/xml"
	"fmt"
)

type Person struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
	City    string   `xml:"city"`
}

func main() {
	xmlData := `
		<person>
			<name>Alice</name>
			<age>25</age>
			<city>New York</city>
		</person>`

	var p Person
	err := xml.Unmarshal([]byte(xmlData), &p)
	if err != nil {
		fmt.Println("解析XML失败:", err)
		return
	}

	fmt.Println("姓名:", p.Name)
	fmt.Println("年龄:", p.Age)
	fmt.Println("城市:", p.City)

	// 将结构体转换为XML
	xmlBytes, err := xml.MarshalIndent(p, "", "\t")
	if err != nil {
		fmt.Println("生成XML失败:", err)
		return
	}

	fmt.Println("生成的XML:")
	fmt.Println(string(xmlBytes))
}

/*输出结果:
姓名: Alice
年龄: 25
城市: New York
生成的XML:
<person>
	<name>Alice</name>
	<age>25</age>
	<city>New York</city>
</person>*/
```

#### 1.5.2  Go和Json

> encoding/json 是go专门用来处理XML文本的标准库。

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
	jsonData := `
		{
			"name": "Alice",
			"age": 25,
			"city": "New York"
		}`

	var p Person
	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return
	}

	fmt.Println("姓名:", p.Name)
	fmt.Println("年龄:", p.Age)
	fmt.Println("城市:", p.City)

	// 将结构体转换为JSON
	jsonBytes, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		fmt.Println("生成JSON失败:", err)
		return
	}

	fmt.Println("生成的JSON:")
	fmt.Println(string(jsonBytes))
}
/*输出结果：
姓名: Alice
年龄: 25
城市: New York
生成的JSON:
{
	"name": "Alice",
	"age": 25,
	"city": "New York"
}
*/
```

#### 1.5.3 Go和正则表达式

> regexp包是Go专门用来处理正则表达式的标准库。

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := "Hello, 12345 World!"
	pattern := "[0-9]+"

	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)

	fmt.Println("匹配结果:")
	for _, match := range matches {
		fmt.Println(match)
	}
}
```

### 1.6  Web服务

#### 1.6.1  Socket编程

> 流式Socket（SOCK_STREAM）：面向连接，对应TCP服务应用。
>
> 数据报式Socket（SOCK_DGRAM）：无连接，对应UDP服务应用。

TCP客户端

```go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {							//检查命令行参数是否为两个
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]							//获取命令行参数中的服务器地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)	//解析TCP地址
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)		//建立TCP连接
	checkError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))// 向服务器发送HTTP HEAD请求
	checkError(err)
	result, err := io.ReadAll(conn)						// 读取服务器响应
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}
func checkError(err error) {							//检查错误
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}


```

TCP服务端

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"strings"
)

func main() {
	service := ":1200" 									// 指定服务地址和端口号
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)		 // 解析TCP地址
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr) 			// 监听TCP连接
	checkError(err)
	for {
		conn, err := listener.Accept() // 接受客户端连接请求
		if err != nil {
			continue
		}
		go handleClient(conn) // 处理客户端请求
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) 	// 设置读取超时时间为2分钟
	request := make([]byte, 128)		 // 限制请求数据最大长度为128字节，以防止洪泛攻击
	defer conn.Close() 									// 函数结束前关闭连接
	for {
		read_len, err := conn.Read(request) 			// 读取客户端请求数据

		if err != nil {
			fmt.Println(err)
			break
		}

    		if read_len == 0 {
    			break // 客户端已关闭连接
    		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
                 // 获取当前时间戳并转换为字符串
    			daytime := strconv.FormatInt(time.Now().Unix(), 10) 
    			conn.Write([]byte(daytime)) // 将时间戳字符串发送回客户端
    		} else {
    			daytime := time.Now().String() // 获取当前时间并转换为字符串
    			conn.Write([]byte(daytime)) // 将时间字符串发送回客户端
    		}

    		request = make([]byte, 128) // 清空缓存数据
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

UDP客户端

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	_, err = conn.Write([]byte("anything"))
	checkError(err)
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)
	fmt.Println(string(buf[0:n]))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
```

UDP服务端

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":1200"                      			 // 定义服务地址和端口号
	udpAddr, err := net.ResolveUDPAddr("udp4", service)   // 解析UDP地址
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr) 			 // 监听UDP连接
	checkError(err)
	for {
		handleClient(conn) 								// 处理客户端请求
	}
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte 									// 定义一个长度为512字节的byte数组
	_, addr, err := conn.ReadFromUDP(buf[0:]) 		// 从UDP连接中读取数据，并获取客户端地址信息
	if err != nil {
		return 								
	}
	daytime := time.Now().String() 							// 获取当前时间并转换为字符串
	conn.WriteToUDP([]byte(daytime), addr) 					// 将时间字符串发送回客户端
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
```

#### 1.6.2 WebSocket

> WebSocket基于socket实现了浏览器和服务器的全双工通信

客户端

```html
<html>
<head></head>
<body>
	<script type="text/javascript">
		var sock = null; // 定义一个变量用于存储WebSocket对象
		var wsuri = "ws://127.0.0.1:1234"; // 定义WebSocket服务器的地址

		window.onload = function() {
			console.log("onload"); // 在加载页面时输出日志

			sock = new WebSocket(wsuri); // 创建WebSocket对象，并与服务器建立连接

			sock.onopen = function() { // 当WebSocket连接成功打开时执行的回调函数
				console.log("connected to " + wsuri); // 输出连接成功的日志
			}

			sock.onclose = function(e) { // 当WebSocket连接关闭时执行的回调函数
				console.log("connection closed (" + e.code + ")"); // 输出连接关闭的日志，显示关闭的原因（code）
			}

			sock.onmessage = function(e) { // 当接收到WebSocket消息时执行的回调函数
				console.log("message received: " + e.data); // 输出接收到的消息
			}
		};

		function send() {
			var msg = document.getElementById('message').value; // 获取输入框中的文本内容作为消息
			sock.send(msg); // 发送消息到服务器
		};
	</script>
	<h1>WebSocket Echo Test</h1>
	<form>
		<p>
			Message: <input id="message" type="text" value="Hello, world!"> <!-- 显示一个输入框，初始值为"Hello, world!" -->
		</p>
	</form>
	<button onclick="send();">Send Message</button> <!-- 显示一个按钮，点击按钮会调用send函数发送消息 -->
</body>
</html>

```

服务端

```go
package main

import (
	"golang.org/x/net/websocket" // 导入websocket包
	"fmt"
	"log"
	"net/http"
)

func Echo(ws *websocket.Conn) { 				// 定义Echo函数，处理WebSocket连接
	var err error

	for {
		var reply string 						// 定义变量reply用于存储接收到的消息

		if err = websocket.Message.Receive(ws, &reply); err != nil { // 接收客户端发送的消息
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply) // 打印接收到的消息

		msg := "Received:  " + reply // 构造要发送给客户端的消息
		fmt.Println("Sending to client: " + msg) // 打印要发送的消息

		if err = websocket.Message.Send(ws, msg); err != nil { // 将消息发送给客户端
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo)) // 将Echo函数注册为处理WebSocket连接的处理器

	if err := http.ListenAndServe(":1234", nil); err != nil { // 启动HTTP服务器并监听指定端口
		log.Fatal("ListenAndServe:", err)
	}
}

```

#### 1.6.3 REST

> REST 是一种架构约束条件和原则，满足这些条件和原则就叫RESTful的。

RESTful架构：

- 每一个 URI 代表一种资源。
- 客户端和服务器之间，传递这种资源的某种表现层。
- 客户端通过四个 HTTP 动词，对服务器端资源进行操作，实现 "表现层状态转化"。

#### 1.6.4  RPC

> RPC（Remote Procedure Call Protocol）：远程过程调用协议。是一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议。



### 1.7  安全和加密

#### 1.7.1  预防CSRF攻击

> CSRF（Cross-site request forgery）：跨站请求伪造。攻击者可以盗用你的登录信息，以你的身份模拟发送各种请求。

CSRF攻击原理：

- 用户在浏览器中登录了一个网站A，并保持了有效的会话。
- 用户在同一浏览器中打开了另一个标签页或者窗口，访问了一个恶意网站B。
- 恶意网站B中的代码包含了对网站A的请求，例如提交表单、发起GET/POST请求等。
- 浏览器自动携带了用户在网站A中的身份凭证（如Cookie）发送了请求，由于浏览器会自动携带相关的凭证，网站A无法判断这个请求是来自用户的意愿还是恶意网站B的伪造。
- 网站A接收到了这个请求，并以为它是用户的合法请求，执行了相应的操作，比如修改用户信息、删除数据等。

预防CSRF攻击的措施：

- 添加CSRF令牌：在每个请求中添加一个随机生成的CSRF令牌，并在服务器端验证该令牌的有效性。攻击者无法获取到正确的CSRF令牌，因此不能伪造合法的请求。
- 同源检测：服务器可以通过检查请求的来源（Referer）确保请求来自合法的网站，但这种方式并不绝对可靠，因为Referer也可以被伪造。
- 使用SameSite Cookie属性：将Cookie的SameSite属性设置为Strict或Lax，限制第三方网站对Cookie的访问，从而减少CSRF攻击的风险。
- 验证HTTP请求头：服务器可以验证请求中的自定义头部字段，以确保请求是合法的。

#### 1.7.2  预防XXS攻击

> XSS  (Cross-Site Scripting)：跨站脚本攻击，它允许攻击者将恶意代码植入到提供给其它用户使用的页面中。

XXS的攻击方式：

- 反射型XSS：攻击者将包含恶意脚本的链接发送给用户，用户点击链接后，恶意脚本被执行，攻击者就可以获取用户的敏感信息。
- 存储型XSS：攻击者将恶意脚本存储到数据库中，当用户访问受影响的页面时，恶意脚本会被执行，攻击者就可以获取用户的敏感信息。
- DOM-based XSS：攻击者通过修改网页中的DOM元素来执行恶意代码，从而获取用户的敏感信息。

XXS攻击的原理：

- 攻击者在目标网站上注入恶意脚本，通常是利用表单、评论、搜索等交互性操作，将脚本代码插入到网页中。
- 用户在浏览器中访问被注入恶意脚本的网站，此时浏览器会自动执行其中的JavaScript代码。
- 恶意脚本可以做很多事情，比如盗取 Cookie、获取用户输入的数据、修改网页内容等，攻击者可以利用这些功能进行各种形式的攻击。

 预防XXS攻击的措施：

- 输入过滤：对用户输入的内容进行过滤和转义，过滤掉一些特殊字符和脚本代码，从而减少恶意脚本注入的风险。
- 使用HTTPOnly Cookie属性：将Cookie的HTTPOnly属性设置为true，限制JavaScript对Cookie的访问，从而减少Cookie泄露的风险。
- CSP（Content Security Policy）：使用CSP来限制网站中可以执行的脚本，只允许来自指定源的脚本被执行，从而减少恶意脚本的注入和执行。
- 验证输出数据：在输出数据时，对数据进行验证和过滤，确保不会执行任何恶意代码。

#### 1.7.3  预防SQL注入

> 程序没有有效过滤用户的输入，而让攻击者成功向服务器提交了恶意的SQL代码。

预防SQL注入的措施：

- 输入过滤：对用户输入的内容进行过滤和转义，过滤掉一些特殊字符和SQL代码，从而减少恶意代码注入的风险。
- 参数化查询：使用参数化查询语句，将用户输入的参数作为参数传递给SQL查询语句，从而避免拼接SQL语句时出现恶意代码。
- 最小权限原则：数据库用户的权限应该尽可能地低，只有必要的最小权限才能够执行相应的操作。
- 安全编码：开发者需要了解和遵守安全编码的规则，比如使用预编译语句、使用ORM等方式来避免注入漏洞。
- 在应用发布之前先用专业的SQL注入检测工具检测
- 避免网站打印出SQL错误信息。

#### 1.7.4  加密和解密

### 1.8  错误处理、调试和测试

#### 1.8.1  错误处理

```go
//error接口
type error interface {
	Error() string
}
//实现error接口的errorString结构体
type errorString struct {
	s string
}
func (e *errorString) Error() string {
	return e.s
}
//errors.New方法的实现
func New(text string) error {
	return &errorString{text}
}

//返回错误的例子
if i<0 {
    return 0,errors.New("这是生成的一个错误！")
}
```

#### 1.8.2  使用GBD调试

> GDB 是 FSF (自由软件基金会) 发布的一个强大的类 UNIX 系统下的程序调试工具。

GDB可以做的事：

- 启动程序，可以按照开发者的自定义要求运行程序。
- 可让被调试的程序在开发者设定的调置的断点处停住。（断点可以是条件表达式）
- 当程序被停住时，可以检查此时程序中所发生的事。
- 动态的改变当前程序的执行环境。

#### 1.8.3  Go如何写测试用例

> Go 语言中自带有一个轻量级的测试框架 `testing` 和自带的 `go test` 命令来实现单元测试和性能测试。

使用testing包的注意事项：

- 文件名必须是 `_test.go` 结尾的，这样在执行 `go test` 的时候才会执行到相应的代码
- 所有的测试用例函数必须是 `Test` 开头
- 测试用例会按照源代码中写的顺序依次执行
- 测试函数 `TestXxx()` 的参数是 `testing.T`，我们可以使用该类型来记录错误或者是测试状态
- 测试格式：`func TestXxx (t *testing.T)`, `Xxx` 部分可以为任意的字母数字的组合，但是首字母不能是小写字母 [a-z]，例如 `Testintdiv` 是错误的函数名。
- 函数中通过调用 `testing.T` 的 `Error`, `Errorf`, `FailNow`, `Fatal`, `FatalIf` 方法，说明测试不通过，调用 `Log` 方法用来记录测试的信息。
- 完成编码后在命令行输入`go test [-v]` 来执行代码

### 1.9  部署和维护

#### 1.9.1  应用日志

> 如果我们想把我们的应用日志保存到文件，然后又能够结合日志实现很多复杂的功能，可以使用第三方开发的日志系统: [logrus (opens new window)](https://github.com/sirupsen/logrus)和 [seelog](https://github.com/cihub/seelog)

#### 1.9.2  网站错误处理

panic-recover

#### 1.9.3应用部署

> 针对 Go 的应用程序部署，我们可以利用第三方工具来管理，例如Supervisord、upstart、daemontools等。

**1.9.3.1  使用Supervisord部署**

> Supervisord 是用 Python 实现的一款非常实用的进程管理工具。supervisord 会帮你把管理的应用程序转成 daemon 程序，而且可以方便的通过命令开启、关闭、重启等操作，而且它管理的进程一旦崩溃会自动重启，这样就可以保证程序执行中断后的情况下有自我修复的功能。

> 使用Supervisord 时，当修改了操作系统的文件描述符之后，需要重启Supervisord，光重启下面的应用程序没用。

#### 1.9.4  备份和恢复

## 2  Gin框架

go的框架其实是可以理解为库，并不是用了某一个框架就不能用别的框架，可以选择性的使用各个库中的优秀组件，进行组合

### 2.1  Gin框架基础

#### 2.1.1  Gin框架简介

Gin框架是用Go语言开发的Web框架，比net/http有着更加强大的功能。它的主要特点有：

（1）高性能：Gin采用了一些优化策略，例如使用了Radix树来快速路由请求，减少了内存的使用。

（2）易用性：Gin提供了许多简洁易懂的API。

（3）可扩展性：Gin可以通过各种中间件实现各种扩展功能（认证等）。

Gin的工作原理可以简单概括为：初始化一个路由器对象，将请求路由到特定的处理程序中然后在处理程序中进行业务逻辑的处理，例如对数据进行操作等然后返回结果给客户端。

> Gin还支持Crash处理、Json验证、路由组、错误处理和内置渲染等。
>
> - Crash处理：Gin 可以 catch 一个发生在 HTTP 请求中的 panic 并 recover 它。
> - Gin 可以解析并验证请求的 JSON，例如检查所需值的存在。
> - 更好地组织路由。
> - Gin 提供了一种方便的方法来收集 HTTP 请求期间发生的所有错误。最终，中间件可以将它们写入日志文件、数据库并通过网络发送。

#### 2.1.2  Gin框架的简单使用

**下载Gin框架**

`go get -u github.com/gin-gonic/gin`

> 使用go get命令会下载到文件GOPATH/pkg/mod/文件里。其它项目可以共用。

**Gin框架的简单使用**

```go
package main
import (
	"github.com/gin-gonic/gin"
)
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	r.Run(":8080")
}


//两种初始化方式：
r:=gin.Default( )     	 //初始化，返回一个带有默认值的路由地址：&Engine{...}
r:=gin.New( )              //返回一个没有默认值（空的）的&Engine{ }
//路由: 当输入localhost:8080/hello时会执行下面这个方法（对应第一个参数）
r.GET("/hello",func(r gin.Context){...})
```

**HTTP 方法**

> 使用时应该尽量遵循其语义

| 方法    | 作用                                                         |
| ------- | ------------------------------------------------------------ |
| GET     | 请求一个指定资源的表示形式，使用GET的请求应该只被用于获取数据 |
| POST    | 用于将实体提交到指定的资源，通常会导致在服务器上的状态变化   |
| PUT     | 用请求有效载荷替换目标资源的所有当前表示                     |
| DELETE  | 删除指定的资源                                               |
| HEAD    | 请求一个与GET请求的响应相同的响应，但没有响应体              |
| CONNECT | 建立一个到由目标资源标识的服务器的隧道                       |
| OPTIONS | 用于描述目标资源的通信选项                                   |
| TRACE   | 沿着到目标资源的路径执行一个消息环回测试                     |
| PATCH   | 用于对资源应用部分修改                                       |
| Any     | 支持所有（Gin框架）                                          |

> RESTful API规范（开发时应遵循此规范）
>
> - GET        表示读取服务器上的资源
> - POST      表示在服务器上创建资源
> - PUT        表示更新或者替换服务器上的资源
> - DELETE  表示删除服务器上的资源
> - PATCH    表示更新/修改资源的一部分

**读取文件**

func ReadFile(filename string)([]byte,error)

ReadFile 从filename指定的文件中读取数据并返回文件的内容。

#### 2.1.3 分组路由

> 在开发时，我们需要进行模块的划分（用户模块，商品模块等），我们可以进行对应的路由分组。

```go
ug := r.Group("/user")
	{
		ug.GET("find", func(ctx *gin.Context) {
			ctx.JSON(200, "user find")
		})
		ug.POST("save", func(ctx *gin.Context) {
			ctx.JSON(200, "user save")
		})
	}
	gg := r.Group("/goods")
	{
		gg.GET("find", func(ctx *gin.Context) {
			ctx.JSON(200, "goods find")
		})
		gg.POST("save", func(ctx *gin.Context) {
			ctx.JSON(200, "goods save")
		})
	}
```

#### 2.1.4  GET请求参数

**GET请求普通参数**

```go
//request url: http://localhost:8080/user/save?id=11&name=zhangsan
r.GET("/user/save", func(ctx *gin.Context) {
		id := ctx.Query("id")						//普通获取
		name,ok := ctx.GetQuery("name")				 //通过ok检查是否能获取到
		address := ctx.DefaultQuery("address", "北京") //参数不存在时，设置默认值
		ctx.JSON(200, gin.H{
			"id":      id,
             "ok":		ok,
			"name":    name,
			"address": address,
		})
	})

//输出到结构体
type User struct {
	Id   int64  `form:"id"`
	Name string `form:"name"`
}
r.GET("/user/save", func(ctx *gin.Context) {
		var user User
  		err := ctx.BindQuery(&user)				//也可以用ctx.ShouldBindQuery()
		if err != nil {
			log.Println(err)
		}
		ctx.JSON(200, user)
})
//BindQuery()如果绑定出错需要进行错误处理；ShouldBindQuery()如果绑定出错可以忽略掉
```

**GET请求数组参数**

```go
//请求url：http://localhost:8080/user/save?address=Beijing&address=shanghai
r.GET("/user/save", func(ctx *gin.Context) {
		address := ctx.QueryArray("address")
   		//address, ok := ctx.GetQueryArray("address")
    	//err := ctx.ShouldBindQuery(&user)
		ctx.JSON(200, address)
    	
	})

```

**GET请求map参数**

```go
//请求url：http://localhost:8080/user/save?ddressMap[home]=Beijing&addressMap[company]=shanghai
r.GET("/user/save", func(ctx *gin.Context) {
		addressMap := ctx.QueryMap("addressMap")
    	//addressMap, _ := ctx.GetQueryMap("addressMap")
		ctx.JSON(200, addressMap)
	})
//map参数 不能使用BindQuery()和ShouldBindQuery()
```

#### 2.1.5  Post请求参数

**Post请求表单参数**

```go
//Postman提交地址：http://localhost:8080/user/save
//提交类型：Body
//提交数据类型：form-data
//提交内容：键 ----  值
//id				---		1001
//name				---	 	Tom
//address 			---		 beijing
//address 			---		 dalian
//addressMap[home]	 ---	  nanning

r.POST("/user/save", func(ctx *gin.Context) {
		id := ctx.PostForm("id")
		name := ctx.PostForm("name")
		address := ctx.PostFormArray("address")
		addressMap := ctx.PostFormMap("addressMap")
		ctx.JSON(200, gin.H{
			"id":         id,
			"name":       name,
			"address":    address,
			"addressMap": addressMap,
		})
	})

r.POST("/user/save", func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBind(&user)
		addressMap, _ := ctx.GetPostFormMap("addressMap")
		user.AddressMap = addressMap
		fmt.Println(err)
		ctx.JSON(200, user)
	})
```

**Post请求json参数**

```json
{
    "id":1111,
    "name":"zhangsan",
    "address": [
        "beijing",
        "shanghai"
    ],
    "addressMap":{
        "home":"beijing"
    }
}
```

```go
r.POST("/user/save", func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBindJSON(&user)
		fmt.Println(err)
		ctx.JSON(200, user)
	})
```

**Post请求路径参数**

```go
//请求url：http://localhost:8080/user/save/111
r.POST("/user/save/:id", func(ctx *gin.Context) {
		ctx.JSON(200, ctx.Param("id"))
	})
```

**Post请求文件参数**

```go
r.POST("/user/save", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			log.Println(err)
		}
		files := form.File
		for _, fileArray := range files {
			for _, v := range fileArray {
				ctx.SaveUploadedFile(v, "./"+v.Filename)
			}

		}
		ctx.JSON(200, form.Value)
	})
```

#### 2.1.6  响应方式

**&gin.Context的常用方法**

| 方法                                                     | 返回值                                                       |
| -------------------------------------------------------- | ------------------------------------------------------------ |
| JSON(code int, obj interface{})                          | JSON 格式的响应                                              |
| String(status int, format string, values ...interface{}) | 字符串响应                                                   |
| HTML(status int, name string, data interface{})          | 页面响应                                                     |
| Data(status int, contentType string, data []byte)        | 二进制数据响应                                               |
| Redirect(status int, location string)                    | 重定向响应                                                   |
| AbortWithStatus(status int)                              | 终止请求处理，并返回指定状态码的响应                         |
| Param(key string) string                                 | 获取指定名称的 URL 参数                                      |
| Query(key string) string                                 | 获取指定名称的查询参数                                       |
| PostForm(key string) string                              | 获取指定名称的表单数据参数                                   |
| BindJSON(obj interface{}) error                          | 绑定 JSON 格式的请求体到指定的结构体对象                     |
| BindQuery(obj interface{}) error                         | 绑定查询参数到指定的结构体对象                               |
| Bind(obj interface{}) error                              | 绑定请求参数（查询参数、表单数据、JSON 等）到指定的结构体对象 |

```go
//字符串响应
r.GET("/user/save", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})
//json响应
r.GET("/user/save", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "this is a %s", "ms string response")
	})
//xml响应
type XmlUser struct {
	Id   int64  `xml:"id"`
	Name string `xml:"name"`
}
r.GET("/user/save", func(ctx *gin.Context) {
		u := XmlUser{
			Id:   11,
			Name: "zhangsan",
		}
		ctx.XML(http.StatusOK, u)
	})
```

#### 2.1.7  渲染模板

**使用Gin框架对现有的前端项目模板进行渲染**

在main.go的目录新建一个template文件夹存放以下模板

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>gin_templates</title>
</head>
<body>
{{.title}}
</body>
</html>
```

渲染这个模板：

```go
func main() {
	r := gin.Default()
	// 模板解析
	r.LoadHTMLFiles("templates/index.tmpl")
	r.GET("/index", func(c *gin.Context) {
		// HTML请求
		// 模板的渲染
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "hello 模板",
		})
	})
	r.Run(":9090") // 启动server
}

```

当有templates文件夹下多个模板文件时

```go
r.LoadHTMLGlob("templates/**")				//解析templates目录下的所有模板文件
```

如果目录为`templates/post/index.tmpl`和`templates/user/index.tmpl`这种，可以这样使用：

```go
router.LoadHTMLGlob("templates/**/*")	  	  // **/* 代表所有子目录下的所有文件
```

如果在模板中引入了静态文件（如.css文件)

```go
//模板文件中引用静态文件： <link rel="stylesheet" href="/css/index.css">
r.Static("/css", "./static/css")			//加载静态文件
```

此时我们可以这样渲染模板：

```go
func main() {
r := gin.Default()
    // 告诉gin框架模板文件引用的静态文件去哪里找,前一个参数是当前项目的哪个目录，
    r.Static("/static", "static")
    // 告诉gin框架去哪里找模板文件
    r.LoadHTMLGlob("templates/*")
    //把模板返回并渲染给浏览器
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
})
    r.Run(":9090")
}
//方式2
```

自定义模板函数

```go
/*模板：
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>gin_templates</title>
</head>
<body>
{{.title | safe}}
</body>
</html>
*/

   // gin框架给模板添加自定义函数
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	// 模板解析,解析templates目录下的所有模板文件
	r.LoadHTMLGlob("templates/**")

	r.GET("/index", func(c *gin.Context) {
		// HTML请求和模板的渲染
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "<a href='http://baidu.com'>跳转到其他地方</a>",
		})
	})
```

#### 2.1.8  Gin的cookie和session

**设置cookie**

> func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)

参数说明：

| 参数名   | 类型   | 说明                                                         |
| -------- | ------ | ------------------------------------------------------------ |
| name     | string | cookie名字                                                   |
| value    | string | cookie值                                                     |
| maxAge   | int    | 有效时间，单位是秒，MaxAge=0 忽略MaxAge属性，MaxAge\<0 相当于删除cookie,通常可以设置-1代表删除，MaxAge>0 多少秒后cookie失效 |
| path     | string | cookie路径                                                   |
| domain   | string | cookie作用域                                                 |
| secure   | bool   | Secure=true，那么这个cookie只能用https协议发送给服务器       |
| httpOnly | bool   | 设置HttpOnly=true的cookie不能被is获取到                      |

```go
r.GET("/cookie", func(c *gin.Context) {
		// 设置cookie
		c.SetCookie("site_cookie", "cookievalue", 3600, "/", "localhost", false, true)
	})
```

**读取cookie**

```go
r.GET("/read", func(c *gin.Context) {
		// 根据cookie名字读取cookie值
		data, err := c.Cookie("site_cookie")
		if err != nil {
			// 直接返回cookie值
			c.String(200,data)
			return
		}
		c.String(200,"not found!")
	})
```

**删除cookie**

```go
r.GET("/del", func(c *gin.Context) {
		// 通过将cookie的MaxAge设置为-1, 达到删除cookie的目的。
		c.SetCookie("site_cookie", "cookievalue", -1, "/", "localhost", false, true)
		c.String(200,"删除cookie")
	})

```

**session**

> 在Gin中，我们可以依赖gin-contrib/sessions中间件处理session

```go
package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 创建基于cookie的存储引擎，secret 参数是用于加密的密钥
	store := cookie.NewStore([]byte("secret"))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c *gin.Context) {
		// 初始化session对象
		session := sessions.Default(c)
		// 通过session.Get读取session值
		// session是键值对格式数据，因此需要通过key查询数据
		if session.Get("hello") != "world" {
			fmt.Println("没读到")
			// 设置session数据
			session.Set("hello", "world")
			session.Save()
		}
		c.JSON(200, gin.H{"hello": session.Get("hello")})
	})
	r.Run(":8080")
}

```

**多session**

```go
package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	sessionNames := []string{"a", "b"}
	r.Use(sessions.SessionsMany(sessionNames, store))

	r.GET("/hello", func(c *gin.Context) {
		sessionA := sessions.DefaultMany(c, "a")
		sessionB := sessions.DefaultMany(c, "b")

		if sessionA.Get("hello") != "world!" {
			sessionA.Set("hello", "world!")
			sessionA.Save()
		}

		if sessionB.Get("hello") != "world?" {
			sessionB.Set("hello", "world?")
			sessionB.Save()
		}

		c.JSON(200, gin.H{
			"a": sessionA.Get("hello"),
			"b": sessionB.Get("hello"),
		})
	})
	r.Run(":8080")
}
```

**基于redis存储引擎的session**

> 如果我们想将session数据保存到redis中，只要将session的存储引擎改成redis即可。

```go
package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 初始化基于redis的存储引擎
	// 参数说明：
	//    第1个参数 - redis最大的空闲连接数
	//    第2个参数 - 数通信协议tcp或者udp
	//    第3个参数 - redis地址, 格式，host:port
	//    第4个参数 - redis密码
	//    第5个参数 - session加密密钥
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8080")
}
```

#### 2.1.9  Gin使用中间件

> **中间件**（Middleware）指的是可以拦截**http请求-响应**生命周期的特殊函数

中间件的常见应用场景如下：

- 请求限速
- api接口签名处理
- 权限校验
- 统一错误处理

**中间件使用**

```go
   r := gin.New()
	// 通过use设置全局中间件
	// 设置日志中间件，主要用于打印请求日志
	r.Use(gin.Logger())
	// 设置Recovery中间件，主要用于拦截paic错误，不至于导致进程崩掉
	r.Use(gin.Recovery())
	r.GET("/test", func(ctx *gin.Context) {
		panic(errors.New("test error"))
	})
	r.Run(":8080")

```

**自定义中间件**

```go
package main
// 导入gin包
import (
"github.com/gin-gonic/gin"
	"log"
	"time"
)

// 自定义个日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 可以通过上下文对象，设置一些依附在上下文对象里面的键/值数据
		c.Set("example", "12345")

		// 在这里处理请求到达控制器函数之前的逻辑
     
		// 调用下一个中间件，或者控制器处理函数，具体得看注册了多少个中间件。
		c.Next()

		// 在这里可以处理请求返回给用户之前的逻辑
		latency := time.Since(t)
		log.Print(latency)

		// 例如，查询请求状态吗
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	r := gin.New()
	// 注册上面自定义的日志中间件
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		// 查询我们之前在日志中间件，注入的键值数据
		example := c.MustGet("example").(string)
		// it would print: "12345"
		log.Println(example)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}

```





### 2.2  Gin框架的实现原理

### Gin拦截器的原理

1. 注册全局中间件：在Gin应用程序初始化时，可以注册一个或多个全局中间件。这些中间件将被应用到所有的请求上。
2. 注册路由级别中间件：在定义路由时，可以为每个路由注册特定的中间件。这些中间件只会应用于与其相关联的路由。
3. 请求阶段：当请求到达服务器时，Gin会按照注册的顺序依次执行中间件。每个中间件都可以在请求处理之前或之后对请求进行处理。
4. 响应阶段：当请求处理完毕后，Gin会按照相反的顺序依次执行中间件。这允许中间件对响应进行修改或处理。
5. 错误处理：如果在处理过程中发生错误，Gin会跳过当前的中间件，直接进入错误处理流程。

### Gin的路由是如何实现的

Gin的路由是通过`gin.Engine`结构体实现的，`gin.Engine`结构体内部包含一个`[]*RouteGroup`类型的路由组数组。每个路由组（`RouteGroup`）又包含一个`[]IRoutes`类型的路由数组和一个`HandlersChain`类型的中间件链。

当我们使用Gin定义路由时，实际上是在创建一个新的`RouteGroup`对象，并将其添加到`gin.Engine`的路由组数组中。然后，我们可以通过`RouteGroup`对象的`GET()`、`POST()`等方法来定义具体的路由，并通过`Use()`、`Handle()`等方法添加相应的中间件。

在请求处理过程中，Gin会遍历路由组数组，匹配请求的URL路径和HTTP方法，找到对应的路由对象。如果找到，则执行该路由对象所包含的中间件链以及最终的请求处理函数；如果没有找到，则返回404错误响应。

Gin的路由匹配算法采用了`httprouter`库，它使用基于前缀树的方法进行高效的路由匹配。这种方法能够快速地找到与请求路径最匹配的路由对象，从而提高请求处理的效率。

### Gin的路由使用的数据结构（字典树），介绍一下字典树

Gin的路由使用了一种高效的数据结构称为`Radix Tree`，而不是传统的字典树（Trie）;它是一种基于前缀树的数据结构，也称为压缩前缀树（Compressed Prefix Tree）。与传统的字典树相比，Radix Tree采用了路径压缩的方式来减少存储空间和提高查找效率。具体来说，在Radix Tree中，相邻的节点如果只有一个子节点的话，就会将它们合并成一个节点，并用边上的字符表示这一段路径上的所有字符。这样，就可以避免在树中重复存储相同的字符，从而减少了内存消耗，提高了路由匹配的性能。

### gin常用中间件

1. Logger：记录HTTP请求的日志信息，包括请求方法、路径、状态码等。
2. Recovery：在发生panic时恢复应用程序并返回一个500 Internal Server Error响应。它可以提高应用程序的稳定性。
3. CORS：处理跨域资源共享（Cross-Origin Resource Sharing），允许不同域上的客户端访问API。
4. Static：用于提供静态文件服务，例如CSS、JavaScript或图像文件。
5. Auth（授权）：用于验证用户身份和权限的中间件，通常使用令牌或会话来实现身份验证和授权。
6. Rate Limiting（限流）：限制API请求的速率，防止滥用和恶意攻击。
7. Request ID：为每个请求生成唯一的请求ID，方便跟踪和调试。
8. JWT Middleware：处理JSON Web Token（JWT）的验证和解析，用于身份验证和授权。
9. Gzip：对HTTP响应进行gzip压缩，减小传输数据的大小。
10. Secure（安全）：应用一些常见的安全措施，如HTTP Strict Transport Security（HSTS）、X-Content-Type-Options、X-XSS-Protection等。