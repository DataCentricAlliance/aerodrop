package main

import (
	"fmt"
	aerospike "github.com/aerospike/aerospike-client-go"
	"sync"
	"time"
)

type AerospikeStorage struct {
	Client    *aerospike.Client
	Namespace string
}

type AeroResponse struct {
	Bins       *aerospike.BinMap `json:"bins"`
	Generation uint32               `json:"version"`
	Ttl        uint32               `json:"ttl"`
	PrimaryKey string            `json:"pk"`
}

var aerospike_connection_lock sync.Mutex
var aerospike_storage *AerospikeStorage

func RecordToAeroResponse(record *aerospike.Record) *AeroResponse {
	var response AeroResponse

	response = AeroResponse{
		Bins:       &record.Bins,
		Generation: record.Generation,
		Ttl:        record.Expiration}

	if record.Key.Value() != nil {
		response.PrimaryKey = record.Key.Value().String()
	}

	return &response
}

func InitAerospikeClient() Storage {
	var Hosts []*aerospike.Host
	var storage *AerospikeStorage
	var err error
	defer aerospike_connection_lock.Unlock()
	aerospike_connection_lock.Lock()

	if aerospike_storage == nil {

		storage = new(AerospikeStorage)
		for _, cfg := range config.Aerospike.Hosts {
			Hosts = append(Hosts, aerospike.NewHost(cfg.Host, cfg.Port))
		}

		policy := aerospike.NewClientPolicy()

		if config.Aerospike.ConnectionQueueSize > 0 {
			policy.ConnectionQueueSize = config.Aerospike.ConnectionQueueSize
		}

		if config.Aerospike.ConnectionTimeout > 0 {
			policy.Timeout = time.Duration(config.Aerospike.ConnectionTimeout) * time.Millisecond
		}

		if storage.Client, err = aerospike.NewClientWithPolicyAndHost(policy, Hosts...); err != nil {
			panic(fmt.Sprintf("Unable to connect to aerospike (%s)\n", err))
		}
		aerospike_storage = storage
	}
	return aerospike_storage
}
