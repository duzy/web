package main

import (
        "./_obj/eBay"
        //"./_obj/web"
        "log"
        "fmt"
)

func main() {
        a := eBay.NewApp(true)
        trading := a.NewTradingService()

        c := trading.NewGetCategoriesCall()
        c.ViewAllNodes = true
        c.DetailLevel = eBay.ReturnAll //"ReturnAll"
        //c.CategorySiteID = ""

        fmt.Printf("fetching categories from eBay...")
        xml, err := a.Invoke(c)
        if err != nil {
                fmt.Printf("\n")
                log.Printf("error: %v", err)
                return
        }
        fmt.Printf(" Done.\n")

        fmt.Printf("parsing categories...")
        resp := new(eBay.GetCategoriesResponse)
        err = a.ParseXMLResponse(resp, xml)
        if err != nil {
                fmt.Printf("\n")
                log.Printf("error: %v", err)
                return
        }
        fmt.Printf(" Done.\n")

        fmt.Printf("Ack: %v\n", resp.Ack)
        fmt.Printf("Version: %v\n", resp.Version)
        fmt.Printf("Timestamp: %v\n", resp.Timestamp)
        fmt.Printf("Build: %v\n", resp.Build)
        fmt.Printf("CategoryCount: %v (%d parsed)\n", resp.CategoryCount, len(resp.CategoryArray.Category))
        fmt.Printf("CategoryVersion: %v\n", resp.CategoryVersion)
        fmt.Printf("UpdateTime: %v\n", resp.UpdateTime)
        fmt.Printf("ReservePriceAllowed: %v\n", resp.ReservePriceAllowed)
        fmt.Printf("MinimumReservePrice: %v\n", resp.MinimumReservePrice)

        // SQL: admin may copy table like this
        // CREATE TABLE categories IF NOT EXISTS
        //     SELECT * FROM DownloadedCategories

        cache, err := eBay.NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil {
                log.Printf("NewDBCache: %v", err);
                return
        }

        defer cache.Close()

        fmt.Printf("saving categories to database...")
        for i := 0; i < len(resp.CategoryArray.Category); i += 1 {
                err = cache.CacheCategory(&(resp.CategoryArray.Category[i]))
                if err != nil {
                        fmt.Printf("\n")
                        log.Printf("CacheCategory: %v", err);
                }
        }
        fmt.Printf("Done.\n")
}

