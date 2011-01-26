package dusell

import (
        "../_obj/web"
        //"../_obj/eBay"
        //"bytes"
        //"fmt"
        "os"
)

type homePageFeed struct {
        title string
        keywords string
        description string
        page_header string
        page_content_left *web.View
        page_content *web.View
        page_content_right *web.View
        page_footer string
}

type homePageContentLeft struct {
}

type homePageContent struct {
}

type homePageContentRight struct {
}

var contentLeft = &homePageContentLeft{}
var content = &homePageContent{}
var contentRight = &homePageContentRight{}

// The internal singleton homePage object.
var feed = &homePageFeed{
title: "Dusell.com",
keywords: "",
description: "",
page_header: "<h1>We list good items</h1>",
page_content_left: nil,
page_content: nil,
page_content_right: nil,
page_footer: "&copy; 2011",
}

func (feed *homePageFeed) PrepareFeed(request *web.Request) (err os.Error) {
        feed.page_content_left = web.NewView(contentLeft)
        feed.page_content = web.NewView(content)
        feed.page_content_right = web.NewView(contentRight)
        return
}

var home, _ = web.NewHtmlPage(TemplateHomePage, feed)

// Get the singleton home page web.ViewModel.
func GetHomePage() *web.HtmlPage { return home }

func (v *homePageContentLeft) Render(w web.RenderBuffer) (err os.Error) {
        return
}

func (v *homePageContent) Render(w web.RenderBuffer) (err os.Error) {
        return
}

func (v *homePageContentRight) Render(w web.RenderBuffer) (err os.Error) {
        return
}

/*
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
*/
