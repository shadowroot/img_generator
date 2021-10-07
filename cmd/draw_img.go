package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"bitbucket.org/gekky/img-generator/img_generator"
)

func main() {
	randImgCreate()
}

func createDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0770)
		return err
	}
	return nil
}

func fileWrite(imgBuff *bytes.Buffer) error {
	dirPath := "../img"
	if err := createDirIfNotExists(dirPath); err != nil {
		return err
	}
	img_path := fmt.Sprintf("%v/%v.svg", dirPath, time.Now())
	fh, err := os.OpenFile(img_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fh.Close()
	if err != nil {
		fmt.Printf("E:%v", err)
		return err
	}
	_, err = fh.Write(imgBuff.Bytes())
	if err != nil {
		fmt.Printf("E:%v\n", err)
	}
	//fmt.Printf("Generated img with Len:%v\n", n)
	return nil
}

func randImgCreate() error {

	r := &img_generator.IMGRequest{
		W:     980,
		H:     560,
		Color: "",
	}

	startTime := time.Now()
	imgParams := img_generator.CreateImage(r.W, r.H, r.Color, true)
	imgBuff, err := imgParams.DrawImage()
	if err != nil {
		return err
	}
	fmt.Println("IMG generator took time:", time.Now().Sub(startTime))
	if err := fileWrite(imgBuff); err != nil {
		return err
	}
	return nil
}
