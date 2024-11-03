package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/otakakot/loginapp/internal/base/firebase"
	"github.com/otakakot/loginapp/internal/base/pocketbase"
	"github.com/otakakot/loginapp/internal/base/supabase"
)

const cookey = "__session__"

type AuthN struct {
	firebase   *firebase.Firebase
	supabase   *supabase.Supabase
	pocketbase *pocketbase.Pocketbase
}

func New(
	firebase *firebase.Firebase,
	supabase *supabase.Supabase,
	pocketbase *pocketbase.Pocketbase,
) *AuthN {
	return &AuthN{
		firebase:   firebase,
		supabase:   supabase,
		pocketbase: pocketbase,
	}
}

func (an *AuthN) Handle(
	rw http.ResponseWriter,
	req *http.Request,
) {
	slog.Info(req.Method + " : " + req.URL.Path + " start")
	defer slog.Info(req.Method + " : " + req.URL.Path + " done")

	switch req.Method {
	case http.MethodGet:
		an.Get(rw, req)

		return
	case http.MethodPost:
		an.Post(rw, req)

		return
	case http.MethodPatch:
		an.Patch(rw, req)

		return
	case http.MethodPut:
		an.Put(rw, req)

		return
	case http.MethodDelete:
		an.Delete(rw, req)

		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

const login = `<!DOCTYPE html>
<html lang="en">

<head>
    <title>Login</title>
</head>

<body>
    <form method="POST" action="/">
		<div>
		<input type="radio" name="base" value="firebase" checked>Firebase</input>
		<input type="radio" name="base" value="supabase">Supabase</input>
		<input type="radio" name="base" value="pocketbase">Pocketbase</input>
		</div>

        <label for="email">Email</label>
        <input type="email" name="email" id="email">

        <label for="password">Password</label>
        <input type="password" name="password" id="password">

        <button type="submit">Login</button>
    </form>
</body>

</html>
`

const session = `<!DOCTYPE html>
<html lang="en">

<head>
	<title>Session</title>
</head>

<body>
	<button id="verify">Verify</button>
	<button id="refresh">Refresh</button>
	<button id="logout">Logout</button>

	<script>
		document.getElementById("verify").
			addEventListener("click", async () => {
				const res = await fetch("/", {
					method: "PATCH",
					credentials: "include",
				});
				if (res.ok) {
					location.href = "/";
				}
			});
		document.getElementById("refresh").
			addEventListener("click", async () => {
				const res = await fetch("/", {
					method: "PUT",
					credentials: "include",
				});
				if (res.ok) {
					location.href = "/";
				}
			});
		document.getElementById("logout").
			addEventListener("click", async () => {
				const res = await fetch("/", {
					method: "DELETE",
					credentials: "include",
				});
				if (res.ok) {
					location.href = "/";
				}
			});
	</script>
</body>

</html>
`

// Get is session
func (an *AuthN) Get(
	rw http.ResponseWriter,
	req *http.Request,
) {
	if _, err := req.Cookie(cookey); err != nil {
		rw.Write([]byte(login))

		return
	}

	rw.Write([]byte(session))
}

// Post is login
func (an *AuthN) Post(
	rw http.ResponseWriter,
	req *http.Request,
) {
	slog.Info(req.FormValue("base"))

	value := ""

	switch req.FormValue("base") {
	case "firebase":
		res, err := an.firebase.Auth(req.Context(), firebase.AuthRequest{
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		})
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)

			rw.Write([]byte(err.Error()))

			return
		}

		value = res.LocalID
	case "supabase":
		res, err := an.supabase.Auth(req.Context(), supabase.AuthRequest{
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		})
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)

			rw.Write([]byte(err.Error()))

			return
		}

		value = res.UserID
	case "pocketbase":
		res, err := an.pocketbase.Auth(req.Context(), pocketbase.AuthRequest{
			Identity: req.FormValue("email"),
			Password: req.FormValue("password"),
		})
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)

			rw.Write([]byte(err.Error()))

			return
		}

		value = res.Record.ID
	default:
		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     cookey,
		Value:    value,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	http.Redirect(rw, req, "/", http.StatusFound)
}

// Patch is verify session
func (an *AuthN) Patch(
	rw http.ResponseWriter,
	req *http.Request,
) {
	if _, err := req.Cookie(cookey); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}
}

// Put is refresh session
func (an *AuthN) Put(
	rw http.ResponseWriter,
	req *http.Request,
) {
	val, err := req.Cookie(cookey)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     cookey,
		Value:    val.Value,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

// Delete is logout
func (an *AuthN) Delete(
	rw http.ResponseWriter,
	req *http.Request,
) {
	http.SetCookie(rw, &http.Cookie{
		Name:   cookey,
		Value:  "",
		MaxAge: -1,
	})
}
