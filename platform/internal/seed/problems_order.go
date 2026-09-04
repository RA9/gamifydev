package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// orderProblems: sorting, ranking, and the tie-breaks and derived keys that
// make an ordering mean something.
var orderProblems = []seedProblem{
	{
		slug: "competition-ranks", title: "Competition ranks", difficulty: "medium", topic: "Sorting",
		statement: `A race is scored the way races are: the highest score is 1st, everyone
tied shares the same rank, and the next rank after a tie **skips**. Two people
in 2nd place means nobody is 3rd.

Write ~ranks(scores)~ returning each competitor's rank, in the same order the
scores were given.

    ranks([10, 8, 8, 5])  ->  [1, 2, 2, 4]

A field of two hundred thousand runners is not unusual, so counting how many
people beat each runner one runner at a time will not finish in time.`,
		timeLimitMs: 4000,
		starters:    py("def ranks(scores):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A tie shares a rank and skips the next", "ranks([10, 8, 8, 5]) == [1, 2, 2, 4]"),
			vis("Everyone distinct", "ranks([3, 1, 2]) == [1, 3, 2]"),
			vis("Nobody raced", "ranks([]) == []"),
			vis("Everyone tied for first", "ranks([5, 5, 5]) == [1, 1, 1]"),
			vis("Fast enough for a huge field", "ranks(list(range(200000)))[0] == 200000"),
			hid("A tie at the bottom", "ranks([9, 1, 1]) == [1, 2, 2]"),
		},
		solution: "def ranks(scores):\n    counts = {}\n    for s in scores:\n        counts[s] = counts.get(s, 0) + 1\n    pos = {}\n    n = 1\n    for s in sorted(counts, reverse=True):\n        pos[s] = n\n        n += counts[s]\n    return [pos[s] for s in scores]\n",
		tooSlow:  "def ranks(scores):\n    return [1 + sum(1 for o in scores if o > s) for s in scores]\n",
	},
	{
		slug: "version-order", title: "Version order", difficulty: "medium", topic: "Sorting",
		statement: `Version strings are dotted numbers, and they do not sort like text:
~1.10~ comes **after** ~1.2~, because ten is more than two.

Write ~sort_versions(versions)~ returning them in ascending order. Compare
segment by segment as numbers. A version with fewer segments is padded with
zeros, so ~1.2~ and ~1.2.0~ are the same version and keep the order they were
given in.

    sort_versions(["1.10", "1.2", "1.2.1"])  ->  ["1.2", "1.2.1", "1.10"]`,
		timeLimitMs: 5000,
		starters:    py("def sort_versions(versions):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Ten sorts after two", `sort_versions(["1.10", "1.2", "1.2.1"]) == ["1.2", "1.2.1", "1.10"]`),
			vis("Nothing to sort", "sort_versions([]) == []"),
			vis("Leading zeros are just numbers", `sort_versions(["1.09", "1.9"]) == ["1.09", "1.9"]`),
			hid("A short version pads with zeros", `sort_versions(["1.2.1", "1.2"]) == ["1.2", "1.2.1"]`),
			hid("Equal versions keep their given order", `sort_versions(["1.2.0", "1.2"]) == ["1.2.0", "1.2"]`),
			hid("Major version wins over everything", `sort_versions(["2.0", "1.99.99"]) == ["1.99.99", "2.0"]`),
		},
		solution: "def sort_versions(versions):\n    n = max((v.count('.') + 1 for v in versions), default=1)\n\n    def key(v):\n        parts = [int(p) for p in v.split('.')]\n        return tuple(parts + [0] * (n - len(parts)))\n\n    return sorted(versions, key=key)\n",
	},
	{
		slug: "ticket-order", title: "Ticket order", difficulty: "medium", topic: "Sorting",
		statement: `A support queue holds ~(id, severity, hours_waiting)~ tickets, severity 1
being the most urgent.

A ticket that has waited **more than 48 hours** is escalated: treat it as one
severity more urgent, though nothing goes above 1.

Order by effective severity first, then by the longest wait, then by id.
Write ~order(tickets)~ returning just the ids.

    order([("a", 2, 10), ("b", 1, 1), ("c", 2, 50)])  ->  ["c", "b", "a"]

Ticket ~c~ escalates to severity 1 and has waited longest, so it goes first.`,
		timeLimitMs: 5000,
		starters:    py("def order(tickets):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Escalation reorders the queue", `order([("a", 2, 10), ("b", 1, 1), ("c", 2, 50)]) == ["c", "b", "a"]`),
			vis("Without escalation, severity decides", `order([("a", 2, 1), ("b", 1, 1)]) == ["b", "a"]`),
			vis("An empty queue", "order([]) == []"),
			hid("Exactly 48 hours does not escalate", `order([("a", 1, 1), ("b", 2, 48)]) == ["a", "b"]`),
			hid("Severity 1 cannot escalate further", `order([("a", 1, 100), ("b", 1, 99)]) == ["a", "b"]`),
			hid("Ids break a complete tie", `order([("z", 1, 5), ("a", 1, 5)]) == ["a", "z"]`),
		},
		solution: "def order(tickets):\n    def key(t):\n        tid, sev, hours = t\n        eff = max(1, sev - 1) if hours > 48 else sev\n        return (eff, -hours, tid)\n\n    return [t[0] for t in sorted(tickets, key=key)]\n",
	},
	{
		slug: "bracket-seeding", title: "Bracket seeding", difficulty: "hard", topic: "Sorting",
		statement: `A knockout bracket is seeded so the two best players can only meet in the
final, the top four only in the semis, and so on.

You build it by doubling: start with seed 1 alone. To go from a bracket of
~size~ to one of ~2 * size~, replace every seed ~s~ with the pair
~s, 2 * size + 1 - s~, keeping them in place.

    1                    (size 1)
    1 2                  (size 2)
    1 4 2 3              (size 4)
    1 8 4 5 2 7 3 6      (size 8)

Write ~pairings(n)~ for a power-of-two ~n~, returning the first-round matches
as a list of ~(a, b)~ pairs read left to right off that final row.

    pairings(4)  ->  [(1, 4), (2, 3)]

A bracket of one player has no matches.`,
		timeLimitMs: 5000,
		starters:    py("def pairings(n):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A bracket of four", "pairings(4) == [(1, 4), (2, 3)]"),
			vis("A bracket of two is one match", "pairings(2) == [(1, 2)]"),
			vis("One player plays nobody", "pairings(1) == []"),
			vis("A bracket of eight", "pairings(8) == [(1, 8), (4, 5), (2, 7), (3, 6)]"),
			hid("Every seed appears exactly once", "sorted([s for p in pairings(16) for s in p]) == list(range(1, 17))"),
			hid("The top seed meets the bottom seed", "pairings(16)[0] == (1, 16)"),
			hid("The second seed lands in the other half of the draw", "[i for i, p in enumerate(pairings(16)) if 2 in p] == [4]"),
		},
		solution: "def pairings(n):\n    order = [1]\n    size = 1\n    while size < n:\n        size *= 2\n        nxt = []\n        for s in order:\n            nxt.append(s)\n            nxt.append(size + 1 - s)\n        order = nxt\n    return [(order[i], order[i + 1]) for i in range(0, len(order) - 1, 2)]\n",
	},
	{
		slug: "merge-streams", title: "Merge streams", difficulty: "medium", topic: "Sorting",
		statement: `Several machines each write a log, and each log is already in time order.
Write ~merge(streams)~ returning every entry in one ordered list.

Entries with the same timestamp are ordered by which stream they came from —
the earlier stream in the list goes first — and two entries from the same
stream keep the order that stream had them in.

    merge([[1, 3], [2, 3]])  ->  [1, 2, 3, 3]

Both 3s are there, and the one from the first stream comes first.`,
		timeLimitMs: 5000,
		starters:    py("def merge(streams):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Interleaves two streams", "merge([[1, 3], [2, 3]]) == [1, 2, 3, 3]"),
			vis("One stream comes back as it was", "merge([[5, 6]]) == [5, 6]"),
			vis("No streams at all", "merge([]) == []"),
			vis("Empty streams contribute nothing", "merge([[], [1], []]) == [1]"),
			hid("A tie goes to the earlier stream", "merge([[2], [1], [1]]) == [1, 1, 2]"),
			hid("Three streams merge cleanly", "merge([[1, 4], [2, 5], [3, 6]]) == [1, 2, 3, 4, 5, 6]"),
		},
		solution: "def merge(streams):\n    tagged = []\n    for i, s in enumerate(streams):\n        for x in s:\n            tagged.append((x, i))\n    tagged.sort(key=lambda p: (p[0], p[1]))\n    return [x for x, _ in tagged]\n",
	},
	{
		slug: "fair-split", title: "Fair split", difficulty: "hard", topic: "Dynamic programming",
		statement: `Two removal vans have to share a load. Write ~split(weights)~ returning
the smallest possible difference between the two vans' totals.

    split([1, 2, 3])  ->  0
    split([8, 1])     ->  7

Every item goes in one van or the other; a van may end up empty. Weights are
whole numbers and never negative.

Thirty items is a realistic load, and trying all 2^30 ways to divide them is
not. What you can afford to track is every total one van could possibly end up
with.`,
		timeLimitMs: 4000,
		starters:    py("def split(weights):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("An even split", "split([1, 2, 3]) == 0"),
			vis("A load that cannot be evened out", "split([8, 1]) == 7"),
			vis("Nothing to move", "split([]) == 0"),
			vis("One item goes in one van", "split([5]) == 5"),
			vis("Fast enough for thirty items", "split([1] * 30) == 0"),
			hid("A near miss", "split([1, 2, 4]) == 1"),
			hid("Zero-weight items change nothing", "split([0, 0, 5]) == 5"),
		},
		solution: "def split(weights):\n    total = sum(weights)\n    reach = {0}\n    for w in weights:\n        reach |= {r + w for r in reach}\n    return min(abs(total - 2 * r) for r in reach)\n",
		tooSlow:  "def split(weights):\n    total = sum(weights)\n    n = len(weights)\n    best = total\n    for mask in range(1 << n):\n        s = 0\n        for i in range(n):\n            if mask >> i & 1:\n                s += weights[i]\n        d = abs(total - 2 * s)\n        if d < best:\n            best = d\n    return best\n",
	},
	{
		slug: "group-leaders", title: "Group leaders", difficulty: "easy", topic: "Sorting",
		statement: `Write ~leaders(entries)~, where each entry is ~(group, name, score)~.
Return a dictionary mapping each group to the name with the highest score in
it. If two names in a group tie, the one earlier in the alphabet leads.

    leaders([("red", "ada", 5), ("red", "bob", 7), ("blue", "cy", 1)])
      ->  {"red": "bob", "blue": "cy"}`,
		timeLimitMs: 5000,
		starters:    py("def leaders(entries):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The highest score leads each group", `leaders([("red", "ada", 5), ("red", "bob", 7), ("blue", "cy", 1)]) == {"red": "bob", "blue": "cy"}`),
			vis("No entries, no groups", "leaders([]) == {}"),
			vis("A group of one", `leaders([("solo", "ada", 3)]) == {"solo": "ada"}`),
			hid("A tie goes to the earlier name", `leaders([("g", "zoe", 5), ("g", "ada", 5)]) == {"g": "ada"}`),
			hid("Negative scores still rank", `leaders([("g", "a", -5), ("g", "b", -1)]) == {"g": "b"}`),
		},
		solution: "def leaders(entries):\n    top = {}\n    for g, name, score in entries:\n        cur = top.get(g)\n        if cur is None or score > cur[0] or (score == cur[0] and name < cur[1]):\n            top[g] = (score, name)\n    return {g: v[1] for g, v in top.items()}\n",
	},
	{
		slug: "foreign-alphabet", title: "A foreign alphabet", difficulty: "medium", topic: "Sorting",
		statement: `A dictionary in another language uses the same letters in a different
order. Write ~sort_words(words, alphabet)~ returning the words in that
language's order.

    sort_words(["ba", "ab"], "ba")  ->  ["ba", "ab"]

Compare letter by letter under the given order. When one word is a prefix of
another, the shorter one comes first — it runs out before the disagreement.

Every letter in every word appears in ~alphabet~.`,
		timeLimitMs: 5000,
		starters:    py("def sort_words(words, alphabet):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Sorts by the given order", `sort_words(["ba", "ab"], "ba") == ["ba", "ab"]`),
			vis("A prefix comes first", `sort_words(["ab", "a"], "ab") == ["a", "ab"]`),
			vis("Nothing to sort", `sort_words([], "abc") == []`),
			hid("The usual alphabet is just one order", `sort_words(["b", "a"], "abcdefghijklmnopqrstuvwxyz") == ["a", "b"]`),
			hid("The disagreement can be late in the word", `sort_words(["zzb", "zza"], "ba" + "zyxwvutsrqponmlkjihgfedc") == ["zzb", "zza"]`),
		},
		solution: "def sort_words(words, alphabet):\n    rank = {c: i for i, c in enumerate(alphabet)}\n    return sorted(words, key=lambda w: [rank[c] for c in w])\n",
	},
	{
		slug: "round-robin-feed", title: "Round-robin feed", difficulty: "easy", topic: "Sorting",
		statement: `A social feed shows one post from each account in turn, then goes round
again. When an account runs out of posts it is simply skipped; the others carry
on.

Write ~feed(accounts)~, where ~accounts~ is a list of lists of posts, returning
the feed in order.

    feed([[1, 2], [3]])  ->  [1, 3, 2]`,
		timeLimitMs: 5000,
		starters:    py("def feed(accounts):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Takes one from each in turn", "feed([[1, 2], [3]]) == [1, 3, 2]"),
			vis("A single account posts in order", "feed([[1, 2, 3]]) == [1, 2, 3]"),
			vis("No accounts, no feed", "feed([]) == []"),
			hid("An exhausted account is skipped, not padded", "feed([[1], [2, 3, 4]]) == [1, 2, 3, 4]"),
			hid("Empty accounts contribute nothing", "feed([[], [1], []]) == [1]"),
		},
		solution: "def feed(accounts):\n    out = []\n    i = 0\n    while True:\n        took = False\n        for a in accounts:\n            if i < len(a):\n                out.append(a[i])\n                took = True\n        if not took:\n            return out\n        i += 1\n",
	},
	{
		slug: "most-improved", title: "Most improved", difficulty: "medium", topic: "Sorting",
		statement: `Two leaderboards, ~before~ and ~after~, list names best-first. Write
~climbers(before, after)~ returning everyone on the new board ordered by how
many places they climbed, biggest climb first, ties broken alphabetically.

Somebody who was not on the old board is treated as having been just off the
end of it — position ~len(before)~.

    climbers(["a", "b", "c"], ["c", "a", "b"])  ->  ["c", "a", "b"]

Only people on the new board are returned; whoever dropped off it is gone.`,
		timeLimitMs: 5000,
		starters:    py("def climbers(before, after):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The biggest climb leads", `climbers(["a", "b", "c"], ["c", "a", "b"]) == ["c", "a", "b"]`),
			vis("A newcomer climbs from just off the board", `climbers(["a"], ["z", "a"]) == ["z", "a"]`),
			vis("An empty new board", `climbers(["a"], []) == []`),
			hid("Nobody moved, so alphabetical", `climbers(["a", "b"], ["a", "b"]) == ["a", "b"]`),
			hid("Someone who dropped off is not reported", `climbers(["a", "b"], ["b"]) == ["b"]`),
		},
		solution: "def climbers(before, after):\n    was = {n: i for i, n in enumerate(before)}\n    gain = {n: was.get(n, len(before)) - i for i, n in enumerate(after)}\n    return sorted(gain, key=lambda n: (-gain[n], n))\n",
	},
	{
		slug: "seat-overflow", title: "Seat overflow", difficulty: "medium", topic: "Simulation",
		statement: `A hall seats people row by row. Each request is ~(name, preferred_row)~,
and requests are honoured in the order they arrive. If the preferred row is
full, the person takes the next row with space; if there is no such row, the
hall cannot seat everyone and the whole plan fails.

Write ~seat(requests, capacity)~, where ~capacity~ lists how many seats each row
has. Return a dictionary from name to row number, or ~None~ if anyone cannot be
seated.

    seat([("a", 0), ("b", 0)], [1, 1])  ->  {"a": 0, "b": 1}

Nobody moves backwards: overflow only ever goes to a later row.`,
		timeLimitMs: 5000,
		starters:    py("def seat(requests, capacity):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Overflow moves to the next row", `seat([("a", 0), ("b", 0)], [1, 1]) == {"a": 0, "b": 1}`),
			vis("Everyone gets their preference", `seat([("a", 1)], [1, 1]) == {"a": 1}`),
			vis("Nobody to seat", "seat([], [1]) == {}"),
			vis("A hall that cannot hold them", `seat([("a", 0)], [0]) is None`),
			hid("Overflow never goes backwards", `seat([("a", 1), ("b", 1)], [5, 1]) is None`),
			hid("Overflow can skip a full row", `seat([("a", 0), ("b", 0), ("c", 0)], [1, 0, 2]) == {"a": 0, "b": 2, "c": 2}`),
		},
		solution: "def seat(requests, capacity):\n    left = list(capacity)\n    out = {}\n    for name, pref in requests:\n        r = pref\n        while r < len(left) and left[r] == 0:\n            r += 1\n        if r >= len(left):\n            return None\n        left[r] -= 1\n        out[name] = r\n    return out\n",
	},
	{
		slug: "league-table", title: "League table", difficulty: "medium", topic: "Sorting",
		statement: `Write ~standings(results)~, where each result is
~(home, away, home_goals, away_goals)~.

A win is 3 points, a draw is 1 point each, a loss is nothing. Goal difference is
goals scored minus goals conceded across every match.

Return the team names ordered by points, then goal difference, then
alphabetically — each highest first except the name.

    standings([("A", "B", 2, 1)])  ->  ["A", "B"]

Every team that appears in a result is in the table, even if it lost everything.`,
		timeLimitMs: 5000,
		starters:    py("def standings(results):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A win puts you top", `standings([("A", "B", 2, 1)]) == ["A", "B"]`),
			vis("A draw is settled by name", `standings([("B", "A", 1, 1)]) == ["A", "B"]`),
			vis("No matches, no table", "standings([]) == []"),
			hid("Goal difference separates equal points", `standings([("A", "C", 5, 0), ("B", "C", 1, 0)]) == ["A", "B", "C"]`),
			hid("Three points beats a better goal difference", `standings([("A", "B", 1, 0), ("C", "D", 9, 0), ("C", "A", 0, 1)]) == ["A", "C", "B", "D"]`),
		},
		solution: "def standings(results):\n    pts = {}\n    gd = {}\n    for h, a, hg, ag in results:\n        for t in (h, a):\n            pts.setdefault(t, 0)\n            gd.setdefault(t, 0)\n        gd[h] += hg - ag\n        gd[a] += ag - hg\n        if hg > ag:\n            pts[h] += 3\n        elif ag > hg:\n            pts[a] += 3\n        else:\n            pts[h] += 1\n            pts[a] += 1\n    return sorted(pts, key=lambda t: (-pts[t], -gd[t], t))\n",
	},
}
