package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "runtime/debug"

    // "encoding/json"

    "fmt"
    // "github.com/aerospike/aerospike-client-go/types"
    "net"
    "os"
)

type MemcacheListener struct{}

var Memcache *MemcacheListener

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    var (
        buf        []byte
        err        error
        statements [][]byte
        reader     *bufio.Reader
        action     string
    )

    defer func() {
        if r := recover(); r != nil {
            if r != "EOF" {
                fmt.Fprintf(conn, "SERVER_ERROR %s\r\n", r)
                fmt.Printf("\n\tRecovery from2: %s\n\t\n", r)
                debug.PrintStack()
            }
            conn.Close()
        }
    }()

    reader = bufio.NewReader(conn)

    for {
        if buf, err = reader.ReadBytes('\n'); err != nil {
            panic(err.Error())
        }

        statements = bytes.SplitN(buf, []byte(" "), 2)

        switch string(statements[0]) {
        case "set":
            MemcacheHandlerSet(conn, reader, buf)
            break
        case "get":
            MemcacheHandlerGet(conn, reader, buf)
            break
        case "delete":
            MemcacheHandlerDelete(conn, reader, buf)
            break

        default:
            panic(fmt.Sprintf("Unknown action %s", action))
        }
    }
}

func MemcacheHandlerSet(conn net.Conn, reader *bufio.Reader, first_line []byte) {

    // var buf []byte
    var (
        query      AeroNew = AeroNew{}
        exptime    int
        value_len  int
        readed_len int
        key        []byte
        key_items  [][]byte
        buf        []byte
        flags      int
        err        error
        i          int
    )

    if i, err = fmt.Sscanf(string(first_line), "set %s %d %d %d\r\n", &key, &flags, &exptime, &value_len); err != nil {
        panic(fmt.Sprintf("unable to parse query - %s", err))
    } else if i != 4 {
        panic(fmt.Sprintf("unable to parse (parsed %d, but 5 must)", i))
    }

    buf = make([]byte, value_len)
    if readed_len, err = reader.Read(buf); err != nil {
        panic(fmt.Sprintf("Unable to read value (%s)", err))
    } else if readed_len != value_len {
        panic("Wrong size")
    }
    reader.ReadBytes('\n')

    key_items = bytes.Split(key, []byte("."))
    if len(key_items) < 3 {
        panic(fmt.Sprintf("unable to parse key (parsed %d (%s), but 3 must)", len(key_items), key))
    }
    query.namespace = string(key_items[0])
    query.set = string(key_items[1])
    query.pk = string(bytes.Join(key_items[2:], []byte(".")))
    query.data.Meta.Ttl = int32(exptime)

    if err = json.Unmarshal(buf, &query.data.Bins); err != nil {
        panic(fmt.Sprintf("unable to load value as json %s", err))
    }
    if len(query.data.Bins) == 0 {
        panic("no bins")
    } else if aerospike_storage.Put(query) {
        conn.Write([]byte("STORED\r\n"))
    } else {
        panic("unable to save")
    }
}

func MemcacheHandlerDelete(conn net.Conn, reader *bufio.Reader, first_line []byte) {

    var (
        query     AeroDelete
        key       []byte
        key_items [][]byte
        err       error
        i         int
    )

    if i, err = fmt.Sscanf(string(first_line), "delete %s\r\n", &key); err != nil {
        panic(fmt.Sprintf("unable to parse query - %s", err))
    } else if i != 1 {
        panic(fmt.Sprintf("unable to parse (parsed %d, but 1 must)", i))
    }

    key_items = bytes.Split(key, []byte("."))

    query = AeroDelete{namespace: string(key_items[0]), set: string(key_items[1])}
    query.pk = string(bytes.Join(key_items[2:], []byte(".")))

    if aerospike_storage.Delete(query) {
        conn.Write([]byte("DELETED\r\n"))
        return
    }
    conn.Write([]byte("NOT_FOUND\r\n"))

}

func MemcacheHandlerGet(conn net.Conn, reader *bufio.Reader, first_line []byte) {

    var (
        query     AeroPK
        key       []byte
        key_items [][]byte
        response  *AeroResponse

        buf []byte
        err error
        i   int
    )

    if i, err = fmt.Sscanf(string(first_line), "get %s\r\n", &key); err != nil {
        panic(fmt.Sprintf("unable to parse query - %s", err))
    } else if i != 1 {
        panic(fmt.Sprintf("unable to parse (parsed %d, but 1 must)", i))
    }
    key_items = bytes.Split(key, []byte("."))
    query = AeroPK{}
    query.namespace = string(key_items[0])
    query.set = string(key_items[1])
    query.pk = []string{string(bytes.Join(key_items[2:], []byte(".")))}

    if response = aerospike_storage.Get(query); response == nil {
        conn.Write([]byte("NOT_FOUND\r\n"))
        return
    }

    if buf, err = json.Marshal(response.Bins); err != nil {
        panic("Unable to serialize response")
    }
    fmt.Fprintf(conn, "VALUE %s 0 %d\r\n%s\r\nEND\r\n", key, len(buf), buf)
    // VALUE <key> <flags> <bytes> [<cas unique>]\r\n
    // <data block>\r\n
    // …
    // END\r\n

}

func (listener *MemcacheListener) Run() bool {
    var (
        listen_socket net.Listener
        conn          net.Conn
        err           error
    )

    go func() {
        if listen_socket, err = net.Listen("tcp", ":"+config.Memcache.Port); err != nil {
            fmt.Println("Error listening:", err.Error())
            os.Exit(1)
        }
        defer listen_socket.Close()

        for {
            if conn, err = listen_socket.Accept(); err == nil {
                go handleRequest(conn)
            } else {
                fmt.Println("Error accepting: ", err.Error())
            }
        }
    }()
    return true
}

func RunMemcacheListener() bool {
    if config.Memcache.Port != "" {
        return Memcache.Run()
    } else {
        return false
    }
}

func init() {
    RegisterListener("memcache", RunMemcacheListener)
}
