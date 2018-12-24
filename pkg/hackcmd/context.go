package hackcmd

import (
	"os"
	"path/filepath"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/utils/exec"

	"github.com/gofunct/hack/pkg/cli"
	"github.com/gofunct/hack/pkg/protoc"
)

// Ctx contains the runtime context of grpai.
type Ctx struct {
	FS     afero.Fs
	Viper  *viper.Viper
	Execer exec.Interface
	IO     *cli.IO

	RootDir   cli.RootDir
	insideApp bool

	Config       Config
	Build        BuildConfig
	ProtocConfig protoc.Config
}

// Config stores general setting params and provides accessors for them.
type Config struct {
	Package string
	Hack   struct {
		ServerDir string
	}
}

// BuildConfig contains the build metadata.
type BuildConfig struct {
	AppName   string
	Version   string
	Revision  string
	BuildDate string
	Prebuilt  bool
}

// Init initializes the runtime context.
func (c *Ctx) Init() error {
	if c.RootDir == "" {
		dir, _ := os.Getwd()
		c.RootDir = cli.RootDir(dir)
	}

	if c.IO == nil {
		c.IO = cli.DefaultIO()
	}

	if c.FS == nil {
		c.FS = afero.NewOsFs()
	}

	if c.Viper == nil {
		c.Viper = viper.New()
	}

	c.Viper.SetFs(c.FS)

	if c.Execer == nil {
		c.Execer = exec.New()
	}

	if c.Build.AppName == "" {
		c.Build.AppName = "hack"
	}

	return errors.WithStack(c.loadConfig())
}

func (c *Ctx) loadConfig() error {
	c.Viper.SetConfigName(c.Build.AppName)
	for dir := c.RootDir.String(); dir != "/"; dir = filepath.Dir(dir) {
		c.Viper.AddConfigPath(dir)
	}

	err := c.Viper.ReadInConfig()
	if err != nil {
		zap.L().Info("failed to find config file", zap.Error(err))
		return nil
	}

	c.insideApp = true
	c.RootDir = cli.RootDir(filepath.Dir(c.Viper.ConfigFileUsed()))

	err = c.Viper.Unmarshal(&c.Config)
	if err != nil {
		zap.L().Warn("failed to parse config", zap.Error(err))
		return errors.WithStack(err)
	}

	err = c.Viper.UnmarshalKey("protoc", &c.ProtocConfig)
	if err != nil {
		zap.L().Warn("failed to parse protoc config", zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}

// IsInsideApp returns true if the current working directory is inside a hack project.
func (c *Ctx) IsInsideApp() bool {
	return c.insideApp
}

// CtxSet is a provider set that includes modules contained in Ctx.
var CtxSet = wire.NewSet(
	ProvideFS,
	ProvideViper,
	ProvideExecer,
	ProvideIO,
	ProvideRootDir,
	ProvideConfig,
	ProvideBuildConfig,
	ProvideProtocConfig,
)

func ProvideFS(c *Ctx) afero.Fs                 { return c.FS }
func ProvideViper(c *Ctx) *viper.Viper          { return c.Viper }
func ProvideExecer(c *Ctx) exec.Interface       { return c.Execer }
func ProvideIO(c *Ctx) *cli.IO                  { return c.IO }
func ProvideRootDir(c *Ctx) cli.RootDir         { return c.RootDir }
func ProvideConfig(c *Ctx) *Config              { return &c.Config }
func ProvideBuildConfig(c *Ctx) *BuildConfig    { return &c.Build }
func ProvideProtocConfig(c *Ctx) *protoc.Config { return &c.ProtocConfig }
