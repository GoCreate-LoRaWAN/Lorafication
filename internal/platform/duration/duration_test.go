// Package duration_test tests the duration package.
package duration_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/22arw/lorafication/internal/platform/duration"
	"gopkg.in/yaml.v3"
)

// Foo is a type that incorporates the Duration type to be tested.
type Foo struct {
	Bar duration.Duration `json:"bar"`
}

// TestDuration_MarshalJSON tests the MarshalJSON receiver function of
// the duration type.
func TestDuration_MarshalJSON(t *testing.T) {
	t.Parallel()

	foo := Foo{
		Bar: duration.Duration{
			Duration: time.Second * 5,
		},
	}

	b, err := json.Marshal(foo)
	if err != nil {
		t.Fatalf("marshal duration type: %v", err)
	}

	if e, a := `{"bar":"5s"}`, string(b); e != a {
		t.Errorf("expected marshaled JSON to be \"%s\", got \"%s\"", e, a)
	}
}

// TestDuration_UnmarshalJSON tests the UnmarshalJSON receiver function
// of the duration type.
func TestDuration_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var unmarshaled Foo
	if err := json.Unmarshal([]byte(`{"bar":"5s"}`), &unmarshaled); err != nil {
		t.Fatalf("unmarshal duration type: %v", err)
	}

	if e, a := time.Second*5, unmarshaled.Bar.Duration; e != a {
		t.Errorf("expected unmarhsaled JSON duration to be %v, got %v", e, a)
	}
}

// TestDuration_MarshalYAML tests the MarshalYAML receiver function of
// the duration type.
func TestDuration_MarshalYAML(t *testing.T) {
	t.Parallel()

	foo := Foo{
		Bar: duration.Duration{
			Duration: time.Second * 5,
		},
	}

	b, err := yaml.Marshal(foo)
	if err != nil {
		t.Fatalf("marshal duration type: %v", err)
	}

	if e, a := `bar: 5s`, strings.TrimSpace(string(b)); e != a {
		t.Errorf("expected marshaled YAML to be \"%s\", got \"%s\"", e, a)
	}
}

// TestDuration_UnmarshalYAML tests the UnmarshalYAML receiver function of
// the duration type.
func TestDuration_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	var unmarshaled Foo
	if err := yaml.Unmarshal([]byte(`bar: 5s`), &unmarshaled); err != nil {
		t.Fatalf("unmarshal duration type: %v", err)
	}

	if e, a := time.Second*5, unmarshaled.Bar.Duration; e != a {
		t.Errorf("expected unmarhsaled YAML duration to be %v, got %v", e, a)
	}
}
