// Copyright 2020 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsImageMatched(t *testing.T) {
	testCases := []struct {
		testName  string
		value     string
		name      string
		isMatched bool
	}{
		{
			testName:  "identical",
			value:     "nginx",
			name:      "nginx",
			isMatched: true,
		},
		{
			testName:  "name is match",
			value:     "nginx:12345",
			name:      "nginx",
			isMatched: true,
		},
		{
			testName:  "name is not a match",
			value:     "apache:12345",
			name:      "nginx",
			isMatched: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			assert.Equal(t, tc.isMatched, IsImageMatched(tc.value, tc.name))
		})
	}
}

func TestSplit(t *testing.T) {
	testCases := []struct {
		testName string
		value    string
		name     string
		tag      string
	}{
		{
			testName: "no tag",
			value:    "nginx",
			name:     "nginx",
			tag:      "",
		},
		{
			testName: "with tag",
			value:    "nginx:1.2.3",
			name:     "nginx",
			tag:      ":1.2.3",
		},
		{
			testName: "with digest",
			value:    "nginx@12345",
			name:     "nginx",
			tag:      "@12345",
		},
		// {
		// 	testName: "no tag with repository",
		// 	value:    "my-repository/nginx@12345",
		// 	name:     "nginx",
		// 	tag:      "@12345",
		// },
		// {
		// 	testName: "with tag with repository",
		// 	value:    "my-repository/nginx@1.2.3",
		// 	name:     "nginx",
		// 	tag:      ":12345",
		// },
		// {
		// 	testName: "with digest wit repository",
		// 	value:    "my-repository/nginx@12345",
		// 	name:     "nginx",
		// 	tag:      "@12345",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			name, tag := Split(tc.value)
			assert.Equal(t, tc.name, name)
			assert.Equal(t, tc.tag, tag)
		})
	}
}

func TestNewImageName(t *testing.T) {
	testCases := []struct {
		testName       string
		origName       string
		searchPattern  string
		newNamePattern string
		resultNewName  string
	}{
		{
			testName:       "simple replace",
			origName:       "adoptopenjdk",
			searchPattern:  "open",
			newNamePattern: "closed",
			resultNewName:  "adoptclosedjdk",
		},
		{
			testName:       "simple regex",
			origName:       "adoptopenjdk",
			searchPattern:  "^adopt",
			newNamePattern: "my",
			resultNewName:  "myopenjdk",
		},
		{
			testName:       "no match",
			origName:       "nginx",
			searchPattern:  "adoptopenjdk",
			newNamePattern: "my",
			resultNewName:  "nginx",
		},
		{
			testName:       "simple regex with capture groups",
			origName:       "nginx",
			searchPattern:  "(.*)",
			newNamePattern: "my-registry/$1",
			resultNewName:  "my-registry/nginx",
		},
		{
			testName:       "complex regex with capture groups",
			origName:       "nginx",
			searchPattern:  "([^/]*/)?(.*)",
			newNamePattern: "my-registry/$2",
			resultNewName:  "my-registry/nginx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			name := NewImageName(tc.searchPattern, tc.origName, tc.newNamePattern)
			assert.Equal(t, tc.resultNewName, name)
		})
	}
}
