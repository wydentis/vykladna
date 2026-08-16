package core_utils_env

import (
	"errors"
	"testing"
	"time"
)

func TestProcess_Basics(t *testing.T) {
	type Config struct {
		Name    string        `env:"NAME"`
		Port    int           `env:"PORT,required"`
		Debug   bool          `env:"DEBUG,,true"`
		Timeout time.Duration `env:"TIMEOUT,,5s"`
		NoTag   string
		Skipped string `env:"-"`
	}

	t.Setenv("PORT", "8080")
	t.Setenv("NAME", "svc")

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Name != "svc" || c.Port != 8080 || c.Debug != true || c.Timeout != 5*time.Second {
		t.Fatalf("unexpected config: %+v", c)
	}
}

func TestProcess_RequiredMissing(t *testing.T) {
	type Config struct {
		Port int `env:"PORT_MISSING_XYZ,required"`
	}
	var c Config
	err := Process(&c, "")
	if !errors.Is(err, ErrRequiredFieldNoInfo) {
		t.Fatalf("expected ErrRequiredFieldNoInfo, got %v", err)
	}
}

func TestProcess_NestedStruct(t *testing.T) {
	type DB struct {
		Host string `env:"HOST,,localhost"`
	}
	type Config struct {
		DB DB `env:"DB"`
	}
	t.Setenv("DB_HOST", "example.com")

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.DB.Host != "example.com" {
		t.Fatalf("unexpected DB.Host: %q", c.DB.Host)
	}
}

func TestProcess_PointerToNestedStruct(t *testing.T) {
	type DB struct {
		Host string `env:"HOST,,localhost"`
	}
	type Config struct {
		DB *DB `env:"DB"`
	}

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.DB == nil || c.DB.Host != "localhost" {
		t.Fatalf("expected allocated *DB with default host, got %+v", c.DB)
	}
}

type customLevel int

func (l *customLevel) UnmarshalText(b []byte) error {
	switch string(b) {
	case "debug":
		*l = 0
	case "info":
		*l = 1
	case "warn":
		*l = 2
	default:
		return errors.New("unknown level: " + string(b))
	}
	return nil
}

func TestProcess_TextUnmarshalerEnvValue(t *testing.T) {
	type Config struct {
		Level customLevel `env:"LEVEL"`
	}
	t.Setenv("LEVEL", "warn")

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Level != 2 {
		t.Fatalf("expected Level=2, got %d", c.Level)
	}
}

func TestProcess_TextUnmarshalerDefault(t *testing.T) {
	// Regression test: default values must go through the same
	// TextUnmarshaler path as real env values, not raw strconv parsing.
	type Config struct {
		Level customLevel `env:"LEVEL_UNSET_XYZ,,info"`
	}

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Level != 1 {
		t.Fatalf("expected default Level=1, got %d", c.Level)
	}
}

func TestProcess_TimeTimeDefault(t *testing.T) {
	// Regression test: time.Time is a struct that implements
	// TextUnmarshaler; a non-empty default must not error as
	// "unsupported type: struct".
	type Config struct {
		StartedAt time.Time `env:"STARTED_AT_XYZ,,2024-01-02T15:04:05Z"`
	}

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	want := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	if !c.StartedAt.Equal(want) {
		t.Fatalf("expected %v, got %v", want, c.StartedAt)
	}
}

func TestProcess_ExplicitEmptyStringDefault(t *testing.T) {
	// Regression test: an explicit empty default ("NAME,,") must be
	// distinguishable from "no default" and actually apply.
	type Config struct {
		Name string `env:"NAME_UNSET_XYZ,,"`
	}

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Name != "" {
		t.Fatalf("expected empty string default, got %q", c.Name)
	}
}

func TestProcess_NoDefaultLeavesZeroValue(t *testing.T) {
	type Config struct {
		Count int `env:"COUNT_UNSET_XYZ"`
	}
	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Count != 0 {
		t.Fatalf("expected zero value, got %d", c.Count)
	}
}

func TestProcess_EmptyEnvStringValueIsSet(t *testing.T) {
	// Regression test: an env var explicitly set to "" for a string
	// field must be applied, not silently dropped.
	type Config struct {
		Name string `env:"NAME_EMPTY_XYZ,,fallback"`
	}
	t.Setenv("NAME_EMPTY_XYZ", "")

	var c Config
	if err := Process(&c, ""); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Name != "" {
		t.Fatalf("expected explicit empty string to override default, got %q", c.Name)
	}
}

func TestProcess_EmptyEnvNumericValueErrors(t *testing.T) {
	// Regression test: an env var explicitly set to "" for a non-string
	// type must produce a parse error, not panic on Set(nil).
	type Config struct {
		Count int `env:"COUNT_EMPTY_XYZ"`
	}
	t.Setenv("COUNT_EMPTY_XYZ", "")

	var c Config
	if err := Process(&c, ""); err == nil {
		t.Fatal("expected error for empty numeric env var, got nil")
	}
}

func TestProcess_InvalidTagFormat(t *testing.T) {
	type Config struct {
		Name string `env:",required"`
	}
	var c Config
	if err := Process(&c, ""); !errors.Is(err, ErrInvalidTagFormat) {
		t.Fatalf("expected ErrInvalidTagFormat, got %v", err)
	}
}

func TestProcess_InvalidRequiredFlag(t *testing.T) {
	type Config struct {
		Name string `env:"NAME,yes"`
	}
	var c Config
	if err := Process(&c, ""); !errors.Is(err, ErrInvalidTagFormat) {
		t.Fatalf("expected ErrInvalidTagFormat, got %v", err)
	}
}

func TestProcess_InvalidInputType(t *testing.T) {
	var notAPointer int
	if err := Process(notAPointer, ""); !errors.Is(err, ErrInvalidInputType) {
		t.Fatalf("expected ErrInvalidInputType, got %v", err)
	}
}

func TestProcess_Prefix(t *testing.T) {
	type Config struct {
		Name string `env:"NAME"`
	}
	t.Setenv("APP_NAME", "prefixed")

	var c Config
	if err := Process(&c, "APP"); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if c.Name != "prefixed" {
		t.Fatalf("expected prefixed lookup, got %q", c.Name)
	}
}
