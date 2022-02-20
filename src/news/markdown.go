package news

import (
	"bytes"
	"io"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var (
	githubLinkPrefix = []byte("https://github.com/xvello/letsblockit/")
)

type releaseNoteRenderer struct {
	*blackfriday.HTMLRenderer
}

func (r *releaseNoteRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.HorizontalRule:
		// Stop parsing on the first hr, items below are not relevant for end users
		return blackfriday.Terminate
	case blackfriday.Link:
		// Match github pull/commit links and shorten the anchor text
		if bytes.HasPrefix(node.LinkData.Destination, githubLinkPrefix) && node.FirstChild != nil {
			textNode := node.FirstChild
			linkParts := strings.Split(string(node.LinkData.Destination), "/")

			if len(linkParts) > 2 && textNode != nil && textNode.Type == blackfriday.Text {
				linkedType := linkParts[len(linkParts)-2]
				linkedId := linkParts[len(linkParts)-1]
				switch linkedType {
				case "pull", "issues":
					node.FirstChild.Literal = append([]byte("#"), linkedId...)
				case "commit":
					if len(linkedId) >= 7 {
						node.FirstChild.Literal = []byte(linkedId[0:7])
					}
				}
			}
		}
	}
	return r.HTMLRenderer.RenderNode(w, node, entering)
}

func initRenderer() blackfriday.Option {
	return blackfriday.WithRenderer(
		&releaseNoteRenderer{
			HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				HeadingLevelOffset: 1,
			}),
		})
}
