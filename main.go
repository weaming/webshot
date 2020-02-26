package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Task struct {
	request    TaskRequest
	resultChan chan TaskResult
}

type TaskRequest struct {
	url     string
	width   int
	quality int
}
type TaskResult struct {
	name string
	data []byte
	err  error
}

var (
	taskQueue chan Task
	imgRender ImageRender

	httpPort       = flag.Int("port", 80, "port to serve as http service")
	binPath        = flag.String("bin", "/usr/local/bin/wkhtmltoimage", "wkhtmltoimage bin path")
	widthDefault   = flag.Int("w", 1200, "output image width")
	qualityDefault = flag.Int("q", 80, "output image quality, maxium is 100")
	workerCount    = flag.Int("workers", 100, "number of works")
)

func readFile(filePath string) []byte {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return dat
}

func reqOr(req *http.Request, name string, dft int) int {
	v := req.FormValue(name)
	if len(v) > 0 {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return dft

}

func processTask(id int) {
	for {
		task := <-taskQueue
		request := task.request
		url := request.url
		resultCh := task.resultChan
		log.Printf("[%d] Processing %s\n", id, url)

		reg, _ := regexp.Compile("[^A-Za-z0-9]+")
		fileName := fmt.Sprintf("%s.%s", reg.ReplaceAllString(url, "-"), "png")
		output := "" // blank to output image to stdout
		html, err := url2html(url)
		if err != nil {
			resultCh <- TaskResult{fileName, []byte{}, err}
		} else {
			bytes, err := imgRender.generateImage(html, "png", output, request.width, request.quality)
			if err != nil {
				log.Printf("[%d] Cannot take a screenshot.\n%v+\n", id, err)
				resultCh <- TaskResult{fileName, []byte{}, err}
			} else {
				resultCh <- TaskResult{fileName, bytes, nil}
			}
		}
		log.Printf("[%d] Finished processing %s\n", id, url)
	}

}

func fromUrl(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	download := len(req.FormValue("dl")) > 0
	if url == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing url"))
	} else {
		resultChan := make(chan TaskResult, 1)
		taskQueue <- Task{
			request:    TaskRequest{url, reqOr(req, "width", *widthDefault), reqOr(req, "quality", *qualityDefault)},
			resultChan: resultChan,
		}
		result := <-resultChan
		if result.err != nil {
			if strings.HasPrefix(url, "http") {
				w.WriteHeader(500)
			} else {
				w.WriteHeader(400)
			}
			w.Write([]byte(result.err.Error()))
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "image/png")
			if download {
				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", result.name))
			}
			w.Write(result.data)
		}
	}
}

func main() {
	flag.Parse()
	imgRender = ImageRender{BinaryPath: binPath}
	taskQueue = make(chan Task, *workerCount)
	for i := 0; i < *workerCount; i++ {
		go processTask(i)
	}
	http.HandleFunc("/", fromUrl)
	log.Printf("listening on %d\n", *httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
}
