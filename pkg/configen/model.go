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

import "github.com/pkg/errors"

var (
	//ErrEngineNotFound engine not found
	ErrEngineNotFound = errors.New("engine not found")
	//ErrTemplateConfigNotFound engine not found
	ErrTemplateConfigNotFound = errors.New("template config not found")
	//ErrPostProcessorNotFound engine not found
	ErrPostProcessorNotFound = errors.New("post processor not found")
)

// Engine enum
type Engine string

const (
	//EngineGolang ...
	EngineGolang = "golang"
)

// TemplateConfig is the model of template config.
// So each template we want to use have to have this config.
// This config lives in the main_pattern folder under the name *config.json*.
type TemplateConfig struct {
	// TemplateEngine
	TemplateEngine Engine `yaml:"engine" enums:"golang"`
	// MainPattern the name of the templates which are loaded for that generation.
	// (e.g. "*.goyaml") this will be replaced by <basefolder>/<template>/<main_pattern>.
	MainPattern string `yaml:"main_pattern"`
	// IncludePattern the name of the templates which are loaded for that generation.
	// These files are mainly for inclusion and can be shared between multiple templates.
	// (e.g. "includes/*.goyaml") this will be replaced by <basefolder>/<include_pattern>.
	IncludePattern string `yaml:"include_pattern"`
	// MainTemplate name of the template file that is used as entry point for the generation.
	// (e.g. "main.goyaml")
	MainTemplate string `yaml:"main_template"`
	// PostProcessors are used to manipulate the generated code after generation.
	PostProcessors []string `yaml:"post_processors"`
	// Format gives the output format of the template.
	// e.g.: json, json5, text (Default is text)
	// This information is also used to find the correct response Content-Type for the sync restcall.
	OutputFormat string `yaml:"output_format"`
}
