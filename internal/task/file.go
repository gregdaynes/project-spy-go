package task

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/djherbis/times"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func ParseFile(fp string) (task Task, err error) {
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
			task.Title += ParseTitle(text)
			task.Priority = Priority(text)
			task.Order = Order(text)
			task.Tags = Tags(text)
			continue
		}
	}

	t, err := times.Stat(fp)
	if err != nil {
		log.Fatal(err.Error())
	}

	task.ModifiedTime = t.ModTime()

	if t.HasBirthTime() {
		task.CreatedTime = t.BirthTime()
	}

	relative, ok := strings.CutPrefix(fp, ".projectSpy/")
	if ok != true {
		log.Fatal("bad time trimming prefix")
	}
	task.RelativePath = relative

	strs := strings.Split(relative, "/")

	task.Lane = strs[0]
	task.Filename = strs[1]

	description, html := ParseDescription(task.RawContents)
	task.Description = description
	task.DescriptionHTML = html

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

func ParseTitle(title string) (parsedTitle string) {
	// remove priority
	r := regexp.MustCompile(`!+`)
	title = r.ReplaceAllString(title, "")

	// remove order
	r = regexp.MustCompile(`(\s\d+)`)
	title = r.ReplaceAllString(title, "")

	// remove tags
	r = regexp.MustCompile(`\[([^\]]+)\]`)
	title = r.ReplaceAllString(title, "")

	return title
}

func ParseDescription(text string) (output, outputHTML string) {
	reChangelog := regexp.MustCompile(`(?:\n---\n\n)(?:(?:\d{4}-\d{2}-\d{2} \d{2}:\d{2}\t.*)\n?)+`)
	reHeader := regexp.MustCompile(`.+\n===+\n|.+\n---+\n|#+\s.+\n`)
	output = text
	output = reChangelog.ReplaceAllString(output, "")
	output = reHeader.ReplaceAllString(output, "")
	output = strings.TrimSpace(output)
	output = strings.Split(output, "\n")[0]

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(output), &buf); err != nil {
		panic(err)
	}

	return output, buf.String()
}
