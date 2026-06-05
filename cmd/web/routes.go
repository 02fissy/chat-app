package main
import(
	"net/http"
	"chat.fisayo.net/ui"
)


func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /chat", app.serveHome)
	mux.HandleFunc("GET /user/signup",app.userSignup)
	mux.HandleFunc("POST /user/signup",app.userSignupPost)
	mux.HandleFunc("GET /user/login",app.userLogin)
	mux.HandleFunc("POST /user/login",app.userLoginPost)
	mux.HandleFunc("GET /user/logout",app.userLogoutPost)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(app.hub, w, r)
	})
	return app.sessionManager.LoadAndSave(mux)
}