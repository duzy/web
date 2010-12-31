package eBay

import (
        "../_obj/web"
        "fmt"
        "os"
)

const (
        SQL_CREATE_CACHE_CATEGORY_TABLE = `
CREATE TABLE IF NOT EXISTS table_eBay_cache_categories(
  categoryId VARCHAR(32) PRIMARY KEY,
  categoryName VARCHAR(128) NOT NULL
)
`
        SQL_INSERT_CACHE_CATEGORY_ROW = `
INSERT INTO table_eBay_cache_categories(categoryId, categoryName) VALUES(?,?)
  ON DUPLICATE KEY UPDATE categoryName=VALUES(categoryName)
`
        SQL_SELECT_CACHE_CATEGORY_ROW = `
SELECT categoryId, categoryName FROM table_eBay_cache_categories
  WHERE categoryId=? LIMIT 1
`
        SQL_CREATE_CACHE_ITEM_TABLE = `
CREATE TABLE IF NOT EXISTS table_eBay_cache_items(
  itemId VARCHAR(32) PRIMARY KEY,
  title VARCHAR(128) NOT NULL,
  primaryCategory VARCHAR(32) NOT NULL,
  galleryURL VARCHAR(256),
  galleryPlusPictureURL VARCHAR(512),
  viewItemURL VARCHAR(256),
  productId VARCHAR(32),
  paymentMethod VARCHAR(32),
  location VARCHAR(512),
  country VARCHAR(32),
  condition_ConditionId VARCHAR(32),
  condition_ConditionDisplayName VARCHAR(256),
  shippingInfo_ShippingServiceCost FLOAT,
  shippingInfo_ShippingServiceCost_CurrencyId CHAR(3),
  shippingInfo_ShippingType VARCHAR(32),
  shippingInfo_ShipToLocations VARCHAR(256),
  shippingInfo_HandlingTime SMALLINT,
  shippingInfo_ExpeditedShipping BIT(1),
  shippingInfo_OneDayShippingAvailable BIT(1),
  sellingStatus_CurrentPrice FLOAT,
  sellingStatus_CurrentPrice_CurrencyId CHAR(3),
  sellingStatus_ConvertedCurrentPrice FLOAT,
  sellingStatus_ConvertedCurrentPrice_CurrencyId CHAR(3),
  sellingStatus_BidCount INT,
  sellingStatus_SellingState VARCHAR(256),
  sellingStatus_TimeLeft VARCHAR(24),
  listingInfo_StartTime VARCHAR(30),
  listingInfo_EndTime VARCHAR(30),
  listingInfo_ListingType VARCHAR(32),
  listingInfo_BestOfferEnabled BIT(1),
  listingInfo_BuyItNowAvailable BIT(1),
  listingInfo_Gift BIT(1),
  returnsAccepted BIT(1),
  autoPay BIT(1)
)
`
        SQL_INSERT_CACHE_ITEM_ROW = `
INSERT INTO table_eBay_cache_items(
  itemId,
  title,
  primaryCategory,
  galleryURL,
  galleryPlusPictureURL,
  viewItemURL,
  productId,
  paymentMethod,
  location,
  country,
  condition_ConditionId,
  condition_ConditionDisplayName,
  shippingInfo_ShippingServiceCost,
  shippingInfo_ShippingServiceCost_CurrencyId,
  shippingInfo_ShippingType,
  shippingInfo_ShipToLocations,
  shippingInfo_HandlingTime,
  shippingInfo_ExpeditedShipping,
  shippingInfo_OneDayShippingAvailable,
  sellingStatus_CurrentPrice,
  sellingStatus_CurrentPrice_CurrencyId,
  sellingStatus_ConvertedCurrentPrice,
  sellingStatus_ConvertedCurrentPrice_CurrencyId,
  sellingStatus_BidCount,
  sellingStatus_SellingState,
  sellingStatus_TimeLeft,
  listingInfo_StartTime,
  listingInfo_EndTime,
  listingInfo_ListingType,
  listingInfo_BestOfferEnabled,
  listingInfo_BuyItNowAvailable,
  listingInfo_Gift,
  returnsAccepted,
  autoPay
) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
  ON DUPLICATE KEY UPDATE 
    title=VALUES(title),
    primaryCategory=VALUES(primaryCategory),
    galleryURL=VALUES(galleryURL),
    galleryPlusPictureURL=VALUES(galleryPlusPictureURL),
    viewItemURL=VALUES(viewItemURL),
    productId=VALUES(productId),
    paymentMethod=VALUES(paymentMethod),
    location=VALUES(location),
    country=VALUES(country),
    condition_ConditionId=VALUES(condition_ConditionId),
    condition_ConditionDisplayName=VALUES(condition_ConditionDisplayName),
    shippingInfo_ShippingServiceCost=VALUES(shippingInfo_ShippingServiceCost),
    shippingInfo_ShippingServiceCost_CurrencyId=VALUES(shippingInfo_ShippingServiceCost_CurrencyId),
    shippingInfo_ShippingType=VALUES(shippingInfo_ShippingType),
    shippingInfo_ShipToLocations=VALUES(shippingInfo_ShipToLocations),
    shippingInfo_HandlingTime=VALUES(shippingInfo_HandlingTime),
    shippingInfo_ExpeditedShipping=VALUES(shippingInfo_ExpeditedShipping),
    shippingInfo_OneDayShippingAvailable=VALUES(shippingInfo_OneDayShippingAvailable),
    sellingStatus_CurrentPrice=VALUES(sellingStatus_CurrentPrice),
    sellingStatus_CurrentPrice_CurrencyId=VALUES(sellingStatus_CurrentPrice_CurrencyId),
    sellingStatus_ConvertedCurrentPrice=VALUES(sellingStatus_ConvertedCurrentPrice),
    sellingStatus_ConvertedCurrentPrice_CurrencyId=VALUES(sellingStatus_ConvertedCurrentPrice_CurrencyId),
    sellingStatus_BidCount=VALUES(sellingStatus_BidCount),
    sellingStatus_SellingState=VALUES(sellingStatus_SellingState),
    sellingStatus_TimeLeft=VALUES(sellingStatus_TimeLeft),
    listingInfo_StartTime=VALUES(listingInfo_StartTime),
    listingInfo_EndTime=VALUES(listingInfo_EndTime),
    listingInfo_ListingType=VALUES(listingInfo_ListingType),
    listingInfo_BestOfferEnabled=VALUES(listingInfo_BestOfferEnabled),
    listingInfo_BuyItNowAvailable=VALUES(listingInfo_BuyItNowAvailable),
    listingInfo_Gift=VALUES(listingInfo_Gift),
    returnsAccepted=VALUES(returnsAccepted),
    autoPay=VALUES(autoPay)
`
        SQL_SELECT_CACHE_ITEM_ROW = `
SELECT
  title,
  primaryCategory,
  galleryURL,
  galleryPlusPictureURL,
  viewItemURL,
  productId,
  paymentMethod,
  location,
  country,
  condition_ConditionId,
  condition_ConditionDisplayName,
  shippingInfo_ShippingServiceCost,
  shippingInfo_ShippingServiceCost_CurrencyId,
  shippingInfo_ShippingType,
  shippingInfo_ShipToLocations,
  shippingInfo_HandlingTime,
  shippingInfo_ExpeditedShipping,
  shippingInfo_OneDayShippingAvailable,
  sellingStatus_CurrentPrice,
  sellingStatus_CurrentPrice_CurrencyId,
  sellingStatus_ConvertedCurrentPrice,
  sellingStatus_ConvertedCurrentPrice_CurrencyId,
  sellingStatus_BidCount,
  sellingStatus_SellingState,
  sellingStatus_TimeLeft,
  listingInfo_StartTime,
  listingInfo_EndTime,
  listingInfo_ListingType,
  listingInfo_BestOfferEnabled,
  listingInfo_BuyItNowAvailable,
  listingInfo_Gift,
  returnsAccepted,
  autoPay
  FROM table_eBay_cache_items
  WHERE itemId=? LIMIT 1
`
)

type dbCache struct {
        db web.Database
}

func _2string(v interface{}) string { return fmt.Sprintf("%v",v) }
func _2int(v interface{}) int { return v.(int) }
func _2float(v interface{}) float { return v.(float) }
func _2bool(v interface{}) bool { return v.(bool) }
func _bool2int(v bool) (n int) { if v { n = 1 } else { n = 0 }; return }

// NewDBCache accepts parameters in this fixed order:
//      host, user, password, database
func NewDBCache(params ...interface{}) (c Cacher, err os.Error) {
        a := []interface{}(params)
        cfg := &web.DatabaseConfig{
        Host: a[0].(string),
        User: a[1].(string),
        Password: a[2].(string),
        Database: a[3].(string),
        }
        dbm := web.GetDBManager()
        db, err := dbm.GetDatabase(cfg)
        if err == nil {
                err = createCacheTables(db)
                if err == nil {
                        dbc := &dbCache{ db }
                        c = Cacher(dbc);
                }
        }
        return
}

func createCacheTables(db web.Database) (err os.Error) {
        sql := SQL_CREATE_CACHE_CATEGORY_TABLE
        sql += ";\n"
        sql += SQL_CREATE_CACHE_ITEM_TABLE
        _, err = db.Query(sql)
        if err != nil { return }
        return
}

func (c *dbCache) Close() (err os.Error) {
        err = c.db.Close()
        return
}

func (c *dbCache) exec(sql string, params ...interface{}) (res web.QueryResult, err os.Error) {
        var stmt web.SQLStatement
        if stmt, err = c.db.Prepare(sql); err != nil { return }

        defer stmt.Close()

        if res, err = stmt.Execute(params...); err != nil { return }
        return
}

func (c *dbCache) get(sql string, params ...interface{}) (row []interface{}, err os.Error) {
        res, err := c.exec(sql, params...)
        if err != nil { return }
        if res.GetRowCount() <= 0 {
                err = os.NewError("not found")
                return
        }

        row, err = res.FetchRow()
        if err != nil { return }
        //if row == nil { err = os.NewError("fatal: FetchRow") }
        return
}

func (c *dbCache) CacheCategory(cat *Category) (err os.Error) {
        res, err := c.exec(SQL_INSERT_CACHE_CATEGORY_ROW,
                cat.CategoryId,
                cat.CategoryName )
        if err != nil { return }
        if res.GetAffectedRows() == 0 /*!= 1*/ {
                err = os.NewError(fmt.Sprintf("%d rows affected", res.GetAffectedRows()))
        }
        return
}

func (c *dbCache) CacheItem(i *Item) (err os.Error) {
        res, err := c.exec(SQL_INSERT_CACHE_ITEM_ROW,
                i.ItemId,
                i.Title,
                i.PrimaryCategory.CategoryId,
                i.GalleryURL,
                i.GalleryPlusPictureURL,
                i.ViewItemURL,
                i.ProductId,
                i.PaymentMethod,
                i.Location,
                i.Country,
                i.Condition.ConditionId,
                i.Condition.ConditionDisplayName,
                i.ShippingInfo.ShippingServiceCost.Amount,
                i.ShippingInfo.ShippingServiceCost.CurrencyId,
                i.ShippingInfo.ShippingType,
                i.ShippingInfo.ShipToLocations,
                i.ShippingInfo.HandlingTime,
                _bool2int(i.ShippingInfo.ExpeditedShipping),
                _bool2int(i.ShippingInfo.OneDayShippingAvailable),
                i.SellingStatus.CurrentPrice.Amount,
                i.SellingStatus.CurrentPrice.CurrencyId,
                i.SellingStatus.ConvertedCurrentPrice.Amount,
                i.SellingStatus.ConvertedCurrentPrice.CurrencyId,
                i.SellingStatus.BidCount,
                i.SellingStatus.SellingState,
                i.SellingStatus.TimeLeft,
                i.ListingInfo.StartTime,
                i.ListingInfo.EndTime,
                i.ListingInfo.ListingType,
                _bool2int(i.ListingInfo.BestOfferEnabled),
                _bool2int(i.ListingInfo.BuyItNowAvailable),
                _bool2int(i.ListingInfo.Gift),
                _bool2int(i.ReturnsAccepted),
                _bool2int(i.AutoPay) )
        if err != nil { return }
        if res.GetAffectedRows() == 0 /*!= 1*/ {
                err = os.NewError(fmt.Sprintf("%d rows affected", res.GetAffectedRows()))
        }
        return
}

func (c *dbCache) GetCategory(id string) (cat *Category, err os.Error) {
        row, err := c.get(SQL_SELECT_CACHE_CATEGORY_ROW, id)
        if err != nil { return }

        cat = &Category{
        CategoryId: _2string(row[0]),
        CategoryName: _2string(row[1]),
        }
        return
}

func (c *dbCache) GetItem(id string) (itm *Item, err os.Error) {
        row, err := c.get(SQL_SELECT_CACHE_ITEM_ROW, id)
        if err != nil { return }

        itm = &Item{
        ItemId: id,
        Title: _2string(row[0]),
        PrimaryCategory: Category{
                CategoryId: _2string(row[1]),
                CategoryName: "",
                },
        GalleryURL: _2string(row[2]),
        GalleryPlusPictureURL: _2string(row[3]),
        ViewItemURL: _2string(row[4]),
        ProductId: _2string(row[5]),
        PaymentMethod: _2string(row[6]),
        Location: _2string(row[7]),
        Country: _2string(row[8]),
        Condition: Condition{
                ConditionId: _2string(row[9]),
                ConditionDisplayName: _2string(row[10]),
                },
        ShippingInfo: ShippingInfo{
                ShippingServiceCost: Money{
                        Amount: _2float(row[11]),
                        CurrencyId: _2string(row[12]),
                        },
                ShippingType: _2string(row[13]),
                ShipToLocations: _2string(row[14]),
                HandlingTime: _2int(row[15]),
                ExpeditedShipping: _2bool(row[16]),
                OneDayShippingAvailable: _2bool(row[17]),
                },
                SellingStatus: SellingStatus{
                        CurrentPrice: Money{
                                Amount: _2float(row[18]),
                                CurrencyId: _2string(row[19]),
                                },
                        ConvertedCurrentPrice: Money{
                                Amount: _2float(row[20]),
                                CurrencyId: _2string(row[21]),
                                },
                        BidCount: _2int(row[22]),
                        SellingState: _2string(row[23]),
                        TimeLeft: _2string(row[24]),
                        },
                ListingInfo: ListingInfo{
                        StartTime: _2string(row[25]),
                        EndTime: _2string(row[26]),
                        ListingType: _2string(row[27]),
                        BestOfferEnabled: _2bool(row[28]),
                        BuyItNowAvailable: _2bool(row[29]),
                        Gift: _2bool(row[30]),
                        },
                ReturnsAccepted: _2bool(row[31]),
                AutoPay: _2bool(row[32]),
        }
        return
}
