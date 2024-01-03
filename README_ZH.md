# fileserver
## 说明
极度简单的HTTP文件服务器，不提供丰富的功能。

## 特性
1. 极度简单。无需安装特定客户端，任何用户可直接用wget/curl进行上传下载操作。
2. 快速部署。单文件静态编译，默认配置运行。
3. Web浏览。支持简单的web目录浏览与下载。
4. 进度查询。支持上传文件过程中，实时进度查询。
5. 不安全。仅在上传时进行简单的密码验证，建议只在可信的内网使用。


## 服务配置
```yaml
# 服务监听地址
ip: "127.0.0.1"
port: 9988

# 静态资源本地根目录
rootPath: "/var/fileserver"

# 上传时需要提供的密码
password: "network123"

# 在主页展示的指引文档路径
docFile: "README_ZH.md"
```

## 文件上传

**新版文件上传接口**
```sh
# POST /rawupload -H 'password:network123' --data-binary url
# 把本地文件"local/file.tar.gz"上传到服务器"some/path/"目录下，文件名为file
curl -X POST -H 'password: network123' --data-binary @local/file.tar.gz 'http://127.0.0.1:9988/rawupload/some/path/file'
```

旧版文件上传接口
```sh
# POST /upload -F "file=@<localfilepath>" -H 'password:network123' url
# 把本地文件"local/file.tar.gz"上传到服务器的"hello"目录下
curl -X POST -F "file=@local/file.tar.gz" -H 'password:network123' http://127.0.0.1:9988/upload/hello/
```

## 文件上传进度查询
**仅支持**新版文件上传接口的上传进度查询。
```sh
# GET /progress url
# 查询some/path/file的上传进度
curl -X GET 'http://127.0.0.1:9988/progress/some/path/file'
# 输出进度信息: 传输百分比 [ 已耗时 / 预计耗时] [ 已传输 / 总大小] 速率
# 3.03% [9.1 s / 300.4 s] [6240880 B / 206167131 B] 0.65 MB/s
```

## 文件下载

```sh
# GET /static
# 从服务器的"hello"目录下载文件"file.tar.gz"
wget http://127.0.0.1:9988/static/hello/file.tar.gz
```

## 浏览器访问
[-> 点此浏览 <- ](http://127.0.0.1:9988/list)
```
GET /list
浏览器访问`http://127.0.0.1:9988/list`
```
