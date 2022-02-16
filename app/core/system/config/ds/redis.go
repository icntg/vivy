package ds

type Redis struct {
	MaxIdle  int    `mapstructure:"max-idle" json:"maxIdle" yaml:"max-idle"`  // redis最大的空闲连接数
	Protocol string `mapstructure:"protocol" json:"protocol" yaml:"protocol"` // 数通信协议tcp或者udp
	Address  string `mapstructure:"address" json:"address" yaml:"address"`    // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
}
