package pipeline

import (
	"sync"
	"time"
)

// CreateContext for manipulating data across units of work.
// The created context supports concurrent operations
func CreateContext() Context {
	return &memMapContext{}
}

type memMapContext struct {
	sync.RWMutex
	value map[Tag]interface{}
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *memMapContext) Get(key Tag) (value interface{}, exists bool) {
	c.RLock()
	defer c.RUnlock()
	value, exists = c.value[key]
	return
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes the map if it was not used previously.
func (c *memMapContext) Set(key Tag, value interface{}) {
	c.Lock()
	if c.value == nil {
		c.value = make(map[Tag]interface{})
	}
	c.value[key] = value
	c.Unlock()
}

// Delete the key and any value assigned to it
func (c *memMapContext) Delete(key Tag) {
	c.Lock()
	delete(c.value, key)
	c.Unlock()
}

// GetString returns the value associated with the key as a string if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetString(key Tag) (s string, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetBool(key Tag) (b bool, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetInt(key Tag) (i int, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetUInt returns the value associated with the key as an uint if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetUInt(key Tag) (i uint, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		i, _ = val.(uint)
	}
	return
}

// GetUInt64 returns the value associated with the key as an uint if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetUInt64(key Tag) (i uint64, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		i, _ = val.(uint64)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetInt64(key Tag) (i64 int64, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64 if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetFloat64(key Tag) (f64 float64, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetTime(key Tag) (t time.Time, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetDuration(key Tag) (d time.Duration, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetByteSlice returns the value associated with the key as a slice of bytes if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetByteSlice(key Tag) (ss []byte, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		ss, _ = val.([]byte)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetStringSlice(key Tag) (ss []string, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetStringMap(key Tag) (sm map[string]interface{}, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetStringMapString(key Tag) (sms map[string]string, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings if possible, and if it exists regardless of the type.
// Note that if the type stored is different from the expected, the value will be nil but the exists will be true
func (c *memMapContext) GetStringMapStringSlice(key Tag) (smss map[string][]string, exists bool) {
	var val interface{}
	if val, exists = c.Get(key); exists && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
