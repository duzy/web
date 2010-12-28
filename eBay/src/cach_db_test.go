package eBay

import (
        "testing"
        "time"
        "fmt"
        "reflect"
)

func TestCacheCategory(t *testing.T) {
        c, err := NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { t.Errorf("NewDBCache: %v", err); return }

        cat := &Category{
        CategoryId: "123456",
        CategoryName: fmt.Sprintf("test category: %v", time.Nanoseconds()),
        }
        err = c.CacheCategory(cat)
        if err != nil { t.Errorf("c.CacheCategory: %v", err); return }

        cat2, err := c.GetCategory(cat.CategoryId)
        if err != nil { t.Errorf("c.GetCategory(%s): %v", cat.CategoryId, err); return }
        if cat2 == nil { t.Errorf("c.GetCategory(%s): no category returned", cat.CategoryId); return }
        if cat.CategoryId != cat2.CategoryId { t.Errorf("wrong category: %v != %v", cat, cat2); return }
        if cat.CategoryName != cat2.CategoryName { t.Errorf("wrong category: %v != %v", cat, cat2); return }
}

func TestCacheItem(t *testing.T) {
        c, err := NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { t.Errorf("NewDBCache: %v", err); return }

        item := &Item{
        ItemId: "234567",
        Title: "test-title",
        PrimaryCategory: Category{
                CategoryId: "95840",
                CategoryName: "test-category",
                },
        GalleryURL: "gallery-url",
        GalleryPlusPictureURL: "gallery-plus-picture-url",
        ViewItemURL: "view-item-url",
        ProductId: "product-id",
        PaymentMethod: "payment-method",
        Location: "location",
        Country: "country",
        Condition: Condition{
                ConditionId: "condition-id",
                ConditionDisplayName: "condition-display-name",
                },
        ShippingInfo: ShippingInfo{
                ShippingServiceCost: Money{
                        Amount: 1000.50,
                        CurrencyId: "USD",
                        },
                ShippingType: "shipping-type",
                ShipToLocations: "ship-to-locations",
                HandlingTime: 100,
                ExpeditedShipping: false,
                OneDayShippingAvailable: true,
                },
                SellingStatus: SellingStatus{
                        CurrentPrice: Money{
                                Amount: 2000.70,
                                CurrencyId: "EUR",
                                },
                        ConvertedCurrentPrice: Money{
                                Amount: 3000.90,
                                CurrencyId: "USD",
                                },
                        BidCount: 200,
                        SellingState: "selling-state",
                        TimeLeft: "time-left",
                        },
                ListingInfo: ListingInfo{
                        StartTime: "start-time",
                        EndTime: "end-time",
                        ListingType: "listing-type",
                        BestOfferEnabled: true,
                        BuyItNowAvailable: false,
                        Gift: true,
                        },
                ReturnsAccepted: false,
                AutoPay: true,
        }

        err = c.CacheItem(item)
        if err != nil { t.Errorf("c.CacheItem: %v", err); return }

        item2, err := c.GetItem(item.ItemId)
        if err != nil { t.Errorf("c.GetItem: %v", err); return }
        if !reflect.DeepEqual(item, item2) {
                t.Errorf("c.GetItem: not equal: %v != %v", item, item2)
                return
        }
}
