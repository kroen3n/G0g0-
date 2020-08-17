package main

import (
    "os"
    "log"
    "io"
    "io/ioutil"
    "archive/tar"
    "bytes"
    "context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
)


func main() {


    //create context and client object
    //documentation: https://godoc.org/github.com/docker/docker/client

    ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        log.Fatal(err, " cannot initiate client")
    }

    // negotiate version if client version 1.41 is too new. 
    cli.NegotiateAPIVersion(ctx)

    //create  tar archive

    buffer := new(bytes.Buffer)
    tar_buff := tar.NewWriter(buffer)
    defer tar_buff.Close()

    //pass the name of Dockerfile and open it

    docker_file := "Dockerfile"

    //here change the path to where your Dockerfile is located 
    //...lazy, did it statically! 

    docker_file_open, err := os.Open("/home/gog0/Dockerfile")

    if err != nil {
        log.Fatal(err, " cannot open Dockerfile")
    }

    //read Dockerfile file
    //ReadAll documentation: https://golang.org/pkg/io/ioutil/#example_ReadAll

    read_docker_file, err := ioutil.ReadAll(docker_file_open)

    if err != nil {
        log.Fatal(err, " cannot read Dockerfile")
    }


    // add files (Dockerfile) to tar archive
    //and write tar body

    tar_header := &tar.Header{
        Name: docker_file,
        Size: int64(len(read_docker_file)),
    }

    err = tar_buff.WriteHeader(tar_header)

    if err != nil {
        log.Fatal(err, " :unable to write tar header")
    }

    _, err = tar_buff.Write(read_docker_file)

    if err != nil {
        log.Fatal(err, " cannot write tar body")
    }
    
    //open archive tar,  so it can be passed to Docker image options

    docker_file_tar_read := bytes.NewReader(buffer.Bytes())

    // Docker options for building the Docker image
    //context - our archive
    // Dockerfile - our Dockerfile 
    // Tags - the name of our Docker image 
    //documentations types.ImageBuildOptions: https://godoc.org/github.com/docker/docker/api/types#ImageBuildOptions

    opt := types.ImageBuildOptions{
	    Context: docker_file_tar_read,
	    Dockerfile: docker_file,
	    Tags: []string{"hielaunewimage"},}
 
    //build image by passing the following:
    //context, the opened archive tar &&  the options

    build_image, err := cli.ImageBuild(ctx, docker_file_tar_read, opt)

    if err != nil {
        log.Fatal(err, " cannit build Docker image")
    }

    defer build_image.Body.Close()

    _, err = io.Copy(os.Stdout, build_image.Body)

    if err != nil {
        log.Fatal(err, "cannot read image build response")
    }
}
