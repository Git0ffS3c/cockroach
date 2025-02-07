// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package tsearch

import (
	"bytes"
	"embed"
	"path/filepath"
	"strings"
)

//go:embed stopwords/*
var stopwordFS embed.FS

var stopwordsMap map[string]map[string]struct{}

func init() {
	stopwordsMap = make(map[string]map[string]struct{})
	dir, err := stopwordFS.ReadDir("stopwords")
	if err != nil {
		panic("error loading stopwords: " + err.Error())
	}
	for _, f := range dir {
		filename := f.Name()
		name := strings.TrimSuffix(filename, ".stop")
		contents, err := stopwordFS.ReadFile(filepath.Join("stopwords", filename))
		if err != nil {
			panic("error loading stopwords: " + err.Error())
		}
		wordList := bytes.Fields(contents)
		stopwordsMap[name] = make(map[string]struct{}, len(wordList))
		for _, word := range wordList {
			stopwordsMap[name][string(word)] = struct{}{}
		}
	}
	// The simple text search config has no stopwords.
	stopwordsMap["simple"] = nil
}
