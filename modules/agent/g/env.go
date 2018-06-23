package g

import (
	"os"
	"log"
	"encoding/json"
)

const SERVICE_TYPE string = "springboot"

type SystemEnv struct {
	Group			string	`json:"group"`
	ServiceType		string	`json:"serviceType"`
	Ip				string	`json:"ip"`
	App				string	`json:"app"`
	Tenant			string	`json:"tenant"`
	Service			string	`json:"service"`
}

var (
	env		*SystemEnv
)

func GetEnv() *SystemEnv {
	return env
}

func InitEnv() {
	service := os.Getenv("SERVICE_CODE")
	group := os.Getenv("SERVICE_DEPLOY_CODE")
	app := os.Getenv("APP_CODE")
	tenant := os.Getenv("TENANT_CODE")
	ip := os.Getenv("KETTY_IP")
	
	e := SystemEnv{
		Group: group,
		ServiceType: SERVICE_TYPE,
		Ip: ip,
		App: app,
		Tenant: tenant,
		Service: service,
	}
	
	env = &e
	j, err := json.Marshal(e)
	if err != nil {
		log.Println("read system env done: ", string(j))
	} else {
		log.Println("read system env failed.", err)
	}
}
