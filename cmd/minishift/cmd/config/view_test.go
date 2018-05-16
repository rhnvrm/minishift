/*
Copyright (C) 2017 Red Hat, Inc.

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

package config

import (
	"bytes"
	"testing"

	"github.com/minishift/minishift/cmd/testing/cli"
	"github.com/stretchr/testify/assert"
)

type templateExpectedOutputTest struct {
	expectedString string
	template       string
}

var (
	configTest = map[string]interface{}{
		"iso-url":   "http://foo.bar/minishift-centos.iso",
		"vm-driver": "kvm",
		"cpus":      4,
		"disk-size": "20g",
		"v":         5,
		"show-libmachine-logs":      true,
		"log_dir":                   "/etc/hosts",
		"ReminderWaitPeriodInHours": 99,
	}

	templateExpectedOutputTestList = []templateExpectedOutputTest{
		{
			expectedString: "- ReminderWaitPeriodInHours: 99\n" +
				"- cpus                 : 4\n" +
				"- disk-size            : 20g\n" +
				"- iso-url              : http://foo.bar/minishift-centos.iso\n" +
				"- log_dir              : /etc/hosts\n" +
				"- show-libmachine-logs : true\n" +
				"- v                    : 5\n" +
				"- vm-driver            : kvm\n",
			template: "- {{.ConfigKey | printf \"%-21s\"}}: {{.ConfigValue}}",
		},
		{
			expectedString: "- ReminderWaitPeriodInHours: 99\n" +
				"- cpus: 4\n" +
				"- disk-size: 20g\n" +
				"- iso-url: http://foo.bar/minishift-centos.iso\n" +
				"- log_dir: /etc/hosts\n" +
				"- show-libmachine-logs: true\n" +
				"- v: 5\n" +
				"- vm-driver: kvm\n",
			template: "- {{.ConfigKey}}: {{.ConfigValue}}",
		},
	}
)

func TestConfigView(t *testing.T) {
	for _, tt := range templateExpectedOutputTestList {
		tmpMinishiftHomeDir := cli.SetupTmpMinishiftHome(t)
		tee := cli.CreateTee(t, true)
		defer cli.TearDown(tmpMinishiftHomeDir, tee)

		template := determineTemplate(tt.template)
		configView(configTest, template, tee.StdoutBuffer)
		assert.Equal(t, tt.expectedString, tee.StdoutBuffer.String())
	}
}

func TestAlignTabbedLines(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		args     args
		expected string
	}{
		{
			args: args{
				lines: []string{
					"small\t: SomeValue",
					"veryLongKeyForTestingAgainstSmallKey\t: AnotherValue",
				},
			},
			expected: "small\t\t\t\t\t: SomeValue\nveryLongKeyForTestingAgainstSmallKey\t: AnotherValue\n",
		},
		{
			args: args{
				lines: []string{
					"smallKeyOne\t: SomeValue",
					"smallKeyTwo\t: AnotherValue",
				},
			},
			expected: "smallKeyOne\t: SomeValue\nsmallKeyTwo\t: AnotherValue\n",
		},
		{
			args: args{
				lines: []string{
					"urlKey\t: 'http://minishift.io/'",
				},
			},
			expected: "urlKey\t: 'http://minishift.io/'\n",
		},
	}
	for _, tt := range tests {
		writer := &bytes.Buffer{}
		alignTabbedLines(writer, tt.args.lines)
		actual := writer.String()
		assert.Equal(t, tt.expected, actual)
	}
}
