## Sql Source/Sink Example

```json
{
  "id": "rule1",
  "name": "Test Condition",
  "graph": {
    "nodes": {
      "sqldemo": {
        "type": "source",
        "nodeType": "sql",
        "props": {
          "url": "mysql://root:@127.0.0.1:3306/test?parseTime=true",
          "interval": 1000,
          "templateSqlQueryCfg" : {
              "templateSql":  "select a,b from t limit 1"
          }
        }
      },
      "sqlout": {
        "type": "sink",
        "nodeType": "sql",
        "props": {
           "url": "mysql://root:@127.0.0.1:3306/test?parseTime=true",
          "table": "t2",
          "fields" :["a","b"]
        }
      }
    },
    "topo": {
      "sources": ["sqldemo"],
      "edges": {
        "sqldemo": ["sqlout"]
      }
    }
  }
}
```

## unnest example

```json
{
  "id": "sqlrule3",
  "name": "sqlrule3",
  "graph": {
    "nodes": {
      "mqttdemo": {
        "type": "source",
        "nodeType": "mqtt",
        "props": {
          "server": "tcp://127.0.0.1:1883",
          "datasource": "/test"
        }
      },
      "pickSignal": {
        "type": "operator",
        "nodeType": "pick",
        "props": {
          "fields": ["data[0].signal"]
        }
      },
      "unnestsignal":{
        "type": "operator",
        "nodeType": "srfunc",
        "props": {
          "expr": "unnest(signal) as signal"
        }
      },
      "sqlout": {
        "type": "sink",
        "nodeType": "sql",
        "props": {
          "url": "mysql://root:@127.0.0.1:3306/test?parseTime=true",
          "table": "test_signal",
          "fields" :["signal"]
        }
      }
    },
    "topo": {
      "sources": ["mqttdemo"],
      "edges": {
        "mqttdemo": ["pickSignal"],
        "pickSignal": ["unnestsignal"],
        "unnestsignal": ["sqlout"]
      }
    }
  }
}
```

## window / fold exmaple

```json
{
  "id": "rule",
  "graph": {
    "nodes": {
      "mqttdemo": {
        "type": "source",
        "nodeType": "mqtt",
        "props": {
          "server": "tcp://127.0.0.1:1883",
          "datasource": "/test"
        }
      },
      "window":  {
        "type": "operator",
        "nodeType": "window",
        "props": {
          "type": "countwindow",
          "size": 5
        }
      },
      "fold_into_list":  {
        "type": "operator",
        "nodeType": "aggfunc",
        "props": {
          "expr": "fold_into_list(a) as signal"
        }
      },
      "logout": {
        "type": "sink",
        "nodeType": "log",
        "props": {
        }
      }
    },
    "topo": {
      "sources": ["mqttdemo"],
      "edges": {
        "mqttdemo": ["window"],
        "window": ["fold_into_list"],
        "fold_into_list": ["logout"]
      }
    }
  }
}
```