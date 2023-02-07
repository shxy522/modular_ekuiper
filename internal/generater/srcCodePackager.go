// copyright 2021 EMQ Technologies Co., Ltd.
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
	"github.com/lf-edge/ekuiper/internal/conf"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
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
	fileFuncMeta struct {
		About     about      `json:"about"`
		Functions []FileFunc `json:"functions"`
	}

	wrapperFunc struct {
		Name        string        `json:"name"`
		Example     string        `json:"example"`
		FilesPath   string        `json:"filesPath"`
		SoPath      string        `json:"soPath"`
		IsAggregate bool          `json:"aggregate"`
		Hint        language      `json:"hint"`
		Args        []interface{} `json:"args"`
		Outputs     []interface{} `json:"outputs"`
		Node        interface{}   `json:"node"`
	}
	wrapperFuncs struct {
		Version      string         `json:"version"`
		PkgName      string         `json:"packagename"`
		About        about          `json:"about"`
		Functions    []*wrapperFunc `json:"functions"`
		Dependencies []string       `json:"dependencies"`
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
	funcMeta               *wrapperFuncs
	packageDir             string
	zipDir                 string
	functionsDir           string
	pkgname                string
	HostIP                 string
	wrapperFileInstanceMap map[string]string
	sourceFilesPath        []string
	soFilesPath            []string
	EtcDir                 string
}

func newPythonCodePackage(u *wrapperFuncs) (*PythonCodePackage, error) {
	p := &PythonCodePackage{
		funcMeta:     u,
		packageDir:   "",
		zipDir:       "",
		functionsDir: "",
		pkgname:      "",
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

	p.functionsDir = path.Join(p.packageDir, "functions")
	_ = os.MkdirAll(p.functionsDir, fs.ModePerm)
	p.zipDir = "web/common/static"
	_ = os.MkdirAll(p.zipDir, fs.ModePerm)
	p.wrapperFileInstanceMap = make(map[string]string)
	return p, nil
}

func (p *PythonCodePackage) generateFunctionConfigFile() error {
	for _, f := range p.funcMeta.Functions {
		funcConfig := fileFuncMeta{
			About:     p.funcMeta.About,
			Functions: []FileFunc{NewFileFunc(f)},
		}

		configFilePath := p.functionsDir + "/" + f.Name + ".json"

		data, err := json.Marshal(funcConfig)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(configFilePath, data, fs.ModePerm)
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
		fileContent, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		configFilePath := p.packageDir + "/" + baseName
		err = ioutil.WriteFile(configFilePath, fileContent, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) copySoFile() error {
	for _, v := range p.soFilesPath {
		baseName := filepath.Base(v)
		fileContent, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		configFilePath := p.packageDir + "/" + baseName
		err = ioutil.WriteFile(configFilePath, fileContent, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PythonCodePackage) generateInstallFile() error {
	//load the template
	fileContent, err := ioutil.ReadFile(path.Join(p.EtcDir, "templates/function/install.sh"))
	if err != nil {
		return err
	}
	configFilePath := p.packageDir + "/install.sh"
	err = ioutil.WriteFile(configFilePath, fileContent, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *PythonCodePackage) generateRequirementFile() error {
	//load the template
	fileContent, err := ioutil.ReadFile(path.Join(p.EtcDir, "templates/function/requirements.tmpl"))
	if err != nil {
		return err
	}
	u := p.funcMeta

	config := map[string]interface{}{
		"dependencies": u.Dependencies,
	}

	var tp *template.Template = nil

	tp, err = template.New("requirementFile").Parse(string(fileContent))
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	configFilePath := p.packageDir + "/requirements.txt"
	err = ioutil.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *PythonCodePackage) generateMainFile() error {
	//load the template
	fileContent, err := ioutil.ReadFile(path.Join(p.EtcDir, "templates/function/main.tmpl"))
	if err != nil {
		return err
	}
	u := p.funcMeta

	config := map[string]interface{}{
		"imports":     p.wrapperFileInstanceMap,
		"packageName": u.PkgName,
	}

	var tp *template.Template = nil

	tp, err = template.New("mainFile").Parse(string(fileContent))
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = tp.Execute(&output, config)
	if err != nil {
		return err
	}

	configFilePath := p.packageDir + "/main.py"
	err = ioutil.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
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
	//load the template
	fileContent, err := ioutil.ReadFile(path.Join(p.EtcDir, "templates/function/configPython.json"))
	if err != nil {
		return err
	}
	u := p.funcMeta

	var funcInstances []string
	for _, v := range p.wrapperFileInstanceMap {
		funcInstances = append(funcInstances, v)
	}

	config := map[string]interface{}{
		"functions": funcInstances,
		"version":   u.Version,
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
	err = ioutil.WriteFile(configFilePath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (f *wrapperFunc) generateFunctionWrapper(p *PythonCodePackage) error {
	//load the template
	fileContent, err := ioutil.ReadFile(path.Join(p.EtcDir, "templates/function/functionPython.tmpl"))
	if err != nil {
		return err
	}

	// get python modules
	var PythonModules string
	baseName := filepath.Base(f.FilesPath)
	if strings.HasSuffix(baseName, ".py") {
		p.sourceFilesPath = append(p.sourceFilesPath, f.FilesPath)
		PythonModules = strings.TrimSuffix(baseName, ".py")
	}

	baseName = filepath.Base(f.SoPath)
	if strings.HasSuffix(baseName, ".so") {
		p.soFilesPath = append(p.soFilesPath, f.SoPath)
	}

	//prepare the config used in template
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

	p.wrapperFileInstanceMap[wrapperFileName] = wrapperFileName

	wrapperPythonPath := p.packageDir + "/" + wrapperFileName + ".py"
	err = ioutil.WriteFile(wrapperPythonPath, output.Bytes(), fs.ModePerm)
	if err != nil {
		return err
	}

	f.Example = strings.ReplaceAll(f.Example, f.Name, wrapperFileName)
	f.Name = wrapperFileName

	return nil
}

func PackageSrcCode(data []byte) (string, error) {
	fcs := &wrapperFuncs{
		Version:      "",
		PkgName:      "",
		About:        about{},
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

	for _, f := range pck.funcMeta.Functions {
		err := f.generateFunctionWrapper(pck)
		if err != nil {
			return "", err
		}
	}

	err = pck.generateFunctionConfigFile()
	if err != nil {
		return "", err
	}

	err = pck.copySourcePythonFile()
	if err != nil {
		return "", err
	}

	err = pck.copySoFile()
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

	err = pck.generateInstallFile()
	if err != nil {
		return "", err
	}

	return pck.generateZipFile()
}
