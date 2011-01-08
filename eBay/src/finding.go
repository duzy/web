package eBay

import (
        "os"
        "io"
        "fmt"
        "xml"
        "json"
        "bytes"
        //"strconv"
        "strings"
        "reflect"
)

// TODO: get rid of this type
type findItemsResponseJSON struct {
        Ack []string
        Version []string
        Timestamp []string
        SearchResult []struct {
                Item []struct {
                        ItemId []string
                        Title []string
                        GlobalId []string
                        PrimaryCategory []struct {
                                CategoryId []string
                                CategoryName []string
                        }
                        GalleryURL []string
                        ViewItemURL []string
                        PaymentMethod []string
                        AutoPay []string
                        Location []string
                        Country []string
                        ShippingInfo []struct {
                                ShippingServiceCost []struct {
                                        CurrencyId string "@currencyId"
                                        Amount string "__value__"
                                }
                                ShippingType []string
                                ShipToLocations []string
                                ExpeditedShipping []string
                                OneDayShippingAvailable []string
                                HandlingTime []string
                        }
                        SellingStatus []struct {
                                CurrentPrice []struct {
                                        CurrencyId string "@currencyId"
                                        Amount string "__value__"
                                }
                                        ConvertedCurrentPrice []struct {
                                        CurrencyId string "@currencyId"
                                        Amount string "__value__"
                                }
                                BidCount []string
                                SellingState []string
                                TimeLeft []string
                        }
                        ListingInfo []struct {
                                BestOfferEnabled []string
                                BuyItNowAvailable []string
                                StartTime []string
                                EndTime []string
                                ListingType []string
                                Gift []string
                        }
                        ReturnsAccepted []string
                        Condition []struct {
                                ConditionId []string
                                ConditionDisplayName []string
                        }
                }
        }
        PaginationOutput []struct {
                PageNumber []string
                EntriesPerPage []string
                TotalPages []string
                TotalEntries []string
        }
        ItemSearchURL []string
}

type FindingService struct {
        app *App
}

type eBayFindingServiceCall struct {
        Affiliate Affiliate                     // Optional
        BuyerPostalCode string                  // Optional
        PaginationInput PaginationInput         // Optional
        SortOrder string                        // Optional

        // TODO: aspectFilter, domainFilter, itemFilter, outputSelector, ...
}

type eBayFindingService_GetVersion struct {
        eBayFindingServiceCall
}

type eBayFindingService_GetHistograms struct {
        eBayFindingServiceCall
        CategoryID string
}

type eBayFindingService_GetSearchKeywordsRecommendation struct {
        eBayFindingServiceCall
        Keywords string
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
        Keywords string
}

type eBayFindingService_FindItemsAdvanced struct {
        eBayFindingServiceCall
}

func (call *eBayFindingServiceCall) SetEntriesPerPage(count int) {
        call.PaginationInput.EntriesPerPage = count
}

func eBayFindingCallOpName(call interface{}) (op string) {
        switch call.(type) {
        case *eBayFindingService_GetVersion:                            op = "getVersion"
        case *eBayFindingService_GetHistograms:                         op = "getHistograms"
        case *eBayFindingService_GetSearchKeywordsRecommendation:       op = "getSearchKeywordsRecommendation"
        case *eBayFindingService_FindItemsByCategory:                   op = "findItemsByCategory"
        case *eBayFindingService_FindItemsByProduct:                    op = "findItemsByProduct"
        case *eBayFindingService_FindItemsIneBayStores:                 op = "findItemsIneBayStores"
        case *eBayFindingService_FindItemsByKeywords:                   op = "findItemsByKeywords"
        case *eBayFindingService_FindItemsAdvanced:                     op = "findItemsAdvanced"
        }
        return
}

func (call *eBayFindingServiceCall) GetURL(app *App) string { return URL_eBayFindingService }

func (call *eBayFindingServiceCall) getHeaders(ncall interface{}, app *App) (h map[string]string) {
        oper := eBayFindingCallOpName(ncall)
        h = map[string]string {
                "X-EBAY-SOA-SERVICE-NAME": "FindingService",
                "X-EBAY-SOA-SERVICE-VERSION": app.ServiceVersion,
                "X-EBAY-SOA-OPERATION-NAME": oper,
                "X-EBAY-SOA-GLOBAL-ID": app.GlobalID,
                "X-EBAY-SOA-SECURITY-APPNAME": app.AppID,
                "X-EBAY-SOA-REQUEST-DATA-FORMAT": "XML",
                "X-EBAY-SOA-RESPONSE-DATA-FORMAT": app.ResponseFormat,
        }
        return
}

func (call *eBayFindingServiceCall) newMessage(ncall interface{}) (io.ReadWriter, int) {
        const xmlns = "http://www.ebay.com/marketplace/search/v1/services"
        name := eBayFindingCallOpName(ncall) + "Request"
        msg := bytes.NewBuffer(make([]byte, 0, 512))
        fmt.Fprintf(msg, "<?xml version='1.0' encoding='UTF-8'?>")
        fmt.Fprintf(msg, "<%s xmlns=\"%s\">", name, xmlns)

        buf := bytes.NewBuffer(make([]byte, 0, 128))
        b0 := bytes.NewBuffer(make([]byte, 0, 128))
        f0 := func(t *reflect.StructField, v reflect.Value) (nxt bool) {
                if v == nil || v.Interface() == nil { return true }

                set := true
                switch a := v.Interface().(type) {
                case string:    set = a != ""
                case int:       set = a != 0
                case bool:      //set = a
                default:
                        fmt.Printf("todo: field: %s = %v\n", t.Name, v);
                }

                if set {
                        name := strings.ToLower(t.Name[0:1])
                        if l := len(t.Name) ; 1 < l { name += t.Name[1:l] }

                        buf.Reset()
                        xml.Escape(buf, []byte(fmt.Sprintf("%v", v.Interface())))
                        fmt.Fprintf(b0, "<%s>%s</%s>", name, buf, name)
                }
                return true
        }
        f := func(t *reflect.StructField, v reflect.Value) (nxt bool) {
                name := strings.ToLower(t.Name[0:1])
                if l := len(t.Name) ; 1 < l { name += t.Name[1:l] }

                switch p := v.(type) {
                case *reflect.StructValue:
                        b0.Reset()
                        ForEachField(p.Interface(), f0) // TODO: errors?
                        if 0 < b0.Len() {
                                fmt.Fprintf(msg, "<%s>%s</%s>", name, b0, name)
                        }
                        return true
                case *reflect.PtrValue:
                        if s, ok := p.Elem().(*reflect.StructValue); ok {
                                b0.Reset()
                                ForEachField(s.Interface(), f0) // TODO: errors?
                                if 0 < b0.Len() {
                                        fmt.Fprintf(msg, "<%s>%s</%s>", name, b0, name)
                                }
                                return true
                        } else {
                                v = p.Elem()
                        }
                }

                b0.Reset()
                nxt = f0(t, v)
                fmt.Fprint(msg, b0)
                return
        }
        ForEachField(call, f)
        ForEachField(ncall, f)
        fmt.Fprintf(msg, "</%s>", name)

        //fmt.Printf("finding: %s\n", msg)
        return io.ReadWriter(msg), msg.Len()
}

func (call *eBayFindingService_GetVersion)                      GetHeaders(app *App) map[string] string { return call.getHeaders(call, app) }
func (call *eBayFindingService_GetHistograms)                   GetHeaders(app *App) map[string] string { return call.getHeaders(call, app) }
func (call *eBayFindingService_GetSearchKeywordsRecommendation) GetHeaders(app *App) map[string] string { return call.getHeaders(call, app) }
func (call *eBayFindingService_FindItemsByKeywords)             GetHeaders(app *App) map[string] string { return call.getHeaders(call, app) }
func (call *eBayFindingService_FindItemsAdvanced)               GetHeaders(app *App) map[string] string { return call.getHeaders(call, app) }

func (call *eBayFindingService_GetVersion)                      GetMessage(app *App) (io.Reader, int) { return call.newMessage(call) }
func (call *eBayFindingService_GetHistograms)                   GetMessage(app *App) (io.Reader, int) { return call.newMessage(call) }
func (call *eBayFindingService_GetSearchKeywordsRecommendation) GetMessage(app *App) (io.Reader, int) { return call.newMessage(call) }
func (call *eBayFindingService_FindItemsByKeywords)             GetMessage(app *App) (io.Reader, int) { return call.newMessage(call) }
func (call *eBayFindingService_FindItemsAdvanced)               GetMessage(app *App) (io.Reader, int) { return call.newMessage(call) }

// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsByCategory.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsByProduct.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsIneBayStores.html
// TODO: http://developer.ebay.com/DevZone/finding/CallRef/findItemsAdvanced.html

func (svc *FindingService) NewGetVersionCall() (call *eBayFindingService_GetVersion)                    { return &eBayFindingService_GetVersion{} }
func (svc *FindingService) NewGetHistogramsCall() (call *eBayFindingService_GetHistograms)              { return &eBayFindingService_GetHistograms{} }
func (svc *FindingService) NewGetSearchKeywordsRecommendationCall() (call *eBayFindingService_GetSearchKeywordsRecommendation) { return &eBayFindingService_GetSearchKeywordsRecommendation{} }
func (svc *FindingService) NewFindItemsByKeywordsCall() (call *eBayFindingService_FindItemsByKeywords)  { return &eBayFindingService_FindItemsByKeywords{} }
func (svc *FindingService) NewFindItemsAdvancedCall() (call *eBayFindingService_FindItemsAdvanced)      { return &eBayFindingService_FindItemsAdvanced{} }

// ParseResponse parse text response of an eBay operation.
func (svc *FindingService) ParseResponse(str string) (res *FindItemsResponse, err os.Error) {
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
func (svc *FindingService) parseXMLResponse(str string) (res *FindItemsResponse, err os.Error) {
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

        res = new(FindItemsResponse)

        var v interface{}
        switch start.Name.Local {
        case "findItemsByCategoryResponse":
                v = &struct {
                        XMLName xml.Name "findItemsByCategoryResponse"
                        *FindItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsByProductResponse":
                v = &struct {
                        XMLName xml.Name "findItemsByProductResponse"
                        *FindItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsIneBayStoresResponse":
                v = &struct {
                        XMLName xml.Name "findItemsIneBayStoresResponse"
                        *FindItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsByKeywordsResponse":
                v = &struct {
                        XMLName xml.Name "findItemsByKeywordsResponse"
                        *FindItemsResponse
                }{ xml.Name{}, res, }
        case "findItemsAdvancedResponse":
                v = &struct {
                        XMLName xml.Name "findItemsAdvancedResponse"
                        *FindItemsResponse
                }{ xml.Name{}, res, }
                // TODO: more response
        }

        if v == nil {
                err = os.NewError(fmt.Sprintf("bad response: %s",start.Name.Local))
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
func (svc *FindingService) parseJSONResponse(str string) (res *FindItemsResponse, err os.Error) {
        ra := make([]*findItemsResponseJSON, 1)

        var v interface{}
        switch t := getJSONResponseType(str[0:]); t {
        case "findItemsByCategoryResponse":
                v = &struct { V []*findItemsResponseJSON "findItemsByCategoryResponse" }{ ra }
        case "findItemsByProductResponse":
                v = &struct { V []*findItemsResponseJSON "findItemsByProductResponse" }{ ra }
        case "findItemsIneBayStoresResponse":
                v = &struct { V []*findItemsResponseJSON "findItemsIneBayStoresResponse" }{ ra }
        case "findItemsByKeywordsResponse":
                v = &struct { V []*findItemsResponseJSON "findItemsByKeywordsResponse" }{ ra }
        case "findItemsAdvancedResponse":
                v = &struct { V []*findItemsResponseJSON "findItemsAdvancedResponse" }{ ra }
        default:
                err = os.NewError("unknown JSON response: '"+t+"'")
                return
        }

        err = json.Unmarshal([]byte(str), v)
        if err == nil { res, err = noJSON(ra[0]) }
        return
}

func noJSON(r *findItemsResponseJSON) (res *FindItemsResponse, err os.Error) {
        res = &FindItemsResponse{}
        res.SearchResult.Item = make([]Item, len(r.SearchResult[0].Item))
        err = RoughAssign(res, r)
        return
}
