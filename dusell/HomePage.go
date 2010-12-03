package dusell

import "../_obj/web"

type homePage struct {
        web.ViewTemplateName
        web.StandardFields
}

var home *homePage = &homePage{
        web.ViewTemplateName{ TemplateHomePage },
        web.StandardFields{ nil },
}

// Get the singleton homePage object.
func GetHomePage() web.ViewModel { return web.ViewModel(home) }

func (h *homePage) MakeFields(app *web.App) (fields interface{}) {
        fields = h.StandardFields.MakeFields(app)
        names := []string{ "name1", "name2", "name3" }
        h.StandardFields.SetField("names", names)
        return
}
