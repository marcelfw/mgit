// Copyright (c) 2014 Marcel Wouters
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS
// OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
// OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package config

import (
	"github.com/marcelfw/mgit/repository"
	"reflect"
	"testing"
)

func TestHardcodedParseCommandLine(t *testing.T) {
	filters := make([]repository.FilterDefinition, 0)

	//func ParseCommandline(filterDefs []repository.FilterDefinition) (command string, args []string, repositoryFilter repository.RepositoryFilter, ok bool) {
	_, _, _, ok := ParseCommandline(make([]string, 0), filters)
	if ok {
		t.Error("Empty command-line should not parse succesfully.")
	}

	command, _, repFilter, ok := ParseCommandline([]string{"list"}, filters)
	if !ok {
		t.Errorf("Expected ok to be true, but got '%v'", ok)
	}
	if command != "list" {
		t.Errorf("Expected command to be 'list', but got '%v'", command)
	}
	st := reflect.ValueOf(repFilter)
	if value := st.FieldByName("rootDirectory"); value.String() != "." {
		t.Errorf("Expended rootDirectory to be '.', got '%v'", value)
	}

	command, _, repFilter, ok = ParseCommandline([]string{"-root", "/", "status"}, filters)
	if !ok {
		t.Errorf("Expected ok to be true, but got %v", ok)
	}
	if command != "status" {
		t.Errorf("Expected command to be 'status', but got '%v'", command)
	}
	st = reflect.ValueOf(repFilter)
	if value := st.FieldByName("rootDirectory"); value.String() != "/" {
		t.Errorf("Expected rootDirectory to be '/', got '%v'", value)
	}

	command, _, repFilter, ok = ParseCommandline([]string{"-depth", "10", "path"}, filters)
	if !ok {
		t.Errorf("Expected ok to be true, but got %v", ok)
	}
	if command != "path" {
		t.Errorf("Expected command to be 'status', but got '%v'", command)
	}
	st = reflect.ValueOf(repFilter)
	if value := st.FieldByName("depth"); value.Int() != 10 {
		t.Errorf("Expected depth to be '10', got '%v'", value)
	}

}
