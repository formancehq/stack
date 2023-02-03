package schema

import (
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(d).String())), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {

	if data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid duration")
	}
	data = data[1 : len(data)-1]

	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}
