package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/djherbis/times"
)

func parseFile(fp string) (task Task, err error) {
	data, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)

	ds := false

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		task.RawContents = task.RawContents + "\n" + text

		if text == "" {
			ds = true
			continue
		}

		if text == "===" {
			ds = true
			continue
		}

		if text == "---" {
			ds = true
			continue
		}

		if ds == false {
			task.Title += text + "\n"
			task.Priority = Priority(task.Title)
			task.Order = Order(task.Title)
			task.Tags = Tags(task.Title)
			continue
		}

		if ds == true {
			task.Description += text + "\n"
			continue
		}
	}

	t, err := times.Stat(fp)
	if err != nil {
		log.Fatal(err.Error())
	}

	task.ModifiedTime = t.ModTime()

	if t.HasBirthTime() {
		log.Println(t.BirthTime())
		task.CreatedTime = t.BirthTime()
	}

	relative, ok := strings.CutPrefix(fp, ".projectSpy")
	if ok != true {
		fmt.Println(relative, fp)
		log.Fatal("bad time trimming prefix")
	}
	task.RelativePath = relative

	strs := strings.Split(relative, "/")

	task.Lane = strs[1]
	task.Filename = strs[2]

	return task, nil
}

func Priority(title string) (priority int) {
	r := regexp.MustCompile(`!+`)
	s := r.FindString(title)

	return len(s)
}

func Order(title string) (order int) {
	r := regexp.MustCompile(`(\d+)`)
	o := r.FindString(title)

	order, err := strconv.Atoi(o)
	if err != nil {
		fmt.Errorf("Error getting order from title %s", title)
	}

	return order
}

func Tags(title string) (tags []string) {
	r := regexp.MustCompile(`\[([^\]]+)\]`)

	m := r.FindAllStringSubmatch(title, -1)
	for _, v := range m {
		tags = append(tags, v[1])
	}

	return tags
}
