package herd

type Manager struct {
	spawnQueue *[][]any
}

func NewManager() Manager {
	spawnQueue := make([][]any, 0, 32)
	return Manager{
		spawnQueue: &spawnQueue,
	}
}

func (c *Manager) Spawn(comps ...any) {
	*c.spawnQueue = append(*c.spawnQueue, comps)
}

func (c *Manager) clear() {
	if len(*c.spawnQueue) != 0 {
		spawnQueue := make([][]any, 0, 32)
		*c.spawnQueue = spawnQueue
	}
}
