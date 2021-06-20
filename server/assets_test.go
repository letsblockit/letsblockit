package server

import (
	"encoding/base64"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type WrappedAssetsTestSuite struct {
	suite.Suite
	assets *wrappedAssets
}

// embed.Open currently returns a seekable, ensure this is true in future go versions
func (s *WrappedAssetsTestSuite) TestSeekable() {
	file, err := s.assets.Open("assets/css/styles.min.css")
	s.NoError(err)
	pos, err := file.Seek(0, io.SeekEnd)
	s.NoError(err)
	s.True(pos > 0)
}

func (s *WrappedAssetsTestSuite) TestOpenDir() {
	// Confirm we can open a stat a regular file
	f, err := s.assets.Open("assets/css/styles.min.css")
	s.NoError(err)
	info, err := f.Stat()
	s.NoError(err)
	s.Equal("styles.min.css", info.Name())

	// Confirm we cannot open its parent directory
	f, err = s.assets.Open("assets/css")
	s.Error(err)
	s.Nil(f)
}

func (s *WrappedAssetsTestSuite) TestHash() {
	out, err := base64.RawURLEncoding.DecodeString(s.assets.hash)
	s.NoError(err)
	s.Len(out, 16) // 128 bits hash expected
}

func (s *WrappedAssetsTestSuite) TestETag() {
	expectedETag := "\"" + s.assets.hash + "\""
	s.Equal(expectedETag, s.assets.eTag)
}

func TestWrappedAssetsTestSuite(t *testing.T) {
	suite.Run(t, &WrappedAssetsTestSuite{
		assets: loadAssets(),
	})
}
