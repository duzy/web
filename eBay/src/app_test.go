package eBay

import (
        "testing"
        "bytes"
        "xml"
        "json"
        //"fmt"
)

func TestXMLUnmarshalFindItemsResponse(t *testing.T) {
        const x = `<?xml version='1.0' encoding='UTF-8'?><findItemsByKeywordsResponse xmlns="http://www.ebay.com/marketplace/search/v1/services"><ack>Success</ack><version>1.8.0</version><timestamp>2010-12-26T04:06:27.315Z</timestamp><searchResult count="3"><item><itemId>120663454616</itemId><title>Auth. NOKIA N8 Touchscreen Unlocked Silver 12mp Cam NEW</title><globalId>EBAY-ENCA</globalId><primaryCategory><categoryId>3312</categoryId><categoryName>Cell Phones &amp; Smartphones</categoryName></primaryCategory><galleryURL>http://thumbs1.ebaystatic.com/pict/1206634546168080_1.jpg</galleryURL><viewItemURL>http://cgi.ebay.com/Auth-NOKIA-N8-Touchscreen-Unlocked-Silver-12mp-Cam-NEW-/120663454616?pt=Cell_Phones</viewItemURL><productId type="ReferenceID">83383670</productId><paymentMethod>PayPal</paymentMethod><autoPay>false</autoPay><postalCode>V3S2W2</postalCode><location>Canada</location><country>CA</country><shippingInfo><shippingServiceCost currencyId="USD">5.99</shippingServiceCost><shippingType>Flat</shippingType><shipToLocations>Worldwide</shipToLocations><expeditedShipping>false</expeditedShipping><oneDayShippingAvailable>false</oneDayShippingAvailable><handlingTime>2</handlingTime></shippingInfo><sellingStatus><currentPrice currencyId="USD">450.0</currentPrice><convertedCurrentPrice currencyId="USD">450.0</convertedCurrentPrice><bidCount>2</bidCount><sellingState>Active</sellingState><timeLeft>P0DT14H34M2S</timeLeft></sellingStatus><listingInfo><bestOfferEnabled>false</bestOfferEnabled><buyItNowAvailable>false</buyItNowAvailable><startTime>2010-12-23T18:40:29.000Z</startTime><endTime>2010-12-26T18:40:29.000Z</endTime><listingType>Auction</listingType><gift>false</gift></listingInfo><returnsAccepted>true</returnsAccepted><condition><conditionId>1000</conditionId><conditionDisplayName>New</conditionDisplayName></condition></item><item><itemId>250743776675</itemId><title>Nokia N8 (Unlocked)</title><globalId>EBAY-US</globalId><primaryCategory><categoryId>3312</categoryId><categoryName>Cell Phones &amp; Smartphones</categoryName></primaryCategory><galleryURL>http://thumbs4.ebaystatic.com/pict/2507437766758080_1.jpg</galleryURL><viewItemURL>http://cgi.ebay.com/Nokia-N8-Unlocked-/250743776675?pt=Cell_Phones</viewItemURL><productId type="ReferenceID">83383670</productId><paymentMethod>PayPal</paymentMethod><autoPay>true</autoPay><postalCode>60714</postalCode><location>Niles,IL,USA</location><country>US</country><shippingInfo><shippingServiceCost currencyId="USD">20.0</shippingServiceCost><shippingType>FlatDomesticCalculatedInternational</shippingType><shipToLocations>Worldwide</shipToLocations><expeditedShipping>true</expeditedShipping><oneDayShippingAvailable>false</oneDayShippingAvailable><handlingTime>4</handlingTime></shippingInfo><sellingStatus><currentPrice currencyId="USD">425.0</currentPrice><convertedCurrentPrice currencyId="USD">425.0</convertedCurrentPrice><bidCount>3</bidCount><sellingState>Active</sellingState><timeLeft>P0DT14H46M54S</timeLeft></sellingStatus><listingInfo><bestOfferEnabled>false</bestOfferEnabled><buyItNowAvailable>false</buyItNowAvailable><startTime>2010-12-16T18:53:21.000Z</startTime><endTime>2010-12-26T18:53:21.000Z</endTime><listingType>Auction</listingType><gift>false</gift></listingInfo><returnsAccepted>false</returnsAccepted><condition><conditionId>1000</conditionId><conditionDisplayName>New</conditionDisplayName></condition></item><item><itemId>250745791164</itemId><title>New Factory Unlocked Nokia N8 GSM Phone  Silver White</title><globalId>EBAY-US</globalId><primaryCategory><categoryId>3312</categoryId><categoryName>Cell Phones &amp; Smartphones</categoryName></primaryCategory><galleryURL>http://thumbs1.ebaystatic.com/pict/2507457911648080_1.jpg</galleryURL><viewItemURL>http://cgi.ebay.com/New-Factory-Unlocked-Nokia-N8-GSM-Phone-Silver-White-/250745791164?pt=Cell_Phones</viewItemURL><productId type="ReferenceID">83383670</productId><paymentMethod>PayPal</paymentMethod><autoPay>true</autoPay><postalCode>60016</postalCode><location>Des Plaines,IL,USA</location><country>US</country><shippingInfo><shippingServiceCost currencyId="USD">0.0</shippingServiceCost><shippingType>Free</shippingType><shipToLocations>US</shipToLocations><expeditedShipping>true</expeditedShipping><oneDayShippingAvailable>false</oneDayShippingAvailable><handlingTime>4</handlingTime></shippingInfo><sellingStatus><currentPrice currencyId="USD">490.0</currentPrice><convertedCurrentPrice currencyId="USD">490.0</convertedCurrentPrice><bidCount>0</bidCount><sellingState>Active</sellingState><timeLeft>P0DT15H2M24S</timeLeft></sellingStatus><listingInfo><bestOfferEnabled>false</bestOfferEnabled><buyItNowAvailable>true</buyItNowAvailable><buyItNowPrice currencyId="USD">539.0</buyItNowPrice><convertedBuyItNowPrice currencyId="USD">539.0</convertedBuyItNowPrice><startTime>2010-12-21T19:08:51.000Z</startTime><endTime>2010-12-26T19:08:51.000Z</endTime><listingType>AuctionWithBIN</listingType><gift>false</gift></listingInfo><returnsAccepted>false</returnsAccepted><galleryPlusPictureURL>http://galleryplus.ebayimg.com/ws/web/250745791164_1_0_1.jpg</galleryPlusPictureURL><condition><conditionId>1000</conditionId><conditionDisplayName>New</conditionDisplayName></condition></item></searchResult><paginationOutput><pageNumber>1</pageNumber><entriesPerPage>3</entriesPerPage><totalPages>1579</totalPages><totalEntries>4735</totalEntries></paginationOutput><itemSearchURL>http://shop.ebay.com/i.html?_nkw=Nokia+N8&amp;_ddo=1&amp;_ipg=3&amp;_pgn=1</itemSearchURL></findItemsByKeywordsResponse>`
        const j = `{"findItemsByKeywordsResponse":[{"ack":["Success"],"version":["1.8.0"],"timestamp":["2010-12-26T06:17:22.587Z"],"searchResult":[{"@count":"3","item":[{"itemId":["120663454616"],"title":["Auth. NOKIA N8 Touchscreen Unlocked Silver 12mp Cam NEW"],"globalId":["EBAY-ENCA"],"primaryCategory":[{"categoryId":["3312"],"categoryName":["Cell Phones & Smartphones"]}],"galleryURL":["http:\/\/thumbs1.ebaystatic.com\/pict\/1206634546168080_1.jpg"],"viewItemURL":["http:\/\/cgi.ebay.com\/Auth-NOKIA-N8-Touchscreen-Unlocked-Silver-12mp-Cam-NEW-\/120663454616?pt=Cell_Phones"],"productId":[{"@type":"ReferenceID","__value__":"83383670"}],"paymentMethod":["PayPal"],"autoPay":["false"],"postalCode":["V3S2W2"],"location":["Canada"],"country":["CA"],"shippingInfo":[{"shippingServiceCost":[{"@currencyId":"USD","__value__":"5.99"}],"shippingType":["Flat"],"shipToLocations":["Worldwide"],"expeditedShipping":["false"],"oneDayShippingAvailable":["false"],"handlingTime":["2"]}],"sellingStatus":[{"currentPrice":[{"@currencyId":"USD","__value__":"450.0"}],"convertedCurrentPrice":[{"@currencyId":"USD","__value__":"450.0"}],"bidCount":["2"],"sellingState":["Active"],"timeLeft":["P0DT12H23M7S"]}],"listingInfo":[{"bestOfferEnabled":["false"],"buyItNowAvailable":["false"],"startTime":["2010-12-23T18:40:29.000Z"],"endTime":["2010-12-26T18:40:29.000Z"],"listingType":["Auction"],"gift":["false"]}],"returnsAccepted":["true"],"condition":[{"conditionId":["1000"],"conditionDisplayName":["New"]}]},{"itemId":["250743776675"],"title":["Nokia N8 (Unlocked)"],"globalId":["EBAY-US"],"primaryCategory":[{"categoryId":["3312"],"categoryName":["Cell Phones & Smartphones"]}],"galleryURL":["http:\/\/thumbs4.ebaystatic.com\/pict\/2507437766758080_1.jpg"],"viewItemURL":["http:\/\/cgi.ebay.com\/Nokia-N8-Unlocked-\/250743776675?pt=Cell_Phones"],"productId":[{"@type":"ReferenceID","__value__":"83383670"}],"paymentMethod":["PayPal"],"autoPay":["true"],"postalCode":["60714"],"location":["Niles,IL,USA"],"country":["US"],"shippingInfo":[{"shippingServiceCost":[{"@currencyId":"USD","__value__":"20.0"}],"shippingType":["FlatDomesticCalculatedInternational"],"shipToLocations":["Worldwide"],"expeditedShipping":["true"],"oneDayShippingAvailable":["false"],"handlingTime":["4"]}],"sellingStatus":[{"currentPrice":[{"@currencyId":"USD","__value__":"425.0"}],"convertedCurrentPrice":[{"@currencyId":"USD","__value__":"425.0"}],"bidCount":["3"],"sellingState":["Active"],"timeLeft":["P0DT12H35M59S"]}],"listingInfo":[{"bestOfferEnabled":["false"],"buyItNowAvailable":["false"],"startTime":["2010-12-16T18:53:21.000Z"],"endTime":["2010-12-26T18:53:21.000Z"],"listingType":["Auction"],"gift":["false"]}],"returnsAccepted":["false"],"condition":[{"conditionId":["1000"],"conditionDisplayName":["New"]}]},{"itemId":["250745791164"],"title":["New Factory Unlocked Nokia N8 GSM Phone  Silver White"],"globalId":["EBAY-US"],"primaryCategory":[{"categoryId":["3312"],"categoryName":["Cell Phones & Smartphones"]}],"galleryURL":["http:\/\/thumbs1.ebaystatic.com\/pict\/2507457911648080_1.jpg"],"viewItemURL":["http:\/\/cgi.ebay.com\/New-Factory-Unlocked-Nokia-N8-GSM-Phone-Silver-White-\/250745791164?pt=Cell_Phones"],"productId":[{"@type":"ReferenceID","__value__":"83383670"}],"paymentMethod":["PayPal"],"autoPay":["true"],"postalCode":["60016"],"location":["Des Plaines,IL,USA"],"country":["US"],"shippingInfo":[{"shippingServiceCost":[{"@currencyId":"USD","__value__":"0.0"}],"shippingType":["Free"],"shipToLocations":["US"],"expeditedShipping":["true"],"oneDayShippingAvailable":["false"],"handlingTime":["4"]}],"sellingStatus":[{"currentPrice":[{"@currencyId":"USD","__value__":"490.0"}],"convertedCurrentPrice":[{"@currencyId":"USD","__value__":"490.0"}],"bidCount":["0"],"sellingState":["Active"],"timeLeft":["P0DT12H51M29S"]}],"listingInfo":[{"bestOfferEnabled":["false"],"buyItNowAvailable":["true"],"buyItNowPrice":[{"@currencyId":"USD","__value__":"539.0"}],"convertedBuyItNowPrice":[{"@currencyId":"USD","__value__":"539.0"}],"startTime":["2010-12-21T19:08:51.000Z"],"endTime":["2010-12-26T19:08:51.000Z"],"listingType":["AuctionWithBIN"],"gift":["false"]}],"returnsAccepted":["false"],"galleryPlusPictureURL":["http:\/\/galleryplus.ebayimg.com\/ws\/web\/250745791164_1_0_1.jpg"],"condition":[{"conditionId":["1000"],"conditionDisplayName":["New"]}]}]}],"paginationOutput":[{"pageNumber":["1"],"entriesPerPage":["3"],"totalPages":["1577"],"totalEntries":["4730"]}],"itemSearchURL":["http:\/\/shop.ebay.com\/i.html?_nkw=Nokia+N8&_ddo=1&_ipg=3&_pgn=1"]}]}`

        v := &findItemsResponse{}
        err := xml.Unmarshal(bytes.NewBufferString(x), v)
        if err != nil { t.Error(err); return }
        if v.Ack != "Success" { t.Error("ack:",v.Ack); return }
        if v.Version != "1.8.0" { t.Error("version:",v.Version); return }
        if v.Timestamp != "2010-12-26T04:06:27.315Z" { t.Error("timestamp:",v.Timestamp); return }
        if len(v.SearchResult.Item) != 3 { t.Error("not 3 items"); return }
        if v.SearchResult.Item[0].ItemId != "120663454616" { t.Error("wrong item[0]: ",v.SearchResult.Item[0]); return }
        if v.SearchResult.Item[1].ItemId != "250743776675" { t.Error("wrong item[1]: ",v.SearchResult.Item[1]); return }
        if v.SearchResult.Item[2].ItemId != "250745791164" { t.Error("wrong item[2]: ",v.SearchResult.Item[2]); return }
        if v.SearchResult.Item[0].ShippingInfo.ShippingServiceCost.CurrencyId != "USD" { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ShippingServiceCost.Amount != 5.99 { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ShippingType != "Flat" { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ShipToLocations != "Worldwide" { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ExpeditedShipping != false { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.OneDayShippingAvailable != false { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.HandlingTime != 2 { t.Error("ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].SellingStatus.CurrentPrice.CurrencyId != "USD" { t.Error("wrong currentprice:",v.SearchResult.Item[0].SellingStatus.CurrentPrice); return }
        if v.SearchResult.Item[0].SellingStatus.CurrentPrice.Amount != 450.0 { t.Error("wrong currentprice:",v.SearchResult.Item[0].SellingStatus.CurrentPrice); return }
        if v.SearchResult.Item[0].SellingStatus.BidCount != 2 { t.Error("SellingStatus:",v.SearchResult.Item[0].SellingStatus); return }
        if v.SearchResult.Item[0].SellingStatus.SellingState != "Active" { t.Error("SellingStatus:",v.SearchResult.Item[0].SellingStatus); return }
        if v.SearchResult.Item[0].SellingStatus.TimeLeft != "P0DT14H34M2S" { t.Error("SellingStatus:",v.SearchResult.Item[0].SellingStatus); return }
        if v.SearchResult.Item[0].ListingInfo.BestOfferEnabled != false { t.Error("ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.BuyItNowAvailable != false { t.Error("ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.StartTime != "2010-12-23T18:40:29.000Z" { t.Error("ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.EndTime != "2010-12-26T18:40:29.000Z" { t.Error("ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.ListingType != "Auction" { t.Error("ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.Gift != false { t.Error("ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.ItemSearchURL != "http://shop.ebay.com/i.html?_nkw=Nokia+N8&_ddo=1&_ipg=3&_pgn=1" { t.Error("itemSearchURL: ",v.ItemSearchURL); return }
        if v.PaginationOutput.PageNumber != 1 { t.Error("paginationOutput.pageNumber:",v.PaginationOutput.PageNumber); return }
        if v.PaginationOutput.EntriesPerPage != 3 { t.Error("paginationOutput.entriesPerPage:",v.PaginationOutput.EntriesPerPage); return }
        if v.PaginationOutput.TotalPages != 1579 { t.Error("paginationOutput.totalPages:",v.PaginationOutput.TotalPages); return }
        if v.PaginationOutput.TotalEntries != 4735 { t.Error("paginationOutput.totalEntries:",v.PaginationOutput.TotalEntries); return }

        jv := &struct {
                V []findItemsJSONResponse "findItemsByKeywordsResponse"
        }{}
        err = json.Unmarshal([]byte(j), jv)
        if err != nil { t.Error(err); return }

        ja := &(jv.V[0])
        if ja.Ack[0] != "Success" { t.Error("json: Ack:",ja.Ack[0]); return }
        if ja.Version[0] != "1.8.0" { t.Error("json: Version:",ja.Version[0]); return }
        if ja.Timestamp[0] != "2010-12-26T06:17:22.587Z" { t.Error("json: Timestamp:",ja.Timestamp[0]); return }
        if len(ja.SearchResult[0].Item) != 3 { t.Error("json: not 3 items"); return }
        if ja.SearchResult[0].Item[0].ItemId[0] != "120663454616" { t.Error("json: item[0]:",ja.SearchResult[0].Item[0]); return }
        if ja.SearchResult[0].Item[1].ItemId[0] != "250743776675" { t.Error("json: item[0]:",ja.SearchResult[0].Item[1]); return }
        if ja.SearchResult[0].Item[2].ItemId[0] != "250745791164" { t.Error("json: item[0]:",ja.SearchResult[0].Item[2]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].ShippingServiceCost[0].CurrencyId != "USD" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].ShippingServiceCost[0].Amount != "5.99" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].ShippingType[0] != "Flat" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].ShipToLocations[0] != "Worldwide" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].ExpeditedShipping[0] != "false" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].OneDayShippingAvailable[0] != "false" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ShippingInfo[0].HandlingTime[0] != "2" { t.Error("ShippingInfo:", ja.SearchResult[0].Item[0].ShippingInfo[0]); return }
        if ja.SearchResult[0].Item[0].SellingStatus[0].CurrentPrice[0].CurrencyId != "USD" { t.Error("wrong currentprice:",ja.SearchResult[0].Item[0].SellingStatus[0].CurrentPrice[0]); return }
        if ja.SearchResult[0].Item[0].SellingStatus[0].CurrentPrice[0].Amount != "450.0" { t.Error("wrong currentprice:",ja.SearchResult[0].Item[0].SellingStatus[0].CurrentPrice[0]); return }
        if ja.SearchResult[0].Item[0].SellingStatus[0].BidCount[0] != "2" { t.Error("SellingStatus:",ja.SearchResult[0].Item[0].SellingStatus[0]); return }
        if ja.SearchResult[0].Item[0].SellingStatus[0].SellingState[0] != "Active" { t.Error("SellingStatus:",ja.SearchResult[0].Item[0].SellingStatus[0]); return }
        if ja.SearchResult[0].Item[0].SellingStatus[0].TimeLeft[0] != "P0DT12H23M7S" { t.Error("SellingStatus:",ja.SearchResult[0].Item[0].SellingStatus[0]); return }
        if ja.SearchResult[0].Item[0].ListingInfo[0].BestOfferEnabled[0] != "false" { t.Error("ListingInfo:",ja.SearchResult[0].Item[0].ListingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ListingInfo[0].BuyItNowAvailable[0] != "false" { t.Error("ListingInfo:",ja.SearchResult[0].Item[0].ListingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ListingInfo[0].StartTime[0] != "2010-12-23T18:40:29.000Z" { t.Error("ListingInfo:",ja.SearchResult[0].Item[0].ListingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ListingInfo[0].EndTime[0] != "2010-12-26T18:40:29.000Z" { t.Error("ListingInfo:",ja.SearchResult[0].Item[0].ListingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ListingInfo[0].ListingType[0] != "Auction" { t.Error("ListingInfo:",ja.SearchResult[0].Item[0].ListingInfo[0]); return }
        if ja.SearchResult[0].Item[0].ListingInfo[0].Gift[0] != "false" { t.Error("ListingInfo:",ja.SearchResult[0].Item[0].ListingInfo[0]); return }
        if ja.ItemSearchURL[0] != "http://shop.ebay.com/i.html?_nkw=Nokia+N8&_ddo=1&_ipg=3&_pgn=1" { t.Error("json: ItemSearchURL:",ja.ItemSearchURL[0]); return }
        if ja.PaginationOutput[0].PageNumber[0] != "1" { t.Error("json: PaginationOutput:",ja.PaginationOutput[0]); return }
        if ja.PaginationOutput[0].EntriesPerPage[0] != "3" { t.Error("json: PaginationOutput:",ja.PaginationOutput[0]); return }
        if ja.PaginationOutput[0].TotalPages[0] != "1577" { t.Error("json: PaginationOutput:",ja.PaginationOutput[0]); return }
        if ja.PaginationOutput[0].TotalEntries[0] != "4730" { t.Error("json: PaginationOutput:",ja.PaginationOutput[0]); return }

        v = noJSON(ja)
        if v.Ack != "Success" { t.Error("nojson: ack:",v.Ack); return }
        if v.Version != "1.8.0" { t.Error("nojson: version:",v.Version); return }
        if v.Timestamp != "2010-12-26T06:17:22.587Z" { t.Error("nojson: timestamp:",v.Timestamp); return }
        if len(v.SearchResult.Item) != 3 { t.Error("nojson: not 3 items"); return }
        if v.SearchResult.Item[0].ItemId != "120663454616" { t.Error("nojson: wrong item[0]: ",v.SearchResult.Item[0]); return }
        if v.SearchResult.Item[1].ItemId != "250743776675" { t.Error("nojson: wrong item[1]: ",v.SearchResult.Item[1]); return }
        if v.SearchResult.Item[2].ItemId != "250745791164" { t.Error("nojson: wrong item[2]: ",v.SearchResult.Item[2]); return }
        if v.SearchResult.Item[0].ShippingInfo.ShippingServiceCost.CurrencyId != "USD" { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ShippingServiceCost.Amount != 5.99 { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ShippingType != "Flat" { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ShipToLocations != "Worldwide" { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.ExpeditedShipping != false { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.OneDayShippingAvailable != false { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].ShippingInfo.HandlingTime != 2 { t.Error("nojson: ShippingInfo:", v.SearchResult.Item[0].ShippingInfo); return }
        if v.SearchResult.Item[0].SellingStatus.CurrentPrice.CurrencyId != "USD" { t.Error("nojson: wrong currentprice:",v.SearchResult.Item[0].SellingStatus.CurrentPrice); return }
        if v.SearchResult.Item[0].SellingStatus.CurrentPrice.Amount != 450.0 { t.Error("nojson: wrong currentprice:",v.SearchResult.Item[0].SellingStatus.CurrentPrice); return }
        if v.SearchResult.Item[0].SellingStatus.BidCount != 2 { t.Error("nojson: SellingStatus:",v.SearchResult.Item[0].SellingStatus); return }
        if v.SearchResult.Item[0].SellingStatus.SellingState != "Active" { t.Error("nojson: SellingStatus:",v.SearchResult.Item[0].SellingStatus); return }
        if v.SearchResult.Item[0].SellingStatus.TimeLeft != "P0DT12H23M7S" { t.Error("nojson: SellingStatus:",v.SearchResult.Item[0].SellingStatus); return }
        if v.SearchResult.Item[0].ListingInfo.BestOfferEnabled != false { t.Error("nojson: ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.BuyItNowAvailable != false { t.Error("nojson: ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.StartTime != "2010-12-23T18:40:29.000Z" { t.Error("nojson: ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.EndTime != "2010-12-26T18:40:29.000Z" { t.Error("nojson: ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.ListingType != "Auction" { t.Error("nojson: ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.SearchResult.Item[0].ListingInfo.Gift != false { t.Error("nojson: ListingInfo:",v.SearchResult.Item[0].ListingInfo); return }
        if v.ItemSearchURL != "http://shop.ebay.com/i.html?_nkw=Nokia+N8&_ddo=1&_ipg=3&_pgn=1" { t.Error("nojson: itemSearchURL: ",v.ItemSearchURL); return }
        if v.PaginationOutput.PageNumber != 1 { t.Error("nojson: paginationOutput.pageNumber:",v.PaginationOutput.PageNumber); return }
        if v.PaginationOutput.EntriesPerPage != 3 { t.Error("nojson: paginationOutput.entriesPerPage:",v.PaginationOutput.EntriesPerPage); return }
        if v.PaginationOutput.TotalPages != 1577 { t.Error("nojson: paginationOutput.totalPages:",v.PaginationOutput.TotalPages); return }
        if v.PaginationOutput.TotalEntries != 4730 { t.Error("nojson: paginationOutput.totalEntries:",v.PaginationOutput.TotalEntries); return }
}

func TestAppParseResponse(t *testing.T) {
        a := NewApp()
        svc := a.NewFindingService()
        xml, err := svc.FindItemsByKeywords("Nokia N8", 3)
        if err != nil { t.Error(err); return }

        res, err := a.ParseResponse(xml)
        //fmt.Printf("%v\n", res)

        if err != nil { t.Error(err); return }
        if res == nil { t.Error("no xml response"); return }
        if len(res.SearchResult.Item) != 3 {
                t.Error("xml: not 3 items:", len(res.SearchResult.Item))
                return
        }

        // switch to JSON format
        a.ResponseFormat = "JSON"

        // test JSON response
        json, err := svc.FindItemsByKeywords("Nokia N8", 3)
        if err != nil { t.Error(err); return }

        //fmt.Printf("json: %v\n",json)

        res, err = a.ParseResponse(json)
        //fmt.Printf("%v\n", res)

        if err != nil { t.Error(err); return }
        if res == nil { t.Error("no json response"); return }
        if len(res.SearchResult.Item) != 3 {
                t.Error("json: not 3 items:", len(res.SearchResult.Item))
                return
        }
}
