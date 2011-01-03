package eBay

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
        CurrencyId string "attr" // USD, EUR, etc...
        Amount float "chardata"
}

type ShippingInfo struct {
        ShippingServiceCost Money
        ShippingType string // Flat, Air ...
        ShipToLocations string
        ExpeditedShipping bool
        OneDayShippingAvailable bool
        HandlingTime int // The number of days it will take the seller to ship this item
}

type SellingStatus struct {
        CurrentPrice Money
        ConvertedCurrentPrice Money
        BidCount int
        SellingState string
        TimeLeft string
}

type ListingInfo struct {
        BestOfferEnabled bool
        BuyItNowAvailable bool
        StartTime string
        EndTime string
        ListingType string
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
}

type Item struct {
        ItemId string
        Title string
        GlobalId string
        PrimaryCategory Category
        GalleryURL string
        GalleryPlusPictureURL string
        ViewItemURL string
        ProductId string
        PaymentMethod string
        AutoPay bool
        Location string
        Country string
        ShippingInfo ShippingInfo
        SellingStatus SellingStatus
        ListingInfo ListingInfo
        ReturnsAccepted bool
        Condition Condition
}

