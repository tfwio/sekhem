package ormus

import "time"

var (
	datasource              string
	datasys                 string
	defaultSessionLength, _ = time.ParseDuration("2h")
)
