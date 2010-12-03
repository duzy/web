package main

import (
        "../_obj/web"
        "./_obj/dusell"
        "flag"
)

var flagPath = flag.String("path", "", "Set PATH_INFO for debug.")

func setupCGIModel(model web.AppModel) {
        if cgi, ok := model.(*web.CGIModel); ok {
                if *flagPath != "" { cgi.Setenv("PATH_INFO", *flagPath) }
        }
}

func main() {
        flag.Parse()

        homeView := web.NewView(dusell.GetHomePage())

        model := web.NewCGIModel()

        setupCGIModel(model)

        app := web.NewApp("DuSell.com", model)
        app.HandleDefault(homeView)
        app.Handle("/home", homeView)
        //app.Handle("/order", OrderHandler)
        //app.Handle("/pay", PaymentHandler)
        //app.Handle("/cats", CatalogsHandler)
        //app.Handle("/cat", GetItemList)
        //app.Handle("/get_item", GetItem)
        //app.Handle("/signin", SigninHandler)
        //app.Handle("/signup", SignupHandler)
        app.Exec()
}
