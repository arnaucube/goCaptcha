package main

import (
	"io/ioutil"
)

var dataset map[string][]string
var categDataset []string

func readDataset(path string) {
	dataset = make(map[string][]string)

	folders, _ := ioutil.ReadDir(path)
	for _, folder := range folders {
		categDataset = append(categDataset, folder.Name())
		folderFiles, _ := ioutil.ReadDir(path + "/" + folder.Name())
		for _, file := range folderFiles {
			dataset[folder.Name()] = append(dataset[folder.Name()], file.Name())
		}
	}
}
