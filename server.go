package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed index.html
var index []byte

func main() {
	cliCmd := append(
		strings.Split(os.Getenv("EXEC"), " "),
		"-x",
		"--audio-format",
		"mp3",
		"--audio-quality",
		"0")

	bind := os.Getenv("BIND")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Write(index)
			return
		}

		cli := append(cliCmd, r.FormValue("link"))
		cmd := exec.Command(cli[0], cli[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()

		var file string
		_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if matched, err := filepath.Match("*.mp3", filepath.Base(path)); err != nil {
				return err
			} else if matched {
				file = path
			}
			return nil
		})

		f, _ := os.Open(file)
		defer func() {
			f.Close()
			os.Remove(file)
		}()

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, filepath.Base(file)))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
	})

	http.ListenAndServe(bind, nil)
}
