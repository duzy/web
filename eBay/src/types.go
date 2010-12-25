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



type Money struct {
        CurrencyId string // USD, EUR, etc...
        Amount float
}

type ShippingInfo struct {
        ServiceCost Money
        Type string // Flat, Air ...
        ShipToLocations string
        ExpeditedShipping bool
        OneDayShippingAvailable bool
        HandlingTime int // The number of days it will take the seller to ship this item
}

type SellingStatus struct {
}

type ListingInfo struct {
}

type Condition struct {
}

type Category struct {
        Id string
        Name string
}

type Item struct {
        Id string
        Title string
        PrimaryCategory *Category
        GalleryURL string
        ViewItemURL string
        ProductId string
        PaymentMethod string
        Location string
        Country string
        ShippingInfo ShippingInfo
        SellingStatus SellingStatus
        ListingInfo ListingInfo
        ReturnsAccepted bool
        GalleryPlusPictureURL string
        Condition Condition
}
