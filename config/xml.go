// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func ParseXMLFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return parseXML(f)
}

// 从一个xml字符串初始化日志系统。
func ParseXMLString(xml string) (*Config, error) {
	return parseXML(strings.NewReader(xml))
}

// 从一个xml格式的reader初始化Config
func parseXML(r io.Reader) (*Config, error) {
	var cfg *Config = nil
	var t xml.Token
	var err error

	d := xml.NewDecoder(r)
	for t, err = d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			c := &Config{
				Parent: cfg,
				Name:   token.Name.Local,
				Attrs:  make(map[string]string),
			}
			for _, v := range token.Attr {
				c.Attrs[v.Name.Local] = v.Value
			}

			if cfg != nil {
				if cfg.Items == nil {
					cfg.Items = make(map[string]*Config)
				}

				if _, found := cfg.Items[token.Name.Local]; found {
					return nil, fmt.Errorf("parseXML:重复的元素名[%v]", token.Name.Local)
				}

				cfg.Items[token.Name.Local] = c
			}
			cfg = c
		case xml.EndElement:
			if cfg.Parent != nil {
				cfg = cfg.Parent
			}
		default: // 可能还有ProcInst,CharData,Comment等用不到的标签
			continue
		}
	} // end for

	if err != io.EOF {
		return nil, err
	}

	if err = check(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
