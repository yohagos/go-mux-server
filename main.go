package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

const index = `
<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>List of Services</title>
		</head>
		<body>
			<h1>List of all Services - Version {{ .Version }}</h1>
			{{ range .WebList }}
			<ul style="list-style-type: none;">
				<li><a href="http://localhost:8888/{{.}}">{{ . }}</a></li>
			</ul>
			{{ end }}
			<br>
			<h1>List of all Docs </h1>
			{{ range .DocsList }}
			<ul style="list-style-type: none;">
				<li><a href="http://localhost:8888/{{.}}">{{ . }}</a></li>
			</ul>
			{{ end }}
		</body>
	</html>
`

var (
	versionDirectory  string
	staticPath        string
	filePathSeparator string

	resourceList  []string
	directoryList []string
	docsPathList  []string

	linksList []string
	docsList  []string

	templates *template.Template
)

func main() {
	versionDirectory = addVersionDirectory()
	filePathSeparator = determineOS()

	resourceList, directoryList = findResources(versionDirectory)
	log.Println("List of services")
	for _, v := range resourceList {
		strs := strings.Split(v, filePathSeparator)
		linksList = append(linksList, strs[len(strs)-1])
		log.Println(" - " + strs[len(strs)-1])
	}

	dir := findDir(resourceList, filePathSeparator)

	docsList = getDocsList(docsPathList, filePathSeparator)

	loadTemplates(dir + filePathSeparator + "*.html")

	router := mux.NewRouter()
	router.HandleFunc("/", indexGETHandler)
	router.HandleFunc("/{findService}", serviceHandler)

	fsJavascript := http.FileServer(http.Dir(staticPath + "/js/service"))
	router.PathPrefix("../js/service").Handler(http.StripPrefix("../js/service", fsJavascript))
	router.PathPrefix("/js/service").Handler(http.StripPrefix("/js/service", fsJavascript))
	router.PathPrefix("js/service").Handler(http.StripPrefix("/js/service", fsJavascript))

	fsConfig := http.FileServer(http.Dir(staticPath + "/app.config"))
	router.PathPrefix("../app.config").Handler(http.StripPrefix("../app.config", fsConfig))
	router.PathPrefix("/app.config").Handler(http.StripPrefix("/app.config", fsConfig))
	router.PathPrefix("app.config").Handler(http.StripPrefix("/app.config", fsConfig))

	http.Handle("/", router)
	fmt.Println()
	log.Println("Server running under: localhost:8888")

	if err := http.ListenAndServe(":8888", router); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}

}

func determineOS() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return "\\"
	default:
		return "/"
	}
}

func isError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func addVersionDirectory() string {
	if len(os.Args[1:]) > 0 {
		for _, v := range os.Args[1:] {
			return v
		}
	}
	log.Fatal("Version was not mentioned. Please enter which version you want to test")
	os.Exit(1)
	return ""
}

func findDir(list []string, seperator string) string {
	var strList []string

	for i := 0; i <= 1; i++ {
		strList = strings.Split(list[i], seperator)
	}

	strList = strList[:len(strList)-1]
	dir := strings.Join(strList[:], seperator)
	return dir
}

func getDocsList(list []string, seperator string) []string {
	var strList []string
	var str string

	for _, v := range list {
		s := strings.Split(v, seperator)
		str = s[len(s)-1]
		strList = append(strList, str)
	}

	for _, v := range strList {
		fmt.Println(v)
	}
	return strList
}

func findResources(versionDirectory string) ([]string, []string) {
	var pathList []string
	var dirList []string
	staticPath, err := os.Getwd()
	isError(err)
	staticPath = staticPath + filePathSeparator + versionDirectory

	err = filepath.Walk(staticPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		matched, _ := regexp.MatchString(".html", path)
		if matched {
			pathList = append(pathList, path)
		}
		return nil
	})
	isError(err)

	err = filepath.Walk(staticPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		matched, _ := regexp.MatchString(".pdf", path)
		if matched {
			docsPathList = append(docsPathList, path)
		}
		return nil
	})
	isError(err)

	err = filepath.Walk(staticPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == true {
			dirList = append(dirList, path)
		}
		return nil
	})

	return pathList, dirList
}

func loadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func indexGETHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content Type", "text/html")
	tmpl, err := templates.New("test").Parse(index)
	if err == nil {
		tmpl.Execute(w, struct {
			Version  string
			WebList  []string
			DocsList []string
		}{
			Version:  versionDirectory,
			WebList:  linksList,
			DocsList: docsList,
		})
	}
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	service := r.URL.RequestURI()[1:]
	path := checkPath(service)
	if strings.Contains(service, ".pdf") {
		var str string
		for _, v := range docsPathList {
			if strings.Contains(v, service) {
				str = v
			}
		}

		pdfFile, err := os.Open(str)
		if err != nil {
			log.Println("Error : ", err)
		}
		defer pdfFile.Close()
		w.Header().Set("Content Type", "application/pdf")
		if _, err := io.Copy(w, pdfFile); err != nil {
			log.Println("Error opening PDF File: ", service)
		}
	}

	fmt.Println(path)
	executeTemplate(w, service, nil)
}

func checkPath(service string) string {
	for _, v := range resourceList {
		matched, _ := regexp.MatchString(service, v)
		if matched {
			return v
		}
	}
	return ""
}

func executeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}
