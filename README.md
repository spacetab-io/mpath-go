# mpath-go

Golang realisation of MPTT (or modified preorder tree traversal) in materialized path way.

## About

It provides interfaces which yor database object should implement.

Your database object should store:
* `path` property as slice of uint64 IDs of materialized path to this object in traversal tree;
* `position` property as integer for determine the order of leafs in tree 

## Usage

Implementation example and tests are in [test file](/mpath_test.go).

```go
package main

import (
    "fmt"

    "github.com/spacetab-io/mpath"
)

type TestItems []*TestItem

type TestItem struct {
    ID       uint64
    Path     []uint64
    Position int
    Siblings TestItems
    Name     string
}

// Leaf interface implementation for TestItem
// ...
// Leafs interface implementation for TestItems
// ...

func main() {
    flatItemsSlice := getTestItems()

    var parent = TestItem{}
    if err := mpath.InitTree(&parent, flatItemsSlice); err != nil {
        panic("error tree init")
    }
    
    fmt.Print(parent)
}

func getTestItems() *TestItems {
    return &TestItems{
        {ID: 1, Position: 0, Name: "item 1", Path: []uint64{1}},
        {ID: 2, Position: 0, Name: "item 2", Path: []uint64{1, 2}},
        {ID: 3, Position: 1, Name: "item 3", Path: []uint64{1, 3}},
        {ID: 4, Position: 0, Name: "item 4", Path: []uint64{1, 2, 4}},
        {ID: 5, Position: 1, Name: "item 5", Path: []uint64{1, 2, 5}},
        {ID: 6, Position: 0, Name: "item 6", Path: []uint64{1, 3, 6}},
    }
}
```

## Tests

    go test ./... -v -race
