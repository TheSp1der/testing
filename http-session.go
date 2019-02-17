package main

import (
	"strconv"
	"time"

	"net/http"

	"github.com/TheSp1der/goerror"
	"github.com/gofrs/uuid"
)

var users = map[string]map[string]string{
	"john": {
		"password":      "baspass",
		"first_name":    "John",
		"last_name":     "Dorian",
		"session_token": "",
	},
	"paul": {
		"password":      "changeme",
		"first_name":    "Paul",
		"last_name":     "Brandaniowitz",
		"session_token": "",
	},
}

func webListener(port int) {
	ws := http.NewServeMux()

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      ws,
		ReadTimeout:  time.Duration(15 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(120 * time.Second),
	}

	ws.HandleFunc("/signin", signIn)
	ws.HandleFunc("/", rootLanding)

	if err := srv.ListenAndServe(); err != nil {
		goerror.Fatal(err)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var (
			login        bool
			sessionToken string
		)

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for k, v := range users {
			if k == r.FormValue("username") && r.FormValue("password") == v["password"] {
				login = true
				t, _ := uuid.NewV4()
				sessionToken = t.String()
				v["session_token"] = sessionToken
			}
		}

		if login {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: time.Now().Add(120 * time.Second),
			})
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		w.Write([]byte("Nope, just Nope!"))
	}
}

func rootLanding(w http.ResponseWriter, r *http.Request) {
	var (
		nocookie bool
		userData map[string]string
	)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			nocookie = true
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	if !nocookie {
		for _, v := range users {
			if cookie.Value == v["session_token"] {
				userData = map[string]string{
					"first_name": v["first_name"],
					"last_name":  v["last_name"],
					"login":      "true",
				}
				break
			}
		}

		if userData["login"] == "true" {
			output := `<!doctype html>`
			output += `<html lang="en">`
			output += `	<head>`
			output += `		<title></title>`
			output += `		<meta http-equiv="X-UA-Compatible" content="IE=edge">`
			output += `	</head>`
			output += `	<body>`
			output += `		<span style="font-weight: bold;">hi</span><br>`
			output += `	</body>`
			output += `</html>`

			w.Write([]byte(output))
			return
		}
	}

	output := `<!doctype html>`
	output += `<html lang="en">`
	output += `	<head>`
	output += `		<title></title>`
	output += `		<meta http-equiv="X-UA-Compatible" content="IE=edge">`
	output += `	</head>`
	output += `	<body>`
	output += `		<form action="/signin" method="post">`
	output += `			Username: <input name="username" type="text"><br>`
	output += `			Password: <input name="password" type="text"><br>`
	output += `			<input type="submit" value="Submit">`
	output += `		</form>`
	output += `	</body>`
	output += `</html>`
	w.Write([]byte(output))
}

func main() {
	webListener(80)
}
