package xsdvalidate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
)

func TestMemoryParseXsd(t *testing.T) {

	Init()

	defer Cleanup()

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < 1000000; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			//start := time.Now()
			//handler, err := NewXsdHandlerUrl("examples/test1_fail.xsd", ParserDefault)
			handler, err := NewXsdHandlerUrl("examples/test1_pass.xsd", ParserDefault)
			if err != nil {
				fmt.Println(err)
			}
			handler.Free()
			//elapsed := time.Since(start)
			//log.Printf("time %s", elapsed)
			<-guard
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestMemoryParseXml(t *testing.T) {

	Init()

	defer Cleanup()

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	//xmlfile := "examples/test1_fail1.xml"
	xmlfile := "examples/test1_pass.xml"

	fxml, err := os.Open(xmlfile)
	if err != nil {
		log.Printf("failed to open file: %s", err)
		return
	}
	defer fxml.Close()

	inXml, err := ioutil.ReadAll(fxml)
	if err != nil {
		log.Printf("failed to read file: %s", err)
		return
	}

	for i := 0; i < 1000000; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(inXml []byte) {
			//start := time.Now()
			xmlhandler, err := NewXmlHandlerMem(inXml, ParserDefault)
			if err != nil {
				fmt.Println(err)
			}
			xmlhandler.Free()
			//elapsed := time.Since(start)
			//log.Printf("time %s", elapsed)
			<-guard
			wg.Done()
		}(inXml)
	}
	wg.Wait()
}
func TestMemoryValidate(t *testing.T) {

	Init()

	defer Cleanup()

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	//var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	//xmlfile := "examples/test1_fail2.xml"
	xmlfile := "examples/test1_pass.xml"

	fxml, err := os.Open(xmlfile)
	if err != nil {
		log.Printf("failed to open file: %s", err)
		return
	}
	defer fxml.Close()

	inXml, err := ioutil.ReadAll(fxml)
	if err != nil {
		log.Printf("failed to read file: %s", err)
		return
	}

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_pass.xsd", ParserDefault)
	if err != nil {
		panic(err)
	}

	defer xsdhandler.Free()

	for i := 0; i < 1000000; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(inXml []byte) {
			//start := time.Now()
			xmlhandler, err := NewXmlHandlerMem(inXml, ParserDefault)
			if err != nil {
				panic(err)
			}
			err = xsdhandler.Validate(xmlhandler, ParserDefault)
			if err != nil {
				//log.Print(err)
			}
			xmlhandler.Free()
			//elapsed := time.Since(start)
			//log.Printf("time %s", elapsed)
			<-guard
			wg.Done()
		}(inXml)
	}
	wg.Wait()
}