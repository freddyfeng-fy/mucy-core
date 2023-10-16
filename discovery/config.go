package discovery

type Etcd struct {
	Address string `mapstructure:"address" json:"address" yaml:"address"`
	Port    int    `mapstructure:"port" json:"port" yaml:"port"`
}
