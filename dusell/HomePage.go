package dusell

import (
        "../_obj/web"
        "../_obj/eBay"
        "bytes"
        "fmt"
        "os"
)

type homePageFeed struct {
        visitCount int
        title string
        keywords string
        description string
        page_header string
        page_content_left web.View
        page_content web.View
        page_content_right web.View
        page_footer string
}

type homePage struct { *web.HtmlPage }

type homePageContentLeft struct {
        html string
}

type homePageContent struct {
        html string
}

type homePageContentRight struct {
        html string
}

var catPath string
var contentLeft = &homePageContentLeft{}
var content = &homePageContent{}
var contentRight = &homePageContentRight{}

func (feed *homePageFeed) Prepare(request *web.Request) (err os.Error) {
        feed.visitCount += 1
        catPath = request.ScriptName + "/cat"
        return
}

func GetHomePage() web.RequestHandler {
        var feed = &homePageFeed{
        visitCount: 0,
        title: "Dusell.com",
        keywords: "",
        description: "",
        page_header: "<h1>We list good items</h1>",
        page_content_left: web.NewView(contentLeft),
        page_content: web.NewView(content),
        page_content_right: web.NewView(contentRight),
        page_footer: "&copy; 2011",
        }

        var home, _ = web.NewHtmlPage(TemplateHomePage, feed)
        return homePage{ home }
}

func (h homePage) HandleSubpath(subpath string, request *web.Request) (handled bool) {
        switch {
        case 4 <= len(subpath) && subpath[0:4] == "/cat":
                handled = true
                if 4 < len(subpath) {
                        content.ListCategory(subpath[5:len(subpath)])
                }
                contentRight.html = subpath
        }
        return
}

func (v *homePageContentLeft) Render(w web.RenderBuffer) (err os.Error) {
        cats, err := GetHotCategories()
        if err != nil {
                // TODO: using some kind of standard error formater
                return
        }

        fmt.Fprintf(w, "<ul>")
        for _, cat := range cats {
                fmt.Fprintf(w, `<li><a href="%s/%s">%s</a></li>`,
                        catPath,
                        cat.CategoryID,
                        cat.CategoryName)
        }
        fmt.Fprintf(w, "</ul>")
        
        return
}

func GetHotCategories() (cats []*eBay.Category, err os.Error) {
        cc, err := eBay.NewDBCache("localhost", "test", "abc", "dusell")
        if err != nil { return }

        cats, err = cc.GetCategoriesByLevel(1)
        return
}

func (v *homePageContentRight) Render(w web.RenderBuffer) (err os.Error) {
        fmt.Fprint(w, v.html)
        return
}

func (v *homePageContent) Render(w web.RenderBuffer) (err os.Error) {
        fmt.Fprint(w, v.html)
        return
}

func (v *homePageContent) ListCategory(cat string) (err os.Error) {
        a := eBay.NewApp()
        svc := a.NewFindingService()
        call := svc.NewFindItemsByCategoryCall()
        call.CategoryId = cat

        s, err := a.Invoke(call)
        if err != nil {
                v.html = err.String()
                return
        }

        resp := new(eBay.FindItemsByCategoryResponse)
        err = a.ParseXMLResponse(resp, s)
        if err != nil {
                v.html = err.String()
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
        v.html = buf.String()
        return
}

