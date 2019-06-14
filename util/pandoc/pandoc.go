package pandoc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"tfw.io/Go/fsindex/util"
)

// Wrapper wraps pandoc of course.
type Wrapper struct {
	cmdFlags []string

	PandocEXE  string
	InputFile  string
	InputType  string
	InputEXT   string
	OutputType string

	DoStandalone        bool
	DoNumericalHeadings bool
	DoTOC               bool
	NoHighlight         bool

	// MediaPath             string
	// WindowsOffice2013Path string
	// LibreOfficePath       string
	// Extensions    []string
	// TemplatesPath string `tools/templates`
}

func (w *Wrapper) addFlag(flag ...string) {
	for _, item := range flag {
		w.cmdFlags = append(w.cmdFlags, item)
	}
}
func (w *Wrapper) addFlags(flag []string) {
	for _, item := range flag {
		w.cmdFlags = append(w.cmdFlags, item)
	}
}

// Create makes a new Wrapper with some default settings,
// and some integrated.
func Create(app string, flags string, extensions string, tpl string) Wrapper {
	w := Wrapper{
		NoHighlight:         true,
		DoTOC:               true,
		DoStandalone:        false,
		DoNumericalHeadings: false,
	}
	w.create(app, flags, extensions, tpl)
	return w
}

func (w *Wrapper) create(
	exeFilePath string,
	pdFlags string,
	pdExtensions string,
	tplFilePath string) {
	w.PandocEXE = exeFilePath
	w.InputEXT = pdExtensions
	w.cmdFlags = []string{}

	if len(w.InputType) == 0 {
		w.InputType = "markdown_mmd"
	}
	if len(w.OutputType) == 0 {
		w.OutputType = "html5"
	}
	if w.DoStandalone {
		w.addFlag("--standalone")
	}
	if w.DoTOC {
		w.addFlag("--toc")
	}
	if w.DoNumericalHeadings {
		w.addFlag("-N")
	}
	if w.NoHighlight {
		w.addFlag("--no-highlight")
	}
	w.addFlags(w.cmdFlags)
	w.addFlag("-f", w.InputType+w.InputEXT)
	w.addFlag("-t", w.OutputType)
	w.addFlag("--template", tplFilePath)
}

// Args gets the command-arguments.
func (w *Wrapper) Args() []string {
	return w.cmdFlags
}

// Do does
func (w *Wrapper) Do(
	pInputFile string,
	pStdOutBuffer *bytes.Buffer,
	pStdErrBuffer *bytes.Buffer,
	pUseStandardFallback bool) error {

	// mInputFilePath := strings.TrimLeft(pInputFilePath, "/")
	mCommand := exec.Command(w.PandocEXE)
	w.addFlag(util.UnixSlash(pInputFile)) // *nix friendly
	mCommand.Args = w.cmdFlags

	if pUseStandardFallback == true {

		if pStdOutBuffer == nil {
			mCommand.Stdout = os.Stdout
		} else {
			mCommand.Stdout = pStdOutBuffer
		}
		if pStdErrBuffer == nil {
			mCommand.Stderr = os.Stderr
		} else {
			mCommand.Stderr = pStdErrBuffer
		}
	} else {
		mCommand.Stdout = pStdOutBuffer
		mCommand.Stderr = pStdErrBuffer
	}
	fmt.Sprintln("Output: ", os.Stderr, w.cmdFlags)
	// fmt.Println(w.cmdFlags)
	// fmt.Sprintln()
	return mCommand.Run()

	// it appears that this only validates existance of the supplied executable.
	// ca := CatArrayPad(p.cmdFlags, " ")
	// fmt.Sprintln(os.Stderr, ca)
	// fmt.Sprintln(os.Stderr, "- running on input:")
	// fmt.Sprintln(os.Stderr, "  - ", inp)

}

// Version i.e. `pandoc --version`
func Version(exe string) {
	cmd := exec.Command(exe, "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
