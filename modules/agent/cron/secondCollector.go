package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func CollectOthers() {

	if !g.Config().Transfer.Enabled {
		return
	}

	if len(g.Config().Transfer.Addrs) == 0 {
		return
	}

	if g.GetEnv().ServiceType == "SPRING_BOOT_WEB" {
		fns := []func(){funcs.SpringMetrics, funcs.SpringHealthMetrics}
		go collect(30, fns)
	}
}

func collectOther(sec int64, fns []func() []*model.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec))
	defer t.Stop()
	for {
		<-t.C

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		mvs := []*model.MetricValue{}

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			for _, mv := range items {
				mvs = append(mvs, mv)
			}
		}

		now := time.Now().Unix()
		for j := 0; j < len(mvs); j++ {
			mvs[j].Step = sec
			mvs[j].Endpoint = hostname
			mvs[j].Timestamp = now
		}

		g.SendToTransfer(mvs)

	}
}
