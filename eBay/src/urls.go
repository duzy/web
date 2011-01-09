package eBay

const (
        URL_eBayShoppingSandbox         = "http://open.api.sandbox.ebay.com/shopping"
        URL_eBayShoppingProduction      = "http://open.api.ebay.com/shopping"
        URL_eBayTradingSandbox          = "https://api.sandbox.ebay.com/ws/api.dll"
        URL_eBayTradingProduction       = "https://api.ebay.com/ws/api.dll"
        URL_eBayMerchandisingSandbox    = "http://svcs.sandbox.ebay.com/MerchandisingService"
        URL_eBayMerchandisingProduction = "http://svcs.ebay.com/MerchandisingService"

        /*
        URL_eBayTrading         = URL_eBayTradingSandbox
        URL_eBayShopping        = URL_eBayShoppingSandbox
        URL_eBayMerchandising   = URL_eBayMerchandisingSandbox
         */
        URL_eBayTrading         = URL_eBayTradingProduction
        URL_eBayShopping        = URL_eBayShoppingProduction
        URL_eBayMerchandising   = URL_eBayMerchandisingProduction

        // APIs only supported in Production
        URL_eBayFinding   = "http://svcs.ebay.com/services/search/FindingService/v1"
        URL_eBayBestMatch = "https://svcs.ebay.com/services/search/BestMatchItemDetailsService/v1"
        URL_eBayFeedback  = "https://svcs.ebay.com/FeedbackService"
)
