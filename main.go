package main

import (
	"fmt"
	"html/template"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

var counter int
var i photo
var path string
var tpl *template.Template

type photo struct {
	Name string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", index)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":" + port, nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	// delete files if necessary
	checkCounter()

	if req.Method == http.MethodPost {
		mf, fh, err := req.FormFile("nf")
		if err != nil {
			fmt.Println(err)
		}
		defer mf.Close()

		// create new file
		path = filepath.Join("public", "pics", fh.Filename)
		nf, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer nf.Close()

		// copy
		mf.Seek(0, 0)
		io.Copy(nf, mf)

		transformImage(path)
		// create struct for view
		i = photo{
			Name: path,
		}
	}

	increment()
	tpl.ExecuteTemplate(w, "index.gohtml", i)
}

func checkCounter() {
	if counter > 0 {
		os.Remove(path)
		i = photo{}
	}
}

func increment() {
	counter = counter + 1
}

func transformImage(f string) {
	// open uploaded file
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 128
	m := resize.Resize(128, 128, img, resize.Lanczos3)

	out, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
