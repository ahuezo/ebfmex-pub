package site

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "net/http"
	"fmt"
)

func init() {
    http.HandleFunc("/backend", GaeLogin)
    http.HandleFunc("/listausuarios", ListaUsuarios)
    http.HandleFunc("/listasesiones", ListaSesiones)
    http.HandleFunc("/test", test)
}

func test(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    if u := user.Current(c); u != nil {
		if ck, err := r.Cookie("ebfmex-pub-sessid-ua"); err == nil {
			fmt.Fprintf(w, "Nombre: %q, Valor: %q\n", ck.Name, ck.Value);
		} else {
	     http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cr, err := r.Cookie("ebfmex-pub-sesscontrol-ua"); err == nil {
			fmt.Fprintf(w, "Nombre: %q, Valor: %q\n", cr.Name, cr.Value);
		} else {
	     http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListaUsuarios(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	/* Verifica si el usuario es interno */
    if u := user.Current(c); u != nil {
		q := datastore.NewQuery("Cta").Order("-FechaHora").Limit(10)
		usuarios := make([]Cta, 0, 10)
		if _, err := q.GetAll(c, &usuarios); err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := listUsersTpl.Execute(w, usuarios); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ListaSesiones(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	/* Verifica si el usuario es interno */
    if u := user.Current(c); u != nil {
		q := datastore.NewQuery("Sess").Limit(10)
		s := make([]Sess, 0, 10)
		if _, err := q.GetAll(c, &s); err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := listSessTpl.Execute(w, s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GaeLogin(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/* Autenticación de usuario interno */
	u := user.Current(c)
    if u == nil {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Location", url)
        w.WriteHeader(http.StatusFound)
        return
    }
    http.Redirect(w, r, "/listausuarios", http.StatusFound)
}

