// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2013 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	fileRegister       map[string]bool
	fileRegisterLock   sync.RWMutex
	readerRegister     map[string]func() io.Reader
	readerRegisterLock sync.RWMutex
)

// RegisterLocalFile adds the given file to the file allowlist,
// so that it can be used by "LOAD DATA LOCAL INFILE <filepath>".
// Alternatively you can allow the use of all local files with
// the DSN parameter 'allowAllFiles=true'
//
//	filePath := "/home/gopher/data.csv"
//	mysql.RegisterLocalFile(filePath)
//	err := db.Exec("LOAD DATA LOCAL INFILE '" + filePath + "' INTO TABLE foo")
//	if err != nil {
//	...
func RegisterLocalFile(filePath string) {
	fileRegisterLock.Lock()
	// lazy map init
	if fileRegister == nil {
		fileRegister = make(map[string]bool)
	}

	fileRegister[strings.Trim(filePath, `"`)] = true
	fileRegisterLock.Unlock()
}

// DeregisterLocalFile removes the given filepath from the allowlist.
func DeregisterLocalFile(filePath string) {
	fileRegisterLock.Lock()
	delete(fileRegister, strings.Trim(filePath, `"`))
	fileRegisterLock.Unlock()
}

// RegisterReaderHandler registers a handler function which is used
// to receive a io.Reader.
// The Reader can be used by "LOAD DATA LOCAL INFILE Reader::<name>".
// If the handler returns a io.ReadCloser Close() is called when the
// request is finished.
//
//	mysql.RegisterReaderHandler("data", func() io.Reader {
//		var csvReader io.Reader // Some Reader that returns CSV data
//		... // Open Reader here
//		return csvReader
//	})
//	err := db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE foo")
//	if err != nil {
//	...
func RegisterReaderHandler(name string, handler func() io.Reader) {
	readerRegisterLock.Lock()
	// lazy map init
	if readerRegister == nil {
		readerRegister = make(map[string]func() io.Reader)
	}

	readerRegister[name] = handler
	readerRegisterLock.Unlock()
}

// DeregisterReaderHandler removes the ReaderHandler function with
// the given name from the registry.
func DeregisterReaderHandler(name string) {
	readerRegisterLock.Lock()
	delete(readerRegister, name)
	readerRegisterLock.Unlock()
}

func deferredClose(err *error, closer io.Closer) {
	closeErr := closer.Close()
	if *err == nil {
		*err = closeErr
	}
}

func (mc *mysqlConn) handleInFileRequest(name string) (err error) {
	var rdr io.Reader
	var data []byte
	packetSize := 16 * 1024 // 16KB is small enough for disk readahead and large enough for TCP
	if mc.maxWriteSize < packetSize {
		packetSize = mc.maxWriteSize
	}

	if idx := strings.Index(name, "Reader::"); idx == 0 || (idx > 0 && name[idx-1] == '/') { // io.Reader
		// The server might return an an absolute path. See issue #355.
		name = name[idx+8:]

		readerRegisterLock.RLock()
		handler, inMap := readerRegister[name]
		readerRegisterLock.RUnlock()

		if inMap {
			rdr = handler()
			if rdr != nil {
				if cl, ok := rdr.(io.Closer); ok {
					defer deferredClose(&err, cl)
				}
			} else {
				err = fmt.Errorf("Reader '%s' is <nil>", name)
			}
		} else {
			err = fmt.Errorf("Reader '%s' is not registered", name)
		}
	} else { // File
		name = strings.Trim(name, `"`)
		fileRegisterLock.RLock()
		fr := fileRegister[name]
		fileRegisterLock.RUnlock()
		if mc.cfg.AllowAllFiles || fr {
			var file *os.File
			var fi os.FileInfo

			if file, err = os.Open(name); err == nil {
				defer deferredClose(&err, file)

				// get file size
				if fi, err = file.Stat(); err == nil {
					rdr = file
					if fileSize := int(fi.Size()); fileSize < packetSize {
						packetSize = fileSize
					}
				}
			}
		} else {
			err = fmt.Errorf("local file '%s' is not registered", name)
		}
	}

	// send content packets
	// if packetSize == 0, the Reader contains no data
	if err == nil && packetSize > 0 {
		data := make([]byte, 4+packetSize)
		var n int
		for err == nil {
			n, err = rdr.Read(data[4:])
			if n > 0 {
				if ioErr := mc.writePacket(data[:4+n]); ioErr != nil {
					return ioErr
				}
			}
		}
		if err == io.EOF {
			err = nil
		}
	}

	// send empty packet (termination)
	if data == nil {
		data = make([]byte, 4)
	}
	if ioErr := mc.writePacket(data[:4]); ioErr != nil {
		return ioErr
	}

	// read OK packet
	if err == nil {
		return mc.readResultOK()
	}

	mc.readPacket()
	return err
}
