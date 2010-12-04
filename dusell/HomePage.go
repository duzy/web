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
        h.StandardFields["headline"] = web.NewStringer(headline, app)
        h.StandardFields["page_content_left"] = web.NewStringer(contentLeft, app)
        h.StandardFields["page_content_center"] = web.NewStringer(contentCenter, app)
        h.StandardFields["page_content_right"] = web.NewStringer(contentRight, app)
        h.StandardFields["page_content_tail"] = web.NewStringer(contentTail, app)
        return
}

var headline = &Headline {
        "TODO:headline",
}

var contentLeft = &ContentLeft {
        "TODO:content_left",
}

var contentCenter = &ContentCenter {
        "TODO:content_center",
}

var contentRight = &ContentRight {
        "TODO: content_right",
}

var contentTail = &ContentTail {
        "TODO: content_tail",
}        

type Headline struct {
        web.TemplateString
}
func (vm *Headline) MakeFields(app *web.App) interface{} { return vm }

type ContentLeft struct {
        web.TemplateString
}
func (vm *ContentLeft) MakeFields(app *web.App) interface{} { return vm }

type ContentCenter struct {
        web.TemplateString
}
func (vm *ContentCenter) MakeFields(app *web.App) interface{} { return vm }

type ContentRight struct {
        web.TemplateString
}
func (vm *ContentRight) MakeFields(app *web.App) interface{} { return vm }

type ContentTail struct {
        web.TemplateString
}
func (vm *ContentTail) MakeFields(app *web.App) interface{} { return vm }
