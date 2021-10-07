package img_generator

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	rnd "math/rand"
	"strings"
	"time"

	"github.com/boombuler/barcode"

	"github.com/rs/zerolog/log"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
)

var (
	width                     = 980
	heigth                    = 560
	iterations                = 300 // 100 300 400
	shape                     = "ellipse"
	brushBaseWidth            = 50
	brushBaseHeigth           = 20
	brushSizeCoef             = float64(0.2)
	businessCardRatio         = 4 / 3
	logoRatio                 = 1
	rotation                  = true
	pixelArt                  = 5
	timeout                   = time.Minute
	id                int64   = 0
	blur              float64 = 6
	debug                     = true
)

type IMGBusinessCard struct {
	contact   *VCardContact
	text      string
	textHash  string
	imgParams *IMGParams
	qrSize    int
}

type IMGlogo struct {
	text      string
	textHash  string
	imgParams *IMGParams
}

type IMG struct {
	text      string
	textHash  string
	imgParams *IMGParams
}

type DrawParams struct {
	iterations int
	shape      string
	//percentageVariation float32
	// baseWidth float32
	// baseHeigth float32
	// rotation bool
	options          map[string]string //part of options
	drawingModelName string            //as key in Redis
}

type DrawMove struct {
	x       int
	y       int
	width   int
	heigth  int
	color   string
	options string
	rotate  float64
}

type QRLocation struct {
	startX   int
	startY   int
	location string //eq. business_card_center
}

type QRIMGP struct {
	imgp        *IMGParams
	blockSize   int
	qrLocation  QRLocation
	fillColor   string
	unfillColor string
	text        string
	transparent bool
}

type VCARDIMGP struct {
	qs           QRIMGP
	textTemplate TextDrawable
	lines        []string
}

type IMGParams struct {
	width         int
	heigth        int
	colorPalette  ColorGenerator
	svg           *svg.SVG
	OutputBuffer  *bytes.Buffer
	timeout       time.Duration //20 sec default
	identityUid   int64
	generatedTime time.Time
	drawParams    *DrawParams
	applyEffects  bool
	effects       []string
	opacity       bool
}

type TextDrawable struct {
	x         int
	y         int
	text      string
	font      string
	options   []string
	fillColor string
	align     string
	spacing   int
}

func logoTextTransform(text string) string {
	strs := strings.Fields(text)
	if len(strs) > 1 {
		ret := ""
		for _, str := range strs {
			ret += string(str[0]) + " "
		}
		return strings.ToUpper(ret)
	}
	if rand.Intn(1) > 0 {
		return strings.ToUpper(string(text[0]))
	}
	return text
}

func CreateEllipseDrawer() *DrawParams {
	return &DrawParams{
		iterations: iterations,
		shape:      shape,
		options: map[string]string{
			"brushSizeCoef": fmt.Sprintf("%v", brushSizeCoef),
		},
		drawingModelName: "",
	}
}

func CreateImage(width int, heigth int, colorText string, opacity bool) *IMGParams {
	outBuff := &bytes.Buffer{}
	cId := id
	id++
	return &IMGParams{
		width:         width,
		heigth:        heigth,
		colorPalette:  ParseTextInputString(colorText),
		drawParams:    CreateEllipseDrawer(),
		generatedTime: time.Now(),
		timeout:       timeout,
		identityUid:   cId,
		OutputBuffer:  outBuff,
		svg:           svg.New(outBuff),
		opacity:       opacity,
		effects:       make([]string, 1),
		applyEffects:  true,
	}
}

func (imgp IMGParams) CreateError(err string) error {
	e := errors.New(err)
	return e
}

func (imgp IMGParams) ErrorHandler(err error) {
	log.Error().Msgf("[ERR] App error=%v, app=%v", err, imgp)
}

func (imgp *IMGParams) getRandomCoords() (int, int) {
	return imgp.width + 50 - rand.Intn(imgp.width+50), imgp.heigth + 50 - rand.Intn(imgp.heigth+50)
}

func (imgp *IMGParams) drawToolParams() (*DrawMove, error) {
	var err error
	var drawMove *DrawMove
	if imgp.drawParams.drawingModelName != "" {
		//rnn
		//key like ImgDrawingModel:rnn:<model name> inside json with weights on neurons

	} else {
		drawMove, err = imgp.randomizedEllipseParams()
	}
	return drawMove, err
}

func (imgp *IMGParams) drawShape(drawMove *DrawMove) error {
	var err error = nil
	switch imgp.drawParams.shape {
	default:
		err = imgp.ellipseDrawing(drawMove)
	}
	return err
}

func (imgp *IMGParams) blurFilter(blur float64) error {
	imgp.svg.Def()
	filterName := "blurFilter"
	imgp.svg.Filter(filterName, "width='100%' heigth='100%'")
	imgp.effects = append(imgp.effects, filterName)
	//SourceGraphic , SourceAlpha
	imgp.svg.FeGaussianBlur(svg.Filterspec{In: "SourceGraphic", Result: "blur"}, blur, blur)
	imgp.svg.Fend()
	imgp.svg.DefEnd()
	return nil
}

func (imgp *IMGParams) initializeFilters() error {
	return imgp.blurFilter(blur)
}

func (imgp *IMGParams) randomizedEllipseParams() (*DrawMove, error) {
	ellipseWidth := brushBaseWidth + int(float64(imgp.width)*rand.Float64()*brushSizeCoef)
	ellipseHeigth := brushBaseHeigth + int(float64(imgp.heigth)*rand.Float64()*brushSizeCoef)
	x, y := imgp.getRandomCoords()
	color, _ := imgp.colorPalette.GenerateTematicColorString()
	log.Debug().Msgf("Color palette: %v", imgp.colorPalette)

	return &DrawMove{
		x:      x,
		y:      y,
		width:  ellipseWidth,
		heigth: ellipseHeigth,
		rotate: rnd.Float64() * 360.0,
		color:  color,
	}, nil
}

func (imgp *IMGParams) ellipseDrawing(drawMove *DrawMove) error {

	//Bad impl
	options := "fill:" + drawMove.color + ";"
	optional_filters := ""
	if imgp.opacity {
		options += fmt.Sprintf(" fill-opacity: %v;", imgp.colorPalette.GetOpacity())
	}
	//Filter apply on each, huge perf degradation in rendering

	if imgp.applyEffects {
		for _, eff := range imgp.effects {
			optional_filters += " filter: url(#" + eff + ");"
		}
	}

	opts := options
	opts += optional_filters
	//Probabilistic filter apply
	/*
		if rand.Float32() > 0.5 {
			opts += optional_filters
		}
	*/

	if drawMove.rotate != 0 {
		//imgp.svg.Rotate(drawMove.rotate)
		imgp.svg.Ellipse(drawMove.x, drawMove.y, drawMove.width, drawMove.heigth, opts, fmt.Sprintf("transform=\"rotate(%v)\"", drawMove.rotate))
		//imgp.svg.Gend()
	} else {
		imgp.svg.Ellipse(drawMove.x, drawMove.y, drawMove.width, drawMove.heigth, opts)
	}

	return nil
}

func (imgp *IMGParams) drawBrush() error {
	drawMove, err := imgp.drawToolParams()
	if err != nil {
		return err
	}
	log.Info().Msgf("DrawMove: %v", drawMove)
	if drawMove != nil {
		imgp.drawShape(drawMove)
	} else {
		return imgp.CreateError("DrawMove params weren't created.")
	}
	return nil
}

func (imgp *IMGParams) createCanvas() error {
	imgp.OutputBuffer = &bytes.Buffer{}
	imgp.svg = svg.New(imgp.OutputBuffer)
	if imgp.width <= 0 || imgp.heigth <= 0 {
		return errors.New("Bad dimensions")
	}
	return nil
}

func (imgp *IMGParams) startDrawing() error {
	err := imgp.createCanvas()
	if err != nil {
		return err
	}
	imgp.svg.Start(imgp.width, imgp.heigth)
	log.Debug().Msg("Drawing started.")
	return nil
}

func (imgp *IMGParams) endDrawing() error {
	imgp.svg.End()
	log.Debug().Msg("Drawing stopped.")
	return nil
}

func (imgp *IMGParams) generateImage() error {
	log.Info().Msg("Started drawing image")
	dp := *imgp.drawParams
	if err := imgp.initializeFilters(); err != nil {
		log.Error().Msgf("ERR: %v", err)
		return nil
	}
	bgColor, _ := imgp.colorPalette.GenerateTematicColorString()
	imgp.svg.Rect(0, 0, imgp.width, imgp.heigth, "fill:"+bgColor)
	imgp.svg.Gid("Ellipses")
	for i := 0; i < dp.iterations; i++ {
		imgp.drawBrush()
	}
	//Testing blur filter
	//imgp.svg.Rect(0, 0, imgp.width, imgp.heigth, "filter: url(#blurFilter);")
	imgp.svg.Gend()
	//imgp.svg.Use(0, 0, "Ellipses", `filter="url(#blurFilter)"`)
	log.Info().Msg("Stopped drawing image")
	return nil
}

func (qs *QRIMGP) drawQR() error {
	qrCode, err := qr.Encode(qs.text, qr.M, qr.Auto)
	if err != nil {
		return err
	}
	// Write QR code to SVG
	return qs.drawBlockQr(qrCode)
}

func (imgp *IMGParams) drawVCard(vCard VCardContact, drawableTemplate TextDrawable) error {
	lines, err := vCard.ToLinesFormat()
	if err != nil {
		return err
	}
	log.Debug().Msgf("Drawing VCARD vCard=%v, textTemplate=%v", vCard, drawableTemplate)
	imgp.svg.Textlines(drawableTemplate.x, drawableTemplate.y, lines, imgp.width/2-drawableTemplate.x, drawableTemplate.spacing, drawableTemplate.fillColor, drawableTemplate.align)

	return nil
}

func (qs *QRIMGP) calculateTextCoords(text *TextDrawable) error {
	switch qs.qrLocation.location {
	case "business_card_center":
		text.x = qs.imgp.width/2 + 5
		text.y = 10
		break
	}
	return nil
}

func (imgp *IMGParams) drawText(text TextDrawable) error {
	opts := ""
	for _, opt := range text.options {
		opts += opt
	}
	imgp.svg.Text(text.x, text.y, text.text, text.font, opts)
	return nil
}

func (qs *QRIMGP) drawBlockQr(qr barcode.Barcode) error {
	rect := qr.Bounds()
	log.Debug().Msgf("Drawing QR with qs=%v, qr=%v", qs, qr)
	switch qs.qrLocation.location {
	case "center":
		qs.qrLocation.startX = (qs.imgp.width - rect.Dx()*qs.blockSize) / 2
		qs.qrLocation.startY = (qs.imgp.heigth - rect.Dy()*qs.blockSize) / 2
		break
	case "business_card_center":
		//if qs.qrLocation.location == "business_card_center" {
		qs.qrLocation.startX = ((qs.imgp.width / 2) - rect.Dx()*qs.blockSize) / 2
		qs.qrLocation.startY = (qs.imgp.heigth - rect.Dy()*qs.blockSize) / 2
		break
	}

	//log.Debug().Msgf("Start y=%v", qs.qrLocation.startY)

	if qr.Metadata().CodeKind == "QR Code" {

		if qs.imgp.colorPalette.IsImageDark() {
			qs.transparent = false
		}

		if qs.fillColor == "" {
			// if imgp.colorPalette.IsImageDark() {
			// 	fillColor = "white"
			// } else {
			// 	fillColor = "black"
			// }
			qs.fillColor = "black"
		}

		if !qs.transparent && qs.unfillColor == "" {
			qs.unfillColor = "white"
			if qs.fillColor == "white" {
				qs.unfillColor = "black"
			}
		}

		maxX := rect.Dx() + 1
		maxY := rect.Dy() + 1

		currY := qs.qrLocation.startY - qs.blockSize
		for y := -1; y < maxY; y++ {
			currX := qs.qrLocation.startX - qs.blockSize
			for x := -1; x < maxX; x++ {
				if x >= 0 && x < rect.Dx() && y >= 0 && y < rect.Dy() {
					if qr.At(x, y) == color.Black {
						qs.imgp.svg.Rect(currX, currY, qs.blockSize, qs.blockSize, "fill:"+qs.fillColor+";stroke:none")
					} else if !qs.transparent && qr.At(x, y) == color.White {
						qs.imgp.svg.Rect(currX, currY, qs.blockSize, qs.blockSize, "fill:"+qs.unfillColor+";stroke:none")
					}
				} else if !qs.transparent || qs.imgp.colorPalette.IsImageDark() {
					qs.imgp.svg.Rect(currX, currY, qs.blockSize, qs.blockSize, "fill:"+qs.unfillColor+";stroke:none")
				}
				currX += qs.blockSize
			}
			currY += qs.blockSize
		}
	}
	return nil
}

func (imgp *IMGParams) DrawImage() (*bytes.Buffer, error) {
	imgp.startDrawing()
	if err := imgp.generateImage(); err != nil {
		return nil, err
	}
	imgp.endDrawing()
	return imgp.OutputBuffer, nil
}

func (imgp *IMGParams) DrawQRCode(text string, x int, y int, qrLocation string, blockSize int, fillColor string, unfillColor string, transparent bool) (*bytes.Buffer, error) {
	imgp.startDrawing()
	if err := imgp.generateImage(); err != nil {
		return nil, err
	}

	qs := QRIMGP{
		imgp:      imgp,
		blockSize: blockSize,
		text:      text,
		qrLocation: QRLocation{
			startX:   x,
			startY:   y,
			location: qrLocation,
		},
		fillColor:   fillColor,
		unfillColor: unfillColor,
		transparent: transparent,
	}

	if err := qs.drawQR(); err != nil {
		return nil, err
	}

	imgp.endDrawing()
	return imgp.OutputBuffer, nil
}

func (imgp *IMGParams) DrawBusinessCard(

	contact VCardContact, x int, y int, blockSize int, fillColor string, unfillColor string, transparent bool,
	padding int, fontfamily string, fontsize int, fontunits string, fontcolor string, textOptions []string,

) (*bytes.Buffer, error) {
	imgp.startDrawing()
	if err := imgp.generateImage(); err != nil {
		return nil, err
	}

	qrText, err := contact.ToQRFormat()
	if err != nil {
		return nil, err
	}

	qs := QRIMGP{
		imgp:      imgp,
		blockSize: blockSize,
		text:      qrText,
		qrLocation: QRLocation{
			startX:   x,
			startY:   y,
			location: "business_card_center",
		},
		fillColor:   fillColor,
		unfillColor: unfillColor,
		transparent: transparent,
	}

	log.Debug().Msg("Drawing QR")
	if err := qs.drawQR(); err != nil {
		return nil, err
	}
	log.Debug().Msg("QR drawn successfully")

	log.Debug().Msg("Obtaining font style")
	textTexplate, err := FontStyleTextDrawable(imgp.svg, fontfamily, fontsize, fontunits, fontcolor, textOptions)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Font style obtained successfully textTemplate=%v", textTexplate)

	log.Debug().Msg("Calculating text coords")
	if err := qs.calculateTextCoords(textTexplate); err != nil {
		return nil, err
	}
	log.Debug().Msg("Text coords calculated successfully")

	log.Debug().Msg("Drawing VCARD")
	if err := imgp.drawVCard(contact, *textTexplate); err != nil {
		return nil, err
	}
	log.Debug().Msg("Drawing VCARD successfully")

	imgp.endDrawing()
	return imgp.OutputBuffer, nil
}

func (imgp *IMGParams) DrawLogo(companyName string, companyDesc string) (*bytes.Buffer, error) {
	imgp.startDrawing()
	if err := imgp.generateImage(); err != nil {
		return nil, err
	}

	fontfamily := ""
	fontsize := 32
	fontunits := ""
	fontcolor := "#fff"
	textOptions := make([]string, 0)

	textTexplate, err := FontStyleTextDrawable(imgp.svg, fontfamily, fontsize, fontunits, fontcolor, textOptions)
	if err != nil {
		return nil, err
	}

	imgp.drawText(*textTexplate)

	imgp.endDrawing()
	return imgp.OutputBuffer, nil
}
