package config

type Conf struct {
	DB map[string]DB `json:"db"`
}

type DB struct {
	Driver       string   `json:"driver"`
	Host         string   `json:"host"`
	Sources      []string `json:"sources"`
	Replicas     []string `json:"replicas"`
	Port         int      `json:"port"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	AuthDatabase string   `json:"auth_database"`
	Database     string   `json:"database"`
	Alias        string   `json:"alias"`
	Options      string   `json:"options"`
}
