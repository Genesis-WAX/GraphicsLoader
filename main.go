package main

import (
	"fmt"
	"github.com/saracen/go7z"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
)

func main() {

	fmt.Println("[WARNING] ONLY RUN THIS ON WINDOWS SERVERS, TO UNINSTALL THIS LOADER'S DEPENDENCIES, GO HERE:")
	fmt.Println("https://github.com/pal1000/mesa-dist-win/releases")
	fmt.Println()
	fmt.Println("Enter anything to continue...")
	var input string
	fmt.Scanln(&input)

	fmt.Println("Making Loader Directory")
	err := os.MkdirAll("loader", os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloading Release From Github")
	err = DownloadFile("./loader/loader.7z", "https://github.com/pal1000/mesa-dist-win/releases/download/21.3.3/mesa3d-21.3.3-release-msvc.7z")
	if err != nil {
		panic(err)
	}

	UnzipLoader()

	fmt.Println("Downloading Loader Bash")
	err = DownloadFile("./loader/runme.cmd", "https://www.genesiswax.io/runme.cmd")
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(path.Join(wd, "/loader/runme"))
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("...Finished, press any key to continue...")
}

func UnzipLoader() {
	fmt.Println("Unzipping Release From Github (this will take a while)")
	sz, err := go7z.OpenReader("./loader/loader.7z")
	if err != nil {
		panic(err)
	}
	defer sz.Close()

	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
		}

		// If empty stream (no contents) and isn't specifically an empty file...
		// then it's a directory.
		if hdr.IsEmptyStream && !hdr.IsEmptyFile {
			if err := os.MkdirAll(path.Join("loader", hdr.Name), os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		// Create file
		f, err := os.Create(path.Join("loader", hdr.Name))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := io.Copy(f, sz); err != nil {
			panic(err)
		}
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}