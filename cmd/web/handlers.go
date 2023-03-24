package main

import "net/http"

func (app *application) HelloWorld(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Hit HelloWorld Handler")
}
