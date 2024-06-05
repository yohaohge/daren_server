package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	fBuf, err := ioutil.ReadFile("./data/9.9.jpg")
	if err != nil {
		panic(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(fBuf))
}
