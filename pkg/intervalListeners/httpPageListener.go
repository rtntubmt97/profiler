package intervalListeners

import (
	"fmt"
	"io/ioutil"
	"net/http"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type HttpPage struct {
	InnerHttpApi      *HttpApi
	CachedStaticPages map[string]*StaticPageHandler
}

func NewHttpPage() *HttpPage {
	ret := new(HttpPage)
	ret.InnerHttpApi = NewHttpApi()

	return ret
}

func (httpPage *HttpPage) Serve(port int) *HttpPage {
	server := http.NewServeMux()
	httpPage.SetupHandler(server)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	}()
	return httpPage
}

func (httpPage *HttpPage) SetupHandler(server *http.ServeMux) {
	httpPage.InnerHttpApi.SetupHandler(server)

	httpPage.configCachedHandlers()
	httpPage.loadCachedHandlers()

	for path, handler := range httpPage.CachedStaticPages {
		server.Handle(path, handler)
	}
}

func (httpPage *HttpPage) configCachedHandlers() {
	staticPages := make(map[string]*StaticPageHandler)
	staticPages["/static/summary.html"] = &StaticPageHandler{
		FilePath:    "./web/static/summary.html",
		contentType: "text/html",
	}
	staticPages["/static/stylesheets/main.css"] = &StaticPageHandler{
		FilePath:    "./web/static/stylesheets/main.css",
		contentType: "text/css",
	}

	httpPage.CachedStaticPages = staticPages
}

func (httpPage *HttpPage) loadCachedHandlers() {
	for _, handler := range httpPage.CachedStaticPages {
		handler.Load()
	}
}

func (httpPage *HttpPage) Listen(profiles map[string]*k.Profile) {
	httpPage.InnerHttpApi.Listen(profiles)
}

type StaticPageHandler struct {
	FilePath    string
	Data        []byte
	contentType string
}

func (handler *StaticPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.Load()
	w.Header().Set("Content-Type", handler.contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(handler.Data)
}

func (handler *StaticPageHandler) Load() {
	data, err := ioutil.ReadFile(handler.FilePath)
	if err == nil {
		handler.Data = data
	} else {
		fmt.Println(err)
	}
}
