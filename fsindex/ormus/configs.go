package ormus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/util"
)

var (
	datasource              string
	datasys                 string
	saltsize                = 48
	defaultSessionLength, _ = time.ParseDuration("2h")
	unknownclient           = "unknown-client"
	dataLogging             = false
)

// SetDataLogging allows you to turn on or off GORM data logging.
func SetDataLogging(value bool) {
	dataLogging = value
}

// SetDefaults allows a external library to set the local datasource.
// Set saltSize to -1 to persist default.
func SetDefaults(source string, sys string, saltSize int) {
	datasource = source
	datasys = sys
	if saltSize != -1 {
		saltsize = saltSize
	}
	EnsureTableUsers()
	EnsureTableSessions()
}

// returns calculated duration or on error the default session length '2hr'
func durationHrs(hr int) time.Duration {
	if result, err := time.ParseDuration(fmt.Sprintf("%vh", hr)); err != nil {
		return result
	}
	return defaultSessionLength
}

func getClientString(client interface{}) string {

	clistr := ""
	// cess := ""
	switch c := client.(type) {
	case *gin.Context:
		clistr = util.ToUBase64(c.ClientIP())
		break
	case *http.Request:
		clistr = util.ToUBase64(c.RemoteAddr)
		break
	case string:
		clistr = util.ToUBase64(c)
		break
	default:
		clistr = util.ToBase64(unknownclient)
		break
	}
	return clistr
}
