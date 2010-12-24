package eBay

import (
        "io/ioutil"
        "http"
        "os"
)

type App struct {
        DEVID string
        AppID string
        CertID string

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
        }
        return
}

func (eb *App) get(call *eBayFindingService) (str string, err os.Error) {
        u := call.GetURL()

        //fmt.Printf("%s\n", u)

        r, _, err := http.Get(u)
        if err == nil {
                var b []byte
                b, err = ioutil.ReadAll(r.Body)
                r.Body.Close()
                str = string(b)
        }
        return
}
