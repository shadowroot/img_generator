package img_generator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	rnd "math/rand"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	colorMap = map[string][]byte{
		"red":     {0xff, 0x00, 0x00},
		"yellow":  {0xff, 0xff, 0x00},
		"green":   {0x0, 0xff, 0x0},
		"purple":  {0xf0, 0x0, 0xf0},
		"magenta": {0xff, 0x0, 0xff},
		"white":   {0xff, 0xff, 0xff},
		"black":   {0x00, 0x00, 0x00},
		"blue":    {0x00, 0x00, 0xff},
	}
	colorMapChanged              = true
	colorMapKeys        []string = nil
	coefColorCloseRange          = float64(0.3) // close derivation
)

func GetColors() []string {
	if colorMapChanged {
		colorMapKeys = make([]string, len(colorMap))
		for key := range colorMap {
			colorMapKeys = append(colorMapKeys, key)
		}
		colorMapChanged = false
	}
	return colorMapKeys
}

type ColorGenerator struct {
	primaryColor string
	dark         bool
	options      string
	colorScheme  []byte
	color        []byte
	colorCoefs   []float64
	chosenScheme bool
	percentage   float64
}

func NewColorGenerator(primaryColor string, dark bool, options string) ColorGenerator {
	log.Debug().Msg("Creating new Color Generator")
	return ColorGenerator{
		primaryColor: primaryColor,
		dark:         dark,
		options:      options,
		colorScheme:  make([]byte, 3),
		color:        make([]byte, 3),
		colorCoefs:   make([]float64, 3),
		chosenScheme: false,
	}
}

func ParseTextInputString(text_input string) ColorGenerator {
	words := strings.Fields(text_input)
	dark := false
	primaryColor := ""
	options := ""
	for _, word := range words {
		word = strings.ToLower(word)
		if word[0] == '#' {
			wLen := len(word)
			if wLen == 4 || wLen == 7 {
				return NewColorGenerator(word, false, "")
			} else {
				log.Error().Msg("Either bad formatted color inserted or multiple colors, only one color could be inserted.")
			}
		}
		if word == "dark" {
			dark = true
		} else {
			_, ok := colorMap[word]
			if ok {
				if primaryColor != "" {
					fmt.Println("Warning: Multiple colors defined, using ", word)
				}
				primaryColor = word
			} else {
				options += word + ";"
			}
		}
	}
	return NewColorGenerator(primaryColor, dark, options)
}

func (cg ColorGenerator) isDark(color []byte) bool {
	dark := true
	for _, val := range color {
		if val > 0x90 {
			dark = false
		}
	}
	return dark
}

func (cg ColorGenerator) parsePrimaryColor() (*[]byte, error) {
	if len(cg.primaryColor) <= 0 {
		return nil, nil
	}
	if len(cg.primaryColor) > 0 && cg.primaryColor[0] == '#' {
		result, err := hex.DecodeString(cg.primaryColor[1:])
		return &result, err
	}
	val, ok := colorMap[cg.primaryColor]
	if ok {
		return &val, nil
	}
	return nil, nil
}

func (cg *ColorGenerator) calculateCG() {
	for i, _ := range cg.colorScheme {
		cg.colorCoefs[i] = float64(cg.colorScheme[i]) / float64(255)
		//Fix for too dark colors - mistake
		// if cg.colorScheme[i] == 0 {
		// 	cg.colorCoefs[i] = 0xa
		// }
	}
}

func (cg *ColorGenerator) generateRandomColor() ([]byte, error) {
	//log.Debug().Msgf("CG: %v", cg)

	if !cg.chosenScheme {
		log.Debug().Msg("Creating new color")

		c, err := cg.parsePrimaryColor()
		if err != nil {
			return nil, err
		}
		if c != nil {
			cg.colorScheme = *c
			cg.dark = cg.isDark(cg.colorScheme)
			log.Debug().Msgf("Using color scheme for color: '%v'", cg.colorScheme)
		} else {
			n, err := rand.Read(cg.colorScheme)
			if n < 3 || err != nil {
				return nil, err
			}
			cg.dark = cg.isDark(cg.colorScheme)
		}
		cg.calculateCG()
		cg.chosenScheme = true
	}

	cntr := 50
	cg.percentage = rnd.Float64()
	for {

		for i := 0; i < len(cg.colorCoefs); i++ {
			cg.color[i] = byte(math.Floor(cg.colorCoefs[i] * cg.percentage * float64(255)))
		}
		if cntr <= 0 || cg.isDark(cg.color) == cg.dark {
			break
		}
		cg.percentage = rnd.Float64()
		cntr--
	}
	log.Debug().Msgf("[Generated color] ColorScheme: %v, Color: %v, Coefs: %v, percentage: %v, chosenScheme: %v", cg.colorScheme, cg.color, cg.colorCoefs, cg.percentage, cg.chosenScheme)
	return cg.color, nil
}

//
func (cg *ColorGenerator) GenerateTematicColorString() (string, error) {
	clr, err := cg.generateRandomColor()
	if err != nil {
		return "", err
	}
	return colorToColorString(clr), nil
}

//
func (cg ColorGenerator) IsImageDark() bool {
	return cg.dark
}

//Opacity of colors
func (cg ColorGenerator) GetOpacity() float64 {
	return cg.percentage
}

func colorToColorString(color []byte) string {
	return fmt.Sprintf("#%02x%02x%02x", color[0], color[1], color[2])
}

func grayScale(color []byte) byte {
	return byte(0.21*float32(color[0]) + 0.72*float32(color[1]) + 0.07*float32(color[2]))
}
