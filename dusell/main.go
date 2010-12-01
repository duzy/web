package main

import "./dusell"
import "../web"

func main() {
        home := dusell.MakeHomePage("home.tpl")
        homeView := web.NewView(home)
        app := web.NewApp("DuSell.com")
        app.Handle("", homeView)
        app.Handle("/home", homeView)
        app.Exec()
}
