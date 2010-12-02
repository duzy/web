package main

import (
        "../_obj/web"
        "./_obj/dusell"
)

func main() {
        homeView := web.NewView(dusell.GetHomePage())

        app := web.NewApp("DuSell.com", web.NewCGIModel())
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
