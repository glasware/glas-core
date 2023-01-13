package config

type Action struct {
	Name     string `hcl:"name,label"`
	Pattern  string `hcl:"pattern"`
	Template string `hcl:"template"`
	Regex    bool   `hcl:"regex,optional"`
}
