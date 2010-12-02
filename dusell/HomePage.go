package dusell

import "../_obj/web"

type homePage struct {
        web.DefaultView // embed fields and methods of web.DefaultView
}

var home *homePage // Singleton object.

// Get the singleton homePage object.
func GetHomePage() (model web.ViewModel) {
        if home == nil {
                home = &homePage{ web.DefaultView{ TemplateHomePage, nil } }
        }

        model = web.ViewModel(home)
        return
}

func (home *homePage) MakeFields(app *web.App) (fields map[string]string) {
        fields = home.DefaultView.MakeFields(app)
        // TODO: ...
        return
}
