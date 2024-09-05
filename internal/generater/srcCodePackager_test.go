package generater

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//func TestGenerate(t *testing.T) {
//	conf.InitConf()
//	file, err := ioutil.ReadFile("testdata/test.json")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fcs := &wrapperFuncs{
//		Version:      "",
//		About:        about{},
//		Functions:    nil,
//		Dependencies: nil,
//	}
//
//	err = json.Unmarshal(file, fcs)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	pck, _ := newPythonCodePackage(fcs)
//
//	for _, f := range pck.funcMeta.Functions {
//		err := f.generateFunctionWrapper(pck)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//	}
//
//	err = pck.generateFunctionConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copySourcePythonFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateMainFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateJsonConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateRequirementFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateInstallFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	dwnPath, err := pck.generateZipFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fmt.Printf("##########%s\n", dwnPath)
//
//	pck.clean()
//
//}
//
//func TestGenerate2(t *testing.T) {
//	conf.InitConf()
//	file, err := ioutil.ReadFile("testdata2/test.json")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fcs := &wrapperFuncs{
//		Version:      "",
//		About:        about{},
//		Functions:    nil,
//		Dependencies: nil,
//	}
//
//	err = json.Unmarshal(file, fcs)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	pck, _ := newPythonCodePackage(fcs)
//
//	for _, f := range pck.funcMeta.Functions {
//		err := f.generateFunctionWrapper(pck)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//	}
//
//	err = pck.generateFunctionConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copySourcePythonFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copyOtherFile()
//
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateMainFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateJsonConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateRequirementFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateInstallFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	dwnPath, err := pck.generateZipFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fmt.Printf("##########%s\n", dwnPath)
//
//	pck.clean()
//
//}
//
//func TestGenerate3(t *testing.T) {
//	conf.InitConf()
//	file, err := ioutil.ReadFile("testdata3/test.json")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fcs := &wrapperFuncs{
//		Version:      "",
//		About:        about{},
//		Functions:    nil,
//		Dependencies: nil,
//	}
//
//	err = json.Unmarshal(file, fcs)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	pck, _ := newPythonCodePackage(fcs)
//
//	for _, f := range pck.funcMeta.Functions {
//		err := f.generateFunctionWrapper(pck)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//	}
//
//	err = pck.generateFunctionConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copySourcePythonFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copyOtherFile()
//
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateMainFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateJsonConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateRequirementFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateInstallFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	dwnPath, err := pck.generateZipFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fmt.Printf("##########%s\n", dwnPath)
//
//	pck.clean()
//
//}
//
//func TestGenerate4(t *testing.T) {
//	conf.InitConf()
//	file, err := ioutil.ReadFile("testdata4/sgoSmooth_lib_crrc.json")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fcs := &wrapperFuncs{
//		Version:      "",
//		About:        about{},
//		Functions:    nil,
//		Dependencies: nil,
//	}
//
//	err = json.Unmarshal(file, fcs)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	pck, _ := newPythonCodePackage(fcs)
//
//	for _, f := range pck.funcMeta.Functions {
//		err := f.generateFunctionWrapper(pck)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//	}
//
//	err = pck.generateFunctionConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copySourcePythonFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.copyOtherFile()
//
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateMainFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateJsonConfigFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateRequirementFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	err = pck.generateInstallFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	dwnPath, err := pck.generateZipFile()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	fmt.Printf("##########%s\n", dwnPath)
//
//	pck.clean()
//
//}

func TestInstallScript(t *testing.T) {
	p := &PythonCodePackage{}
	p.EtcDir = "."
	p.packageDir = "."
	require.NoError(t, p.generateInstallFile("active", condaInstallTemplate))
	c, err := os.ReadFile("./install.sh")
	require.NoError(t, err)
	result := `#!/bin/sh

cur=$(dirname "$0")
echo "Base path $cur"
conda install --name active --yes --file $cur/requirements.txt
echo "Done"`
	require.Equal(t, result, string(c))
	os.Remove("./install.sh")
}

func TestFunctionWrapper(t *testing.T) {
	w := &wrapperFunc{}
	p := &PythonCodePackage{}
	p.EtcDir = "."
	p.packageDir = "."
	p.functions.wrapperFileInstanceMap = make(map[string]string)
	w.Name = "apply_butter_filter"
	w.HasInitModel = true
	w.FilesPath = "testdata/butterFilter.py"
	w.Args = []interface{}{
		1, 2, 3, 4, 5,
	}
	require.NoError(t, w.generateFunctionWrapper(p, functionTemplate))
	result := `# coding=utf-8
from typing import List, Any
from ekuiper import Function, Context

from butterFilter import apply_butter_filter
from butterFilter import init_model

class APPLY_BUTTER_FILTER(Function):

    def __init__(self):
        init_model()

    def validate(self, args: List[Any]):
        if len(args) != 5:
            return "require 5 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return apply_butter_filter(args[0], args[1], args[2], args[3], args[4])

    def is_aggregate(self):
        return False


apply_butter_filter_wrapper = APPLY_BUTTER_FILTER()


`
	c, err := os.ReadFile("./apply_butter_filter_wrapper.py")
	require.NoError(t, err)
	require.Equal(t, result, string(c))
	os.Remove("./apply_butter_filter_wrapper.py")
}
