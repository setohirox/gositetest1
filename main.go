package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

type PageData struct {
	Hello   string
	Title   string
	URLpath string
}

func shouldOpenBrowser() bool {
	// air実行かどうかはコマンドパスで検出（tmpフォルダなどを含む）
	exe := strings.ToLower(os.Args[0])
	noBrowser := os.Getenv("NO_BROWSER")
	log.Println("実行パス:", exe)
	log.Println("NO_BROWSER:", noBrowser)
	return !strings.Contains(exe, "tmp") && noBrowser == ""
}

func handlerTOP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	data := PageData{
		Hello:   "setohirox",
		Title:   "Hello setohirox",
		URLpath: "hello",
	}
	tmpl.Execute(w, data)
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	data := PageData{
		Hello:   "Again",
		Title:   "Hello Again",
		URLpath: "",
	}
	tmpl.Execute(w, data)
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		return // Linuxなど
	}

	exec.Command(cmd, args...).Start()
}

func main() {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	println("サーバーが起動しました。")
	println("Ctrl + Cで終了")
	log.SetOutput(f)
	log.Println("サーバーが起動しました")

	http.HandleFunc("/", handlerTOP)
	http.HandleFunc("/hello", handlerHello)

	if shouldOpenBrowser() {
		openBrowser("http://localhost:8080")
	}

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("サーバーが終了しました")
}
