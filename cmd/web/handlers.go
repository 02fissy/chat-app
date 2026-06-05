package main

import (
	"fmt"
	"net/http"
    "chat.fisayo.net/internal/validator"
    "chat.fisayo.net/internal/models"
    "errors"
)
type userSignupForm struct {
    Name                string `form:"name"`
    Phone             string `form:"phone"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}
type userLoginForm struct {
     Name                string `form:"name"`
    Phone             string `form:"phone"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}

func (app *application) serveHome(w http.ResponseWriter, r *http.Request) { 
    data := app.newTemplateData(r)
	err := app.render(w, r, http.StatusOK, "chat.html", data)
    if err != nil {
        app.logger.Error(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
func (app *application) home(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
	err := app.render(w, r, http.StatusOK, "home.html", data)
    if err != nil {
        app.logger.Error(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
func (app *application) userSignup(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = userSignupForm{}
	err := app.render(w, r, http.StatusOK, "signup.html", data)
    if err != nil {
        app.logger.Error(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
     var form userSignupForm
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }
    form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Phone), "phone", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
    form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
        return
    }
    err = app.users.Insert(form.Name, form.Phone, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrDuplicatePhone) {
            form.AddFieldError("phone", "Phone number is already in use")
            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
        } else {
            app.serverError(w, r, err)
        }
        return
    }
    app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
    data.Form = userLoginForm{}
    app.render(w, r, http.StatusOK, "login.html", data)
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
    var form userLoginForm
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
        return
    }
    id, err := app.users.Authenticate(form.Name, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            form.AddNonFieldError("Name or password is incorrect")
            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
        } else {
            app.serverError(w, r, err)
        }
        return
    }
    err = app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.serverError(w, r, err)
        return
    }
    app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
    http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Log out user")
}
