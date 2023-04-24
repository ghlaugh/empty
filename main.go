package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"empty/routers"
	"github.com/urfave/negroni"
	"github.com/15125505/zlog/log"
	_ "empty/controllers" // 这个是为了保证这个包被引入，否则不会调用init函数
)

func init() {
	log.Log.SetLogFile("logs/empty")
	log.Log.SetFileColor(true)
	log.Log.SetAdditionalErrorFile(true)
	log.Log.SetLogLevel(log.LevelDebug)
}

func main() {
	r := mux.NewRouter()
	routers.CreateHandle(r)
	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("public")), negroni.HandlerFunc(routers.PreProcess))
	n.UseHandler(r)
	n.Run(":15000")
}
