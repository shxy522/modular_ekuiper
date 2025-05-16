## Kafka Source/Sink Example

```json
{
  "id": "rule1",
  "name": "Test Condition",
  "graph": {
    "nodes": {
      "kafkademo": {
        "type": "source",
        "nodeType": "kafka",
        "props": {
          "brokers": "127.0.0.1:29092",
          "datasource": "test9"
        }
      },
      "kafkaout": {
        "type": "sink",
        "nodeType": "kafka",
        "props": {
           "brokers": "127.0.0.1:29092",
          "topic": "test10"
        }
      }
    },
    "topo": {
      "sources": ["kafkademo"],
      "edges": {
        "kafkademo": ["kafkaout"]
      }
    }
  }
}
```

Kakfa Source 配置项参考: https://ekuiper.org/docs/zh/latest/guide/sources/plugin/kafka.html

Kafka Sink 配置项参考: https://ekuiper.org/docs/zh/latest/guide/sinks/plugin/kafka.html

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
              "templateSql":  "select * from t limit 1"
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

SQL Source 配置项参考:  https://ekuiper.org/docs/zh/latest/guide/sources/plugin/sql.html
SQL SINK 配置项参考: https://ekuiper.org/docs/zh/latest/guide/sinks/plugin/sql.html

## influx

```json
{
  "id": "rule1",
  "name": "Test Condition",
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
      "influxout": {
        "type": "sink",
        "nodeType": "influx",
        "props": {
          "addr": "http://localhost:8086",
          "username": "",
          "password": "",
          "measurement": "test",
          "databasename": "mydb",
          "tags": "{\"tag1\":\"value1\"}",
          "fields": ["a", "b"]
        }
      }
    },
    "topo": {
      "sources": ["mqttdemo"],
      "edges": {
        "mqttdemo": ["influxout"]
      }
    }
  }
}
```

influx 配置项参考： https://ekuiper.org/docs/zh/latest/guide/sinks/plugin/influx.html

## influx2

```json
{
  "id": "rule1",
  "name": "Test Condition",
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
      "influx2out": {
        "type": "sink",
        "nodeType": "influx2",
        "props": {
          "addr": "http://localhost:8086",
          "token": "test_token",
          "org": "admin",
          "measurement": "test",
          "bucket": "admin",
          "tags": "{\"tag1\":\"value1\"}",
          "fields": ["a", "b"]
        }
      }
    },
    "topo": {
      "sources": ["mqttdemo"],
      "edges": {
        "mqttdemo": ["influx2out"]
      }
    }
  }
}
```

influxv2 配置项参考： https://ekuiper.org/docs/zh/latest/guide/sinks/plugin/influx2.html

## tdengine v3

```json
{
  "id": "rule1",
  "name": "Test Condition",
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
      "tdengineout": {
        "type": "sink",
        "nodeType": "tdengine",
        "props": {
          "host": "hostname",
          "port": 6030,
          "database": "dab",
          "table": "{{.table}}",
          "tsfieldname": "ts",
          "fields": [
            "f1",
            "f2"
          ],
          "sTable": "myStable",
          "tagFields": [
            "f3",
            "f4"
          ]
        }
      }
    },
    "topo": {
      "sources": [
        "mqttdemo"
      ],
      "edges": {
        "mqttdemo": [
          "tdengineout"
        ]
      }
    }
  }
}
```

tdengines 配置项： https://ekuiper.org/docs/zh/v1.14/guide/sinks/plugin/tdengine.html

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
