package src

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

type testCase struct {
	flagsAndArgs   []string // As the cli command would be written
	stdinPath      string
	expectedOutput string
	err            error // set to non nil if an error is expected
}

// Used to setup and run the cli table testing
func runTableTest(t *testing.T, testCases []testCase) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err.Error())
	}

	os.Stdin = pr
	os.Stdout = pw

	for _, tC := range testCases {

		app := SetupCli("dev")
		app.ExitErrHandler = func(cCtx *cli.Context, err error) {
			// Empty to not let the cli handle the errors
		}

		_, err := pw.WriteString(tC.stdinPath)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = app.Run(append([]string{"dummy"}, tC.flagsAndArgs...))
		if !errors.Is(err, tC.err) {
			t.Fatal(err.Error())
		}

		b := make([]byte, 4096)
		n, err := pr.Read(b)
		if err != nil {
			t.Fatal(err.Error())
		}

		stdout := string(b[:n])

		// Convert OS specific paths over to posix
		osSpecificPaths := strings.Split(stdout, "\n")
		var stdoutPosixPaths []string
		for _, path := range osSpecificPaths {
			stdoutPosixPaths = append(stdoutPosixPaths, filepath.ToSlash(path))
		}
		stdoutPosix := strings.Join(stdoutPosixPaths, "\n")

		if stdoutPosix != tC.expectedOutput+"\n" {
			t.Errorf("\nargs: '%s'\nstdin: '%s'\nresulted in:\n'%s'\ninstead of\n'%s'", strings.Join(tC.flagsAndArgs, " "), tC.stdinPath, stdout, tC.expectedOutput)
		}

	}

}

func TestRegexpFlag(t *testing.T) {

	testCases := []testCase{
		{
			flagsAndArgs:   []string{"d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-e", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-regexp", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{".*2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f2.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-e", ".*2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f2.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-e", "d1", "-e", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestLineStringsFlag(t *testing.T) {

	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-X", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"--line-strings", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-X", "d2f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f2.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-X", "d3"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-X", "d3f1.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f1.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-X", "d1", "-X", "d2f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f2.f",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestFileFlag(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-f", "./test/input-test1"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-f", "./test/input-test2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f1.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-f", "./test/input-test3"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
			err:            nil,
		},
	}
	runTableTest(t, testCases)

}

func TestCombinePatternSources(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-e", ".*f1.f", "-X", "d2f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f1.f",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestLeftFlag(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-l", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-l", "f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d1f2.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"--left", "f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d1f2.f",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestBaseDirectoryFlag(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-b", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-b", "d1"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"--base-directory", "d1"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestBaseDirectoryWithOtherFlags(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-ba", "d2.*.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2\ntest/d1/d2",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-bl", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-b", "-X", "d2f1.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"-b", "--color", "-X", "d2f1.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestMatchAllFlag(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-a", "d2.*.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f\ntest/d1/d2/d2f2.f",
			err:            nil,
		},
		{
			flagsAndArgs:   []string{"--match-all", "d2.*.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f\ntest/d1/d2/d2f2.f",
			err:            nil,
		},
	}

	runTableTest(t, testCases)
}

func TestNoMatchOutputFlag(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"-O", "abc", "-e", "definitely_wrong"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "abc",
			err:            ErrNoMatch,
		},
		{
			flagsAndArgs:   []string{"-O", ".", "-e", "definitely_wrong"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: ".",
			err:            ErrNoMatch,
		},
		{
			flagsAndArgs:   []string{"--no-match", ".", "-e", "definitely_wrong"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: ".",
			err:            ErrNoMatch,
		},
	}

	runTableTest(t, testCases)
}
