package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// tableProblems: dictionaries and sets — counting things, looking things up,
// and the tie-breaks that decide what "the most" means.
var tableProblems = []seedProblem{
	{
		slug: "loot-tally", title: "Loot tally", difficulty: "easy", topic: "Dictionaries",
		statement: `After a raid you have a pile of drops. Write ~tally(drops, threshold)~
returning every item that dropped at least ~threshold~ times, as a list of
~(item, count)~ pairs.

The list is ordered by count, highest first. Items with the same count are
ordered alphabetically.

    tally(["ore", "gem", "ore", "ore", "gem"], 2)
      ->  [("ore", 3), ("gem", 2)]`,
		timeLimitMs: 5000,
		starters:    py("def tally(drops, threshold):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Orders by count", `tally(["ore", "gem", "ore", "ore", "gem"], 2) == [("ore", 3), ("gem", 2)]`),
			vis("Drops below the threshold are left out", `tally(["ore", "gem", "ore"], 2) == [("ore", 2)]`),
			vis("Nothing dropped", "tally([], 1) == []"),
			hid("A tie is broken alphabetically", `tally(["b", "a"], 1) == [("a", 1), ("b", 1)]`),
			hid("A threshold nothing reaches", `tally(["a", "a"], 5) == []`),
		},
		solution: "def tally(drops, threshold):\n    counts = {}\n    for d in drops:\n        counts[d] = counts.get(d, 0) + 1\n    return sorted(((k, v) for k, v in counts.items() if v >= threshold),\n                  key=lambda p: (-p[1], p[0]))\n",
	},
	{
		slug: "unmatched-item", title: "The unmatched item", difficulty: "easy", topic: "Dictionaries",
		statement: `A warehouse scans every crate on the way in and on the way out. Exactly
one item does not balance: it appears a different number of times in the two
lists.

Write ~unmatched(inbound, outbound)~ returning that item.

    unmatched(["a", "b", "c"], ["a", "c"])  ->  "b"

The lists are in no particular order, and an item may appear many times.`,
		timeLimitMs: 5000,
		starters:    py("def unmatched(inbound, outbound):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the item that never left", `unmatched(["a", "b", "c"], ["a", "c"]) == "b"`),
			vis("Order does not matter", `unmatched(["x", "y"], ["y"]) == "x"`),
			vis("An item can go out without coming in", `unmatched([], ["z"]) == "z"`),
			hid("Repeats count, not just presence", `unmatched(["a", "a", "b"], ["a", "b"]) == "a"`),
			hid("The odd one out is buried in the middle", `unmatched(["p", "q", "r", "s"], ["s", "p", "q"]) == "r"`),
		},
		solution: "def unmatched(inbound, outbound):\n    counts = {}\n    for x in inbound:\n        counts[x] = counts.get(x, 0) + 1\n    for x in outbound:\n        counts[x] = counts.get(x, 0) - 1\n    for k, v in counts.items():\n        if v != 0:\n            return k\n",
	},
	{
		slug: "recipe-cost", title: "Recipe cost", difficulty: "easy", topic: "Dictionaries",
		statement: `Write ~cost(recipe, prices)~, where ~recipe~ maps an ingredient to how
much of it you need and ~prices~ maps an ingredient to its unit price.

Return what the recipe costs. If any ingredient has no price, you cannot cost
the recipe at all — return ~-1~ rather than quietly leaving it out.

    cost({"flour": 2, "salt": 1}, {"flour": 30, "salt": 5})  ->  65
    cost({"saffron": 1}, {"flour": 30})                      ->  -1`,
		timeLimitMs: 5000,
		starters:    py("def cost(recipe, prices):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Costs a full recipe", `cost({"flour": 2, "salt": 1}, {"flour": 30, "salt": 5}) == 65`),
			vis("An unpriced ingredient spoils the total", `cost({"saffron": 1}, {"flour": 30}) == -1`),
			vis("An empty recipe costs nothing", `cost({}, {"flour": 30}) == 0`),
			hid("A missing price is not treated as free", `cost({"flour": 1, "gold": 1}, {"flour": 30}) == -1`),
			hid("Prices for things you do not need are ignored", `cost({"salt": 3}, {"salt": 2, "pepper": 99}) == 6`),
		},
		solution: "def cost(recipe, prices):\n    total = 0\n    for k, q in recipe.items():\n        if k not in prices:\n            return -1\n        total += prices[k] * q\n    return total\n",
	},
	{
		slug: "poll-winner", title: "Poll winner", difficulty: "medium", topic: "Dictionaries",
		statement: `Write ~winner(votes)~ returning the option with the most votes.

When two options tie, the winner is the one that was **voted for first** — the
option that got on the board earliest, not the one that finished the count
first, and not whichever the dictionary happens to yield.

    winner(["blue", "red", "red", "blue"])  ->  "blue"

Both have two votes, and blue was voted for first. An empty poll has no
winner: return ~None~.`,
		timeLimitMs: 5000,
		starters:    py("def winner(votes):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The clear majority wins", `winner(["a", "b", "a"]) == "a"`),
			vis("A tie goes to whoever was voted for first", `winner(["blue", "red", "red", "blue"]) == "blue"`),
			vis("Nobody voted", "winner([]) is None"),
			vis("A single vote decides it", `winner(["only"]) == "only"`),
			hid("First vote, not first to reach the count", `winner(["x", "y", "y", "x"]) == "x"`),
			hid("A three-way tie still resolves", `winner(["c", "b", "a"]) == "c"`),
		},
		solution: "def winner(votes):\n    counts = {}\n    first = {}\n    for i, v in enumerate(votes):\n        counts[v] = counts.get(v, 0) + 1\n        if v not in first:\n            first[v] = i\n    if not counts:\n        return None\n    return min(counts, key=lambda k: (-counts[k], first[k]))\n",
	},
	{
		slug: "wasted-bytes", title: "Wasted bytes", difficulty: "easy", topic: "Sets",
		statement: `A backup holds the same file under many names. Each upload is a
~(name, size, digest)~ triple; two uploads with the same digest hold identical
bytes.

Write ~wasted(uploads)~ returning how many bytes are spent on copies — that is,
everything except the first upload of each digest.

    wasted([("a.txt", 10, "x"), ("b.txt", 10, "x"), ("c.txt", 5, "y")])  ->  10`,
		timeLimitMs: 5000,
		starters:    py("def wasted(uploads):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A second copy is waste", `wasted([("a.txt", 10, "x"), ("b.txt", 10, "x"), ("c.txt", 5, "y")]) == 10`),
			vis("All distinct files waste nothing", `wasted([("a", 1, "p"), ("b", 2, "q")]) == 3 - 3`),
			vis("Nothing uploaded", "wasted([]) == 0"),
			hid("Three copies waste two of them", `wasted([("a", 4, "x"), ("b", 4, "x"), ("c", 4, "x")]) == 8`),
			hid("The name has nothing to do with it", `wasted([("same", 7, "p"), ("same", 7, "q")]) == 0`),
		},
		solution: "def wasted(uploads):\n    seen = set()\n    total = 0\n    for name, size, digest in uploads:\n        if digest in seen:\n            total += size\n        else:\n            seen.add(digest)\n    return total\n",
	},
	{
		slug: "complete-kits", title: "Complete kits", difficulty: "easy", topic: "Dictionaries",
		statement: `A kit is a recipe of parts. Write ~kits(parts, kit)~ returning how many
complete kits you can build from the parts you have.

    kits({"bolt": 10, "nut": 4}, {"bolt": 2, "nut": 1})  ->  4

You run out of nuts first. A part you have none of stops you completely, and a
kit that requires nothing is not a kit — return 0 for it.`,
		timeLimitMs: 5000,
		starters:    py("def kits(parts, kit):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The scarcest part decides", `kits({"bolt": 10, "nut": 4}, {"bolt": 2, "nut": 1}) == 4`),
			vis("A missing part means none at all", `kits({"bolt": 10}, {"bolt": 1, "nut": 1}) == 0`),
			vis("A kit of nothing is not a kit", `kits({"bolt": 10}, {}) == 0`),
			hid("Spare parts you do not need are ignored", `kits({"bolt": 4, "glue": 99}, {"bolt": 2}) == 2`),
			hid("Not quite enough for one", `kits({"bolt": 1}, {"bolt": 2}) == 0`),
		},
		solution: "def kits(parts, kit):\n    if not kit:\n        return 0\n    return min(parts.get(p, 0) // n for p, n in kit.items())\n",
	},
	{
		slug: "attendance-streaks", title: "Attendance streaks", difficulty: "easy", topic: "Dictionaries",
		statement: `Write ~streaks(log)~, where ~log~ maps each student to a list of booleans
— one per day, in order, ~True~ for present.

Return a dictionary mapping each student to their longest run of consecutive
days present.

    streaks({"ada": [True, True, False, True]})  ->  {"ada": 2}

A student who was never there has a streak of 0.`,
		timeLimitMs: 5000,
		starters:    py("def streaks(log):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds the longest run", `streaks({"ada": [True, True, False, True]}) == {"ada": 2}`),
			vis("Never present is a streak of nothing", `streaks({"bob": [False, False]}) == {"bob": 0}`),
			vis("An empty log has no students", "streaks({}) == {}"),
			hid("Every student is reported", `streaks({"a": [True], "b": []}) == {"a": 1, "b": 0}`),
			hid("A run that reaches the last day counts", `streaks({"c": [False, True, True]}) == {"c": 2}`),
		},
		solution: "def streaks(log):\n    out = {}\n    for name, days in log.items():\n        best = run = 0\n        for d in days:\n            run = run + 1 if d else 0\n            if run > best:\n                best = run\n        out[name] = best\n    return out\n",
	},
	{
		slug: "hold-queue", title: "Hold queue", difficulty: "medium", topic: "Dictionaries",
		statement: `A library records holds in the order they were placed, as
~(title, person)~ pairs across every title at once.

Write ~position(holds, title, person)~ returning that person's place in the
queue for that title, counting from 1. Return 0 if they are not waiting for it.

Nobody queues twice for the same book: a repeat request from someone already in
that queue is ignored and does not move anyone.

    position([("Dune", "ada"), ("Emma", "bob"), ("Dune", "bob")], "Dune", "bob")  ->  2`,
		timeLimitMs: 5000,
		starters:    py("def position(holds, title, person):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts only holds on that title", `position([("Dune", "ada"), ("Emma", "bob"), ("Dune", "bob")], "Dune", "bob") == 2`),
			vis("First in line", `position([("Dune", "ada")], "Dune", "ada") == 1`),
			vis("Not waiting at all", `position([("Dune", "ada")], "Dune", "zoe") == 0`),
			vis("Nobody is waiting for anything", `position([], "Dune", "ada") == 0`),
			hid("A repeat hold does not move the queue", `position([("D", "a"), ("D", "a"), ("D", "b")], "D", "b") == 2`),
			hid("A repeat hold keeps the original place", `position([("D", "a"), ("D", "b"), ("D", "a")], "D", "a") == 1`),
		},
		solution: "def position(holds, title, person):\n    seen = set()\n    n = 0\n    for t, p in holds:\n        if t != title or p in seen:\n            continue\n        seen.add(p)\n        n += 1\n        if p == person:\n            return n\n    return 0\n",
	},
	{
		slug: "tag-cloud-sizes", title: "Tag cloud sizes", difficulty: "medium", topic: "Dictionaries",
		statement: `A tag cloud draws each tag at one of five sizes. The rarest tag is size
1, the most common is size 5, and everything between is spread evenly across
that range:

    size = 1 + (count - lowest) * 4 // (highest - lowest)

If every tag has the same count there is no range to spread across, and every
tag is drawn at size 3.

Write ~sizes(counts)~ returning a dictionary from tag to size.

    sizes({"go": 1, "python": 5})  ->  {"go": 1, "python": 5}`,
		timeLimitMs: 5000,
		starters:    py("def sizes(counts):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The extremes take the extreme sizes", `sizes({"go": 1, "python": 5}) == {"go": 1, "python": 5}`),
			vis("All equal means all middling", `sizes({"a": 3, "b": 3}) == {"a": 3, "b": 3}`),
			vis("No tags, no sizes", "sizes({}) == {}"),
			hid("A lone tag has no range either", `sizes({"solo": 9}) == {"solo": 3}`),
			hid("The middle lands where the formula says", `sizes({"a": 0, "b": 2, "c": 4}) == {"a": 1, "b": 3, "c": 5}`),
		},
		solution: "def sizes(counts):\n    if not counts:\n        return {}\n    lo, hi = min(counts.values()), max(counts.values())\n    if lo == hi:\n        return {k: 3 for k in counts}\n    return {k: 1 + (v - lo) * 4 // (hi - lo) for k, v in counts.items()}\n",
	},
	{
		slug: "keep-the-last", title: "Keep the last", difficulty: "medium", topic: "Dictionaries",
		statement: `A playlist has picked up duplicates. Remove them — but keep the **last**
time each track appears, not the first, and leave the survivors in the order
those last appearances happen.

Write ~dedupe(tracks)~.

    dedupe(["a", "b", "a", "c"])  ->  ["b", "a", "c"]

The first ~a~ goes; the second one stays where it is.`,
		timeLimitMs: 5000,
		starters:    py("def dedupe(tracks):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Keeps the later copy in its own place", `dedupe(["a", "b", "a", "c"]) == ["b", "a", "c"]`),
			vis("Nothing to remove", `dedupe(["x", "y"]) == ["x", "y"]`),
			vis("An empty playlist", "dedupe([]) == []"),
			hid("All the same leaves one", `dedupe(["a", "a", "a"]) == ["a"]`),
			hid("Several duplicates at once", `dedupe(["a", "b", "b", "a"]) == ["b", "a"]`),
		},
		solution: "def dedupe(tracks):\n    last = {}\n    for i, t in enumerate(tracks):\n        last[t] = i\n    return [t for i, t in enumerate(tracks) if last[t] == i]\n",
	},
	{
		slug: "roster-diff", title: "Roster diff", difficulty: "easy", topic: "Sets",
		statement: `Write ~diff(before, after)~ comparing two team rosters and returning a
tuple of three sorted lists: who joined, who left, and who stayed.

    diff(["ada", "bob"], ["bob", "cy"])  ->  (["cy"], ["ada"], ["bob"])

Each list is sorted alphabetically. A name repeated in a roster is still just
one person.`,
		timeLimitMs: 5000,
		starters:    py("def diff(before, after):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Joined, left and stayed", `diff(["ada", "bob"], ["bob", "cy"]) == (["cy"], ["ada"], ["bob"])`),
			vis("Nobody moved", `diff(["a"], ["a"]) == ([], [], ["a"])`),
			vis("A team from nothing", `diff([], ["a", "b"]) == (["a", "b"], [], [])`),
			hid("Results are sorted, not in roster order", `diff(["z", "y"], ["y", "z"]) == ([], [], ["y", "z"])`),
			hid("A repeated name is one person", `diff(["a", "a"], []) == ([], ["a"], [])`),
		},
		solution: "def diff(before, after):\n    b, a = set(before), set(after)\n    return (sorted(a - b), sorted(b - a), sorted(a & b))\n",
	},
	{
		slug: "restock-list", title: "Restock list", difficulty: "easy", topic: "Dictionaries",
		statement: `Write ~restock(stock, sales, level)~. Start from the ~stock~ dictionary,
subtract one for every item in the ~sales~ list, then return the sorted names of
everything left at or below ~level~.

An item can go negative — that is a backorder, not an error — and an item sold
that was never in stock starts from zero.

    restock({"tea": 3, "jam": 1}, ["tea", "jam"], 0)  ->  ["jam"]`,
		timeLimitMs: 5000,
		starters:    py("def restock(stock, sales, level):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Reports what ran down to the level", `restock({"tea": 3, "jam": 1}, ["tea", "jam"], 0) == ["jam"]`),
			vis("Nothing sold, nothing to restock", `restock({"tea": 3}, [], 0) == []`),
			vis("Everything is low", `restock({"a": 1, "b": 1}, ["a", "b"], 0) == ["a", "b"]`),
			hid("An item can go into backorder", `restock({"tea": 0}, ["tea"], 0) == ["tea"]`),
			hid("Selling something never stocked", `restock({}, ["ghost"], 0) == ["ghost"]`),
		},
		solution: "def restock(stock, sales, level):\n    left = dict(stock)\n    for item in sales:\n        left[item] = left.get(item, 0) - 1\n    return sorted(k for k, v in left.items() if v <= level)\n",
	},
	{
		slug: "best-shared-day", title: "Best shared day", difficulty: "medium", topic: "Sets",
		statement: `Write ~best_day(days, availability)~, where ~days~ lists the candidate
days in the order they should be considered and ~availability~ maps each person
to the set of days they can make.

Return the day the most people can make. If several days tie, take the one
earliest in ~days~.

    best_day(["mon", "tue"], {"ada": {"tue"}, "bob": {"mon", "tue"}})  ->  "tue"

Days nobody mentions are still candidates — with nobody available. If there are
no days at all, return ~None~.`,
		timeLimitMs: 5000,
		starters:    py("def best_day(days, availability):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Picks the most popular day", `best_day(["mon", "tue"], {"ada": {"tue"}, "bob": {"mon", "tue"}}) == "tue"`),
			vis("A tie goes to the earlier day", `best_day(["mon", "tue"], {"ada": {"mon", "tue"}}) == "mon"`),
			vis("No candidate days", "best_day([], {}) is None"),
			hid("A day nobody can make can still win", `best_day(["mon", "tue"], {}) == "mon"`),
			hid("Days outside the candidate list are ignored", `best_day(["mon"], {"a": {"sun"}, "b": {"mon"}}) == "mon"`),
		},
		solution: "def best_day(days, availability):\n    best = None\n    best_n = -1\n    for d in days:\n        n = sum(1 for a in availability.values() if d in a)\n        if n > best_n:\n            best, best_n = d, n\n    return best\n",
	},
	{
		slug: "pair-up", title: "Pair up", difficulty: "medium", topic: "Dictionaries",
		statement: `A buddy scheme pairs people off. Write ~pair_up(pairs)~ turning a list of
~(a, b)~ pairs into a dictionary where each name maps to their buddy — both
directions.

The scheme only works if everybody has exactly one buddy. If any name turns up
in more than one pair, or is paired with themselves, the whole scheme is
invalid: return ~None~.

    pair_up([("ada", "bob")])  ->  {"ada": "bob", "bob": "ada"}
    pair_up([("ada", "bob"), ("ada", "cy")])  ->  None`,
		timeLimitMs: 5000,
		starters:    py("def pair_up(pairs):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Maps both directions", `pair_up([("ada", "bob")]) == {"ada": "bob", "bob": "ada"}`),
			vis("Somebody in two pairs is invalid", `pair_up([("ada", "bob"), ("ada", "cy")]) is None`),
			vis("Nobody to pair is fine", "pair_up([]) == {}"),
			hid("Pairing with yourself is invalid", `pair_up([("ada", "ada")]) is None`),
			hid("A clash on the second name is caught too", `pair_up([("a", "b"), ("c", "b")]) is None`),
			hid("Several valid pairs", `pair_up([("a", "b"), ("c", "d")]) == {"a": "b", "b": "a", "c": "d", "d": "c"}`),
		},
		solution: "def pair_up(pairs):\n    out = {}\n    for a, b in pairs:\n        if a == b or a in out or b in out:\n            return None\n        out[a] = b\n        out[b] = a\n    return out\n",
	},
}
