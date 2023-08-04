package herd

type Manager struct {
	spawnQueue []any
}

func newManager() *Manager {
	spawnQueue := make([]any, 0, 32)
	return &Manager{
		spawnQueue: spawnQueue,
	}
}

func (c *Manager) Spawn(bundle any) {
	c.spawnQueue = append(c.spawnQueue, bundle)
}

func (c *Manager) clear() {
	if len(c.spawnQueue) != 0 {
		c.spawnQueue = make([]any, 0, 32)
	}
}
