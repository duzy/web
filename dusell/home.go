package main

import (
        "../_obj/web"
        "./_obj/dusell"
)

func dusell_run(app *web.App) {
        homeView := web.NewView(dusell.GetHomePage())
        app.HandleDefault(homeView)
        //app.Handle("/order", OrderHandler)
        //app.Handle("/pay", PaymentHandler)
        //app.Handle("/cats", CatalogsHandler)
        //app.Handle("/cat", GetItemList)
        //app.Handle("/get_item", GetItem)
        //app.Handle("/signin", SigninHandler)
        //app.Handle("/signup", SignupHandler)
        app.Exec()
}
