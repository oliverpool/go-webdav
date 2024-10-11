package caldavtester

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCardDAV(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:8008")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	s := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.BasicAuth())

			// TODO: create backend
		}),
	}
	go s.Serve(ln)

	var debug bool
	for _, a := range os.Args {
		if strings.HasPrefix(a, "-test.run=") {
			debug = true
			break
		}
	}

	logFile := "carddav.log"
	if debug {
		logFile = "debug.log"
	}
	f, err := os.Create(logFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := f.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	args := []string{
		"run", "--offline", "--quiet", "../", "--",
		"testcaldav.py",
		// relative to the ccs-caldavtester folder
		"-x", "scripts/tests/CardDAV",
		"-s", "../serverinfo.xml",
	}
	if debug {
		args = append(args,
			"--print-details-onfail",
			"--stop",
		)
	}

	cmd := exec.Command("nix", args...)
	cmd.Dir = "ccs-caldavtester"
	cmd.Stdout = io.MultiWriter(f, os.Stdout)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// t.Log(stderr.String())
		t.Fatal(err)
	}

}
