/*
 * Copyright 2020 RtBrick Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.  You may obtain a copy
 * of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

package main

import (
	"net/http"

	"github.com/leitstand/leitstand-template-engine/pkg/requestlog"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	_ "github.com/leitstand/leitstand-template-engine/pkg/statik"
)

func (app *application) routes(serveFromFileSystem bool) (http.Handler, error) {
	router := mux.NewRouter()
	router.NewRoute().Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/public/", http.StatusTemporaryRedirect)
	})
	if !serveFromFileSystem {
		staticServer := http.FileServer(app.staticFS)
		router.NewRoute().Name("public").PathPrefix("/public/").Handler(http.StripPrefix("/public", staticServer))
	} else {
		fileServer := http.FileServer(http.Dir("./web/src"))
		router.PathPrefix("/public").Handler(http.StripPrefix("/public", fileServer))
	}

	app.jobApplication.Routes("/template-engine/api/v1", router)
	app.restApplication.Routes(router)
	_ = app.printAllRoutes(router)
	loggedRouter := requestlog.NewHandler(new(requestLogger), router)
	return loggedRouter, nil
}

func (app *application) printAllRoutes(router *mux.Router) error {
	return router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		pathRegexp, _ := route.GetPathRegexp()
		queriesTemplates, _ := route.GetQueriesTemplates()
		queriesRegexps, _ := route.GetQueriesRegexp()
		methods, _ := route.GetMethods()
		log.Debug().Str("route", pathTemplate).
			Str("path_regexp", pathRegexp).
			Strs("queries_templates", queriesTemplates).
			Strs("queries_regexps", queriesRegexps).
			Strs("methods", methods).Msgf("route %s", pathTemplate)

		return nil
	})
}

type requestLogger struct {
}

// Log method takes the request log an sends it to the log infrastructure
func (cl *requestLogger) Log(le *requestlog.Entry) {
	log.Info().
		Time("received_time", le.ReceivedTime).
		Str("method", le.RequestMethod).
		Str("url", le.RequestURL).
		Int64("header_size", le.RequestHeaderSize).
		Int64("body_size", le.RequestBodySize).
		Str("agent", le.UserAgent).
		Str("referer", le.Referer).
		Str("proto", le.Proto).
		Str("remote_ip", le.RemoteIP).
		Str("server_ip", le.ServerIP).
		Int("status", le.Status).
		Int64("resp_header_size", le.ResponseHeaderSize).
		Int64("resp_body_size", le.ResponseBodySize).
		Dur("latency", le.Latency).
		Msg("")
}
