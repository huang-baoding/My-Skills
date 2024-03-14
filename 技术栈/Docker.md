# Docker
## 1  Docker面试题
### 1.1  什么是Docker

Docker 可以将应用程序与其依赖环境打包成一个容器，这个容器在任何Docker平台上都可以运行，并且一个Docker可以运行多个独立的程序，解决了开发、运维环境不一致的问题。总的来说，Docker 是一个开源的容器化平台，在软件开发、部署和运行方面具有重要的作用。

### 1.2  Docker和传统虚拟机的区别

|              | 虚拟机                                               | Docker                                                       |
| ------------ | ---------------------------------------------------- | ------------------------------------------------------------ |
| **实现原理** | 虚拟出一套硬件，然后在这个基础上虚拟出完整的操作系统 | 直接虚拟出一个操作系统运行在宿主机的操作系统上，比虚拟机少了一个抽象层；这个虚拟的操作系统相当于一个进程，比虚拟机更加的轻量所以启动速度更快，更容易管理对CPU资源的利用率更高 |
| **部署**     | 部署过程较为繁琐                                     | Docker可以将应用程序和依赖项打包成一个镜像，并在任何支持Docker的环境中部署，部署过程更加灵活。 |
| **隔离性**   | 每个虚拟机拥有独立的内核及用户空间，隔离性较高       | Docker的容器则是共享宿主机的内核，通过命名空间实现用户空间的隔离，通过控制组限制访问的资源等实现隔离，隔离性较低。 |

### 1.3  Docker三个基本概念

* 镜像

  镜像是一个只读的静态文件，包含了程序运行所需的一切资源，包括操作系统、程序和环境等。

* 容器

  容器是基于镜像创建的运行实例，它是一个隔离的进程空间。

* 仓库

  仓库是用于存储和组织镜像的地方。

> 镜像和容器在系统中的存储分为镜像层和容器层，镜像层是只读的，里面又包含多个文件系统层和唯一标识符，而容器层是可写的，它记录了容器自身的变化。当一个镜像启动了多个容器，它们的镜像层都是共用的，容器层则不同。

### 1.4  Docker的架构

Docker基于C/S架构,主要组成是客户端、服务端守护进程、镜像、容器、仓库。客户端用于和服务端交互，服务端则执行客户端的交互请求对容器进行管理，镜像是可运行程序的模板，可以存放于仓库，容器则是镜像的运行实例，一个镜像可以运行多个容器。

> 更详细一些：服务端包括Docker Server和Engine，Docker Server提供接口，路由，Controller等服务，而Engine才真是对容器进行管理。

### 1.5  Docker的持久化方式

Docker的持久化方式主要是容器卷和绑定挂载：

* 容器卷：创建容器卷时，Docker 会在主机上的指定位置分配空间。这个容器卷由Docker进行管理可以给多个容器共享，即使容器被删除，容器卷的数据依然存在。
* 绑定挂载：将主机上的任何目录挂在到容器内的目录，主机和容器可以双向修改数据。

当有敏感信息或者不需要存储的缓存信息时也会用到临时文件挂载（tmpfs挂载）：

* 临时文件挂载可以在容器运行时快速且安全地在内存中创建临时的文件系统，但需要在容器终止前数据持久化，否则数据会丢失。

  >需要在docker run命令中添加**--tmpfs**标签

还可以用docker社区的第三方工具进行分布式文件系统的持久化等等。

### 1.6  什么是虚悬镜像

仓库和标签都是\<none>的镜像叫做虚悬镜像。

> 查询所有虚悬镜像：docker image ls -f dangling=true
>
> 删除所有虚悬镜像：docker image prune

> 在构建或者删除镜像时发生错误的话可能会产生虚悬镜像。

### 1.7  创建镜像的方式有哪些

* **基于现有容器创建镜像**：对一个运行中的容器进行了修改后，可以commit成一个新的镜像。

- **使用 Dockerfile 创建镜像**：编写Dockerfile 文件，然后使用 `docker build` 命令 可以创建新镜像。

* **使用第三方工具创建新镜像**。


### 1.8  进入和退出容器的方式

* 退出：

exit 容器停止
ctrl+p+q 容器不停止

* 两种重新进入正在运行的容器的方式：

docker exec -it ID（exit不会停止容器）
docker attach -it ID（exit停止容器）

>  有的Docker容器后台运行，就必须有一个前台进程,所以一般用 -d后台启动的程序，再用exec进入对应容器实例


###   1.9  什么是harbor
Harbor是企业级的镜像仓库，提供了UI、权限控制等功能，并支持RESTful API。

> Harbor实质 上是对docker registry做 了封装，扩展了自己的业务模板。



## 2  Docker基础
### 2.1  安装Docker

配置环境 yum -y install gcc
  yum -y install gcc-c++
安装需要的软件包：yum install -y yum-utils
设置镜像仓（国内阿里云仓库）
yum-cnfig-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
更新yum软件包索引 yum makecache fast
安装docker ce （引擎）：yum -y install docker-ce-cli containerd.io
查看是否安装成功：docker version
启动docker：systemctl start docker
配置阿里云镜像加速器（不然之后下载配置会超时）

### 2.2  运行容器

` docker run [options] imageID [command][ARG]`
options说明：
--name 容器新名字
-d后台运行并返回容器ID，也即启动守护式容器（后台运行）
-i以交互模式运行容器，通常与-t同时使用
-t为容器重新分配一个伪输入终端，通常与-i同时使用也即启动交互式容器（前台有伪终端，等待交互）
-P随机端口映射
-p指定端口映射
docker run -it ubuntu /bin/bash   -it是启动交互式终端，/bin/bash是交互式命令接口

> docker run 的时候会查看本地有没有镜像没有的话先从镜像仓库上拉下来再run一个容器；没有写标签默认是latest

### 2.3  查看正在运行的容器

`docker ps [options]`

-a :列出当前所有正在运行的容器+历史上运行过的
-l :显示最近创建的容器。
-n：显示最近n个创建的容器。
-q :静默模式，只显示容器编号。

### 2.4  拷贝文件和导入导出镜像

拷贝：
docker cp 容器ID:容器内路径 目的主机路径
导出容器export：导出整个容器的内容留作为一个tar归档文件
docker export 容器ID > 文件名.tar //导出的tar文件放在当前目录（在目的主机执行导出命令）
导入import：从tar包中的内容创建一个新的文件系统再导入为镜像
cat 文件名.tar | docker import - （镜像用户/)镜像名：版本号

### 2.5本地镜像发布到阿里云：

aliyun.com->控制台->个人实例->创建命名空间->创建本地仓库->创建后显示如何使用本仓库的命令包括（1）登录阿里云（2）拉取本仓库的镜像（3）上传镜像到本仓库
上传步骤： （1）登录阿里云：
docker login --username=aliyun1214553053 registry.cn-hangzhou.aliyuncs.com
（2）推送到阿里云的仓库
$ docker login --username=aliyun1214553053 registry.cn-hangzhou.aliyuncs.com
$ docker tag [ImageId] registry.cn-hangzhou.aliyuncs.com/baoding2/shangchuanubuntu:[镜像版本号]
$ docker push registry.cn-hangzhou.aliyuncs.com/baoding2/shangchuanubuntu:[镜像版本号]
上传后在此仓库上把镜像拉下来：
$ docker pull registry.cn-hangzhou.aliyuncs.com/baoding2/shangchuanubuntu:[镜像版本号]

### 2.6  私人仓库

1. Dockerhub、阿里云这样的公共镜像仓库可能不太方便，涉及机密的公司不可能提供镜像给公网，所以需要创建一个本地私人仓库供给团队使用，基于公司内部项目构建镜像。
  （Docker Registry是官方提供的工具镜像，可用于构建本地私有镜像仓库）
  （1）拉取registry镜像：
  docker pull registry
  （2）运行registry镜像容器：
  docker run -d -p 5000:5000 -v /home/baoding/myregistry/:/tem/registry --privileged=true registry
  默认情况下，仓库被创建在容器的/var/lib/reegistry目录下，建议自行用容器卷，方便与宿主机联调
  （3）查询本地仓库有无镜像：
  curl -XGET http://192.168.111.100:5000/v2/_catalog (选择ens33的inet)
  （4）把镜像打包成符合规范的镜像：
  docker tag 16ecd2772934 192.168.111.100:5000/cangkuredis:1.2
  （5）修改本地配置文件使之支持http：
  vim /etc/docker/daemon.json    添加： ,”insecure-registres”:[192.168.111.100:5000”]
  (修改后不生效可以重启docker）
  （6）把本地镜像push到本地仓库：
  docker push 192.168.111.100:5000/cangkuredis:1.2 
  （7）从本地仓库上pull镜像：
  docker pull 192.168.111.100:5000/cangkuredis:1.2
2. ​


### 2.7  容器卷

Docker挂载主机目录访问cannot opent directory:Pemission denied
解决办法：在挂载目录后多加一个--privileged=true参数即可（没出错也尽量加上）
docker run -it --priviligede=true -v /宿主机绝对路径目录：/容器内目录 镜像名

- 特点： 1.数据卷可在容器之间共享或重用数据

2.卷中的更改可以直接实时生效
（主机和容器的容器卷目录会同时变化，容器停的时候依然会同步修改）
3.数据卷中的更改不会包含在镜像的更新中
4.数据卷的生命周期一直持续到没有容器使用它为止
查看是否挂载成功：docker inspect 容器ID    (在Mounts里查看）

- 启动容器时附带容器卷：

docker run -it --name=u5 --privileged=true -v /tmp/dockertmp1:/tmp/docker2tmp ubuntu
docker run -it --name=u5 --privileged=true -v /tmp/dockertmp1:/tmp/docker2tmp:ro ubuntu
ro：read only 在容器内只能读 在主机可以读写

- 容器2继承容器1的容器卷  

docker run -it --privileged=true --volumes-from 父类 --name u2 ubuntu
如：docker run -it --privilegeed=true --volumes-from u5 --name u6 ubuntu
主机和u5容器卷的内容u6也有

### 2.8  Docker上安装应用整体步骤

（1）seach 镜像名
（2）pull 镜像名：tag
（3）docker images 查看镜像
（4）docker run -it（-d） -p......   启动镜像
（5）停止容器
（6）移除容器

> tomcat运行不起来的解决办法：把容器里的目录/usr/local/tomcat下的webapps文件删除；并把webapps.dist重命名为webapps

> 启动mysql的命令可在dockerhun官网中查看：https://hub.docker.com/_/mysql找到

在docker中运行的mysql容器可在外部用连接软件进行连接
插入中文的时候会报错，需要修改第37集第10min
docker安装mysql并run
出容器后，建议请先修改完字符集编码后再新建mysql数据库-表-插数据
运行的时候一定要挂容器卷

## 3  DockerFile

> Dockerfile是用来构建Docker镜像的文本文件，是由一条条构建镜像所需的指令和参数构成 的脚本。DockerFile每条指令都会创建一个新的镜像层并对镜像进行提交

（1）FROM：基础镜像，当前新镜像是基于哪个镜像的，指定一个已经存在的镜像作为模板，第一条必须是from
（2）MAINTAINER：镜像维护者的姓名和邮箱地址
（3）RUN：容器构建时需要的命令（shell或exec）
（4）EXPOSE：当前容器对外暴露出的端口
（5）WORKDIR：指定在创建容器后，终端默认登陆的进来工作目录，一个落脚点
（6）USER：指定该镜像以什么样的用户去执行，如果都不指定，默认是root
（7）ENV：用来在构建镜像过程中设置环境变量
（8）ADD：将宿主机目录下的文件拷贝进镜像且会自动处理URL和解压tar压缩包
（9）COPY：类似ADD，拷贝文件和目录到镜像中。将从构建上下文目录中 <源路径> 的文件/目录复制到新的一层的镜像内的 <目标路径> 位置。COPY src dest 、 COPY ["src", "dest"] 、 \<src源路径>：源文件或者源目录 、 \<dest目标路径>：容器内的指定路径，该路径不用事先建好，路径不存在的话，会自动创建。
（10）VOLUME：容器数据卷，用于数据保存和持久化工作
（11）CMD：指定容器启动后的要干的事情。Dockerfile 中可以有多个 CMD 指令，但只有最后一个生效，CMD 会被 docker run 之后的参数替换
（12）ENTRYPOINT：类似于 CMD 指令，但是ENTRYPOINT不会被docker run后面的命令覆盖，而且这些命令行参数会被当作参数送给 ENTRYPOINT 指令指定的程序（如果 Dockerfile 中如果存在多个 ENTRYPOINT 指令，仅最后一个生效。）

```dockerfile
# 基础镜像
FROM mysql:latest

# 维护者信息
MAINTAINER Your Name <your.email@example.com>

# 设置环境变量
ENV MYSQL_ROOT_PASSWORD=rootpassword
ENV MYSQL_DATABASE=exampledb
ENV MYSQL_USER=exampleuser
ENV MYSQL_PASSWORD=examplepassword

# 在容器中创建工作目录
WORKDIR /app

# 添加初始化数据库的 SQL 脚本到工作目录
ADD init.sql /app/

# 暴露 MySQL 默认端口
EXPOSE 3306

# 指定容器启动时执行的命令
CMD ["mysqld"]

# 指定容器启动时执行的入口点
ENTRYPOINT ["docker-entrypoint.sh"]
```



* 写好Dockerfile后需执行docker build -t centosjava8:1.5 . 形成一个新镜像（.是当前目录，在其之前有一个空格）
* 微服务打包成jar包->在同一目录下创建Dockerfile->docker build形成镜像

## 4  Docker Compose

> Docker Compose负责对Docker容器集群的快速编排

37. Docker Compose可以管理多个Docker容器组成一个应用。配置docker-compose.yml文件，写好多个容器之间的调用关系。然后一个命令就能同时启动或关闭这些容器。（如一个项目需要用到多个软件）

38. Docker Compose使用步骤：
    (1)编写Dockerfile定义各个微服务应用并构建镜像
    (2)配置docker-compose定义完整的业务单元      //vim 
    (3)-compose.yml
    (4)执行docker-compose up命令，完成一件部署上线（等价于一次运行多个docker run）
    39.docker-compose.yml例子

    ```yaml
    version: '3'   
    services:
      mysql:
        image: mysql:5.7
        command: --default-authentication-plugin=mysql_native_password
        container_name: mysql
        hostname: mysqlServiceHost
        network_mode: bridge
        ports:
        - "3306:3306"
        #restart: always
        restart: on-failure
        volumes:
        - ./mysql:/var/lib/mysql
        - ./my.cnf:/etc/mysql/conf.d/my.cnf
        - ./mysql/init:/docker-entrypoint-initdb.d/
        - ./shop.sql:/docker-entrypoint-initdb.d/shop.sql
        environment:
        - MYSQL_ROOT_PASSWORD=a123456
        - MYSQL_USER=root
        - MYSQL_PASSWORD=a123456
        - MYSQL_DATABASE=shop
      redis:
        image: redis:3
        container_name: redis
        hostname: redisServiceHost
        network_mode: bridge
        restart: on-failure
        ports:
        - "6379:6379"

      golang:
        build: .
        restart: on-failure
        network_mode: bridge
        ports:
        - "8080:8080"
        links:
        - mysql
        - redis
        volumes:
        - /Users/mac/go/src/gitee.com/shirdonl/LeastMall:/go
        tty: true
    ```
    ​



## 5  Docker命令集

Docker命令集

docker commit提交添加功能后的容器使之成为一个新的镜像（镜像-容器-添加内容-容器2.0-commit-镜像2.0)[依然还在本地镜像仓库]
docker commit -m=“提交的描述信息” -a=“作者” 容器ID 目标镜像名：标签



docker version

systemctl start docker

systemctl stop docker

systemctl enable docker					//开机启动

systemctl status docker

docker info									//查看Docker概要说明

docker --help/docker 具体命令 --help

docker images （-a所有，-q只显示ID）

docker search (--limit 5) redis					//在阿里云上搜redis镜像

docker pull centos(:TAG)						//在阿里云上拉取镜像

docker system df								//查看docker占用空间

docker rmi （-f）镜像ID

docker rm (-f) 容器ID

docker run [options] imageID [command] [ARG]

docker run -it --name=mycentos centos /bin/bash //不命名自动分配

docker run -d redis:6.0.8						//后台守护式容器，如redis

docker ps	(-a所有,-l最近，-n n个，-a只显示ID)//列出正在运行的容器

exit / ctrl+p+q		//退出容器 exit容器停止，ctrl p q容器不停止

docker start 器ID或容器名					//启动已经停止的容器

docker exec -it 容器ID /bin/bash	//重新进入正在运行的容器，exit不会停止容器

docker logs 	容器ID或容器名			//查看容器日志

docker top 容器ID或容器名				//查看容器内运行的进程

docker inspect 容器ID					//查看容器内部细节

docker attach 容器ID /bin/bash //也是重新进入，但是exit会停止容器

docker restart 容器ID或容器名

docker stop容器ID或容器名									//停止

docker kill 容器ID或容器名									//强制停止

docker cp  容器ID:容器内路径 目的主机路径

docker export 容器ID > 文件名.tar

cat 文件名.tar | docker import - (用户名/)镜像名:tag

14.docker commit -m="提交的描述信息" -a="作者" 容器ID 要创建的目标镜像名:[标签名]

15.docker run -it --privileged=true -v /宿主机绝对路径目录:/容器内目录

16.redis-cli -p 6379

docker exec -it 容器ID /bin/bash

 

 

要退出终端，直接输入 exit:

docker ps [OPTIONS]

   OPTIONS说明OPTIONS说明（常用）：

 

退出容器

  两种退出方式

​    exit

​      run进去容器，exit退出，容器停止

​    ctrl+p+q

​      run进去容器，ctrl+p+q退出，容器不停止

强制停止容器

  docker kill 容器ID或容器名

删除已停止的容器

  docker rm 容器ID

​    一次性删除多个容器实例

​      docker rm -f $(docker ps -a -q)

​      docker ps -a -q | xargs docker rm

从容器内拷贝文件到主机上

  容器→主机

  docker cp  容器ID:容器内路径 目的主机路径

导入和导出容器

  export 导出容器的内容留作为一个tar归档文件[对应import命令]

  import 从tar包中的内容创建一个新的文件系统再导入为镜像[对应export]

  案例

​    docker export 容器ID > 文件名.tar

​    cat 文件名.tar | docker import - 镜像用户/镜像名:镜像版本号

 

docker commit提交容器副本使之成为一个新的镜像docker commit -m="提交的描述信息" -a="作者" 容器ID 要创建的目标镜像名:[标签名]

运行一个带有容器卷存储功能的容器实例

   docker run -it --privileged=true -v /宿主机绝对路径目录:/容器内目录      镜像名

 

 

 

 

 

在centos上卸载旧版Docker

 sudo yum remove docker \

​                  docker-client \

​                  docker-client-latest \

​                  docker-common \

​                  docker-latest \

​                  docker-latest-logrotate \

​                  docker-logrotate \

​                  docker-engine

 

 

Docker Compose常用命令

docker-compose -h			

查看帮助。

docker-compose up			

启动所有docker-compose服务

docker-compose up -d	

启动所有docker-compose服务并后台运行

docker-compose down		

停止并删除容器、网络、卷、镜像

docker-compose exec yml里面的服务id			

进入器实例内部，docker-compose exec docker-compose.yml文件中写的服务id/bin/bash

docker-compose ps		

展示当前docker-compose编排过的运行的所有容器。	

docker-compose top

展示当前docker-compose编排过的容器进程

docker-compose logs yml里面的服务id

查看容器输出日志。

docker-compose config

检查配置。

docker-compose config -q

检查配置，有问题才有输出。

docker-compose restart

docker-compose start

docker-compose stop



