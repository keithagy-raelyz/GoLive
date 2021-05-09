package app

import (
	"GoLive/cache"
	"GoLive/db"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Data struct {
	User      db.User
	Merchant  db.Merchant
	Merchants []db.Merchant
	Error     Error
	Products  []db.Product
	Cart      cache.Cart
}
type Error struct {
	ErrMsg string
}

func (a *App) displayLogin(w http.ResponseWriter, r *http.Request) {
	// Check session; if already logged in then redirect to home page
	if _, alreadyLoggedIn := a.HaveValidSessionCookie(r); alreadyLoggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// ExecuteTemplate
	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	data := Data{}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) validateUserLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	// Need to hash password
	fmt.Println(username)
	fmt.Println(password)
	foundUser, err := a.db.GetUser(username)
	if err != nil {
		fmt.Println(err, "error in get user")
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
		if err != nil {
			log.Fatal(err)
		}
		data := Data{Error: Error{"Piece of shit who can't even remember your password / username"}}
		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if correct := PWcompare(password, foundUser.Password); !correct {
		fmt.Println(foundUser.Password, "not correct")
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/loginBody.html", "templates/error.html")
		if err != nil {
			log.Fatal(err)
		}
		data := Data{Error: Error{"Piece of shit who can't even remember your password / username"}}
		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	newsessionKey := uuid.NewV4().String()
	newsession := cache.NewUserSession(
		newsessionKey,
		time.Now().Add(cache.SessionLife),
		foundUser,
		nil)
	a.cacheManager.AddtoCache(newsession)

	t, err := template.ParseFiles("templates/base.html", "templates/footer.html", "templates/navbar.html", "templates/body.html", "templates/error.html")
	if err != nil {
		log.Fatal(err)
	}
	p, _ := a.db.GetAllProducts()
	data := Data{User: foundUser, Products: p}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (a *App) validateMerchantLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	// Need to hash password

	foundMerch, err := a.db.GetMerchant(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		return
	}
	if correct := PWcompare(password, foundMerch.User.Password); !correct {
		w.WriteHeader(http.StatusForbidden)
		//w.Write([]byte("403 - Invalid Login Credentials"))
		return
	}

	newsessionKey := uuid.NewV4().String()
	newsession := cache.NewMerchSession(newsessionKey,
		time.Now().Add(cache.SessionLife),
		foundMerch)
	a.cacheManager.AddtoCache(newsession)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	sessionKey, err := r.Cookie("sessionCookie")
	sessionKeyVal := sessionKey.String()
	params := mux.Vars(r)
	// Verify valid user type
	userType := params["userType"]
	if err != nil {
		// No Session Cookie
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	a.cacheManager.RemoveFromCache(sessionKeyVal, userType)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (a *App) allUser(w http.ResponseWriter, r *http.Request) {

	users, err := a.db.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Get All Users OK"))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Verify valid merchant ID
	userID := params["userid"]

	// Merchant ID supplied
	// Show all products under merchID; if invalid merchID handle error
	// TODO GetInventory errors need to be multiplexed to differentiate between invalid merchant and empty store
	var u db.User
	u, err := a.db.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		// Invalid merchant ID inputted
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No user for input USERID"))
		return
	}

	fmt.Println(u, "printing user")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Valid merchant ID, displaying store"))
}

func (a *App) postUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := strings.ToLower(r.FormValue("email"))
	pw1 := r.FormValue("pw1")
	pw2 := r.FormValue("pw2")

	var u db.User

	u.Name = username
	u.Email = email
	if u.Name == "" {
		//TODO proper error handling
		log.Fatal()
	}
	if u.Email == "" {
		//TODO proper error handling
		log.Fatal()
	}
	err := a.db.CheckUser(u)
	if err != nil {
		//send the err msg back (err = errmsg)
		fmt.Println(err, "CHeckUser error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}
	if pw1 != pw2 {
		//t.ParseFiles("./templates/errorRegister.html")
		//data := Data{nil, ErrorMsg{"color:red", "Password entered are different"}, ErrorMsg{"display:none", ""}, ErrorMsg{"display:block", ""}, 0, nil}
		//t.Execute(w, data)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Repeated PW"))
		return
	}
	err = a.db.CreateUser(u, pw1)
	if err != nil {
		//send the err msg back (err = errmsg)
		fmt.Println(err, "createuser error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("201 - User Creation Successful"))
	//TODO setcookie here thanks
	http.Redirect(w, r, "/", 201)
}

func (a *App) putUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status Not Found"))
		// TODO display error message serve a proper template redirecting to registry of all merchants
		return
	}

	username := r.URL.Query().Get("Username")
	if username == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Username cannot be empty"))
		return
	}

	email := r.URL.Query().Get("Email")
	if email == "" {
		//TODO proper error handling
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Email cannot be empty"))
		return
	}
	var u db.User
	u.Name = username
	u.Email = email
	u.Id = userID
	err := a.db.UpdateUser(u)
	if err != nil {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Updated Successfully"))
}

func (a *App) delUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, ok := params["userid"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Status not Found"))
		return
	}
	err := a.db.DeleteUser(userID)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		//TODO proper error handling in template
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Unprocessable Entity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Deleted Successfully"))
}

//func (a *App) changeUserPw(w http.ResponseWriter, r *http.Request){
//	//TODO pw change
//	pw1:= r.URL.Query().Get("Pw1")
//	if pw1== "" {
//		//TODO proper error handling
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("400 - Password cannot be empty"))
//		return
//	}
//	pw2:= r.URL.Query().Get("Pw2")
//	if pw2== "" {
//		//TODO proper error handling
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("400 - Password cannot be empty"))
//		return
//	}
//	if pw1 != pw2 {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("400 - Password must be different"))
//		return
//	}
//}
