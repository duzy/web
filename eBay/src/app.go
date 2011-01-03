package eBay

import (
        "os"
        "io"
        "io/ioutil"
        "http"
        "bufio"
        "strings"
        "fmt"
        "net"
        "crypto/tls"
        "encoding/base64"
)

type Cacher interface {
        Close() (err os.Error)
        CacheCategory(cat *Category) (err os.Error)
        CacheItem(itm *Item) (err os.Error)
        GetCategory(id string) (cat *Category, err os.Error)
        GetItem(id string) (itm *Item, err os.Error)
}

type App struct {
        DEVID string
        AppID string
        CertID string
        UserToken string

        ServiceVersion string
        GlobalID string
        ResponseFormat string

        cache Cacher
}

type eBayServiceCall interface {
        GetHeaders(app *App) map[string]string
        GetURL(app *App) string
        GetMessage(app *App) (io.Reader, int)
}

func NewApp(sandbox bool) (app *App) {
        if sandbox {
                app = &App {
                DEVID: "c5f14b63-0bf9-405f-8c5c-efaaba2b4a02",
                AppID: "dusellco-da1b-434b-9d10-2448ee5fc58a",
                CertID: "4d6382a0-aad7-4b93-92d1-4f558471c576",
                UserToken: `AgAAAA**AQAAAA**aAAAAA**hXIgTQ**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4CoAZCCpA2dj6x9nY+seQ**xnMBAA**AAMAAA**s7Ft71dh0bgNFPn/VC3cGVUGv7B162v7Dq1lzTamuW/qLuycuUz4A4Io82+wv4ESB3DHXqywAirzNZ23cbWSzGccPnKe9yei6S433LhOpWHKimqztXvPnPFFQfiI6fJa/+Wfd4RAKjkGu0ymbi0tJ4WF1Xd50Iz9ZJxtHUNgkHM6qjnUe/q9SkY+cVoL25jFX6lyBzTUTklsKd/ASBLsFqItuUL45v7kpLAba/MqcNe35PFIGQ61nNKs+nAUNrE7mMizmAd0eXCsIhdtcC75fERplsvZGNxD+GudZMJjagWoUhIcD49yDvOVl9AmRqi72NjDiCTXqk8B2Hv/I/FMe4Ig5vVjRduzxR4AlwyxyP/ZhEqX951GEML2mCJiaGyNz+vlJYkFeccrxYE8QymmJ+iSABDZx8Qmcz8s7LFn++YBFoGjtLpMgOyzH/JSlbefpzh8JaaDFoE0P2u17e6/wEjfU9bjibBY5Evb2qFCEKjBorTB6U+fsf4ST8WZCWItHXbJjgyNsI3DmuOXkYTpKGj+HsnlKDuJmMPtmOrgkVXaBysH30u7WHyDYtAQNBE4s9Nr8DLhPemWK78y52layS2xzc/qFRtJnsWQ5AZeRNgJXx9M4PQZD36VqNFCMozHBKQB6HMNL/hx/DlkDs2likpYbr0ksSjBkCpnkfJCDCR9tVAityqW27sz4ukYkKKXGYynlqiuS0Ds6hnIwUQy2H2uwSnEz8oADfKfZzq07HKOyvV9CyhL8jHZFQ1Uwcl5`,
                GlobalID: "EBAY-US",
                ServiceVersion: "1.8.0",
                ResponseFormat: "XML",
                }
        } else {
                app = &App {
                DEVID: "c5f14b63-0bf9-405f-8c5c-efaaba2b4a02",
                AppID: "dusellco-2abe-4ae8-8bc6-5fd8dc98b37e",
                CertID: "87aab9ab-375c-41e5-bf14-9702fec7dec3",
                UserToken: ``,
                GlobalID: "EBAY-US",
                ServiceVersion: "1.8.0",
                ResponseFormat: "XML",
                }
        }
        return
}

func (app *App) get(call eBayServiceCall) (str string, err os.Error) {
        u := call.GetURL(app)
        r, _, err := http.Get(u)
        if err == nil {
                var b []byte
                b, err = ioutil.ReadAll(r.Body)
                r.Body.Close()
                str = string(b)
        }
        return
}

func (app *App) post_(call eBayServiceCall) (str string, err os.Error) {
        t := "text/xml" // TODO: text/json
        u := call.GetURL(app)
        msg, _ := call.GetMessage(app)
        r, err := http.Post(u, t, msg)
        if err == nil {
                var b []byte
                b, err = ioutil.ReadAll(r.Body)
                r.Body.Close()
                str = string(b)
        }
        return
}

func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

type nopCloser struct { io.Reader }
func (nopCloser) Close() os.Error { return nil }

func (app *App) post(call eBayServiceCall) (str string, err os.Error) {
        msg, ml := call.GetMessage(app)

        var req http.Request
        req.Method = "POST"
        req.ProtoMajor = 1
        req.ProtoMinor = 1
        req.Close = true
        req.Body = nopCloser{ msg }
        req.Header = call.GetHeaders(app)
        if req.Header == nil { req.Header = make(map[string]string) }
        req.Header["Content-Type"] = "text/xml" // TODO: text/json
        req.Header["Content-Length"] = fmt.Sprintf("%d", ml)
        req.ContentLength = int64(ml)
        //req.TransferEncoding = []string{"chunked"}
        req.URL, err = http.ParseURL(call.GetURL(app))
        if err != nil { return }

        if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
                err = os.NewError("unsupported protocol scheme: "+req.URL.Scheme)
                return
        }

        addr := req.URL.Host
        if !hasPort(addr) { addr += ":" + req.URL.Scheme }

        info := req.URL.RawUserinfo
        if 0 < len(info) {
                enc := base64.URLEncoding
                encoded := make([]byte, enc.EncodedLen(len(info)))
                enc.Encode(encoded, []byte(info))
                req.Header["Authorization"] = "Basic " + string(encoded)
        }

        var conn io.ReadWriteCloser
        if req.URL.Scheme == "http" {
                conn, err = net.Dial("tcp", "", addr)
                if err != nil { return }
        } else { // https
                conn, err = tls.Dial("tcp", "", addr, nil)
                if err != nil { return }

                h := req.URL.Host
                if hasPort(h) {
                        h = h[0:strings.LastIndex(h, ":")]
                }
                if err = conn.(*tls.Conn).VerifyHostname(h); err != nil {
                        return
                }
        }

        if err = req.Write(conn); err != nil { conn.Close(); return }

        reader := bufio.NewReader(conn)
        resp, err := http.ReadResponse(reader, req.Method)
        if err != nil { conn.Close(); return }

        var buf []byte
        buf, err = ioutil.ReadAll(resp.Body)
        conn.Close()

        str = string(buf)
        return
}

// NewFindingService returns a new eBay.FindingService.
func (app *App) NewFindingService() (p *FindingService) {
        p = &FindingService{ app }
        return
}

func (app *App) NewTradingService() (p *TradingService) {
        p = &TradingService{ app }
        return
}

func (app *App) Invoke(call eBayServiceCall) (str string, err os.Error) {
        return app.post(call)
}

func (app *App) GetCache() Cacher { return app.cache }
func (app *App) SetCache(cache Cacher) (prev Cacher) {
        prev = app.cache
        app.cache = cache
        return
}
