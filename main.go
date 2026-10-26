package main

import (
	"firstProgram/internal/color"
	"fmt"

	"github.com/joefazee/httplib"
)

func main() {
	redText := "This is red text"
	blueText := "This is blue text"
	magentaText := "This is magenta text"
	cyanText := "This is cyan text"
	boldRedText := "This is bold red text"
	BoldUnderlineBlueText := "This is bold underline blue text"

	fmt.Println(color.Text(redText, color.Red))
	fmt.Println(color.Text(blueText, color.Blue))
	fmt.Println(color.Text(magentaText, color.Magenta))
	fmt.Println(color.Text(cyanText, color.Cyan))
	fmt.Println("==============================")
	fmt.Println(color.Text(boldRedText, color.Red, color.Bold))
	fmt.Println(color.Text(magentaText, color.Magenta, color.Bold))
	fmt.Println(color.Text(BoldUnderlineBlueText, color.Blue, color.Underline))
	fmt.Println(color.Text(cyanText, color.Cyan, color.Underline))
	fmt.Println()

	rb := httplib.NewRequestBuilder("GET", "https://dummyjson.com/test")

	var res struct {
		Status string `json:"status"`
		Method string `json:"method"`
	}
	response, err := rb.Send()
	if err == nil {
		err = rb.ReadJSONResponse(response, &res)
	}

	fmt.Println(res)
}
