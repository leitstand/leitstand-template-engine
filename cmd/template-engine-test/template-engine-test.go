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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/leitstand/leitstand-template-engine/pkg/configen"
	"github.com/leitstand/leitstand-template-engine/pkg/util"

	"github.com/google/go-cmp/cmp"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func main() {
	templatePath := flag.String("templatePath", ".", "Template main folder")
	template := flag.String("template", "", "Template name")
	test := flag.String("test", "", "Name of the test")
	skipValidation := flag.Bool("n", false, "skips the validation")
	//logging
	debug := flag.Bool("debug", false, "turn on debug logging")
	console := flag.Bool("console", true, "turn on pretty console logging")
	flag.Parse()

	if *templatePath == "" || *template == "" || *test == "" {
		flag.Usage()
		return
	}
	initializeLogger(debug, console)
	r := configen.NewRepository(*templatePath)

	var variables map[string]interface{}
	variablesFile := fmt.Sprintf("%s/%s/%s_variables.json", *templatePath, *template, *test)
	err := util.ReadJSONObject(variablesFile, &variables)
	if err != nil {
		log.Error().Err(err).Msg("can't find variables file")
		return
	}

	got, format, err := r.GenerateFile(*template, variables)
	if err != nil {
		log.Error().Err(err).Msg("error in generation")
		return
	}

	gotFile := fmt.Sprintf("%s/%s/%s_got.%s", *templatePath, *template, *test, format)
	_ = ioutil.WriteFile(gotFile, got, os.ModePerm)

	log.Info().Msgf("Wrote file %s!", gotFile)

	if !*skipValidation {
		resultFile := fmt.Sprintf("%s/%s/%s_result.%s", *templatePath, *template, *test, format)
		want, err := ioutil.ReadFile(resultFile)
		if err != nil {
			log.Error().Err(err).Msg("can't find result file")
			return
		}

		switch format {
		case "json":
			if diff := cmp.Diff(transformToJSONObject(want), transformToJSONObject(got)); diff != "" {
				log.Warn().Msgf("mismatch (-want +got):\n%s", diff)
				return
			}
		case "json5":
			if diff := cmp.Diff(transformToJSON5Object(want), transformToJSON5Object(got)); diff != "" {
				log.Warn().Msgf("mismatch (-want +got):\n%s", diff)
				return
			}
		case "txt":
			fallthrough
		default:
			if diff := cmp.Diff(string(want), string(got)); diff != "" {
				log.Warn().Msgf("mismatch (-want +got):\n%s", diff)
				return
			}
		}
		log.Info().Msg("Success!")
	}
}
func initializeLogger(debug, console *bool) {
	var w io.Writer
	w = os.Stderr
	if *console {
		w = zerolog.ConsoleWriter{Out: os.Stderr}
	}
	log.Logger = zerolog.New(w).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
func transformToJSONObject(value []byte) interface{} {
	var v interface{}
	if err := json.Unmarshal(value, &v); err != nil {
		return fmt.Sprintf("not parseable (%s)", value) // use unparseable input as the output
	}
	return v
}
func transformToJSON5Object(value []byte) interface{} {
	var v interface{}
	if err := json5.Unmarshal(value, &v); err != nil {
		return value // use unparseable input as the output
	}
	return v
}
