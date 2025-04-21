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

```shell
mysql> show create table test_signal;
+-------------+--------------------------------------------------------------------------------------------------------+
| Table       | Create Table                                                                                           |
+-------------+--------------------------------------------------------------------------------------------------------+
| test_signal | CREATE TABLE `test_signal` (
  `data_signal` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 |
+-------------+--------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)

```

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
        "nodeType": "function",
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

dataInput

```json
{
  "data":[
    {
      "signal": [1.0,2.0,3.0]
    }
  ]
}
```