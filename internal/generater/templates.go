package generater

var (
	functionTemplate = `# coding=utf-8
from typing import List, Any
from ekuiper import Function, Context

from {{ .imports }} import {{.functionName}}
{{if .initModel}}from {{ .imports }} import init_model{{end}}

class {{.functionClassName}}(Function):

    def __init__(self):
        {{ if .initModel}}init_model(){{else}}pass{{end}}

    def validate(self, args: List[Any]):
        if len(args) != {{.parasLen}}:
            return "require {{.parasLen}} parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return {{.functionCallName}}

    def is_aggregate(self):
        return {{.isAggr}}


{{.functionWrapperName}} = {{.functionClassName}}()


`

	mainTemplate = `
#  Copyright 2022 EMQ Technologies Co., Ltd.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

from ekuiper import plugin, PluginConfig

{{ range $index, $value := .functionImports -}}
    from {{ $index }} import {{ $value }}
{{end}}

{{ range $index, $value := .sourceImports -}}
    from {{ $index }} import {{ $value }}
{{end}}

{{ range $index, $value := .sinkImports -}}
    from {{ $index }} import {{ $value }}
{{end}}

if __name__ == '__main__':
    funcDict = {
    {{ range $index, $value := .functionImports }}
        "{{$value}}": lambda: {{$value}},
    {{end}}
    }

    sourceDict = {
    {{ range $index, $value := .sourceImports }}
        "{{$index}}": lambda: {{$value}}(),
    {{end}}
    }

    sinkDict = {
    {{ range $index, $value := .sinkImports }}
        "{{$index}}": lambda: {{$value}}(),
    {{end}}
    }

    c = PluginConfig("{{.packageName}}", sourceDict, sinkDict, funcDict)

    plugin.start(c)


`

	installTemplate = `#!/bin/sh

cur=$(dirname "$0")
echo "Base path $cur"
conda install --name {{.env}} --yes --file $cur/requirements.txt
echo "Done"`

	requirementsTemplate = `{{ range $index, $value := .dependencies -}}
    {{ $value }}
{{end}}`
)
