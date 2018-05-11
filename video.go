package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"html"
)

func saveHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}

		uploadedFile, err := os.Create("./save/" + part.FileName())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			uploadedFile.Close()
			redirectToErrorPage(w,r)
			return
		}

		_, err = io.Copy(uploadedFile, part)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			uploadedFile.Close()
			redirectToErrorPage(w,r)
			return
		}
	}
	http.Redirect(w,r,"/uploadComplate",http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var templatefile = template.Must(template.ParseFiles("./html/upload.html"))
	templatefile.Execute(w, "upload.html")
}

func errorPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"%s","<p>Internal Server Error</p>")
}

func redirectToErrorPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w,r,"/errorPage",http.StatusFound)
}

func mainSample(w http.ResponseWriter, r *http.Request){
		dir := "./save"
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		var paths []string
		for _, file := range files {
			paths = append(paths, file.Name())
		}

		for _, path := range paths {
			fmt.Println(path)
		}

	value_post := paths

	output := `
<html>
    <head>
        <title>bouno</title>
    </head>
    <body>
		<form method="post" action="/upload">
 		   <input id="upload" type="submit" value="アップロード画面へ">
		</form>

		<a href="http://localhost:8080/play" target="_blank">
			<ul>
            	<li>` + html.EscapeString(value_post[0]) + `</li>
        	</ul></a>
    </body>
</html>
`
	fmt.Fprintf(w, "%s", output)
}

func playHandler(w http.ResponseWriter, r *http.Request){
	var templatefile = template.Must(template.ParseFiles("./html/play.html"))
	templatefile.Execute(w, "play.html")
}

func uploadCompHandler(w http.ResponseWriter, r *http.Request){
	var templatefile = template.Must(template.ParseFiles("./html/uploadComplate.html"))
	templatefile.Execute(w, "uploadComplate.html")
}

func main() {
	//ハンドラの登録
	http.HandleFunc("/errorPage", errorPageHandler)
	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/uploadComplate", uploadCompHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/main", mainSample)

	//サーバーの開始
	http.ListenAndServe("localhost:8080", nil)
}
