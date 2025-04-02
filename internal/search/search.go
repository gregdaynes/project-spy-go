package search

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gosimple/slug"
	"golang.org/x/net/html"
	"projectspy.dev/internal/task"
)

type searchEntry []string

func Data(taskLanes task.Lanes) string {
	searchData := make([]searchEntry, 0)

	for i := 0; i < len(taskLanes); i++ {
		lane := taskLanes[i]

		for j := 0; j < len(lane.Tasks); j++ {
			t := lane.Tasks[j]

			entry := searchEntry{}
			entry = append(entry, strings.ToLower(t.Title+" "+stripTags(t.Description)))
			entry = append(entry, slug.Make(lane.Name+"-"+t.Filename))
			searchData = append(searchData, entry)
		}
	}
	searchJSON, _ := json.Marshal(searchData)
	return string(searchJSON)
}

func renderNode(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(c, buf)
	}
}

// stripTags removes HTML tags from a string.
func stripTags(htmlStr string) string {
	doc, err := html.Parse(bytes.NewReader([]byte(htmlStr)))
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	renderNode(doc, &buf)
	return buf.String()
}
