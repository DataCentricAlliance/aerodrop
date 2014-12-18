// +build test
package main

import (
    "gopkg.in/fatih/set.v0"
    "encoding/json"
    "github.com/bsm/openrtb"
    // "github.com/stretchr/testify/mock"
    google "./realtime-bidding.pb"
    "code.google.com/p/goprotobuf/proto"
    "encoding/hex"
    // "testing"
)

var facetz_id *FacetzUser
var Segments *set.Set
var bid_request_body []byte
var bid_request *google.BidRequest
var bid_request_wrong_facetz_id *google.BidRequest
var openrtb_response openrtb.Response
var openrtb_response_body []byte
var openrtb_response_body_bad []byte
var openrtb_response_body_adm []byte
var openrtb_response_body_low_price []byte
var default_context *Context
const DEBUG = false


func CreateContext(facetz_id *FacetzUser, bid_request *google.BidRequest) *Context {
    var backend *Backend
    var config *Config
    var context *Context
    backend = DefaultBackend
    config = &Config{
        Backends: []*Backend{backend},
        Micros: 1000000,
    }
    context = &Context{
        system_stat: NewTestStat(),
        facetz_id: facetz_id,
        bid_request: bid_request,
        config: config,
    }
    return context
}


func CreateTestConfigBackend(name string, priority int, urls []*Url) *Backend {
    var backend *Backend

    backend = &Backend{
        Name: name,
        Priority: priority,
        Urls: urls,
        Stat: NewTestStat(),
    }
    return backend
}


func CreateTestConfigRule(name string, priority int, rule_to_backend []*RuleToBackend, segments ...[]int) *Rule {
    var rule Rule
    var Sets []*set.Set

    for _, segment := range segments {
        Sets = append(Sets, set.New(IntArrayToInterArray(segment)...))
    }

    rule = Rule{
        Name: name,
        Priority: priority,
        SetSegments: Sets,
        Backends: rule_to_backend,
    }
    return &rule

}

func CreateTestConfigDefaultBackends() []*Backend {
    var Backends []*Backend
    Backends = []*Backend{
        CreateTestConfigBackend(
            "Example",
            90000,
            []*Url{&Url{Url: "http://ya.ru", Format: "openrtb2.2"}},
        ),
    }
    return Backends
}


func CreateTestConfigDefaultRules() []*Rule {
    var Rules []*Rule
    Rules = []*Rule{
        CreateTestConfigRule(
            "Example",
            100,
            nil,
            []int{1,2,3},
        ),
    }
    return Rules
}


func CreateTestConfig(Backend_timeout int, Backends []*Backend, Rules []*Rule) *Config {
    var config *Config 
    for _, backend := range Backends {
        ConfigSetHandlerForBackend(backend)
    }
    config = &Config{
        Backend_timeout: Backend_timeout,
        Backends: Backends,
        Rules: Rules,
        Micros: 1000000,
        SystemReportTo: "drop",
        Storage: "local",
    }
    return config
}


func CreateDefaultTestConfig() *Config{
    return CreateTestConfig(90000, CreateTestConfigDefaultBackends(), CreateTestConfigDefaultRules())
}


func CreateEmptyLocalStorage() Storage {
    return GetLocalStorage()
}


func CreateLocalStorage(data map[string] map[string] map[string] interface{}) Storage {
    return UpdateLocalStorage(GetLocalStorage(), data)
}
 

func init() {
    var wrong string = "QUFBQUJCQkJDQ0NDRERERA=="

    Segments = set.New(
        65, 78, 85, 242, 259, 956, 963, 1336, 1337, 1345, 1346, 1491, 1953, 1958, 2004, 2005,
        2011, 2061, 2065, 2066, 2079, 2081, 2082, 99601, 99606, 111702, 174102, 174104,
        174108, 1117014, 1117042, 1117069, 1117081, 1117082, 11170120, 11170133,
        11170148, 11170150, 11170152, 11170178, 11170182, 11170206, 11170208,
        11170241, 11170337, 11170347, 11170410, 11170772, 11170781, 11170894,
        11170984, 92302097, 92402097, 111701091, 111701142, 111701222, 111701430,
        111701726, 111701857, 111703068, 111703412, 111704366, 111706166, 111708361,
        111709519, 111709520, 111709531, 174302271, 174302371, 174302571,
    )
    facetz_id = &FacetzUser{
        Id:       "9c855c1f-b645-435b-b5e9-5695cd830081",
        Segments: Segments,
    }
    bid_request = new(google.BidRequest)
    bid_request_wrong_facetz_id = new(google.BidRequest)


    bid_request_body, _ = hex.DecodeString(
        "12104145286126d115c1e20ec2331e33a11c2203bc8a0732604d6f7a696c6c612f352e3" +
        "0202857696e646f77733b20553b2057696e646f7773204e5420352e313b20656e2d5553" +
        "3b2072763a312e382e312e3270726529204765636b6f2f3230303730313138204669726" +
        "5666f782f322e302e302e327072655a1c687474703a2f2f7777772e796f75747562652e" +
        "636f6d2f73686f77736202656e6a0808b30415e3ca0c3c6a0808ba041560e8b03e6a080" +
        "8fb0215dd9fb43d6a0808fc02150427973d6a07084215032bfb3e722d08b60110a40318" +
        "d804220308020332099801c201bb021ce6014a0910f9c2e5082890e52279f8bdee01000" +
        "000007801a00101aa0116674f31426d416667514c61564966366f3636385f3341b20107" +
        "596f7574756265f801b8e84b8a020432373533b802d6863dc802f603d10242bdee01000" +
        "00000")

    proto.Unmarshal(bid_request_body, bid_request)
    proto.Unmarshal(bid_request_body, bid_request_wrong_facetz_id)
    bid_request_wrong_facetz_id.GoogleUserId = &wrong

    openrtb_response_body = []byte("{\"cur\": \"RUB\", \"seatbid\": [{\"bid\": [{\"price\": 0.6, \"adm\": \"<a href='http://blabla.ru/?url=%%CLICK_URL_ESC%%'><img src='http://blabla.ru/1.gif?url=%%CLICK_URL_ESC%%'/></a>\", \"ext\": {\"url\": \"http://shop-zilla.ru/?url=%%CLICK_URL_ESC%%\"}, \"impid\": \"182\", \"id\": \"249621612006772521663818705005073365670\"}]}], \"id\": \"g\\u0000(\\u001b\\ufffd\\ufffd\\ufffd0\\ufffd\\ufffdxD<\\u0004\\ufffd[\", \"bidid\": \"1bb9245cc40e44789063b350525bd049\"}")
    openrtb_response_body_low_price = []byte("{\"cur\": \"RUB\", \"seatbid\": [{\"bid\": [{\"price\": 0.569, \"adm\": \"<a href='http://blabla.ru'><img src='http://blabla.ru/1.gif?url=%%CLICK_URL_ESC%%'/></a>\", \"ext\": {\"url\": \"http://shop-zilla.ru/?url=%%CLICK_URL_ESC%%\"}, \"impid\": \"182\", \"id\": \"249621612006772521663818705005073365670\"}]}], \"id\": \"g\\u0000(\\u001b\\ufffd\\ufffd\\ufffd0\\ufffd\\ufffdxD<\\u0004\\ufffd[\", \"bidid\": \"1bb9245cc40e44789063b350525bd049\"}")
    openrtb_response_body_bad = []byte("{\"cur\": \"RUB\", \"seatbid\": [{\"bid\": [{\"price\": 0.6, \"impid\": \"182\", \"id\": \"249621612006772521663818705005073365670\"}]}], \"id\": \"g\\u0000(\\u001b\\ufffd\\ufffd\\ufffd0\\ufffd\\ufffdxD<\\u0004\\ufffd[\", \"bidid\": \"1bb9245cc40e44789063b350525bd049\"}")
    proto.Unmarshal(bid_request_body, bid_request)
    json.Unmarshal(openrtb_response_body, &openrtb_response)
    default_context = CreateContext(facetz_id, bid_request)
}
