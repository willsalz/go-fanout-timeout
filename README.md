# go-fanout-timeout
an experiment in timing out concurrent fanout in go

```console
$ watch -n 0.5 "go run main.go --num-workers=1000 --work-duration=49 --timeout=50 2>/dev/null"
```
