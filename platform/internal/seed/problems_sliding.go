package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// windowProblems: two pointers, sliding windows and prefix sums — the family of
// tricks whose whole point is not redoing work you have already done.
var windowProblems = []seedProblem{
	{
		slug: "busiest-window", title: "Busiest window", difficulty: "medium", topic: "Sliding window",
		statement: `A turnstile counts arrivals each minute. Write ~busiest(counts, k)~
returning the starting minute of the ~k~-minute stretch with the most arrivals.
If two stretches tie, take the earlier one.

    busiest([1, 5, 2, 3], 2)  ->  1

If ~k~ is bigger than the log, or not positive, there is no such stretch:
return ~-1~.

Logs run to hundreds of thousands of minutes. Adding up each window from
scratch does the same additions over and over — a window is the last one with
one number swapped out and one swapped in.`,
		timeLimitMs: 4000,
		starters:    py("def busiest(counts, k):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the busiest stretch", "busiest([1, 5, 2, 3], 2) == 1"),
			vis("A tie goes to the earlier stretch", "busiest([2, 2, 2], 1) == 0"),
			vis("A window longer than the log", "busiest([1, 2], 5) == -1"),
			vis("A window of nothing", "busiest([1, 2], 0) == -1"),
			vis("Fast enough for a long log", "busiest(list(range(200000)), 5000) == 195000"),
			hid("The whole log is one window", "busiest([4, 1], 2) == 0"),
			hid("The busiest stretch is at the end", "busiest([1, 1, 9], 1) == 2"),
		},
		solution: "def busiest(counts, k):\n    n = len(counts)\n    if k <= 0 or k > n:\n        return -1\n    total = sum(counts[:k])\n    best, best_at = total, 0\n    for i in range(k, n):\n        total += counts[i] - counts[i - k]\n        if total > best:\n            best, best_at = total, i - k + 1\n    return best_at\n",
		tooSlow:  "def busiest(counts, k):\n    n = len(counts)\n    if k <= 0 or k > n:\n        return -1\n    best, best_at = None, 0\n    for i in range(n - k + 1):\n        s = sum(counts[i:i + k])\n        if best is None or s > best:\n            best, best_at = s, i\n    return best_at\n",
	},
	{
		slug: "calm-stretch", title: "Calm stretch", difficulty: "easy", topic: "Sliding window",
		statement: `A instrument is calm while it never jumps by more than ~tol~ from one
reading to the next. Write ~calm(readings, tol)~ returning the length of the
longest calm stretch.

    calm([1, 2, 10], 3)  ->  2

A single reading is a stretch of one — it has not jumped anywhere. No readings
at all is a stretch of zero.`,
		timeLimitMs: 5000,
		starters:    py("def calm(readings, tol):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the calm run", "calm([1, 2, 10], 3) == 2"),
			vis("One reading is a stretch of one", "calm([7], 1) == 1"),
			vis("No readings", "calm([], 1) == 0"),
			vis("Everything is calm", "calm([1, 2, 3, 4], 1) == 4"),
			hid("A jump exactly at the tolerance is still calm", "calm([1, 4], 3) == 2"),
			hid("A drop counts as much as a rise", "calm([10, 1, 2], 3) == 2"),
		},
		solution: "def calm(readings, tol):\n    if not readings:\n        return 0\n    best = run = 1\n    for a, b in zip(readings, readings[1:]):\n        run = run + 1 if abs(b - a) <= tol else 1\n        if run > best:\n            best = run\n    return best\n",
	},
	{
		slug: "balanced-stretch", title: "Balanced stretch", difficulty: "hard", topic: "Prefix sums",
		statement: `A market log records each day as ~"up"~ or ~"down"~. Write
~balanced(days)~ returning the length of the longest stretch with exactly as
many up days as down days.

    balanced(["up", "down", "up"])  ->  2

Checking every stretch is too slow for a log of a hundred thousand days. The way
in: keep a running score, +1 for up and -1 for down. Two days with the **same**
running score have a balanced stretch between them — so all you need to remember
is the first day each score was seen.`,
		timeLimitMs: 4000,
		starters:    py("def balanced(days):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the balanced stretch", `balanced(["up", "down", "up"]) == 2`),
			vis("Nothing balances", `balanced(["up", "up"]) == 0`),
			vis("An empty log", "balanced([]) == 0"),
			vis("The whole log balances", `balanced(["up", "up", "down", "down"]) == 4`),
			vis("Fast enough for a long log", `balanced(["up", "down"] * 50000) == 100000`),
			hid("The stretch need not start at the beginning", `balanced(["up", "up", "down"]) == 2`),
			hid("A long imbalance around a short balance", `balanced(["up", "up", "up", "down", "up"]) == 2`),
		},
		solution: "def balanced(days):\n    first = {0: -1}\n    score = 0\n    best = 0\n    for i, d in enumerate(days):\n        score += 1 if d == 'up' else -1\n        if score in first:\n            if i - first[score] > best:\n                best = i - first[score]\n        else:\n            first[score] = i\n    return best\n",
		tooSlow:  "def balanced(days):\n    best = 0\n    for i in range(len(days)):\n        score = 0\n        for j in range(i, len(days)):\n            score += 1 if days[j] == 'up' else -1\n            if score == 0 and j - i + 1 > best:\n                best = j - i + 1\n    return best\n",
	},
	{
		slug: "average-alarm", title: "Average alarm", difficulty: "easy", topic: "Sliding window",
		statement: `A monitor raises an alarm whenever the average of the last ~k~ readings is
**strictly above** ~limit~. Write ~alarms(readings, k, limit)~ returning how
many windows of exactly ~k~ readings raise one.

    alarms([1, 5, 1], 2, 2)  ->  2

Both windows average 3. If the log is shorter than ~k~, no window exists, so no
alarms. Work in whole numbers: comparing the sum against ~limit * k~ says the
same thing as comparing the average against ~limit~, without any rounding to
argue about.`,
		timeLimitMs: 5000,
		starters:    py("def alarms(readings, k, limit):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Both windows trip it", "alarms([1, 5, 1], 2, 2) == 2"),
			vis("A quiet log raises nothing", "alarms([1, 1, 1], 2, 5) == 0"),
			vis("Too short for a window", "alarms([9], 2, 0) == 0"),
			hid("Exactly at the limit is not above it", "alarms([2, 2], 2, 2) == 0"),
			hid("Only some windows trip it", "alarms([9, 0, 0, 0], 2, 2) == 1"),
		},
		solution: "def alarms(readings, k, limit):\n    n = len(readings)\n    if k <= 0 or k > n:\n        return 0\n    total = sum(readings[:k])\n    count = 1 if total > limit * k else 0\n    for i in range(k, n):\n        total += readings[i] - readings[i - k]\n        if total > limit * k:\n            count += 1\n    return count\n",
	},
	{
		slug: "affordable-pairs", title: "Affordable pairs", difficulty: "medium", topic: "Two pointers",
		statement: `A price list is already sorted, cheapest first. Write ~pairs(prices, budget)~
returning how many pairs of two **different** items cost ~budget~ or less
together.

    pairs([1, 2, 3, 4], 5)  ->  4

Those are 1+2, 1+3, 1+4 and 2+3. A list of two hundred thousand prices makes
checking every pair impossible — but the list is sorted, and that is the whole
gift. If the cheapest and the dearest fit together, so does the cheapest with
everything in between.`,
		timeLimitMs: 4000,
		starters:    py("def pairs(prices, budget):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts the affordable pairs", "pairs([1, 2, 3, 4], 5) == 4"),
			vis("Nothing is affordable", "pairs([9, 9], 5) == 0"),
			vis("One item makes no pair", "pairs([1], 100) == 0"),
			vis("Everything is affordable", "pairs([1, 1, 1], 5) == 3"),
			vis("Fast enough for a long list", "pairs(list(range(200000)), 10) == 30"),
			hid("Exactly on budget counts", "pairs([2, 3], 5) == 1"),
			hid("An empty list", "pairs([], 5) == 0"),
		},
		solution: "def pairs(prices, budget):\n    i, j = 0, len(prices) - 1\n    n = 0\n    while i < j:\n        if prices[i] + prices[j] <= budget:\n            n += j - i\n            i += 1\n        else:\n            j -= 1\n    return n\n",
		tooSlow:  "def pairs(prices, budget):\n    n = 0\n    for i in range(len(prices)):\n        for j in range(i + 1, len(prices)):\n            if prices[i] + prices[j] <= budget:\n                n += 1\n    return n\n",
	},
	{
		slug: "range-totals", title: "Range totals", difficulty: "medium", topic: "Prefix sums",
		statement: `Write ~totals(values, queries)~, where each query is a pair ~(lo, hi)~
asking for the sum of ~values[lo]~ through ~values[hi]~ **inclusive**. Return
one answer per query, in order.

    totals([1, 2, 3], [(0, 2), (1, 1)])  ->  [6, 2]

There can be tens of thousands of queries over a list of a hundred thousand
values. Answering each one by adding up its range does the same additions again
and again; answering all of them from one pass over the values does not.`,
		timeLimitMs: 4000,
		starters:    py("def totals(values, queries):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Answers each query", "totals([1, 2, 3], [(0, 2), (1, 1)]) == [6, 2]"),
			vis("No queries", "totals([1, 2], []) == []"),
			vis("A single-element range", "totals([5], [(0, 0)]) == [5]"),
			vis("Fast enough for many queries", "totals(list(range(100000)), [(0, 99999)] * 20000) == [4999950000] * 20000"),
			hid("Overlapping ranges", "totals([1, 1, 1, 1], [(0, 1), (1, 3), (0, 3)]) == [2, 3, 4]"),
			hid("Negative values sum too", "totals([-1, 5], [(0, 1)]) == [4]"),
		},
		solution: "def totals(values, queries):\n    prefix = [0]\n    for v in values:\n        prefix.append(prefix[-1] + v)\n    return [prefix[hi + 1] - prefix[lo] for lo, hi in queries]\n",
		tooSlow:  "def totals(values, queries):\n    return [sum(values[lo:hi + 1]) for lo, hi in queries]\n",
	},
	{
		slug: "worst-drawdown", title: "Worst drawdown", difficulty: "medium", topic: "Prefix sums",
		statement: `A drawdown is how far a price has fallen from the highest point it had
reached **before** it. Write ~drawdown(prices)~ returning the worst one.

    drawdown([5, 3, 6, 2])  ->  4

The 2 comes after a peak of 6. A price that only ever rises has a drawdown of 0,
and so does an empty list.`,
		timeLimitMs: 5000,
		starters:    py("def drawdown(prices):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Measures from the peak before it", "drawdown([5, 3, 6, 2]) == 4"),
			vis("A rising market never falls", "drawdown([1, 2, 3]) == 0"),
			vis("No prices at all", "drawdown([]) == 0"),
			vis("A single price", "drawdown([7]) == 0"),
			hid("A later peak does not count backwards", "drawdown([9, 1, 100]) == 8"),
			hid("The worst fall is the first one", "drawdown([10, 1, 5, 4]) == 9"),
		},
		solution: "def drawdown(prices):\n    peak = None\n    worst = 0\n    for p in prices:\n        if peak is None or p > peak:\n            peak = p\n        if peak - p > worst:\n            worst = peak - p\n    return worst\n",
	},
	{
		slug: "pivot-points", title: "Pivot points", difficulty: "medium", topic: "Prefix sums",
		statement: `A pivot is a position where everything to its left weighs exactly as much
as everything to its right. The value at the pivot itself belongs to neither
side.

Write ~pivots(weights)~ returning every pivot position, in order.

    pivots([1, 2, 3, 3])  ->  [2]

An end position is a pivot when the other side sums to zero — including when
that side is empty. Lists run to hundreds of thousands of entries, so re-adding
each side at every position will not finish.`,
		timeLimitMs: 4000,
		starters:    py("def pivots(weights):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the balance point", "pivots([1, 2, 3, 3]) == [2]"),
			vis("A single item balances nothing against nothing", "pivots([9]) == [0]"),
			vis("No items, no pivots", "pivots([]) == []"),
			vis("Fast enough for a long list", "pivots([0] * 200000) == list(range(200000))"),
			hid("There can be more than one", "pivots([0, 0, 0]) == [0, 1, 2]"),
			hid("There can be none", "pivots([1, 9]) == []"),
			hid("Negative weights can balance", "pivots([1, -1, 0, 3, -3]) == [2]"),
		},
		solution: "def pivots(weights):\n    total = sum(weights)\n    left = 0\n    out = []\n    for i, w in enumerate(weights):\n        if left == total - left - w:\n            out.append(i)\n        left += w\n    return out\n",
		tooSlow:  "def pivots(weights):\n    return [i for i in range(len(weights))\n            if sum(weights[:i]) == sum(weights[i + 1:])]\n",
	},
	{
		slug: "log-throttle", title: "Log throttle", difficulty: "medium", topic: "Sliding window",
		statement: `A logger accepts at most ~cap~ messages in any ~window~ of time. Messages
arrive in time order. When one arrives and ~cap~ messages have already been
accepted within the last ~window~ ticks, it is dropped — and a dropped message
does not count against the limit either, because it was never logged.

A message at time ~t~ counts against a message at time ~u~ only while
~u > t - window~.

Write ~dropped(times, window, cap)~ returning how many are lost.

    dropped([1, 2, 3], 10, 2)  ->  1`,
		timeLimitMs: 5000,
		starters:    py("def dropped(times, window, cap):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The third message is over the cap", "dropped([1, 2, 3], 10, 2) == 1"),
			vis("Spread out enough to all get through", "dropped([1, 20, 40], 10, 1) == 0"),
			vis("Nothing logged", "dropped([], 10, 1) == 0"),
			hid("A message ageing out frees a slot", "dropped([1, 2, 11], 10, 2) == 0"),
			hid("A dropped message does not hold a slot", "dropped([1, 1, 1, 30], 10, 1) == 2"),
			hid("A cap of nothing drops everything", "dropped([1, 2], 10, 0) == 2"),
		},
		solution: "def dropped(times, window, cap):\n    kept = []\n    lost = 0\n    for t in times:\n        while kept and kept[0] <= t - window:\n            kept.pop(0)\n        if len(kept) < cap:\n            kept.append(t)\n        else:\n            lost += 1\n    return lost\n",
	},
	{
		slug: "cooldown-check", title: "Cooldown check", difficulty: "easy", topic: "Dictionaries",
		statement: `A game gives each ability a cooldown: once used, it cannot be used again for
~cooldown~ turns. Measure the gap as the distance between the two turns — turn
0 and turn 3 are 3 apart, so a cooldown of 3 permits that and a cooldown of 4
does not.

Write ~allowed(actions, cooldown)~ returning whether a whole sequence of turns
is legal.

    allowed(["a", "b", "a"], 2)  ->  True
    allowed(["a", "a"], 2)       ->  False`,
		timeLimitMs: 5000,
		starters:    py("def allowed(actions, cooldown):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Far enough apart", `allowed(["a", "b", "a"], 2) is True`),
			vis("Used again too soon", `allowed(["a", "a"], 2) is False`),
			vis("A gap exactly equal to the cooldown is legal", `allowed(["a", "a"], 1) is True`),
			vis("Nothing to check", "allowed([], 5) is True"),
			vis("All different abilities", `allowed(["a", "b", "c"], 10) is True`),
			hid("Exactly at the cooldown is allowed", `allowed(["a", "b", "c", "a"], 3) is True`),
			hid("One turn short is not", `allowed(["a", "b", "c", "a"], 4) is False`),
			hid("A cooldown of zero permits anything", `allowed(["a", "a"], 0) is True`),
		},
		solution: "def allowed(actions, cooldown):\n    last = {}\n    for i, a in enumerate(actions):\n        if a in last and i - last[a] < cooldown:\n            return False\n        last[a] = i\n    return True\n",
	},
	{
		slug: "neighbour-smoothing", title: "Neighbour smoothing", difficulty: "easy", topic: "Sliding window",
		statement: `Smooth a noisy signal: replace each reading with the whole-number average
of itself and its immediate neighbours, rounding **down**.

Readings at the ends have only one neighbour, so they average over two values
rather than three. Every average is taken from the **original** readings, not
from the smoothed ones you are producing.

Write ~smooth(readings)~.

    smooth([1, 2, 3])  ->  [1, 2, 2]`,
		timeLimitMs: 5000,
		starters:    py("def smooth(readings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Smooths a short signal", "smooth([1, 2, 3]) == [1, 2, 2]"),
			vis("A single reading is its own average", "smooth([5]) == [5]"),
			vis("Nothing to smooth", "smooth([]) == []"),
			vis("A flat signal stays flat", "smooth([4, 4, 4]) == [4, 4, 4]"),
			hid("Averages come from the original readings", "smooth([0, 9, 0, 9, 0]) == [4, 3, 6, 3, 4]"),
			hid("Rounding is always down", "smooth([0, 1]) == [0, 0]"),
		},
		solution: "def smooth(readings):\n    n = len(readings)\n    out = []\n    for i in range(n):\n        window = readings[max(0, i - 1):min(n, i + 2)]\n        out.append(sum(window) // len(window))\n    return out\n",
	},
	{
		slug: "peak-hours", title: "Peak hours", difficulty: "easy", topic: "Sliding window",
		statement: `Write ~peaks(counts)~ returning the positions that are strictly higher than
every neighbour they have.

A position at either end has only one neighbour, and only has to beat that one.
A lone reading has no neighbours at all, so it is a peak.

    peaks([1, 3, 2])  ->  [1]`,
		timeLimitMs: 5000,
		starters:    py("def peaks(counts):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the peak", "peaks([1, 3, 2]) == [1]"),
			vis("A lone reading is a peak", "peaks([5]) == [0]"),
			vis("Nothing at all", "peaks([]) == []"),
			vis("A peak at the end", "peaks([1, 2]) == [1]"),
			hid("A plateau is not a peak", "peaks([1, 2, 2, 1]) == []"),
			hid("Two separate peaks", "peaks([3, 1, 3, 1]) == [0, 2]"),
		},
		solution: "def peaks(counts):\n    out = []\n    n = len(counts)\n    for i, v in enumerate(counts):\n        if (i == 0 or v > counts[i - 1]) and (i == n - 1 or v > counts[i + 1]):\n            out.append(i)\n    return out\n",
	},
}
