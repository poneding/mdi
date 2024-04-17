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
	"github.com/poneding/mdi/pkg/mdi"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate markdown index",
	Long:  `Generate markdown index`,
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
	genCmd.Flags().StringVarP(&indexOpt.WorkDir, "workdir", "d", ".", "Specify the directory to generate markdown index.")
	genCmd.Flags().StringVarP(&indexOpt.IndexTitle, "index-title", "t", "", "Specify the title of markdown index, default is title of markdown index file or current directory name.")
	genCmd.Flags().StringVarP(&indexOpt.IndexFile, "index-file", "f", "./zz_generated_mdi.md", "Specify the markdown index file, default is `zz_generated_mdi.md`.")
	genCmd.Flags().BoolVar(&indexOpt.InheritGitIgnore, "inherit-gitignore", true, "Use `.gitignore` file as ignore file, default is `true`.")
	genCmd.Flags().BoolVar(&genOpt.Override, "override", false, "Override markdown existing index file, default is `false`.")
	genCmd.Flags().BoolVarP(&genOpt.Recursive, "recursive", "r", false, "Recursively generate markdown index in subdirectories, default is `false`.")
	genCmd.Flags().BoolVar(&genOpt.Nav, "nav", false, "Generate navigation in markdown file, default is `false`.")
	genCmd.Flags().BoolVarP(&genOpt.Verbose, "verbose", "v", false, "Show verbose log, default is `false`.")

	rootCmd.AddCommand(genCmd)
}
