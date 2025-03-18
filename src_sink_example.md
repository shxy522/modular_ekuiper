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