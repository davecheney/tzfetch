package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"archive/tar"
	"compress/gzip"
	"flag"
	"log"
	"os"
)

var (
	verbose = flag.Bool("v", false, "verbose")
	cwd     = flag.String("C", mustCwd(), "base directory")
)

func main() {
	flag.Parse()
	for _, url := range flag.Args() {
		if err := fetchAndUnpack(url); err != nil {
			log.Fatal("failed to unpack %v: %v", url, err)
		}
	}
}

func fetchAndUnpack(url string) error {
	r, err := fetch(url)
	if err != nil {
		return err
	}
	defer r.Close()
	return unpack(r)
}

func fetch(url string) (io.ReadCloser, error) {
	debugf("fetching: %v", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetch %q return status: %d %v", url, resp.StatusCode, resp.Status)
	}
	return resp.Body, nil
}

func unpack(r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		path := filepath.Join(*cwd, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			debugf("making directory: %v", path)
			if err := os.MkdirAll(path, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			debugf("unpacking file: %v", path)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tr)
			f.Close()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot handle tar file with type: %d", hdr.Typeflag)
		}
	}
	return nil
}

func debugf(format string, args ...interface{}) {
	if *verbose {
		log.Printf(format, args...)
	}
}

func mustCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
