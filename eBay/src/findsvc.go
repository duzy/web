package eBay

import (
        "http"
        "fmt"
        "os"
)

type eBayFindingService struct {
        endpoint string
        version string
        appid string // Production AppID
        globalid string
        query string
        op string
        entries int
        format string
}

func newFindingServiceCall(op string) (api *eBayFindingService) {
        api = &eBayFindingService{
        endpoint: URL_eBayFindingService,
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
func (api *eBayFindingService) GetURL() (u string) {
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

func (eb *App) FindItemsByKeywords(keywords string, count int) (str string, err os.Error) {
        call := newFindingServiceCall("findItemsByKeywords")
        call.appid = eb.AppID
        call.query = keywords
        call.entries = count
        call.format = eb.ResponseFormat
        str, err = eb.get(call)
        return
}
