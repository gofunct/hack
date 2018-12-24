package cli_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/gofunct/hack/pkg/cli"
)

func TestUI(t *testing.T) {
	defer func(b bool) { color.NoColor = b }(color.NoColor)
	color.NoColor = true

	want := `  ➜  section 1
  ▸  subsection 1.1
     ✔  created
     ╌  skipped
     ✔  ok

  ▸  subsection 1.2
     ✗  failure

  ➜  section 2
     ✗  fail!!!
`

	out := new(bytes.Buffer)
	ui := cli.NewUI(&cli.IO{Out: out, In: new(bytes.Buffer)})

	ui.Section("section 1")
	ui.Subsection("subsection 1.1")
	ui.ItemSuccess("created")
	ui.ItemSkipped("skipped")
	ui.ItemSuccess("ok")
	ui.Subsection("subsection 1.2")
	ui.ItemFailure("failure")
	ui.Section("section 2")
	ui.ItemFailure("fail!!!")

	if got := out.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

type errReader struct {
}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("failed to read")
}

func TestUI_Confirm(t *testing.T) {
	type TestContext struct {
		in, out, err *bytes.Buffer
		ui           cli.UI
	}

	createTestContext := func() *TestContext {
		in := new(bytes.Buffer)
		out := new(bytes.Buffer)
		return &TestContext{
			in:  in,
			out: out,
			ui:  cli.NewUI(&cli.IO{Out: out, In: in}),
		}
	}

	cases := []struct {
		test      string
		errMsgCnt int
		input     string
		output    bool
	}{
		{
			test:      "with inputs 'n'",
			errMsgCnt: 0,
			input:     "n\n",
			output:    false,
		},
		{
			test:      "with inputs 'Y'",
			errMsgCnt: 0,
			input:     "Y\n",
			output:    true,
		},
		{
			test:      "with inputs 2 invalid chars and Y",
			errMsgCnt: 2,
			input:     "y\nN\nY\n",
			output:    true,
		},
		{
			test:      "with inputs 1 invalid chars and n",
			errMsgCnt: 2,
			input:     "N\ny\nn\n",
			output:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			ctx := createTestContext()
			ctx.in.WriteString(c.input)

			ok, err := ctx.ui.Confirm(c.test)

			if got, want := ctx.out.String(), c.test; !strings.HasPrefix(got, want) {
				t.Errorf("Confirm() wrote %q, want %q", got, want)
			}

			wantErrMsg := "input must be Y or n\n"
			if got, want := strings.Count(ctx.out.String(), wantErrMsg), c.errMsgCnt; got != want {
				t.Errorf("Confirm() wrote %q as error %d times, want %d times", wantErrMsg, got, want)
			}

			if err != nil {
				t.Errorf("Confirm() should not return errors, but returned %v", err)
			}

			if got, want := ok, c.output; got != want {
				t.Errorf("Confirm() returned %t, want %t", got, want)
			}
		})
	}

	t.Run("when failed to read", func(t *testing.T) {
		ui := cli.NewUI(&cli.IO{Out: new(bytes.Buffer), In: &errReader{}})

		ok, err := ui.Confirm("test")

		if err == nil {
			t.Error("Confirm() should return an error")
		}

		if ok {
			t.Error("Confirm() should return false")
		}
	})
}
