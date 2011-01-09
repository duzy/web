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

/*
 FeedbackRatingStar values:
        None            - No graphic displayed, feedback score 0-9.
        Yellow          - Yellow Star, feedback score 10-49.
        Blue            - Blue Star, feedback score 50-99.
        Turquoise       - Turquoise Star, feedback score 100-499.
        Purple          - Purple Star, feedback score 500-999.
        Red             - Red Star, feedback score 1,000-4,999.
        Green           - Green Star, feedback score 5,000-9,999.
        YellowShooting  - Yellow Shooting Star, feedback score 10,000-24,999.
        TurquoiseShooting - Turquoise Shooting Star, feedback score 25,000-49,999.
        PurpleShooting  - Purple Shooting Star, feedback score 50,000-99,999.
        RedShooting     - Red Shooting Star, feedback score 100,000-499,000 and above.
        GreenShooting   - Green Shooting Star, feedback score 500,000-999,000 and above.
        SilverShooting  - Silver Shooting Star, feedback score 1,000,000 or more.
 */
type SellerInfo struct {
        FeedbackRatingStar string
        FeedbackScore int
        PositiveFeedbackPercent float
        SellerUserName string
        TopRatedSeller bool
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

type Storefront struct {
        StoreName string
        StoreURL string
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
        SecondaryCategory Category
        GalleryURL string
        GalleryPlusPictureURL string
        ViewItemURL string
        ProductId string
        PaymentMethod string
        Location string
        Country string
        Condition Condition
        SellerInfo SellerInfo
        StoreInfo Storefront
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
        ItemSearchURL string
        PaginationOutput PaginationOutput

        ErrorMessage *struct{
                Error []struct {
                        Category string // see http://www.developer.ebay.com/DevZone/finding/CallRef/types/ErrorCategory.html
                        Domain string
                        ErrorId string
                        ExceptionId string
                        Message string
                        Serverity string
                        Subdomain string
                }
        }

        SearchResult struct { Item []Item }
}

type FindItemsByCategoryResponse FindItemsResponse
type FindItemsByKeywordsResponse FindItemsResponse

func (resp *FindItemsByCategoryResponse) MakeXMLParseBuffer() (interface{}, os.Error) {
        return &struct {
                XMLName xml.Name "findItemsByCategoryResponse"
                *FindItemsByCategoryResponse
        }{ xml.Name{}, resp, }, nil
}

func (resp *FindItemsByKeywordsResponse) MakeXMLParseBuffer() (interface{}, os.Error) {
        return &struct {
                XMLName xml.Name "findItemsByKeywordsResponse"
                *FindItemsByKeywordsResponse
        }{ xml.Name{}, resp, }, nil
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
