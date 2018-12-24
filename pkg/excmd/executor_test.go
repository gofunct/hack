package excmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/gofunct/hack/pkg/cli"
	"github.com/gofunct/hack/pkg/excmd"
)

func TestExecutor_WithConnectedIO(t *testing.T) {
	cases := []struct {
		test   string
		cmd    string
		opts   []excmd.Option
		out    string
		stdout string
	}{
		{
			test: "not connected",
			cmd:  "/bin/sh",
			opts: []excmd.Option{excmd.WithArgs("-c", "echo foo")},
			out:  "foo\n",
		},
		{
			test:   "connected",
			cmd:    "/bin/sh",
			opts:   []excmd.Option{excmd.WithArgs("-c", "read i && echo $i-$i"), excmd.WithIOConnected()},
			out:    "foo-foo\n",
			stdout: "foo-foo\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			stdin := bytes.NewBufferString("foo\n")

			execer := excmd.NewExecutor(&cli.IO{Out: stdout, Err: stderr, In: stdin})

			out, err := execer.Exec(context.TODO(), tc.cmd, tc.opts...)
			if err != nil {
				t.Errorf("returned %v, want nil", err)
			}

			if got, want := string(out), tc.out; got != want {
				t.Errorf("returned %q, want %q", got, want)
			}

			if got, want := stdout.String(), tc.stdout; got != want {
				t.Errorf("printed %q, want %q", got, want)
			}
		})
	}
}
