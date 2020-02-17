package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	expr := regexp.MustCompile("^[[:space:]]*- ")

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	currDir := &dir{name: "."}
	currLevel := 0

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		if !expr.MatchString(line) {
			continue
		}

		indent := strings.Index(line, "- ")

		lvl := indent / 2.0
		for lvl < currLevel {
			currDir = pop(currDir)
			currLevel--
		}

		if lvl >= currLevel {
			currDir = push(currDir, line[indent+2:])
			currLevel++
		}
	}

	err = scanner.Err()
	if err != nil {
		panic(err)
	}

	for currDir.parent != nil {
		currDir = currDir.parent
	}

	walk(currDir)
}

func walk(d *dir) {

	currDir := d
	path := d.name
	for currDir.parent != nil {
		path = currDir.parent.name + string(filepath.Separator) + path
		currDir = currDir.parent
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		println(err.Error())
	}

	for _, child := range d.children {
		walk(child)
	}
}

type dir struct {
	name     string
	parent   *dir
	children []*dir
}

func push(d *dir, childName string) *dir {
	child := dir{
		parent: d,
		name:   childName,
	}
	d.children = append(d.children, &child)
	return &child
}

func pop(d *dir) *dir {
	if d.parent == nil {
		panic("nil parent: " + d.name)
	}
	return d.parent
}
