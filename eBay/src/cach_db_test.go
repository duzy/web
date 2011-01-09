package eBay

import (
        "testing"
        "time"
        "fmt"
        "reflect"
)

func TestCacheCategoriesByLevel(t *testing.T) {
        c, err := NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { t.Errorf("NewDBCache: %v", err); return }

        defer c.Close()

        _, err = c.GetCategoriesByLevel(1)
        if err != nil {
                t.Errorf("GetCategoriesByLevel: %v", err);
                return
        }

        //fmt.Printf("cats: %v\n", cats)
}

func TestCacheCategory(t *testing.T) {
        c, err := NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { t.Errorf("NewDBCache: %v", err); return }

        defer c.Close()

        cat := &Category{
        CategoryID: "1234567890",
        CategoryLevel: 5,
        CategoryName: fmt.Sprintf("test category: %v", time.Nanoseconds()),
        CategoryParentID: "<none>",
        AutoPayEnabled: true,
        BestOfferEnabled: true,
        }
        err = c.CacheCategory(cat)
        if err != nil { t.Errorf("c.CacheCategory: %v", err); return }

        cat2, err := c.GetCategory(cat.CategoryID)
        if err != nil { t.Errorf("c.GetCategory(%s): %v", cat.CategoryID, err); return }
        if cat2 == nil { t.Errorf("c.GetCategory(%s): no category returned", cat.CategoryID); return }
        if cat.CategoryID != cat2.CategoryID { t.Errorf("wrong category: %v != %v", cat, cat2); return }
        if cat.CategoryLevel != cat2.CategoryLevel { t.Errorf("wrong category: %v != %v", cat, cat2); return }
        if cat.CategoryParentID != cat2.CategoryParentID { t.Errorf("wrong category: %v != %v", cat, cat2); return }
        if cat.AutoPayEnabled != cat2.AutoPayEnabled { t.Errorf("wrong category: %v != %v", cat, cat2); return }
        if cat.BestOfferEnabled != cat2.BestOfferEnabled { t.Errorf("wrong category: %v != %v", cat, cat2); return }
}

func TestCacheItem(t *testing.T) {
        c, err := NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { t.Errorf("NewDBCache: %v", err); return }

        defer c.Close()

        item := &Item{
        ItemId: "234567",
        Title: fmt.Sprintf("test-title: %v", time.Nanoseconds()),
        PrimaryCategory: Category{
                CategoryID: "95840",
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
        if err != nil { t.Errorf("c.GetItem(%s): %v", item.ItemId, err); return }
        if item.ItemId != item2.ItemId {
                t.Error("c.GetItem: not equal")
                t.Errorf("c.GetItem: %v", item)
                t.Errorf("c.GetItem: %v", item2)
        }
        if !reflect.DeepEqual(item, item2) {
                t.Error("c.GetItem: not equal")
                t.Errorf("c.GetItem: %v", item)
                t.Errorf("c.GetItem: %v", item2)
                return
        }
}
