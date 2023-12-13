# CDProxy

A proxy & cache for files on a content delivery network like the Discord CDN

## What?

CDProxy lets you submit URLs of hosted files and returns an URL that you can use to download that file at a later point. If the file is smaller than the set max size CDProxy will try to download it and store it on an SFTP server. When an expiry is set the stored file will be deleted once it has expired. If no stored file is found, CDProxy will try to download from the original URL on the fly.

## Why?

In September 2023 Discord has announced to introduce authentication to their CDN used for file uploads on the platform. This means that CDN URLs from Discord will expire 24 hours (the timing may change) after it has been obtained from Discord.

Some services heavily rely on storing these URLs and accessing at a later point. CDProxy can help mitigating the affects of this change by caching files for an extended period.

## Installation

```shell
go install github.com/merlinfuchs/cdproxy@latest
```

## Configuration

CDProxy will look for a file called `config.yaml` containing configuration.

```yaml
host: 127.0.0.1
port: 8080
public_url: http://localhost:8080 # Where CDProxy is exposed, used to generate the URL for the file

db_file_Name: cdproxy.db
download_timeout: 30 # timeout in seconds for downloading from the original url

sftp_host: localhost:22
sftp_user: ""
sftp_password: ""
brotli_compression_level: 7

default_max_size: 104857600 # 100MB, in bytes
default_expiry: 0 # in seconds, no expiry by default
default_original_expiry: 86400 # when the original url will become invalid, default 24 hours

max_queue_size: 100 # Number of files that can wait to be processed in the queue
max_workers: 8 # Number of files that can be processed at once, default number of CPU cores
```

## Usage

### 1. Start the server

```shell
cdproxy
```

### 2. Submit a file

```shell
POST /submit
```

```json
{
  "original_url": "https://cdn.discordapp.com/...",
  "original_expires_at": null, // ISO timestamp when the original url will expire (optiona, defaults to config value)
  "expires_at": null, // ISO timestamp when the file expires (optional, defaults to config value)
  "size": 42, // If you already know the size of the file you can set it here, this will prevent CDProxy from having to download it at all if it's too big (optional)
  "max_size": 1000, // If the size of the file is lower than this it will be stored (optional, defaults to config value)
  "metadata": { "user_id": "123" }, // Any metadata for the file (optional)
  "wait": false // Whether to wait for the file to be processed or return instantly (optional)
}
```

### 3. Download file

```shell
GET /download/<file_id>
```

### 4. Profit
