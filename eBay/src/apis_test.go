package eBay

import (
        "testing"
        "fmt"
)

func TestAPIFindItemsByKeywords(t *testing.T) {
        eb := NewEBay()
        eb.ResponseFormat = "XML" // "JSON"
        s := eb.FindItemsByKeywords("iPhone", 3)
        fmt.Printf("%s\n", s)
}
