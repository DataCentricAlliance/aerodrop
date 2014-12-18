package main

import (
    "fmt"
    "time"
)

type Listener interface {
    Run()
}

type ListenerFunType func() bool
type ListenerMapType map[string]ListenerFunType

var listeners ListenerMapType

func RegisterListener(name string, listener_getter ListenerFunType) {
    listeners[name] = listener_getter
}

func init() {
    listeners = make(ListenerMapType, 1)
}

func run_listener() {
    var count = 0
    var name string
    var fn ListenerFunType

    for name, fn = range listeners {
        if fn() {
            fmt.Printf("RUN %s\n", name)
            count += 1
        } else {
            fmt.Printf("FAIL %s\n", name)
        }
    }

    fmt.Printf("Listeners: %d\n", count)
    if count > 0 {
        for _ = range time.Tick(30 * time.Second) {
        }
    }
}
