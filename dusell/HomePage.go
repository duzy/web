package dusell

import (
        "../_obj/web"
        //"fmt"
)

type homePage struct {
        web.TemplateFile
        web.StandardFields
}

// The internal singleton homePage object.
var home *homePage = &homePage{
        web.TemplateFile(TemplateHomePage),
        make(web.StandardFields),
}

// Get the singleton home page web.ViewModel.
func GetHomePage() web.ViewModel { return web.ViewModel(home) }

func (h *homePage) MakeFields(app *web.App) (fields interface{}) {
        fields = h.StandardFields.MakeFields(app)
        h.StandardFields["keywords"] = "DuSell Online Store"
        h.StandardFields["description"] = "DuSell Online Store do what you buy"
        h.StandardFields["page-header"] = web.NewStringer(header, app)
        h.StandardFields["page-content"] = web.NewStringer(content, app)
        h.StandardFields["page-content-left"] = web.NewStringer(contentLeft, app)
        h.StandardFields["page-content-right"] = web.NewStringer(contentRight, app)
        h.StandardFields["page-footer"] = web.NewStringer(footer, app)
        return
}

var header = &PageHeader {
        "TODO:page header",
}

var content = &PageContent {
        "TODO:page content",
}

var contentLeft = &PageContentLeft {
        "TODO:page content left",
}

var contentRight = &PageContentRight {
        "TODO:page content right",
}

var footer = &PageContent {
        "TODO:page footer",
}

type PageHeader struct {
        web.TemplateString
}

type PageContent struct {
        web.TemplateString
}

type PageContentLeft struct {
        web.TemplateString
}

type PageContentRight struct {
        web.TemplateString
}

type PageFooter struct {
        web.TemplateString
}

