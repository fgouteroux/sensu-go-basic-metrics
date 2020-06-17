# sensu-go-basic-metrics

Get basic os system metrics in graphite format.

## Usage

All metrics collectors support the scheme option.

```
metrics-memory -h
Usage of metrics-memory:
  -scheme string
    	Metric naming scheme, text to prepend to metric.
```

Output:

```
myprefix.memory.total 16700456960 1592300994
myprefix.memory.used 6063837184 1592300994
myprefix.memory.cached 7503912960 1592300994
myprefix.memory.free 2086256640 1592300994
myprefix.memory.active 9296670720 1592300994
myprefix.memory.inactive 3047002112 1592300994
myprefix.memory.swaptotal 4294963200 1592300994
myprefix.memory.swapused 99090432 1592300994
myprefix.memory.swapfree 4195872768 1592300994
```

## References

**Sources**
- [go-osstat](https://github.com/mackerelio/go-osstat)

**Linux Metric**
- cpu
- disk
- disk-usage
- interface
- loadavg
- memory
- netstat-tcp
- uptime
- vmstat

**Windows Metrics**
- memory
- netstat-tcp
- uptime
