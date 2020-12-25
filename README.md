# Nomad Prometheus Exporter

### Exported Metrics

| Metric | Meaning | Labels |
| ------ | ------- | ------ |
| nomad_task_status | Get status of job | job, taskgroup, task |

### Job status codes

```go
	statuses := map[string]float64{
		"running":  1,
		"pending":  2,
		"complete": 3,
		"failed":   4,
		"lost":     5,
	}
```

### Query Examples

```
count(nomad_task_status == 2) # count of pending tasks
count(nomad_task_status != 1) # count of non-running tasks
```

### Nomad Client Env Varables

- The Nomad client environment variables can be found in the [client source](https://github.com/hashicorp/nomad/blob/master/api/api.go#L243)
- ADDR=:5577
- RECORD_COMPLETED

### Docker

```
docker run --rm --network host -it alexes/nomad-exporter
```
