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
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/leitstand/leitstand-template-engine/pkg/configen"
	"github.com/leitstand/leitstand-template-engine/pkg/job"
	jobRest "github.com/leitstand/leitstand-template-engine/pkg/job/rest"
	"github.com/leitstand/leitstand-template-engine/pkg/options"
	"github.com/leitstand/leitstand-template-engine/pkg/rest"
	"github.com/leitstand/leitstand-template-engine/pkg/util"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rakyll/statik/fs"
)

// VERSION is the app-global version string, which should be substituted with a
// real value during build.
var VERSION = "UNKNOWN"

type application struct {
	restApplication *rest.Application
	jobApplication  *jobRest.Application
	staticFS        http.FileSystem
}

// @title leitstand-template-engine API
// @version 0.1
// @description _Copyright (C) 2020, RtBrick, Inc._
// @description
// @description Licensed under the Apache License, Version 2.0 (the "License"); you may not
// @description use this file except in compliance with the License.  You may obtain a copy
// @description of the License at
// @description
// @description   http://www.apache.org/licenses/LICENSE-2.0
// @description
// @description Unless required by applicable law or agreed to in writing, software
// @description distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// @description WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// @description License for the specific language governing permissions and limitations under
// @description the License.
// @contact.name Engineering (RtBrick)
// @contact.url http://www.rtbrick.com
// @contact.email eng@rtbrick.com
// @basepath /
// @host
// @tag.name template-engine
// host {{baseurl}}
func main() {
	configFile := flag.String("config", "/etc/rtbrick/leitstand-template-engine/config.json", "Configuration for the leitstand-template-engine")
	serveFromFileSystem := flag.Bool("fs", false, "Serves from filesystem, is only used for development")
	//logging
	debug := flag.Bool("debug", false, "turn on debug logging")
	console := flag.Bool("console", true, "turn on pretty console logging")
	nocolor := flag.Bool("nocolor", true, "turn of color of color output")
	//version
	versionFlag := flag.Bool("version", false, "Returns the software version")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: v%s (built with %s)\n", VERSION, runtime.Version())
		return
	}
	initializeLogger(*debug, *console, *nocolor)

	fileName, _ := filepath.Abs(*configFile)
	opts := &options.Options{}

	err := util.ReadJSONObject(fileName, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("startup error occurred")
	}

	if _, err := os.Stat(opts.TemplatePath); os.IsNotExist(err) {
		log.Fatal().Err(err).Str("folder", opts.TemplatePath).Msg("Folder does not exist")
	}

	// Initialize a new instance of application containing the dependencies.
	jobRepository := job.NewDefaultRepository("/template-engine/api/v1/jobs")
	jobApplication := jobRest.NewApplication(jobRepository)

	configenRepository := configen.NewRepository(opts.TemplatePath)
	restApplication := rest.NewApplication(configenRepository, jobRepository)

	staticFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("startup error occurred")
	}

	app := &application{
		restApplication: restApplication,
		jobApplication:  jobApplication,
		staticFS:        staticFS,
	}

	handler, err := app.routes(*serveFromFileSystem)
	if err != nil {
		log.Fatal().Err(err).Msg("startup error occurred")
	}

	s := &Server{
		Handler: handler,
		Opts:    opts,
	}
	s.ListenAndServe()
}

func initializeLogger(debug, console, nocolor bool) {
	var w io.Writer
	w = os.Stderr
	if console {
		w = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    nocolor,
			TimeFormat: "2006-01-02 15:04:05 MST",
		}
	}

	log.Logger = zerolog.New(w).With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
