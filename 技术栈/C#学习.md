 # C#编辑器

## 1  VSCode

### 1.1  初始化项目

`dotnet new console [-n newFileName]`

初始化后的项目将有如下文件和文件夹：

* Program.cs：这是主要的代码文件，它包含了应用程序的入口点和主要逻辑。
* obj文件夹 ：这是用于存储编译过程中的中间文件和生成的可执行文件的文件夹。
* .csproj 文件：这是项目文件，它用于定义项目的结构、依赖关系和构建设置。
* bin 文件夹：这是用于存储最终生成的可执行文件的文件夹。
* .vscode 文件夹：这个文件夹包含一些用于配置和管理VSCode的设置文件。

>  编写代码后使用`dotnet build`来编译代码，使用`dotnet run`来运行代码。

### 1.2 初始化Web API项目

`dotnet new webapi [--framework net6.0]`

初始化后的项目将有如下文件和文件夹：

* Program.cs：这是项目的入口点，负责创建和运行主机，以及配置应用程序的启动和生命周期。


* Controllers/：包含所有控制器类的文件夹，用于处理请求和响应。
* appsettings.json：这是项目的设置文件，负责存储应用程序的配置信息，如连接字符串，日志级别，环境变量等。
* WeatherForecast.cs：这是项目的模型文件，定义了一个表示天气预报的类，包含日期，温度，摘要等属性。
*  Properties：这是一个文件夹，存放项目的属性文件，如 launchSettings.json，该文件负责定义项目的启动设置，如应用程序 URL，环境变量，启动参数等。
*  obj：这是一个文件夹，存放项目的中间文件，如项目文件的快照，编译输出，生成日志等。该文件夹通常不需要手动修改或查看，而是由 dotnet 命令或 IDE 自动管理。
*  bin：这是一个文件夹，存放项目的二进制文件，如可执行文件，动态链接库，符号文件等。
*   .csproj：这是一个文件，用于定义项目的元数据，如项目的名称，目标框架，包引用，项目引用等。

## 2  Visual Studio

### 2.1  初始化项目

### 2.2 初始化Web API项目

创建项目 -->  使用ASP.NET Core Web API-->输入项目名字选择文件夹

-->把模型示例文件WeatherForcastController.cs和WeatherForcast.cs文件删掉

-->在Controllers文件夹中添加新控制器文件ValuesController.cs文件：

```c#
using Microsoft.AspNetCore.Http;					//引入命名空间
using Microsoft.AspNetCore.Mvc;

namespace WebApiStudy1.Controllers
{
    [Route("api/[controller]/[action]")]			//初始化访问配置
    [ApiController]
    public class ValuesController : ControllerBase	 //必须要继承控制器基类
    {
        [HttpGet]									//处理的是GET请求
        public string Test()						//访问这个方法时会触发的逻辑代码
        {
            return "Hello Wrorld!";
        }
    }
}
```

-->选择浏运行浏览器并点击运行会在浏览器上打开一个Swagger的网页如`https://localhost:44307/swagger/index.html`，可以简单测试。

-->也可以将‘https://localhost:44307/api/Values/Test’在新的浏览器页面打开，会触发Test方法逻辑代码。

-->也可以在PostMan中输入‘https://localhost:44307/api/Values/Test’直接发送GET请求触发Test方法逻辑代码。



# C#学习

# 1  C#概述

C#是.NET 框架中的一种面向对象、类型安全、静态类型&解释型 语言，用于构建各种类型的应用程序。

.NET是整个框架的名称，包括多种编程语言、运行时和库。

.NETFramework是一个.NET 框架的版本，主要用于构建 Windows 应用程序。

ASP.NET是.NET 框架的一部分，用于构建 Web 应用程序。

ASP.NETCore是.NET 框架的新一代 Web 框架，跨平台并具有更高的性能。

ML.NET是.NET 框架的机器学习框架，用于在应用程序中集成机器学习功能。

## 1.1  基础命令

初始化C#文件夹：dotnet new console [-n fileName]

## 1.2  C#代码结构

最基础的C#的代码结构（应包括命名空间、类、主方法）

```C#
using System;          //引入命名空间

namespace MyNamespace  //此程序命名空间
{
    class Program
    {
        static void Main(string[] args) //主方法，程序的入口；唯一。
        {
            try
            {
                     Console.WriteLine("Hello, C#!");
            }
            catch ()
            {
            }
        }
    }
}
```





# 2  C#的数据类型

## 2.1  C#数据类型分类

|                                | 分类                                                | 具体类型                                           |
| ------------------------------ | --------------------------------------------------- | -------------------------------------------------- |
| 值类型                         | 整数类型                                            | sbyte、byte、short、ushort、int、uint、long、ulong |
| 浮点数类型                     | float、double、decimal                              |                                                    |
| 字符类型                       | char                                                |                                                    |
| 布尔类型                       | bool                                                |                                                    |
| 结构体                         | struct                                              |                                                    |
| 枚举类型                       | enum                                                |                                                    |
| 引用类型                       | 任何类型                                            | object                                             |
| 字符串类型                     | string                                              |                                                    |
| 指针类型                       | int number = 42; int* pointer =  &number;（unsafe） |                                                    |
| 类                             | class C {…}                                         |                                                    |
| 动态类型                       | dynamic                                             |                                                    |
| 数组                           | int[]、string[]                                     |                                                    |
| 集合                           | 动态数组：List<>                                    |                                                    |
| 映射：Dictionary\<TKey,TValue> |                                                     |                                                    |
| 无重复元素的集合：HashSet\<T>  |                                                     |                                                    |
| 先进先出队列Queue\<T>          |                                                     |                                                    |
| 后进先出栈：Stack\<T>          |                                                     |                                                    |
| 双向链表：LinkedList\<T>       |                                                     |                                                    |
| 接口                           | Interface                                           |                                                    |
| 委托                           | Delegate                                            |                                                    |
| 事件                           | Event                                               |                                                    |
| 匿名类型                       |                                                     |                                                    |
| 泛型                           | 允许我们使用占位符来定义类和方法                    |                                                    |
| Lambda表达式                   |                                                     |                                                    |

## 2.2  C#的类型安全

通过类型检查和限制确保变量和操作在使用时始终与其声明的类型一致，以减少类型相关的错误并提高代码的可靠性。

## 2.3  结构体类型

结构体是一种自定义的数据结构，适用于轻量级的数据。

## 2.4  枚举类型

枚举类型是一个包含一组命名常量的新数据类型，通常表示为从 0 开始的递增数字。

## 2.5  object

object 和 System.Object 是同一概念的不同表示方式。System.Object 是 C# 中的根基类，所有类（都派生自 System.Object。object则是关键字，使用 object 作为返回值类型或方法参数类型，意味着这个方法可以接受或返回任何类型的对象。一定程度上实现通用性和灵活性，但也可能导致类型安全性和性能方面的问题。因此，建议在能够确定具体类型的情况下，尽量使用具体的类型。只有在确实需要处理多种类型的情况下，才考虑使用 object 类型。

## 2.6  动态类型

动态类型是指变量的类型在运行时才确定，可以提高灵活性，但是会降低类型安全性。

 

## 2.7  数组

固定大小且元素类型一样。

| int[] numbers = new int[5] { 1, 2, 3, 4, 5 };                                                                                                     |
| ------------------------------------------------------------------------------------------------------------------------------------------------- |
| int[,] matrix = new int[3, 3] {                   { 1, 2, 3 },                   { 4, 5, 6 },                   { 7, 8, 9 }                    }; |

## 2.8  动态数组

动态数组可以动态增加和减少元素个数，可以存储任意类型的元素，包括数组类型。ArrayList类提供了多种属性和方法来操作动态数组，例如添加、删除、插入、排序、搜索、复制等。

## 2.9  映射

 

 

## static

变量实例化后才可以被访问；但是加了static之后不实例化也可以访问

 

## 接口

### 接口的特点：

* 接口一般以“I”开头命名，可以声明方法、属性等成员，这些成员只能是“抽象”的，所以不能直接对接口进行实例化，这些成员的访问权限默认为 public。
* 接口一旦被实现（被一个类继承）， 类就必须实现接口中的所有成员，除非派生类本身也是抽象类。
* C# 支持接口之间的多重继承，实现子接口时必须实现所有父接口的成员。
* 接口类型可以引用任何实现该接口的对象，无论是直接实现该接口的类还是实现了该接口的子类都可以。





## 泛型

泛型可以用来定义通用类或方法，它允许我们使用占位符来定义类和方法，使用时将这些占位符替换为指定的类型

泛型的特性：

泛型最常见的用途是创建集合类

可以最大限度地重用代码、保护类型的安全性以及提高性能

泛型数据类型中所用类型的信息可在运行时通过使用反射来获取

System.Collections.Generic 命名空间中的泛型集合类可以用来代替System.Collections

泛型委托：

可以使用类型参数定义泛型委托delegate T NumberChanger<T>(T n);



泛型的使用：

* 泛型类（Generic Class）：使用一个或多个类型参数来定义的类。例如，List<T> 是一个泛型类，其中的 T 可以被替换为任何有效的类型。
* 泛型结构（Generic Structure）：与泛型类类似，但是定义为结构体。
* 泛型接口（Generic Interface）：使用一个或多个类型参数来定义的接口。它可以在多个类中实现，并且可以根据需要指定具体的类型。
* 泛型方法（Generic Method）：在方法签名中使用类型参数的方法。它可以在静态类、非泛型类中声明，并且可以根据需要指定具体的类型。

使用泛型的代码例子

```c#
csharp
public class MyGenericClass<T>
{
    private T genericMember;

    public MyGenericClass(T value)
    {
        genericMember = value;
    }

    public T GenericMethod(T parameter)
    {
        // 在这里可以使用泛型类型参数 T 进行操作
        return genericMember;
    }
}
```



 

 

# 封装

访问权限，可以用来修饰类中的变量和方法等成员（可以限制其它成员的访问）

public：公共的，所有对象都可以访问，但是需要引用命名空间；

private：私有的，类的内部才可以访问；

internal：内部的，同一个命名空间的对象可以访问；

protected：受保护的，类的内部或类的父类和子类中可以访问；

Protected internal：protected 和 internal 的并集，符合任意一条都可以访问。

 

 

# 重构

代码重构的

代码抽离可实现代码重用、模块化和可维护性。



面向接口编程



# ASP.NET Core Web API

> [ASP.NET](http://asp.net/) Core Web API 是基于 [ASP.NET](http://asp.net/) Core 框架的工具集，<u>用于构建和发布 HTTP 服务</u>，以便与客户端应用程序进行通信。它主要用于构建支持 RESTful 架构风格的 Web 服务，可以用于创建各种类型的 Web 应用程序，包括单页应用程序（SPA）、移动应用程序后端等。

 

 