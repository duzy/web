package main

import (
        "./_obj/eBay"
        //"./_obj/web"
        "log"
)

func main() {
        a := eBay.NewApp(false)
        trading := a.NewTradingService()

        c := trading.NewGetCategoriesCall()
        c.ViewAllNodes = true
        c.DetailLevel = eBay.ReturnAll //"ReturnAll"
        //c.CategorySiteID = ""

        _, err := a.Invoke(c)
        if err != nil {
                log.Printf("error: %v", err)
                return
        }

        
}

