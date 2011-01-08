package eBay

import (
        "os"
        "xml"
)

// See http://developer.ebay.com/DevZone/XML/docs/Reference/eBay/types/DetailLevelCodeType.html
// DetailLevelCodeType
const (
        ItemReturnAttributes = "ItemReturnAttributes"
        ItemReturnCategories = "ItemReturnCategories"
        ItemReturnDescription = "ItemReturnDescription"
        ReturnAll = "ReturnAll"
        ReturnHeaders = "ReturnHeaders"
        ReturnMessages = "ReturnMessages"
        ReturnSummary = "ReturnSummary"
)

// see http://developer.ebay.com/DevZone/XML/docs/Reference/eBay/types/WarningLevelCodeType.html
//  WarningLevelCodeType
const (
        Low = "Low"
        High = "High"
)

type Affiliate struct {
        CustomId string
        NetworkId string
        TrackingId string
}

type PaginationInput struct {
        EntriesPerPage int
        PageNumber int
}

type PaginationOutput struct {
        PageNumber int
        EntriesPerPage int
        TotalPages int
        TotalEntries int
}

type Money struct {
        Amount float "chardata"
        CurrencyId string "attr" // USD, EUR, etc...
}

type ShippingInfo struct {
        ShippingServiceCost Money
        ShippingType string // Flat, Air ...
        ShipToLocations string
        HandlingTime int // The number of days it will take the seller to ship this item
        ExpeditedShipping bool
        OneDayShippingAvailable bool
}

type SellingStatus struct {
        CurrentPrice Money
        ConvertedCurrentPrice Money
        SellingState string
        TimeLeft string
        BidCount int
}

type ListingInfo struct {
        StartTime string
        EndTime string
        ListingType string
        BestOfferEnabled bool
        BuyItNowAvailable bool
        Gift bool
}

type Condition struct {
        ConditionId string
        ConditionDisplayName string
}

type Category struct {
        CategoryID string
        CategoryLevel int
        CategoryName string
        CategoryParentID string
        AutoPayEnabled bool
        BestOfferEnabled bool
}

type Item struct {
        // NOTE: Do NOT change the order, order must the same as cach_db SQLs
        ItemId string
        Title string
        GlobalId string
        PrimaryCategory Category
        GalleryURL string
        GalleryPlusPictureURL string
        ViewItemURL string
        ProductId string
        PaymentMethod string
        Location string
        Country string
        Condition Condition
        ShippingInfo ShippingInfo
        SellingStatus SellingStatus
        ListingInfo ListingInfo
        ReturnsAccepted bool
        AutoPay bool
}

// TODO: ...
type FindItemsResponse struct {
        Ack string
        Version string
        Timestamp string
        SearchResult struct { Item []Item }
        ItemSearchURL string
        PaginationOutput PaginationOutput
}

type GetCategoriesResponse struct {
        Ack string
        Version string
        Timestamp string
        Build string

        UpdateTime string
        CategoryCount int
        CategoryVersion int
        CategoryArray struct{ Category []Category }
        ReservePriceAllowed bool
        MinimumReservePrice float
}

func (resp *GetCategoriesResponse) MakeXMLParseBuffer() (interface{}, os.Error) {
        return &struct {
                XMLName xml.Name "GetCategoriesResponse"
                *GetCategoriesResponse
        }{ xml.Name{}, resp, }, nil
}
