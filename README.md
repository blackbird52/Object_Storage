# 分布式对象存储

## 组件

- go version：1.15
- 消息队列：RabbitMQ
- 元数据服务：ElasticSearch

## REST 接口

PUT /objects/<object_name>

- Request Header
  - Digest: SHA-256=\<Base64 code\>
  - Content-Length: \<length\>

GET /objects/<object_name>?version=<version_id>

DELETE /objects/<object_name>



GET /locate/<object_name>

GET /versions/

GET /versions/<object_name>

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