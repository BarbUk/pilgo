package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/cmd/internal/command"
)

var _ command.Interface = new(checkCmd)

func TestCheck(t *testing.T) {
	t.Run("Execute", testCheckExecute)
	t.Run("SetFlags", testCheckSetFlags)
}

func testCheckExecute(t *testing.T) {
	testCases := []struct {
		name string
		cmd  checkCmd
		want string
		err  error
	}{
		{
			name: "check",
			cmd:  checkCmd{},
			err:  nil,
		},
	}
	for _, tc := range testCases {
		testdata := filepath.Join("testdata", t.Name())
		t.Run(tc.name, func(t *testing.T) {
			tc.cmd.config = filepath.Join(testdata, "pilgrim.json")
			before, err := ioutil.ReadFile(tc.cmd.config)
			if err != nil {
				t.Fatal(err)
			}
			tc.cmd.cwd = filepath.Join(testdata, "targets")
			var bd strings.Builder
			if want, got := tc.err, tc.cmd.Execute(&bd); !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			golden := readFile(t, filepath.Join("testdata", t.Name())+".golden")
			if want, got := golden, bd.String(); got != want {
				t.Errorf(
					`"show" command output mismatch (-want +got):\n%s`,
					cmp.Diff(want, got),
				)
			}
			after, err := ioutil.ReadFile(tc.cmd.config)
			if err != nil {
				t.Fatal(err)
			}
			// This guarantees the config file has only been read, not written.
			if want, got := before, after; bytes.Compare(got, want) != 0 {
				t.Errorf("%s has been modified after command", tc.cmd.config)
			}
		})
	}
}

func testCheckSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(checkCmd{})
	testCases := []struct {
		flags map[string]string
		want  checkCmd
	}{
		{
			flags: nil,
			want:  checkCmd{},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  checkCmd
				fset = flag.NewFlagSet("show", flag.PanicOnError)
				args = make([]string, 0, len(tc.flags))
			)
			for name, value := range tc.flags {
				args = append(args, fmt.Sprintf("-%s=%s", name, value))
			}
			cmd.SetFlags(fset)
			t.Logf("args: %v", args)
			if err := fset.Parse(args); err != nil {
				t.Fatal(err)
			}
			if want, got := tc.want, cmd; !cmp.Equal(got, want, allowUnexported) {
				t.Errorf(
					"(*showCmd).SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}