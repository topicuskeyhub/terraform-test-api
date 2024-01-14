package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.6.6")),
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

func main() {
	tf = setupTerraform()

	mux := http.NewServeMux()
	mux.HandleFunc("/apply", tfApply)
	mux.HandleFunc("/import", tfImport)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
