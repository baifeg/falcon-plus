package g

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

const SERVICE_TYPE string = "SPRING_BOOT_WEB"
const DEFAULT_STEP int = 30

type SystemEnv struct {
	Template string         `json:"template"`
	Uri      []string       `json:"uri"`
	Step     int            `json:"step"`
	Prefix   []string       `json:"prefix"`
	Mapping  map[string]int `json:"mapping"`
	Counter  []string       `json:"counter"`
	Ignore   []string       `json:"ignore"`

	Group   string `json:"group"`
	Ip      string `json:"ip"`
	App     string `json:"app"`
	Tenant  string `json:"tenant"`
	Service string `json:"service"`

	HbsAddr       string   `json:"hbsAddr"`       // host:port
	TransferAddrs []string `json:"transferAddrs"` // host:port
}

var (
	env *SystemEnv
)

func GetEnv() *SystemEnv {
	return env
}

func InitEnv() {
	template := os.Getenv("COLLECT_TEMPLATE")
	uri := toArray(os.Getenv("COLLECT_URI"))
	step, err := strconv.Atoi(os.Getenv("COLLECT_STEP"))
	if err != nil {
		step = DEFAULT_STEP
	}
	prefix := toArray(os.Getenv("COLLECT_PREFIX"))
	mapping := toMap(so.Getenv("COLLECT_MAPPING"))
	counter := toArray(os.Getenv("COLLECT_COUNTER"))
	ignore := toArray(os.Getenv("COLLECT_IGNORE"))

	service := os.Getenv("SERVICE_CODE")
	group := os.Getenv("SERVICE_DEPLOY_CODE")
	app := os.Getenv("APP_CODE")
	tenant := os.Getenv("TENANT_CODE")
	ip, ipOk := os.LookupEnv("KETTY_IP")

	hbsAddr := os.Getenv("HBS_ADDR")
	transferAddrs := toArray(os.Getenv("TRANSFER_ADDR"))
	if !ipOk {
		ip = ""
	}

	e := SystemEnv{
		Template:      template,
		Uri:           uri,
		Step:          step,
		Prefix:        prefix,
		Mapping:       mapping,
		Counter:       counter,
		Ignore:        ignore,
		Group:         group,
		Ip:            ip,
		App:           app,
		Tenant:        tenant,
		Service:       service,
		HbsAddr:       hbsAddr,
		TransferAddrs: transferAddrs,
	}

	env = &e
	j, err := json.Marshal(e)
	if err == nil {
		log.Println("read system env done: ", string(j))
	} else {
		log.Fatalln("read system env failed.", err)
	}
}

func toArray(str string) []string {
	str = strings.Trim(str, ",")
	return strings.Split(str, ",")
}

func toMap(str string) map[string]int {
	str = strings.Trim(str, ",")
	mapping := make(map[string]int)
	for _, item := range strings.Split(str, ",") {
		if strings.Contains(item, "=") {
			kv := strings.Split(str, "=")
			k := kv[0]
			v, err := strconv.Atoi(kv[1])
			if err == nil {
				mapping[k] = v
			}
		}
	}
	return mapping
}
