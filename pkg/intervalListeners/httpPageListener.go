package intervalListeners

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	k "github.com/rtntubmt97/profiler/pkg/kernel"
)

var resourcePrefix string

func init() {
	resourcePrefix = "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static"
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Download resource via internet...")
		return
	}
	currentPath := path.Dir(filename)
	moduleRelatePath := "/pkg/intervalListeners"
	prjPath := currentPath[0 : len(currentPath)-len(moduleRelatePath)]
	_, err := os.Stat(prjPath + "/web/pivot.txt")
	if err == nil {
		resourcePrefix = prjPath + "/web/static"
	}
}

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
		FilePath: resourcePrefix + "/summary.html",
		// FilePath: "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static/summary.html",
		// FilePath:    "web/static/summary.html",
		contentType: "text/html",
	}
	staticPages["/"] = staticPages["/static/summary.html"]
	staticPages["/static/stylesheets/main.css"] = &StaticPageHandler{
		FilePath: resourcePrefix + "/stylesheets/main.css",
		// FilePath: "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static/stylesheets/main.css",
		// FilePath:    "web/static/stylesheets/main.css",
		contentType: "text/css",
	}
	staticPages["/static/js/main.js"] = &StaticPageHandler{
		FilePath: resourcePrefix + "/js/main.js",
		// FilePath: "https://raw.githubusercontent.com/rtntubmt97/profiler/master/web/static/js/main.js",
		// FilePath:    "web/static/js/main.js",
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
