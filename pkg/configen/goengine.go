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

package configen

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"

	"github.com/rs/zerolog/log"
)

// GoEngine Template Engine
type GoEngine struct {
}

// Ensure, that GoEngine does implement TemplateEngine.
var _ TemplateEngine = &GoEngine{}

//newGoEngine creates a new code generation repository
func newGoEngine() (*GoEngine, error) {
	return &GoEngine{}, nil
}

// GenerateFile executes a template and adds a variable set
func (r GoEngine) GenerateFile(config *TemplateConfig, data map[string]interface{}) ([]byte, string, error) {
	t := template.New("base").Funcs(sprig.TxtFuncMap())
	templates, err := t.ParseGlob(config.MainPattern)
	if err != nil {
		return nil, "", err
	}
	if len(config.IncludePattern) > 0 {
		templates, err = templates.ParseGlob(config.IncludePattern)
	}
	if err != nil {
		return nil, "", err
	}
	result, err := r.executeTemplate(config.MainTemplate, templates, data)
	return result, config.OutputFormat, err
}

func (r *GoEngine) executeTemplate(templateName string, template *template.Template, data interface{}) ([]byte, error) {
	log.Debug().Str("template_name", templateName).Msg("Execute")
	var tpl bytes.Buffer
	err := template.ExecuteTemplate(&tpl, templateName, data)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return tpl.Bytes(), nil
}
