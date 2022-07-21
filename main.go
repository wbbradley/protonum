package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var protoDir = flag.String("dir", ".", "directory containing protobuf models")

func main() {
	flag.Parse()
	_ = filepath.Walk(*protoDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			fmt.Println(path)
			bytes, _ := ioutil.ReadFile(path)
			renumberedAll := renumberAllMessages(string(bytes))
			ioutil.WriteFile(path, []byte(renumberedAll), 0644)
		}

		return nil
	})
}

func renumberAllMessages(text string) string {
	messages := regexp.MustCompile(`(?sm)^message.*?{$.*?^}`).FindAllString(text, -1)

	for _, msg := range messages {
		renumbered := renumberFields(msg)
		text = strings.ReplaceAll(text, msg, renumbered)
	}

	return text
}

func renumberFields(text string) string {
	loc := regexp.MustCompile(`(?sm){$.*^}`).FindStringIndex(text)
	body := text[loc[0]:loc[1]]
	i := 1
	replaced := regexp.MustCompile(`(?sm)\d+;$`).ReplaceAllStringFunc(body, func(s string) string {
		s = fmt.Sprintf("%d;", i)
		i++
		return s
	})
	return strings.ReplaceAll(text, body, replaced)
}
