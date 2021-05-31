package metrics

type Publisher interface {
	PublishFlushing(count int)
}
