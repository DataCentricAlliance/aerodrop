package main

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
)

type AeroIndex struct {
	namespace string
	set       string
	name      string
	Key       string `json: "key"`
	Type      string `json: "type"`
}

func (storage *AerospikeStorage) CreateIndex(query AeroIndex) bool {
	var (
		err        error
		index_type aerospike.IndexType
	)

	switch query.Type {
	case "STRING":
		index_type = aerospike.STRING
		break
	case "NUMERIC":
		index_type = aerospike.NUMERIC
		break
	default:
		panic(fmt.Sprintf("Unknown Index type %s", query.Type))
	}
	if _, err = storage.Client.CreateIndex(nil, query.namespace, query.set, query.name, query.Key, index_type); err != nil {
		panic(err)
	}
	return true
}

func (storage *AerospikeStorage) DropIndex(query AeroIndex) bool {
	var (
		err error
	)

	if err = storage.Client.DropIndex(nil, query.namespace, query.set, query.name); err != nil {
		panic(err)
	}
	return true
}
