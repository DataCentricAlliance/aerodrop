# Build

## Debian/Ubuntu
* apt-get install -qq -yy git mercurial  bison make curl
* wget -q https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer -O /tmp/gvm-insta
ller
* bash /tmp/gvm-installer
* source ~/.gvm/scripts/gvm
* gvm install go1.4
* gvm use go1.4 --global
* go get github.com/mattn/gom
* gom install
* gom build
* cp aerospike_client /usr/local/sbin
* cp config.yaml /etc/aerospike_client.yaml
* edit /etc/aerospike_client.yaml
* run /usr/local/sbin/aerospike_client --config=/etc/aerospike_client.yaml

## Mac OS X (homebrew)
* brew install git mercurial  bison make curl
* wget -q https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer -O /tmp/gvm-insta
ller
* bash /tmp/gvm-installer
* source ~/.gvm/scripts/gvm
* gvm install go1.4
* gvm use go1.4 --global
* go get github.com/mattn/gom
* gom install
* gom build
* run ./aerospike_client --config=config.yaml


# HTTP Interface for Aerospike

## Example HTTP Request/Response

### Response explain
    bins - user data
    ttl - is meta argument with record's ttl
    version - record's generation
    pk - primary key of record

### Index

##### Explain
* Url for create new index is  /v1/index/disk/demo/index_7782
    * v1 - version
    * index - action
    * namespace - disk
    * set - demo
    * index name - index_7782


#### Create Index

##### Example Request
```
PUT /v1/index/disk/demo/index_7782 HTTP/1.1
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

### Item

##### Explain

* Url for create new index is  /v1/item/disk/demo/7782
    * v1 - version
    * item - action
    * namespace - disk
    * set - demo
    * primary key - 7782

#### New/Update Item

##### Example Request

```
PUT /v1/item/disk/demo/7782 HTTP/1.1
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

#### Get by Primary Key

##### Example Request

```
GET /v1/item/disk/demo/7782 HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip, deflate
Accept: */*
User-Agent: HTTPie/0.8.0

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:19 GMT
Content-Length: 56

{"bins":{"user_id_7782":"7782"},"version":1,"ttl":3600, "pk": "7782"}
```

#### Batch Request by Primary Key

You can pass multiple Id that comma separated
Server will return list of record or null when no record found

##### Example Request

```
GET /v1/item/disk/demo/5,1,2,3,4,5,5 HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip, deflate
Accept: */*
User-Agent: HTTPie/0.8.0

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:20 GMT
Content-Length: 286

[null, {"bins":{"user_id":"xxxx"},"version":1,"ttl":159189797, "pk": "1"},{"bins":{"bin":"post_id_2","user_id":"user_1"},"version":3,"ttl":156955627, "pk": "2"},{"bins":{"bin":"post_id_3","user_id":"user_1"},"version":1,"ttl":156955620, "pk": "3"},{"bins":{"bin":"post_id_4","user_id":"user_1"},"version":1,"ttl":156955624, "pk": "4"}, null, null]
```



### Query

##### Explain

* Url for create new index is  /v1/query/disk/demo/?user_id_7782=7782
    * v1 - version
    * query - action
    * namespace - disk
    * set - demo
    * primary key - 7782
    * filter statements - user_id_7782=7782, also you can pass multiple statements (example: user_id=user_x&action=new)

#### Example Request

```
GET /v1/query/disk/demo/?user_id_7782=7782 HTTP/1.1
Host: localhost:8080
Accept-Encoding: gzip, deflate
Accept: */*
User-Agent: HTTPie/0.8.0

HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Dec 2014 15:34:20 GMT
Content-Length: 63

[{"bins":{"user_id_7782":"7782"},"version":1,"ttl":156616459, "pk": "7782"}]
```

# Know Problems
* aerospike-client-go doesn't support custom digest for key

