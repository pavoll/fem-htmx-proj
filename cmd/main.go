package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"html/template"
	"io"
)

type Templates struct {
	templates *template.Template
}

type Count struct {
	Count int
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Contacts []Contact

type Data struct {
	Contacts Contacts
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Andrew Boyle", "ab@gmail.com"),
			newContact("Sarah Connor", "sc@gmail.com"),
			newContact("John Connor", "jc@gmail.com"),
			newContact("Winston Smith", "ws@gmail.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func newFormData() FormData {
	return FormData{
		Values: map[string]string{},
		Errors: map[string]string{},
	}
}

func (c *Contacts) hasEmail(email string) bool {
	for _, contact := range *c {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := newPage()
	e.Renderer = NewTemplates()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.Contacts.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			page.Form = formData
			return c.Render(400, "form", page.Form)
		}

		page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))

		return c.Render(200, "index", page)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
