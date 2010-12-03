package main

import (
        "../_obj/web"
        "./_obj/duzyinfo"
        "flag"
)

var flagPath = flag.String("path", "", "Set CGI PATH_INFO for debug.")

func setupCGIModel(model web.AppModel) {
        if cgi, ok := model.(*web.CGIModel); ok {
                if *flagPath != "" { cgi.Setenv("PATH_INFO", *flagPath) }
        }
}

func main() {
        flag.Parse()

        homeView := web.NewView(duzyinfo.GetHomePage())

        model := web.NewCGIModel()

        setupCGIModel(model)
        
        app := web.NewApp("Duzy Chan", model)
        app.HandleDefault(homeView)
        app.Exec()
}
