package filters

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestModeTransformer(t *testing.T) {
	input := `
! youtube-cleanup
www.youtube.com###clarify-box
www.youtube.com###chat:remove()

! youtube-mixes
www.youtube.com##ytd-browse ytd-rich-item-renderer:has(#video-title-link[href*="&start_radio=1"])
www.youtube.com##ytd-search ytd-radio-renderer:style(display: none)
www.youtube.com##ytd-watch-next-secondary-results-renderer ytd-compact-radio-renderer
www.youtube.com##ytd-player div.videowall-endscreen a[data-is-list=true]

! youtube-recommendations
www.youtube.com##.ytp-ce-element

`
	expected := `
! youtube-cleanup
www.youtube.com###clarify-box:style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)
www.youtube.com###chat:style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)

! youtube-mixes
www.youtube.com##ytd-browse ytd-rich-item-renderer:has(#video-title-link[href*="&start_radio=1"]):style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)
www.youtube.com##ytd-search ytd-radio-renderer:style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)
www.youtube.com##ytd-watch-next-secondary-results-renderer ytd-compact-radio-renderer:style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)
www.youtube.com##ytd-player div.videowall-endscreen a[data-is-list=true]:style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)

! youtube-recommendations
www.youtube.com##.ytp-ce-element:style(background: rgba(255,63,63,0.20) !important; border: 1px solid red !important)

`
	for chunkSize := 1; chunkSize < 256; chunkSize += 7 {
		t.Run(fmt.Sprintf("chunks %d", chunkSize), func(t *testing.T) {
			in := []byte(input)
			out := new(bytes.Buffer)
			trans := NewTestModeTransformer(out)

			for {
				if len(in) < chunkSize+1 {
					break
				}
				_, err := trans.Write(in[:chunkSize])
				require.NoError(t, err)
				in = in[chunkSize:]
			}
			_, err := trans.Write(in)
			require.NoError(t, err)

			assert.Equal(t, expected, out.String())
		})
	}
}
