// Copyright 2024 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/lf-edge/ekuiper/internal/conf"
	"github.com/lf-edge/ekuiper/internal/pkg/httpx"
)

type (
	about struct {
		Author      author   `json:"author"`
		Description language `json:"description"`
	}
	author struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Company string `json:"company"`
		Website string `json:"website"`
	}
	language struct {
		English string `json:"en_US"`
		Chinese string `json:"zh_CN"`
	}

	Value struct {
		Value interface{} `json:"value"`
		Label language    `json:"label"`
	}

	Args struct {
		Control string   `json:"control"`
		Type    string   `json:"type"`
		Label   language `json:"label"`
		Values  []Value  `json:"values"`
	}
	//node struct {
	//	Category string   `json:"category"`
	//	Icon     string   `json:"icon"`
	//	Label    language `json:"label"`
	//	Color    string   `json:"color"`
	//}
	FileFunc struct {
		Name        string        `json:"name"`
		Example     string        `json:"example"`
		IsAggregate bool          `json:"aggregate"`
		Hint        language      `json:"hint"`
		Args        []interface{} `json:"args"`
		Outputs     []interface{} `json:"outputs"`
		Node        interface{}   `json:"node"`
	}

	fileMeta struct {
		About     about      `json:"about"`
		Functions []FileFunc `json:"functions"`
	}

	sourceMeta struct {
		About about       `json:"about"`
		Node  interface{} `json:"node"`
	}

	sinkMeta struct {
		About about       `json:"about"`
		Node  interface{} `json:"node"`
	}

	wrapperSink struct {
		Name      string      `json:"name"`
		FilesPath string      `json:"filesPath"`
		ClassName string      `json:"className"`
		Node      interface{} `json:"node"`
	}

	wrapperSource struct {
		Name      string      `json:"name"`
		FilesPath string      `json:"filesPath"`
		ClassName string      `json:"className"`
		Node      interface{} `json:"node"`
	}

	wrapperFunc struct {
		Name          string        `json:"name"`
		Example       string        `json:"example"`
		FilesPath     string        `json:"filesPath"`
		OtherFilePath []string      `json:"otherFilePath"`
		IsAggregate   bool          `json:"aggregate"`
		Hint          language      `json:"hint"`
		Args          []interface{} `json:"args"`
		Outputs       []interface{} `json:"outputs"`
		Node          interface{}   `json:"node"`
		HasInitModel  bool          `json:"initModel"`
	}
	wrapperMeta struct {
		Version        string           `json:"version"`
		PkgName        string           `json:"packagename"`
		About          about            `json:"about"`
		Functions      []*wrapperFunc   `json:"functions"`
		Sources        []*wrapperSource `json:"sources"`
		Sinks          []*wrapperSink   `json:"sinks"`
		Dependencies   []string         `json:"dependencies"`
		VirtualEnvType string           `json:"virtualEnvType"`
		Env            string           `json:"env"`
	}
)

func NewFileFunc(w *wrapperFunc) FileFunc {
	return FileFunc{
		Name:        w.Name,
		Example:     w.Example,
		IsAggregate: w.IsAggregate,
		Hint:        w.Hint,
		Args:        w.Args,
		Outputs:     w.Outputs,
		Node:        w.Node,
	}
}

type PythonCodePackage struct {
	meta           *wrapperMeta
	packageDir     string
	zipDir         string
	pkgname        string
	HostIP         string
	EtcDir         string
	otherFilesPath []string

	sourceFilesPath []string
	functions       functionsWrapper
	sources         sourcesWrapper
	sinks           sinksWrapper
}

type functionsWrapper struct {
	functionsDir           string
	wrapperFileInstanceMap map[string]string
}

type sourcesWrapper struct {
	sourcesDir        string
	sourceInstanceMap map[string]string
}

type sinksWrapper struct {
	sinksDir        string
	sinkInstanceMap map[string]string
}

func newPythonCodePackage(u *wrapperMeta) (*PythonCodePackage, error) {
	p := &PythonCodePackage{
		meta:       u,
		packageDir: "",
		zipDir:     "",
		pkgname:    "",

		functions: functionsWrapper{},
		sinks:     sinksWrapper{},
		sources:   sourcesWrapper{},
	}

	etcDir, err := conf.GetConfLoc()
	if err != nil {
		return nil, err
	}
	p.EtcDir = etcDir
	IP, _ := ExternalIP()
	p.HostIP = IP.String()

	p.pkgname = u.PkgName
	p.packageDir = u.PkgName

	p.functions.functionsDir = path.Join(p.packageDir, "functions")
	_ = os.MkdirAll(p.functions.functionsDir, fs.ModePerm)
	p.sinks.sinksDir = path.Join(p.packageDir, "sinks")
	_ = os.MkdirAll(p.sinks.sinksDir, fs.ModePerm)
	p.sources.sourcesDir = path.Join(p.packageDir, "sources")
	_ = os.MkdirAll(p.sources.sourcesDir, fs.ModePerm)

	p.zipDir = "web/common/static"
	_ = os.MkdirAll(p.zipDir, fs.ModePerm)
	p.functions.wrapperFileInstanceMap = make(map[string]string)
	p.sources.sourceInstanceMap = make(map[string]string)
	p.sinks.sinkInstanceMap = make(map[string]string)
	return p, nil
}

func (p *PythonCodePackage) generateFunctionConfigFile() error {
	for _, f := range p.meta.Functions {
		funcConfig := fileMeta{
			About:     p.meta.About,
			Functions: []FileFunc{NewFileFunc(f)},
		}

		configFilePath := p.functions.functionsDir + "/" + f.Name + ".json"

		data, err := json.Marshal(funcConfig)
		if err != nil {
			return err
		}

		err = os.WriteFile(configFilePath, data, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) genSourcesConfigFile() error {
	for _, f := range p.meta.Sources {
		srcConfig := sourceMeta{
			About: p.meta.About,
			Node:  f.Node,
		}
		configFilePath := p.sources.sourcesDir + "/" + f.Name + ".json"
		data, err := json.Marshal(srcConfig)
		if err != nil {
			return err
		}
		err = os.WriteFile(configFilePath, data, fs.ModePerm)
		if err != nil {
			return err
		}

		yamlFilePath := p.sources.sourcesDir + "/" + f.Name + ".yaml"
		data, err = yaml.Marshal(map[string]interface{}{"default": nil})
		if err != nil {
			return err
		}
		err = os.WriteFile(yamlFilePath, data, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) genSinksConfigFile() error {
	for _, f := range p.meta.Sinks {
		funcConfig := sinkMeta{
			About: p.meta.About,
			Node:  f.Node,
		}
		configFilePath := p.sinks.sinksDir + "/" + f.Name + ".json"
		data, err := json.Marshal(funcConfig)
		if err != nil {
			return err
		}
		err = os.WriteFile(configFilePath, data, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) clean() {
	_ = os.RemoveAll(p.packageDir)
}

func (p *PythonCodePackage) copySourcePythonFile() error {
	for _, v := range p.sourceFilesPath {
		baseName := filepath.Base(v)
		file, err := httpx.ReadFile(v)
		if err != nil {
			return err
		}
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		baseFilePath := "plugins/portable/" + p.packageDir + "/"

		config := map[string]interface{}{
			"BASEPATH": baseFilePath,
		}

		var tp *template.Template = nil

		tp, err = template.New("pythonCodeWrapper").Parse(string(fileContent))
		if err != nil {
			return err
		}
		var output bytes.Buffer
		err = tp.Execute(&output, config)
		if err != nil {
			return err
		}

		configFilePath := p.packageDir + "/" + baseName
		err = os.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) copyOtherFile() error {
	for _, v := range p.otherFilesPath {
		baseName := filepath.Base(v)
		file, err := httpx.ReadFile(v)
		if err != nil {
			return err
		}
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		configFilePath := p.packageDir + "/" + baseName
		err = os.WriteFile(configFilePath, fileContent, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) generateInstallFile(env, tmpl string) error {
	// load the template
	fileContent := tmpl
	config := map[string]interface{}{
		"env": env,
	}
	tp, err := template.New("installScript").Parse(fileContent)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	configFilePath := p.packageDir + "/install.sh"
	err = os.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *PythonCodePackage) generateRequirementFile() error {
	// load the template
	var err error
	fileContent := requirementsTemplate
	u := p.meta

	config := map[string]interface{}{
		"dependencies": u.Dependencies,
	}

	var tp *template.Template = nil

	tp, err = template.New("requirementFile").Parse(fileContent)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	configFilePath := p.packageDir + "/requirements.txt"
	err = os.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *PythonCodePackage) generateMainFile() error {
	// load the template
	var err error
	fileContent := mainTemplate
	u := p.meta

	config := map[string]interface{}{
		"functionImports": p.functions.wrapperFileInstanceMap,
		"sourceImports":   p.sources.sourceInstanceMap,
		"sinkImports":     p.sinks.sinkInstanceMap,
		"packageName":     u.PkgName,
	}

	var tp *template.Template = nil

	tp, err = template.New("mainFile").Parse(fileContent)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	configFilePath := p.packageDir + "/main.py"
	err = os.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *PythonCodePackage) generateZipFile() (string, error) {
	pkgZip := p.zipDir + "/" + p.pkgname + ".zip"
	err := Zip(pkgZip, p.packageDir)
	if err != nil {
		return "", err
	}
	downloadPath := fmt.Sprintf("http://%s:%d/%s", conf.Config.Basic.RestIp, conf.Config.Basic.RestPort, pkgZip)
	return downloadPath, nil
}

func (p *PythonCodePackage) generateJsonConfigFile() error {
	// load the template
	fileContent, err := os.ReadFile(path.Join(p.EtcDir, "templates/function/configPython.json"))
	if err != nil {
		return err
	}
	u := p.meta

	var funcInstances []string
	for _, v := range p.functions.wrapperFileInstanceMap {
		funcInstances = append(funcInstances, v)
	}

	var sourceInstances []string
	for key := range p.sources.sourceInstanceMap {
		sourceInstances = append(sourceInstances, key)
	}

	var sinkInstances []string
	for key := range p.sinks.sinkInstanceMap {
		sinkInstances = append(sinkInstances, key)
	}

	config := map[string]interface{}{
		"functions":      funcInstances,
		"sources":        sourceInstances,
		"sinks":          sinkInstances,
		"version":        u.Version,
		"virtualEnvType": u.VirtualEnvType,
		"env":            u.Env,
	}

	var tp *template.Template = nil

	tp, err = template.New("jsonConfigFile").Parse(string(fileContent))
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	configFilePath := p.packageDir + "/" + u.PkgName + ".json"
	err = os.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (f *wrapperFunc) generateFunctionWrapper(p *PythonCodePackage, tmpl string) error {
	// load the template
	fileContent := tmpl
	var err error

	// get python modules
	var PythonModules string
	baseName := filepath.Base(f.FilesPath)
	if strings.HasSuffix(baseName, ".py") {
		p.sourceFilesPath = append(p.sourceFilesPath, f.FilesPath)
		PythonModules = strings.TrimSuffix(baseName, ".py")
	}

	for _, file := range f.OtherFilePath {
		p.otherFilesPath = append(p.otherFilesPath, file)
	}

	// prepare the config used in template
	wrapperFileName := f.Name + "_wrapper"
	var args []string
	for k := 0; k < len(f.Args); k++ {
		args = append(args, "args["+strconv.Itoa(k)+"]")
	}

	funcCallName := f.Name + "(" + strings.Join(args, ", ") + ")"
	aggStr := ""
	if f.IsAggregate {
		aggStr = "True"
	} else {
		aggStr = "False"
	}
	config := map[string]interface{}{
		"imports":             PythonModules,
		"functionName":        f.Name,
		"functionClassName":   strings.ToUpper(f.Name),
		"functionCallName":    funcCallName,
		"functionWrapperName": wrapperFileName,
		"parasLen":            len(f.Args),
		"isAggr":              aggStr,
		"initModel":           f.HasInitModel,
	}

	var tp *template.Template = nil

	tp, err = template.New("pythonCodeWrapper").Parse(fileContent)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	p.functions.wrapperFileInstanceMap[wrapperFileName] = wrapperFileName

	wrapperPythonPath := p.packageDir + "/" + wrapperFileName + ".py"
	err = os.WriteFile(wrapperPythonPath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}

	f.Example = strings.ReplaceAll(f.Example, f.Name, wrapperFileName)
	f.Name = wrapperFileName

	return nil
}

func (f *wrapperSource) generateSource(p *PythonCodePackage) error {
	p.sourceFilesPath = append(p.sourceFilesPath, f.FilesPath)
	p.sources.sourceInstanceMap[f.Name] = f.ClassName
	return nil
}

func (f *wrapperSink) generateSink(p *PythonCodePackage) error {
	p.sourceFilesPath = append(p.sourceFilesPath, f.FilesPath)
	p.sinks.sinkInstanceMap[f.Name] = f.ClassName
	return nil
}

func PackageSrcCode(data []byte) (string, error) {
	fcs := &wrapperMeta{
		Version:      "",
		PkgName:      "",
		About:        about{},
		Sources:      nil,
		Sinks:        nil,
		Functions:    nil,
		Dependencies: nil,
	}

	err := json.Unmarshal(data, fcs)
	if err != nil {
		return "", err
	}

	pck, err := newPythonCodePackage(fcs)
	if err != nil {
		return "", err
	}

	defer pck.clean()

	if err := generateFunctions(pck); err != nil {
		return "", err
	}

	if err := generateSources(pck); err != nil {
		return "", err
	}

	if err := generateSinks(pck); err != nil {
		return "", err
	}

	err = pck.copySourcePythonFile()
	if err != nil {
		return "", err
	}

	err = pck.copyOtherFile()
	if err != nil {
		return "", err
	}

	err = pck.generateMainFile()
	if err != nil {
		return "", err
	}

	err = pck.generateJsonConfigFile()
	if err != nil {
		return "", err
	}

	err = pck.generateRequirementFile()
	if err != nil {
		return "", err
	}

	err = pck.generateInstallFile(fcs.Env, installTemplate)
	if err != nil {
		return "", err
	}

	return pck.generateZipFile()
}

func generateFunctions(pck *PythonCodePackage) error {
	for _, f := range pck.meta.Functions {
		err := f.generateFunctionWrapper(pck, functionTemplate)
		if err != nil {
			return err
		}
	}
	return pck.generateFunctionConfigFile()
}

func generateSources(pck *PythonCodePackage) error {
	for _, f := range pck.meta.Sources {
		if err := f.generateSource(pck); err != nil {
			return err
		}
	}
	return pck.genSourcesConfigFile()
}

func generateSinks(pck *PythonCodePackage) error {
	for _, f := range pck.meta.Sinks {
		if err := f.generateSink(pck); err != nil {
			return err
		}
	}
	return pck.genSinksConfigFile()
}
