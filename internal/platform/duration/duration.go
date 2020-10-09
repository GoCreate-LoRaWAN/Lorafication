// Package duration exposes a custom Duration type that can be used on struct
// fields meant to be marshaled and unmarshaled from JSON and YAML files.
package duration

import (
	"encoding/json"
	"errors"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration is a type that can be used in place of time.Duration with
// the added benefit of JSON marshaling and unmarshaling.
type Duration struct {
	time.Duration
}

// MarshalJSON implements the interface needed to marshal the duration
// type into a JSON structure.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements the interface needed to unmarshal the Duration
// type from a JSON structure to the Duration type.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// MarshalYAML implements the interface needed to marshal the duration
// type into a YAML structure.
func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// UnmarshalYAML implements the interface needed to unmarshal the duration
// from a YAML representation into the Duration type.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var err error

	d.Duration, err = time.ParseDuration(value.Value)
	if err != nil {
		return err
	}

	return nil
}

func (d *Duration) IsEmpty() bool {
	return d.Duration.Milliseconds() == 0
}
