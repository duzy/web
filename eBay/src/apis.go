package eBay

import (
        "io/ioutil"
        "http"
        "fmt"
)

const (
        eBayEndpoint = "http://svcs.ebay.com/services/search/FindingService/v1"
)

type EBayAPI struct {
        endpoint string
        version string
        appid string // Production AppID
        globalid string
        query string
        op string
        entries int
        format string
}

type EBay struct {
        //DEVID string
        AppID string
        //CertID string

        ResponseFormat string
}

func NewEBay() (eb *EBay) {
        eb = &EBay {
        AppID: "dusellco-2abe-4ae8-8bc6-5fd8dc98b37e",
        }
        return
}

func NewCall(op string) (api *EBayAPI) {
        api = &EBayAPI{
        endpoint: eBayEndpoint,
        version: "1.0.0",
        appid: "",
        globalid: "EBAY-US",
        query: "",
        op: op,
        entries: 0,
        }
        return
}

// construct the HTTP GET call URL
func (api *EBayAPI) GetURL() (u string) {
        u  = api.endpoint + "?"
        u += "OPERATION-NAME=" + api.op
        u += "&SERVICE-VERSION=" + api.version
        u += "&SECURITY-APPNAME=" + api.appid
        u += "&GLOBAL-ID=" + api.globalid
        if api.format != "" {
                u += "&RESPONSE-DATA-FORMAT=" + api.format
        }
        u += "&keywords=" + http.URLEscape(api.query)
        u += "&paginationInput.entriesPerPage=" + fmt.Sprint(api.entries)
        return
}

func (eb *EBay) FindItemsByKeywords(keywords string, count int) (str string) {
        call := NewCall("findItemsByKeywords")
        call.appid = eb.AppID
        call.query = keywords
        call.entries = count
        call.format = eb.ResponseFormat
        str = eb.get(call)
        return
}

func (eb *EBay) get(call *EBayAPI) (str string) {
        u := call.GetURL()
        fmt.Printf("%s\n", u)
        r, _, err := http.Get(u)
        if err == nil {
                var b []byte
                b, err = ioutil.ReadAll(r.Body)
                r.Body.Close()
                str = string(b)
        }
        return
}
