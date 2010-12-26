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
        //"fmt"
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

func (eb *App) NewFindingService() (p *FindingService) {
        p = &FindingService{ eb }
        return
}

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

func (eb *App) parseXMLResponse(str string) (res *findItemsResponse, err os.Error) {
        r := new(findItemsResponse)
        err = xml.Unmarshal(bytes.NewBufferString(str), r)
        if err == nil { res = r }
        return
}

func (eb *App) parseJSONResponse(str string) (res *findItemsResponse, err os.Error) {
        r := new(findItemsJSONResponse)
        err = json.Unmarshal([]byte(str), r)
        if err == nil {
                res = noJSON(r)
        }
        return
}

func stoi(s string) (i int) { i, _ = strconv.Atoi(s); return }
func stof(s string) (f float) { f, _ = strconv.Atof(s); return }
func stob(s string) (b bool) { b, _ = strconv.Atob(s); return }
func noJSON(r *findItemsJSONResponse) (res *findItemsResponse) {
        res = &findItemsResponse{
        Ack: r.V[0].Ack[0],
        Version: r.V[0].Version[0],
        Timestamp: r.V[0].Timestamp[0],
        //SearchResult: { make([]Item, len(r.V[0].SearchResult[0].Item)) },
        ItemSearchURL: r.V[0].ItemSearchURL[0],
        PaginationOutput: PaginationOutput{
                PageNumber: stoi(r.V[0].PaginationOutput[0].PageNumber[0]),
                EntriesPerPage: stoi(r.V[0].PaginationOutput[0].EntriesPerPage[0]),
                TotalPages: stoi(r.V[0].PaginationOutput[0].TotalPages[0]),
                TotalEntries: stoi(r.V[0].PaginationOutput[0].TotalEntries[0]),
                },
        }

        res.SearchResult.Item = make([]Item, len(r.V[0].SearchResult[0].Item))
        
        for n, i := range r.V[0].SearchResult[0].Item {
                res.SearchResult.Item[n] = Item{
                ItemId: i.ItemId[0],
                Title: i.Title[0],
                //GlobalId: i.GlobalId[0],
                PrimaryCategory: Category{
                        Id: i.PrimaryCategory[0].CategoryId[0],
                        Name: i.PrimaryCategory[0].CategoryName[0],
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
                        Id: i.Condition[0].ConditionId[0],
                        DisplayName: i.Condition[0].ConditionDisplayName[0],
                        },
                }
        }//for (items)
        return
}
