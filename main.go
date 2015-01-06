package main

import (
	"github.com/op/go-logging"
	"io/ioutil"
	"log"
	"runtime"
	"syscall"
)

var logger = logging.MustGetLogger("logger")
var format = logging.MustStringFormatter("%{time:15:04:05.000000}|%{shortfunc}|%{level:.8s}|%{message}")

func main() {
	// #TODO: wrap standart logging to go-logging with custom level
	log.SetOutput(ioutil.Discard)

	r := syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &r)
	runtime.GOMAXPROCS(runtime.NumCPU())
	ReadConfig()
	InitAerospikeClient()
	run_listener() // last call
}
