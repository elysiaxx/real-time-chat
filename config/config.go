package config

type Config struct {
	PostgresqlConnectionString string
	MongoConnectionString      string
	SocketServerHost           string
	SocketServerPort           string
	WebServerHost              string
	WebServerPort              string
}

func DefaulConfig() *Config {
	return &Config{
		PostgresqlConnectionString: "host=localhost user=postgres password=example dbname=real-time-chat port=5432 sslmode=disable",
		MongoConnectionString:      "mongodb://root:example@localhost:27017/chat-app",
		SocketServerHost:           "0.0.0.0",
		SocketServerPort:           "9909",
		WebServerHost:              "0.0.0.0",
		WebServerPort:              "8808",
	}
}
