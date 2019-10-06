package io

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"
)

// getCSPHeader generates the CSP header value for the index html file
func (s SPA) getCSPHeader() (string, error) {
	baseCSP := "script-src 'self' "
	hashes := []string{}

	indexFile, err := s.FileSystem.Open(s.Index)
	if err != nil {
		return "", err
	}

	out, err := ioutil.ReadAll(indexFile)
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(bytes.NewReader(out))
	if err != nil {
		return "", err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil {
			hashes = append(hashes, generateNodeHash(n.FirstChild))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return baseCSP + strings.Join(hashes, " "), nil
}

func generateNodeHash(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	str := html.UnescapeString(buf.String())
	sum := sha256.Sum256([]byte(str))
	sha1Hash := base64.StdEncoding.EncodeToString(sum[:])

	return fmt.Sprintf("'sha256-%s'", sha1Hash)
}
