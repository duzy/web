package main

import "./dusell"
import "../_obj/web"

func main() {
        home := dusell.MakeHomePage("home.tpl")
        homeView := web.NewView(home)
        app := web.NewApp("DuSell.com", web.NewCGIModel())
        app.Handle("", homeView)
        app.Handle("/home", homeView)
        app.Exec()
}
