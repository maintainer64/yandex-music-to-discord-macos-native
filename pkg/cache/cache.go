package cache

type IFixedQueueItem interface {
	GetCacheKey() string
}

type FixedQueue struct {
	queue []IFixedQueueItem
	size  int
}

func NewFixedQueue(size int) *FixedQueue {
	return &FixedQueue{
		size:  size,
		queue: make([]IFixedQueueItem, 0, size),
	}
}

func (fq *FixedQueue) Get(key string) IFixedQueueItem {
	for _, item := range fq.queue {
		if item.GetCacheKey() == key {
			return item
		}
	}

	return nil
}

func (fq *FixedQueue) Push(item IFixedQueueItem) {
	fq.queue = append(fq.queue, item)
	if len(fq.queue) > fq.size {
		fq.queue = fq.queue[1:]
	}
}
