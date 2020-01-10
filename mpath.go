package mpath

import (
	"fmt"
)

//Leafs is Leaf collection interface
type Leafs interface {
	GetLeafs() []Leaf
}

//Leaf interface for making tree from flat Leafs slice
type Leaf interface {
	GetID() uint64
	SetID(uint64)
	GetPosition() int
	SetPosition(int)
	GetPath() []uint64
	GetPathFromIdx(*int) []uint64
	GetLeafByID(leafs Leafs, id uint64) Leaf
	GetLeafOrMakeNew(leafs Leafs, id uint64, position int) Leaf
	GetSiblings() Leafs
	AppendSiblings(interface{})
	MakeRoot(leafs Leafs, leaf Leaf) (pathIndex *int)
	GetRootPathIndex() *int
}

//InitTree creates Leaf tree (index) from flat Leafs slice
func InitTree(tree Leaf, leafs Leafs) error {
	for _, leaf := range leafs.GetLeafs() {
		if err := parsePath(tree, leafs, leaf); err != nil {
			return err
		}
	}
	return nil
}

func parsePath(index Leaf, items Leafs, leaf Leaf) error {
	rootPathIdx := index.MakeRoot(items, leaf)
	if rootPathIdx == nil {
		return fmt.Errorf("no root for path for itemID %d", leaf.GetID())
	}

	itemRootID := leaf.GetPath()[*rootPathIdx]

	// pass child items adding because of wrong root element
	if index.GetID() != itemRootID {
		return nil
	}

	addSibling(index, items, leaf.GetPathFromIdx(rootPathIdx), leaf.GetPosition())

	return nil
}

func addSibling(parent Leaf, leafs Leafs, path []uint64, position int) {
	if len(path) == 0 {
		return
	}

	index := getSiblingsIndex(parent, leafs, path[0], position)

	if len(path) > 1 {
		addSibling(index, leafs, path[1:], position)
	}
}

func getSiblingsIndex(parent Leaf, leafs Leafs, id uint64, position int) Leaf {
	for _, leaf := range parent.GetSiblings().GetLeafs() {
		if leaf.GetID() == id {
			return leaf
		}
	}

	parent.AppendSiblings(parent.GetLeafOrMakeNew(leafs, id, position))

	return getSiblingsIndex(parent, leafs, id, position)
}
