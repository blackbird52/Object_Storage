# 分布式对象存储

## 组件

- go version：1.15
- 消息队列：RabbitMQ
- 元数据服务：ElasticSearch
- 数据冗余：RS纠删码
- 断点续传

## REST 接口

```
PUT /objects/<object_name>
请求头部：
  - Digest: SHA-256=<Base64 code>
  - Content-Length: <length>
请求正文：
	- 对象的内容
```

```
GET /objects/<object_name>?version=<version_id>
响应正文：
	- 对象的数据：这个参数可以告诉接口服务客户端需要的是该对象的第一个版本，默认是最新的那个
```

```
GET /objects/<object_name>
请求头部：
	- Range: bytes=<first>-
响应头部：
	- Content-Range: bytes <first>-<size>/<size>
响应正文：
	- 从 first 开始的对象内容
```

```
GET /versions/
响应正文：
	- 所有对象的所有版本
	
GET /versions/<object_name>
响应正文：
	- 指定对象的所有版本
```

```
HEAD /temp/<token>
响应头部：
	- Content-Range: <token 当前的上传字节数>
```

```
DELETE /objects/<object_name>
```

## 元数据

```json
{
  "mappings": {
    "objects": {
      "properties": {
        "name": {"type": "string", "index": "not analyzed"},
        "version": {"type": "integer"},
        "size": {"type": "integer"},
        "hash": {"type": "string"}
      }
    }
  }
}
```

## 环境变量

$STORAGE_ROOT：存储根目录

$LISTEN_ADDRESS：监听端口

$ES_SERVER：ElasticSearch 服务器地址