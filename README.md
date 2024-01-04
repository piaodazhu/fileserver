# fileserver

[中文](./README_ZH.md)

## Description
A simple HTTP file server with minimal features.

## Features
1. Extremely simple. No need to install specific clients; any user can directly upload and download using wget/curl.
2. Fast deployment. Single-file static compilation, runs with default configurations.
3. Web browsing. Supports simple web directory browsing and downloading.
4. Progress tracking. Supports real-time progress tracking during file uploads.
5. Insecure. Provides only simple password validation, rate limiting, and concurrency control. It is recommended for use only in trusted intranets.

## Service Configuration
```yaml
# Service listening address
ip: "127.0.0.1"
port: 9988
# Local root directory for static resources
rootPath: "/var/fileserver"
# Password required for uploads
password: "network123"
# Document displayed on the homepage
docFile: "README_EN.md"
# Maximum file upload size limit (MB)
maxFileSize: 4096
# Total size limit for root directory storage (MB)
maxStorageSize: 20480
# Maximum concurrency limit (requests/s)
maxConcurrency: 5
# Maximum service queuing time (s)
maxQueuing: 5
# Maximum requests per second limit (requests/s)
maxLimit: 100
# Maximum burst rate per second (requests/s)
maxBurst: 10
```

## File Upload

```sh
# POST /rawupload -H 'password:network123' --data-binary url
# Upload the local file "local/file.tar.gz" to the "some/path/" directory on the server with the filename "file"
curl -X POST -H 'password: network123' --data-binary @local/file.tar.gz 'http://127.0.0.1:9988/rawupload/some/path/file'
```

## File Upload Progress Query
```sh
# GET /progress url
# Query the upload progress of "some/path/file"
curl -X GET 'http://127.0.0.1:9988/progress/some/path/file'
# Output progress information: Transfer percentage [Elapsed time / Estimated time] [Transferred / Total size] Rate
# 3.03% [9.1 s / 300.4 s] [6240880 B / 206167131 B] 0.65 MB/s
```

## File Download

```sh
# GET /static
# Download the file "file.tar.gz" from the "hello" directory on the server
wget http://127.0.0.1:9988/static/hello/file.tar.gz
```

## File Deletion

```sh
# DELETE /delete
# Delete the file "file.tar.gz" from the "hello" directory on the server
curl -X DELETE http://10.108.30.85:9988/delete/hello/file.tar.gz
```

## Browser Access
[-> Click here to browse <- ](http://127.0.0.1:9988/list)
```
GET /list
Access via browser `http://127.0.0.1:9988/list`
```