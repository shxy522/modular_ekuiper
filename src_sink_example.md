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
