package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// intervalProblems: things with a start and an end, and the sweeps and greedy
// choices that make sense of a pile of them.
var intervalProblems = []seedProblem{
	{
		slug: "merge-bookings", title: "Merge bookings", difficulty: "medium", topic: "Intervals",
		statement: `Write ~merge_bookings(bookings)~, collapsing overlapping ~(start, end)~
bookings into the fewest that cover the same time, sorted by start.

Bookings that merely touch — one ending exactly where the next begins — become
one booking. Bookings arrive in no particular order.

    merge_bookings([(1, 3), (2, 4), (6, 8)])  ->  [(1, 4), (6, 8)]
    merge_bookings([(1, 2), (2, 3)])          ->  [(1, 3)]`,
		timeLimitMs: 5000,
		starters:    py("def merge_bookings(bookings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Overlapping bookings become one", "merge_bookings([(1, 3), (2, 4), (6, 8)]) == [(1, 4), (6, 8)]"),
			vis("Touching bookings merge too", "merge_bookings([(1, 2), (2, 3)]) == [(1, 3)]"),
			vis("Nothing booked", "merge_bookings([]) == []"),
			vis("Separate bookings stay separate", "merge_bookings([(1, 2), (5, 6)]) == [(1, 2), (5, 6)]"),
			hid("Order of arrival does not matter", "merge_bookings([(6, 8), (1, 3), (2, 4)]) == [(1, 4), (6, 8)]"),
			hid("A booking swallowed whole by another", "merge_bookings([(1, 10), (2, 3)]) == [(1, 10)]"),
			hid("A chain of overlaps collapses to one", "merge_bookings([(1, 3), (2, 5), (4, 9)]) == [(1, 9)]"),
		},
		solution: "def merge_bookings(bookings):\n    out = []\n    for s, e in sorted(bookings):\n        if out and s <= out[-1][1]:\n            out[-1] = (out[-1][0], max(out[-1][1], e))\n        else:\n            out.append((s, e))\n    return out\n",
	},
	{
		slug: "double-booked", title: "Double booked", difficulty: "medium", topic: "Intervals",
		statement: `Write ~clash(bookings)~ returning the first two bookings that genuinely
overlap — as a tuple ~(earlier, later)~ ordered by start time — or ~None~ if the
diary is clean.

A booking that ends exactly when another begins is not a clash. "First" means
the clash whose earlier booking starts soonest.

    clash([(1, 3), (2, 4)])  ->  ((1, 3), (2, 4))
    clash([(1, 2), (2, 3)])  ->  None`,
		timeLimitMs: 5000,
		starters:    py("def clash(bookings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the overlap", "clash([(1, 3), (2, 4)]) == ((1, 3), (2, 4))"),
			vis("Touching is not clashing", "clash([(1, 2), (2, 3)]) is None"),
			vis("An empty diary", "clash([]) is None"),
			vis("A single booking clashes with nothing", "clash([(1, 5)]) is None"),
			hid("Order of arrival does not matter", "clash([(2, 4), (1, 3)]) == ((1, 3), (2, 4))"),
			hid("A booking inside another is a clash", "clash([(1, 10), (2, 3)]) == ((1, 10), (2, 3))"),
			hid("The earliest clash is reported", "clash([(1, 9), (2, 3), (5, 6)]) == ((1, 9), (2, 3))"),
		},
		solution: "def clash(bookings):\n    ivs = sorted(bookings)\n    for a, b in zip(ivs, ivs[1:]):\n        if b[0] < a[1]:\n            return (a, b)\n    return None\n",
	},
	{
		slug: "rooms-needed", title: "Rooms needed", difficulty: "hard", topic: "Intervals",
		statement: `Write ~rooms(bookings)~ returning the fewest rooms that could hold every
booking — which is the largest number of bookings happening at the same moment.

A booking that ends exactly when another starts can reuse the room.

    rooms([(1, 3), (2, 4), (3, 5)])  ->  2

A diary can hold a hundred thousand bookings, so counting how many are live at
each booking's start time, one booking at a time, will not finish. Think of the
starts and ends as one sequence of events and walk it in order.`,
		timeLimitMs: 4000,
		starters:    py("def rooms(bookings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Two at once at the busiest moment", "rooms([(1, 3), (2, 4), (3, 5)]) == 2"),
			vis("An empty diary needs no rooms", "rooms([]) == 0"),
			vis("Bookings back to back share a room", "rooms([(1, 2), (2, 3), (3, 4)]) == 1"),
			vis("Everything at once", "rooms([(1, 9), (1, 9), (1, 9)]) == 3"),
			vis("Fast enough for a full diary", "rooms([(i, i + 2) for i in range(100000)]) == 2"),
			hid("A long booking spanning many short ones", "rooms([(0, 100), (1, 2), (3, 4)]) == 2"),
		},
		solution: "def rooms(bookings):\n    events = []\n    for s, e in bookings:\n        events.append((s, 1))\n        events.append((e, -1))\n    events.sort()\n    live = best = 0\n    for _, delta in events:\n        live += delta\n        if live > best:\n            best = live\n    return best\n",
		tooSlow:  "def rooms(bookings):\n    best = 0\n    for s, _ in bookings:\n        n = sum(1 for a, b in bookings if a <= s < b)\n        if n > best:\n            best = n\n    return best\n",
	},
	{
		slug: "free-slots", title: "Free slots", difficulty: "medium", topic: "Intervals",
		statement: `A day runs from 0 to ~day_end~. Write ~free(busy, day_end, min_len)~
returning the gaps of at least ~min_len~ where nothing is booked, in order.

Busy periods arrive unsorted and may overlap each other.

    free([(2, 4)], 10, 2)  ->  [(0, 2), (4, 10)]

The gap before the first booking and the one after the last both count.`,
		timeLimitMs: 5000,
		starters:    py("def free(busy, day_end, min_len):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Gaps either side of a booking", "free([(2, 4)], 10, 2) == [(0, 2), (4, 10)]"),
			vis("A completely empty day is one gap", "free([], 8, 1) == [(0, 8)]"),
			vis("Short gaps are not offered", "free([(1, 2)], 3, 2) == []"),
			hid("Overlapping busy periods count once", "free([(2, 6), (3, 4)], 10, 2) == [(0, 2), (6, 10)]"),
			hid("A fully booked day has no gaps", "free([(0, 10)], 10, 1) == []"),
			hid("A gap exactly the minimum length counts", "free([(0, 5)], 7, 2) == [(5, 7)]"),
		},
		solution: "def free(busy, day_end, min_len):\n    out = []\n    t = 0\n    for s, e in sorted(busy):\n        if s - t >= min_len:\n            out.append((t, s))\n        if e > t:\n            t = e\n    if day_end - t >= min_len:\n        out.append((t, day_end))\n    return out\n",
	},
	{
		slug: "talks-attended", title: "Talks attended", difficulty: "hard", topic: "Greedy",
		statement: `A conference runs talks at overlapping times. Write ~attend(talks)~
returning the most talks you could sit through, given you must stay for a whole
talk and cannot be in two places at once. A talk starting exactly when another
ends is fine.

    attend([(1, 3), (2, 4), (3, 5)])  ->  2

The instinct is to take the talk that starts soonest, or the shortest one.
Neither is right. What you actually want at every moment is to be free again as
early as possible.`,
		timeLimitMs: 5000,
		starters:    py("def attend(talks):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Two of the three fit", "attend([(1, 3), (2, 4), (3, 5)]) == 2"),
			vis("Nothing on", "attend([]) == 0"),
			vis("Talks back to back all fit", "attend([(1, 2), (2, 3), (3, 4)]) == 3"),
			vis("Everything clashes with everything", "attend([(1, 9), (2, 8), (3, 7)]) == 1"),
			hid("Earliest start is the wrong choice", "attend([(1, 10), (2, 3), (4, 5)]) == 2"),
			hid("Order of arrival does not matter", "attend([(3, 5), (1, 3), (2, 4)]) == 2"),
			hid("A single talk", "attend([(0, 1)]) == 1"),
		},
		solution: "def attend(talks):\n    n = 0\n    last = None\n    for s, e in sorted(talks, key=lambda t: t[1]):\n        if last is None or s >= last:\n            n += 1\n            last = e\n    return n\n",
	},
	{
		slug: "cover-the-day", title: "Cover the day", difficulty: "hard", topic: "Greedy",
		statement: `Volunteers each offer a ~(start, end)~ shift. Write ~covers(shifts, day_end)~
returning whether their shifts together cover every moment from 0 to
~day_end~ with nobody unattended.

    covers([(0, 5), (4, 10)], 10)  ->  True
    covers([(0, 4), (5, 10)], 10)  ->  False

A day of length 0 is covered by nobody at all. Shifts may overlap and arrive in
any order.

The greedy move: from wherever you are covered up to, take the shift that
reaches furthest among those that start at or before that point.`,
		timeLimitMs: 5000,
		starters:    py("def covers(shifts, day_end):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Overlapping shifts cover the day", "covers([(0, 5), (4, 10)], 10) is True"),
			vis("A gap in the middle", "covers([(0, 4), (5, 10)], 10) is False"),
			vis("A day of no length needs nobody", "covers([], 0) is True"),
			vis("Nobody covers a real day", "covers([], 5) is False"),
			hid("The day must be covered from the very start", "covers([(1, 10)], 10) is False"),
			hid("Reaching short of the end fails", "covers([(0, 9)], 10) is False"),
			hid("The longest reach is not always the first shift", "covers([(0, 2), (1, 3), (1, 10)], 10) is True"),
			hid("Shifts touching exactly still cover", "covers([(0, 5), (5, 10)], 10) is True"),
		},
		solution: "def covers(shifts, day_end):\n    ivs = sorted(shifts)\n    reach = 0\n    i = 0\n    while reach < day_end:\n        best = reach\n        while i < len(ivs) and ivs[i][0] <= reach:\n            if ivs[i][1] > best:\n                best = ivs[i][1]\n            i += 1\n        if best == reach:\n            return False\n        reach = best\n    return True\n",
	},
	{
		slug: "longest-free-gap", title: "Longest free gap", difficulty: "easy", topic: "Intervals",
		statement: `Write ~gap(meetings)~ returning the longest stretch of free time between
meetings — from the moment you are free until the next one begins.

Meetings arrive unsorted and may overlap. Time before the first meeting and
after the last does not count: this is the gap *between* meetings only.

    gap([(1, 2), (5, 6)])  ->  3

Fewer than two meetings means no gap at all.`,
		timeLimitMs: 5000,
		starters:    py("def gap(meetings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The gap between two meetings", "gap([(1, 2), (5, 6)]) == 3"),
			vis("One meeting has no gap", "gap([(1, 2)]) == 0"),
			vis("No meetings at all", "gap([]) == 0"),
			vis("Back to back leaves no gap", "gap([(1, 2), (2, 3)]) == 0"),
			hid("A meeting inside another leaves no gap", "gap([(1, 10), (2, 3)]) == 0"),
			hid("The longest of several gaps", "gap([(0, 1), (2, 3), (9, 10)]) == 6"),
			hid("Order of arrival does not matter", "gap([(5, 6), (1, 2)]) == 3"),
		},
		solution: "def gap(meetings):\n    best = 0\n    end = None\n    for s, e in sorted(meetings):\n        if end is not None and s - end > best:\n            best = s - end\n        end = e if end is None else max(end, e)\n    return best\n",
	},
	{
		slug: "shared-availability", title: "Shared availability", difficulty: "medium", topic: "Two pointers",
		statement: `Two people each give their free time as a sorted list of ~(start, end)~
ranges that do not overlap each other. Write ~shared(a, b)~ returning when they
are both free, as a sorted list of ranges.

A range of zero length is not a meeting: ranges that merely touch produce
nothing.

    shared([(1, 5)], [(3, 7)])  ->  [(3, 5)]
    shared([(1, 3)], [(3, 5)])  ->  []

Both lists can be long. Because both are already sorted you never need to look
backwards — whichever range ends first can be retired.`,
		timeLimitMs: 5000,
		starters:    py("def shared(a, b):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The overlapping part", "shared([(1, 5)], [(3, 7)]) == [(3, 5)]"),
			vis("Touching is not overlapping", "shared([(1, 3)], [(3, 5)]) == []"),
			vis("One person is never free", "shared([(1, 5)], []) == []"),
			vis("Neither is ever free", "shared([], []) == []"),
			hid("One range can overlap several", "shared([(0, 10)], [(1, 2), (3, 4)]) == [(1, 2), (3, 4)]"),
			hid("Several overlaps in order", "shared([(1, 4), (6, 9)], [(2, 7)]) == [(2, 4), (6, 7)]"),
			hid("No overlap at all", "shared([(1, 2)], [(5, 6)]) == []"),
		},
		solution: "def shared(a, b):\n    out = []\n    i = j = 0\n    while i < len(a) and j < len(b):\n        lo = max(a[i][0], b[j][0])\n        hi = min(a[i][1], b[j][1])\n        if lo < hi:\n            out.append((lo, hi))\n        if a[i][1] < b[j][1]:\n            i += 1\n        else:\n            j += 1\n    return out\n",
	},
	{
		slug: "worst-lateness", title: "Worst lateness", difficulty: "hard", topic: "Greedy",
		statement: `A workshop has jobs, each ~(duration, deadline)~. One job runs at a time,
starting at time 0, with no gaps. A job's lateness is when it finishes minus its
deadline — negative if it finishes early.

You choose the order. Write ~worst(jobs)~ returning the smallest possible value
of the **largest** lateness across all jobs.

    worst([(4, 2), (1, 3)])  ->  2

Running the short job first finishes both sooner, but the tight deadline slips
further. There is one ordering rule that is always optimal here, and it is worth
convincing yourself why swapping any two adjacent jobs out of that order never
helps.

No jobs at all has a worst lateness of 0.`,
		timeLimitMs: 5000,
		starters:    py("def worst(jobs):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The tight deadline goes first", "worst([(4, 2), (1, 3)]) == 2"),
			vis("No jobs, no lateness", "worst([]) == 0"),
			vis("Finishing early gives negative lateness", "worst([(1, 10)]) == -9"),
			hid("Everything comfortably on time", "worst([(1, 5), (1, 6)]) == -4"),
			hid("Order genuinely matters — the other way round gives 3", "worst([(3, 3), (1, 1)]) == 1"),
			hid("Several jobs sharing a deadline", "worst([(1, 4), (1, 4), (1, 4)]) == -1"),
		},
		solution: "def worst(jobs):\n    t = 0\n    late = None\n    for dur, due in sorted(jobs, key=lambda j: j[1]):\n        t += dur\n        if late is None or t - due > late:\n            late = t - due\n    return 0 if late is None else late\n",
	},
	{
		slug: "room-numbers", title: "Room numbers", difficulty: "medium", topic: "Intervals",
		statement: `Bookings arrive in start order and each takes the **lowest-numbered** room
that is free. Rooms are numbered from 0, and a new room is opened only when
every existing one is busy. A room is free again the moment its booking ends.

Write ~assign(bookings)~ returning the room each booking gets.

    assign([(1, 3), (2, 4), (3, 5)])  ->  [0, 1, 0]

The third booking starts exactly when the first ends, so room 0 is free again.`,
		timeLimitMs: 5000,
		starters:    py("def assign(bookings):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Reuses the lowest free room", "assign([(1, 3), (2, 4), (3, 5)]) == [0, 1, 0]"),
			vis("Nothing booked", "assign([]) == []"),
			vis("Back-to-back bookings share one room", "assign([(1, 2), (2, 3), (3, 4)]) == [0, 0, 0]"),
			vis("Everything at once opens a room each", "assign([(1, 9), (1, 9), (1, 9)]) == [0, 1, 2]"),
			hid("A freed low room beats a freed high one", "assign([(0, 5), (1, 2), (3, 9)]) == [0, 1, 1]"),
			hid("One long booking holds its room throughout", "assign([(0, 100), (1, 2), (3, 4)]) == [0, 1, 1]"),
		},
		solution: "def assign(bookings):\n    free_at = []\n    out = []\n    for start, end in bookings:\n        room = None\n        for r, e in enumerate(free_at):\n            if e <= start:\n                room = r\n                break\n        if room is None:\n            free_at.append(end)\n            room = len(free_at) - 1\n        else:\n            free_at[room] = end\n        out.append(room)\n    return out\n",
	},
}
