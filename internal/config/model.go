package config

type Config struct {
	KAFKA_BROKER_URL      string `mapstructure:"KAFKA_BROKER_URL"`
	PostgresDB            string `mapstructure:"POSTGRES_DB"`
	PostgresUser          string `mapstructure:"POSTGRES_USER"`
	PostgresPass          string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDSN           string `mapstructure:"POSTGRES_DSN"`
	PostgresHost          string `mapstructure:"POSTGRES_HOST"`
	PostgresPort          string `mapstructure:"POSTGRES_PORT"`
	RedisAddr             string `mapstructure:"REDIS_ADDR"`
	RedisPassword         string `mapstructure:"REDIS_PASSWORD"`
	PORT                  string `mapstructure:"PORT"`
	ORIGINS               string `mapstructure:"ORIGINS"`
	JwtSecret             string `mapstructure:"JWT_SECRET"`
	RefreshTokenSecret    string `mapstructure:"REFRESH_TOKEN_SECRET"`
	JwtTokenRefreshExpiry int64  `mapstructure:"JWT_TOKEN_REFRESH_EXPIRY"`
	JwtTokenExpiry        int64  `mapstructure:"JWT_TOKEN_EXPIRY"`
	EnableDBQueryLogging  bool   `mapstructure:"ENABLE_DB_QUERY_LOGGING"`

	PstkSecretHash string `mapstructure:"PSTK_SECRET_HASH"`
	FlwSecretHash  string `mapstructure:"FLW_SECRET_HASH"`
}
