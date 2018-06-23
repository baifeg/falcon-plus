package funcs

import (
	"encoding/json"
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"log"
	"net/http"
)

var (
	mainMetrics    = []string{"mem.free", "threads"}
	healthMetrics  = []string{"status"}
	METRICS_PREFIX = "spring.actuator.metrics."
	HEALTH_PREFIX  = "spring.actuator.health."
)

func SpringMetrics() (L []*model.MetricValue) {
	url := fmt.Sprintf("http://127.0.0.1:%d/actuator/metrics", g.GetEnv().ServicePort)
	return actuatorInfo(url, METRICS_PREFIX, mainMetrics)
}

func SpringHealthMetrics() (L []*model.MetricValue) {
	url := fmt.Sprintf("http://127.0.0.1:%d/actuator/health", g.GetEnv().ServicePort)
	return actuatorInfo(url, HEALTH_PREFIX, healthMetrics)
}

func actuatorInfo(url string, metricPrefix string, metrics []string) (L []*model.MetricValue) {
	env := g.GetEnv()
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("get actuator metrics failed.", err)
		return
	}

	defer res.Body.Close()
	var m map[string]interface{}
	json.NewDecoder(res.Body).Decode(&m)

	tags := fmt.Sprintf("group=%s,tenant=%s,app=%s,service=%s", env.Group, env.Tenant, env.App, env.Service)
	for _, metric := range metrics {
		if value, ok := m[metric]; ok {
			name := fmt.Sprintf("%s%s", metricPrefix, metric)
			L = append(L, GaugeValue(name, value, tags))
		}
	}
	return
}
