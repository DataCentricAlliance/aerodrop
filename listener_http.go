package main

import (
    "encoding/json"
    "fmt"
    "github.com/aerospike/aerospike-client-go/types"
    "net/http"
    "net/url"
    "runtime/debug"
    "strings"
)

type HttpListener struct{}

var Http *HttpListener

func HttpHandlerV1New(w http.ResponseWriter, req *http.Request, namespace string, set string, pk string) {
    var query AeroNew = AeroNew{namespace: namespace, set: set, pk: pk}
    var decoder *json.Decoder = json.NewDecoder(req.Body)
    decoder.Decode(&query.data)
    if len(query.data.Bins) == 0 {
        panic("no bins")
    } else if aerospike_storage.Put(query) {
        w.WriteHeader(200)
    } else {
        panic("unable to save")
    }
}

func HttpHandlerV1Index(w http.ResponseWriter, req *http.Request, namespace string, set string, name string) {
    var query AeroIndex = AeroIndex{namespace: namespace, set: set, name: name}
    var decoder *json.Decoder = json.NewDecoder(req.Body)

    decoder.Decode(&query)

    if query.Key == "" || query.Type == "" {
        panic("Need key and type")
    }
    aerospike_storage.CreateIndex(query)
}

func HttpHandlerV1IndexRemove(w http.ResponseWriter, req *http.Request, namespace string, set string, name string) {
    aerospike_storage.DropIndex(AeroIndex{namespace: namespace, set: set, name: name})
}

func HttpHandlerV1Get(w http.ResponseWriter, req *http.Request, namespace string, set string, pk string) {
    var query AeroPK = AeroPK{namespace: namespace, set: set, pk: strings.Split(pk, ",")}
    var response *AeroResponse
    var responses *[]*AeroResponse
    var encoder *json.Encoder = json.NewEncoder(w)

    if len(query.pk) == 1 {
        if response = aerospike_storage.Get(query); response == nil {
            w.WriteHeader(404)
            return
        }
        encoder.Encode(response)
    } else {
        if responses = aerospike_storage.BatchGet(query); responses == nil {
            w.WriteHeader(404)
            return
        }
        encoder.Encode(responses)
    }
}

func HttpHandlerV1Remove(w http.ResponseWriter, req *http.Request, namespace string, set string, pk string) {
    if aerospike_storage.Delete(AeroDelete{namespace: namespace, set: set, pk: pk}) {
        w.WriteHeader(200)
    } else {
        w.WriteHeader(404)
    }
}

func HttpHandlerV1Query(w http.ResponseWriter, req *http.Request, namespace string, set string) {
    var (
        query     AeroQuery = AeroQuery{namespace: namespace, set: set, queries: make(map[string]AeroQueryCond, 0)}
        cond      AeroQueryCond
        repsonses *[]*AeroResponse
        encoder   *json.Encoder = json.NewEncoder(w)
    )

    for query_name, query_value := range req.URL.Query() {
        switch len(query_value) {
        case 2:
            cond = AeroQueryCond{between: query_value}
            break
        case 1:
            cond = AeroQueryCond{value: query_value[0]}
            break
        default:
            panic(fmt.Sprintf("Wrong Key %s", query_name))
        }
        query.queries[query_name] = cond
    }
    repsonses = aerospike_storage.Query(query)
    if *repsonses != nil {
        encoder.Encode(*repsonses)
    } else {
        w.WriteHeader(404)
    }
}

func ParseUrlHTTP(Url string) (
    action string, namespace string, set string, parts []string) {
    u, err := url.Parse(Url)
    if err != nil {
        panic(err)
    }
    parts = strings.Split(u.Path, "/")[2:]
    if len(parts) < 2 {
        panic("Wrong URL")
    }
    action = parts[0]
    namespace = parts[1]
    parts = parts[2:]
    if len(parts) > 0 {
        set = parts[0]
        parts = parts[1:]
    }
    return
}

func HttpHandlerV1(w http.ResponseWriter, req *http.Request) {
    var parts []string
    var namespace string
    var action string
    var set string
    var aerospike_error_code int

    w.Header().Set("Content-Type", "application/json")

    defer func() {
        if r := recover(); r != nil {
            if aerospike_error, ok := r.(types.AerospikeError); ok {
                aerospike_error_code = int(aerospike_error.ResultCode())
                w.Header().Set("X-Aerospike-Error-Code", fmt.Sprintf("%d", aerospike_error_code))
            }
            if aerospike_error_code == 200 {
                w.Header().Set("Content-Type", "text/plain")
                w.WriteHeader(200)
            } else if r == "timeout" {
                w.Header().Set("Content-Type", "text/plain")
                w.WriteHeader(408)
            } else {
                w.Header().Set("X-Error", fmt.Sprintf("%s", r))
                w.Header().Set("Content-Type", "text/plain")
                w.WriteHeader(500)
                fmt.Printf("\n\tRecovery from: %s\n\t\n", r)
                debug.PrintStack()
            }
        }
        req.Body.Close()
    }()

    action, namespace, set, parts = ParseUrlHTTP(req.URL.Path)

    switch action {
    case "index":
        if len(parts) != 1 || parts[0] == "" {
            panic("Need Index name")
        }
        if req.Method == "PUT" {
            HttpHandlerV1Index(w, req, namespace, set, parts[0])
            break
        }
        if req.Method == "DELETE" {
            HttpHandlerV1IndexRemove(w, req, namespace, set, parts[0])
            break
        }
        panic(fmt.Sprintf("Unknown method %s for action %s", req.Method, action))

    case "item":
        if len(parts) != 1 || parts[0] == "" {
            panic("Need PK")
        }
        if req.Method == "GET" {
            HttpHandlerV1Get(w, req, namespace, set, parts[0])
            break
        }
        if req.Method == "PUT" {
            HttpHandlerV1New(w, req, namespace, set, parts[0])
            break
        }
        if req.Method == "DELETE" {
            HttpHandlerV1Remove(w, req, namespace, set, parts[0])
            break
        }
        panic(fmt.Sprintf("Unknown method %s for action %s", req.Method, action))

    case "query":
        HttpHandlerV1Query(w, req, namespace, set)
        break
    default:
        panic(fmt.Sprintf("Unknown action %s", action))
    }
}

func (listener *HttpListener) Run() bool {
    go func() {
        http.HandleFunc("/20141217/", HttpHandlerV1)
        http.HandleFunc("/v1/", HttpHandlerV1)
        http.ListenAndServe(fmt.Sprintf(":%s", config.Http.Port), nil)
    }()
    return true
}

func RunHTTPListener() bool {
    if config.Http.Port != "" {
        return Http.Run()
    } else {
        return false
    }
}

func init() {
    RegisterListener("http", RunHTTPListener)
}
