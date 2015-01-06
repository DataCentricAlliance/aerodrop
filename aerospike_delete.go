package main

import (
	"github.com/aerospike/aerospike-client-go"
)

type AeroDelete struct {
	namespace string
	set       string
	pk        string
}

func (storage *AerospikeStorage) Delete(query AeroDelete) bool {
	var (
		err     error
		existed bool
		Key     *aerospike.Key
	)
	Key, _ = aerospike.NewKey(query.namespace, query.set, query.pk)
	if existed, err = storage.Client.Delete(nil, Key); err != nil || !existed {
		return false
	}
	return true
}
