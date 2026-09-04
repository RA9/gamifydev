package pyharness

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseOutputSeparatesResultsFromTheLearnersOwnPrinting(t *testing.T) {
	stdout := "hello\n" +
		ResultMarker + "[0, true]\n" +
		"still going\n" +
		ResultMarker + "[1, false]\n" +
		ResultMarker + "[2, true]\n"

	got, logs := ParseOutput(stdout, 3)
	if want := []bool{true, false, true}; !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %v, want %v", got, want)
	}
	if want := []string{"hello", "still going"}; !reflect.DeepEqual(logs, want) {
		t.Fatalf("logs = %v, want %v", logs, want)
	}
}

func TestAResultThatNeverArrivedIsAFailure(t *testing.T) {
	// A run killed partway through must keep what it proved and claim nothing
	// about the rest. This is what lets a checkpoint feed a huge input to the
	// last check — a solution too slow for it still gets told which of the
	// earlier requirements it met, instead of a wall of failures that reads
	// like the program was wrong everywhere.
	got, _ := ParseOutput(ResultMarker+"[0, true]\n"+ResultMarker+"[1, true]\n", 4)
	if want := []bool{true, true, false, false}; !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %v, want %v", got, want)
	}
}

func TestNoOutputAtAllFailsEveryCheck(t *testing.T) {
	// The program crashed or was killed before a single check reported.
	got, logs := ParseOutput("", 3)
	if want := []bool{false, false, false}; !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %v, want %v", got, want)
	}
	if len(logs) != 0 {
		t.Fatalf("logs = %v, want none", logs)
	}
}

func TestAnOutOfRangeIndexIsIgnored(t *testing.T) {
	// Learner output that happens to start with the marker must not be able to
	// write past the end of the results slice.
	got, _ := ParseOutput(ResultMarker+"[9, true]\n"+ResultMarker+"[0, true]\n", 2)
	if want := []bool{true, false}; !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %v, want %v", got, want)
	}
}

func TestBothBuildersReportInTheSameShape(t *testing.T) {
	// ParseOutput reads whatever these two print. If they ever disagree on the
	// format, one language silently grades every submission as a failure.
	tests := []string{"1 == 1", "2 == 3"}
	for name, prog := range map[string]string{
		"Build":           Build("x = 1", tests),
		"BuildAssertions": BuildAssertions(map[string]string{"_out": "hi"}, map[string]int{"_exit": 0}, tests),
	} {
		if !strings.Contains(prog, `_gd_json.dumps([_gd_i, _gd_r])`) {
			t.Errorf("%s does not report one indexed result per check", name)
		}
		if !strings.Contains(prog, "flush=True") {
			t.Errorf("%s does not flush, so a killed run loses its partial results", name)
		}
	}
}
