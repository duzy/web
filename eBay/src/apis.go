package eBay

import (
        "os"
        "io"
        "io/ioutil"
        "http"
)

type eBayServiceCall interface {
        GetURL(app *App) string
        GetMessage(app *App) io.Reader
}

type App struct {
        DEVID string
        AppID string
        CertID string

        ServiceVersion string
        GlobalID string
        ResponseFormat string
}

func NewApp() (eb *App) {
        eb = &App {
                // Sandbox Key Set
                /*
        DEVID: "c5f14b63-0bf9-405f-8c5c-efaaba2b4a02",
        AppID: "dusellco-da1b-434b-9d10-2448ee5fc58a",
        CertID: "4d6382a0-aad7-4b93-92d1-4f558471c576",
                 */

                // Production Key Set
        DEVID: "c5f14b63-0bf9-405f-8c5c-efaaba2b4a02",
        AppID: "dusellco-2abe-4ae8-8bc6-5fd8dc98b37e",
        CertID: "87aab9ab-375c-41e5-bf14-9702fec7dec3",

        GlobalID: "EBAY-US",
        ServiceVersion: "1.8.0",
        }
        return
}

func (eb *App) get(call eBayServiceCall) (str string, err os.Error) {
        u := call.GetURL(eb)
        r, _, err := http.Get(u)
        if err == nil {
                var b []byte
                b, err = ioutil.ReadAll(r.Body)
                r.Body.Close()
                str = string(b)
        }
        return
}

func (eb *App) post(call eBayServiceCall) (str string, err os.Error) {
        t := "text/xml" // TODO: text/json
        u := call.GetURL(eb)
        r, err := http.Post(u, t, call.GetMessage(eb))
        if err == nil {
                var b []byte
                b, err = ioutil.ReadAll(r.Body)
                r.Body.Close()
                str = string(b)
        }
        return
}

func (eb *App) GetVersion() (str string, err os.Error) {
        call := &eBayFindingService_GetVersion{}
        str, err = eb.post(call)
        return
}

func (eb *App) FindItemsByKeywords(keywords string, count int) (str string, err os.Error) {
        call := &eBayFindingService_FindItemsByKeywords{}
        call.setEntriesPerPage(count)
        call.keywords = keywords
        str, err = eb.post(call)
        //fmt.Printf("%s\n%s\n", call.GetURL(eb), call.GetMessage(eb))
        return
}

func (eb *App) FindItemsAdvanced() (str string, err os.Error) {
        // TODO: ...
        return
}