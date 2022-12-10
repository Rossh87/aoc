package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 1. build tree from input
// 2. recursive dfs to get sum of every dir.  If sum <= 100,000, add to total
// calc size of every directory.  If size <= 100,000, add size to total sum

type fileType string

const (
	file fileType = "f"
	dir  fileType = "d"
)

type tree struct {
	kind     fileType
	path     string
	children *[]tree
	size     int64
	parent   *tree
}

func getDirSizes(t tree) []tree {
	var setSizes func(node tree) int64

	res := []tree{}

	setSizes = func(node tree) int64 {
		currNode := node

		for _, child := range *node.children {
			if child.kind == dir {
				s := setSizes(child)
				currNode.size += int64(s)
				continue
			}

			if child.kind == file {
				currNode.size += child.size
			}
		}

		res = append(res, currNode)
		return currNode.size
	}

	setSizes(t)

	return res
}

func sumDirs(dirs []tree) int64 {
	var sum int64 = 0

	for _, d := range dirs {
		if d.size <= 100000 {
			sum += d.size
		}
	}

	return sum
}

type fs struct {
	tree     *tree
	currNode *tree
}

func (f *fs) cd(dest string) {
	if dest == ".." {
		// if parent is nil, we're at root, so just stay there
		if f.currNode.parent == nil {
			return
		}

		f.currNode = f.currNode.parent

	} else {
		childPath := childPathName(*f.currNode, dest)
		// we start searching at currNode, NOT AT ROOT. Consider that
		// '/a' and '/a/b/a' are NOT THE SAME DIRECTORY.  If we always start
		// searching at root of tree, we will always get the shallowest occurence
		// of a node with the specified name, not necessarily the node that is the
		// immediate child of the current node.
		for _, child := range *f.currNode.children {

			if child.path == childPath {
				f.currNode = &child
				break
			}
		}

	}
}

func childPathName(parentNode tree, name string) string {
	if parentNode.path == "/" {
		return "/" + name
	}

	return parentNode.path + "/" + name
}

func (f *fs) appendFile(name string, size int64) {
	childPath := childPathName(*f.currNode, name)
	for _, child := range *f.currNode.children {
		if childPath == child.path {
			return
		}
	}

	*f.currNode.children = append(*f.currNode.children, tree{path: childPath, size: size, children: &[]tree{}, kind: file, parent: f.currNode})
}

func (f *fs) appendDir(name string) {
	childPath := childPathName(*f.currNode, name)
	for _, child := range *f.currNode.children {
		if childPath == child.path {
			return
		}

	}
	*f.currNode.children = append(*f.currNode.children, tree{path: childPath, children: &[]tree{}, kind: dir, parent: f.currNode})
}

func (f *fs) ingestLine(ln string) {
	fields := strings.Fields(ln)

	switch fields[0] {
	case "$":
		if fields[1] == "ls" {
			return
		}
		if fields[1] == "cd" {
			f.cd(fields[2])
		}
	case "dir":
		f.appendDir(fields[1])
	default:
		size64, _ := strconv.ParseInt(fields[0], 10, 64)

		f.appendFile(fields[1], size64)
	}
}

func (f *fs) solve() int64 {
	return sumDirs(getDirSizes(*f.tree))
}

const (
	totalAvailable int64 = 7e7
	required       int64 = 3e7
)

func (f *fs) solvePartTwo() int64 {
	dirs := getDirSizes(*f.tree)

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].size < dirs[j].size
	})

	rootUse := dirs[len(dirs)-1].size

	currUnused := totalAvailable - rootUse

	for _, d := range dirs {
		availableIfDirDeleted := currUnused + d.size

		if availableIfDirDeleted >= required {
			return d.size
		}
	}

	return 0
}

func newFS() *fs {
	root := tree{
		path:     "/",
		kind:     dir,
		size:     0,
		parent:   nil,
		children: &[]tree{},
	}

	f := fs{
		currNode: &root,
		tree:     &root,
	}

	return &f
}

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	f := newFS()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		f.ingestLine(line)
	}

	fmt.Println(f.solvePartTwo())
}
