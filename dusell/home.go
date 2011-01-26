package main

import (
        "../_obj/web"
        "./_obj/dusell"
)

func dusell_run(app *web.App) {
        app.HandleDefault(dusell.GetHomePage())
        app.Exec()
}
