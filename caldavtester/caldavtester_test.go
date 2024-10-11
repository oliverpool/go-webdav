package caldavtester

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestServer(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:8008")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	s := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: create backend
		}),
	}
	go s.Serve(ln)

	f, err := os.Create("test.log")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := f.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
	cmd := exec.Command("nix", "run", "--offline", "--quiet", "../", "--",
		"testcaldav.py",
		// relative to the ccs-caldavtester folder
		"-x", "scripts/tests/CardDAV",
		"-s", "../serverinfo.xml",
		"--print-details-onfail",
		"--stop",
	)
	cmd.Dir = "ccs-caldavtester"
	cmd.Stdout = io.MultiWriter(f, os.Stdout)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// t.Log(stderr.String())
		t.Fatal(err)
	}

}
