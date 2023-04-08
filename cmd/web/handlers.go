package main

import "net/http"

func (app *application) HelloWorld(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit HelloWorld Handler")
}

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "terminal", nil); err != nil {
		app.errorLog.Println(err)
	}
}
