package main

import (
        "../_obj/web"
        //"flag"
        //"fmt"
        //"os"
)

/*
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

        app, err := web.NewApp("config.json")
        if err != nil { panic(fmt.Sprintf("NewApp: %v", err)) }
        if app == nil { panic("NewApp: failed") }
        setupCGIModel(app.GetModel())

        dusell_run(app)
}
*/

func main() {
        fcgi, _ := web.NewFCGIModel()
        app, _ := web.NewApp(fcgi)
        dusell_run(app)
}
