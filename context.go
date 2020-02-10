package pipeline

import (
	"time"
)

// Tag for a given context key to be used as an identifier.
type Tag string

// Context contract for manipulating data across units of work. A context can be used to write, read or delete data
// A context should always support concurrent operations, since units of work aren't necessarily sequential
type Context interface {
	// Set is used to store a new key/value pair exclusively for this context.
	// It also lazy initializes the map if it was not used previously.
	Set(key Tag, value interface{})

	// Get returns the value for the given key, ie: (value, true).
	// If the value does not exists it returns (nil, false) or the default ("", false)
	Get(key Tag) (value interface{}, exists bool)

	// Delete the key and any value assigned to it
	Delete(key Tag)

	// GetString returns the value associated with the key as a string if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetString(key Tag) (s string, exists bool)

	// GetBool returns the value associated with the key as a boolean if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetBool(key Tag) (b bool, exists bool)

	// GetInt returns the value associated with the key as an int if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetInt(key Tag) (i int, exists bool)

	// GetUInt returns the value associated with the key as an uint if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetUInt(key Tag) (i uint, exists bool)

	// GetUInt64 returns the value associated with the key as an uint64 if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetUInt64(key Tag) (i uint64, exists bool)

	// GetInt64 returns the value associated with the key as an int64 if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetInt64(key Tag) (i64 int64, exists bool)

	// GetFloat64 returns the value associated with the key as a float64 if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetFloat64(key Tag) (f64 float64, exists bool)

	// GetTime returns the value associated with the key as time if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetTime(key Tag) (t time.Time, exists bool)

	// GetDuration returns the value associated with the key as a duration if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetDuration(key Tag) (d time.Duration, exists bool)

	// GetByteSlice returns the value associated with the key as a slice of bytes if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetByteSlice(key Tag) (ss []byte, exists bool)

	// GetStringSlice returns the value associated with the key as a slice of strings if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetStringSlice(key Tag) (ss []string, exists bool)

	// GetStringMap returns the value associated with the key as a map of interfaces if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetStringMap(key Tag) (sm map[string]interface{}, exists bool)

	// GetStringMapString returns the value associated with the key as a map of strings if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetStringMapString(key Tag) (sms map[string]string, exists bool)

	// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings if possible, and if it exists regardless of the type.
	// Note that if the type stored is different from the expected, the value will be nil / default and won't exist
	GetStringMapStringSlice(key Tag) (smss map[string][]string, exists bool)
}
