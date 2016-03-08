package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var (
	Params struct {
		LagThreshold int
		PingInterval int
		URLsFile     string
		URL			 string
	}
	tomlFile = "healthmon.toml"
)

type (
	TargetData struct{
		Records_total int		`json:"records_total"`
		Records_filtered int	`json:"records_filtered"`
		Data []Target			`json:"data"`
	}

	Target struct{
		Alternative_ips 	string	`json:"alternative_ips"`
		Host_ip				string	`json:"host_ip"`
		Host_name			string	`json:"host_name"`
		Host_path			string	`json:"host_path"`
		Id					int		`json:"id"`
		Ip_url				string	`json:"ip_url"`
		Max_owtf_rank		int		`json:"max_owtf_rank"`
		Max_user_rank		int		`json:"max_user_rank"`
		Port_number			string	`json:"port_number"`
		Scope				bool	`json:"scope"`
		Target_url			string	`json:"target_url"`
		Top_domain			string	`json:"top_domain"`
		Top_url				string	`json:"top_url"`
		Url_scheme			string	`json:"url_scheme"`
	}
)


func init() {
	if _, err := toml.DecodeFile(tomlFile, &Params); err != nil {
		log.Fatal(err)
	}
}
