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
)

const (
	indexComment     = "<!-- Index generated by MDI -->"
	navComment       = "<!-- Nav generated by MDI -->"
	bottomNavComment = "<!-- Bottom nav generated by MDI -->"
)

var defaultIndexFile = "index.md"
var mdExts = []string{".md"}

type IndexOption struct {
	WorkDir    string
	IndexTitle string
	IndexFile  string
	Exclude    []string
	chains     []*index
}

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

func NewIndex(idxOpt *IndexOption) *index {
	if includeFile(idxOpt.Exclude, idxOpt.WorkDir) {
		return nil
	}

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
		if includeFile(idxOpt.Exclude, subFile) {
			continue
		}

		if f.IsDir() {
			if hasMdFile(subFile) {
				subIndexOpt := &IndexOption{
					WorkDir:    subFile,
					IndexTitle: readTitle(path.Join(subFile, defaultIndexFile)),
					IndexFile:  path.Join(subFile, defaultIndexFile),
					Exclude:    idxOpt.Exclude,
					chains:     append(idxOpt.chains, idx), // append chains in sub index option
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

	var writeIndex bool
	if file, err := os.Open(idx.file); err != nil && os.IsNotExist(err) {
		if os.IsNotExist(err) {
			writeIndex = true
		}
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			firstLine := strings.TrimSpace(scanner.Text())
			if firstLine == indexComment {
				writeIndex = true
			}
		}
	}

	if writeIndex || genOpt.Override {
		err := os.WriteFile(idx.file, []byte(fmt.Sprintf("%s\n%s\n# %s\n%s", indexComment, idx.getIndexNav(), idx.title, content)), 0644)
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

	indexNav := "> "
	// index not included, so loop to len-1
	for i := 0; i < len(idx.chains)-1; i++ {
		backpath := strings.Repeat("../", len(idx.chains)-i-1) + path.Base(idx.chains[i].file)
		indexNav += fmt.Sprintf("[%s](%s) / ", idx.chains[i].title, getLink(backpath))
	}
	indexNav += idx.title + "\n"
	return indexNav
}

func (idx *index) getEntryNavPrefix() string {
	navPrefix := "> "
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
			if len(lines) > 0 && lines[0] == navComment && strings.HasPrefix(lines[1], "> [") {
				// update nav
				lines[1] = navPrefix + readTitle(entry.file)
			} else {
				// insert nav
				lines = append([]string{navComment, navPrefix + readTitle(entry.file) + "\n"}, lines...)
			}

			if len(lines) > 5 {
				if lines[len(lines)-4] == bottomNavComment {
					lines = lines[:len(lines)-4]
				}
				if lines[len(lines)-6] == bottomNavComment {
					lines = lines[:len(lines)-6]
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
		result += fmt.Sprintf("> [上篇：%s](%s)\n", e.prev.title, getLink(path.Base(e.prev.file)))
	}
	if e.next != nil {
		if result != "" {
			result += ">\n"
		}
		result += fmt.Sprintf("> [下篇：%s](%s)\n", e.next.title, getLink(path.Base(e.next.file)))
	}
	if result != "" {
		result = bottomNavComment + "\n---\n" + result
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

	dirEntries, _ := os.ReadDir(dir)
	for _, de := range dirEntries {
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

func includeFile(paths []string, file string) bool {
	for _, p := range paths {
		s, err := filepath.Rel(p, file)
		if err == nil && !strings.HasPrefix(s, "..") {
			return true
		}
	}
	return false
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
