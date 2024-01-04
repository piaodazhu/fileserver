# Fileserver
[中文](./README_ZH.md)
## Introduction
An extremely simple HTTP file server with minimal features.

## Features
1. Extremely simple. No need to install specific clients; any user can directly upload and download using wget/curl.
2. Quick deployment. Single file, statically compiled, and runs with default configuration.
3. Web browsing. Supports basic web directory browsing and downloading.
4. Progress query. Supports real-time progress query during file upload.
5. Not secure. Simple password verification is performed only during upload. It is recommended for use only in trusted local networks.

## Service Configuration
```yaml
# Service listening address
ip: "127.0.0.1"
port: 9988
# Local root directory for static resources
rootPath: "/var/fileserver"
# Password required for upload
password: "network123"
# Path to the guide document displayed on the homepage
docFile: "README.md"
# Single file size limit (MB)
maxFileSize: 4096
# Total storage space limit (MB)
maxStorageSize: 20480
```

## File Upload

**New version file upload interface**
```sh
# POST /rawupload -H 'password:network123' --data-binary url
# Upload the local file "local/file.tar.gz" to the "some/path/" directory on the server, with the filename as "file"
curl -X POST -H 'password: network123' --data-binary @local/file.tar.gz 'http://127.0.0.1:9988/rawupload/some/path/file'
```

Old version file upload interface
```sh
# POST /upload -F "file=@<localfilepath>" -H 'password:network123' url
# Upload the local file "local/file.tar.gz" to the "hello" directory on the server
curl -X POST -F "file=@local/file.tar.gz" -H 'password:network123' http://127.0.0.1:9988/upload/hello/
```

## File Upload Progress Query
**Only supports** progress query for the new version file upload interface.
```sh
# GET /progress url
# Query the upload progress of some/path/file
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

## File Delete

```sh
# DELETE /delete
# Delete the file "file.tar.gz" from the "hello" directory on the server
curl -X DELETE http://10.108.30.85:9988/delete/hello/file.tar.gz
```

## Browser Access
[-> Click here to browse <- ](http://127.0.0.1:9988/list)
```
GET /list
Access via browser at `http://127.0.0.1:9988/list`
```