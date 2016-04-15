# DDLOG

# TODO
- Improve Testing of fancy ui.

# Future
- Separate agent and viewer. Agent runs in the background, and independent tool visualizes its output.
- Optimize stats, I was lazy and used locks in many places that I could have used atomic ints.
- Rank pages based on EWMA as opposed to all time.
- Remove alert jitter. If load continually alternates above and below the threshold, the current system will display a ton of alerts.
- Calculate metrics based on status code.
