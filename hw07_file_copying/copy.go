package main

import (
	"errors"
	"io"
	"log"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const chunkSize = int64(1024)

type srcParams struct {
	limit     int64
	offset    int64
	chunkSize int64
}

func validate(src *os.File, offset int64, limit int64) (srcParams, error) {
	params := srcParams{
		limit:     limit,
		offset:    offset,
		chunkSize: chunkSize,
	}
	fileInfo, err := src.Stat()
	if err != nil {
		return params, err
	}
	size := fileInfo.Size()
	if offset > size {
		return params, ErrOffsetExceedsFileSize
	}
	if !fileInfo.Mode().IsRegular() {
		return params, ErrUnsupportedFile
	}
	size -= offset
	if limit == 0 || limit > size {
		params.limit = size
	}
	if chunkSize > params.limit {
		params.chunkSize = params.limit
	}
	return params, err
}

func openFiles(fromPath, toPath string) (*os.File, *os.File, error) {
	src, err := os.OpenFile(fromPath, os.O_RDONLY, 0o744)
	if err != nil {
		return nil, nil, err
	}
	dst, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o744)
	if err != nil {
		return nil, nil, err
	}
	return src, dst, nil
}

func execute(src, dst *os.File, params srcParams) error {
	bar := pb.StartNew(int(math.Ceil(float64(params.limit) / float64(params.chunkSize))))
	defer bar.Finish()
	var count int64
	for {
		if count >= params.limit {
			break
		}
		w, err := io.CopyN(dst, src, params.chunkSize)
		count += w
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		bar.Increment()
		if errors.Is(err, io.EOF) {
			break
		}
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, dst, err := openFiles(fromPath, toPath)
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			log.Println(err)
		}
	}(src)
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			log.Println(err)
		}
	}(dst)
	if err != nil {
		return err
	}
	params, err := validate(src, offset, limit)
	if err != nil {
		return err
	}
	_, err = src.Seek(params.offset, io.SeekStart)
	if err != nil {
		return err
	}
	return execute(src, dst, params)
}
