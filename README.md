# go-fanout-timeout
an experiment in timing out concurrent fanout in go

```console
Every 0.5s: go run main.go --num-workers=1000 --work-duration=49 --timeout=50 2>/dev/null

Timeout after 50.034049ms
Processed 992 jobs out of 1000
```
