package dusell

import (
        "../_obj/web"
        "../_obj/eBay"
        "bytes"
        "fmt"
        "os"
)

type homePage struct {
        web.TemplateFile
        web.StandardFields
}

type PageHeader struct {
        web.TemplateString
}

type PageContent struct {
        web.TemplateString
}

type PageContentLeft struct {
        //web.TemplateString
        catPath string 
}

type PageContentRight struct {
        web.TemplateString
}

type PageFooter struct {
        web.TemplateString
}

// The internal singleton homePage object.
var home *homePage = &homePage{
        web.TemplateFile(TemplateHomePage),
        make(web.StandardFields),
}

var header = &PageHeader {
        "<h1>We list good items</h1>",
}

var contentLeft = &PageContentLeft {}
var content = &PageContent {}
var contentRight = &PageContentRight {}
var footer = &PageContent {
        "&copy; 2011",
}

// Get the singleton home page web.ViewModel.
func GetHomePage() web.ViewModel { return web.ViewModel(home) }

func (h *homePage) MakeFields(app *web.App) (fields interface{}) {
        contentLeft.catPath = app.ScriptName() + "/cat"

        fields = h.StandardFields.MakeFields(app)
        h.StandardFields["keywords"] = "DuSell Online Store"
        h.StandardFields["description"] = "DuSell Online Store do what you buy"
        h.StandardFields["page-header"] = web.NewStringer(header, app)
        h.StandardFields["page-content-left"] = web.NewStringer(contentLeft, app)
        h.StandardFields["page-content"] = web.NewStringer(content, app)
        h.StandardFields["page-content-right"] = web.NewStringer(contentRight, app)
        h.StandardFields["page-footer"] = web.NewStringer(footer, app)
        return
}

func (h *homePage) HandleSubpath(subpath string, app *web.App) (handled bool) {
        switch {
        case 4 <= len(subpath) && subpath[0:4] == "/cat":
                handled = true
                if 4 < len(subpath) {
                        //content.TemplateString = web.TemplateString("TODO: list category " + subpath[5:len(subpath)])
                        a := eBay.NewApp()
                        svc := a.NewFindingService()
                        call := svc.NewFindItemsByCategoryCall()
                        call.CategoryId = subpath[5:len(subpath)]

                        s, err := a.Invoke(call)
                        if err != nil {
                                content.TemplateString = web.TemplateString(err.String())
                                return
                        }

                        resp := new(eBay.FindItemsByCategoryResponse)
                        err = a.ParseXMLResponse(resp, s)
                        if err != nil {
                                content.TemplateString = web.TemplateString(err.String())
                                return
                        }

                        buf := bytes.NewBuffer(make([]byte, 0, 128))
                        fmt.Fprintf(buf, `<table>`)
                        for i:=0; i < len(resp.SearchResult.Item); i += 1 {
                                itm := &(resp.SearchResult.Item[i])
                                fmt.Fprintf(buf, `<tr>`)
                                fmt.Fprintf(buf, `<td><a href="%s"><img src="%s"/></a></td>`, itm.ViewItemURL, itm.GalleryURL)
                                fmt.Fprintf(buf, `<td>%s</td>`, itm.Title)
                                fmt.Fprintf(buf, `</tr>`)
                        }
                        fmt.Fprintf(buf, `</table>`)
                        content.TemplateString = web.TemplateString(buf.String())
                }
                contentRight.TemplateString = web.TemplateString(subpath)
        }
        return
}

func (view *PageContentLeft) GetTemplate() string {
        temp := bytes.NewBufferString("")
        cats, err := GetHotCategories()
        if err != nil {
                // TODO: using some kind of standard error formater
                return err.String()
        }

        fmt.Fprintf(temp, "<ul>")
        for _, cat := range cats {
                fmt.Fprintf(temp, `<li><a href="{catPath}/%s">%s</a></li>`,
                        cat.CategoryID,
                        cat.CategoryName)
        }
        fmt.Fprintf(temp, "</ul>")
        
        return temp.String()
}

func GetHotCategories() (cats []*eBay.Category, err os.Error) {
        cc, err := eBay.NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { return }

        cats, err = cc.GetCategoriesByLevel(1)
        return
}
