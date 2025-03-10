package bug1

import "sync"

type Counter struct {
	n    int64
	lock sync.Mutex 
}


func (c *Counter) Inc() {
	c.lock.Lock()     
	c.n++
	c.lock.Unlock()   
}

