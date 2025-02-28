package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	keys     map[*ListItem]Key
	// Я не хочу хранить в ListItem.Value ключ и значение, поэтому сделал этот словарь чтобы можно было получить
	// ключ по *ListItem
}

func (c *lruCache) Set(key Key, value any) bool {
	// если элемент уже есть - перемещаем в начало
	listItem, exists := c.items[key]
	if exists {
		listItem.Value = value
		c.queue.MoveToFront(listItem)

		return true
	}

	// добавляем элемент в начало
	listItem = c.queue.PushFront(value)
	c.items[key] = listItem
	c.keys[listItem] = key

	// если размер очереди > вместимости кэша - удаляем последний элемент
	if c.capacity < c.queue.Len() {
		lastListItem := c.queue.Back()
		keyToRemove := c.keys[lastListItem]
		delete(c.items, keyToRemove)
		delete(c.keys, lastListItem)
		c.queue.Remove(lastListItem)
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	listItem, exists := c.items[key]
	if exists {
		c.queue.MoveToFront(listItem)
		return listItem.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.keys = make(map[*ListItem]Key, c.capacity)
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}
