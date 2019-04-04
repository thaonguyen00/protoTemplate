package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var flags struct {
	FileIn  string
	FileOut string
	Entity  string
}

func main() {
	app := cli.NewApp()
	app.Name = "protoTemplate"
	app.Usage = "add generic service code to a proto file"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "input, i",
			Usage:       "inputfile",
			Destination: &flags.FileIn,
		},
		cli.StringFlag{
			Name:        "output,o",
			Usage:       "output file",
			Value:       "./structNames.txt",
			Destination: &flags.FileOut,
		},
		cli.StringFlag{
			Name:        "Entity, e",
			Usage:       "Entity name",
			Value:       "Template",
			Destination: &flags.Entity,
		},
	}
	app.Action = launch
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

var oldCode = `syntax = "proto3";
package proto;`

var newCode = `syntax = "proto3";
package proto;

import "google/api/annotations.proto";

option java_multiple_files = true;
service Template {
  rpc GetTemplate (TemplateRequest) returns (Template) {
    option (google.api.http) = {
             get:"/api/v1/template/{templateguid}"
         };
  }
  rpc GetTemplates (TemplatesRequest) returns (Templates) {
    option (google.api.http).get = "/api/v1/templates";
  }

  rpc HealthCheck (HealthRequest) returns (HealthReply) {
    option (google.api.http).get = "/api/v1/health";
  }
}
message HealthRequest{}

message HealthReply{
  string status = 1;
}
message TemplateRequest {
  string templateguid = 1;
  string attributes = 2;

}
message TemplatesRequest {
  string search = 1;
  string attributes = 2;
  int64 offset = 3;
  int64 limit = 4;
  string sortKey = 5;
  bool descendingOrder = 6;

}
message Templates {
  repeated Template templates =1;
}
`

func launch(_ *cli.Context) error {
	file, _ := filepath.Abs(flags.FileIn)
	fmt.Println("worked:", file)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "Failed to read file.")
	}

	content := string(b)
	content = strings.Replace(content, oldCode, newCode, -1)
	content = strings.Replace(content, "Template", flags.Entity, -1)
	content = strings.Replace(content, "template", strings.ToLower(flags.Entity), -1)

	var fOut *os.File
	if flags.FileOut == "" {
		fOut, _ = os.Create(flags.FileIn)
	} else {
		fOut, _ = os.Create(flags.FileOut)
	}
	defer fOut.Close()

	fOut.WriteString(content)
	return nil
}
