package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type tftpStuff struct {
	username string
	password string
	role     string
}

func (t tftpStuff) readHandler(filename string, rf io.ReaderFrom) error {
	log.Infof("Got request for %s", filename)
	someString := fmt.Sprintf("username %s password 0 %s role %s\n", t.username, t.password, t.role)
	myReader := strings.NewReader(someString)
	_, err := rf.ReadFrom(myReader)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func (t tftpStuff) writeHandler(filenameRaw string, wt io.WriterTo) error {
	log.Infof("Got request for %s", filenameRaw)
	_, filename := filepath.Split(filenameRaw)
	if filename == "" {
		log.Infof("empty filename after cleanup: %s", filenameRaw)
		return fmt.Errorf("empty filename after cleanup: %s", filenameRaw)
	}
	filename = fmt.Sprintf("%s_%s", filename, randomString(5))

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	log.Infof("%s: %d bytes received", filename, n)
	return nil
}
