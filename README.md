# CDProxy

A proxy & cache for files on a content delivery network like the Discord CDN

## What?

CDProxy lets you submit URLs of hosted files and returns an URL that you can use to download that file at a later point. If the file is smaller than the set max size it will try to download it and store it locally. When an expiry is set the local file will be deleted once it has expired. If no local file is found, CDProxy will try to download from the original URL on the fly.

## Why?

In September 2023 Discord has announced to introduce authentication to their CDN used for file uploads on the platform. This means that CDN URLs from Discord will expire 24 hours (the timing may change) after it has been obtained from Discord.

Some services heavily rely on storing these URLs and accessing at a later point. CDProxy can help mitigating the affects of this change by caching files for an extended period.

## Installation

```shell
go install github.com/merlinfuchs/cdproxy@latest
```

## Configuration

CDProxy will for a file called `config.yaml` containing configuration.

```yaml
host: 127.0.0.1
port: 8080
public_url: http://localhost:8080 # Where CDProxy is exposed, used to generate the URL for the file

db_file_Name: cdproxy.db
file_path: ./files
download_timeout: 30 # timeout in seconds for downloading from the original url

default_max_size: 104857600 # 100MB, in bytes
default_expiry: 0 # in seconds, no expiry by default
default_original_expiry: 86400 # when the original url will become invalid, default 24 hours

max_queue_size: 100 # Number of files that can wait to be processed in the queue
max_workers: 8 # Number of files that can be processed at once, default number of CPU cores
brotli_compression_level: 7
```
