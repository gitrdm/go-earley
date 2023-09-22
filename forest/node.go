package forest

type Node interface {
	node()
	Origin() int
	Location() int
}

type Group interface {
	Children() []Node
}

type group struct {
	children []Node
}

func NewGroup(children ...Node) Group {
	return &group{
		children: children,
	}
}

func (g group) Children() []Node {
	return g.children
}

type internal struct {
	alternatives []Group
}

type Internal interface {
	Alternatives() []Group
	AddUniqueFamily(w, v Node)
}

func (i *internal) AddUniqueFamily(w, v Node) {
	childCount := 1
	if v != nil {
		childCount += 1
	}
	for _, group := range i.alternatives {

		if len(group.Children()) != childCount {
			continue
		}
		if i.isMatchedSubtree(w, v, group) {
			return
		}
	}

	group := &group{}
	group.children = append(group.children, w)
	if childCount > 1 {
		group.children = append(group.children, v)
	}
	i.alternatives = append(i.alternatives, group)
}

func (i *internal) isMatchedSubtree(first, second Node, group Group) bool {

	firstCompare := group.Children()[0]

	if first != firstCompare {
		return false
	}

	if second == nil {
		return true
	}

	secondCompare := group.Children()[1]

	return secondCompare == second
}
