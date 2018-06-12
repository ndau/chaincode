package vm

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSinglePushPop(t *testing.T) {
	st := newStack()
	assert.Equal(t, st.Depth(), 0)
	err := st.Push(NewNumber(123))
	assert.Nil(t, err)
	assert.Equal(t, st.Depth(), 1)
	v, err := st.Pop()
	assert.Nil(t, err)
	assert.Equal(t, v.String(), "123")
	assert.Equal(t, st.Depth(), 0)
}

func pushMulti(t *testing.T, st *Stack, values ...int64) {
	for _, v := range values {
		err := st.Push(NewNumber(v))
		assert.Nil(t, err)
	}
}

func checkMulti(t *testing.T, st *Stack, values ...int64) {
	for _, v := range values {
		n, err := st.Pop()
		assert.Nil(t, err)
		assert.Equal(t, n.String(), strconv.FormatInt(v, 10))
	}
}

func TestOverflow(t *testing.T) {
	st := newStack()
	for i := 0; i < maxStackDepth+2; i++ {
		err := st.Push(NewNumber(int64(i)))
		if err != nil {
			assert.Equal(t, i, maxStackDepth+1)
			return
		}
	}
	assert.Fail(t, "overflow never occurred")
}

func TestMultiPushPop(t *testing.T) {
	st := newStack()
	pushMulti(t, st, 2, 6, -7, 99)
	assert.Equal(t, st.Depth(), 4)
	checkMulti(t, st, 99, -7, 6, 2)
	_, err := st.Pop()
	assert.NotNil(t, err)
}

func TestPopAt(t *testing.T) {
	st := newStack()
	pushMulti(t, st, 1, 2, 3, 4, 5)
	n, err := st.PopAt(3)
	assert.Equal(t, n.String(), "2")
	assert.Nil(t, err)
	n, err = st.PopAt(0)
	assert.Equal(t, n.String(), "5")
	assert.Nil(t, err)
	n, err = st.PopAt(3)
	assert.NotNil(t, err)
	assert.Equal(t, st.Depth(), 3)
	checkMulti(t, st, 4, 3, 1)
}

func TestGet(t *testing.T) {
	st := newStack()
	pushMulti(t, st, 1, 2, 3, 4, 5)
	n, err := st.Get(3)
	assert.Nil(t, err)
	assert.Equal(t, n.String(), "2")
	assert.Equal(t, st.Depth(), 5)
	checkMulti(t, st, 5, 4, 3, 2, 1)
}

func TestString(t *testing.T) {
	st := newStack()
	assert.Equal(t, "|== Empty", st.String())
	pushMulti(t, st, 1, 2, 3, 4, 5)
	assert.Equal(t, "|== 5\n|== 4\n|== 3\n|== 2\n|== 1", st.String())
}

func listOfStructs() List {
	l := NewList()
	for i := int64(0); i < 5; i++ {
		s := NewStruct().Append(NewNumber(i)).Append(NewNumber(5 - i)).Append(NewBytes([]byte("hi")))
		l = l.Append(s)
	}
	return l
}

func TestPopAsList1(t *testing.T) {
	st := newStack()
	st.Push(NewList())
	l1, err := st.PopAsList()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), l1.Len())
}

func TestPopAsList2(t *testing.T) {
	st := newStack()
	st.Push(listOfStructs())
	l1, err := st.PopAsList()
	assert.Nil(t, err)
	assert.Equal(t, int64(5), l1.Len())
}

func TestPopAsListFail(t *testing.T) {
	st := newStack()
	st.Push(NewNumber(0))
	_, err := st.PopAsList()
	assert.NotNil(t, err)
	_, err = st.PopAsList()
	assert.NotNil(t, err)
}

func TestPopAsListOfStructs(t *testing.T) {
	st := newStack()
	ls := listOfStructs()
	st.Push(ls)
	l1, err := st.PopAsListOfStructs(0)
	assert.Nil(t, err)
	assert.Equal(t, int64(5), l1.Len())
}

func TestPopAsListOfStructsFail(t *testing.T) {
	st := newStack()
	st.Push(NewNumber(0))
	_, err := st.PopAsListOfStructs(2)
	st.Push(listOfStructs())
	_, err = st.PopAsListOfStructs(2)

	l := NewList()
	l = l.Append(NewNumber(0))
	st.Push(l)
	_, err = st.PopAsListOfStructs(0)
	assert.NotNil(t, err)

	l2 := NewList()
	l2 = l2.Append(NewStruct().Append(NewNumber(1))).Append(NewNumber(0))
	st.Push(l2)
	_, err = st.PopAsListOfStructs(0)
	assert.NotNil(t, err)

	l3 := NewList()
	l3 = l3.Append(NewStruct().Append(NewNumber(1))).Append(NewStruct().Append(NewList()))
	st.Push(l3)
	_, err = st.PopAsListOfStructs(0)
	assert.NotNil(t, err)

	l4 := NewList()
	l4 = l4.Append(NewStruct().Append(NewNumber(1))).Append(NewStruct())
	st.Push(l4)
	_, err = st.PopAsListOfStructs(0)
	assert.NotNil(t, err)

	_, err = st.PopAsListOfStructs(0)
	assert.NotNil(t, err)
}
