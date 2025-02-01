package configs

// Config Is a wrapper that wraps both GeneralConfigs and EmailConfigs
type Config struct {
	ApplicationConfigs GeneralConfigs `yaml:"generalConfigs"`
	Email              EmailConfigs   `yaml:"email"`
}

// GeneralConfigs are a set of application specific configurations
type GeneralConfigs struct {
	Port     string `yaml:"port"`
	LogLevel int    `yaml:"logLevel"`
}

// EmailConfigs is a set of configurable parameters from your email service provider
type EmailConfigs struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}
