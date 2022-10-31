package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal_web_day9/connection"
	"personal_web_day9/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	// route path public folder
	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//route path folder upload
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")

	route.HandleFunc("/form-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", middleware.UploadFile(AddProject)).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	route.HandleFunc("/form-register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/formEditProject/{id}", formEditProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", edit).Methods("POST")

	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running on port 3307")
	http.ListenAndServe("localhost:3307", route)
}

type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = SessionData{}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	IsLogin  bool
}

type Project struct {
	Id              int
	Title           string
	Desc            string
	Image           string
	StartDate       time.Time
	EndDate         time.Time
	StartDateFormat string
	EndDateFormat   string
	Duration        string
	Author          string
	Html            string
	HtmlIcon        string
	Css             string
	CssIcon         string
	Javascript      string
	JavascriptIcon  string
	Bootstrap       string
	BootstrapIcon   string
	IsLogin         bool
}

var dataProject = []Project{}

func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Conten-type", "text/html; charset=utf8")
	var tmplt, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	// FOR STORING SESSION COOKIES
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["isLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["isLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	// QUERY GET DATA FROM DB
	data, err := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, image, description, duration, html_icon, css_icon, js_icon, bs_icon, tb_user.name as author FROM public.tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id")
	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	var result []Project

	for data.Next() {
		var each = Project{}
		// SCAN PROCESS, LINKING EACH DATA FROM DB WITH STRUCT
		err := data.Scan(&each.Id, &each.Title, &each.Image, &each.Desc, &each.Duration, &each.HtmlIcon, &each.CssIcon, &each.JavascriptIcon, &each.BootstrapIcon, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	respondData := map[string]interface{}{
		"DataSession": Data,
		"Projects":    result,
	}

	w.WriteHeader(http.StatusOK)
	tmplt.Execute(w, respondData)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf8")
	var tmplt, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	tmplt.Execute(w, nil)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf8")
	var tmplt, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	tmplt.Execute(w, nil)
}

func AddProject(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)

	}

	title := r.PostForm.Get("inputTitle")
	desc := r.PostForm.Get("inputDesc")
	startDate := r.PostForm.Get("inputStartDate")
	endDate := r.PostForm.Get("inputEndDate")
	html := (r.FormValue("inputHtml"))
	htmlIcon := ""
	if html == "htmlchecked" {
		htmlIcon = "fa-brands fa-html5"
	}
	css := (r.FormValue("inputCss"))
	cssIcon := ""
	if css == "csschecked" {
		cssIcon = "fa-brands fa-css3"
	}
	js := (r.FormValue("inputJs"))
	jsIcon := ""
	if js == "jschecked" {
		jsIcon = "fa-brands fa-js"
	}
	bootstrap := (r.FormValue("inputBootstrap"))
	bootstrapIcon := ""
	if bootstrap == "bootstrapchecked" {
		bootstrapIcon = "fa-brands fa-bootstrap"
	}

	layout := "2006-01-02"
	startDateParse, _ := time.Parse(layout, startDate)
	endDateParse, _ := time.Parse(layout, endDate)

	hour := 1
	day := hour * 24
	week := hour * 24 * 7
	month := hour * 24 * 30
	year := hour * 24 * 365

	diffHour := endDateParse.Sub(startDateParse).Hours()
	var diffHours int = int(diffHour)
	// fmt.Println(diffHour)

	days := diffHours / day
	weeks := diffHours / week
	months := diffHours / month
	years := diffHours / year

	var duration string

	if int(diffHours) < week {
		duration = strconv.Itoa(int(days)) + " Day(s)"
	} else if int(diffHours) < month {
		duration = strconv.Itoa(int(weeks)) + " Week(s)"
	} else if int(diffHours) < year {
		duration = strconv.Itoa(int(months)) + " Month(s)"
	} else if int(diffHours) > year {
		duration = strconv.Itoa(int(years)) + " Year(s)"
	}

	// context dari file upload
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	author := session.Values["Id"].(int)
	fmt.Println(author)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(title, description, image, author_id, start_date, end_date, duration, html_icon, css_icon, js_icon, bs_icon) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ", title, desc, image, author, startDate, endDate, duration, htmlIcon, cssIcon, jsIcon, bootstrapIcon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	fmt.Println(dataProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Conten-type", "text/html; charset=utf8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tmplt, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message: " + err.Error()))
		return
	}

	ProjectDetail := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, title, image, description, start_date, end_date, duration, html_icon, css_icon, js_icon, bs_icon, tb_user.name as author FROM public.tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id WHERE tb_project.id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Image, &ProjectDetail.Desc, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.HtmlIcon, &ProjectDetail.CssIcon, &ProjectDetail.JavascriptIcon, &ProjectDetail.BootstrapIcon, &ProjectDetail.Author)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail.StartDateFormat = ProjectDetail.StartDate.Format("02 Jan 2006")
	ProjectDetail.EndDateFormat = ProjectDetail.EndDate.Format("02 Jan 2006")

	response := map[string]interface{}{
		"Project": ProjectDetail,
	}
	w.WriteHeader(http.StatusOK)
	tmplt.Execute(w, response)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM public.tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message: " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmplt, err = template.ParseFiles("views/form-register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message: " + err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	tmplt.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("inputName")
	email := r.PostForm.Get("inputEmail")
	password := r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message" + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.AddFlash("Succes Registered, Log In!", "message")
	session.Save(r, w)
	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)

}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-type", "text/html; charset=utf8")
	var tmplt, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		w.Write([]byte("message: " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["isLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["isLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {
			flashes = append(flashes, f1.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	tmplt.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("inputEmail")
	password := r.PostForm.Get("inputPassword")

	user := User{}

	// EMAIL CHECKING
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM public.tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Email belum terdaftar!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	fmt.Println(user)

	// PASSWORD CHECKING
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Password Salah!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	// FOR STORING DATA IN SESSION
	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["Id"] = user.Id //  as a relation between tb_user and tb_project
	session.Values["isLogin"] = true
	session.Options.MaxAge = 10800 // 10800 seconds = 3 hours

	session.AddFlash("Login Success, Go To Home Page", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	var tmplt, err = template.ParseFiles("views/edit-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectEdit := Project{}

	// QUERY GET DATA PROJECT FROM DB BY ID
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, image, description, start_date, end_date, html_icon, css_icon, js_icon, bs_icon FROM public.tb_project WHERE id = $1", id).Scan(&ProjectEdit.Id, &ProjectEdit.Title, &ProjectEdit.Image, &ProjectEdit.Desc, &ProjectEdit.StartDate, &ProjectEdit.EndDate, &ProjectEdit.HtmlIcon, &ProjectEdit.CssIcon, &ProjectEdit.JavascriptIcon, &ProjectEdit.BootstrapIcon)

	ProjectEdit.StartDateFormat = ProjectEdit.StartDate.Format("2006-01-02")
	ProjectEdit.EndDateFormat = ProjectEdit.EndDate.Format("2006-01-02")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	ProjectEdit.StartDateFormat = ProjectEdit.StartDate.Format("2006-01-02")
	ProjectEdit.EndDateFormat = ProjectEdit.EndDate.Format("2006-01-02")

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

	tmplt.Execute(w, response)
}

func edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	inputTitle := r.PostForm.Get("inputTitle")
	inputStartdate := r.PostForm.Get("inputStartDate")
	inputEnddate := r.PostForm.Get("inputEndDate")
	inputDeskripsi := r.PostForm.Get("inputDesc")
	html := (r.FormValue("inputHtml"))
	htmlIcon := ""
	if html == "htmlchecked" {
		htmlIcon = "fa-brands fa-html5"
	}
	css := (r.FormValue("inputCss"))
	cssIcon := ""
	if css == "csschecked" {
		cssIcon = "fa-brands fa-css3"
	}
	js := (r.FormValue("inputJs"))
	jsIcon := ""
	if js == "jschecked" {
		jsIcon = "fa-brands fa-js"
	}
	bootstrap := (r.FormValue("inputBootstrap"))
	bootstrapIcon := ""
	if bootstrap == "bootstrapchecked" {
		bootstrapIcon = "fa-brands fa-bootstrap"
	}

	// dataContext := r.Context().Value("dataFile")
	// inputGambar := dataContext.(string)

	layout := "2006-01-02"
	startDateParse, _ := time.Parse(layout, inputStartdate)
	endDateParse, _ := time.Parse(layout, inputEnddate)

	hour := 1
	day := hour * 24
	week := hour * 24 * 7
	month := hour * 24 * 30
	year := hour * 24 * 365

	diffHour := endDateParse.Sub(startDateParse).Hours()
	var diffHours int = int(diffHour)

	days := diffHours / day
	weeks := diffHours / week
	months := diffHours / month
	years := diffHours / year

	var duration string

	if int(diffHours) < week {
		duration = strconv.Itoa(int(days)) + " Day(s)"
	} else if int(diffHours) < month {
		duration = strconv.Itoa(int(weeks)) + " Week(s)"
	} else if int(diffHours) < year {
		duration = strconv.Itoa(int(months)) + " Month(s)"
	} else if int(diffHours) > year {
		duration = strconv.Itoa(int(years)) + " Year(s)"
	}

	// dataContext := r.Context().Value("dataFile")
	// Image := dataContext.(string)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title = $1, description = $2, start_date = $3, end_date = $4, duration = $5, html_icon=$6, css_icon=$7, js_icon=$8, bs_icon=$9 WHERE id = $10", inputTitle, inputDeskripsi, startDateParse, endDateParse, duration, htmlIcon, cssIcon, jsIcon, bootstrapIcon, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
