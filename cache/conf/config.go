package conf

type Config struct {
	DebugMode bool
}

var Conf Config = Config{
	DebugMode: true,
}
