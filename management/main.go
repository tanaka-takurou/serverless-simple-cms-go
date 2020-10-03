package main

import (
	"io"
	"os"
	"log"
	"bytes"
	"context"
	"net/http"
	"html/template"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type PageData struct {
	Title   string
	ApiPath string
	ProductUrl string
}

type Response events.APIGatewayProxyResponse

const title string = "CMS Management Page"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	tmp := template.New("tmp")
	var dat PageData
	p := request.PathParameters
	q := request.QueryStringParameters
	page := p["proxy"]
	if len(page) < 1 {
		page = q["page"]
	}
	buf := new(bytes.Buffer)
	fw := io.Writer(buf)
	dat.Title = title
	dat.ApiPath = os.Getenv("API_PATH")
	dat.ProductUrl = os.Getenv("FRONT_URL")
	if page == "signup" {
		tmp = getTemplate("templates/login/signup.html", 2)
	} else if page == "activation" {
		tmp = getTemplate("templates/login/activation.html", 1)
	} else if page == "changepass" {
		tmp = getTemplate("templates/login/changepass.html", 2)
	} else if page == "top" {
		tmp = getTemplate("templates/top.html", 0)
	} else if page == "add" {
		tmp = getTemplate("templates/add.html", 0)
	} else if page == "fix" {
		tmp = getTemplate("templates/fix.html", 0)
	} else if page == "setting" {
		tmp = getTemplate("templates/setting.html", 0)
	} else if page == "dynamodb" {
		tmp = getTemplate("templates/dynamodb.html", 0)
	} else if page == "s3" {
		tmp = getTemplate("templates/s3.html", 0)
	} else if page == "js" {
		tmp = getTemplate("templates/js.html", 0)
	} else if page == "css" {
		tmp = getTemplate("templates/css.html", 0)
	} else {
		tmp = getTemplate("templates/login/login.html", 1)
	}
	if e := tmp.ExecuteTemplate(fw, "base", dat); e != nil {
		log.Fatal(e)
	}
	res := Response{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            string(buf.Bytes()),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
	return res, nil
}

func getTemplate(file string, baseType int) *template.Template {
	templates := template.New("templates")
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	}
	if baseType == 0 {
		templates = template.Must(template.New("templates").Funcs(funcMap).ParseFiles(
			"templates/view.html",
			"templates/index.html",
			"templates/header.html",
			file,
		))
	} else if baseType == 1 {
		templates = template.Must(template.New("templates").Funcs(funcMap).ParseFiles(
			"templates/view.html",
			"templates/login/index.html",
			"templates/login/header.html",
			file,
		))
	} else {
		templates = template.Must(template.New("templates").Funcs(funcMap).ParseFiles(
			"templates/view.html",
			"templates/login/index.html",
			"templates/login/header.html",
			"templates/login/passwarning.html",
			file,
		))
	}
	return templates
}

func main() {
	lambda.Start(HandleRequest)
}
