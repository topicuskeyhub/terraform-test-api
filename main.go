package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"slices"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type Message struct {
	Args   []string
	Opts   []string
	Config string
	Vars   map[string]string
}

var tf *tfexec.Terraform

var workingDir string

func setupTerraform() *tfexec.Terraform {
	var err error
	workingDir, err = os.MkdirTemp("/tmp", "terraapiwork")
	if err != nil {
		log.Fatalf("error creating workdir: %s", err)
	}
	installer := &releases.LatestVersion{
		Product: product.Terraform,
		Constraints: version.MustConstraints(version.NewConstraint(">= 1.11.3")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}
	log.Printf("Terraform runs in %s", execPath)
	log.Printf("Working dir is %s", workingDir)

	ret, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}
	return ret
}

func rebuild(w http.ResponseWriter, r *http.Request) {
	openapi, err := os.Create("/work/sdk-go/openapi.json")
	if err != nil {
		log.Printf("cannot open openapi file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := openapi.Close(); err != nil {
			log.Fatalf("cannot close file: %s", err)
		}
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = openapi.Write(data)
	if err != nil {
		log.Printf("cannot write openapi file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("/rebuild.sh")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("build failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		log.Printf("build failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	in := bufio.NewScanner(stdout)

	for in.Scan() {
		log.Print(in.Text())
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("build failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeConfig(r *http.Request) Message {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("cannot read body: %s", err)
	}

	var msg Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Fatalf("cannot unmarshal body: %s", err)
	}

	mainTf, err := os.Create(workingDir + "/main.tf")
	if err != nil {
		log.Fatalf("cannot create main.tf: %s", err)
	}
	defer mainTf.Close()
	_, err = mainTf.WriteString(msg.Config)
	if err != nil {
		log.Fatalf("cannot write main.tf: %s", err)
	}

	varsTf, err := os.Create(workingDir + "/test.auto.tfvars")
	if err != nil {
		log.Fatalf("cannot create test.auto.tfvars: %s", err)
	}
	defer varsTf.Close()
	for n, v := range msg.Vars {
		_, err = varsTf.WriteString(n + " = \"" + v + "\"\n")
		if err != nil {
			log.Fatalf("cannot test.auto.tfvars: %s", err)
		}
	}
	return msg
}

func tfApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not supported " + r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msg := writeConfig(r)
	err := tf.ApplyJSON(r.Context(), w, tfexec.RefreshOnly(slices.Contains(msg.Opts, "-refresh-only")))
	tf.SetStdout(nil)
	if err != nil {
		log.Println(err)
	}
}

func tfImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not supported " + r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	msg := writeConfig(r)
	err := tf.Import(r.Context(), msg.Args[0], msg.Args[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

func tfOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method not supported " + r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	output, err := tf.Output(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not supported " + r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dir, _ := os.ReadDir(workingDir)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{workingDir, d.Name()}...))
	}
}

func main() {
	tf = setupTerraform()

	mux := http.NewServeMux()
	mux.HandleFunc("/rebuild", rebuild)
	mux.HandleFunc("/apply", tfApply)
	mux.HandleFunc("/import", tfImport)
	mux.HandleFunc("/output", tfOutput)
	mux.HandleFunc("/cleanup", cleanup)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
