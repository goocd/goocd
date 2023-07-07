//go:build ignore
// +build ignore

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {

	keepF := flag.Bool("keep", false, "Keep temporary directory with svd files instead of deleting it upon exit")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Printf("Expected exactly one argument, the path to an .atpack file")
	}

	inPath := args[0]

	// .atpack files are just zip files
	f, err := os.Open(inPath)
	if err != nil {
		log.Fatalf("Error opening file %q: %v", inPath, err)
	}
	defer f.Close()
	st, err := f.Stat()
	r, err := zip.NewReader(f, st.Size())
	if err != nil {
		log.Fatalf("Error creating zip reader: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "go-from-atpack")
	if err != nil {
		log.Fatalf("Error creating temp dir: %v", err)
	}
	log.Printf("Writing svd files to temp directory: %s", tmpDir)
	if !strings.Contains(tmpDir, "go-from-atpack") { // just in case...
		panic(fmt.Errorf("tmpDir looks wrong! %q", tmpDir))
	}
	if !*keepF {
		defer os.RemoveAll(tmpDir)
	}

	// loop over everything, find all .svd files, write them to tmpDir
	for _, zf := range r.File {
		ext := strings.ToLower(path.Ext(zf.Name))
		if ext != ".svd" {
			continue
		}
		_, fn := path.Split(zf.Name)
		zff, err := zf.Open()
		if err != nil {
			log.Fatalf("Error while opening zip file for entry %q: %v", zf.Name, err)
		}
		b, err := ioutil.ReadAll(zff)
		if err != nil {
			log.Fatalf("Error while reading zip file for entry %q: %v", zf.Name, err)
		}
		zff.Close()
		// write them to temp dir
		wpath := filepath.Join(tmpDir, fn)
		err = os.WriteFile(wpath, b, 0644)
		if err != nil {
			log.Fatalf("Error while writing svd to temp file for entry %q: %v", zf.Name, err)
		}
		log.Printf("Wrote: %s", wpath)
	}

	// now fire off gen-device-svd.go so it can do the rest of the
	// work to convert from svd to .go files
	genDeviceSVDCmdLine := []string{
		"go", "run", "gen-device-svd.go", "-source=" + inPath, tmpDir, "../sam",
	}
	cmd := exec.Command(genDeviceSVDCmdLine[0], genDeviceSVDCmdLine[1:]...)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error from os.Getwd: %v", err)
	}
	cmd.Dir = filepath.Join(wd, "../../tools")
	log.Printf("gen-device-svd command (wd=%s): %v", cmd.Dir, genDeviceSVDCmdLine)
	b, err := cmd.CombinedOutput()
	log.Printf("gen-device-svd output: %s", b)
	if err != nil {
		log.Fatalf("gen-device-svd execution error: %v", err)
	}
}
