package lrucache

const (
	cacheSize = 3000
)

type Node stuct {
	url string
	prev *Node
	next *Node
}

type Cache struct {
	mru   *Node
	lru   *Node
	length int
}

func (c *Cache) addToCache(n *Node) *Node {
	if c.length == 0 {
		c.mru = n
		c.lru = n
		c.length++
		return nil
	}
	prevMru := c.mru
	c.mru = n
	c.mru.next = prevMru
	prevMru.prev = n
	c.length++
	if c.length > cacheSize {
		deletedNode := c.tail
		newTail := c.lru.prev
		newTail.next = nil
		c.lru = newTail
		c.length--
		return deletedNode
	}
	return nil
}

func (c *Cache) moveToMru(n *Node) {
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