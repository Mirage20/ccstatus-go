package config

import (
	"testing"

	"github.com/knadh/koanf/v2"
)

func TestGet(t *testing.T) {
	type TestConfig struct {
		Name    string `yaml:"name"`
		Value   int    `yaml:"value"`
		Enabled bool   `yaml:"enabled"`
	}

	tests := []struct {
		name         string
		setupKoanf   func() *koanf.Koanf
		path         string
		defaultValue interface{}
		expected     interface{}
	}{
		{
			name: "loads struct config successfully",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("test.name", "loaded")
				_ = k.Set("test.value", 42)
				_ = k.Set("test.enabled", true)
				return k
			},
			path: "test",
			defaultValue: TestConfig{
				Name:    "default",
				Value:   0,
				Enabled: false,
			},
			expected: TestConfig{
				Name:    "loaded",
				Value:   42,
				Enabled: true,
			},
		},
		{
			name: "returns default when path doesn't exist",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("other.name", "something")
				return k
			},
			path: "nonexistent",
			defaultValue: TestConfig{
				Name:    "default",
				Value:   10,
				Enabled: true,
			},
			expected: TestConfig{
				Name:    "default",
				Value:   10,
				Enabled: true,
			},
		},
		{
			name:       "handles nil koanf",
			setupKoanf: func() *koanf.Koanf { return nil },
			path:       "test",
			defaultValue: TestConfig{
				Name: "default",
			},
			expected: TestConfig{
				Name: "default",
			},
		},
		{
			name: "handles empty path",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("test.name", "loaded")
				return k
			},
			path: "",
			defaultValue: TestConfig{
				Name: "default",
			},
			expected: TestConfig{
				Name: "default",
			},
		},
		{
			name: "loads primitive string",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("stringval", "hello world")
				return k
			},
			path:         "stringval",
			defaultValue: "default string",
			expected:     "hello world",
		},
		{
			name: "loads primitive int",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("intval", 123)
				return k
			},
			path:         "intval",
			defaultValue: 0,
			expected:     123,
		},
		{
			name: "loads primitive bool",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("boolval", true)
				return k
			},
			path:         "boolval",
			defaultValue: false,
			expected:     true,
		},
		{
			name: "partially overrides struct fields",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("test.name", "partial")
				// Value and Enabled not set - should keep defaults
				return k
			},
			path: "test",
			defaultValue: TestConfig{
				Name:    "default",
				Value:   99,
				Enabled: true,
			},
			expected: TestConfig{
				Name:    "partial",
				Value:   99,
				Enabled: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{k: tt.setupKoanf()}

			switch v := tt.defaultValue.(type) {
			case TestConfig:
				result := Get(reader, tt.path, v)
				if result != tt.expected {
					t.Errorf("Get() = %+v, want %+v", result, tt.expected)
				}
			case string:
				result := Get(reader, tt.path, v)
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			case int:
				result := Get(reader, tt.path, v)
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			case bool:
				result := Get(reader, tt.path, v)
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestGetComponent(t *testing.T) {
	type ComponentConfig struct {
		Template string `yaml:"template"`
		Icon     string `yaml:"icon"`
		Color    string `yaml:"color"`
		ShowZero bool   `yaml:"show_zero"`
	}

	tests := []struct {
		name          string
		setupKoanf    func() *koanf.Koanf
		componentName string
		defaultValue  ComponentConfig
		expected      ComponentConfig
	}{
		{
			name: "loads component config with components prefix",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("components.mycomponent.template", "{{.Icon}} {{.Name}}")
				_ = k.Set("components.mycomponent.icon", "🚀")
				_ = k.Set("components.mycomponent.color", "cyan")
				_ = k.Set("components.mycomponent.show_zero", true)
				return k
			},
			componentName: "mycomponent",
			defaultValue: ComponentConfig{
				Template: "{{.Default}}",
				Icon:     "",
				Color:    "white",
				ShowZero: false,
			},
			expected: ComponentConfig{
				Template: "{{.Icon}} {{.Name}}",
				Icon:     "🚀",
				Color:    "cyan",
				ShowZero: true,
			},
		},
		{
			name: "returns default when component doesn't exist",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("components.other.template", "other")
				return k
			},
			componentName: "mycomponent",
			defaultValue: ComponentConfig{
				Template: "{{.Default}}",
				Icon:     "⚡",
				Color:    "yellow",
				ShowZero: false,
			},
			expected: ComponentConfig{
				Template: "{{.Default}}",
				Icon:     "⚡",
				Color:    "yellow",
				ShowZero: false,
			},
		},
		{
			name: "handles missing components section",
			setupKoanf: func() *koanf.Koanf {
				k := koanf.New(".")
				_ = k.Set("other", "value")
				return k
			},
			componentName: "mycomponent",
			defaultValue: ComponentConfig{
				Template: "{{.Default}}",
				Icon:     "",
				Color:    "white",
				ShowZero: false,
			},
			expected: ComponentConfig{
				Template: "{{.Default}}",
				Icon:     "",
				Color:    "white",
				ShowZero: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &Reader{k: tt.setupKoanf()}
			result := GetComponent(reader, tt.componentName, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("GetComponent() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// TestGetComponentMapOverrideSemantics documents (and locks) how a partial
// map[string]string override in YAML interacts with a defaultConfig() that
// ships a non-empty map, as used by the effort/model/account components'
// icons and colors maps.
//
// koanf's UnmarshalWithConf does not set mapstructure's ZeroFields, and
// mapstructure's decodeMapFromMap reuses the destination map instead of
// replacing it when it is non-nil - so overriding one key MERGES into the
// existing map and leaves sibling keys untouched, rather than replacing the
// whole map.
func TestGetComponentMapOverrideSemantics(t *testing.T) {
	type mapConfig struct {
		Colors map[string]string `yaml:"colors"`
	}

	k := koanf.New(".")
	_ = k.Set("components.mycomponent.colors.high", "red")

	reader := &Reader{k: k}
	got := GetComponent(reader, "mycomponent", &mapConfig{
		Colors: map[string]string{"high": "yellow", "low": "gray"},
	})

	if got.Colors["high"] != "red" {
		t.Errorf("override not applied: Colors[high] = %q, want %q", got.Colors["high"], "red")
	}
	if got.Colors["low"] != "gray" {
		t.Errorf(
			"partial map override replaced the whole map instead of merging: Colors[low] = %q, want %q (unchanged default)",
			got.Colors["low"],
			"gray",
		)
	}
}
