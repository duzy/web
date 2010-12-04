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
        h.StandardFields["page_content_left"] = "content_left"
        h.StandardFields["page_content_center"] = "content_center"
        h.StandardFields["page_content_right"] = "content_right"
        h.StandardFields["page_content_tail"] = "tail"
        return
}

type Headline struct {
        web.TemplateString
}

var headline = &Headline {
        "TODO:headline",
}

func (vm *Headline) MakeFields(app *web.App) interface{} {
        return interface{}(vm)
}
