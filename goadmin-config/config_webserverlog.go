package goadmin_config

type WebServerLog struct {
	Enable          bool     `yaml:"enable"`
	LogIDShowheader bool     `yaml:"log_id_show_header" mapstructure:"log_id_show_header"`
	LogPath         string   `yaml:"log_path" mapstructure:"log_path"`
	LogFormat       string   `yaml:"log_format" mapstructure:"log_format"`
	Output          string   `yaml:"output"`
	SkipPaths       []string `yaml:"skip_paths" mapstructure:"skip_paths"` // 忽略哪些路径
	HostName        string   `yaml:"host_name" mapstructure:"host_name"`
	LocalIP         string   `yaml:"local_ip"  mapstructure:"local_ip"`
}
