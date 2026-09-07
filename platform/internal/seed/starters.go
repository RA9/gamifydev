package seed

import "fmt"

// Starter builders for the common input shapes.
//
// Every problem in the bank is answerable in C, Go or Python, and each of those
// needs the same three things written before the interesting part begins: read
// the input, do the work, print the answer. Writing that by hand three hundred
// times would bury the problems, and worse, would let the three drift — a Go
// starter that reads a different format from the C one is a bug the learner
// pays for.
//
// So the parsing and the printing are generated from the input shape, and each
// problem supplies only the hint that says what to do in between.

// triCount builds starters for the commonest shape: a count, then that many
// integers.
func triCount(hint string) map[string]string {
	return tri(
		fmt.Sprintf(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        int v;
        if (scanf("%%d", &v) != 1) return 1;
        /* %s */
    }
    printf("%%d\n", 0);
    return 0;
}
`, hint),
		fmt.Sprintf(`package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	values := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	// %s
	fmt.Println(0)
}
`, hint),
		fmt.Sprintf(`import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    values = [int(x) for x in data[1:1 + n]]
    # %s
    print(0)

main()
`, hint))
}

// triCountWith is triCount preceded by one extra integer on the first line —
// a threshold, a window, a limit.
func triCountWith(param, hint string) map[string]string {
	return tri(
		fmt.Sprintf(`#include <stdio.h>

int main(void) {
    int n, %s;
    if (scanf("%%d %%d", &n, &%s) != 2) return 1;
    for (int i = 0; i < n; i++) {
        int v;
        if (scanf("%%d", &v) != 1) return 1;
        /* %s */
    }
    printf("%%d\n", 0);
    return 0;
}
`, param, param, hint),
		fmt.Sprintf(`package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, %s int
	fmt.Fscan(in, &n, &%s)
	values := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	// %s
	fmt.Println(0)
}
`, param, param, hint),
		fmt.Sprintf(`import sys

def main():
    data = sys.stdin.read().split()
    n, %s = int(data[0]), int(data[1])
    values = [int(x) for x in data[2:2 + n]]
    # %s
    print(0)

main()
`, param, hint))
}

// triScalars builds starters for a fixed number of integers on one line.
func triScalars(names []string, hint string) map[string]string {
	cDecl, cFmt, cAddr := "", "", ""
	goDecl, goAddr := "", ""
	pyUnpack := ""
	for i, n := range names {
		sep := ", "
		if i == 0 {
			sep = ""
		}
		cDecl += sep + n
		cFmt += map[bool]string{true: "%d", false: " %d"}[i == 0]
		cAddr += ", &" + n
		goDecl += sep + n
		goAddr += ", &" + n
		pyUnpack += sep + n
	}
	idx := ""
	for i := range names {
		if i > 0 {
			idx += ", "
		}
		idx += fmt.Sprintf("int(data[%d])", i)
	}
	return tri(
		fmt.Sprintf(`#include <stdio.h>

int main(void) {
    long long %s;
    if (scanf("%s"%s) != %d) return 1;
    /* %s */
    printf("%%lld\n", 0LL);
    return 0;
}
`, cDecl, replaceAll(cFmt, "%d", "%lld"), cAddr, len(names), hint),
		fmt.Sprintf(`package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var %s int
	fmt.Fscan(in%s)
	// %s
	fmt.Println(0)
}
`, goDecl, goAddr, hint),
		fmt.Sprintf(`import sys

def main():
    data = sys.stdin.read().split()
    %s = %s
    # %s
    print(0)

main()
`, pyUnpack, idx, hint))
}

// triWords builds starters for a count followed by that many whitespace-
// separated words.
func triWords(hint string) map[string]string {
	return tri(
		fmt.Sprintf(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%%d", &n) != 1) return 1;
    static char word[100001][32];
    for (int i = 0; i < n; i++) {
        if (scanf("%%31s", word[i]) != 1) return 1;
    }
    /* %s */
    printf("%%d\n", 0);
    return 0;
}
`, hint),
		fmt.Sprintf(`package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := range words {
		fmt.Fscan(in, &words[i])
	}
	// %s
	fmt.Println(0)
}
`, hint),
		fmt.Sprintf(`import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    words = data[1:1 + n]
    # %s
    print(0)

main()
`, hint))
}

// replaceAll is strings.ReplaceAll, local so this file needs no import beyond
// fmt for what is otherwise pure formatting.
func replaceAll(s, old, new string) string {
	out := ""
	for {
		i := indexOf(s, old)
		if i < 0 {
			return out + s
		}
		out += s[:i] + new
		s = s[i+len(old):]
	}
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// triLine builds starters for a single line of text, read whole — spaces and
// all — rather than word by word.
func triLine(hint string) map[string]string {
	return tri(
		fmt.Sprintf(`#include <stdio.h>
#include <string.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    /* %s */
    printf("\n");
    return 0;
}
`, hint),
		fmt.Sprintf(`package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	// %s
	fmt.Println("")
}
`, hint),
		fmt.Sprintf(`import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    # %s
    print("")

main()
`, hint))
}

// triLines builds starters for a count followed by that many whole lines.
func triLines(hint string) map[string]string {
	return tri(
		fmt.Sprintf(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%%d", &n) != 1) return 1;
    getchar();
    static char line[1005][1005];
    for (int i = 0; i < n; i++) {
        if (!fgets(line[i], sizeof line[i], stdin)) line[i][0] = '\0';
        line[i][strcspn(line[i], "\n")] = '\0';
    }
    /* %s */
    printf("\n");
    return 0;
}
`, hint),
		fmt.Sprintf(`package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	lines := make([]string, n)
	for i := range lines {
		s, _ := in.ReadString('\n')
		lines[i] = strings.TrimRight(s, "\r\n")
	}
	// %s
	fmt.Println("")
}
`, hint),
		fmt.Sprintf(`import sys

def main():
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    lines = raw[1:1 + n]
    # %s
    print("")

main()
`, hint))
}
