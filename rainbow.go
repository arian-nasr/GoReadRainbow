package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
}

func decompressBlob(blobData []byte) []byte {
	b := bytes.NewReader(blobData)
	z, err := zlib.NewReader(b)
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()
	data, err := ioutil.ReadAll(z)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

type BlobXML struct {
	BlobID int `xml:"blobid,attr"`
	Size int `xml:"size,attr"`
	Compression string `xml:"compression,attr"`
}

func getRainbowBlobData(data []byte, blobid int) []byte {
	blobSearchString := fmt.Sprintf(`<BLOB blobid="%d"`, blobid)
	// startPos will either give index of blob or will be -1 if blob not in file
	startPos := bytes.Index(data, []byte(blobSearchString))
	endPos := bytes.Index(data[startPos:], []byte(">"))+startPos
	xmlString := string(data[startPos:endPos+1])+"</BLOB>"
	var blobXML BlobXML
	err := xml.Unmarshal([]byte(xmlString), &blobXML)
	if err != nil {
		log.Fatal(err)
	}
	blobData := data[endPos+2:endPos+2+blobXML.Size]
	if blobXML.Compression == "qt" {
		data = decompressBlob(blobData[4:])
	}
	return data
}

// Get rainbow header
func getRainbowHeader(file *os.File) string {
	header := ""
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		linetext := scanner.Text()
		if linetext != "<!-- END XML -->" {
			header += linetext
		} else {
			break
		}
	}
	return header
}

func readRainbow(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	blob := getRainbowBlobData(data, 10)
	fmt.Println(blob)


	// TEST - Print the header
	//fmt.Println(getRainbowHeader(file))

	//blobSearchString := getRainbowBlobData("test", 1)
	//fmt.Println(blobSearchString)


	file.Close()
}