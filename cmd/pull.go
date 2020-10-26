// Package cmd is the command line utility for frodo
/*
Copyright © 2020 Theo Salvo <buzzsurfr>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/buzzsurfr/frodo/schema"
	"github.com/containerd/containerd/reference"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	targetRef     string
	keepOldFiles  bool
	pathTraversal bool
	output        string
	verbose       bool

	debug     bool
	configs   []string
	username  string
	password  string
	insecure  bool
	plainHTTP bool
}

var pullOpts pullOptions

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")

		pullOpts.targetRef = args[0]

		ctx := context.Background()

		// Logging (logrus)
		if pullOpts.debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		resolver := newResolver(pullOpts.username, pullOpts.password, pullOpts.insecure, pullOpts.plainHTTP, pullOpts.configs...)
		store := content.NewFileStore(pullOpts.output)
		defer store.Close()
		store.DisableOverwrite = pullOpts.keepOldFiles
		store.AllowPathTraversalOnWrite = pullOpts.pathTraversal

		desc, artifacts, err := oras.Pull(ctx, resolver, pullOpts.targetRef, store,
			oras.WithAllowedMediaTypes([]string{schema.ProtoV2.String(), schema.ProtoV3.String()}),
			oras.WithPullStatusTrack(os.Stdout),
		)
		if err != nil {
			if err == reference.ErrObjectRequired {
				fmt.Println("image reference format is invalid. Please specify <name:tag|name@digest>")
			}
			os.Exit(1)
		}
		if len(artifacts) == 0 {
			fmt.Println("Downloaded empty artifact")
		}
		fmt.Println("Pulled", pullOpts.targetRef)
		fmt.Println("Digest:", desc.Digest)

	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pullCmd.Flags().BoolVarP(&pullOpts.keepOldFiles, "keep-old-files", "k", false, "do not replace existing files when pulling, treat them as errors")
	pullCmd.Flags().BoolVarP(&pullOpts.pathTraversal, "allow-path-traversal", "T", false, "allow storing files out of the output directory")
	pullCmd.Flags().StringVarP(&pullOpts.output, "output", "o", "", "output directory")
	pullCmd.Flags().BoolVarP(&pullOpts.verbose, "verbose", "v", false, "verbose output")

	pullCmd.Flags().BoolVarP(&pullOpts.debug, "debug", "d", false, "debug mode")
	pullCmd.Flags().StringArrayVarP(&pullOpts.configs, "config", "c", nil, "auth config path")
	pullCmd.Flags().StringVarP(&pullOpts.username, "username", "u", "", "registry username")
	pullCmd.Flags().StringVarP(&pullOpts.password, "password", "p", "", "registry password")
	pullCmd.Flags().BoolVarP(&pullOpts.insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	pullCmd.Flags().BoolVarP(&pullOpts.plainHTTP, "plain-http", "", false, "use plain http and not https")
}
