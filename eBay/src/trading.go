package eBay

import (
        //"os"
        "io"
        "fmt"
        "xml"
        "bytes"
        "reflect"
)

type TradingService struct {
        app *App
}

type eBayTradingCall struct {
        DetailLevel string // see DetailLevelCodeType
        ErrorLanguage string
        MessageID string
        OutputSelector string
        Version string
        WarningLevel string // see WarningLevelCodeType
}

func (call *eBayTradingCall) GetHeaders(app *App) (h map[string]string) {
        h = make(map[string]string)
        h["X-EBAY-API-COMPATIBILITY-LEVEL"] = "697" //"613"
        h["X-EBAY-API-DEV-NAME"] = app.DEVID
        h["X-EBAY-API-APP-NAME"] = app.AppID
        h["X-EBAY-API-CERT-NAME"] = app.CertID
        h["X-EBAY-API-CALL-NAME"] = "GetCategories"
        h["X-EBAY-API-SITEID"] = "0"
        return
}

func (call *eBayTradingCall) GetURL(app *App) string { return URL_eBayTrading }

func eBayTradingCallOpName(call interface{}) (op string) {
        switch call.(type) {
        case *eBayTradingCall_GetCategories: op = "GetCategories"
        }
        return
}

func (call *eBayTradingCall) newMessage(ncall interface{}, app *App) (msg *bytes.Buffer, ml int) {
        userToken := bytes.NewBuffer(make([]byte, 0, 128))
        xml.Escape(userToken, []uint8(app.UserToken))

        name := eBayTradingCallOpName(ncall) + "Request"
        if name == "" { return nil, 0 }

        xmlns := "urn:ebay:apis:eBLBaseComponents"
        msg = bytes.NewBuffer(make([]byte, 0, 128))
        fmt.Fprintf(msg, `<?xml version="1.0" encoding="UTF-8" ?>`)
        fmt.Fprintf(msg, `<%s xmlns="%s">`, name, xmlns)
        fmt.Fprintf(msg, `<RequesterCredentials>`)
        fmt.Fprintf(msg, `<eBayAuthToken>%s</eBayAuthToken>`, userToken)
        fmt.Fprintf(msg, `</RequesterCredentials>`)

        f := func(t *reflect.StructField, v reflect.Value) (nxt bool) {
                if t.Anonymous { return true }

                //fmt.Printf("field: %s = %v\n", t.Name, v);

                set := true
                switch a := v.Interface().(type) {
                case string:    set = a != ""
                case int:       //set = a != 0
                case bool:      //set = a
                default:
                        fmt.Printf("todo: field: %s = %v\n", t.Name, v);
                }

                if set { fmt.Fprintf(msg, `<%s>%v</%s>`, t.Name, v.Interface(), t.Name) }
                return true // continue with the next field
        }
        ForEachField(call, f)
        ForEachField(ncall, f)
        
        fmt.Fprintf(msg, `</%s>`, name)
        ml = msg.Len()

        //fmt.Printf("request: %v\n", msg);
        return
}

// see http://developer.ebay.com/DevZone/XML/docs/Reference/eBay/GetCategories.html
type eBayTradingCall_GetCategories struct {
        eBayTradingCall
        CategoryParent string
        CategorySiteID string
        LevelLimit int
        ViewAllNodes bool
}

func (call *eBayTradingCall_GetCategories) GetMessage(app *App) (io.Reader, int) { return call.newMessage(call, app) }

func (svc *TradingService) NewGetCategoriesCall() *eBayTradingCall_GetCategories { return &eBayTradingCall_GetCategories{} }
