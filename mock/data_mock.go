package mock

type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}

type Iterator struct {
	data     []string
	position int
	closed   bool
}

func NewIterator(data []string) *Iterator {
	return &Iterator{
		data:     data,
		position: -1,
	}
}

func (it *Iterator) HasNext() bool {
	return !it.closed && it.position+1 < len(it.data)
}

func (it *Iterator) Next() string {
	it.position++
	return it.data[it.position]
}

func (it *Iterator) Close() {
	it.closed = true
}
