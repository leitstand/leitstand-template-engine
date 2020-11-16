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
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/tidwall/pretty"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Repository to generate files via templates
type Repository struct {
	path string
}

// NewRepository creates a new code generation repository
func NewRepository(templatePath string) *Repository {
	return &Repository{path: templatePath}
}

// TemplateEngine allow to generate files
type TemplateEngine interface {
	// GenerateFile executes a template and added a variable set
	GenerateFile(config *TemplateConfig, data map[string]interface{}) ([]byte, string, error)
}

func (r *Repository) newEngine(name Engine) (TemplateEngine, error) {
	switch name {
	case EngineGolang:
		return newGoEngine()
	case "":
		return newGoEngine()
	}
	return nil, errors.WithMessage(ErrEngineNotFound, string(name))
}

// GenerateFile executes a template and added a variable set
func (r Repository) GenerateFile(templateFolder string, variables map[string]interface{}) ([]byte, string, error) {
	config, err := parseConfigFile(r.path, templateFolder)
	if err != nil {
		return nil, "", err
	}
	engine, err := r.newEngine(config.TemplateEngine)
	if err != nil {
		return nil, "", err
	}
	result, format, err := engine.GenerateFile(config, variables)
	if err != nil {
		return result, format, err
	}
	for _, postProcessorName := range config.PostProcessors {
		processor, ok := postProcessors[postProcessorName]
		if !ok {
			return nil, "", errors.WithMessage(ErrPostProcessorNotFound, postProcessorName)
		}
		result, err = processor(result)
		if err != nil {
			return result, format, err
		}
	}
	return result, format, err
}
func parseConfigFile(templatePath string, templateFolder string) (*TemplateConfig, error) {
	configFile := fmt.Sprintf("%s/%s/config.yaml", templatePath, templateFolder)
	configFileFD, err := os.Open(configFile)
	if err != nil {
		log.Error().Err(err).Str("config_file", configFile).Msg("not able to read config file for templating")
		return nil, errors.WithMessage(ErrTemplateConfigNotFound, configFile)
	}
	log.Debug().Str("config_file", configFile).Msg("successfully opened")
	defer func() { _ = configFileFD.Close() }()

	byteValue, err := ioutil.ReadAll(configFileFD)
	if err != nil {
		return nil, err
	}
	config := &TemplateConfig{}
	err = yaml.Unmarshal(byteValue, config)
	if err != nil {
		return nil, err
	}
	if len(config.MainPattern) > 0 {
		config.MainPattern = fmt.Sprintf("%s/%s/%s", templatePath, templateFolder, config.MainPattern)
	}
	if len(config.IncludePattern) > 0 {
		config.IncludePattern = fmt.Sprintf("%s/%s", templatePath, config.IncludePattern)
	}
	return config, err
}

type postProcess func(value []byte) ([]byte, error)

var (
	postProcessors = map[string]postProcess{
		"removeTrailingCommas": removeTrailingCommas,
		"removeEmptyLines":     removeEmptyLines,
		"prettyJSON":           prettyJSON,
		"uglyJSON":             uglyJSON,
	}
	removeTrailingCommasPattern = regexp.MustCompile(`(\s*),(\s*[}\]])`)
	removeEmptyLinesPattern     = regexp.MustCompile(`\n(( )*\n)+`)
)

func removeTrailingCommas(value []byte) ([]byte, error) {
	return []byte(removeTrailingCommasPattern.ReplaceAllString(string(value), "$1$2")), nil
}

func removeEmptyLines(value []byte) ([]byte, error) {
	return []byte(removeEmptyLinesPattern.ReplaceAllString(string(value), "\n")), nil
}

func prettyJSON(value []byte) ([]byte, error) {
	out := pretty.Pretty(value)
	return out, nil
}

func uglyJSON(value []byte) ([]byte, error) {
	out := pretty.UglyInPlace(value)
	return out, nil
}
