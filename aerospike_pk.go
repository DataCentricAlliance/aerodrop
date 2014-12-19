package main

import (
    "github.com/aerospike/aerospike-client-go"
    "time"
)

type AeroPK struct {
    namespace string
    set       string
    pk        []string
}

func (storage *AerospikeStorage) BatchGet(query AeroPK) *[]*AeroResponse {
    var (
        records []*aerospike.Record
        record  *aerospike.Record
        err     error
        Key     *aerospike.Key
        Keys    []*aerospike.Key
        Bins    []*AeroResponse
    )

    policy := aerospike.NewPolicy()
    policy.Timeout = time.Duration(config.Aerospike.ReadTimeout) * time.Millisecond

    for idx := range query.pk {
        Key, _ = aerospike.NewKey(query.namespace, query.set, query.pk[idx])
        Keys = append(Keys, Key)
    }

    if records, err = storage.Client.BatchGet(policy, Keys); err != nil {
        panic("timeout")
        return nil
    }

    if len(records) == 0 {
        return nil
    }
    for idx := range records {
        record = records[idx]
        if record == nil {
            Bins = append(Bins, nil)
        } else {
            Bins = append(Bins, RecordToAeroResponse(record))
        }
    }
    return &Bins
}

func (storage *AerospikeStorage) Get(query AeroPK) *AeroResponse {
    var (
        record *aerospike.Record
        err    error
        Key    *aerospike.Key
    )
    policy := aerospike.NewPolicy()
    policy.Timeout = time.Duration(config.Aerospike.ReadTimeout) * time.Millisecond
    Key, _ = aerospike.NewKey(query.namespace, query.set, query.pk[0])

    if record, err = storage.Client.Get(policy, Key); err != nil {
        panic("timeout")
        return nil
    }
    if record == nil {
        return nil
    }
    return RecordToAeroResponse(record)
}
