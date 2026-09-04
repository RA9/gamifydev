package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// basicProblems: loops, conditionals and arithmetic. The on-ramp.
//
// Even here the problems are invented rather than borrowed. A beginner who
// searches their way past "sum this list" has learned that searching works,
// which is the opposite of what an on-ramp is for.
var basicProblems = []seedProblem{
	{
		slug: "double-tiles", title: "Double tiles", difficulty: "easy", topic: "Loops",
		statement: `A board game scores each tile you land on, and some tiles are gold —
a gold tile is worth twice its face value.

Write ~score(values, gold)~, where ~gold~ is a list of booleans the same length
as ~values~, saying which tiles were gold. Return the total.

    score([3, 5, 2], [False, True, False])  ->  15

An empty board scores 0.`,
		timeLimitMs: 5000,
		starters:    py("def score(values, gold):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Doubles the gold tiles", "score([3, 5, 2], [False, True, False]) == 15"),
			vis("No gold means the plain total", "score([1, 2, 3], [False, False, False]) == 6"),
			vis("An empty board scores nothing", "score([], []) == 0"),
			hid("All gold doubles everything", "score([4, 6], [True, True]) == 20"),
			hid("Handles negative tiles", "score([-3, 4], [True, False]) == -2"),
		},
		solution: "def score(values, gold):\n    return sum(v * 2 if g else v for v, g in zip(values, gold))\n",
	},
	{
		slug: "tide-crossings", title: "Tide crossings", difficulty: "easy", topic: "Loops",
		statement: `A sensor records the water level once an hour. A crossing is any hour
where the level was strictly below the flood mark and the next reading is at or
above it.

Write ~crossings(levels, mark)~ returning how many crossings there were.

    crossings([1, 2, 5, 4, 6], 5)  ->  2

Note what this does not count: staying above the mark isn't a crossing, and
falling back below it isn't either. Only the moment it goes up through.`,
		timeLimitMs: 5000,
		starters:    py("def crossings(levels, mark):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts each rise through the mark", "crossings([1, 2, 5, 4, 6], 5) == 2"),
			vis("Staying above is not a crossing", "crossings([6, 7, 8], 5) == 0"),
			vis("A single reading cannot cross", "crossings([9], 5) == 0"),
			vis("No readings, no crossings", "crossings([], 5) == 0"),
			hid("Landing exactly on the mark counts", "crossings([4, 5], 5) == 1"),
			hid("Falling back below does not count", "crossings([6, 1, 6], 5) == 1"),
		},
		solution: "def crossings(levels, mark):\n    return sum(1 for a, b in zip(levels, levels[1:]) if a < mark <= b)\n",
	},
	{
		slug: "quiet-hours", title: "Quiet hours", difficulty: "easy", topic: "Loops",
		statement: `A noise monitor logs a reading each hour. An hour counts as restful only
if it belongs to a stretch of **three or more** consecutive hours that are all
at or below the limit — a single quiet hour between two loud ones is not rest.

Write ~restful(readings, limit)~ returning how many hours qualify.

    restful([1, 1, 1, 9, 2, 2], 3)  ->  3

The first three hours form a stretch of three, so all three count. The last two
are quiet but only two long, so neither does.`,
		timeLimitMs: 5000,
		starters:    py("def restful(readings, limit):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts a stretch of three", "restful([1, 1, 1, 9, 2, 2], 3) == 3"),
			vis("A short stretch counts for nothing", "restful([1, 1, 9, 1, 1], 3) == 0"),
			vis("An empty log has no restful hours", "restful([], 3) == 0"),
			vis("A long stretch counts in full", "restful([0, 0, 0, 0, 0], 3) == 5"),
			hid("A stretch running to the end still counts", "restful([9, 1, 1, 1], 3) == 3"),
			hid("Two separate stretches both count", "restful([1, 1, 1, 9, 1, 1, 1], 3) == 6"),
		},
		solution: "def restful(readings, limit):\n    total = 0\n    run = 0\n    for r in list(readings) + [limit + 1]:\n        if r <= limit:\n            run += 1\n        else:\n            if run >= 3:\n                total += run\n            run = 0\n    return total\n",
	},
	{
		slug: "mind-changes", title: "Mind changes", difficulty: "easy", topic: "Loops",
		statement: `A committee votes over and over until it settles. Write
~changes(votes)~ returning how many times the vote differed from the one before
it.

    changes(["yes", "yes", "no", "yes"])  ->  2

A list with fewer than two votes has no changes.`,
		timeLimitMs: 5000,
		starters:    py("def changes(votes):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts each switch", `changes(["yes", "yes", "no", "yes"]) == 2`),
			vis("A unanimous run never changes", `changes(["no", "no", "no"]) == 0`),
			vis("One vote cannot change", `changes(["yes"]) == 0`),
			vis("No votes at all", "changes([]) == 0"),
			hid("Works for numbers as well as strings", "changes([1, 2, 2, 3]) == 2"),
		},
		solution: "def changes(votes):\n    return sum(1 for a, b in zip(votes, votes[1:]) if a != b)\n",
	},
	{
		slug: "stamp-card", title: "Stamp card", difficulty: "easy", topic: "Arithmetic",
		statement: `A coffee shop stamps your card for every drink. Every fifth drink is
free — the 5th, the 10th, the 15th, and so on, counting from the start of the
list.

Write ~bill(prices)~ returning what you actually pay.

    bill([3, 3, 3, 3, 3])  ->  12

The fifth drink is free whatever it costs, so ordering the expensive one fifth
is a strategy. That is not your problem to solve, only to charge for.`,
		timeLimitMs: 5000,
		starters:    py("def bill(prices):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The fifth drink is free", "bill([3, 3, 3, 3, 3]) == 12"),
			vis("Fewer than five drinks are all charged", "bill([2, 4]) == 6"),
			vis("An empty order costs nothing", "bill([]) == 0"),
			vis("The tenth is free too", "bill([1] * 10) == 8"),
			hid("The free drink is the fifth, whatever it costs", "bill([1, 1, 1, 1, 99]) == 4"),
		},
		solution: "def bill(prices):\n    return sum(p for i, p in enumerate(prices, 1) if i % 5)\n",
	},
	{
		slug: "fare-cap", title: "Fare cap", difficulty: "easy", topic: "Arithmetic",
		statement: `A tram charges 2 for any trip, plus 1 for every kilometre beyond the
third. No single trip costs more than 8, and no day costs more than 20 however
much you ride.

Write ~fare(distances)~, taking a list of whole-kilometre trips, and returning
what the day costs.

    fare([2, 5, 20])  ->  14

That is 2 + 4 + 8: the last trip hits the per-trip cap.`,
		timeLimitMs: 5000,
		starters:    py("def fare(distances):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Charges the base, the extra, and the trip cap", "fare([2, 5, 20]) == 14"),
			vis("A short trip is just the base fare", "fare([1]) == 2"),
			vis("No trips, no fare", "fare([]) == 0"),
			vis("The daily cap holds", "fare([20] * 10) == 20"),
			hid("Exactly three kilometres adds nothing", "fare([3]) == 2"),
			hid("The daily cap applies to the total, not each trip", "fare([8, 8, 8]) == 20"),
		},
		solution: "def fare(distances):\n    return min(20, sum(min(8, 2 + max(0, d - 3)) for d in distances))\n",
	},
	{
		slug: "battery-log", title: "Battery log", difficulty: "easy", topic: "Simulation",
		statement: `A phone starts the day at 50%. The log is a list of whole-number
changes: negative for drain, positive for charge. The battery never goes below
0 or above 100 — an event that would push it past either end just leaves it
there.

Write ~battery(events)~ returning the level at the end of the day.

    battery([-30, -40, 60])  ->  60

The second event would take it to -20, so it stops at 0, and the charge lifts it
from there.`,
		timeLimitMs: 5000,
		starters:    py("def battery(events):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Clamps at empty before charging again", "battery([-30, -40, 60]) == 60"),
			vis("A quiet day changes nothing", "battery([]) == 50"),
			vis("Cannot charge past full", "battery([80, 20]) == 100"),
			hid("Cannot drain past empty", "battery([-90, -90]) == 0"),
			hid("Clamping happens per event, not at the end", "battery([-90, 10, 10]) == 20"),
		},
		solution: "def battery(events):\n    level = 50\n    for e in events:\n        level = max(0, min(100, level + e))\n    return level\n",
	},
	{
		slug: "shift-pay", title: "Shift pay", difficulty: "easy", topic: "Arithmetic",
		statement: `A workshop pays 10 an hour for the first 40 hours of the week, 20 an
hour for the next 10, and 30 an hour beyond 50.

Write ~pay(hours)~ returning the week's wage.

    pay(45)  ->  500
    pay(52)  ->  660

Hours are whole numbers and never negative.`,
		timeLimitMs: 5000,
		starters:    py("def pay(hours):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Pays the middle band", "pay(45) == 500"),
			vis("Pays the top band", "pay(52) == 660"),
			vis("A short week is all base rate", "pay(10) == 100"),
			vis("No hours, no pay", "pay(0) == 0"),
			hid("Exactly 40 hours stays on the base rate", "pay(40) == 400"),
			hid("Exactly 50 hours fills the middle band", "pay(50) == 600"),
		},
		solution: "def pay(hours):\n    return 10 * min(hours, 40) + 20 * max(0, min(hours, 50) - 40) + 30 * max(0, hours - 50)\n",
	},
	{
		slug: "matching-dice", title: "Matching dice", difficulty: "easy", topic: "Loops",
		statement: `A game is played with two dice. Write ~longest_match(rolls)~, taking a
list of ~(a, b)~ pairs, and returning the length of the longest unbroken run of
rolls where both dice showed the same number.

    longest_match([(1, 1), (2, 2), (3, 4), (5, 5)])  ->  2

No rolls means no run.`,
		timeLimitMs: 5000,
		starters:    py("def longest_match(rolls):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the longer of two runs", "longest_match([(1, 1), (2, 2), (3, 4), (5, 5)]) == 2"),
			vis("No matches at all", "longest_match([(1, 2), (3, 4)]) == 0"),
			vis("An empty game has no run", "longest_match([]) == 0"),
			hid("A run at the very end counts", "longest_match([(1, 2), (6, 6), (6, 6), (6, 6)]) == 3"),
			hid("Every roll matching is one long run", "longest_match([(4, 4)] * 5) == 5"),
		},
		solution: "def longest_match(rolls):\n    best = run = 0\n    for a, b in rolls:\n        run = run + 1 if a == b else 0\n        if run > best:\n            best = run\n    return best\n",
	},
	{
		slug: "row-of-seedlings", title: "Row of seedlings", difficulty: "easy", topic: "Arithmetic",
		statement: `You plant the first seedling at position 0 of a bed, then one every
~gap~ centimetres, for as long as the next one still fits inside a bed
~length~ centimetres long. A seedling exactly at the far end fits.

Write ~plants(length, gap)~ returning how many you plant.

    plants(10, 3)  ->  4

Those sit at 0, 3, 6 and 9. A bed of negative length holds nothing.`,
		timeLimitMs: 5000,
		starters:    py("def plants(length, gap):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Fits four in a ten-centimetre bed", "plants(10, 3) == 4"),
			vis("A bed with no room still holds the first", "plants(0, 3) == 1"),
			vis("A negative bed holds nothing", "plants(-1, 3) == 0"),
			hid("A seedling exactly at the end fits", "plants(9, 3) == 4"),
			hid("A huge bed does not need a loop", "plants(10 ** 12, 7) == 142857142858"),
		},
		solution: "def plants(length, gap):\n    if length < 0:\n        return 0\n    return length // gap + 1\n",
	},
	{
		slug: "working-days-ahead", title: "Working days ahead", difficulty: "medium", topic: "Arithmetic",
		statement: `A workshop is open every day except Sunday. Days of the week are
numbered 0 for Monday through 6 for Sunday.

Write ~advance(start, days)~: starting on day ~start~, move forward until you
have passed ~days~ **working** days, and return the day you land on. Sundays are
stepped over and never counted.

    advance(0, 6)  ->  0

Six working days from Monday is Tuesday, Wednesday, Thursday, Friday, Saturday,
then Monday again — Sunday is skipped.

~days~ can run into the billions, so stepping one day at a time will not finish
in time. Six working days is exactly one week.`,
		timeLimitMs: 4000,
		starters:    py("def advance(start, days):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Six working days is a whole week", "advance(0, 6) == 0"),
			vis("One day forward from Monday", "advance(0, 1) == 1"),
			vis("Sunday is stepped over", "advance(5, 1) == 0"),
			vis("Going nowhere lands where you started", "advance(3, 0) == 3"),
			vis("Fast enough for a billion days", "advance(0, 10 ** 9) == 4"),
			hid("Starting on a Sunday still works", "advance(6, 1) == 0"),
			hid("Five days from Monday reaches Saturday", "advance(0, 5) == 5"),
		},
		solution: "def advance(start, days):\n    d = (start + (days // 6) * 7) % 7\n    for _ in range(days % 6):\n        d = (d + 1) % 7\n        if d == 6:\n            d = 0\n    return d\n",
		tooSlow:  "def advance(start, days):\n    d = start\n    while days > 0:\n        d = (d + 1) % 7\n        if d != 6:\n            days -= 1\n    return d\n",
	},
	{
		slug: "biggest-swing", title: "Biggest swing", difficulty: "easy", topic: "Loops",
		statement: `Write ~swing(readings)~ returning the largest jump between two readings
that sit next to each other — the size of the jump, never its direction.

    swing([3, 10, 4])  ->  7

Fewer than two readings means there is no jump at all, which is 0.`,
		timeLimitMs: 5000,
		starters:    py("def swing(readings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the biggest jump", "swing([3, 10, 4]) == 7"),
			vis("A drop counts the same as a rise", "swing([10, 1]) == 9"),
			vis("One reading has no jump", "swing([5]) == 0"),
			vis("No readings at all", "swing([]) == 0"),
			hid("The biggest jump may be at the end", "swing([1, 2, 3, 40]) == 37"),
			hid("A flat run swings by nothing", "swing([7, 7, 7]) == 0"),
		},
		solution: "def swing(readings):\n    return max((abs(a - b) for a, b in zip(readings, readings[1:])), default=0)\n",
	},
	{
		slug: "shortest-queue", title: "Shortest queue", difficulty: "easy", topic: "Loops",
		statement: `Each checkout lane holds a list of baskets, and each basket is a number
of items. Write ~pick(lanes)~ returning the index of the lane with the fewest
items in total. If two lanes tie, take the one further left.

    pick([[3, 4], [10], [1, 1, 1]])  ->  2

There is always at least one lane. A lane with nobody in it holds zero items,
which is as short as a lane gets.`,
		timeLimitMs: 5000,
		starters:    py("def pick(lanes):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Picks the lane with fewest items", "pick([[3, 4], [10], [1, 1, 1]]) == 2"),
			vis("Counts items, not people", "pick([[9], [1, 1]]) == 1"),
			vis("An empty lane wins", "pick([[5], []]) == 1"),
			vis("A single lane is the answer", "pick([[8]]) == 0"),
			hid("A tie goes to the leftmost lane", "pick([[2, 2], [4], [1, 3]]) == 0"),
		},
		solution: "def pick(lanes):\n    return min(range(len(lanes)), key=lambda i: (sum(lanes[i]), i))\n",
	},
	{
		slug: "parking-meter", title: "Parking meter", difficulty: "easy", topic: "Arithmetic",
		statement: `The first 30 minutes are free. After that you pay 50 cents for every
15-minute block you have **started** — one minute into a block costs the whole
block. A day never costs more than 900 cents.

Write ~cost(minutes)~ returning the charge in cents.

    cost(31)  ->  50
    cost(46)  ->  100

Watch the word *started*: 45 minutes is exactly one block past the free half
hour, and 46 minutes is into the second.`,
		timeLimitMs: 5000,
		starters:    py("def cost(minutes):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The free half hour costs nothing", "cost(30) == 0"),
			vis("One minute over starts a block", "cost(31) == 50"),
			vis("A started second block is charged in full", "cost(46) == 100"),
			vis("A long stay hits the daily cap", "cost(10000) == 900"),
			hid("Exactly one block past the free time", "cost(45) == 50"),
			hid("Arriving and leaving costs nothing", "cost(0) == 0"),
		},
		solution: "def cost(minutes):\n    paid = max(0, minutes - 30)\n    blocks = -(-paid // 15)\n    return min(900, 50 * blocks)\n",
	},
}
