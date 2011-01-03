package eBay

import (
        "os"
        "io"
        "fmt"
        "xml"
        "json"
        "bytes"
        "strconv"
        "strings"
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

func (call *eBayFindingServiceCall) GetHeaders(app *App) (h map[string]string) {
        return
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

func (call *eBayFindingServiceCall) getMessage(name string) (io.ReadWriter, int) {
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
        return io.ReadWriter(msg), msg.Len()
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

func (api *eBayFindingService_GetVersion) GetMessage(app *App) (io.Reader, int) {
        api.clear()

        // http://developer.ebay.com/DevZone/finding/CallRef/getVersion.html
        msg, ml := api.getMessage("getVersionRequest")
        fmt.Fprint(msg, "</getVersionRequest>")
        return io.Reader(msg), ml
}

func (api *eBayFindingService_GetHistograms) GetMessage(app *App) (io.Reader, int) {
        api.clear()

        // http://developer.ebay.com/DevZone/finding/CallRef/getHistograms.html
        msg, ml := api.getMessage("getHistogramsRequest")
        fmt.Fprintf(msg, "<categoryId>%s</categoryId>", api.categoryId)
        fmt.Fprint(msg, "</getHistogramsRequest>")
        return io.Reader(msg), ml
}

func (api *eBayFindingService_GetSearchKeywordsRecommendation) GetMessage(app *App) (io.Reader, int) {
        keywords := bytes.NewBuffer(make([]byte, 0, 128))
        xml.Escape(keywords, []byte(api.keywords))

        api.clear()

        // http://developer.ebay.com/DevZone/finding/CallRef/getSearchKeywordsRecommendation.html
        msg, ml := api.getMessage("getSearchKeywordsRecommendationRequest")
        fmt.Fprintf(msg, "<Keywords>%s</Keywords>", keywords)
        fmt.Fprint(msg, "</getSearchKeywordsRecommendationRequest>")
        return io.Reader(msg), ml
}

// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsByCategory.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsByProduct.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsIneBayStores.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsAdvanced.html

func (api *eBayFindingService_FindItemsByKeywords) GetMessage(app *App) (io.Reader, int) {
        keywords := bytes.NewBuffer(make([]byte, 0, 128))
        xml.Escape(keywords, []byte(api.keywords))

        // http://developer.ebay.com/DevZone/finding/CallRef/findItemsByKeywords.html
        msg, ml := api.getMessage("findItemsByKeywordsRequest")
        fmt.Fprintf(msg, "<keywords>%s</keywords>", keywords)
        fmt.Fprint(msg, "</findItemsByKeywordsRequest>")
        return io.Reader(msg), ml
}

func (svc *FindingService) GetVersion() (str string, err os.Error) {
        call := &eBayFindingService_GetVersion{}
        str, err = svc.app.post(call)
        return
}

func (svc *FindingService) FindItemsByKeywords(keywords string, count int) (str string, err os.Error) {
        call := &eBayFindingService_FindItemsByKeywords{}
        call.setEntriesPerPage(count)
        call.keywords = keywords
        str, err = svc.app.post(call)
        //fmt.Printf("%s\n%s\n", call.GetURL(eb), call.GetMessage(eb))
        return
}

func (svc *FindingService) FindItemsAdvanced() (str string, err os.Error) {
        // TODO: ...
        return
}

// ParseResponse parse text response of an eBay operation.
func (svc *FindingService) ParseResponse(str string) (res *findItemsResponse, err os.Error) {
        switch svc.app.ResponseFormat {
        case "JSON":
                res, err = svc.parseJSONResponse(str)
        case "XML":
                res, err = svc.parseXMLResponse(str)
        default:
                err = os.NewError("unknown data format '"+svc.app.ResponseFormat+"'")
        }
        return
}

// parseXMLResponse parse XML format response
func (svc *FindingService) parseXMLResponse(str string) (res *findItemsResponse, err os.Error) {
        p := xml.NewParser(bytes.NewBufferString(str))

        var start *xml.StartElement
        for {
                tok, err := p.Token()
                if err != nil { return }
                if t, ok := tok.(xml.StartElement); ok {
                        start = &t
                        break;
                }
        }

        if start == nil {
                err = os.NewError("no xml.StartElement found")
                return
        }

        res = new(findItemsResponse)

        var v interface{}
        switch start.Name.Local {
        case "findItemsByCategoryResponse":
                v = &struct {
                        XMLName xml.Name "findItemsByCategoryResponse"
                        *findItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsByProductResponse":
                v = &struct {
                        XMLName xml.Name "findItemsByProductResponse"
                        *findItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsIneBayStoresResponse":
                v = &struct {
                        XMLName xml.Name "findItemsIneBayStoresResponse"
                        *findItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsByKeywordsResponse":
                v = &struct {
                        XMLName xml.Name "findItemsByKeywordsResponse"
                        *findItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsAdvancedResponse":
                v = &struct {
                        XMLName xml.Name "findItemsAdvancedResponse"
                        *findItemsResponse
                }{ xml.Name{}, res, }
                // TODO: more response
        }

        if v == nil {
                err = os.NewError(fmt.Sprintf("don't know how to parse '%s'",start.Name))
                return
        }

        err = p.Unmarshal(v, start)
        return
}

func getJSONResponseType(str string) (typ string) {
        n := strings.Index(str, `"` /*`{"`*/)
        if n == -1 { return }

        str = str[n+1 /*2*/:]
        n = strings.Index(str, `":`)
        if n == -1 { return }

        typ = str[0:n]
        return
}

// parseJSONResponse parse JSON format response
func (svc *FindingService) parseJSONResponse(str string) (res *findItemsResponse, err os.Error) {
        ra := make([]*findItemsJSONResponse, 1)

        var v interface{}
        switch t := getJSONResponseType(str[0:]); t {
        case "findItemsByCategoryResponse":
                v = &struct { V []*findItemsJSONResponse "findItemsByCategoryResponse" }{ ra }
        case "findItemsByProductResponse":
                v = &struct { V []*findItemsJSONResponse "findItemsByProductResponse" }{ ra }
        case "findItemsIneBayStoresResponse":
                v = &struct { V []*findItemsJSONResponse "findItemsIneBayStoresResponse" }{ ra }
        case "findItemsByKeywordsResponse":
                v = &struct { V []*findItemsJSONResponse "findItemsByKeywordsResponse" }{ ra }
        case "findItemsAdvancedResponse":
                v = &struct { V []*findItemsJSONResponse "findItemsAdvancedResponse" }{ ra }
        default:
                err = os.NewError("unknown JSON response: '"+t+"'")
                return
        }

        err = json.Unmarshal([]byte(str), v)
        if err == nil { res = noJSON(ra[0]) }
        return
}

func stoi(s string) (i int) { i, _ = strconv.Atoi(s); return }
func stof(s string) (f float) { f, _ = strconv.Atof(s); return }
func stob(s string) (b bool) { b, _ = strconv.Atob(s); return }
func noJSON(r *findItemsJSONResponse) (res *findItemsResponse) {
        res = &findItemsResponse{
        Ack: r.Ack[0],
        Version: r.Version[0],
        Timestamp: r.Timestamp[0],
        //SearchResult: { make([]Item, len(r.SearchResult[0].Item)) },
        ItemSearchURL: r.ItemSearchURL[0],
        PaginationOutput: PaginationOutput{
                PageNumber: stoi(r.PaginationOutput[0].PageNumber[0]),
                EntriesPerPage: stoi(r.PaginationOutput[0].EntriesPerPage[0]),
                TotalPages: stoi(r.PaginationOutput[0].TotalPages[0]),
                TotalEntries: stoi(r.PaginationOutput[0].TotalEntries[0]),
                },
        }

        res.SearchResult.Item = make([]Item, len(r.SearchResult[0].Item))
        
        for n, i := range r.SearchResult[0].Item {
                res.SearchResult.Item[n] = Item{
                ItemId: i.ItemId[0],
                Title: i.Title[0],
                //GlobalId: i.GlobalId[0],
                PrimaryCategory: Category{
                        CategoryID: i.PrimaryCategory[0].CategoryId[0],
                        CategoryName: i.PrimaryCategory[0].CategoryName[0],
                        },
                GalleryURL: i.GalleryURL[0],
                ViewItemURL: i.ViewItemURL[0],
                PaymentMethod: i.PaymentMethod[0],
                AutoPay: stob(i.AutoPay[0]),
                Location: i.Location[0],
                Country: i.Country[0],
                ShippingInfo: ShippingInfo{
                        ShippingServiceCost: Money{
                                        i.ShippingInfo[0].ShippingServiceCost[0].CurrencyId,
                                        stof(i.ShippingInfo[0].ShippingServiceCost[0].Amount),
                                },
                        ShippingType: i.ShippingInfo[0].ShippingType[0],
                        ShipToLocations: strings.Join(i.ShippingInfo[0].ShipToLocations,","),
                        ExpeditedShipping: stob(i.ShippingInfo[0].ExpeditedShipping[0]),
                        OneDayShippingAvailable: stob(i.ShippingInfo[0].OneDayShippingAvailable[0]),
                        HandlingTime: stoi(i.ShippingInfo[0].HandlingTime[0]),
                        },
                SellingStatus: SellingStatus{
                        CurrentPrice: Money{
                                        i.SellingStatus[0].CurrentPrice[0].CurrencyId,
                                        stof(i.SellingStatus[0].CurrentPrice[0].Amount),
                                },
                        ConvertedCurrentPrice: Money{
                                        i.SellingStatus[0].ConvertedCurrentPrice[0].CurrencyId,
                                        stof(i.SellingStatus[0].ConvertedCurrentPrice[0].Amount),
                                },
                        BidCount: stoi(i.SellingStatus[0].BidCount[0]),
                        SellingState: i.SellingStatus[0].SellingState[0],
                        TimeLeft: i.SellingStatus[0].TimeLeft[0],
                        },
                ListingInfo: ListingInfo{
                        BestOfferEnabled: stob(i.ListingInfo[0].BestOfferEnabled[0]),
                        BuyItNowAvailable: stob(i.ListingInfo[0].BuyItNowAvailable[0]),
                        StartTime: i.ListingInfo[0].StartTime[0],
                        EndTime: i.ListingInfo[0].EndTime[0],
                        ListingType: i.ListingInfo[0].ListingType[0],
                        Gift: stob(i.ListingInfo[0].Gift[0]),
                        },
                ReturnsAccepted: stob(i.ReturnsAccepted[0]),
                Condition: Condition{
                        ConditionId: i.Condition[0].ConditionId[0],
                        ConditionDisplayName: i.Condition[0].ConditionDisplayName[0],
                        },
                }
        }//for (items)
        return
}
