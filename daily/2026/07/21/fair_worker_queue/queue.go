package fairworkerqueue

import "fmt"

type Item struct {
	Tenant string
	Value  string
}

type Queue struct {
	queues map[string][]string
	order  []string
	next   int
	size   int
}

func New() *Queue {
	return &Queue{
		queues: make(map[string][]string),
	}
}

func (q *Queue) Enqueue(tenant string, value string) error {
	if tenant == "" {
		return fmt.Errorf("tenant must not be empty")
	}

	if _, ok := q.queues[tenant]; !ok {
		q.order = append(q.order, tenant)
	}
	q.queues[tenant] = append(q.queues[tenant], value)
	q.size++
	return nil
}

func (q *Queue) Next() (Item, bool) {
	if q.size == 0 || len(q.order) == 0 {
		return Item{}, false
	}

	start := q.next
	for scanned := 0; scanned < len(q.order); scanned++ {
		idx := (start + scanned) % len(q.order)
		tenant := q.order[idx]
		items := q.queues[tenant]
		if len(items) == 0 {
			continue
		}

		item := Item{
			Tenant: tenant,
			Value:  items[0],
		}
		q.queues[tenant] = items[1:]
		q.size--
		q.next = (idx + 1) % len(q.order)
		return item, true
	}

	return Item{}, false
}

func (q *Queue) Len() int {
	return q.size
}
