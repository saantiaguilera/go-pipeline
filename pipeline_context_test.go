package pipeline

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TestKeyTag  Tag = "test-key"
	TestKey2Tag Tag = "test-key2"
)

func TestSet_GivenAContext_WhenSettingAValueTwice_ThenItGetsOverwritten(t *testing.T) {
	ctx := CreateContext()

	value, exists := ctx.Get(TestKeyTag)

	assert.Equal(t, nil, value)
	assert.False(t, exists)

	ctx.Set(TestKeyTag, 1234)

	value, exists = ctx.Get(TestKeyTag)

	assert.Equal(t, 1234, value)
	assert.True(t, exists)
}

func TestSet_GivenAContext_WhenSettingAValue_ThenItGetsOnlyOnThatIndex(t *testing.T) {
	ctx := CreateContext()
	ctx.Set(TestKeyTag, 1234)
	ctx.Set(TestKeyTag, 5678)

	value, exists := ctx.Get(TestKeyTag)

	assert.Equal(t, 5678, value)
	assert.True(t, exists)
}

func TestSet_GivenAContext_WhenDeletingAValue_ThenItGetsDeleted(t *testing.T) {
	ctx := CreateContext()
	ctx.Set(TestKeyTag, 1234)
	ctx.Delete(TestKeyTag)

	value, exists := ctx.Get(TestKeyTag)

	assert.Equal(t, nil, value)
	assert.False(t, exists)
}

func TestSet_GivenAContext_WhenDeletingAValueTwice_ThenNothingHappens(t *testing.T) {
	ctx := CreateContext()
	ctx.Set(TestKeyTag, 1234)
	ctx.Delete(TestKeyTag)
	ctx.Delete(TestKeyTag)

	value, exists := ctx.Get(TestKeyTag)

	assert.Equal(t, nil, value)
	assert.False(t, exists)
}

func TestPipelineContext_GetAnonymous_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 1},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: 1, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 1},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: nil, exists: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.Get(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetString_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: "123"},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: "123", exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: "123"},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: "", exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: "", exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: "", exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetString(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetBool_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: true, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: false, exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: false, exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: false, exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetBool(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetInt_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 12},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: 12, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: 0, exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: "asdf"},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: 0, exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: 0, exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetInt(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetInt64_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: int64(123)},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: int64(123), exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: int64(123)},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: int64(0), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: int64(0), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: int64(0), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetInt64(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetFloat64_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: float64(123)},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: float64(123), exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: float64(0), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: float64(0), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: float64(0), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetFloat64(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetTime_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: time.Unix(int64(123), int64(123))},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Unix(int64(123), int64(123)), exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: time.Unix(int64(123), int64(123))},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Time{}, exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Time{}, exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Time{}, exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetTime(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetDuration_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: time.Second},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Second, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: time.Second},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Duration(0), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Duration(0), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: time.Duration(0), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetDuration(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetByteSlice_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: []byte("abc")},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []byte("abc"), exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: []byte("abc")},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []byte(nil), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []byte(nil), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []byte(nil), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetByteSlice(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetStringSlice_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: []string{"a", "b"}},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []string{"a", "b"}, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: []string{"a"}},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []string(nil), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []string(nil), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: []string(nil), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetStringSlice(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetStringMap_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: map[string]interface{}{}},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]interface{}{}, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]interface{}(nil), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]interface{}(nil), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]interface{}(nil), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetStringMap(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetStringMapString_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: map[string]string{}},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]string{}, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]string(nil), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]string(nil), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string]string(nil), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetStringMapString(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}

func TestPipelineContext_GetStringMapStringSlice_Suite(t *testing.T) {
	tests := []struct {
		name  string
		store struct {
			key   Tag
			value interface{}
		}
		retrieve Tag
		want     struct {
			value  interface{}
			exists bool
		}
	}{
		{
			name: "given a stored value, when retrieving it, then it exists and it's the same",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: map[string][]string{}},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string][]string{}, exists: true},
		},
		{
			name: "given a stored value, when retrieving something else, then it doesnt exists and it's null",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: true},
			retrieve: TestKey2Tag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string][]string(nil), exists: false},
		},
		{
			name: "given a stored value of a different type, when retrieving it, then it exists but is default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: 123},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string][]string(nil), exists: true},
		},
		{
			name: "given a stored nil value, when retrieving it, then it exist and it's default value",
			store: struct {
				key   Tag
				value interface{}
			}{key: TestKeyTag, value: nil},
			retrieve: TestKeyTag,
			want: struct {
				value  interface{}
				exists bool
			}{value: map[string][]string(nil), exists: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := CreateContext()
			ctx.Set(tt.store.key, tt.store.value)
			if got, exists := ctx.GetStringMapStringSlice(tt.retrieve); !reflect.DeepEqual(got, tt.want.value) || exists != tt.want.exists {
				t.Errorf("expected - got = {%v %v}, want %v", got, exists, tt.want)
			}
		})
	}
}
