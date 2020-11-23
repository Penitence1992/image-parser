package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	goyaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"k8s.io/client-go/util/jsonpath"
	"os"
	"path/filepath"
)

var gitCommit = ""
var buildStamp = ""

func main() {

	fmt.Printf("Git Commit : %s\n", gitCommit)
	fmt.Printf("Build Stamp : %s\n", buildStamp)

	if len(os.Args) == 1 {
		log.Fatal("请输入文件路径")
	}

	filePath := os.Args[1]

	if filePath == "-" {
		log.Infof("read content in stdin")
		info, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("获取输入流失败:%v", err)
		}
		if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
			log.Fatalf("The command is intended to work with pipes.\n\"Usage: fortune | iparse\n")
		}
		contentParse(os.Stdin)
	} else {
		log.Infof("file path: %v", filePath)
		info, err := os.Stat(filePath)

		if err != nil {
			log.Fatalf("获取路径信息失败:%v", err)
		}

		if info.IsDir() {
			filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
				// 不再遍历
				if !info.IsDir() {
					fileParse(path)
				}
				return nil
			})
		} else {
			fileParse(filePath)
		}
	}
}

func fileParse(filePath string) {
	context, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal("read file error : ", err)
	}

	r := bytes.NewReader(context)

	contentParse(r)
}

func contentParse(reader io.Reader) {

	dec := goyaml.NewDecoder(reader)

	var doc interface{}

	for dec.Decode(&doc) == nil {
		docBytes, err := goyaml.Marshal(doc)
		printError("marshal document error", err)
		jsonDoc, err := yaml.YAMLToJSON(docBytes)
		var jsonInterface interface{}

		printError("convert json to interface fail : ", json.Unmarshal(jsonDoc, &jsonInterface))

		parse(jsonInterface)
	}
}

func parse(jsonDoc interface{}) {
	p := jsonpath.New("image")

	if err := p.Parse("{..image}"); err != nil {
		log.Fatal("json path parse fail :", err)
	}

	ress, err := p.FindResults(jsonDoc)

	if err != nil {
		return
	}

	for _, res := range ress {
		for _, filed := range res {
			_, _ = os.Stdout.WriteString(fmt.Sprintf("%s\n", filed.Interface().(string)))
		}
	}
}

func printError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
