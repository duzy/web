package duzyinfo

import (
        "../_obj/web"
)

type homePage struct {
        title string
        web.TemplateName
}

var home *homePage = &homePage{
        "Duzy Chan",
        web.TemplateName(TemplateHomePage),
}
func GetHomePage() web.ViewModel { return web.ViewModel(home) }

func (h *homePage) MakeFields(app *web.App) interface{} { return(h) }
