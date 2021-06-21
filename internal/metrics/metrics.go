package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var processedByCRUDHandler *prometheus.CounterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "processed_by_crud_handler",
		Help: "Number of processed By CRUD Handler",
	},
	[]string{"handler"}, // labels
)

func RegisterMetrics() {
	prometheus.MustRegister(processedByCRUDHandler)
}

func incrementByHandler(handler string, count int) {
	processedByCRUDHandler.With(prometheus.Labels{"handler": handler}).Add(float64(count))
}

// count - количество созданных сущностей
func IncrementSuccessfulCreate(count int) {
	incrementByHandler("create", count)
}

// count - количество прочитанных сущностей
func IncrementSuccessfulRead(count int) {
	incrementByHandler("read", count)
}

// count - количество обновлённых сущностей
func IncrementSuccessfulUpdate(count int) {
	incrementByHandler("update", count)
}

// count - количество удалённых сущностей
func IncrementSuccessfulDelete(count int) {
	incrementByHandler("delete", count)
}
