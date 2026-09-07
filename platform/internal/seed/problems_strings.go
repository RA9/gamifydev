package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// stringProblems: reading text a character at a time, and the state you have to
// carry while you do it. In C that means no string type and a buffer you own,
// which is exactly why these are worth writing three times.
var stringProblems = []seedProblem{
	{
		slug: "emphasis-pass", title: "Emphasis pass", difficulty: "easy", topic: "Strings",
		statement: `A chat client renders emphasis by case. A word ending in ~!~ is shouted
(uppercased); a word ending in ~?~ is muttered (lowercased); everything else is
left exactly as it was.

**Input.** One line. Words are separated by single spaces, and punctuation stays
on the word.

**Output.** The rendered line.

    Input                   Output
    Look out! is it Safe?   Look OUT! is it safe?`,
		timeLimitMs: 4000,
		starters:    triLine("rewrite each word by its final character"),
		tests: []store.Check{
			exact("Shouts and mutters in one pass", "Look out! is it Safe?\n", "Look OUT! is it safe?"),
			exact("Plain words are untouched", "nothing to see\n", "nothing to see"),
			exact("An empty line stays empty", "\n", ""),
			exactHid("Punctuation inside a word does not count", "wait!for it\n", "wait!for it"),
			exactHid("A word that is only punctuation still works", "well ! ?\n", "well ! ?"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int n = strlen(line), start = 0;
    for (int i = 0; i <= n; i++) {
        if (line[i] != ' ' && line[i] != '\0') continue;
        int end = i;
        char last = end > start ? line[end - 1] : '\0';
        for (int j = start; j < end; j++) {
            if (last == '!') putchar(toupper((unsigned char)line[j]));
            else if (last == '?') putchar(tolower((unsigned char)line[j]));
            else putchar(line[j]);
        }
        if (line[i] == ' ') putchar(' ');
        start = i + 1;
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

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
	words := strings.Split(line, " ")
	for i, w := range words {
		switch {
		case strings.HasSuffix(w, "!"):
			words[i] = strings.ToUpper(w)
		case strings.HasSuffix(w, "?"):
			words[i] = strings.ToLower(w)
		}
	}
	fmt.Println(strings.Join(words, " "))
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    out = []
    for w in line.split(" "):
        if w.endswith("!"):
            out.append(w.upper())
        elif w.endswith("?"):
            out.append(w.lower())
        else:
            out.append(w)
    print(" ".join(out))

main()
`},
	},
	{
		slug: "redact-long-numbers", title: "Redact long numbers", difficulty: "medium", topic: "Strings",
		statement: `A support tool hides anything that looks like an account number before a
transcript is shared. A run of **six or more** digits is replaced by the same
number of ~#~ characters. Shorter runs — years, quantities, door numbers — are
left alone.

**Input.** One line.

**Output.** The redacted line.

    Input                       Output
    call 5551234 or ext 12345   call ####### or ext 12345

The replacement is the same length as what it hides, so the transcript still
lines up.`,
		timeLimitMs: 4000,
		starters:    triLine("replace runs of six or more digits with the same number of hashes"),
		tests: []store.Check{
			exact("Hides a long run and keeps a short one", "call 5551234 or ext 12345\n", "call ####### or ext 12345"),
			exact("A year is not an account number", "since 1999\n", "since 1999"),
			exact("Nothing to redact", "no digits here\n", "no digits here"),
			exact("Exactly six digits is long enough", "id 123456\n", "id ######"),
			exactHid("Five digits is not", "id 12345\n", "id 12345"),
			exactHid("A run at the very end is caught", "ref 9876543\n", "ref #######"),
			exactHid("Two long runs in one line", "1234567 and 7654321\n", "####### and #######"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int n = strlen(line), i = 0;
    while (i < n) {
        if (isdigit((unsigned char)line[i])) {
            int j = i;
            while (j < n && isdigit((unsigned char)line[j])) j++;
            int run = j - i;
            for (int k = 0; k < run; k++) putchar(run >= 6 ? '#' : line[i + k]);
            i = j;
        } else {
            putchar(line[i++]);
        }
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

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
	var b strings.Builder
	for i := 0; i < len(line); {
		if line[i] >= '0' && line[i] <= '9' {
			j := i
			for j < len(line) && line[j] >= '0' && line[j] <= '9' {
				j++
			}
			if j-i >= 6 {
				b.WriteString(strings.Repeat("#", j-i))
			} else {
				b.WriteString(line[i:j])
			}
			i = j
		} else {
			b.WriteByte(line[i])
			i++
		}
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    out = []
    i = 0
    while i < len(line):
        if line[i].isdigit():
            j = i
            while j < len(line) and line[j].isdigit():
                j += 1
            run = j - i
            out.append("#" * run if run >= 6 else line[i:j])
            i = j
        else:
            out.append(line[i])
            i += 1
    print("".join(out))

main()
`},
	},
	{
		slug: "badge-initials", title: "Badge initials", difficulty: "medium", topic: "Strings",
		statement: `A conference badge shows a person's initials. Name particles — ~van~,
~der~, ~de~, ~den~, ~da~, ~di~, ~bin~, ~al~ — are skipped, but only when written
in lower case. A capitalised ~Van~ is somebody's actual name and gets an initial
like any other word.

**Input.** One line holding the name. There may be more than one space between
words.

**Output.** The uppercase initials, with nothing between them. Nobody has no
initials, so an empty name prints an empty line.

    Input                   Output
    ada van der lovelace    AL

    Input                   Output
    Vincent Van Gogh        VVG`,
		timeLimitMs: 4000,
		starters:    triLine("take the first letter of every word that is not a lower-case particle"),
		tests: []store.Check{
			exact("Skips lower-case particles", "ada van der lovelace\n", "AL"),
			exact("A capitalised particle is a name", "Vincent Van Gogh\n", "VVG"),
			exact("An ordinary name is untouched", "Grace Hopper\n", "GH"),
			exact("Nobody has no initials", "\n", ""),
			exactHid("Extra spaces do not create empty initials", "  Ada   Lovelace \n", "AL"),
			exactHid("A particle-looking word that is not one", "Alan Turing\n", "AT"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

static int particle(const char *w) {
    static const char *p[] = {"van", "der", "de", "den", "da", "di", "bin", "al"};
    for (int i = 0; i < 8; i++) if (strcmp(w, p[i]) == 0) return 1;
    return 0;
}

int main(void) {
    static char line[100005], word[1005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int n = strlen(line), start = 0;
    for (int i = 0; i <= n; i++) {
        if (i < n && !isspace((unsigned char)line[i])) continue;
        int len = i - start;
        if (len > 0 && len < 1000) {
            memcpy(word, line + start, len);
            word[len] = '\0';
            int lower = 1;
            for (int j = 0; j < len; j++) if (isupper((unsigned char)word[j])) lower = 0;
            if (!(lower && particle(word))) putchar(toupper((unsigned char)word[0]));
        }
        start = i + 1;
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var particles = map[string]bool{"van": true, "der": true, "de": true, "den": true,
	"da": true, "di": true, "bin": true, "al": true}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	var b strings.Builder
	for _, w := range strings.Fields(strings.TrimRight(line, "\r\n")) {
		if w == strings.ToLower(w) && particles[w] {
			continue
		}
		b.WriteString(strings.ToUpper(w[:1]))
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

PARTICLES = {"van", "der", "de", "den", "da", "di", "bin", "al"}

def main():
    line = sys.stdin.readline().rstrip("\n")
    out = []
    for w in line.split():
        if w.islower() and w in PARTICLES:
            continue
        out.append(w[0].upper())
    print("".join(out))

main()
`},
	},
	{
		slug: "ticker-window", title: "Ticker window", difficulty: "medium", topic: "Strings",
		statement: `A departure board scrolls a message right to left through a narrow
window. The message repeats forever with exactly three spaces between the end of
one pass and the start of the next.

**Input.** The message on the first line, then ~width t~ on the second.

**Output.** The ~width~ characters showing at tick ~t~. At tick 0 the window
starts at the first character; each tick moves the message one place left.

    Input        Output
    HELLO        HEL
    3 0

~t~ can be enormous — far past two billion — and the window can be wider than
the message.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    static char text[100005];
    if (!fgets(text, sizeof text, stdin)) text[0] = '\0';
    text[strcspn(text, "\n")] = '\0';
    long long width, t;
    if (scanf("%lld %lld", &width, &t) != 2) return 1;
    /* print width characters of the looping message, starting at tick t */
    printf("\n");
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	text, _ := in.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	var width, t int
	fmt.Fscan(in, &width, &t)
	// print width characters of the looping message, starting at tick t
	fmt.Println("")
}
`, `import sys

def main():
    lines = sys.stdin.read().split("\n")
    text = lines[0]
    width, t = (int(x) for x in lines[1].split()[:2])
    # print width characters of the looping message, starting at tick t
    print("")

main()
`),
		tests: []store.Check{
			exact("The start of the message", "HELLO\n3 0\n", "HEL"),
			exact("Scrolled into the gap and round again", "HELLO\n3 6\n", "  H"),
			exact("A window wider than the message", "AB\n6 0\n", "AB   A"),
			exactHid("A tick far in the future", "HELLO\n5 8000000000\n", "HELLO"),
			exactHid("A window of nothing shows nothing", "HELLO\n0 3\n", ""),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char text[100005];
    if (!fgets(text, sizeof text, stdin)) text[0] = '\0';
    text[strcspn(text, "\n")] = '\0';
    long long width, t;
    if (scanf("%lld %lld", &width, &t) != 2) return 1;
    long long len = (long long)strlen(text) + 3;
    for (long long i = 0; i < width; i++) {
        long long k = (t + i) % len;
        putchar(k < (long long)strlen(text) ? text[k] : ' ');
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	text, _ := in.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	var width, t int
	fmt.Fscan(in, &width, &t)
	loop := text + "   "
	var b strings.Builder
	for i := 0; i < width; i++ {
		b.WriteByte(loop[(t+i)%len(loop)])
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

def main():
    lines = sys.stdin.read().split("\n")
    text = lines[0]
    width, t = (int(x) for x in lines[1].split()[:2])
    loop = text + "   "
    print("".join(loop[(t + i) % len(loop)] for i in range(width)))

main()
`},
	},
	{
		slug: "message-segments", title: "Message segments", difficulty: "medium", topic: "Strings",
		statement: `A messaging gateway bills by segment. Most characters cost one unit, but
any character listed as special costs two. A message of 160 units or fewer is a
single segment; anything longer is split into segments of 153 units, and a
part-filled final segment still counts.

**Input.** The message on the first line, then the special characters as one
run of characters on the second (possibly empty).

**Output.** How many segments are billed. An empty message is never sent, so it
costs nothing.

    Input      Output
    hello      1
    (blank)`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    static char text[100005], special[256];
    if (!fgets(text, sizeof text, stdin)) text[0] = '\0';
    text[strcspn(text, "\n")] = '\0';
    if (!fgets(special, sizeof special, stdin)) special[0] = '\0';
    special[strcspn(special, "\n")] = '\0';
    /* count units, then segments */
    printf("%d\n", 0);
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	text, _ := in.ReadString('\n')
	special, _ := in.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	special = strings.TrimRight(special, "\r\n")
	// count units, then segments
	fmt.Println(0)
}
`, `import sys

def main():
    lines = sys.stdin.read().split("\n")
    text = lines[0]
    special = lines[1] if len(lines) > 1 else ""
    # count units, then segments
    print(0)

main()
`),
		tests: []store.Check{
			tokens("A short message is one segment", "hello\n\n", "1"),
			tokens("One unit over splits it", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\n", "2"),
			tokens("An empty message is never sent", "\n\n", "0"),
			tokens("Special characters cost double", "$$$$$\n$\n", "1"),
			hid2("Exactly 160 units still fits in one", "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n$\n", "1"),
			hid2("One unit past 160 splits", "$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$a\n$\n", "2"),
			hid2("Only the listed characters cost double", "ab\n$\n", "1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char text[100005], special[256];
    if (!fgets(text, sizeof text, stdin)) text[0] = '\0';
    text[strcspn(text, "\n")] = '\0';
    if (!fgets(special, sizeof special, stdin)) special[0] = '\0';
    special[strcspn(special, "\n")] = '\0';
    long long units = 0;
    for (int i = 0; text[i]; i++) units += strchr(special, text[i]) && text[i] ? 2 : 1;
    long long segs;
    if (units == 0) segs = 0;
    else if (units <= 160) segs = 1;
    else segs = (units + 152) / 153;
    printf("%lld\n", segs);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	text, _ := in.ReadString('\n')
	special, _ := in.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	special = strings.TrimRight(special, "\r\n")
	units := 0
	for _, c := range text {
		if strings.ContainsRune(special, c) {
			units += 2
		} else {
			units++
		}
	}
	switch {
	case units == 0:
		fmt.Println(0)
	case units <= 160:
		fmt.Println(1)
	default:
		fmt.Println((units + 152) / 153)
	}
}
`,
			"python": `import sys

def main():
    lines = sys.stdin.read().split("\n")
    text = lines[0]
    special = lines[1] if len(lines) > 1 else ""
    units = sum(2 if c in special else 1 for c in text)
    if units == 0:
        print(0)
    elif units <= 160:
        print(1)
    else:
        print(-(-units // 153))

main()
`},
	},
	{
		slug: "sentence-caps", title: "Sentence caps", difficulty: "medium", topic: "Strings",
		statement: `Capitalise the first letter of every sentence and change nothing else —
not the middle of words, not the case of anything already written.

A sentence ends at ~.~, ~!~ or ~?~. The next letter after one of those, however
far away it is, starts the next sentence.

**Input.** One line.

**Output.** The corrected line.

    Input                            Output
    hello there. how are you?        Hello there. How are you?

Characters that are not letters never get capitalised, and never stop the search
for the letter that does.`,
		timeLimitMs: 4000,
		starters:    triLine("capitalise the first letter after each sentence ending"),
		tests: []store.Check{
			exact("Capitalises each sentence", "hello there. how are you? fine!\n", "Hello there. How are you? Fine!"),
			exact("Leaves the rest of the words alone", "iPhone and iPad.\n", "IPhone and iPad."),
			exact("An empty line stays empty", "\n", ""),
			exactHid("Skips past spaces and quotes to find the letter", "stop.  'go' now\n", "Stop.  'Go' now"),
			exactHid("An already capital letter is left as it is", "One. Two.\n", "One. Two."),
			exactHid("Trailing punctuation with nothing after it", "done.\n", "Done."),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int start = 1;
    for (int i = 0; line[i]; i++) {
        if (start && isalpha((unsigned char)line[i])) {
            line[i] = toupper((unsigned char)line[i]);
            start = 0;
        } else if (line[i] == '.' || line[i] == '!' || line[i] == '?') {
            start = 1;
        }
    }
    printf("%s\n", line);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	out := []rune(strings.TrimRight(line, "\r\n"))
	start := true
	for i, ch := range out {
		switch {
		case start && unicode.IsLetter(ch):
			out[i] = unicode.ToUpper(ch)
			start = false
		case ch == '.' || ch == '!' || ch == '?':
			start = true
		}
	}
	fmt.Println(string(out))
}
`,
			"python": `import sys

def main():
    out = list(sys.stdin.readline().rstrip("\n"))
    start = True
    for i, ch in enumerate(out):
        if start and ch.isalpha():
            out[i] = ch.upper()
            start = False
        elif ch in ".!?":
            start = True
    print("".join(out))

main()
`},
	},
	{
		slug: "strip-markup", title: "Strip markup", difficulty: "medium", topic: "Strings",
		statement: `Remove every ~<...>~ span along with its angle brackets.

The catch is the stray bracket. A ~<~ with no ~>~ anywhere after it is not the
start of a tag — it is a less-than sign somebody typed, and it stays. A tag ends
at the first ~>~ after it.

**Input.** One line.

**Output.** The line with tags removed.

    Input                 Output
    a <b>bold</b> move    a bold move

    Input      Output
    5 < 6      5 < 6`,
		timeLimitMs: 4000,
		starters:    triLine("copy characters through, skipping any complete tag"),
		tests: []store.Check{
			exact("Removes tags", "a <b>bold</b> move\n", "a bold move"),
			exact("A stray bracket stays", "5 < 6\n", "5 < 6"),
			exact("Nothing to strip", "plain\n", "plain"),
			exact("The whole line is one tag", "<hr>\n", ""),
			exactHid("A tag ends at the first closing bracket after it", "a < b <i>c</i>\n", "a c"),
			exactHid("A stray bracket after a real tag", "a <b> c < d\n", "a  c < d"),
			exactHid("A closing bracket on its own is ordinary text", "3 > 2\n", "3 > 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int i = 0;
    while (line[i]) {
        if (line[i] == '<') {
            char *close = strchr(line + i, '>');
            if (close) { i = (int)(close - line) + 1; continue; }
        }
        putchar(line[i++]);
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

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
	var b strings.Builder
	for i := 0; i < len(line); {
		if line[i] == '<' {
			if j := strings.IndexByte(line[i:], '>'); j >= 0 {
				i += j + 1
				continue
			}
		}
		b.WriteByte(line[i])
		i++
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    out = []
    i = 0
    while i < len(line):
        if line[i] == "<":
            j = line.find(">", i)
            if j != -1:
                i = j + 1
                continue
        out.append(line[i])
        i += 1
    print("".join(out))

main()
`},
	},
	{
		slug: "one-transposition", title: "One transposition", difficulty: "medium", topic: "Strings",
		statement: `A spell checker wants to know whether a typo is just two neighbouring
letters typed the wrong way round.

**Input.** Two lines: the word as typed, then the word meant.

**Output.** ~yes~ if swapping exactly one pair of **adjacent** characters in the
first produces the second, and ~no~ otherwise.

    Input   Output
    form    yes
    from

A word is not a transposition of itself: the swap has to change something.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    static char a[1005], b[1005];
    if (scanf("%1000s %1000s", a, b) != 2) return 1;
    /* print yes when one adjacent swap turns a into b */
    printf("no\n");
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b string
	fmt.Fscan(in, &a, &b)
	// print yes when one adjacent swap turns a into b
	fmt.Println("no")
}
`, `import sys

def main():
    parts = sys.stdin.read().split()
    a, b = parts[0], parts[1]
    # print yes when one adjacent swap turns a into b
    print("no")

main()
`),
		tests: []store.Check{
			tokens("Two neighbours swapped", "form\nfrom\n", "yes"),
			tokens("A word is not a swap of itself", "form\nform\n", "no"),
			tokens("Different lengths cannot be a swap", "form\nforms\n", "no"),
			tokens("Two changes that are not a swap", "cat\ndog\n", "no"),
			hid2("Non-adjacent letters swapped does not count", "abcd\ndbca\n", "no"),
			hid2("Swapping two identical letters changes nothing", "aab\naab\n", "no"),
			hid2("A swap at the very end", "teh\nthe\n", "yes"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char a[1005], b[1005];
    if (scanf("%1000s %1000s", a, b) != 2) return 1;
    int n = strlen(a);
    if ((int)strlen(b) != n || strcmp(a, b) == 0) { printf("no\n"); return 0; }
    int first = -1, second = -1, count = 0;
    for (int i = 0; i < n; i++) {
        if (a[i] != b[i]) {
            count++;
            if (first < 0) first = i; else second = i;
        }
    }
    int ok = count == 2 && second == first + 1 &&
             a[first] == b[second] && a[second] == b[first];
    printf("%s\n", ok ? "yes" : "no");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b string
	fmt.Fscan(in, &a, &b)
	if len(a) != len(b) || a == b {
		fmt.Println("no")
		return
	}
	var diff []int
	for i := range a {
		if a[i] != b[i] {
			diff = append(diff, i)
		}
	}
	ok := len(diff) == 2 && diff[1] == diff[0]+1 &&
		a[diff[0]] == b[diff[1]] && a[diff[1]] == b[diff[0]]
	if ok {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
`,
			"python": `import sys

def main():
    parts = sys.stdin.read().split()
    a, b = parts[0], parts[1]
    if len(a) != len(b) or a == b:
        print("no")
        return
    diff = [i for i in range(len(a)) if a[i] != b[i]]
    ok = (len(diff) == 2 and diff[1] == diff[0] + 1
          and a[diff[0]] == b[diff[1]] and a[diff[1]] == b[diff[0]])
    print("yes" if ok else "no")

main()
`},
	},
	{
		slug: "slug-registry", title: "Slug registry", difficulty: "hard", topic: "Strings",
		statement: `A blog turns every post title into a URL slug. The first post to claim a
slug keeps it; each later post claiming the same one gets ~-2~, then ~-3~, and
so on appended.

**Input.** A count ~n~, then ~n~ requested slugs, one per line.

**Output.** The slug each post ends up with, in order, one per line.

    Input     Output
    3         post
    post      post-2
    post      post-3
    post

The hard part is the collision you did not cause. If a post is *called* ~post-2~,
that slug is taken by the time the second ~post~ needs it — and whatever you hand
out has to be free too. No slug may ever be issued twice.`,
		timeLimitMs: 4000,
		starters:    triWords("issue each slug, appending the smallest free number on a collision"),
		tests: []store.Check{
			tokens("Numbers the repeats", "3\npost\npost\npost\n", "post post-2 post-3"),
			tokens("Distinct titles keep their slugs", "2\na\nb\n", "a b"),
			tokens("Nothing to assign", "0\n", ""),
			tokens("Steps over a slug somebody already took", "3\npost\npost-2\npost\n", "post post-2 post-3"),
			hid2("A taken slug pushes the next one further", "4\np\np-2\np-3\np\n", "p p-2 p-3 p-4"),
			hid2("A collision on the generated slug too", "3\np\np\np-2\n", "p p-2 p-2-2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

#define MAX 100000

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char used[MAX][64];
    static char counts[MAX][64];
    static int countVal[MAX];
    int usedN = 0, countsN = 0;
    for (int i = 0; i < n; i++) {
        char want[64], cand[64];
        if (scanf("%63s", want) != 1) return 1;
        int taken = 0;
        for (int j = 0; j < usedN; j++) if (strcmp(used[j], want) == 0) taken = 1;
        if (!taken) {
            strcpy(used[usedN++], want);
            strcpy(counts[countsN], want);
            countVal[countsN++] = 1;
            printf("%s\n", want);
            continue;
        }
        int slot = -1;
        for (int j = 0; j < countsN; j++) if (strcmp(counts[j], want) == 0) slot = j;
        if (slot < 0) { slot = countsN; strcpy(counts[slot], want); countVal[slot] = 1; countsN++; }
        int k = countVal[slot];
        for (;;) {
            k++;
            snprintf(cand, sizeof cand, "%s-%d", want, k);
            int clash = 0;
            for (int j = 0; j < usedN; j++) if (strcmp(used[j], cand) == 0) clash = 1;
            if (!clash) break;
        }
        countVal[slot] = k;
        strcpy(used[usedN++], cand);
        printf("%s\n", cand);
    }
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	used := map[string]bool{}
	counts := map[string]int{}
	for i := 0; i < n; i++ {
		var want string
		fmt.Fscan(in, &want)
		if !used[want] {
			used[want] = true
			counts[want] = 1
			fmt.Fprintln(out, want)
			continue
		}
		k := counts[want]
		if k == 0 {
			k = 1
		}
		var cand string
		for {
			k++
			cand = want + "-" + strconv.Itoa(k)
			if !used[cand] {
				break
			}
		}
		counts[want] = k
		used[cand] = true
		fmt.Fprintln(out, cand)
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    used, counts, out = set(), {}, []
    for want in data[1:1 + n]:
        if want not in used:
            used.add(want)
            counts[want] = 1
            out.append(want)
            continue
        k = counts.get(want, 1)
        while True:
            k += 1
            cand = want + "-" + str(k)
            if cand not in used:
                break
        counts[want] = k
        used.add(cand)
        out.append(cand)
    print("\n".join(out))

main()
`},
	},
	{
		slug: "letter-drift", title: "Letter drift", difficulty: "medium", topic: "Strings",
		statement: `A toy cipher shifts each letter further than the last: the first letter
moves 0 places along the alphabet, the second 1, the third 2, and so on, wrapping
from ~z~ to ~a~. Case is kept.

Characters that are not letters pass through untouched **and do not advance the
shift** — the count follows the letters, not the positions.

**Input.** One line.

**Output.** The drifted line.

    Input   Output
    abc     ace

    Input   Output
    a-bc    a-ce`,
		timeLimitMs: 4000,
		starters:    triLine("shift each letter by how many letters came before it"),
		tests: []store.Check{
			exact("Each letter drifts one further", "abc\n", "ace"),
			exact("Punctuation does not advance the shift", "a-bc\n", "a-ce"),
			exact("An empty line stays empty", "\n", ""),
			exact("Case is kept", "AB\n", "AC"),
			exactHid("Wraps round the end of the alphabet", "zz\n", "za"),
			exactHid("The first letter never moves", "q\n", "q"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int k = 0;
    for (int i = 0; line[i]; i++) {
        unsigned char c = line[i];
        if (isalpha(c)) {
            int base = isupper(c) ? 'A' : 'a';
            putchar((c - base + k) % 26 + base);
            k++;
        } else {
            putchar(c);
        }
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

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
	var b strings.Builder
	k := 0
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c >= 'a' && c <= 'z':
			b.WriteByte(byte((int(c-'a')+k)%26) + 'a')
			k++
		case c >= 'A' && c <= 'Z':
			b.WriteByte(byte((int(c-'A')+k)%26) + 'A')
			k++
		default:
			b.WriteByte(c)
		}
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    out = []
    k = 0
    for ch in line:
        if ch.isalpha():
            base = ord("A") if ch.isupper() else ord("a")
            out.append(chr((ord(ch) - base + k) % 26 + base))
            k += 1
        else:
            out.append(ch)
    print("".join(out))

main()
`},
	},
	{
		slug: "readable-list", title: "Readable list", difficulty: "easy", topic: "Strings",
		statement: `Join names the way a sentence would: commas between all but the last two,
and ~and~ before the last.

**Input.** A count ~n~, then ~n~ names, one per line.

**Output.** The joined phrase. An empty list has nothing to say, so it prints an
empty line.

    Input     Output
    3         Ada, Grace and Kay
    Ada
    Grace
    Kay`,
		timeLimitMs: 4000,
		starters:    triWords("join the names with commas and a final and"),
		tests: []store.Check{
			exact("One name stands alone", "1\nAda\n", "Ada"),
			exact("Two names take an and", "2\nAda\nGrace\n", "Ada and Grace"),
			exact("Three names take commas and an and", "3\nAda\nGrace\nKay\n", "Ada, Grace and Kay"),
			exact("Nobody says nothing", "0\n", ""),
			exactHid("A long list keeps one and", "4\na\nb\nc\nd\n", "a, b, c and d"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char name[100000][32];
    for (int i = 0; i < n; i++) if (scanf("%31s", name[i]) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (i > 0) printf(i == n - 1 ? " and " : ", ");
        printf("%s", name[i]);
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	names := make([]string, n)
	for i := range names {
		fmt.Fscan(in, &names[i])
	}
	switch {
	case n == 0:
		fmt.Println("")
	case n == 1:
		fmt.Println(names[0])
	default:
		fmt.Println(strings.Join(names[:n-1], ", ") + " and " + names[n-1])
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    names = data[1:1 + n]
    if not names:
        print("")
    elif len(names) == 1:
        print(names[0])
    else:
        print(", ".join(names[:-1]) + " and " + names[-1])

main()
`},
	},
	{
		slug: "column-widths", title: "Column widths", difficulty: "medium", topic: "Strings",
		statement: `Before a table can be printed, each column needs a width: the length of
the longest cell in it.

**Input.** A count ~n~, then ~n~ rows. Each row is a count ~k~ followed by ~k~
cells, written as words. A cell of ~-~ means an empty cell of width 0.

**Output.** One width per column, separated by spaces, counting the widest row.
No rows means no output.

    Input       Output
    2           4 2
    2 name id
    1 ada

Rows may be ragged — a row that stops early has nothing in the later columns and
contributes no width to them.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int width[1000] = {0}, columns = 0;
    for (int i = 0; i < n; i++) {
        int k;
        if (scanf("%d", &k) != 1) return 1;
        for (int j = 0; j < k; j++) {
            char cell[64];
            if (scanf("%63s", cell) != 1) return 1;
            int len = strcmp(cell, "-") == 0 ? 0 : (int)strlen(cell);
            /* widen column j */
        }
    }
    for (int j = 0; j < columns; j++) printf(j ? " %d" : "%d", width[j]);
    printf("\n");
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var width []int
	for i := 0; i < n; i++ {
		var k int
		fmt.Fscan(in, &k)
		for j := 0; j < k; j++ {
			var cell string
			fmt.Fscan(in, &cell)
			length := len(cell)
			if cell == "-" {
				length = 0
			}
			_ = length
			// widen column j
		}
	}
	for j, w := range width {
		if j > 0 {
			fmt.Print(" ")
		}
		fmt.Print(w)
	}
	fmt.Println()
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    pos = 1
    width = []
    for _ in range(int(data[0])):
        k = int(data[pos]); pos += 1
        for j in range(k):
            cell = data[pos + j]
            length = 0 if cell == "-" else len(cell)
            # widen column j
        pos += k
    print(" ".join(str(w) for w in width))

main()
`),
		tests: []store.Check{
			tokens("Takes the longest cell per column", "2\n2 name id\n1 ada\n", "4 2"),
			tokens("A single row is its own widths", "1\n2 abc d\n", "3 1"),
			tokens("No rows, no columns", "0\n", ""),
			hid2("A later row can be the wider one", "2\n1 a\n1 bbb\n", "3"),
			hid2("An empty cell has width zero", "1\n2 - xx\n", "0 2"),
			hid2("A ragged row adds columns", "2\n1 a\n2 b cc\n", "1 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int width[1000] = {0}, columns = 0;
    for (int i = 0; i < n; i++) {
        int k;
        if (scanf("%d", &k) != 1) return 1;
        if (k > columns) columns = k;
        for (int j = 0; j < k; j++) {
            char cell[64];
            if (scanf("%63s", cell) != 1) return 1;
            int len = strcmp(cell, "-") == 0 ? 0 : (int)strlen(cell);
            if (len > width[j]) width[j] = len;
        }
    }
    for (int j = 0; j < columns; j++) printf(j ? " %d" : "%d", width[j]);
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var width []int
	for i := 0; i < n; i++ {
		var k int
		fmt.Fscan(in, &k)
		for len(width) < k {
			width = append(width, 0)
		}
		for j := 0; j < k; j++ {
			var cell string
			fmt.Fscan(in, &cell)
			length := len(cell)
			if cell == "-" {
				length = 0
			}
			if length > width[j] {
				width[j] = length
			}
		}
	}
	for j, w := range width {
		if j > 0 {
			fmt.Print(" ")
		}
		fmt.Print(w)
	}
	fmt.Println()
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    pos = 1
    width = []
    for _ in range(int(data[0])):
        k = int(data[pos]); pos += 1
        while len(width) < k:
            width.append(0)
        for j in range(k):
            cell = data[pos + j]
            length = 0 if cell == "-" else len(cell)
            if length > width[j]:
                width[j] = length
        pos += k
    print(" ".join(str(w) for w in width))

main()
`},
	},
	{
		slug: "keysmash-filter", title: "Keysmash filter", difficulty: "medium", topic: "Strings",
		statement: `A signup form rejects names that look like somebody leaned on the
keyboard. The rule, exactly:

* four or more consonants in a row, or
* the same character three or more times in a row

Anything else is accepted. Comparisons ignore case, and the vowels are
~a e i o u~.

**Input.** One line.

**Output.** ~yes~ if the text breaks either rule, ~no~ otherwise.

    Input     Output
    asdfgh    yes

    Input     Output
    hello     no

The rule is the specification, not a judgement: ~rhythms~ trips it, and that is
the correct answer.`,
		timeLimitMs: 4000,
		starters:    triLine("check both rules and print yes or no"),
		tests: []store.Check{
			tokens("Four consonants in a row", "asdfgh\n", "yes"),
			tokens("An ordinary word passes", "hello\n", "no"),
			tokens("Three of the same character", "aaa\n", "yes"),
			tokens("An empty line is fine", "\n", "no"),
			hid2("Case does not matter", "ASDFGH\n", "yes"),
			hid2("Three consonants are not enough", "string\n", "no"),
			hid2("A repeated non-letter counts too", "a!!!b\n", "yes"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int run = 0, rep = 1, bad = 0;
    for (int i = 0; line[i]; i++) {
        char c = tolower((unsigned char)line[i]);
        if (isalpha((unsigned char)c) && !strchr("aeiou", c)) {
            if (++run >= 4) bad = 1;
        } else {
            run = 0;
        }
        if (i > 0) {
            rep = (c == tolower((unsigned char)line[i - 1])) ? rep + 1 : 1;
            if (rep >= 3) bad = 1;
        }
    }
    printf("%s\n", bad ? "yes" : "no");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	low := strings.ToLower(strings.TrimRight(line, "\r\n"))
	run, rep, bad := 0, 1, false
	for i := 0; i < len(low); i++ {
		c := low[i]
		isLetter := c >= 'a' && c <= 'z'
		if isLetter && !strings.ContainsRune("aeiou", rune(c)) {
			run++
			if run >= 4 {
				bad = true
			}
		} else {
			run = 0
		}
		if i > 0 {
			if c == low[i-1] {
				rep++
			} else {
				rep = 1
			}
			if rep >= 3 {
				bad = true
			}
		}
	}
	if bad {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
`,
			"python": `import sys

def main():
    low = sys.stdin.readline().rstrip("\n").lower()
    run, rep, bad = 0, 1, False
    for i, ch in enumerate(low):
        if ch.isalpha() and ch not in "aeiou":
            run += 1
            if run >= 4:
                bad = True
        else:
            run = 0
        if i > 0:
            rep = rep + 1 if ch == low[i - 1] else 1
            if rep >= 3:
                bad = True
    print("yes" if bad else "no")

main()
`},
	},
	{
		slug: "headline-case", title: "Headline case", difficulty: "medium", topic: "Strings",
		statement: `A headline capitalises every word except the small ones — ~a an the of and
in on to for~ — which stay lower case. The first and last word are always
capitalised whatever they are.

Every word is otherwise normalised: first letter up, the rest down.

**Input.** One line.

**Output.** The headline.

    Input                   Output
    the lord of the rings   The Lord of the Rings

    Input                   Output
    gone WITH the wind      Gone With the Wind`,
		timeLimitMs: 4000,
		starters:    triLine("capitalise every word but the small ones in the middle"),
		tests: []store.Check{
			exact("Small words stay small", "the lord of the rings\n", "The Lord of the Rings"),
			exact("Shouting is normalised", "gone WITH the wind\n", "Gone With the Wind"),
			exact("One word is always capitalised", "of\n", "Of"),
			exactHid("A small word at the end is capitalised", "what are you waiting for\n", "What Are You Waiting For"),
			exactHid("A small word at the start is capitalised", "a tale of two cities\n", "A Tale of Two Cities"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

static int small(const char *w) {
    static const char *s[] = {"a", "an", "the", "of", "and", "in", "on", "to", "for"};
    for (int i = 0; i < 9; i++) if (strcmp(w, s[i]) == 0) return 1;
    return 0;
}

int main(void) {
    static char line[100005], word[1000][64];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int count = 0;
    char *tok = strtok(line, " \t");
    while (tok && count < 1000) { snprintf(word[count++], 64, "%s", tok); tok = strtok(NULL, " \t"); }
    for (int i = 0; i < count; i++) {
        char low[64];
        snprintf(low, sizeof low, "%s", word[i]);
        for (int j = 0; low[j]; j++) low[j] = tolower((unsigned char)low[j]);
        if (i > 0 && i < count - 1 && small(low)) {
            printf("%s%s", i ? " " : "", low);
        } else {
            low[0] = toupper((unsigned char)low[0]);
            printf("%s%s", i ? " " : "", low);
        }
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var small = map[string]bool{"a": true, "an": true, "the": true, "of": true,
	"and": true, "in": true, "on": true, "to": true, "for": true}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	words := strings.Fields(strings.TrimRight(line, "\r\n"))
	for i, w := range words {
		low := strings.ToLower(w)
		if i > 0 && i < len(words)-1 && small[low] {
			words[i] = low
		} else {
			words[i] = strings.ToUpper(low[:1]) + low[1:]
		}
	}
	fmt.Println(strings.Join(words, " "))
}
`,
			"python": `import sys

SMALL = {"a", "an", "the", "of", "and", "in", "on", "to", "for"}

def main():
    words = sys.stdin.readline().rstrip("\n").split()
    out = []
    for i, w in enumerate(words):
        low = w.lower()
        if 0 < i < len(words) - 1 and low in SMALL:
            out.append(low)
        else:
            out.append(low[:1].upper() + low[1:])
    print(" ".join(out))

main()
`},
	},
}
