/*
Package mpath is golang realisation of MPTT (modified preorder tree traversal) in materialized path way. It includes
interfaces with methods that objects should implement to make an ordered tree from a flat slice of them.
*/
package mpath

import (
	"fmt"
)

//Leafs is Leaf collection interface
type Leafs interface {

	//GetLeafs return Leafs slice
	GetLeafs() []Leaf
}

//Leaf interface for making tree from flat Leafs slice
type Leaf interface {

	//GetID returns an uint64 ID of Leaf object
	GetID() uint64

	//SetID sets ID to an empty Leaf object
	SetID(id uint64)

	//GetPosition returns position of Leafs branch
	GetPosition() int

	//SetPosition sets position of Leaf in its branch
	SetPosition(pos int)

	//GetPath returns Leafs path as Leafs IDs slice
	GetPath() []uint64

	//GetPathFromIdx returns Leaf path chunk started from passed index of element in path slice
	GetPathFromIdx(index *int) []uint64

	//GetLeafByID returns Leaf from Leafs by its ID
	GetLeafByID(leafs Leafs, id uint64) Leaf

	//GetLeafOrMakeNew returns Leaf from Leafs by its id or makes new Leaf object with ID and Position property
	GetLeafOrMakeNew(leafs Leafs, id uint64, position int) Leaf

	//GetSiblings return Leaf siblings as Leafs
	GetSiblings() Leafs

	//AppendSiblings append Leafs siblings to current Leaf
	AppendSiblings(interface{})

	//MakeRoot creates root Leaf to a Leafs tree and return path index of this root Leaf in Leafs path
	MakeRoot(leafs Leafs, leaf Leaf) (pathIndex *int)

	//GetRootPathIndex returns index of root Leaf id in Leaf path slice
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
