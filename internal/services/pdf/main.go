package pdf

import (
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func HtmlToPdfBytes(content string) ([]byte, error) {
	g, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	g.AddPage(wkhtml.NewPageReader(strings.NewReader(content)))
	if err := g.Create(); err != nil {
		return nil, err
	}
	return g.Bytes(), nil
}
