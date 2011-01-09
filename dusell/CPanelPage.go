package dusell

import (
        "../_obj/web"
        "../_obj/eBay"
        "strings"
        "fmt"
)

type cpanelPage struct {
        web.TemplateFile
        base string // the base dir in which this scrit locates
        script string // the path of the script
        content interface{}
        rightside interface{}

        eBay *eBay.App
}

type cpaneleBayTest struct {
        //web.TemplateString
}

type cpanelSync struct {
        web.TemplateString
}

type cpanelCategories struct {
        web.TemplateString
}

var cpanel *cpanelPage = &cpanelPage {
        web.TemplateFile(TemplateCPanelPage),
        "", // base dir
        "", // script name
        "", // content
        "", // right side content
        nil,
}

func GetCPanelPage() web.ViewModel {
        if cpanel.eBay == nil {
                cpanel.eBay = eBay.NewApp()
        }
        return web.ViewModel(cpanel)
}

func (cp *cpanelPage) MakeFields(app *web.App) interface{} {
        cp.script = app.ScriptName()
        n := strings.LastIndex(cp.script, "/")
        if n != -1 {
                cp.base = cp.script[0:n]
        } else {
                // this should not happen
                cp.base = cp.script[0:]
        }
        
        if cp.content == "" {
                cp.content = "TODO: default content"
                cp.rightside = "TODO: default commands"
        }
        return cp
}

func (cp *cpanelPage) HandleSubpath(subpath string, app *web.App) (handled bool) {
        switch subpath {
        case "/eBay": // eBay test
                handled = true
                eBayView := &cpaneleBayTest{}
                cp.content = web.NewStringer(eBayView, app)
                cp.rightside = "TODO: eBay test..."
        case "/sync":
                handled = true
                syncView := &cpanelSync{ "TODO: ..." }
                cp.content = web.NewStringer(syncView, app)
                cp.rightside = "TODO: side column for sync view"
        case "/categories":
                handled = true
                catsView := &cpanelCategories{}
                cp.content = web.NewStringer(catsView, app)
                cp.rightside = "TODO: categories commands"
        case "/items":
                handled = true
                cp.content = "TODO: show items"
                cp.rightside = "TODO: items commands"
        case "/sells":
                handled = true
                cp.content = "TODO: show sells"
                cp.rightside = "TODO: sells commands"
        }
        return
}

func (v *cpaneleBayTest) GetTemplate() (str string) {
        eb := eBay.NewApp()
        svc := eb.NewFindingService()
        call := svc.NewFindItemsByKeywordsCall()
        call.Keywords = "Nokia N8"
        call.SetEntriesPerPage(3)
        xml, err := eb.Invoke(call)
        if err == nil {
                resp := new(eBay.FindItemsByKeywordsResponse)
                err := eb.ParseXMLResponse(resp, xml)
                if err != nil {
                        str = fmt.Sprintf("<b>ERROR:</b> %v", err)
                        return
                }
                str = "<table>"
                for n, i := range resp.SearchResult.Item {
                        str += "<tr>"
                        str += fmt.Sprintf("<td>%v</td>", n)
                        str += fmt.Sprintf(`<td><a href="%s"><img src="%s"/></a></td>`, i.ViewItemURL, i.GalleryURL)
                        str += fmt.Sprintf("<td>%s</td>", i.Title)
                        str += "</tr>"
                }
                str += "</table>"
        } else {
                str = fmt.Sprintf("<b>ERROR:</b> %v", err)
        }
        return
}

func (view *cpanelCategories) GetTemplate() (str string) {
        cache, err := eBay.NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil {
                str = fmt.Sprintf("<b>ERROR</b>: ", err)
                return
        }

        cats, err := cache.GetCategoriesByLevel(1)
        /*
        if err != nil {
                str = fmt.Sprintf("<b>ERROR</b>: ", err)
                return
        }
         */

        str = "<table>"
        for _, c := range cats {
                str += "<tr>"
                str += fmt.Sprintf(`<td><a href="%s">%s</a></td>`, c.CategoryID, c.CategoryID)
                str += fmt.Sprintf("<td>%s</td>", c.CategoryName)
                str += fmt.Sprintf(`<td><a href="%s">%s</a></td>`, c.CategoryParentID, c.CategoryParentID)
                str += "</tr>"
        }
        str += "</table>"
        return
}
