package main

import (
        "../_obj/web"
        "./_obj/dusell"
        "flag"
        "os"
)

var flagPath = flag.String("path", "", "Set PATH_INFO for debug.")

func setupCGIModel(model web.AppModel) {
        if cgi, ok := model.(*web.CGIModel); ok {
                if *flagPath != "" { cgi.Setenv("PATH_INFO", *flagPath) }
                if cgi.Getenv("SCRIPT_NAME") == "" {
                        cgi.Setenv("SCRIPT_NAME", os.Args[0])
                }
        }
}

func main() {
        flag.Parse()

        homeView := web.NewView(dusell.GetHomePage())
        cpanelView := web.NewView(dusell.GetCPanelPage())

        model := web.NewCGIModel()

        setupCGIModel(model)

        app := web.NewApp("Dusell - find what you want", model)
        app.HandleDefault(homeView)
        app.Handle("/cp", cpanelView)
        //app.Handle("/order", OrderHandler)
        //app.Handle("/pay", PaymentHandler)
        //app.Handle("/cats", CatalogsHandler)
        //app.Handle("/cat", GetItemList)
        //app.Handle("/get_item", GetItem)
        //app.Handle("/signin", SigninHandler)
        //app.Handle("/signup", SignupHandler)
        app.Exec()
}
