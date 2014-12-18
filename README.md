# HTTP Interface for Aerospike

## Example HTTP Request/Response

### Create Index

```
PUT /20141217/index/disk/demo/index_7782 HTTP/1.1
Content-Length: 41
Accept-Encoding: gzip, deflate
Accept: application/json
User-Agent: HTTPie/0.8.0
Host: localhost:8080
Content-Type: application/json; charset=utf-8

{"key": "user_id_7782", "type": "STRING"}

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:19 GMT
Content-Length: 0
```

### New/Update Item

```
PUT /20141217/item/disk/demo/7782 HTTP/1.1
Content-Length: 57
Accept-Encoding: gzip, deflate
Accept: application/json
User-Agent: HTTPie/0.8.0
Host: localhost:8080
Content-Type: application/json; charset=utf-8

{"meta": {"ttl": 3600}, "bins": {"user_id_7782": "7782"}}

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:19 GMT
Content-Length: 0
```

### Get by Primary Key

```
GET /20141217/item/disk/demo/7782 HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip, deflate
Accept: */*
User-Agent: HTTPie/0.8.0

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:19 GMT
Content-Length: 56

{"bins":{"user_id_7782":"7782"},"version":1,"ttl":3600}
```

### Query by filter

```
GET /20141217/query/disk/demo/?user_id_7782=7782 HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip, deflate
Accept: */*
User-Agent: HTTPie/0.8.0

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:20 GMT
Content-Length: 63

[{"bins":{"user_id_7782":"7782"},"version":1,"ttl":156616459}]
```

### Batch Request by Primary Key

```
GET /20141217/item/disk/demo/1,2,3,4 HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip, deflate
Accept: */*
User-Agent: HTTPie/0.8.0

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:20 GMT
Content-Length: 286

[{"bins":{"user_id":"xxxx"},"version":1,"ttl":159189797},{"bins":{"bin":"post_id_2","user_id":"user_1"},"version":3,"ttl":156955627},{"bins":{"bin":"post_id_3","user_id":"user_1"},"version":1,"ttl":156955620},{"bins":{"bin":"post_id_4","user_id":"user_1"},"version":1,"ttl":156955624}]
```

