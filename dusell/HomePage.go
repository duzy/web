package dusell

import "../web"

type homePage struct {
        web.DefaultView // embed fields and methods of web.DefaultView
}

func MakeHomePage(tpl string) (model web.Model) {
        //home := homePage{ web.DefaultView{tpl,nil} }
        home := new(homePage)
        home.Template = tpl
        model = web.Model(home)
        return
}

func (home *homePage) MakeFields(app *web.App) (fields map[string]string) {
        fields = home.DefaultView.MakeFields(app)
        // TODO: ...
        return
}
