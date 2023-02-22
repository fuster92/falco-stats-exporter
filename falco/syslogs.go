package falco

import (
	"encoding/json"
	"fmt"
)

type Status struct {
	Events      int `json:"events"`
	Drops       int `json:"drops"`
	Preemptions int `json:"preemptions"`
}

type SystemLogLine struct {
	Sample int    `json:"sample"`
	Cur    Status `json:"cur"`
	Delta  Status `json:"delta"`
}

func ParseSingleLine(line string) (SystemLogLine, error) {
	falcoLogLine := SystemLogLine{}
	err := json.Unmarshal([]byte(line), &falcoLogLine)
	if err != nil {
		return SystemLogLine{}, fmt.Errorf("couldn't parse log line: %v", err)
	}
	return falcoLogLine, nil
}
