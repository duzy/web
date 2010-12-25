package eBay

import (
        "os"
        "io"
        "fmt"
        "xml"
        "bytes"
)

type FindingService struct {
        app *App
}

type eBayFindingServiceCall struct {
        affiliate *Affiliate                    // Optional
        buyerPostalCode string                  // Optional
        paginationInput *PaginationInput        // Optional
        sortOrder string                        // Optional

        // TODO: aspectFilter, domainFilter, itemFilter, outputSelector, ...
}

type eBayFindingService_GetVersion struct {
        eBayFindingServiceCall
}

type eBayFindingService_GetHistograms struct {
        eBayFindingServiceCall
        categoryId string
}

type eBayFindingService_GetSearchKeywordsRecommendation struct {
        eBayFindingServiceCall
        keywords string
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

func (call *eBayFindingServiceCall) clear() {
        call.affiliate = nil
        call.buyerPostalCode = ""
        call.paginationInput = nil
        call.sortOrder = ""
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

func (call *eBayFindingService_GetVersion)              GetURL(app *App) string { return call.getURL(app, "getVersion") }
func (call *eBayFindingService_GetHistograms)           GetURL(app *App) string { return call.getURL(app, "getHistograms") }
func (call *eBayFindingService_GetSearchKeywordsRecommendation) GetURL(app *App) string { return call.getURL(app, "getSearchKeywordsRecommendation") }
func (call *eBayFindingService_FindItemsByCategory)     GetURL(app *App) string { return call.getURL(app, "findItemsByCategory") }
func (call *eBayFindingService_FindItemsByProduct)      GetURL(app *App) string { return call.getURL(app, "findItemsByProduct") }
func (call *eBayFindingService_FindItemsIneBayStores)   GetURL(app *App) string { return call.getURL(app, "findItemsIneBayStores") }
func (call *eBayFindingService_FindItemsByKeywords)     GetURL(app *App) string { return call.getURL(app, "findItemsByKeywords") }
func (call *eBayFindingService_FindItemsAdvanced)       GetURL(app *App) string { return call.getURL(app, "findItemsAdvanced") }

func (api *eBayFindingService_GetVersion) GetMessage(app *App) io.Reader {
        api.clear()

        // http://developer.ebay.com/DevZone/finding/CallRef/getVersion.html
        msg := api.getMessage("getVersionRequest")
        fmt.Fprint(msg, "</getVersionRequest>")
        return io.Reader(msg)
}

func (api *eBayFindingService_GetHistograms) GetMessage(app *App) io.Reader {
        api.clear()

        // http://developer.ebay.com/DevZone/finding/CallRef/getHistograms.html
        msg := api.getMessage("getHistogramsRequest")
        fmt.Fprintf(msg, "<categoryId>%s</categoryId>", api.categoryId)
        fmt.Fprint(msg, "</getHistogramsRequest>")
        return io.Reader(msg)
}

func (api *eBayFindingService_GetSearchKeywordsRecommendation) GetMessage(app *App) io.Reader {
        keywords := bytes.NewBuffer(make([]byte, 0, 128))
        xml.Escape(keywords, []byte(api.keywords))

        api.clear()

        // http://developer.ebay.com/DevZone/finding/CallRef/getSearchKeywordsRecommendation.html
        msg := api.getMessage("getSearchKeywordsRecommendationRequest")
        fmt.Fprintf(msg, "<keywords>%s</keywords>", keywords)
        fmt.Fprint(msg, "</getSearchKeywordsRecommendationRequest>")
        return io.Reader(msg)
}

// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsByCategory.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsByProduct.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsIneBayStores.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsAdvanced.html

func (api *eBayFindingService_FindItemsByKeywords) GetMessage(app *App) io.Reader {
        keywords := bytes.NewBuffer(make([]byte, 0, 128))
        xml.Escape(keywords, []byte(api.keywords))

        // http://developer.ebay.com/DevZone/finding/CallRef/findItemsByKeywords.html
        msg := api.getMessage("findItemsByKeywordsRequest")
        fmt.Fprintf(msg, "<keywords>%s</keywords>", keywords)
        fmt.Fprint(msg, "</findItemsByKeywordsRequest>")
        return io.Reader(msg)
}

func (eb *FindingService) GetVersion() (str string, err os.Error) {
        call := &eBayFindingService_GetVersion{}
        str, err = eb.app.post(call)
        return
}

func (eb *FindingService) FindItemsByKeywords(keywords string, count int) (str string, err os.Error) {
        call := &eBayFindingService_FindItemsByKeywords{}
        call.setEntriesPerPage(count)
        call.keywords = keywords
        str, err = eb.app.post(call)
        //fmt.Printf("%s\n%s\n", call.GetURL(eb), call.GetMessage(eb))
        return
}

func (eb *FindingService) FindItemsAdvanced() (str string, err os.Error) {
        // TODO: ...
        return
}
