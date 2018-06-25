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
	METRICS_PREFIX = "spring.metrics."
	HEALTH_PREFIX  = "spring.health."
)

func SpringMetrics() (L []*model.MetricValue) {
	url := fmt.Sprintf("http://127.0.0.1:%d/actuator/metrics", g.GetEnv().ServicePort)
	return actuatorInfo(url, METRICS_PREFIX, mainMetrics)
}

func SpringHealthMetrics() (L []*model.MetricValue) {
	url := fmt.Sprintf("http://127.0.0.1:%d/actuator/health", g.GetEnv().ServicePort)
	items := actuatorInfo(url, HEALTH_PREFIX, healthMetrics)
	for _, item := range items {
		if item.Metric == "spring.health.status" {
			if item.Value == "UP" {
				item.Value = 1
			} else {
				item.Value = 0
			}
		}
	}
	return items
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
	result := translate(m)

	tags := fmt.Sprintf("group=%s,tenant=%s,app=%s,service=%s", env.Group, env.Tenant, env.App, env.Service)
	for _, metric := range metrics {
		if value, ok := result[metric]; ok {
			name := fmt.Sprintf("%s%s", metricPrefix, metric)
			L = append(L, GaugeValue(name, value, tags))
		}
	}
	return
}

// 将嵌套的json转成key-value
func translate(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	pop("", m, result)
	return result
}

func pop(key string, value interface{}, m map[string]string) {
	switch value.(type) {
	case map[string]interface{}:
		r, _ := value.(map[string]interface{})
		for k, v := range r {
			var prefix string
			if len(key) == 0 {
				prefix = k
			} else {
				prefix = key + "." + k
			}
			pop(prefix, v, m)
		}
	default:
		m[key] = fmt.Sprint(value)
	}
}
