package eBay

import (
        "testing"
        "strings"
        //"fmt"
)

func TestAPIFindItemsByKeywords(t *testing.T) {
        {
                eb := NewApp()
                eb.ResponseFormat = "XML" // "JSON"

                s, err := eb.FindItemsByKeywords("iPhone", 3)
                if err != nil { t.Error(err); goto finish }

                //fmt.Printf("%s\n", s)

                n := strings.Index(s, "<findItemsByKeywordsResponse ")
                if n == -1 { t.Error("no tag found: findItemsByKeywordsResponse"); goto finish }

                s = s[n+30:len(s)]
                n = strings.Index(s, "<item>")
                if n == -1 { t.Error("no tag found: <item> [1]"); goto finish }

                s = s[n+6:len(s)]
                n = strings.Index(s, "<item>")
                if n == -1 { t.Error("no tag found: <item> [2]"); goto finish }

                s = s[n+6:len(s)]
                n = strings.Index(s, "<item>")
                if n == -1 { t.Error("no tag found: <item> [3]"); goto finish }
        }
        { // JSON test
                eb := NewApp()
                eb.ResponseFormat = "JSON"

                s, err := eb.FindItemsByKeywords("Nokia N9", 3)
                if err != nil { t.Error(err); goto finish }

                //fmt.Printf("%s\n", s)

                n := strings.Index(s, "\"findItemsByKeywordsResponse\":")
                if n == -1 { t.Error("no prop: \"findItemsByKeywordsResponse\""); goto finish }

                s = s[n+20:len(s)]
                n = strings.Index(s, "\"item\":")
                if n == -1 { t.Error("no prop: \"item\""); goto finish }

                s = s[n+7:len(s)]
                n = strings.Index(s, "[{\"itemId\":")
                if n == -1 { t.Error("no prop: \"itemId\" [1]"); goto finish }

                s = s[n+10:len(s)]
                n = strings.Index(s, "{\"itemId\":")
                if n == -1 { t.Error("no prop: \"itemId\" [2]"); goto finish }

                s = s[n+10:len(s)]
                n = strings.Index(s, "{\"itemId\":")
                if n == -1 { t.Error("no prop: \"itemId\" [3]"); goto finish }
        }
finish:
}
