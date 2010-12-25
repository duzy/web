package main

import (
        "../_obj/web"
        "./_obj/dusell"
)

func dusell_run(app *web.App) {
        cpanelView := web.NewView(dusell.GetCPanelPage())
        app.Handle("/cp", cpanelView)
        app.Exec()
}
