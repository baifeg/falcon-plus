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
	metrics        = []string{"mem.free", "threads"}
	METRICS_PREFIX = "spring.actuator.metrics."
)

func SpringBootMetrics() (L []*model.MetricValue) {
	env := g.GetEnv()
	url := fmt.Sprintf("http://127.0.0.1:%d/actuator/metrics", env.ServicePort)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("get actuator metrics failed.", err)
		return
	}

	defer res.Body.Close()
	var m map[string]string
	json.NewDecoder(res.Body).Decode(&m)

	tags := fmt.Sprintf("group=%s,tenant=%s,app=%s,service=%s", env.Group, env.Tenant, env.App, env.Service)
	for _, metric := range metrics {
		if value, ok := m[metric]; ok {
			L = append(L, GaugeValue(METRICS_PREFIX+metric, value, tags))
		}
	}
	return
}
