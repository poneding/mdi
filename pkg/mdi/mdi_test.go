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

import "testing"

func TestIncludeFile(t *testing.T) {
	paths := []string{"**/test", "test-01.md", "hello/*"}
	testdata := []struct {
		file     string
		expected bool
	}{
		{"hello/index.md", true},
		{"hello1/index.md", false},
		{"dir1/test/index.md", true},
		{"dir1/sub1/test/index.md", true},
		{"dir2/test/index.md", true},
		{"test-01.md", true},
		{"test-02.md", false},
	}

	for _, d := range testdata {
		actual := includeFileV2(paths, d.file)
		if actual != d.expected {
			t.Errorf("includeFile(%q, %q) = %v, expected %v", d.file, paths, actual, d.expected)
		}
	}
	t.Log("TestIncludeFile passed")
}

func TestIncludeFileV2(t *testing.T) {
	paths := []string{""}
	testdata := []struct {
		file     string
		expected bool
	}{
		{"hello/index.md", true},
		{"hello/test/index.md", true},
		{"test-01.md", true},
		{"test-02.md", true},
	}

	for _, d := range testdata {
		actual := includeFileV2(paths, d.file)
		if actual != d.expected {
			t.Errorf("includeFile(%q, %q) = %v, expected %v", d.file, paths, actual, d.expected)
		}
	}
	t.Log("TestIncludeFile passed")
}
