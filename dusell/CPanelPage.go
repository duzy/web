package dusell

import (
        "../_obj/web"
)

type cpanelPage struct {
        web.TemplateFile
        script string
        content string
        rightside string
}

var cpanel *cpanelPage = &cpanelPage {
        web.TemplateFile(TemplateCPanelPage),
        "", // script name
        "", // content
        "", // right side
}

func GetCPanelPage() web.ViewModel { return web.ViewModel(cpanel) }

func (cp *cpanelPage) HandleSubpath(subpath string, app *web.App) (handled bool) {
        switch subpath {
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

func (cp *cpanelPage) MakeFields(app *web.App) interface{} {
        cp.script = app.ScriptName()
        if cp.content == "" {
                cp.content = "TODO: default content"
                cp.rightside = "TODO: default commands"
        }

        cp.rightside  = "raw-cookie: <br/>" + app.RawCookie()
        cp.rightside += "<hr/>"
        cp.rightside += "cookies: <br/>"
        for _, c := range app.Cookies() {
                cp.rightside += c.Name + "=" + c.Content + "<br/>"
        }

        test := app.Session().Get("test")

        cp.rightside += "<hr/>"
        cp.rightside += "test: " + test

        if len(test)%2 == 0 { test += "1" } else { test += "0" }
        app.Session().Set("test", test)
        return cp
}
