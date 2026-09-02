# Complexity Classes P and NP

Some problems have fast algorithms. Some problems have no fast algorithm that anyone has ever found — and, remarkably, nobody has managed to prove that no fast algorithm exists. That gap is the most famous unsolved question in computer science. This lesson gives you an honest, plain-language account of P, NP and their relatives, including the definition that gets misquoted more than any other in the field.

## Decision problems and "fast"

Two ground rules before the definitions, because both classes are defined narrowly on purpose.

**These classes are about decision problems** — questions with a yes/no answer. Not "what's the shortest route?" but "is there a route shorter than 500 km?" Almost any problem can be rephrased this way, and the yes/no form makes the theory clean.

**"Fast" means polynomial time.** An algorithm is polynomial if its running time is `O(n^k)` for some fixed constant k: `O(n)`, `O(n²)`, `O(n log n)`, even `O(n¹⁰⁰)` all count. `O(2ⁿ)` and `O(n!)` do not — no constant exponent will bound them.

That's a strange-looking line to draw, since `O(n¹⁰⁰)` is hopeless in practice. But it's drawn there because polynomials behave well: add two polynomials, multiply them, or feed one into another, and you still get a polynomial. Exponentials don't stay put like that. The line separates "grows with a fixed exponent" from "grows *in* the exponent," and that turns out to be the join in the rock.

:::key
**Polynomial time** = `O(n^k)` for some constant k. It's a rough stand-in for "tractable," not a promise of speed. The point of the boundary is that it's stable and machine-independent, not that n¹⁰⁰ is usable.
:::

## P: problems you can solve quickly

**P** is the set of decision problems that can be *solved* by an algorithm running in polynomial time.

Nearly everything you've written lives here. Is this number in the sorted array? (`O(log n)`.) Is this list sorted? (`O(n)`.) Is there a path from A to B in this graph? (Breadth-first search, polynomial.) Is there a route from A to B shorter than 500 km? (Dijkstra's algorithm, polynomial.)

P is, roughly, the class of problems computers are actually good at.

## NP: problems you can check quickly

Here is the definition people get wrong, so read it twice.

**NP** is the set of decision problems where, *if the answer is yes*, someone can hand you a proposed solution and you can **verify** it in polynomial time.

NP is about **verifying**, not solving. The letters stand for *nondeterministic polynomial time*, a piece of theoretical machinery from an older definition that turns out to be equivalent.

:::warning
NP does **not** stand for "non-polynomial." It is not the class of hard problems, and it is not the opposite of P. Every problem in P is also in NP. If you remember one thing from this lesson, make it this.
:::

The proposed solution is called a **certificate** or **witness**. Consider Sudoku on a general n×n board. Solving one can take a very long time. But if someone hands you a completed grid, checking it — every row, every column, every box — is quick and mechanical. Solving is hard; checking is easy. That gap is the whole idea of NP.

:::analogy
NP is the class of problems that are like a jigsaw puzzle. Finding the arrangement may take you all weekend. Deciding whether a *finished* picture is correct takes one glance. Every problem in NP has that shape: hard to find, easy to check.
:::

**Why P is inside NP.** If you can solve a problem in polynomial time, you can verify a proposed answer in polynomial time too — ignore what you were handed, solve it yourself, and compare. So `P ⊆ NP`. Every easy problem is also easy to check.

```text
P    = "I can find the answer quickly"
NP   = "if you show me an answer, I can check it quickly"
P is contained in NP. Whether they are EQUAL is the open question.
```

## NP-hard and NP-complete

Two more terms, and they're often swapped by mistake.

**NP-hard** means: at least as hard as *every* problem in NP. Formally, every problem in NP can be transformed into this one in polynomial time — so a fast algorithm for this one would give you a fast algorithm for all of them. An NP-hard problem does not have to be in NP itself; it can even be a problem no algorithm can solve at all.

**NP-complete** means both things at once: **in NP** *and* **NP-hard**. These are the hardest problems in NP — the ones that stand or fall together. Crack any single NP-complete problem with a polynomial algorithm and you have cracked all of them, and proved P = NP.

```text
P              solvable in polynomial time
NP             a yes-answer is verifiable in polynomial time
NP-hard        at least as hard as everything in NP (may not be in NP)
NP-complete    in NP AND NP-hard - the hardest problems inside NP
```

The first problem shown to be NP-complete was **boolean satisfiability** (SAT): given a logical formula of ANDs, ORs and NOTs over true/false variables, is there any assignment making the whole thing true? A proposed assignment is trivially checkable, and every other NP problem can be encoded into one.

## Concrete NP-complete problems

These are worth recognising by name, because you will bump into disguised versions of them in real work. Each is stated in its decision form.

**Travelling salesman (decision version).** Given cities with distances between them, is there a round trip visiting every city exactly once with total length at most k? Checking a proposed tour means adding up its distances — fast. Finding the best one is not.

**0/1 knapsack.** Given items with weights and values and a bag that holds W kilograms, is there a selection worth at least V? Checking a proposed selection is a sum and a comparison.

**Graph colouring.** Given a graph and a number k, can you colour the nodes with k colours so that no two connected nodes share a colour? Even k = 3 is NP-complete. This is exam timetabling, radio frequency assignment and register allocation in a compiler, all wearing different clothes.

**Longest path.** Given a graph and a number k, is there a simple path (no repeated nodes) of length at least k? Note the cruelty here: *shortest* path is solidly in P, solved by Dijkstra's algorithm every time you open a maps app. Flip one word and it becomes NP-complete.

:::example
Knapsack has a well-known dynamic-programming solution that runs in `O(nW)`, where W is the bag's capacity. Doesn't that make it polynomial? No — and the reason is worth knowing. The input encodes W as a *number*, which takes only about log W digits to write down. So `O(nW)` is exponential in the size of the input text. It's called **pseudo-polynomial**: genuinely useful when W is small, no help at all when W is astronomical.
:::

## The open question, and why it matters

Nobody knows whether P = NP.

If **P = NP**, then every problem whose solutions are quick to check is also quick to solve. Scheduling, protein folding, circuit design, optimal packing — a whole category of currently intractable problems would fall.

If **P ≠ NP**, then some problems are genuinely, permanently hard to solve even though answers are easy to check, and our failure to find fast algorithms reflects reality rather than a lack of cleverness.

Most researchers expect P ≠ NP, largely because enormous effort has gone into finding fast algorithms for NP-complete problems and none has appeared. But expecting is not proving. This is a genuinely open question, and anyone who tells you the answer is settled is mistaken.

**Why it matters for cryptography.** Much of the encryption protecting your bank details rests on a problem being hard: given a large number that's the product of two big primes, find those primes. Multiplying is easy; factoring back is believed to be hard. That asymmetry is what makes the padlock in your browser meaningful.

Factoring is in NP (hand me the two factors and I'll multiply them to check), and it isn't known to be NP-complete — it may sit in an in-between region. But if P = NP, then everything in NP is in P, factoring included, and that style of encryption stops working. A proof of P = NP with an actual efficient algorithm attached would be one of the most disruptive results in the history of computing.

:::tip
"NP-complete" is not a reason to give up. It means no known algorithm solves *every* instance optimally in polynomial time. Real work routes around it constantly: approximation algorithms that get within a few percent, heuristics that are good enough, solvers exploiting structure in real-world instances, or simply noticing that your n is 30 and brute force is fine.
:::

## A brief word on co-NP

One more class, mentioned so the name isn't a mystery when you meet it.

NP is built around verifying **yes** answers. **co-NP** is the mirror: problems where a **no** answer has a short, quickly checkable proof.

The two aren't obviously the same. Take satisfiability. If a formula *is* satisfiable, proving it is easy — show the assignment. If it *isn't* satisfiable, what short proof convinces me? There's no obvious one; you'd seem to need to rule out every assignment. Whether NP = co-NP is another open question, and like P vs NP, most researchers doubt it.

## Check Your Understanding

:::quiz
Q: What does it mean for a problem to be in NP?
- It cannot be solved in polynomial time
- A proposed solution can be verified in polynomial time *
- It requires exponential time to solve
- It is harder than every problem in P
E: NP is defined by fast verification of a yes-answer, not by slow solving. Every problem in P is also in NP, because being able to solve it quickly lets you check an answer quickly.
:::

:::quiz
Q: Which statement about NP-complete problems is correct?
- They are in NP and every NP problem reduces to them *
- They are known to require exponential time
- They are outside NP by definition
- They have been proven impossible to verify quickly
E: NP-complete means in NP *and* NP-hard. Nobody has proven they need exponential time — that is precisely the open P vs NP question.
:::

:::match
Q: Match each class to its definition.
- P | Solvable in polynomial time
- NP | A yes-answer is verifiable in polynomial time
- NP-hard | At least as hard as every problem in NP
- co-NP | A no-answer is verifiable in polynomial time
E: NP-complete is the intersection of NP and NP-hard, and P sits entirely inside NP.
:::

## Talk about it

> Shortest path is in P and easily solved every day; longest simple path is NP-complete. Both are questions about paths in the same graph. In your own words, why do you think one direction is so much harder than the other — what can a shortest-path algorithm safely assume that a longest-path algorithm cannot?

## What's next

You've now met the deepest ideas in the subject. To finish the course we'll go back to ground level and assemble the small mathematical toolkit that everything here quietly depends on. In **The Math You Actually Need** you'll get logarithms, powers of two, the arithmetic series, a little counting and a little probability — each one explained concretely, and each one tied to a complexity result you've already seen.
