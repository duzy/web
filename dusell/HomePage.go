package dusell

import "../web"

type homePage struct {
        web.DefaultView;
}

func MakeHomePage(tpl string) (model web.Model) {
        home := new(homePage)
        home.Template = tpl
        model = web.Model(home)
        return
}

func (home *homePage) MakeFields(app *web.App) (fields map[string]string) {
        return
}

