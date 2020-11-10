package leftist

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/theodesp/go-heaps"
)

func TestLeftistHeapInteger(t *testing.T) {
	heap := New().Init()

	numbers := []int{4, 3, 2, 5}

	for _, number := range numbers {
		heap.Insert(Int(number))
	}

	sort.Ints(numbers)

	for _, number := range numbers {
		if Int(number) != heap.DeleteMin().(go_heaps.Integer) {
			t.Fail()
		}
	}
}

func BenchmarkLeftistHeapInsert(b *testing.B) {
	for _, size := range []int{1, 10, 100, 1000, 10000, 100000, 1000000} {
		benchmarkLeftistHeapInsertRand(b, size)
	}
}

func benchmarkLeftistHeapInsertRand(b *testing.B, size int) {
	var numbers_rand = make([]int, size)
	rand.Seed(int64(size % 64))
	for i := 0; i < size; i++ {
		numbers_rand[i] = rand.Intn(size)
	}

	var numbers_rep = make([]int, size)
	for i := 0; i < size; i++ {
		numbers_rep[i] = 1
	}

	var numbers_incr = make([]int, size)
	for i := 0; i < size; i++ {
		numbers_incr[i] = i
	}

	var numbers_decr = make([]int, size)
	for i := size; i < 0; i-- {
		numbers_decr[i] = i
	}
	heap1 := New().Init()
	heap2 := New().Init()
	heap3 := New().Init()
	heap4 := New().Init()

	b.ResetTimer()

	b.Run(fmt.Sprintf("Insert_Rand_%d_Item", size), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				heap1.Insert(Int(numbers_rand[j]))
			}
		}
	})

	b.Run(fmt.Sprintf("Insert_Rep_%d_Item", size), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				heap2.Insert(Int(numbers_rep[j]))
			}
		}
	})

	b.Run(fmt.Sprintf("Insert_Incr_%d_Item", size), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				heap3.Insert(Int(numbers_incr[j]))
			}
		}
	})

	b.Run(fmt.Sprintf("Insert_decr_%d_Item", size), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < size; j++ {
				heap4.Insert(Int(numbers_decr[j]))
			}
		}
	})
}

func TestLeftistHeapString(t *testing.T) {
	heap := &LeftistHeap{}

	strs := []string{"a", "ccc", "bb", "d"}

	for _, str := range strs {
		heap.Insert(Str(str))
	}

	sort.Strings(strs)

	for _, str := range strs {
		if Str(str) != heap.DeleteMin().(go_heaps.String) {
			t.Fail()
		}
	}
}

func TestLeftistHeap(t *testing.T) {
	heap := &LeftistHeap{}

	numbers := []int{4, 3, -1, 5, 9}

	for _, number := range numbers {
		heap.Insert(Int(number))
	}

	if heap.FindMin() != Int(-1) {
		t.Fail()
	}

	heap.Clear()
	if heap.FindMin() != nil {
		t.Fail()
	}
}

func Int(value int) go_heaps.Integer {
	return go_heaps.Integer(value)
}

func Str(value string) go_heaps.String {
	return go_heaps.String(value)
}
