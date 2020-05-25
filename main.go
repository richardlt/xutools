package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "xutools",
		Commands: []*cli.Command{
			{
				Name:   "pretty",
				Action: runPretty,
			},
			{
				Name:   "sort-duration",
				Action: runSortDuration,
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "limit",
						Value: 30,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%+v\n", errors.WithStack(err))
	}
}

func runPretty(c *cli.Context) error {
	ts, err := readAll(c.Args().Slice())
	if err != nil {
		return err
	}

	for i := range ts {
		for _, s := range ts[i].TestSuite {
			var counter string
			var details string
			colorFunc := Green
			if s.Skip > 0 {
				colorFunc = Yellow
				details += Yellow(" %d skipped test(s)", s.Skip)
			}
			if s.Failures > 0 || s.Errors > 0 {
				colorFunc = Red
				details += Red(" %d failed test(s)", s.Failures+s.Errors)
			}
			if s.Failures > 0 || s.Errors > 0 || s.Skip > 0 {
				counter = colorFunc(fmt.Sprintf("%d/%d", s.Tests-s.Failures-s.Errors-s.Skip, s.Tests))
			} else {
				counter = colorFunc(fmt.Sprintf("%d", s.Tests))
			}
			if details != "" {
				details = " with" + details
			}
			fmt.Printf("%s: %s successful tests in %fs%s\n", s.Name, counter, s.Time, details)
		}
	}

	return nil
}

type displayTestCase struct {
	name  string
	time  float64
	suite string
}

func runSortDuration(c *cli.Context) error {
	limit := c.Int64("limit")

	ts, err := readAll(c.Args().Slice())
	if err != nil {
		return err
	}

	var allTestCases []displayTestCase
	for i := range ts {
		for j := range ts[i].TestSuite {
			for k := range ts[i].TestSuite[j].TestCases {
				allTestCases = append(allTestCases, displayTestCase{
					name:  ts[i].TestSuite[j].TestCases[k].Name,
					time:  ts[i].TestSuite[j].TestCases[k].Time,
					suite: ts[i].TestSuite[j].Name,
				})
			}
		}
	}

	sort.Slice(allTestCases, func(i, j int) bool {
		return allTestCases[i].time > allTestCases[j].time
	})
	for i := range allTestCases {
		if int64(i) == limit {
			break
		}
		fmt.Printf("%s/%s - %f\n", allTestCases[i].suite, allTestCases[i].name, allTestCases[i].time)
	}

	return nil
}
