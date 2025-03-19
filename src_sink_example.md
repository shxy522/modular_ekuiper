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
