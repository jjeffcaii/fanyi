package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	lpad  = " "
	token = "C78E00EE3519FD0E34684C1318982F7D"
)

var input string

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage:")
		fmt.Println("  fanyi hello")
		os.Exit(1)
	}
	input = strings.ToLower(args[0])
}

func main() {
	body, err := seek(input)
	if err != nil {
		panic(err)
	}
	var word Word
	err = xml.NewDecoder(bytes.NewBuffer(body)).Decode(&word)
	if err != nil {
		panic(err)
	}
	word.Print()
}

func seek(word string) (body []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf("http://dict-co.iciba.com/api/dictionary.php?w=%s&key=%s&type=xml", word, token))
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

type Word struct {
	Key         string   `xml:"key"`
	Ps          []string `xml:"ps"`
	Pos         []string `xml:"pos"`
	Acceptation []string `xml:"acceptation"`
	Sent        []Sent   `xml:"sent"`
}

var (
	cKey   = color.New(color.FgHiWhite)
	cMut   = color.New(color.FgHiBlack)
	cPs    = color.New(color.FgMagenta)
	cAcc   = color.New(color.FgGreen)
	cTrans = color.New(color.FgCyan)
	cHigh  = color.New(color.FgYellow)
)

func (s Word) Print() {
	fmt.Printf("\n%s", lpad)
	cKey.Printf("%s", s.Key)
	fmt.Print("  ")
	if len(s.Ps) > 1 {
		cPs.Printf("英[ %s ]  美[ %s ]", s.Ps[0], s.Ps[1])
	}
	fmt.Print("\n\n")

	if len(s.Ps) > 0 {
		for i := 0; i < len(s.Pos); i++ {
			fmt.Print(lpad)
			cMut.Print("- ")
			cAcc.Printf("%s %s", s.Pos[i], strings.TrimSpace(s.Acceptation[i]))
			fmt.Print("\n")
		}
		fmt.Print("\n")
	}

	for k, v := range s.Sent {
		fmt.Print(lpad)
		sp := strings.Split(strings.TrimSpace(v.Orig), input)
		cMut.Printf("%d. ", k+1)
		for _, vv := range sp[:len(sp)-1] {
			cMut.Print(vv)
			cHigh.Print(input)
		}
		cMut.Print(sp[len(sp)-1])
		fmt.Printf("\n%s   ", lpad)
		cTrans.Printf("%s", strings.TrimSpace(v.Trans))
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

type Sent struct {
	Orig  string `xml:"orig"`
	Trans string `xml:"trans"`
}
