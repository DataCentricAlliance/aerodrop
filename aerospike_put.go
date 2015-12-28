package main

import (
	aerospike "github.com/aerospike/aerospike-client-go"
	"time"
)

type AeroPut struct {
	namespace string
	set       string
	pk        string
	data      struct {
		Bins map[string]interface{} `json:"bins"`
		Meta struct {
			Ttl        uint32 `json:"ttl"`
			Generation uint32 `json:"version"`
		} `json:"meta"`
	}
}

func (storage *AerospikeStorage) Put(query AeroPut) bool {
	var (
		err    error
		Key    *aerospike.Key
		policy *aerospike.WritePolicy
	)
	policy = aerospike.NewWritePolicy(query.data.Meta.Generation, query.data.Meta.Ttl)
	policy.Timeout = time.Duration(config.Aerospike.WriteTimeout) * time.Millisecond
	policy.SendKey = true
	Key, _ = aerospike.NewKey(query.namespace, query.set, query.pk)
	if err = storage.Client.Put(policy, Key, query.data.Bins); err != nil {
		panic(err)
	}
	return true
}
