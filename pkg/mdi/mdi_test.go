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
	paths := []string{"**/target", "**/hello", "test-01.md"}
	testdata := []struct {
		file     string
		expected bool
	}{
		{"rust-lang-book/rlb01_hello-cargo/target/index.md", true},
		{"hello/index.md", true},
		{"dir1/hello/index.md", true},
		{"dir1/sub1/hello/index.md", true},
		{"dir1/hello2/index.md", false},
		{"test-01.md", true},
		{"test-02.md", false},
	}

	for _, d := range testdata {
		actual := includeFile(paths, d.file)
		if actual != d.expected {
			t.Errorf("includeFile(%q, %q) = %v, expected %v", d.file, paths, actual, d.expected)
		}
	}
	t.Log("TestIncludeFile passed")
}
