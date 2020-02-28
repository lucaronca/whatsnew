package filehandler

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestGetFileContents(t *testing.T) {
	request, err := ioutil.ReadFile("./test/file_request.json")

	contents, err := getFileContents(string(request))

	if err != nil {
		t.Errorf("getFileContents returned an err: %v", err)
	}

	valid := isValidXML(contents)

	if valid == false {
		t.Errorf("getFileContents returned a not valid xml")
	}
}

func isValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) != nil
}
