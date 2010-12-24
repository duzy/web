package eBay

import (
        "fmt"
        "xml"
        "io"
        "bytes"
)

type eBayFindingServiceCall struct {
        affiliate *Affiliate                    // Optional
        buyerPostalCode string                  // Optional
        paginationInput *PaginationInput        // Optional
        sortOrder string                        // Optional

        // TODO: aspectFilter, domainFilter, itemFilter, outputSelector, ...
}

type eBayFindingService_FindItemsByCategory struct {
        eBayFindingServiceCall
}

type eBayFindingService_FindItemsByProduct struct {
        eBayFindingServiceCall
}

type eBayFindingService_FindItemsIneBayStores struct {
        eBayFindingServiceCall
}

type eBayFindingService_FindItemsByKeywords struct {
        eBayFindingServiceCall
        keywords string
}

type eBayFindingService_FindItemsAdvanced struct {
        eBayFindingServiceCall
}

func (call *eBayFindingServiceCall) setEntriesPerPage(count int) {
        if call.paginationInput == nil {
                call.paginationInput = &PaginationInput{}
        }
        call.paginationInput.EntriesPerPage = count
}

func (call *eBayFindingServiceCall) getMessage(name string) io.ReadWriter {
        xmlns := "http://www.ebay.com/marketplace/search/v1/services"
        msg := bytes.NewBuffer(make([]byte, 0, 512))
        fmt.Fprintf(msg, "<?xml version='1.0' encoding='UTF-8'?>")
        fmt.Fprintf(msg, "<%s xmlns=\"%s\">", name, xmlns)
        if call.affiliate != nil {
                // TODO: xml.Escape(writer, bytes)...
                fmt.Fprint(msg, "<affiliate>")
                if call.affiliate.CustomId != "" {
                        fmt.Fprintf(msg, "<customId>%s</customId>", call.affiliate.CustomId)
                }
                if call.affiliate.NetworkId != "" {
                        fmt.Fprintf(msg, "<networkId>%s</networkId>", call.affiliate.NetworkId)
                }
                if call.affiliate.TrackingId != "" {
                        fmt.Fprintf(msg, "<trackingId>%s</trackingId>", call.affiliate.TrackingId)
                }
                fmt.Fprint(msg, "</affiliate>")
        }
        if call.buyerPostalCode != "" {
                fmt.Fprintf(msg, "<buyerPostalCode>%s</buyerPostalCode>", call.buyerPostalCode)
        }
        if call.paginationInput != nil {
                fmt.Fprint(msg, "<paginationInput>")
                if 0 < call.paginationInput.EntriesPerPage {
                        fmt.Fprintf(msg, "<entriesPerPage>%v</entriesPerPage>", call.paginationInput.EntriesPerPage)
                }
                if 0 < call.paginationInput.PageNumber {
                        fmt.Fprintf(msg, "<pageNumber>%v</pageNumber>", call.paginationInput.PageNumber)
                }
                fmt.Fprint(msg, "</paginationInput>")
        }
        if call.sortOrder != "" {
                fmt.Fprintf(msg, "<sortOrder>%s</sortOrder>", call.sortOrder)
        }
        return io.ReadWriter(msg)
}

func (call *eBayFindingServiceCall) getURL(app *App, oper string) string {
        u := URL_eBayFindingService
        u += "?OPERATION-NAME=" + oper
        u += "&SERVICE-VERSION=" + app.ServiceVersion
        u += "&SECURITY-APPNAME=" + app.AppID
        u += "&GLOBAL-ID=" + app.GlobalID
        if app.ResponseFormat != "" {
                u += "&RESPONSE-DATA-FORMAT=" + app.ResponseFormat
        }
        // u += "&REST-PAYLOAD"
        return u
}

func (call *eBayFindingService_FindItemsByCategory)   GetURL(app *App) string { return call.getURL(app, "findItemsByCategory") }
func (call *eBayFindingService_FindItemsByProduct)    GetURL(app *App) string { return call.getURL(app, "findItemsByProduct") }
func (call *eBayFindingService_FindItemsIneBayStores) GetURL(app *App) string { return call.getURL(app, "findItemsIneBayStores") }
func (call *eBayFindingService_FindItemsByKeywords)   GetURL(app *App) string { return call.getURL(app, "findItemsByKeywords") }
func (call *eBayFindingService_FindItemsAdvanced)     GetURL(app *App) string { return call.getURL(app, "findItemsAdvanced") }

func (api *eBayFindingService_FindItemsByKeywords) GetMessage(app *App) io.Reader {
        keywords := bytes.NewBuffer(make([]byte, 0, 128))
        xml.Escape(keywords, []byte(api.keywords))

        // http://developer.ebay.com/DevZone/finding/CallRef/findItemsByKeywords.html
        msg := api.getMessage("findItemsByKeywordsRequest")
        fmt.Fprintf(msg, "<keywords>%s</keywords>", keywords)
        fmt.Fprintf(msg, "</findItemsByKeywordsRequest>")
        return io.Reader(msg)
}
