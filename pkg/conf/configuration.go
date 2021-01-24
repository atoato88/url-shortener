package conf

var Conf = Init()

type Configuration struct {
	Port  int
	Debug bool
}

func Init() Configuration {
	c := Configuration{}
	c.Init()
	return c
}

func (c *Configuration) Init() {
	c.Port = 3000
	c.Debug = false
}
