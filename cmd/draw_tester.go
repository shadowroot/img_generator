// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/labstack/echo"
// 	"github.com/rs/zerolog/log"

// 	"bitbucket.org/gekky/img-generator/img_generator"
// )

// func main() {
// 	e := echo.New()
// 	e.GET("/img", imgHandler)
// 	e.GET("/qr", qrHandler)
// 	e.GET("/vcard", vCardHandler)
// 	e.POST("/img", imgHandler)
// 	e.POST("/qr", qrHandler)
// 	e.POST("/vcard", vCardHandler)
// 	e.Logger.Fatal(e.Start(":1323"))
// }

// func createDirIfNotExists(path string) error {
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		err = os.Mkdir(path, 0770)
// 		return err
// 	}
// 	return nil
// }

// func fileWrite(imgBuff *bytes.Buffer) error {
// 	dirPath := "../img"
// 	if err := createDirIfNotExists(dirPath); err != nil {
// 		return err
// 	}
// 	img_path := fmt.Sprintf("%v/%v.svg", dirPath, time.Now())
// 	fh, err := os.OpenFile(img_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	defer fh.Close()
// 	if err != nil {
// 		fmt.Printf("E:%v", err)
// 		return err
// 	}
// 	_, err = fh.Write(imgBuff.Bytes())
// 	if err != nil {
// 		fmt.Printf("E:%v\n", err)
// 	}
// 	//fmt.Printf("Generated img with Len:%v\n", n)
// 	return nil
// }

// func imgHandler(c echo.Context) error {

// 	//color := c.QueryParam("color")
// 	// w := c.QueryParam("w")
// 	// h := c.QueryParam("h")
// 	r := &img_generator.IMGRequest{
// 		W:     980,
// 		H:     560,
// 		Color: "",
// 	}

// 	startTime := time.Now()
// 	if err := c.Bind(r); err != nil {
// 		log.Error().Msgf("e=%v", err)
// 		return err
// 	}
// 	imgParams := img_generator.CreateImage(r.W, r.H, r.Color, true)
// 	imgBuff, err := imgParams.DrawImage()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("IMG generator took time:", time.Now().Sub(startTime))
// 	if err := fileWrite(imgBuff); err != nil {
// 		return err
// 	}
// 	return c.Blob(http.StatusOK, "image/svg+xml", imgBuff.Bytes())
// }

// func vCardHandler(c echo.Context) error {
// 	startTime := time.Now()
// 	vCard := &img_generator.VCardRequest{
// 		IMG: img_generator.IMGRequest{
// 			W:     980,
// 			H:     560,
// 			Color: "",
// 		},
// 		QR: img_generator.QRRequestParam{
// 			X:           5,
// 			Y:           5,
// 			Location:    "business_card_center",
// 			BlockSize:   10,
// 			FillColor:   "",
// 			UnfillColor: "",
// 			Transparent: true,
// 			Opacity:     true,
// 			Text:        "",
// 		},
// 		VCARD: img_generator.VCardContactReq{
// 			FN:    nil,
// 			N:     nil,
// 			TITLE: nil,
// 			ORG:   nil,
// 			TEL:   nil,
// 			EMAIL: nil,
// 			URL:   nil,
// 			KEY:   nil,
// 		},
// 		Padding:    5,
// 		FontFamily: "Over the Rainbow",
// 		FontSize:   32,
// 		FontUnits:  "pt",
// 		FontColor:  "#fff",
// 		Options:    nil,
// 	}

// 	if err := c.Bind(vCard); err != nil {
// 		return err
// 	}

// 	log.Debug().Msgf("[vCard request]: %v", vCard)
// 	log.Debug().Msgf("[QueryParams]: %v", c.QueryParams())

// 	imgParams := img_generator.CreateImage(vCard.IMG.W, vCard.IMG.H, vCard.IMG.Color, vCard.QR.Opacity)
// 	vCardContact := img_generator.NewVCardContact(vCard.VCARD.FN, vCard.VCARD.N, vCard.VCARD.TITLE, vCard.VCARD.ORG, vCard.VCARD.TEL, vCard.VCARD.EMAIL, vCard.VCARD.URL, vCard.VCARD.KEY)

// 	log.Debug().Msgf("vCard=%v", vCard)

// 	imgBuff, err := imgParams.DrawBusinessCard(
// 		*vCardContact,
// 		vCard.QR.X,
// 		vCard.QR.Y,
// 		vCard.QR.BlockSize,
// 		vCard.QR.FillColor,
// 		vCard.QR.UnfillColor,
// 		vCard.QR.Transparent,
// 		vCard.Padding,
// 		vCard.FontFamily,
// 		vCard.FontSize,
// 		vCard.FontUnits,
// 		vCard.FontColor,
// 		vCard.Options,
// 	)
// 	if err != nil {
// 		log.Error().Err(err)
// 		return err
// 	}
// 	fmt.Println("IMG generator took time:", time.Now().Sub(startTime))
// 	if err := fileWrite(imgBuff); err != nil {
// 		log.Error().Err(err)
// 		return err
// 	}

// 	return c.Blob(http.StatusOK, "image/svg+xml", imgBuff.Bytes())
// }

// func qrHandler(c echo.Context) error {
// 	startTime := time.Now()
// 	qrR := &img_generator.QRRequestParam{
// 		IMG: img_generator.IMGRequest{
// 			W:     980,
// 			H:     560,
// 			Color: "",
// 		},
// 		X:           5,
// 		Y:           5,
// 		Location:    "center",
// 		BlockSize:   10,
// 		FillColor:   "",
// 		UnfillColor: "",
// 		Transparent: true,
// 		Opacity:     true,
// 		Text:        "",
// 	}

// 	if err := c.Bind(qrR); err != nil {
// 		return err
// 	}

// 	log.Debug().Msgf("[QR request]: %v", qrR)
// 	log.Debug().Msgf("[QueryParams]: %v", c.QueryParams())

// 	imgParams := img_generator.CreateImage(qrR.IMG.W, qrR.IMG.H, qrR.IMG.Color, qrR.Opacity)

// 	imgBuff, err := imgParams.DrawQRCode(
// 		qrR.Text,
// 		qrR.X,
// 		qrR.Y,
// 		qrR.Location,
// 		qrR.BlockSize,
// 		qrR.FillColor,
// 		qrR.UnfillColor,
// 		qrR.Transparent,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("IMG generator took time:", time.Now().Sub(startTime))
// 	if err := fileWrite(imgBuff); err != nil {
// 		return err
// 	}

// 	return c.Blob(http.StatusOK, "image/svg+xml", imgBuff.Bytes())
// }
