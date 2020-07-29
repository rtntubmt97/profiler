package intervalListeners

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

type HttpPage struct {
	port              int
	name              string
	InnerHttpApi      *HttpApis
	CachedStaticPages map[string]*StaticPageHandler
}

func NewHttpPage(port int, name string) *HttpPage {
	ret := new(HttpPage)
	ret.name = name
	ret.port = port
	ret.InnerHttpApi = NewHttpApi(port, name)

	return ret
}

func (httpPage *HttpPage) Serve() *HttpPage {
	server := http.NewServeMux()
	httpPage.SetupHandler(server)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", httpPage.port), server)
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
		FilePath:    "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static/summary.html",
		contentType: "text/html",
	}
	staticPages["/"] = staticPages["/static/summary.html"]
	staticPages["/static/stylesheets/main.css"] = &StaticPageHandler{
		FilePath:    "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static/stylesheets/main.css",
		contentType: "text/css",
	}
	staticPages["/static/js/main.js"] = &StaticPageHandler{
		FilePath:    "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static/js/main.js",
		contentType: "application/javascript",
	}

	httpPage.CachedStaticPages = staticPages
}

func (httpPage *HttpPage) loadCachedHandlers() {
	for _, handler := range httpPage.CachedStaticPages {
		handler.Load()
	}
}

func (httpPage *HttpPage) Listen(profiles map[string]*k.Profile, startTime time.Time) {
	httpPage.InnerHttpApi.Listen(profiles, startTime)
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
	var data []byte
	var err error
	if strings.HasPrefix(handler.FilePath, "http") {
		webRsp, _ := http.Get(handler.FilePath)
		data, err = ioutil.ReadAll(webRsp.Body)
	} else {
		data, err = ioutil.ReadFile(handler.FilePath)
	}
	if err == nil {
		handler.Data = data
	} else {
		fmt.Println(err)
	}
}
