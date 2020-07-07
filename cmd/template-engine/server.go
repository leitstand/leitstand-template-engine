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
	"strings"
	"time"

	"github.com/leitstand/leitstand-template-engine/pkg/options"

	"github.com/rs/zerolog/log"
)

// Server for http or https
type Server struct {
	Handler http.Handler
	Opts    *options.Options
}

// ListenAndServe starts the server in listening mode
func (s *Server) ListenAndServe() {
	s.serveHTTP()
}

func (s *Server) serveHTTP() {
	addr := s.Opts.HTTPAddress

	log.Info().Str("listen_addr", addr).
		Msg("listening on")

	server := &http.Server{
		Addr:              addr,
		Handler:           s.Handler,
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		log.Error().Err(err).Msg("http.ListenAndServe()")
	}
	log.Info().Str("listen_addr", addr).
		Msg("closing")
}
