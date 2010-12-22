package main

import (
        "../_obj/web"
        "./_obj/dusell"
        "flag"
        "fmt"
        "os"
)

var flagPath = flag.String("path", "", "Set PATH_INFO for debug.")
var flagSid = flag.String("session", "", "Set session id for debug.")

func setupCGIModel(model web.AppModel) {
        if cgi, ok := model.(*web.CGIModel); ok {
                if *flagPath != "" { cgi.Setenv("PATH_INFO", *flagPath) }
                if *flagSid != "" {
                        cookie := "__web_cid=" + *flagSid
                        cgi.Setenv("HTTP_COOKIE", cookie)
                }
                if cgi.Getenv("SCRIPT_NAME") == "" {
                        cgi.Setenv("SCRIPT_NAME", os.Args[0])
                }
        }
}

func main() {
        flag.Parse()

        defer web.CGIHandleErrors()

        homeView := web.NewView(dusell.GetHomePage())
        cpanelView := web.NewView(dusell.GetCPanelPage())

        app, err := web.NewApp("config.json")
        if err != nil { panic(fmt.Sprintf("NewApp: %v", err)) }
        if app == nil { panic("NewApp: failed") }
        setupCGIModel(app.GetModel())

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
