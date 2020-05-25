package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

func readAll(paths []string) ([]TestFile, error) {
	ts := make([]TestFile, 0, len(paths))

	for i := range paths {
		buf, err := ioutil.ReadFile(paths[i])
		if err != nil {
			return nil, errors.WithStack(err)
		}
		var t TestFile
		if err := xml.Unmarshal(buf, &t); err != nil {
			log.Printf("%v\n", errors.Wrapf(err, "cannot unmarshal xml from file at %s", paths[i]))
			continue
		}
		if len(t.TestSuite) == 0 {
			log.Printf("no test suite for file at %s\n", paths[i])
			continue
		}
		ts = append(ts, t)
	}

	return ts, nil
}

type TestFile struct {
	XMLName   xml.Name    `xml:"testsuites"`
	TestSuite []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Tests     int64      `xml:"tests,attr"`
	Failures  int64      `xml:"failures,attr"`
	Errors    int64      `xml:"errors,attr"`
	Skip      int64      `xml:"skip,attr"`
	Time      float64    `xml:"time,attr"`
	Name      string     `xml:"name,attr"`
	TestCases []TestCase `xml:"testcase"`
}

type TestCase struct {
	Name string  `xml:"name,attr"`
	Time float64 `xml:"time,attr"`
}
