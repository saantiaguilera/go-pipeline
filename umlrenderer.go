package pipeline

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
)

const (
	// Base64 Encoding maps
	mapper = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	// Base URL for plant UML graphs creation
	baseURL = "http://www.plantuml.com/plantuml"

	// UMLFormatPNG OutputFormat for graph renderings (a UMLFormatPNG image will be created)
	UMLFormatPNG UMLOutputFormat = "png"
	// UMLFormatSVG OutputFormat for graph renderings (an UMLFormatSVG image will be created)
	UMLFormatSVG UMLOutputFormat = "svg"
	// UMLFormatRaw OutputFormat for graph renderings (a file with the raw contents will be created)
	UMLFormatRaw UMLOutputFormat = "raw"
	// UMLFormatTXT OutputFormat for graph renderings (an ASCII Art will be created)
	UMLFormatTXT UMLOutputFormat = "txt"
)

type (
	// UMLOutputFormat for graph renderings
	UMLOutputFormat string

	// UMLOptions available when drawing a graph
	UMLOptions struct {
		// Type of the drawing, by default we will use UMLFormatSVG
		Type UMLOutputFormat
		// Base URL to use for retrieving Plant UML graphs, by default we will use http://www.plantuml.com/plantuml/
		BaseURL string
	}

	// UMLRenderer allows us to render graphs into an UML diagram output
	UMLRenderer struct {
		Options UMLOptions
	}

	umlRendererGraph interface {
		Graph

		String() string
	}
)

// NewUMLRenderer creates an UML renderer for drawing graphs as specified
func NewUMLRenderer(options UMLOptions) *UMLRenderer {
	if len(options.Type) == 0 {
		options.Type = UMLFormatSVG
	}

	if len(options.BaseURL) == 0 {
		options.BaseURL = baseURL
	}

	return &UMLRenderer{
		Options: options,
	}
}

// Render draws in UML activity the given step, and writes it to the given file
func (u *UMLRenderer) Render(graphDiagram umlRendererGraph, output io.WriteCloser) error {
	content := graphDiagram.String()

	if u.Options.Type == UMLFormatRaw {
		_, err := io.WriteString(output, content)
		return err
	}
	return u.renderUml([]byte(content), output)
}

// Render as  UML the contents, writing them into the File
func (u *UMLRenderer) renderUml(content []byte, output io.WriteCloser) error {
	content = u.deflate(content)
	url := fmt.Sprintf("%s/%s/~1%s", u.Options.BaseURL, u.Options.Type, u.base64Encode(content))

	response, err := http.Get(url)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("status code %d while trying to create the graph through %s", response.StatusCode, url)
	}

	_, err = io.Copy(output, response.Body)

	if err != nil {
		return err
	}

	return output.Close()
}

// Encode in standard B64 the given input
func (u *UMLRenderer) base64Encode(input []byte) string {
	var buffer bytes.Buffer
	inputLength := len(input)
	for i := 0; i < 3-inputLength%3; i++ {
		input = append(input, byte(0))
	}

	for i := 0; i < inputLength; i += 3 {
		b1, b2, b3 := input[i], input[i+1], input[i+2]

		b4 := b3 & 0x3f
		b3 = ((b2 & 0xf) << 2) | (b3 >> 6)
		b2 = ((b1 & 0x3) << 4) | (b2 >> 4)
		b1 = b1 >> 2

		for _, b := range []byte{b1, b2, b3, b4} {
			buffer.WriteByte(mapper[b])
		}
	}
	return buffer.String()
}

// Deflate compression algorithm
func (u *UMLRenderer) deflate(content []byte) []byte {
	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	_, _ = w.Write(content)
	_ = w.Close()
	return b.Bytes()
}
