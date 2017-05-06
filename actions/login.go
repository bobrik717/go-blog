package actions

import (
	"go-book/blog/session"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
	"time"
)

func LoginIndexHandler(rnd render.Render) {
	rnd.HTML (200,"login",nil)
}

func LoginHandler(rnd render.Render, r *http.Request, w http.ResponseWriter) {
	login := r.FormValue("login")
	//password := r.FormValue("password")

	sessionId := session.SessionMain.Init(login)

	cookie := &http.Cookie{
		Name: session.COOKIE_NAME,
		Value: sessionId,
		Expires: time.Now().Add( 5 * time.Minute),
	}

	http.SetCookie(w, cookie)

	rnd.Redirect("/")
}
