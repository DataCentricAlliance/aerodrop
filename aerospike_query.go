package main

import (
    "github.com/aerospike/aerospike-client-go"
    "time"
)

type AeroQueryCond struct {
    value   string
    between []string
}

type AeroQuery struct {
    namespace string
    set       string
    queries   map[string]AeroQueryCond
}

func (storage *AerospikeStorage) Query(query AeroQuery) *[]*AeroResponse {
    var (
        record    *aerospike.Record
        stm       *aerospike.Statement
        err       error
        recordset *aerospike.Recordset
        Bins      []*AeroResponse
        policy    *aerospike.QueryPolicy
    )

    policy = aerospike.NewQueryPolicy()
    policy.Timeout = time.Duration(config.Aerospike.ReadTimeout) * time.Millisecond

    stm = aerospike.NewStatement(query.namespace, query.set)
    for name := range query.queries {
        if len(query.queries[name].value) > 0 {
            stm.Addfilter(aerospike.NewEqualFilter(name, query.queries[name].value))
        } else {
            panic("Not supported")
        }
    }

    if recordset, err = storage.Client.Query(policy, stm); err != nil {
        panic("timeout")
        return nil
    }
    for record = range recordset.Records {
        Bins = append(Bins, RecordToAeroResponse(record))
    }
    return &Bins
}
