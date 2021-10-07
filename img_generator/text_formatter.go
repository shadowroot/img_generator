package img_generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	svg "github.com/ajstarks/svgo"
	"github.com/common-nighthawk/go-figure"
)

var (
	gwfURI  = "http://fonts.googleapis.com/css?family="
	fontfmt = "<style type=\"text/css\">\n<![CDATA[\n%s]]>\n</style>\n"
	//gfmt      = "white;font-size:36pt;text-anchor:middle"
	fontCachePath                 = "../fonts/"
	fontCacheFilemode os.FileMode = 0666
	fontCacheFileEnd              = ".tmp"
)

func AsciiART(text string, w io.Writer) error {
	figure.Write(w, figure.NewFigure(text, "puffy", true))
	return nil
}

func AsciiARTFont(text string, font string, w io.Writer) error {
	figure.Write(w, figure.NewFigure(text, font, true))
	return nil
}

func ObtainFont(font string, canvas *svg.SVG) error {
	font = strings.ToLower(font)
	if content, err := ioutil.ReadFile(fontCachePath + font + fontCacheFileEnd); err == nil {
		fntDef := string(content)
		canvas.Def()
		fmt.Fprint(canvas.Writer, fntDef)
		canvas.DefEnd()
		return nil
	}
	r, err := http.Get(gwfURI + url.QueryEscape(font))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	if err != nil || r.StatusCode != http.StatusOK {
		return err
	}

	buff := &bytes.Buffer{}
	if _, err := fmt.Fprintf(buff, fontfmt, b); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fontCachePath+font+fontCacheFileEnd, buff.Bytes(), fontCacheFilemode); err != nil {
		return err
	}
	canvas.Def()
	fmt.Fprint(canvas.Writer, buff.Bytes())
	canvas.DefEnd()
	return nil
}

func FontStyleTextDrawable(canvas *svg.SVG, fontfamily string, fontsize int, fontunits string, color string, options []string) (*TextDrawable, error) {
	textTemplate := &TextDrawable{
		text:      "",
		x:         0,
		y:         0,
		font:      fontfamily,
		options:   options,
		spacing:   fontsize + 10,
		fillColor: color,
	}
	if err := ObtainFont(fontfamily, canvas); err != nil {
		return nil, err
	}
	if textTemplate.options == nil {
		textTemplate.options = make([]string, 1)
	}

	if fontunits == "" {
		fontunits = "pt"
	}

	textTemplate.options = append(textTemplate.options, fmt.Sprintf("font-size: %v %v", fontsize, fontunits))
	textTemplate.options = append(textTemplate.options, "font-family: "+fontfamily)
	return textTemplate, nil
}

func TextFormatFromTemplate(tpl *template.Template, values map[string][]string) (string, error) {
	buff := &bytes.Buffer{}
	err := tpl.Execute(buff, values)
	return buff.String(), err
}
