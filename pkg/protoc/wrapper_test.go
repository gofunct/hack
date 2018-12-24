package protoc_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/gofunct/hack/pkg/tool"
	"github.com/spf13/afero"
	"k8s.io/utils/exec"
	"k8s.io/utils/exec/testing"

	"github.com/gofunct/hack/pkg/cli"
	"github.com/gofunct/hack/pkg/protoc"
)

func TestWrapper_Exec(t *testing.T) {
	dieIf := func(t *testing.T, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	calls := [][]string{}

	execer := &testingexec.FakeExec{CommandScript: make([]testingexec.FakeCommandAction, 6)}
	for i := range execer.CommandScript {
		execer.CommandScript[i] = func(cmd string, args ...string) exec.Cmd {
			calls = append(calls, append([]string{cmd}, args...))
			return testingexec.InitFakeCmd(&testingexec.FakeCmd{
				CombinedOutputScript: []testingexec.FakeCombinedOutputAction{
					func() ([]byte, error) { return []byte("\n"), nil },
				},
			}, cmd, args...)
		}
	}

	rootDir := cli.RootDir("/go/src/awesomeapp")
	binDir := rootDir.BinDir()
	protosDir := rootDir.Join("api", "protos")

	fs := afero.NewMemMapFs()
	dieIf(t, fs.MkdirAll(binDir, 0755))
	dieIf(t, fs.MkdirAll(protosDir, 0755))
	dieIf(t, afero.WriteFile(fs, rootDir.Join("api", "should_be_ignored.proto"), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, rootDir.Join("api", "should_be_ignored_proto"), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, filepath.Join(protosDir, "book.proto"), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, filepath.Join(protosDir, "types", "users.proto"), []byte{}, 0644))

	cfg := &protoc.Config{
		ImportDirs: []string{
			rootDir.Join("vendor", "github.com", "grpc-ecosystem", "grpc-gateway"),
			protosDir,
		},
		ProtosDir: "./api/protos",
		OutDir:    "./api",
		Plugins: []*protoc.Plugin{
			{Name: "go", Args: map[string]interface{}{"plugins": "grpc"}},
			{Name: "grpc-gateway", Args: map[string]interface{}{"logtostderr": true}},
			{Name: "swagger", Args: map[string]interface{}{"logtostderr": true}},
		},
	}

	wrapper := protoc.NewWrapper(cfg, fs, execer, cli.NopUI, &fakeToolRepository{}, rootDir)

	err := wrapper.Exec(context.TODO())
	if err != nil {
		t.Fatalf("returned %v, want nil", err)
	}

	cupaloy.SnapshotT(t, calls)
}

type fakeToolRepository struct {
	tool.Repository
}

func (*fakeToolRepository) BuildAll(context.Context) error { return nil }
