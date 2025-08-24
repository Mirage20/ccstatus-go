package core

// Config holds application configuration.
type Config struct {
	data map[string]interface{}
}

// NewConfig creates a new configuration.
func NewConfig() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a configuration value.
func (c *Config) Get(key string, defaultValue interface{}) interface{} {
	if value, exists := c.data[key]; exists {
		return value
	}
	return defaultValue
}

// GetBool retrieves a boolean configuration value.
func (c *Config) GetBool(key string, defaultValue bool) bool {
	value := c.Get(key, defaultValue)
	if boolValue, ok := value.(bool); ok {
		return boolValue
	}
	return defaultValue
}

// GetInt retrieves an integer configuration value.
func (c *Config) GetInt(key string, defaultValue int) int {
	value := c.Get(key, defaultValue)
	if intValue, ok := value.(int); ok {
		return intValue
	}
	return defaultValue
}

// GetInt64 retrieves an int64 configuration value.
func (c *Config) GetInt64(key string, defaultValue int64) int64 {
	value := c.Get(key, defaultValue)
	if intValue, ok := value.(int64); ok {
		return intValue
	}
	return defaultValue
}

// GetString retrieves a string configuration value.
func (c *Config) GetString(key string, defaultValue string) string {
	value := c.Get(key, defaultValue)
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return defaultValue
}

// Set stores a configuration value.
func (c *Config) Set(key string, value interface{}) {
	c.data[key] = value
}
