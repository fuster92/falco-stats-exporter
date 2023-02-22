package falco

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var singleFalcoLine = "{\"sample\": 2971, \"cur\": {\"events\": 48342737, \"drops\": 131484, \"preemptions\": 0}, \"delta\": {\"events\": 8890, \"drops\": 0, \"preemptions\": 0}, \"drop_pct\": 0}"
var singleFalcoLineWithError = "{\"sample\": 2971, \"cur\": {\"events\": 48342737, \"drops\": 131484, \"preemptions\": 0}, \"delta\": {\"events\": 8890, \"drops\": 0, \"preemptions\": 0}, \"drop_pct\": 0},"

func TestParseSingleLine(t *testing.T) {

	falcoMetricsLine, err := ParseSingleLine(singleFalcoLine)
	assert.NoError(t, err)
	assert.Equal(t, 8890, falcoMetricsLine.Delta.Events)
	assert.Equal(t, 2971, falcoMetricsLine.Sample)
	assert.Equal(t, 131484, falcoMetricsLine.Cur.Drops)
}
