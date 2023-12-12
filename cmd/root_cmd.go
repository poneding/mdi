/*
Copyright 2023 Pone Ding <poneding@gmail.com>.

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
	"os"

	"github.com/poneding/mdi/pkg/mdi"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdi",
	Short: "mdi is a command line tool used to recursively generate markdown indexes in directories.",
	Long:  `mdi is a command line tool used to recursively generate markdown indexes in directories. version: ` + version,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var indexOpt = &mdi.IndexOption{}

var genOpt = &mdi.GenerationOption{}

func run() {
	mdi.NewIndex(indexOpt).Generate(genOpt)
}

func init() {
	rootCmd.Flags().StringVarP(&indexOpt.WorkDir, "workdir", "d", ".", "Specify the directory to generate markdown index.")
	rootCmd.Flags().StringVarP(&indexOpt.IndexTitle, "index-title", "t", "", "Specify the title of markdown index, default is title of markdown index file or current directory name.")
	rootCmd.Flags().StringVarP(&indexOpt.IndexFile, "index-file", "f", "./index.md", "Specify the markdown index file, default is `index.md`.")
	rootCmd.Flags().BoolVar(&indexOpt.InheritGitIgnore, "inherit-gitignore", true, "Use `.gitignore` file as ignore file, default is `true`.")
	rootCmd.Flags().BoolVar(&genOpt.Override, "override", false, "Override markdown existing index file, default is `false`.")
	rootCmd.Flags().BoolVarP(&genOpt.Recursive, "recursive", "r", false, "Recursively generate markdown index in subdirectories, default is `false`.")
	rootCmd.Flags().BoolVar(&genOpt.Nav, "nav", false, "Generate navigation in markdown file, default is `false`.")
	rootCmd.Flags().BoolVarP(&genOpt.Verbose, "verbose", "v", false, "Show verbose log, default is `false`.")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
