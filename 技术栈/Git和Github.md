# Git、Github、GitLab

## 1 mianshiti

### 1.1  什么是Git，和其它版本控制系统（如SVN）有什么区别

Git 是分布式版本控制系统，每个客户端保存的是完整的项目，即使服务器宕机也可以进行开发。而集中化版本（SVN）控制系统有集中管理的服务器，利于项目的管理和授权等，缺点就是，如果服务器宕机那么客户端就无法协同工作，也无法更新提交。

### 1.2   提交时发生冲突，你能解释冲突是如何产生的吗？你是如何解决的？



## 2  Git

### 2.1  常用命令

Git首次安装的时候要先设置用户签名，否则无法提交；同时可以设置自己的邮箱

| 命令                                 | 作用                                                         | 参数   |
| ------------------------------------ | ------------------------------------------------------------ | ------ |
| git config --global user.name        | 设置用户名                                                   |        |
| git config --global user.email       | 设置用户邮件                                                 |        |
| git config (--global) --list         | 查看设置的用户信息，--global表示查看全局，不加表示查看当前仓库的用户信息 |        |
| git init                             | 新建项目的时候要先初始化本地仓库                             |        |
| git clone 远程地址                   | 将远程仓库的内容克隆到本地(自动生成.git文件）                |        |
| git remove -v                        | 查看当前所有远程地址别名                                     |        |
| git remote add 别名 远程地址         |                                                              | 起别名 |
| git status                           | 查看哪些更改还没有暂存                                       |        |
| git add filename或 -A表示添加全部    | 添加到暂存区                                                 |        |
| git commit -m“日志信息”filename或-a  | 提交到本地库                                                 |        |
| git push 别名 分支                   | 推送本地分支上的内容到远程仓库                               |        |
| git pull \|远程库地址别名 远程分支名 |                                                              |        |
| git push –-all origin\|（提交全部）  |                                                              |        |
| git log                              |                                                              |        |
| git reflog                           | 查看历史记录                                                 |        |
| git reset --hard 版本号              | 版本穿梭                                                     |        |
| git branch 分支名                    | 创建分支                                                     |        |
| git checkout 分支名                  | 切换分支                                                     |        |
| git merge 分支名                     | 把指定分支合并到当前分支上                                   |        |
| rm --cached \<file>                  | 丢弃工作目录中对文件的修改，将其还原为最近一次提交的状态     |        |
| git restore \<file>                  | 用于取消暂存区中对文件的修改，将其还原为最近一次提交的状态，并保留工作目录中的修改 |        |
| git restore --staged \<file>         | 用于停止跟踪某个文件，将其从暂存区中移除，但保留在工作目录中 |        |
| git revert                           | 撤销某个提交: 创建一个新的提交，该提交将会抵消先前的提交效果 |        |

//将远程仓库对于分支最新内容拉下来后与当前本地分支直接合并

 ### 2.2 修改后的项目再次提交

* git status
* git add -A
* git commit -m“日志信息” -a
* git push --all origin

## 3  Github

### 3.1  用Git操作Github的文件

GitHub上的文件可以用Git克隆或拉取一个副本到本地，修改本地副本后需要添加到暂存区然后提交到本地仓库最后才能推送到Github

GitHub -> pull或clone -> 修改 -> add -> commit -> push

### 3.1  pull和clone的区别

- pull是获取最新版本代码
- clone是获取完整代码副本（包括历史提交信息）

> pull会把工作区全都更新掉,可以先用fetch拉到本地仓库,使用diff对比工作区，没问题后再合并过来.

### 3.3  提交时忽略文件

不想提交到远程仓库的文件：在.gitignore文件里添加文件的名字和后缀

## 4  GitLab 

## 5  VSCode与Git和Github



### 在vscode上clone一个项目

打开vscode -> 选择新建窗口 -> 打开源代码管理界面”Ctrl + Shift + G“ -> 点击”克隆仓库“ -> 选择从GitHub上克隆 -> vscode会提示登录github账号，登录后可以轻松克隆自己的项目 -> 也可以直接输入仓库的URL来克隆 -> 直接选择克隆到桌面即可





->->->->->->->->->->->->->->->->

在vscode上连接使用github

写好项目后可直接推送到github（会自动初始化）

 

修改之后可能会提交不了也再次拉取不了，可以在当前文件夹使用gitbash来操作代替