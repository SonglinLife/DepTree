package depstr

type Queue struct{
	nodes []interface{}
}

func NewQueue() *Queue{
	return &Queue{
		nodes: []interface{}{},
	}
}

func (q *Queue) Push(node interface{}) {
	q.nodes = append(q.nodes, node)
}

func (q *Queue)Pop()interface{}  {
	node := q.nodes[0]
	q.nodes = q.nodes[1:]
	return node
}

func (q *Queue) Empty()bool{
	return len(q.nodes) == 0
}



type Stack struct {
	nodes []interface{}
}

func NewStack() *Stack {
	return &Stack{
		nodes: []interface{}{},
	}
}

func (s *Stack) Push(node interface{}) {
	s.nodes = append(s.nodes, node)
}

func (s *Stack) Pop() (node interface{}) {
	l := len(s.nodes)
	node = s.nodes[l-1]
	s.nodes = s.nodes[:l-1]
	return node
}

func (s *Stack) Empty() bool {
	return len(s.nodes) == 0
}