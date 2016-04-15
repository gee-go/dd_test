# DDLOG
```bash
# Usage
go get github.com/gee-go/ddlog/...

ddlog <log file> # analyze log

ddlog_gen <log file> # generate random log
```

```
Usage of ddlog:
  -alert int
    	Trigger alert when visit count exceeds this number over the given window (default 100)
  -fmt string
    	Log format to parse (default "{remote} {ident} {auth} [{time}] \"{request}\" {status} {size}")
  -plain
    	Use non-fancy output
  -time string
    	Time format to parse (default "02/Jan/2006:15:04:05 -0700")
  -window duration
    	Duration to monitor alert count over (default 2m0s)
```


# TODO
- Improve Testing of fancy ui.

# Future
- Separate agent and viewer. Agent runs in the background, and independent tool visualizes its output.
- Optimize stats, I was lazy and used locks in many places that I could have used atomic ints.
- Rank pages based on EWMA as opposed to all time.
- Remove alert jitter. If load continually alternates above and below the threshold, the current system will display a ton of alerts.
- Calculate metrics based on status code.
