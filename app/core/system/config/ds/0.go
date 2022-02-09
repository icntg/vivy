package ds

type DataSource struct {
	MySQL   MySQL   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	MongoDB MongoDB `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
}
