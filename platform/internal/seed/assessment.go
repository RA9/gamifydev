package seed

import (
	"context"

	"github.com/RA9/gamifydev/platform/internal/placement"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// The placement diagnostic's item bank.
//
// Sizing: the bank holds six items per topic and a sitting draws four, so no two
// learners reliably see the same paper — the practical mitigation against a
// leaked bank (PRD §03). Rotating or extending the bank is just editing this
// slice; a learner's finished attempt keeps its own record of what was asked.
//
// Every item is a single-answer multiple choice. `Answer` is the 0-based index
// into `Options` and never reaches the browser.

const placementSlug = "placement"

var placementItems = []store.Item{
	// ---------------------------------------------------------------- programming
	{
		Topic: placement.TopicProgramming, Difficulty: 1,
		Prompt: "What does a function's `return` statement do?",
		Options: []string{
			"Prints the value to the screen",
			"Sends a value back to whatever called the function",
			"Stops the whole program",
			"Stores the value in a global variable",
		},
		Answer:      1,
		Explanation: "`return` hands a value back to the caller and ends the function. Printing is a separate action.",
	},
	{
		Topic: placement.TopicProgramming, Difficulty: 1,
		Prompt: "What does this print?",
		Code:   "total = 0\nfor n in [1, 2, 3, 4]:\n    total = total + n\nprint(total)",
		Lang:   "python",
		Options: []string{
			"4", "10", "24", "1234",
		},
		Answer:      1,
		Explanation: "The loop adds each number to the running total: 1 + 2 + 3 + 4 = 10.",
	},
	{
		Topic: placement.TopicProgramming, Difficulty: 2,
		Prompt: "How many times does this loop run?",
		Code:   "for i in range(0, 10, 2):\n    print(i)",
		Lang:   "python",
		Options: []string{
			"2", "5", "10", "20",
		},
		Answer:      1,
		Explanation: "`range(0, 10, 2)` yields 0, 2, 4, 6, 8 — five values. The stop value is excluded.",
	},
	{
		Topic: placement.TopicProgramming, Difficulty: 2,
		Prompt: "What does this print?",
		Code:   "def add(a, b=10):\n    return a + b\n\nprint(add(5))",
		Lang:   "python",
		Options: []string{
			"5", "10", "15", "An error — b is missing",
		},
		Answer:      2,
		Explanation: "`b` has a default of 10, so calling `add(5)` uses it: 5 + 10 = 15.",
	},
	{
		Topic: placement.TopicProgramming, Difficulty: 2,
		Prompt: "A recursive function that never reaches its base case will:",
		Options: []string{
			"Return zero",
			"Run forever until the call stack overflows",
			"Skip the recursive call automatically",
			"Be rejected by the compiler",
		},
		Answer:      1,
		Explanation: "Each call adds a stack frame. Without a base case the recursion keeps going until the stack is exhausted.",
	},
	{
		Topic: placement.TopicProgramming, Difficulty: 3,
		Prompt: "What does this print?",
		Code:   "def f(xs):\n    xs.append(4)\n\nnums = [1, 2, 3]\nf(nums)\nprint(len(nums))",
		Lang:   "python",
		Options: []string{
			"3", "4", "0", "An error — nums is unchanged inside f",
		},
		Answer:      1,
		Explanation: "Lists are passed by reference, so `f` mutates the caller's list. It now holds four items.",
	},

	// ------------------------------------------------------------ data structures
	{
		Topic: placement.TopicDataStructures, Difficulty: 1,
		Prompt: "Which structure gives you constant-time access to an element by its index?",
		Options: []string{
			"Linked list", "Array", "Queue", "Binary search tree",
		},
		Answer:      1,
		Explanation: "Array elements sit in contiguous memory, so the address of index i is a single calculation away.",
	},
	{
		Topic: placement.TopicDataStructures, Difficulty: 1,
		Prompt: "A stack removes items in which order?",
		Options: []string{
			"First in, first out",
			"Last in, first out",
			"Smallest value first",
			"Random order",
		},
		Answer:      1,
		Explanation: "A stack is LIFO — the most recently pushed item is the first one popped. It's what backs the call stack.",
	},
	{
		Topic: placement.TopicDataStructures, Difficulty: 2,
		Prompt: "What is the *average* time to look up a key in a hash table?",
		Options: []string{
			"O(1)", "O(log n)", "O(n)", "O(n log n)",
		},
		Answer:      0,
		Explanation: "Hashing jumps straight to a bucket, so average lookup is constant. The worst case degrades to O(n) when everything collides.",
	},
	{
		Topic: placement.TopicDataStructures, Difficulty: 2,
		Prompt: "In a binary search tree, where does a value smaller than the current node belong?",
		Options: []string{
			"In the right subtree",
			"In the left subtree",
			"At the root",
			"Anywhere — order doesn't matter",
		},
		Answer:      1,
		Explanation: "The BST invariant is left < node < right. That ordering is what makes searching possible.",
	},
	{
		Topic: placement.TopicDataStructures, Difficulty: 2,
		Prompt: "Why can inserting into the middle of an array be slow?",
		Options: []string{
			"The array has to be re-sorted",
			"Every later element must shift by one position",
			"Arrays cannot be modified after creation",
			"The hash has to be recomputed",
		},
		Answer:      1,
		Explanation: "Contiguous storage means making room in the middle requires shifting everything after it — O(n).",
	},
	{
		Topic: placement.TopicDataStructures, Difficulty: 3,
		Prompt: "A binary search tree that receives values in already-sorted order degenerates into:",
		Options: []string{
			"A balanced tree with O(log n) search",
			"A structure shaped like a linked list, with O(n) search",
			"A hash table",
			"A heap",
		},
		Answer:      1,
		Explanation: "Each new value goes to the same side, so the tree becomes a chain and search degrades to a linear scan. This is why balanced trees exist.",
	},

	// ---------------------------------------------------------------- complexity
	{
		Topic: placement.TopicComplexity, Difficulty: 1,
		Prompt: "What does Big O notation describe?",
		Options: []string{
			"The exact running time in seconds",
			"How the running time grows as the input grows",
			"How much memory a language uses",
			"The number of lines of code",
		},
		Answer:      1,
		Explanation: "Big O is about growth rate as input size increases, not absolute time on a particular machine.",
	},
	{
		Topic: placement.TopicComplexity, Difficulty: 1,
		Prompt: "What is the time complexity of this code?",
		Code:   "for i in range(n):\n    for j in range(n):\n        print(i, j)",
		Lang:   "python",
		Options: []string{
			"O(1)", "O(n)", "O(n^2)", "O(log n)",
		},
		Answer:      2,
		Explanation: "Nested loops each running n times multiply: n * n = O(n^2).",
	},
	{
		Topic: placement.TopicComplexity, Difficulty: 2,
		Prompt:      "Simplified to Big O, `3n^2 + 500n + 9000` is:",
		Options:     []string{"O(n)", "O(n^2)", "O(3n^2)", "O(9000)"},
		Answer:      1,
		Explanation: "Constants and lower-order terms drop; only the fastest-growing term matters as n grows.",
	},
	{
		Topic: placement.TopicComplexity, Difficulty: 2,
		Prompt:      "Which of these grows the fastest as n increases?",
		Options:     []string{"O(n log n)", "O(n^2)", "O(2^n)", "O(n)"},
		Answer:      2,
		Explanation: "Exponential growth outruns any polynomial. At n = 60, 2^n is already astronomically large.",
	},
	{
		Topic: placement.TopicComplexity, Difficulty: 2,
		Prompt: "An algorithm that halves its search range on every step runs in:",
		Options: []string{
			"O(1)", "O(log n)", "O(n)", "O(n^2)",
		},
		Answer:      1,
		Explanation: "Repeated halving reaches one element in log2(n) steps — the signature of logarithmic time.",
	},
	{
		Topic: placement.TopicComplexity, Difficulty: 3,
		Prompt: "NP is the set of problems where:",
		Options: []string{
			"No polynomial-time algorithm can possibly exist",
			"A proposed solution can be verified in polynomial time",
			"The algorithm is non-deterministic in practice",
			"The running time is always exponential",
		},
		Answer:      1,
		Explanation: "NP is about fast *verification*, not slow solving. Every problem in P is also in NP.",
	},

	// ---------------------------------------------------------------- algorithms
	{
		Topic: placement.TopicAlgorithms, Difficulty: 1,
		Prompt: "What must be true of a list before you can binary search it?",
		Options: []string{
			"It must be sorted",
			"It must contain only numbers",
			"It must have an even length",
			"It must be stored in a linked list",
		},
		Answer:      0,
		Explanation: "Binary search decides which half to discard by comparing with the middle, which only works on ordered data.",
	},
	{
		Topic: placement.TopicAlgorithms, Difficulty: 2,
		Prompt:      "Which traversal explores a graph level by level, using a queue?",
		Options:     []string{"Depth-first search", "Breadth-first search", "Binary search", "Quick sort"},
		Answer:      1,
		Explanation: "BFS uses a queue to visit all neighbours at the current distance before going deeper. DFS uses a stack.",
	},
	{
		Topic: placement.TopicAlgorithms, Difficulty: 2,
		Prompt:      "What is the worst-case time complexity of quicksort?",
		Options:     []string{"O(n)", "O(n log n)", "O(n^2)", "O(2^n)"},
		Answer:      2,
		Explanation: "A consistently bad pivot (e.g. already-sorted input with a naive pivot) gives O(n^2). Its *average* is O(n log n).",
	},
	{
		Topic: placement.TopicAlgorithms, Difficulty: 2,
		Prompt:      "Merge sort's worst-case time complexity is:",
		Options:     []string{"O(n)", "O(n log n)", "O(n^2)", "It varies with the input"},
		Answer:      1,
		Explanation: "Merge sort always splits in half and merges linearly, so it is O(n log n) in every case — at the cost of O(n) extra space.",
	},
	{
		Topic: placement.TopicAlgorithms, Difficulty: 3,
		Prompt: "Dijkstra's shortest-path algorithm fails when the graph has:",
		Options: []string{
			"More than 1000 nodes",
			"Negative edge weights",
			"Any cycles at all",
			"Unweighted edges",
		},
		Answer:      1,
		Explanation: "Dijkstra finalises a node once reached, assuming no later path can be shorter. A negative edge breaks that assumption; Bellman-Ford handles it.",
	},
	{
		Topic: placement.TopicAlgorithms, Difficulty: 3,
		Prompt: "A greedy algorithm takes the best local option at each step. This:",
		Options: []string{
			"Always produces the optimal answer",
			"Produces the optimal answer only for some problems",
			"Is the same as backtracking",
			"Never produces the optimal answer",
		},
		Answer:      1,
		Explanation: "Greedy is provably optimal for some problems (Dijkstra, Kruskal) and wrong for others — coin change with {1,3,4} making 6 is the classic failure.",
	},

	// ----------------------------------------------------------------- computers
	{
		Topic: placement.TopicComputers, Difficulty: 1,
		Prompt:      "What is the decimal value of the binary number `1011`?",
		Options:     []string{"7", "9", "11", "13"},
		Answer:      2,
		Explanation: "8 + 0 + 2 + 1 = 11. The place values from the left are 8, 4, 2, 1.",
	},
	{
		Topic: placement.TopicComputers, Difficulty: 1,
		Prompt: "Order these from fastest to slowest to access:",
		Options: []string{
			"CPU register, RAM, CPU cache, disk",
			"CPU register, CPU cache, RAM, disk",
			"CPU cache, CPU register, disk, RAM",
			"RAM, CPU register, CPU cache, disk",
		},
		Answer:      1,
		Explanation: "The memory hierarchy runs register → cache → RAM → disk, each roughly an order of magnitude slower and larger than the last.",
	},
	{
		Topic: placement.TopicComputers, Difficulty: 2,
		Prompt:      "How many distinct values can a single byte (8 bits) represent?",
		Options:     []string{"8", "64", "128", "256"},
		Answer:      3,
		Explanation: "Each bit doubles the possibilities: 2^8 = 256 distinct values.",
	},
	{
		Topic: placement.TopicComputers, Difficulty: 2,
		Prompt:      "What is the result of the bitwise expression `5 & 3`?",
		Options:     []string{"1", "3", "7", "8"},
		Answer:      0,
		Explanation: "0101 AND 0011 compares column by column and keeps only bits set in both: 0001 = 1.",
	},
	{
		Topic: placement.TopicComputers, Difficulty: 2,
		Prompt: "Why does `0.1 + 0.2` not equal exactly `0.3` in most languages?",
		Options: []string{
			"A bug in the language's maths library",
			"0.1 and 0.2 have no exact representation in binary floating point",
			"The result is rounded to two decimal places",
			"Addition is not associative",
		},
		Answer:      1,
		Explanation: "Like 1/3 in decimal, 0.1 is a repeating fraction in binary. IEEE 754 stores the nearest representable value, so tiny errors appear.",
	},
	{
		Topic: placement.TopicComputers, Difficulty: 3,
		Prompt: "What does UTF-8 provide that ASCII does not?",
		Options: []string{
			"Faster string comparison",
			"Representation of characters beyond the original 128, using 1–4 bytes",
			"Automatic compression of text",
			"Case-insensitive sorting",
		},
		Answer:      1,
		Explanation: "UTF-8 is a variable-width encoding covering all Unicode code points while staying byte-compatible with ASCII for the first 128.",
	},

	// ----------------------------------------------------------------------- web
	{
		Topic: placement.TopicWeb, Difficulty: 1,
		Prompt:      "What does HTML provide for a web page?",
		Options:     []string{"Its structure and content", "Its colours and layout", "Its server logic", "Its database schema"},
		Answer:      0,
		Explanation: "HTML marks up structure and meaning; CSS handles presentation and JavaScript handles behaviour.",
	},
	{
		Topic: placement.TopicWeb, Difficulty: 1,
		Prompt: "What does the DOM represent?",
		Options: []string{
			"The CSS rules of a page",
			"A tree of the page's elements that scripts can read and change",
			"The network requests a page makes",
			"The server's file system",
		},
		Answer:      1,
		Explanation: "The Document Object Model is the live tree of nodes the browser builds from HTML, which JavaScript can query and modify.",
	},
	{
		Topic: placement.TopicWeb, Difficulty: 2,
		Prompt: "In JavaScript, how does `===` differ from `==`?",
		Options: []string{
			"`===` is faster but otherwise identical",
			"`===` compares without converting types, `==` converts first",
			"`===` only works on numbers",
			"There is no difference",
		},
		Answer:      1,
		Explanation: "`==` performs type coercion (so `\"1\" == 1` is true) while `===` requires the types to match as well as the values.",
	},
	{
		Topic: placement.TopicWeb, Difficulty: 2,
		Prompt:      "Which HTTP status code means the requested resource was not found?",
		Options:     []string{"200", "301", "404", "500"},
		Answer:      2,
		Explanation: "404 is the classic not-found response. 200 is success, 301 a permanent redirect, 500 a server error.",
	},
	{
		Topic: placement.TopicWeb, Difficulty: 2,
		Prompt:      "In CSS, which selector has the highest specificity?",
		Options:     []string{"A tag selector like `p`", "A class selector like `.intro`", "An ID selector like `#header`", "The universal selector `*`"},
		Answer:      2,
		Explanation: "Specificity climbs from universal → tag → class → ID, so an ID selector wins against the others.",
	},
	{
		Topic: placement.TopicWeb, Difficulty: 3,
		Prompt: "What does `display: flex` on a container do?",
		Options: []string{
			"Hides the container",
			"Lays its direct children out along one axis, with control over alignment and spacing",
			"Makes the container scrollable",
			"Applies a CSS animation",
		},
		Answer:      1,
		Explanation: "Flexbox is a one-dimensional layout model: children become flex items arranged along the main axis.",
	},

	// ------------------------------------------------------------------- backend
	{
		Topic: placement.TopicBackend, Difficulty: 1,
		Prompt:      "Which HTTP method is conventionally used to create a new resource?",
		Options:     []string{"GET", "POST", "DELETE", "HEAD"},
		Answer:      1,
		Explanation: "POST submits data to create something. GET should be safe and read-only.",
	},
	{
		Topic: placement.TopicBackend, Difficulty: 1,
		Prompt: "What does a database primary key guarantee?",
		Options: []string{
			"The column is sorted",
			"Each row has a unique, non-null identifier",
			"The column stores only numbers",
			"The table cannot be modified",
		},
		Answer:      1,
		Explanation: "A primary key uniquely identifies each row and cannot be null, which is what lets other tables reference it.",
	},
	{
		Topic: placement.TopicBackend, Difficulty: 2,
		Prompt: "What does this SQL return?",
		Code:   "SELECT name FROM users WHERE age >= 18;",
		Lang:   "sql",
		Options: []string{
			"Every column for users aged 18 or over",
			"The name column for users aged 18 or over",
			"The number of users aged 18 or over",
			"All users, sorted by age",
		},
		Answer:      1,
		Explanation: "The SELECT list picks the columns and the WHERE clause filters the rows.",
	},
	{
		Topic: placement.TopicBackend, Difficulty: 2,
		Prompt: "What is a database index for?",
		Options: []string{
			"Storing a backup copy of the table",
			"Speeding up lookups, at the cost of slower writes and more space",
			"Enforcing that a column is not null",
			"Compressing the table on disk",
		},
		Answer:      1,
		Explanation: "An index is a secondary structure that makes matching rows fast to find, but must be updated on every write.",
	},
	{
		Topic: placement.TopicBackend, Difficulty: 2,
		Prompt: "An HTTP 500 response means:",
		Options: []string{
			"The client sent a malformed request",
			"The server hit an error while handling the request",
			"The resource has moved",
			"Authentication is required",
		},
		Answer:      1,
		Explanation: "5xx codes are server-side failures. 4xx codes indicate a problem with the client's request.",
	},
	{
		Topic: placement.TopicBackend, Difficulty: 3,
		Prompt: "Why should a web server never store a user's password directly?",
		Options: []string{
			"Passwords take up too much space",
			"A database breach would expose every user's actual password",
			"Databases cannot store long strings",
			"It makes login slower",
		},
		Answer:      1,
		Explanation: "Passwords are stored as salted hashes so that a leaked database does not hand an attacker usable credentials.",
	},
}

// seedPlacement installs (or refreshes) the diagnostic and its item bank.
func seedPlacement(ctx context.Context, st *store.Store) (int, error) {
	id, err := st.UpsertAssessment(ctx, store.Assessment{
		Slug:  placementSlug,
		Title: "Placement diagnostic",
		Kind:  "placement",
		// 28 questions at roughly a minute each, with room to think.
		TimeLimitS: 45 * 60,
		PerTopic:   4,
	})
	if err != nil {
		return 0, err
	}
	if err := st.ReplaceItems(ctx, id, placementItems); err != nil {
		return 0, err
	}
	return len(placementItems), nil
}
