# GoFile

一个匿名分享文件的工具，采用B/S架构，用户上传文件后会生成一个URL，通过这个URL可以下载文件，文件会在过期后自动从服务器物理删除，减少服务器的储存压力。

当然肯定是不支持IE的  :no_entry_sign:

# 开发

## 服务端


用下面的命令解决依赖。

```sh
go mod tidy
```

然后就是编译。

```sh
go build -o gofile main.go
```

## web端

进入到web目录后运行下面命令解决依赖。

```sh
npm i
```

然后打包。

```sh
npm run build
```

会在web目录下的`dist`生成打包好的web页面。


## 打包脚本

如果go和npm的依赖都已经解决了，可以直接运行根路径下的`make.sh`脚本，这样会在`dist/GoFile-${version}-{$kernel}-${platform}`路径下生成对应平台的程序，并且打包好`tag.xz`压缩包。



# 启动

非常简单，无论是手动编译还是用的编译脚本，只需要运行`gofile`二进制文件即可，默认会监听本地的`1323`端口，直接访问`http://localhost:1323`即可看到效果。

> 更多参数请`gofile -h`查看。


