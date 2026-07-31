package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d *Duration) Set(value string) error {
	if value == "" {
		return nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("parse duration %q: %w", value, err)
	}
	*d = Duration(parsed)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return err
	}
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		return d.Set(v)
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return fmt.Errorf("invalid duration %s: a bare number is integer nanoseconds; use the string form (e.g. %q) for fractions", data, "1.5s")
		}
		*d = Duration(n)
		return nil
	default:
		return fmt.Errorf("invalid duration: %s", data)
	}
}
