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

package mdi

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

var defaultIndexFile = "zz_generated_mdi.md"
var mdExts = []string{".md"}

var defaultIndexOption = &IndexOption{
	WorkDir:    ".",
	IndexTitle: "Index",
	IndexFile:  defaultIndexFile,
}

var defaultGenerationOption = &GenerationOption{
	Override:  false,
	Recursive: true,
	Nav:       false,
}

type index struct {
	workDir  string
	file     string
	title    string
	content  string
	chains   []*index
	children []*index
	entries  []*entry
}

type IndexOption struct {
	WorkDir          string
	IndexTitle       string
	IndexFile        string
	InheritGitIgnore bool
	chains           []*index
	rootExcludes     *[]string
	subExcludes      []string
}

type GenerationOption struct {
	Override  bool
	Recursive bool
	Nav       bool
	Verbose   bool
}

type entry struct {
	title string
	file  string
	prev  *entry
	next  *entry
}

func (idxOpt *IndexOption) RootExcludes() []string {
	if idxOpt.rootExcludes == nil {
		idxOpt.rootExcludes = &[]string{}
		// .mdiignore
		*idxOpt.rootExcludes = append(*idxOpt.rootExcludes, getIgnoreEntry(path.Join(idxOpt.WorkDir, ".mdiignore"))...)

		// .gitignore
		if idxOpt.InheritGitIgnore {
			*idxOpt.rootExcludes = append(*idxOpt.rootExcludes, getIgnoreEntry(path.Join(idxOpt.WorkDir, ".gitignore"))...)
		}
	}

	// fmt.Printf("(*idxOpt.rootExcludes): %v\n", (*idxOpt.rootExcludes))

	return *idxOpt.rootExcludes
}

func getSubExcludes(subDir string) []string {
	var result []string
	// sub .mdiignore
	result = append(result, getIgnoreEntry(path.Join(subDir, ".mdiignore"))...)
	// sub .gitignore
	result = append(result, getIgnoreEntry(path.Join(subDir, ".gitignore"))...)
	return result
}

func NewIndex(idxOpt *IndexOption) *index {
	files, err := os.ReadDir(idxOpt.WorkDir)
	if err != nil {
		panic(err)
	}

	// validate opt
	if idxOpt == nil {
		idxOpt = defaultIndexOption
	}

	if idxOpt.WorkDir == "" {
		idxOpt.WorkDir = defaultIndexOption.WorkDir
	}
	if fi, err := os.Stat(idxOpt.WorkDir); os.IsNotExist(err) && !fi.IsDir() {
		panic(fmt.Sprintf("invalid work dir: %s", idxOpt.WorkDir))
	}
	if idxOpt.IndexTitle == "" {
		idxOpt.IndexTitle = defaultIndexOption.IndexTitle
	}
	if idxOpt.IndexFile == "" {
		idxOpt.IndexFile = path.Join(idxOpt.WorkDir, defaultIndexFile)
	}

	idx := &index{
		workDir:  idxOpt.WorkDir,
		file:     idxOpt.IndexFile,
		title:    idxOpt.IndexTitle,
		children: make([]*index, 0),
		entries:  make([]*entry, 0),
	}
	// set self as chain tail
	idx.chains = append(idxOpt.chains, idx)

	for _, f := range files {
		subFile := path.Join(idxOpt.WorkDir, f.Name())
		if matchFile(idxOpt.RootExcludes(), subFile) {
			continue
		}

		if f.IsDir() {
			if hasMdFile(subFile) {
				subIndexOpt := &IndexOption{
					WorkDir:      subFile,
					IndexTitle:   readTitle(path.Join(subFile, defaultIndexFile)),
					IndexFile:    path.Join(subFile, defaultIndexFile),
					rootExcludes: idxOpt.rootExcludes,
					subExcludes:  getSubExcludes(subFile),
					chains:       append(idxOpt.chains, idx), // append chains in sub index option
				}
				subIdx := NewIndex(subIndexOpt)
				if subIdx != nil {
					idx.children = append(idx.children, subIdx)
				}
			}
		} else {
			if slices.Contains(mdExts, path.Ext(f.Name())) && f.Name() != path.Base(idx.file) {
				idx.entries = append(idx.entries, &entry{
					title: readTitle(subFile),
					file:  subFile,
				})
			}
		}
	}

	for i := 0; i < len(idx.entries); i++ {
		if i > 0 {
			idx.entries[i].prev = idx.entries[i-1]
		}
		if i < len(idx.entries)-1 {
			idx.entries[i].next = idx.entries[i+1]
		}
	}

	return idx
}

func (idx *index) Generate(genOpt *GenerationOption) {
	if idx == nil {
		return
	}

	for _, subIdx := range idx.children {
		if genOpt.Recursive {
			subIdx.Generate(genOpt)
		}
	}
	content := parseContent(idx, idx.workDir, "", 0)

	if genOpt.Override {
		err := os.WriteFile(idx.file, []byte(fmt.Sprintf("%s# %s\n%s", idx.getIndexNav(), idx.title, content)), 0644)
		if err != nil {
			fmt.Printf("ERROR: failed to write index file: %s", err)
		} else {
			if genOpt.Verbose {
				fmt.Printf("OK: generated index file: %s\n", idx.file)
			}
		}
	} else {
		if genOpt.Verbose {
			fmt.Printf("SKIP: index file conflict: %s, use --override=true to override it\n", idx.file)
		}
	}

	if genOpt.Nav {
		idx.decorateEntry()
	}
}

func (idx *index) getIndexNav() string {
	if len(idx.chains) <= 1 {
		return ""
	}

	var indexNav string
	// index not included, so loop to len-1
	for i := 0; i < len(idx.chains)-1; i++ {
		backpath := strings.Repeat("../", len(idx.chains)-i-1) + path.Base(idx.chains[i].file)
		indexNav += fmt.Sprintf("[%s](%s) / ", idx.chains[i].title, getLink(backpath))
	}
	indexNav += idx.title + "\n\n"
	return indexNav
}

func (idx *index) getEntryNavPrefix() string {
	var navPrefix string
	for i := 0; i < len(idx.chains); i++ {
		backpath := strings.Repeat("../", len(idx.chains)-i-1) + path.Base(idx.chains[i].file)
		navPrefix += fmt.Sprintf("[%s](%s) / ", idx.chains[i].title, getLink(backpath))
	}
	return navPrefix
}

func (idx *index) decorateEntry() {
	for _, entry := range idx.entries {
		if s, _ := filepath.Rel(idx.file, entry.file); s == "." {
			continue
		}
		b, err := os.ReadFile(entry.file)
		if err == nil {
			if len(b) == 0 {
				continue
			}

			lines := strings.Split(string(b), "\n")

			navPrefix := idx.getEntryNavPrefix()
			if strings.HasPrefix(lines[0], "[") {
				// update nav
				lines[0] = navPrefix + readTitle(entry.file)
			} else {
				// insert nav
				lines = append([]string{navPrefix + readTitle(entry.file) + "\n"}, lines...)
			}

			if len(lines) > 4 {
				if lines[len(lines)-3] == "---" {
					lines = lines[:len(lines)-3]
				}
				if lines[len(lines)-5] == "---" {
					lines = lines[:len(lines)-5]
				}
			}

			// update bottom nav
			bottomNav := entry.getBottomNav()
			if bottomNav != "" {
				lines = append(lines, bottomNav)
			}

			updated := strings.Join(lines, "\n")
			f, err := os.Create(entry.file)
			if err == nil {
				defer f.Close()
				f.WriteString(updated)
			}
		}
	}
}

func (e *entry) getBottomNav() string {
	var result string
	if e.prev != nil {
		result += fmt.Sprintf("[» %s](%s)\n", e.prev.title, getLink(path.Base(e.prev.file)))
	}
	if e.next != nil {
		if result != "" {
			result += "\n"
		}
		result += fmt.Sprintf("[« %s](%s)\n", e.next.title, getLink(path.Base(e.next.file)))
	}
	if result != "" {
		result = "---\n" + result
	}
	return result
}

func parseContent(idx *index, workdir, content string, depth int) string {
	for _, subIdx := range idx.children {
		relPath, _ := filepath.Rel(workdir, subIdx.file)
		if depth == 0 {
			content += fmt.Sprintf("\n## [%s](%s)\n", subIdx.title, getLink(relPath))
		} else {
			content += fmt.Sprintf("\n%s- [%s](%s)", strings.Repeat("  ", depth-1), subIdx.title, getLink(relPath))
		}

		content += parseContent(subIdx, workdir, "", depth+1)
	}

	for _, entry := range idx.entries {
		relPath, _ := filepath.Rel(workdir, entry.file)
		if depth == 0 {
			content += fmt.Sprintf("\n[%s](%s)\n", entry.title, getLink(relPath))
		} else {
			content += fmt.Sprintf("\n%s- [%s](%s)", strings.Repeat("  ", depth-1), entry.title, getLink(relPath))
		}
	}

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return content
}

func getLink(file string) string {
	return strings.ReplaceAll(file, " ", "%20")
}

var dirHasMdFileMap = make(map[string]bool)

func hasMdFile(dir string) bool {
	if v, ok := dirHasMdFileMap[dir]; ok {
		return v
	}

	subExcludes := getSubExcludes(dir)

	dirEntries, _ := os.ReadDir(dir)
	for _, de := range dirEntries {
		if matchFile(subExcludes, de.Name()) {
			continue
		}
		if !de.IsDir() {
			if slices.Contains(mdExts, path.Ext(de.Name())) && de.Name() != defaultIndexFile {
				dirHasMdFileMap[path.Join(dir, de.Name())] = true
				dirHasMdFileMap[dir] = true
				return true
			}
		} else {
			if hasMdFile(path.Join(dir, de.Name())) {
				dirHasMdFileMap[dir] = true
				return true
			}
		}
	}

	return dirHasMdFileMap[dir]
}

func matchFile(paths []string, file string) bool {
	patterns := []gitignore.Pattern{}

	for _, p := range paths {
		patterns = append(patterns, gitignore.ParsePattern(p, nil))
	}
	m := gitignore.NewMatcher(patterns)
	return m.Match(strings.Split(file, "/"), true)
}

var fileTitleMap = make(map[string]string)

func readTitle(file string) string {
	if v, ok := fileTitleMap[file]; ok {
		return v
	}

	if !slices.Contains(mdExts, path.Ext(file)) {
		fileTitleMap[file] = path.Base(file)
		return fileTitleMap[file]
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fileTitleMap[file] = path.Base(path.Dir(file))
		return fileTitleMap[file]
	}

	f, err := os.Open(file)
	if err != nil {
		fileTitleMap[file] = path.Base(file)
		return fileTitleMap[file]
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		cut, ok := strings.CutPrefix(s.Text(), "# ")
		if ok && len(cut) > 0 {
			fileTitleMap[file] = strings.TrimSpace(cut)
			return fileTitleMap[file]
		}
	}
	fileTitleMap[file] = path.Base(file)
	return fileTitleMap[file]
}

func getIgnoreEntry(ignoreFile string) []string {
	var result []string
	mdiignore, err := os.Stat(ignoreFile)
	if err == nil && !mdiignore.IsDir() {
		f, err := os.Open(ignoreFile)
		if err == nil {
			defer f.Close()
			s := bufio.NewScanner(f)
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if line != "" && !strings.HasPrefix(line, "#") {
					result = append(result, s.Text())
				}
			}
		}
	}
	return result
}

func Clean(workDir, rootMDIFile string) {
	if workDir == "" {
		workDir = "."
	}
	files, err := os.ReadDir(workDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			Clean(path.Join(workDir, f.Name()), defaultIndexFile)
		} else {
			if f.Name() == rootMDIFile {
				os.Remove(path.Join(workDir, f.Name()))
			}

			if slices.Contains(mdExts, path.Ext(f.Name())) {
				b, err := os.ReadFile(path.Join(workDir, f.Name()))
				if err == nil {
					if len(b) == 0 {
						continue
					}

					lines := strings.Split(string(b), "\n")

					if len(lines) > 4 {
						if lines[len(lines)-3] == "---" {
							lines = lines[:len(lines)-3]
						}
						if lines[len(lines)-5] == "---" {
							lines = lines[:len(lines)-5]
						}
					}
					if len(lines) > 1 {
						if strings.HasPrefix(lines[0], "[") && lines[1] == "" {
							lines = lines[2:]
						}
					}

					updated := strings.Join(lines, "\n")
					updatedFile, err := os.Create(path.Join(workDir, f.Name()))
					if err == nil {
						defer updatedFile.Close()
						updatedFile.WriteString(updated)
					}
				}
			}
		}

	}
}
