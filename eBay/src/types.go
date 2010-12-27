package eBay

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
        CategoryId string
        CategoryName string
}

type Item struct {
        ItemId string
        Title string
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

type findItemsResponse struct {
        Ack string
        Version string
        Timestamp string
        SearchResult struct { Item []Item }
        ItemSearchURL string
        PaginationOutput PaginationOutput
}

type findItemsJSONResponse struct {
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
