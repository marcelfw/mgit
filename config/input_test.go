// Copyright (c) 2014 Marcel Wouters

package config

import (
	"github.com/marcelfw/mgit/repository"
	"reflect"
	"testing"
)

func TestHardcodedParseCommandLine(t *testing.T) {
	filters := make([]repository.FilterDefinition, 0)

	_, _, _, _, ok := ParseCommandline(make([]string, 0), filters)
	if ok {
		t.Error("Empty command-line should not parse succesfully.")
	}

	command, _, _, repFilter, ok := ParseCommandline([]string{"list"}, filters)
	if !ok {
		t.Errorf("Expected ok to be true, but got '%v'", ok)
	}
	if command != "list" {
		t.Errorf("Expected command to be 'list', but got '%v'", command)
	}
	st := reflect.ValueOf(repFilter)
	if value := st.FieldByName("rootDirectory"); value.String() != "." {
		t.Errorf("Expected rootDirectory to be '.', got '%v'", value)
	}

	command, _, _, repFilter, ok = ParseCommandline([]string{"-root", "/", "status"}, filters)
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

	command, _, _, repFilter, ok = ParseCommandline([]string{"-depth", "10", "path"}, filters)
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

	command, _, _, repFilter, ok = ParseCommandline([]string{"-root", ".", "status"}, filters)
	if !ok {
		t.Errorf("Expected ok to be true, but got %v", ok)
	}
	if command != "status" {
		t.Errorf("Expected command to be 'status', but got '%v'", command)
	}
	st = reflect.ValueOf(repFilter)
	if value := st.FieldByName("rootDirectory"); value.String() != "." {
		t.Errorf("Expected rootDirectory to be '.', got '%v'", value)
	}
}
