package lrucache

// Node - a single node of the cache
type Node struct {
	url  string
	prev *Node
	next *Node
}

// Cache - The cache ordered from most recently used (mru) to
// least recently used (lru)
type Cache struct {
	mru    *Node
	lru    *Node
	length int
}

// cacheSize - The maximum size of the cache.
const cacheSize = 3000

// AddToCache - Create a node with the new url and places it in
// the most recently used (mru) position. If the new node causes
// the cache to exceed its max size, the least recently used (lru)
// node is removed and the url value associated with that node is
// returned for deletion from the cache node map.
func (c *Cache) AddToCache(url string) (*Node, *string) {
	n := &Node{url: url}
	if c.length == 0 {
		c.mru = n
		c.lru = n
		c.length++
		return n, nil
	}
	prevMru := c.mru
	c.mru = n
	c.mru.next = prevMru
	prevMru.prev = n
	c.length++
	if c.length > cacheSize {
		deletedNode := c.lru
		newTail := c.lru.prev
		newTail.next = nil
		c.lru = newTail
		c.length--
		return n, &deletedNode.url
	}
	return n, nil
}

// MoveToMru - Moves a node to the most recently used (mru) position
func (c *Cache) MoveToMru(n *Node) {
	if c.length == 1 {
		return
	}
	if c.mru == n {
		return
	}
	left := n.prev
	right := n.next
	left.next = right
	right.prev = left
	prevMru := c.mru
	prevMru.prev = n
	c.mru = n
	n.next = prevMru
	n.prev = nil
	return
}
