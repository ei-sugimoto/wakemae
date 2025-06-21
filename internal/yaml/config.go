package yaml

type Config struct {
	DNS DNS `yaml:"dns"`
}

type DNS struct {
	BindAddress string `yaml:"bind_address"`
	Upstream    string `yaml:"upstream"`
	Timeout     string `yaml:"timeout"`
}

var DefaultConfig = Config{
	DNS: DNS{
		BindAddress: "127.0.0.1:53",
		Upstream:    "1.1.1.1:53",
		Timeout:     "5s",
	},
}
