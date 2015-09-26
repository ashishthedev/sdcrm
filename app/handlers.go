package TODO

import (
	"appengine"
	"fmt"
	"html/template"
	"net/http"
)

type gaeHandler func(c appengine.Context, w http.ResponseWriter, r *http.Request) error

func (h gaeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := h(c, w, r); err != nil {
		c.Errorf("%v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Orders(c appengine.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
func helloworld(c appengine.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Hello world2")
	return nil
}

type urlStruct struct {
	handler      func(w http.ResponseWriter, r *http.Request)
	templatePath string
}

type apiStruct struct {
	handler func(w http.ResponseWriter, r *http.Request)
}

var urlMaps map[string]urlStruct
var apiMaps map[string]apiStruct

var PRE_BUILT_TEMPLATES = make(map[string]*template.Template)
var PAGE_NOT_FOUND_TEMPLATE = template.Must(template.ParseFiles("templates/pageNotFound.html"))

func initDynamicHTMLUrlMaps() {
	http.HandleFunc("/dynamic", apiNotImplementedHandler)
}

func initStaticAdminHTMLUrlMaps() {
	urlMaps := map[string]urlStruct{}

	for path, urlBlob := range urlMaps {
		templatePath := urlBlob.templatePath
		PRE_BUILT_TEMPLATES[path] = template.Must(template.ParseFiles(templatePath))
	}

	for path, urlBlob := range urlMaps {
		http.HandleFunc(path, urlBlob.handler)
	}
	return
}

func initStaticHTMLUrlMaps() {
	http.Handle("/hello", gaeHandler(helloworld))

	urlMaps := map[string]urlStruct{
		"/newOrder": {generalPageHandler, "templates/newOrder.html"},
	}

	for path, urlBlob := range urlMaps {
		templatePath := urlBlob.templatePath
		PRE_BUILT_TEMPLATES[path] = template.Must(template.ParseFiles(templatePath))
	}

	for path, urlBlob := range urlMaps {
		http.HandleFunc(path, urlBlob.handler)
	}
	return
}

func initRootApiMaps() {
	apiMaps := map[string]apiStruct{
		"/a/api/echo": {echoAPIHandler},
		"/a/api/":     {apiNotImplementedHandler},
		"/api/":       {apiNotImplementedHandler},
	}

	for path, apiBlob := range apiMaps {
		http.HandleFunc(path, apiBlob.handler)
	}

	return
}

func init() {
	initRootApiMaps()
	initStaticHTMLUrlMaps()
	initStaticAdminHTMLUrlMaps()
	initDynamicHTMLUrlMaps()
	return
}

func generalPageHandler(w http.ResponseWriter, r *http.Request) {
	t := PRE_BUILT_TEMPLATES[r.URL.Path]
	if t == nil {
		t = PAGE_NOT_FOUND_TEMPLATE
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func apiNotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, r.URL.Path+" not implemented", http.StatusNotImplemented)
	return
}

func echoAPIHandler(w http.ResponseWriter, r *http.Request) {
	WriteJson(&w, "test")
	return
}
