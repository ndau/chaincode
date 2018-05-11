package vm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSinglePushPop(t *testing.T) {
	st := newStack()
	assert.Equal(t, st.Depth(), 0)
	err := st.Push(newNumber(123))
	assert.Nil(t, err)
	assert.Equal(t, st.Depth(), 1)
	v, err := st.Pop()
	assert.Nil(t, err)
	assert.Equal(t, v.String(), "123")
	assert.Equal(t, st.Depth(), 0)
}

func TestMultiPushPop(t *testing.T) {
	st := newStack()
	values := []int64{2, 6, -7, 99}
	for _, v := range values {
		err := st.Push(newNumber(v))
		assert.Nil(t, err)
	}
	assert.Equal(t, st.Depth(), len(values))
	fmt.Println(st.String())
	for i := 0; i < len(values); i++ {
		v, err := st.Pop()
		assert.Nil(t, err)
		assert.Equal(t, v.String(), strconv.FormatInt(values[len(values)-i-1], 10))
	}
}
