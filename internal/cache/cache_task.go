package cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"todo-list/internal/task"
)

type taskCacheData struct {
	tasks      []task.Task
	lastUpdate time.Time
}

type TaskCache struct {
	cache        atomic.Value
	ttl          time.Duration
	isRefreshing bool
	mu           sync.Mutex
	loader       func() ([]task.Task, error)
}

func NewTaskCache(ttl time.Duration, loader func() ([]task.Task, error)) *TaskCache {
	return &TaskCache{
		ttl:    ttl,
		loader: loader,
	}
}

func (c *TaskCache) Init() error {
	tasks, err := c.loader()
	if err != nil {
		return err
	}
	c.Set(tasks)
	return nil
}

func (c *TaskCache) refreshOnce() {
	c.mu.Lock()
	if c.isRefreshing {
		c.mu.Unlock()
		return
	}
	c.isRefreshing = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.isRefreshing = false
		c.mu.Unlock()
	}()

	tasks, err := c.loader()
	if err == nil {
		c.Set(tasks)
	}
}

func (c *TaskCache) isExpired(data taskCacheData) bool {
	return time.Since(data.lastUpdate) > c.ttl
}

func (c *TaskCache) Get() ([]task.Task, error) {
	raw := c.cache.Load()
	if raw == nil {
		return nil, errors.New("cache is empty")
	}
	data, ok := raw.(taskCacheData)
	if !ok {
		return nil, errors.New("task cache data is not of type TaskCacheData")
	}
	if c.isExpired(data) {
		go c.refreshOnce()
	}

	return data.tasks, nil
}

func (c *TaskCache) Set(tasks []task.Task) {
	data := taskCacheData{
		tasks:      tasks,
		lastUpdate: time.Now(),
	}
	c.cache.Store(data)
}
