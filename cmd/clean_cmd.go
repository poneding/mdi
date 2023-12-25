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

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean",
	Long:  `Clean`,
	Run: func(cmd *cobra.Command, args []string) {
		mdi.Clean(indexOpt.WorkDir, indexOpt.IndexFile)
	},
}

func init() {
	cleanCmd.Flags().StringVarP(&indexOpt.WorkDir, "workdir", "d", ".", "Specify the directory to clean markdown index.")
	cleanCmd.Flags().StringVarP(&indexOpt.IndexFile, "index-file", "f", "./zz_gneratered_mdi.md", "Specify the markdown index file, default is `zz_gneratered_mdi.md`.")

	rootCmd.AddCommand(cleanCmd)
}
