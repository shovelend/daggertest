package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// build the container using a dockerfile
	f := client.Directory().WithFile(".", client.Host().Workdir().File("Dockerfile"))
	container := client.Container().Build(f)

	container.Export(ctx, "test.tar")

	container.Publish(ctx, "fh1-harbor01.dun.fh/findmypast/attila-test")

	pulled := client.Container().From("fh1-harbor01.dun.fh/findmypast/attila-test")
	result := pulled.Exec(dagger.ContainerExecOpts{
		Args: []string{"node", "--version"},
	}).Stdout()

	fmt.Println(result.Contents(ctx))

	if err != nil {
		return err
	}
	fmt.Println("build finished")
	return nil
}
