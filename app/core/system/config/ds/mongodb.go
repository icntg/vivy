package ds

type MongoDB struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     uint16 `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
