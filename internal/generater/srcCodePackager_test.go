package generater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGenerate(t *testing.T) {
	file, err := ioutil.ReadFile("testdata/test.json")
	if err != nil {
		t.Error(err)
		return
	}

	fcs := &wrapperFuncs{
		Version:      "",
		About:        about{},
		Functions:    nil,
		Dependencies: nil,
	}

	err = json.Unmarshal(file, fcs)
	if err != nil {
		t.Error(err)
		return
	}

	pck := newPythonCodePackage(fcs)

	for _, f := range pck.funcMeta.Functions {
		err := f.generateFunctionWrapper(pck)
		if err != nil {
			t.Error(err)
			return
		}
	}

	err = pck.generateFunctionConfigFile()
	if err != nil {
		t.Error(err)
		return
	}

	err = pck.copySourcePythonFile()
	if err != nil {
		t.Error(err)
		return
	}

	err = pck.generateMainFile()
	if err != nil {
		t.Error(err)
		return
	}

	err = pck.generateJsonConfigFile()
	if err != nil {
		t.Error(err)
		return
	}

	err = pck.generateRequirementFile()
	if err != nil {
		t.Error(err)
		return
	}

	err = pck.generateInstallFile()
	if err != nil {
		t.Error(err)
		return
	}

	dwnPath, err := pck.generateZipFile()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("##########%s\n", dwnPath)

	pck.clean()

}
