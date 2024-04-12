package src

import (
	"os"
	"strings"
	"testing"
)

type testCase struct {
	flagsAndArgs   []string // As the cli command would be written
	stdinPath      string
	expectedOutput string
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

		app := SetupCli()

		_, err := pw.WriteString(tC.stdinPath)
		if err != nil {
			t.Fatal(err.Error())
		}

		err = app.Run(append([]string{"dummy"}, tC.flagsAndArgs...))
		if err != nil {
			t.Fatal(err.Error())
		}

		b := make([]byte, 4096)
		n, err := pr.Read(b)
		if err != nil {
			t.Fatal(err.Error())
		}

		stdout := string(b[:n])

		if stdout != tC.expectedOutput+"\n" {
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
		},
		{
			flagsAndArgs:   []string{"-e", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
		},
		{
			flagsAndArgs:   []string{"-regexp", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
		},
		{
			flagsAndArgs:   []string{".*2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f2.f",
		},
		{
			flagsAndArgs:   []string{"-e", ".*2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f2.f",
		},
		{
			flagsAndArgs:   []string{"-e", "d1", "-e", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
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
		},
		{
			flagsAndArgs:   []string{"--line-strings", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
		},
		{
			flagsAndArgs:   []string{"-X", "d2f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f2.f",
		},
		{
			flagsAndArgs:   []string{"-X", "d3"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3",
		},
		{
			flagsAndArgs:   []string{"-X", "d3f1.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d3/d3f1.f",
		},
		{
			flagsAndArgs:   []string{"-X", "d1", "-X", "d2f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f2.f",
		},
	}

	runTableTest(t, testCases)
}

func TestFileFlag(t *testing.T) {
	testCases := []testCase{
		{
			flagsAndArgs:   []string{"d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f",
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
		},
		{
			flagsAndArgs:   []string{"-l", "f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d1f2.f",
		},
		{
			flagsAndArgs:   []string{"--left", "f2.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d1f2.f",
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
		},
		{
			flagsAndArgs:   []string{"-b", "d1"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1",
		},
		{
			flagsAndArgs:   []string{"--base-directory", "d1"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1",
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
		},
		{
			flagsAndArgs:   []string{"-bl", "d2"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1",
		},
		{
			flagsAndArgs:   []string{"-b", "-X", "d2f1.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
		},
		{
			flagsAndArgs:   []string{"-b", "--color", "-X", "d2f1.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2",
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
		},
		{
			flagsAndArgs:   []string{"--match-all", "d2.*.f"},
			stdinPath:      "test/d1/d2/d3",
			expectedOutput: "test/d1/d2/d2f1.f\ntest/d1/d2/d2f2.f",
		},
	}

	runTableTest(t, testCases)
}
