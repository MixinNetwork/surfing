package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MixinNetwork/surfing/durable"
	"github.com/MixinNetwork/surfing/middlewares"
	"github.com/MixinNetwork/surfing/routes"
	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/handlers"
	"github.com/unrolled/render"
)

func StartHTTP(db *durable.Database) error {
	router := httptreemux.New()
	routes.RegisterHanders(router)
	routes.RegisterRoutes(router)
	handler := middlewares.Authenticate(router)
	handler = middlewares.Constraint(handler)
	handler = middlewares.Context(handler, db, render.New())
	handler = handlers.ProxyHeaders(handler)

	log.Println("http service running at: 7001")
	return http.ListenAndServe(fmt.Sprintf(":%d", 7001), handler)
}
