package setting

type UserInfo struct {
	UserName     string `mapstructure:"username" json:"username" yaml:"username"`
	PassWd       string `mapstructure:"passwd" json:"passwd" yaml:"passwd"`
	TimeInterval int64  `mapstructure:"time_interval" json:"time_interval" yaml:"time_interval"`
}
type Config struct {
	Account []UserInfo `mapstructure:"account" json:"account" yaml:"account"`
}

var Cnf Config
