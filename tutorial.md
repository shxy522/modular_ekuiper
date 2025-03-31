## eKuiper 概念

https://ekuiper.org/docs/zh/latest/

1. 了解 eKuiper 流式概念： source configuration/stream/rule/sink
2. Quick Start https://ekuiper.org/docs/zh/latest/installation.html
3. 图规则 https://ekuiper.org/docs/zh/latest/guide/rules/graph_rule.html

常用 operator:

pick/function

source -> operator -> sink

```json
{
  "a": 1
}
```

```json
{
  "type": "operator",
  "nodeType": "function",
  "props": {
    "expr": "a * 2 as a2"
  }
}
````

```json
{
  "a": 1,
  "a2": 2
}
```

mqttsource -> operator -> mqtt sink 动手实操

```json
{
  "a": [1,2,3]
}
```

```json
{
  "type": "operator",
  "nodeType": "function",
  "props": {
    "expr": "normalize(a) as na"
  }
}
````

```json
{
  "a": [1,2,3],
  "na": new value
}
```


中车场景

4. python 插件功能
https://ekuiper.org/docs/zh/latest/extension/portable/overview.html
https://ekuiper.org/docs/zh/latest/extension/portable/python_sdk.html
跑通 pysam.zip 的例子

2. 安装、删除 python 插件 Example

POST /plugins/portables

```json
{
    "name":"pysam",
    "file":"file:///Users/yisa/Downloads/Github/emqx/private/modular_ekuiper/_build/pysam.zip"
}
```

创建 pysam rule

```json
{
  "id": "rule1",
  "name": "Test Condition",
  "graph": {
    "nodes": {
      "pyjsonsource": {
        "type": "source",
        "nodeType": "mqtt",
        "props": {
        }
      },
      "revertoperator": {
        "type": "operator",
        "nodeType": "function",
        "props": {
          "expr": "revert(name) as rname"
        }
      },
      "printout": {
        "type": "sink",
        "nodeType": "mqtt",
        "props": {
        }
      }
    },
    "topo": {
      "sources": ["pyjsonsource"],
      "edges": {
        "pyjsonsource": ["revertoperator"],
        "revertoperator": ["printout"]
      }
    }
  }
}
```

## 常用 REST API

1. 创建、删除、更新、查看 rules
https://ekuiper.org/docs/zh/latest/api/restapi/rules.html

POST /rules

查看规则状态

https://ekuiper.org/docs/zh/latest/api/restapi/rules.html#%E8%8E%B7%E5%8F%96%E8%A7%84%E5%88%99%E7%9A%84%E7%8A%B6%E6%80%81

2. 安装、删除 python 插件 Example

POST /plugins/portables

## 本地编译

仓库: https://github.com/shxy522/modular_ekuiper
可以通过 Makefile 本地编译， make build
最好能掌握打断点调试的方法

cmd/kuiperd/main.go 程序入口

working directory -> make build_prepare

## 常见问题

1. 插件安装后, 部署规则时引用该报错

查看现场环境的 plugins/portable 是否存在对应插件

让中车提供插件后本地安装插件并安装对应规则进行复现

plugins.zip
rule

1. 创建规则失败、并且通常有他们自己的插件

让中车提供插件后本地安装插件并安装规则进行复现调试


```golang
func PlanByGraph(rule *api.Rule) (*topo.Topo, error) {
```

3. 查看日志的方法

目前都是让中车提供

docker container ekuiper -> docker exec -it -> docker cp 

## 定制函数 agg_by_key

代码入口
builtins["agg_by_key"] = builtinFunc{

测试入口
TestAggByKey

example
```json
[
{"a":1,"b":3,"c":5},
{"a":2,"b":4,"c": 6}
]
```

agg_by_key(data,"a,c")

```json
[
[1,2], [5,6]
]
```