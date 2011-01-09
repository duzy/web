package eBay

import (
        "testing"
        "strings"
        //"fmt"
        "os"
)

func TestTradingCallGetCategories(t *testing.T) {
        a := NewApp()

        trading := a.NewTradingService()
        c := trading.NewGetCategoriesCall()
        c.CategoryParent = "20081"
        c.LevelLimit = 1
        c.ViewAllNodes = true
        c.DetailLevel = ReturnAll

        xml, err := a.Invoke(c)
        if err != nil { t.Errorf("%v\n", err); return }

        f, err := os.Open("testout.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
        if err == nil {
                var n int
                n, err = f.WriteString(xml) //fmt.Fprint(f, xml)
                if err != nil { t.Errorf("can't save testout.xml: %v", err) }
                if err = f.Close(); err != nil { t.Errorf("close: %v", err) }
                t.Logf("test: %d bytes(%d bytes) saved into testout.xml", n, len(xml))
        } else {
                t.Errorf("can't save testout.xml: %v", err);
        }

        n := strings.Index(xml, "<Errors>")
        if n != -1 { t.Errorf("GetCategoriesResponse:<Errors>: %v", xml); return }

        s := xml[0:len(xml)]
        n = strings.Index(s, `<GetCategoriesResponse`)
        if n == -1 { t.Errorf("GetCategoriesResponse: %v", xml); return }

        s = s[n:len(s)]
        n = strings.Index(s, `xmlns="urn:ebay:apis:eBLBaseComponents"`)
        if n == -1 { t.Error("GetCategoriesResponse"); return }

        s = s[n:len(s)]
        n = strings.Index(s, `<Ack>Success</Ack>`)
        if n == -1 { t.Errorf("GetCategoriesResponse: %v", xml); return }
}
