package main

import (
	"DictBddProject/dictionary"
	"fmt"
	"os"
)

func main() {
	dict, err := dictionary.New("./data")

	hanleErr(err)
	defer dict.Close()

	dict.AddEntry("Todo 1", "Information 1")
	dict.AddEntry("Todo 2", "Information 2")

	dict.AddEntry("Todo 3", "Information 3")
	list, err := dict.List()
	fmt.Println(list, " Check list.")

	_, err = dict.Get("Todo 1")

	hanleErr(err)

}

func hanleErr(err error) {
	if err != nil {
		fmt.Println("Dictionary error")
		os.Exit(1)
	}
}
