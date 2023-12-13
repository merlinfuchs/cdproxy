# CDProxy

A proxy & cache for files on a content delivery network like the Discord CDN

## What?

CDProxy lets you submit URLs of hosted files and returns an URL that you can use to download that file at a later point. If the file is smaller than the set max size it will try to download it and store it locally. When an expiry is set the local file will be deleted once it has expired. If no local file is found, CDProxy will try to download from the original URL on the fly.

## Why?

In September 2023 Discord has announced to introduce authentication to their CDN used for file uploads on the platform. This means that CDN URLs from Discord will expire 24 hours (the timing may change) after it has been obtained from Discord.

Some services heavily rely on storing these URLs and accessing at a later point. CDProxy can help mitigating the affects of this change by caching files for an extended period.

## Configuration (WIP)

CDProxy will for a file called `config.yaml` containing configuration.

```yaml
default_expiry: 0
default_original_expiry: 86400 # 1 day

max_queue_size: 100
max_workers: 8

public_url: http://localhost:8080

brotli_compression_level: 4
```
