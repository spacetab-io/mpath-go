package mpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	ID       uint64
	Path     []uint64
	Position int
	Siblings TestItems
	Name     string
}

type TestItems []*TestItem

func (ii TestItems) GetLeafs() []Leaf {
	var items []Leaf
	for _, i := range ii {
		items = append(items, i)
	}

	return items
}

func (ti *TestItem) GetID() uint64 {
	return ti.ID
}

func (ti *TestItem) SetID(id uint64) {
	ti.ID = id
}

func (ti TestItem) GetPath() []uint64 {
	return ti.Path
}

func (ti TestItem) GetPathFromIdx(idx *int) []uint64 {
	return ti.Path[*idx+1:]
}

func (ti TestItem) GetPosition() int {
	return ti.Position
}

func (ti *TestItem) SetPosition(position int) {
	ti.Position = position
}

func (ti *TestItem) GetSiblings() Leafs {
	return ti.Siblings
}

func (ti *TestItem) AppendSiblings(index interface{}) {
	idx := index.(*TestItem)
	ti.Siblings = append(ti.Siblings, idx)
}

func (ti *TestItem) GetLeafByID(items Leafs, id uint64) Leaf {
	for _, item := range items.GetLeafs() {
		if item.GetID() == id {
			return item
		}
	}

	return nil
}

func (ti *TestItem) GetLeafOrMakeNew(items Leafs, id uint64, position int) Leaf {
	item := ti.GetLeafByID(items, id)

	if item == nil {
		item = &TestItem{ID: id, Position: position}
	}

	return item
}

func (ti *TestItem) MakeRoot(items Leafs, item Leaf) *int {
	if ti.GetID() != 0 {
		return ti.GetRootPathIndex()
	}
	for idx, id := range item.GetPath() {
		root := ti.GetLeafByID(items, id)
		if root != nil {
			*ti = *root.(*TestItem)
			return &idx
		}
	}

	return nil
}

func (ti *TestItem) GetRootPathIndex() *int {
	for idx, id := range ti.GetPath() {
		if id == ti.GetID() {
			return &idx
		}
	}

	return nil
}

var (
	itemsFull = &TestItems{
		{ID: 1, Position: 0, Name: "item 1", Path: []uint64{1}},
		{ID: 2, Position: 0, Name: "item 2", Path: []uint64{1, 2}},
		{ID: 3, Position: 1, Name: "item 3", Path: []uint64{1, 3}},
		{ID: 4, Position: 0, Name: "item 4", Path: []uint64{1, 2, 4}},
		{ID: 5, Position: 1, Name: "item 5", Path: []uint64{1, 2, 5}},
		{ID: 6, Position: 0, Name: "item 6", Path: []uint64{1, 3, 6}},
	}
	itemsPart = &TestItems{
		{ID: 2, Position: 0, Name: "item 2", Path: []uint64{1, 2}},
		{ID: 4, Position: 0, Name: "item 4", Path: []uint64{1, 2, 4}},
		{ID: 5, Position: 1, Name: "item 5", Path: []uint64{1, 2, 5}},
		{ID: 6, Position: 0, Name: "item 6", Path: []uint64{1, 3, 6}},
	}
	itemsError = &TestItems{
		{ID: 7, Position: 0, Name: "item 2", Path: []uint64{1, 2}},
		{ID: 8, Position: 0, Name: "item 4", Path: []uint64{1, 2, 4}},
		{ID: 9, Position: 0, Name: "item 6", Path: []uint64{1, 3, 6}},
	}
	testTreeFull = &TestItem{
		Position: 0,
		ID:       1,
		Name:     "item 1",
		Path:     []uint64{1},
		Siblings: []*TestItem{
			{
				Position: 0,
				ID:       2,
				Name:     "item 2",
				Path:     []uint64{1, 2},
				Siblings: []*TestItem{
					{
						Position: 0,
						ID:       4,
						Name:     "item 4",
						Path:     []uint64{1, 2, 4},
						Siblings: nil,
					},
					{
						Position: 1,
						ID:       5,
						Name:     "item 5",
						Path:     []uint64{1, 2, 5},
						Siblings: nil,
					},
				},
			},
			{
				Position: 1,
				ID:       3,
				Name:     "item 3",
				Path:     []uint64{1, 3},
				Siblings: []*TestItem{
					{
						Position: 0,
						ID:       6,
						Name:     "item 6",
						Path:     []uint64{1, 3, 6},
						Siblings: nil,
					},
				},
			},
		},
	}
	testTreePart = &TestItem{
		Position: 0,
		ID:       2,
		Name:     "item 2",
		Path:     []uint64{1, 2},
		Siblings: []*TestItem{
			{
				Position: 0,
				ID:       4,
				Name:     "item 4",
				Path:     []uint64{1, 2, 4},
				Siblings: nil,
			},
			{
				Position: 1,
				ID:       5,
				Name:     "item 5",
				Path:     []uint64{1, 2, 5},
				Siblings: nil,
			},
		},
	}
)

func TestInitTree(t *testing.T) {
	t.Parallel()
	t.Run("init tree full", func(t *testing.T) {
		var parent = TestItem{}
		err := InitTree(&parent, itemsFull)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		assert.Equal(t, testTreeFull, &parent)
	})
	t.Run("init tree part", func(t *testing.T) {
		var parent = TestItem{}
		err := InitTree(&parent, itemsPart)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		assert.Equal(t, testTreePart, &parent)
	})
	t.Run("init tree error", func(t *testing.T) {
		var index = TestItem{}
		err := InitTree(&index, itemsError)
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}
