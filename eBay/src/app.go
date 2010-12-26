package eBay

import (
        "os"
        "io"
        "io/ioutil"
        "http"
        "bytes"
        "xml"
        "json"
        "strconv"
        "strings"
        "fmt"
)

type App struct {
        DEVID string
        AppID string
        CertID string

        ServiceVersion string
        GlobalID string
        ResponseFormat string
}

type eBayServiceCall interface {
        GetURL(app *App) string
        GetMessage(app *App) io.Reader
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
        ResponseFormat: "XML",
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

// NewFindingService returns a new eBay.FindingService.
func (eb *App) NewFindingService() (p *FindingService) {
        p = &FindingService{ eb }
        return
}

// ParseResponse parse text response of an eBay operation.
func (eb *App) ParseResponse(str string) (res *findItemsResponse, err os.Error) {
        switch eb.ResponseFormat {
        case "JSON":
                res, err = eb.parseJSONResponse(str)
        case "XML":
                res, err = eb.parseXMLResponse(str)
        default:
                err = os.NewError("unknown data format '"+eb.ResponseFormat+"'")
        }
        return
}

// parseXMLResponse parse XML format response
func (eb *App) parseXMLResponse(str string) (res *findItemsResponse, err os.Error) {
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
func (eb *App) parseJSONResponse(str string) (res *findItemsResponse, err os.Error) {
        ra := make([]*findItemsJSONResponse, 1)

        var v interface{}
        switch t := getJSONResponseType(str[0:]); t {
        case "findItemsByCategoryResponse":
                v = &struct {
                        V []*findItemsJSONResponse "findItemsByCategoryResponse"
                }{ ra }
        case "findItemsByProductResponse":
                v = &struct {
                        V []*findItemsJSONResponse "findItemsByProductResponse"
                }{ ra }
        case "findItemsIneBayStoresResponse":
                v = &struct {
                        V []*findItemsJSONResponse "findItemsIneBayStoresResponse"
                }{ ra }
        case "findItemsByKeywordsResponse":
                v = &struct {
                        V []*findItemsJSONResponse "findItemsByKeywordsResponse"
                }{ ra }
        case "findItemsAdvancedResponse":
                v = &struct {
                        V []*findItemsJSONResponse "findItemsAdvancedResponse"
                }{ ra }
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
                        CategoryId: i.PrimaryCategory[0].CategoryId[0],
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
