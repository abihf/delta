package delta

// Configuration customize delta behaviour
type Configuration struct {
	EncodeResponse bool
}

var globalConfig = &Configuration{
	EncodeResponse: false,
}

// Configure delta
func Configure(conf *Configuration) {
	globalConfig = conf
}
