package eBay

import (
        "../_obj/web"
        "fmt"
        "os"
        "strings"
        //"reflect"
)

const (
        SQL_CREATE_CACHE_CATEGORY_TABLE = `
CREATE TABLE IF NOT EXISTS table_eBay_cache_categories(
  CategoryID VARCHAR(10) PRIMARY KEY,
  CategoryLevel INT,
  CategoryName VARCHAR(30) NOT NULL,
  CategoryParentID VARCHAR(10),
  AutoPayEnabled TINYINT DEFAULT 1,
  BestOfferEnabled TINYINT DEFAULT 1
)
`
        SQL_INSERT_CACHE_CATEGORY_ROW = `
INSERT INTO table_eBay_cache_categories(
  CategoryID,
  CategoryLevel,
  CategoryName,
  CategoryParentID,
  AutoPayEnabled,
  BestOfferEnabled
)
VALUES(?,?,?,?,?,?)
ON DUPLICATE KEY UPDATE
  CategoryLevel=VALUES(CategoryLevel),
  CategoryName=VALUES(CategoryName),
  CategoryParentID=VALUES(CategoryParentID),
  AutoPayEnabled=VALUES(AutoPayEnabled),
  BestOfferEnabled=VALUES(BestOfferEnabled)
`
        SQL_SELECT_CACHE_CATEGORY_ROW = `
SELECT
  CategoryID,
  CategoryLevel,
  CategoryName,
  CategoryParentID,
  AutoPayEnabled,
  BestOfferEnabled
FROM table_eBay_cache_categories
WHERE CategoryID=? LIMIT 1
`
        SQL_SELECT_CACHE_CATEGORY_ROW_BY = `
SELECT
  CategoryID,
  CategoryLevel,
  CategoryName,
  CategoryParentID,
  AutoPayEnabled,
  BestOfferEnabled
FROM table_eBay_cache_categories
WHERE 
`
        SQL_CREATE_CACHE_ITEM_TABLE = `
CREATE TABLE IF NOT EXISTS table_eBay_cache_items(
  ItemId VARCHAR(19) PRIMARY KEY,
  Title VARCHAR(128) NOT NULL,
  GlobalId VARCHAR(16),
  PrimaryCategory$CategoryID VARCHAR(10) NOT NULL,
  PrimaryCategory$CategoryName VARCHAR(30) NOT NULL,
  SecondaryCategory$CategoryID VARCHAR(10) NOT NULL,
  SecondaryCategory$CategoryName VARCHAR(30) NOT NULL,
  GalleryURL VARCHAR(256),
  GalleryPlusPictureURL VARCHAR(512),
  ViewItemURL VARCHAR(256),
  ProductId VARCHAR(32),
  PaymentMethod VARCHAR(32),
  Location VARCHAR(512),
  Country VARCHAR(32),
  Condition$ConditionId VARCHAR(32),
  Condition$ConditionDisplayName VARCHAR(256),
  SellerInfo$FeedbackRatingStar VARCHAR(20),
  SellerInfo$FeedbackScore INT,
  SellerInfo$PositiveFeedbackPercent FLOAT,
  SellerInfo$SellerUserName VARCHAR(30),
  SellerInfo$TopRatedSeller TINYINT,
  StoreInfo$StoreName VARCHAR(200),
  StoreInfo$StoreURL VARCHAR(256),
  ShippingInfo$ShippingServiceCost$Amount FLOAT,
  ShippingInfo$ShippingServiceCost$CurrencyId CHAR(3),
  ShippingInfo$ShippingType VARCHAR(32),
  ShippingInfo$ShipToLocations VARCHAR(256),
  ShippingInfo$HandlingTime SMALLINT,
  ShippingInfo$ExpeditedShipping TINYINT,
  ShippingInfo$OneDayShippingAvailable TINYINT,
  SellingStatus$CurrentPrice$Amount FLOAT,
  SellingStatus$CurrentPrice$CurrencyId CHAR(3),
  SellingStatus$ConvertedCurrentPrice$Amount FLOAT,
  SellingStatus$ConvertedCurrentPrice$CurrencyId CHAR(3),
  SellingStatus$SellingState VARCHAR(256),
  SellingStatus$TimeLeft VARCHAR(24),
  SellingStatus$BidCount INT,
  ListingInfo$StartTime VARCHAR(30),
  ListingInfo$EndTime VARCHAR(30),
  ListingInfo$ListingType VARCHAR(32),
  ListingInfo$BestOfferEnabled TINYINT,
  ListingInfo$BuyItNowAvailable TINYINT,
  ListingInfo$Gift TINYINT,
  ReturnsAccepted TINYINT,
  AutoPay TINYINT
)
`
        SQL_INSERT_CACHE_ITEM_ROW = `
INSERT INTO table_eBay_cache_items(
  ItemId,
  Title,
  GlobalId,
  PrimaryCategory$CategoryID,
  PrimaryCategory$CategoryName,
  SecondaryCategory$CategoryID,
  SecondaryCategory$CategoryName,
  GalleryURL,
  GalleryPlusPictureURL,
  ViewItemURL,
  ProductId,
  PaymentMethod,
  Location,
  Country,
  Condition$ConditionId,
  Condition$ConditionDisplayName,
  SellerInfo$FeedbackRatingStar,
  SellerInfo$FeedbackScore,
  SellerInfo$PositiveFeedbackPercent,
  SellerInfo$SellerUserName,
  SellerInfo$TopRatedSeller,
  StoreInfo$StoreName,
  StoreInfo$StoreURL,
  ShippingInfo$ShippingServiceCost$Amount,
  ShippingInfo$ShippingServiceCost$CurrencyId,
  ShippingInfo$ShippingType,
  ShippingInfo$ShipToLocations,
  ShippingInfo$HandlingTime,
  ShippingInfo$ExpeditedShipping,
  ShippingInfo$OneDayShippingAvailable,
  SellingStatus$CurrentPrice$Amount,
  SellingStatus$CurrentPrice$CurrencyId,
  SellingStatus$ConvertedCurrentPrice$Amount,
  SellingStatus$ConvertedCurrentPrice$CurrencyId,
  SellingStatus$SellingState,
  SellingStatus$TimeLeft,
  SellingStatus$BidCount,
  ListingInfo$StartTime,
  ListingInfo$EndTime,
  ListingInfo$ListingType,
  ListingInfo$BestOfferEnabled,
  ListingInfo$BuyItNowAvailable,
  ListingInfo$Gift,
  ReturnsAccepted,
  AutoPay
)
VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
ON DUPLICATE KEY UPDATE 
  Title=VALUES(Title),
  GlobalId=VALUES(GlobalId),
  PrimaryCategory$CategoryID=VALUES(PrimaryCategory$CategoryID),
  PrimaryCategory$CategoryName=VALUES(PrimaryCategory$CategoryName),
  SecondaryCategory$CategoryID=VALUES(SecondaryCategory$CategoryID),
  SecondaryCategory$CategoryName=VALUES(SecondaryCategory$CategoryName),
  GalleryURL=VALUES(GalleryURL),
  GalleryPlusPictureURL=VALUES(GalleryPlusPictureURL),
  ViewItemURL=VALUES(ViewItemURL),
  ProductId=VALUES(ProductId),
  PaymentMethod=VALUES(PaymentMethod),
  Location=VALUES(Location),
  Country=VALUES(Country),
  Condition$ConditionId=VALUES(Condition$ConditionId),
  Condition$ConditionDisplayName=VALUES(Condition$ConditionDisplayName),
  SellerInfo$FeedbackRatingStar=VALUES(SellerInfo$FeedbackRatingStar),
  SellerInfo$FeedbackScore=VALUES(SellerInfo$FeedbackScore),
  SellerInfo$PositiveFeedbackPercent=VALUES(SellerInfo$PositiveFeedbackPercent),
  SellerInfo$SellerUserName=VALUES(SellerInfo$SellerUserName),
  SellerInfo$TopRatedSeller=VALUES(SellerInfo$TopRatedSeller),
  StoreInfo$StoreName=VALUES(StoreInfo$StoreName),
  StoreInfo$StoreURL=VALUES(StoreInfo$StoreURL),
  ShippingInfo$ShippingServiceCost$Amount=VALUES(ShippingInfo$ShippingServiceCost$Amount),
  ShippingInfo$ShippingServiceCost$CurrencyId=VALUES(ShippingInfo$ShippingServiceCost$CurrencyId),
  ShippingInfo$ShippingType=VALUES(ShippingInfo$ShippingType),
  ShippingInfo$ShipToLocations=VALUES(ShippingInfo$ShipToLocations),
  ShippingInfo$HandlingTime=VALUES(ShippingInfo$HandlingTime),
  ShippingInfo$ExpeditedShipping=VALUES(ShippingInfo$ExpeditedShipping),
  ShippingInfo$OneDayShippingAvailable=VALUES(ShippingInfo$OneDayShippingAvailable),
  SellingStatus$CurrentPrice$Amount=VALUES(SellingStatus$CurrentPrice$Amount),
  SellingStatus$CurrentPrice$CurrencyId=VALUES(SellingStatus$CurrentPrice$CurrencyId),
  SellingStatus$ConvertedCurrentPrice$Amount=VALUES(SellingStatus$ConvertedCurrentPrice$Amount),
  SellingStatus$ConvertedCurrentPrice$CurrencyId=VALUES(SellingStatus$ConvertedCurrentPrice$CurrencyId),
  SellingStatus$BidCount=VALUES(SellingStatus$BidCount),
  SellingStatus$SellingState=VALUES(SellingStatus$SellingState),
  SellingStatus$TimeLeft=VALUES(SellingStatus$TimeLeft),
  ListingInfo$StartTime=VALUES(ListingInfo$StartTime),
  ListingInfo$EndTime=VALUES(ListingInfo$EndTime),
  ListingInfo$ListingType=VALUES(ListingInfo$ListingType),
  ListingInfo$BestOfferEnabled=VALUES(ListingInfo$BestOfferEnabled),
  ListingInfo$BuyItNowAvailable=VALUES(ListingInfo$BuyItNowAvailable),
  ListingInfo$Gift=VALUES(ListingInfo$Gift),
  ReturnsAccepted=VALUES(ReturnsAccepted),
  AutoPay=VALUES(AutoPay)
`
        SQL_SELECT_CACHE_ITEM_ROW = `
SELECT
  Title,
  GlobalId,
  PrimaryCategory$CategoryID,
  PrimaryCategory$CategoryName,
  SecondaryCategory$CategoryID,
  SecondaryCategory$CategoryName,
  GalleryURL,
  GalleryPlusPictureURL,
  ViewItemURL,
  ProductId,
  PaymentMethod,
  Location,
  Country,
  Condition$ConditionId,
  Condition$ConditionDisplayName,
  SellerInfo$FeedbackRatingStar,
  SellerInfo$FeedbackScore,
  SellerInfo$PositiveFeedbackPercent,
  SellerInfo$SellerUserName,
  SellerInfo$TopRatedSeller,
  StoreInfo$StoreName,
  StoreInfo$StoreURL,
  ShippingInfo$ShippingServiceCost$Amount,
  ShippingInfo$ShippingServiceCost$CurrencyId,
  ShippingInfo$ShippingType,
  ShippingInfo$ShipToLocations,
  ShippingInfo$HandlingTime,
  ShippingInfo$ExpeditedShipping,
  ShippingInfo$OneDayShippingAvailable,
  SellingStatus$CurrentPrice$Amount,
  SellingStatus$CurrentPrice$CurrencyId,
  SellingStatus$ConvertedCurrentPrice$Amount,
  SellingStatus$ConvertedCurrentPrice$CurrencyId,
  SellingStatus$SellingState,
  SellingStatus$TimeLeft,
  SellingStatus$BidCount,
  ListingInfo$StartTime,
  ListingInfo$EndTime,
  ListingInfo$ListingType,
  ListingInfo$BestOfferEnabled,
  ListingInfo$BuyItNowAvailable,
  ListingInfo$Gift,
  ReturnsAccepted,
  AutoPay
FROM table_eBay_cache_items
WHERE itemId=? LIMIT 1
`
)

type dbCache struct {
        db web.Database
}

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
                dbc := &dbCache{ db }
                err = dbc.createCacheTables()
                if err == nil {
                        c = Cacher(dbc);
                }
        }
        return
}

func (c *dbCache) createCacheTables() (err os.Error) {
        sql := SQL_CREATE_CACHE_CATEGORY_TABLE
        sql += ";\n"
        sql += SQL_CREATE_CACHE_ITEM_TABLE
        _, err = c.exec(sql)
        if err != nil { return }
        return
}

func (c *dbCache) Close() (err os.Error) {
        err = c.db.Close()
        return
}

func (c *dbCache) exec_(sql string, params ...interface{}) (res web.QueryResult, err os.Error) {
        var stmt web.SQLStatement
        if stmt, err = c.db.Prepare(sql); err != nil { return }
        if stmt == nil { err = os.NewError("failed preparing statement"); return }

        defer stmt.Close()

        if res, err = stmt.Execute(params...); err != nil { return }
        return
}

func (c *dbCache) exec(sql string, params ...interface{}) (res web.QueryResult, err os.Error) {
        q := sql
        if n1, n2 := strings.Count(q, "?"), len(params); n1 != n2 {
                err = os.NewError(fmt.Sprintf("mismatched of SQL parameters: %d != %d", n1, n2))
                return
        }
        for _, a := range params {
                if v, ok := a.(bool); ok {
                        if v { a = 1 } else { a = 0 }
                }
                s := c.db.Escape(fmt.Sprintf("%v", a))
                q = strings.Replace(q, "?", `"` + s + `"`, 1)
        }
        res, err = c.db.Query(q)
        if err != nil && strings.Index(err.String(), "[2014]") != -1 {
                // reconnect and give another try
                if err = c.db.Reconnect(); err == nil {
                        res, err = c.db.Query(q)
                }
        }
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
        named, err := FieldsToArray(cat)
        if err != nil { return }
        //fmt.Printf("category: %v\n", nf)

        fields := GetValues(named)

        res, err := c.exec(SQL_INSERT_CACHE_CATEGORY_ROW, fields...)
        if err != nil { return }
        if res.GetAffectedRows() == 0 /*!= 1*/ {
                //err = os.NewError(fmt.Sprintf("%d rows affected", res.GetAffectedRows()))
        }
        return
}

func (c *dbCache) CacheItem(i *Item) (err os.Error) {
        named, err := FieldsToArrayFlat(i)
        if err != nil { return }

        //fields := GetValues(named)
        fields := make([]interface{}, 0, len(named))
        for i := 0; i < len(named); i += 1 {
                //fmt.Printf("%v\n", named[i].Name)
                switch named[i].Name {
                case "PrimaryCategory.CategoryLevel":
                case "PrimaryCategory.CategoryParentID":
                case "PrimaryCategory.AutoPayEnabled":
                case "PrimaryCategory.BestOfferEnabled":
                case "SecondaryCategory.CategoryLevel":
                case "SecondaryCategory.CategoryParentID":
                case "SecondaryCategory.AutoPayEnabled":
                case "SecondaryCategory.BestOfferEnabled":
                default:
                        fields = append(fields, named[i].Value)
                }
        }

        _, err = c.exec(SQL_INSERT_CACHE_ITEM_ROW, fields...)
        if err != nil { return }
        
        return
}

func (c *dbCache) GetCategory(id string) (cat *Category, err os.Error) {
        res, err := c.exec(SQL_SELECT_CACHE_CATEGORY_ROW, id)
        if err != nil { return }

        cat = &Category{}
        err = RoughAssignQueryResult(cat, res)
        return
}

// TODO: use category-filter to make projection
func (c *dbCache) GetCategoriesByLevel(level int) (cats []*Category, err os.Error) {
        res, err := c.exec(SQL_SELECT_CACHE_CATEGORY_ROW_BY + "CategoryLevel=? ORDER BY CategoryName", level)
        if err != nil { return }

        cats = make([]*Category, res.GetRowCount())

        for i := 0; i < len(cats); i += 1 {
                cats[i] = &Category{}
                err = RoughAssignQueryResult(cats[i], res)
                //if err != nil { return }
        }
        
        return
}

func (c *dbCache) GetItem(id string) (itm *Item, err os.Error) {
        res, err := c.exec(SQL_SELECT_CACHE_ITEM_ROW, id)
        if err != nil { return }

        itm = &Item{ ItemId: id }
        err = RoughAssignQueryResult(itm, res)
        return
}
