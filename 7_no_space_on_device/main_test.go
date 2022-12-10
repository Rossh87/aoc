package main

import (
	"bufio"
	"strings"
	"testing"
)

func expect[T comparable](want, got T, t *testing.T) {
	if got != want {
		t.Fatalf("expected %v, but received %v\n%+v", want, got, got)
	}
}
func TestGetDirSizes(t *testing.T) {
	f := newFS()

	f.appendDir("a")
	f.appendFile("a.txt", 100)
	f.cd("a")
	f.appendDir("a")
	f.cd("a")
	f.appendFile("lowera.txt", 100)

	got := getDirSizes(*f.tree)

	want := []int64{100, 100, 200}

	for idx, wanted := range want {
		expect(wanted, got[idx].size, t)
	}
}

func TestSumDirs(t *testing.T) {
	given := []tree{{path: "/a", kind: dir, size: 1000, children: &[]tree{}},
		{path: "/b", kind: dir, size: 100100, children: &[]tree{}}, {path: "/c", kind: dir, size: 2000, children: &[]tree{}}}

	want := int64(3000)

	got := sumDirs(given)

	expect(want, got, t)
}

var testInput = `$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k`

func TestMain(t *testing.T) {
	f := newFS()

	scanner := bufio.NewScanner(strings.NewReader(testInput))

	for scanner.Scan() {
		ln := scanner.Text()

		f.ingestLine(ln)
	}

	expect(95437, f.solve(), t)
}

func TestMainPartTwo(t *testing.T) {
	f := newFS()

	scanner := bufio.NewScanner(strings.NewReader(testInput))

	for scanner.Scan() {
		ln := scanner.Text()

		f.ingestLine(ln)
	}

	expect(24933642, f.solvePartTwo(), t)
}
