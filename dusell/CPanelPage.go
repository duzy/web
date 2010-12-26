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
}

func GetCPanelPage() web.ViewModel { return web.ViewModel(cpanel) }

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
                eBayView := &cpaneleBayTest{}
                cp.content = web.NewStringer(eBayView, app)
                cp.rightside = "TODO: eBay test..."
        case "/sync":
                syncView := &cpanelSync{ "TODO: ..." }
                cp.content = web.NewStringer(syncView, app)
                cp.rightside = "TODO: side column for sync view"
        case "/categories":
                cp.content = "TODO: show categories"
                cp.rightside = "TODO: categories commands"
        case "/items":
                cp.content = "TODO: show items"
                cp.rightside = "TODO: items commands"
        case "/sells":
                cp.content = "TODO: show sells"
                cp.rightside = "TODO: sells commands"
        }
        handled = true
        return
}

func (v *cpaneleBayTest) GetTemplate() (str string) {
        eb := eBay.NewApp()
        fs := eb.NewFindingService()
        xml, err := fs.FindItemsByKeywords("Nokia N8", 3)
        if err == nil {
                r, err := eb.ParseResponse(xml)
                if err != nil {
                        str = fmt.Sprintf("<b>ERROR:</b> %v", err)
                        return
                }
                str = "<table>"
                for n, i := range r.SearchResult.Item {
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
