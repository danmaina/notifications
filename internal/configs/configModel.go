package configs

type Config struct {
	ApplicationConfigs GeneralConfigs
	Email              EmailConfigs
	RabbitMQ           RabbitMQConfigs
}

type GeneralConfigs struct {
	Port     string
	LogLevel int
}

type EmailConfigs struct {
	Username string
	Password string
	Host     string
	Port     string
}

type RabbitMQConfigs struct {
	URL string
}
