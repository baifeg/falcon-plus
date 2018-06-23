package g

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const SERVICE_TYPE string = "SPRING_BOOT_WEB"

type SystemEnv struct {
	Group       string `json:"group"`
	ServiceType string `json:"serviceType"`
	Ip          string `json:"ip"`
	App         string `json:"app"`
	Tenant      string `json:"tenant"`
	Service     string `json:"service"`
	ServicePort int    `json:"servicePort"`
}

var (
	env *SystemEnv
)

func GetEnv() *SystemEnv {
	return env
}

func InitEnv() {
	service := os.Getenv("SERVICE_CODE")
	group := os.Getenv("SERVICE_DEPLOY_CODE")
	portKey := strings.ToUpper(group)
	portKey = strings.Replace(portKey, "-", "_", -1) + "_SERVICE_PORT"
	servicePort, portOk := os.LookupEnv(portKey)
	if !portOk {
		servicePort = "0"
	}
	app := os.Getenv("APP_CODE")
	tenant := os.Getenv("TENANT_CODE")
	ip := os.Getenv("KETTY_IP")
	serviceType, ok := os.LookupEnv("SERVICE_TYPE")
	if !ok {
		serviceType = SERVICE_TYPE
	}

	e := SystemEnv{
		Group:       group,
		ServiceType: serviceType,
		Ip:          ip,
		App:         app,
		Tenant:      tenant,
		Service:     service,
		ServicePort: int(servicePort),
	}

	env = &e
	j, err := json.Marshal(e)
	if err == nil {
		log.Println("read system env done: ", string(j))
	} else {
		log.Fatalln("read system env failed.", err)
	}
}
