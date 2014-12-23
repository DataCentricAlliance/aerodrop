# Change history

## 2014.12.23 / 0.0.5 some refactoring
* **Refactoring **
    * Rename files
    * Rename expiration to ttl

## 2014.12.21 / 0.0.4 add memcache listener
* **Features**
    * Support memcache protocol, but only get/set/delete commands

## 2014.12.19 / 0.0.3 add RecordToAeroResponse
* ** Refactoring **

## 2014.12.18 / 0.0.2 pk
* **Features**
	* Просим aerospike хранить Primary Key в оригинальном виде, отдаем его при запросах

## 2014.12.18 / 0.0.1 first
* **Features**
	* реализован предварительная версия HTTP интерфейса
	* JSON для запросов и ответов
	* поддержка Get, BatchGet, Query, Put, CreateIndex и DropIndex
