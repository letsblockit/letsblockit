package news

import (
	"bytes"
	"io"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var (
	githubLinkPrefix      = []byte("https://github.com/")
	githubRepoName        = "letsblockit"
	templateNameSeparator = []byte(":")
	templateLinkPrefix    = []byte("https://letsblock.it/filters/")
)

type releaseNoteRenderer struct {
	officialInstance bool
	templateExists   templateExists
	*blackfriday.HTMLRenderer
}

func (r *releaseNoteRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type { //nolint:exhaustive
	case blackfriday.HorizontalRule:
		if r.officialInstance {
			// Stop parsing on the first hr, items below are not relevant on the official instance
			return blackfriday.Terminate
		}
	case blackfriday.Item:
		// Match for template names at the start of list items and make them a link
		hasText := node.FirstChild != nil && node.FirstChild.Type == blackfriday.Paragraph &&
			node.FirstChild.FirstChild != nil && node.FirstChild.FirstChild.Type == blackfriday.Text
		if hasText {
			textNode := node.FirstChild.FirstChild

			maybeName, _, foundSeparator := bytes.Cut(textNode.Literal, templateNameSeparator)
			if foundSeparator && r.templateExists(string(maybeName)) {
				linkNode := blackfriday.NewNode(blackfriday.Link)
				linkNode.LinkData.Destination = append(linkNode.LinkData.Destination, templateLinkPrefix...)
				linkNode.LinkData.Destination = append(linkNode.LinkData.Destination, maybeName...)

				linkTextNode := blackfriday.NewNode(blackfriday.Text)
				linkTextNode.Literal = maybeName
				linkNode.AppendChild(linkTextNode)

				textNode.InsertBefore(linkNode)
				textNode.Literal = textNode.Literal[len(maybeName):]
			}
		}
	case blackfriday.Link:
		// Match github pull/commit links and shorten the anchor text
		if bytes.HasPrefix(node.LinkData.Destination, githubLinkPrefix) && node.FirstChild != nil {
			textNode := node.FirstChild
			linkParts := strings.Split(string(bytes.TrimPrefix(node.LinkData.Destination, githubLinkPrefix)), "/")

			if len(linkParts) >= 4 && textNode != nil && textNode.Type == blackfriday.Text {
				linkedType := linkParts[len(linkParts)-2]
				linkedId := linkParts[len(linkParts)-1]
				repoName := linkParts[1]
				sameRepo := repoName == githubRepoName

				switch linkedType {
				case "pull", "issues":
					if sameRepo {
						node.FirstChild.Literal = append([]byte("#"), linkedId...)
					} else {
						node.FirstChild.Literal = append([]byte(repoName+"#"), linkedId...)
					}
				case "commit":
					if sameRepo {
						node.FirstChild.Literal = append([]byte("@"), []byte(linkedId[0:7])...)
					} else {
						node.FirstChild.Literal = append([]byte(repoName+"@"), []byte(linkedId[0:7])...)
					}
				case "compare":
					if sameRepo {
						node.FirstChild.Literal = []byte(linkedId)
					}
				}
			}
		}
	}
	return r.HTMLRenderer.RenderNode(w, node, entering)
}

func initRenderer(officialInstance bool, tp templateExists) blackfriday.Option {
	return blackfriday.WithRenderer(
		&releaseNoteRenderer{
			officialInstance: officialInstance,
			templateExists:   tp,
			HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				HeadingLevelOffset: 2,
			}),
		})
}
