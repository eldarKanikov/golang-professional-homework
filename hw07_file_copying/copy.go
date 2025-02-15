package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeValues        = errors.New("limit and offset must be equal 0 (default) or positive")
)

const bufferSize = 1024

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 || offset < 0 {
		fmt.Printf("error passing incorrect input : %v\n", ErrNegativeValues)
		return ErrNegativeValues
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		fmt.Printf("error getting fileStat : %v\n", err)
		return err
	}
	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	sFile, err := os.Open(fromPath)
	if err != nil {
		fmt.Printf("error Opening file : %v\n", err)
		return err
	}
	defer sFile.Close()
	reader := bufio.NewReader(sFile)

	dFile, _ := os.Create(toPath)
	defer dFile.Close()
	writer := bufio.NewWriter(dFile)

	if limit == 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	bar := pb.Full.Start64(limit)
	barWriter := bar.NewProxyWriter(writer)

	if offset > 0 {
		_, err := sFile.Seek(offset, io.SeekStart)
		if err != nil {
			fmt.Printf("error seeking file : %v\n", err)
			return err
		}
	}

	buffer := make([]byte, bufferSize)

	var totalCopied int64

	for totalCopied < limit {
		remain := limit - totalCopied

		if remain < bufferSize {
			buffer = make([]byte, remain)
		}

		read, err := reader.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error reading bytes file : %v\n", err)
			return err
		}
		_, err = barWriter.Write(buffer[:read])
		if err != nil {
			fmt.Printf("error writing bytes file : %v\n", err)
			return err
		}
		err = writer.Flush()
		if err != nil {
			fmt.Printf("error flushing the writer : %v\n", err)
			return err
		}
		totalCopied += int64(read)
	}
	bar.Finish()
	return nil
}
