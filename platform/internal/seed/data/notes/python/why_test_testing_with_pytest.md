# Why Test? Testing with pytest

You wrote a function. It works — you tried it once in the REPL and the answer looked right. But "looked right once" and "stays right forever" are two very different promises. This lesson is about the second promise: automated tests that check your code for you, over and over, every time you change something.

## Why Automated Tests Matter

Every program you write will change. You'll add a feature, rename a variable, speed up a slow loop, or fix one bug and accidentally create another. The scary part isn't the change — it's not *knowing* whether the change broke something that used to work. That silent breakage has a name: a **regression**.

A test is a small piece of code that runs your real code and checks the answer. Instead of re-typing inputs into the REPL after every edit, the tests do it — in a fraction of a second, every time you ask. Tests buy you three things:

- **They catch regressions.** Change one line, run the tests, and you instantly see if anything that used to pass now fails.
- **They make refactoring safe.** You can rewrite a messy function with confidence, because a green test suite means the behavior didn't change even though the code did.
- **They document behavior.** A test literally spells out "given *this* input, the answer should be *that*." A new teammate can read your tests to learn what your code is supposed to do — and unlike a comment, a test can't lie, because it runs.

:::analogy
Think of tests like a smoke detector in your house. You don't install one because you expect a fire tonight — you install it so that *if* a fire ever starts, you find out immediately instead of waking up to a disaster. A test suite is a smoke detector for your code: it sits quietly until the moment something breaks, then it screams.
:::

## The Shape of a Test: Arrange, Act, Assert

Almost every good test follows the same three-beat rhythm. Once you see it, you'll see it everywhere. It's called **Arrange–Act–Assert**.

1. **Arrange** — set up the inputs and any data the code needs.
2. **Act** — call the function you're testing, once.
3. **Assert** — check that the result is what you expected.

Here's the whole shape for a `total()` function that adds up a list of numbers:

```python
def total(items):
    result = 0
    for n in items:
        result += n
    return result


def test_total_adds_numbers():
    prices = [10, 20, 5]        # Arrange
    answer = total(prices)      # Act
    assert answer == 35         # Assert
```

Read that test function out loud: "Arrange a list of prices, act by calling `total`, assert the answer is 35." Three lines, three beats. When a test gets confusing, it's usually because those beats got tangled — keep them separate and each test stays easy to read.

:::key
Arrange, Act, Assert. Set up the inputs, call the function once, then check the result. If you can't point at each beat in your test, it's doing too much.
:::

## Writing a Test with plain assert

The heart of every test is the `assert` statement. It's built into Python — no imports, no ceremony. You give `assert` an expression that should be `True`:

```python
assert 2 + 2 == 4      # True → nothing happens, the test keeps going
assert 2 + 2 == 5      # False → raises AssertionError, the test fails
```

That's the entire mechanism. **A passing test raises nothing.** If every `assert` in your test function is true, the function runs to the end quietly and the test is counted as a pass. The moment an `assert` is false, Python raises an `AssertionError`, the function stops, and the test is counted as a fail.

This is exactly how the labs in this course work. You write functions named `test_something` full of `assert` lines, and the checker runs each one. If your test function returns without raising, you pass; if an assertion fails, you see the failure. It's the same plain `assert` you just learned.

```python
def test_total_of_empty_list_is_zero():
    assert total([]) == 0


def test_total_of_single_item():
    assert total([42]) == 42
```

:::tip
You can put several `assert` lines in one test, but the first one that fails stops the function — later asserts in the same test won't even run. That's a good reason to keep each test focused on one idea, so a failure points at exactly one thing.
:::

## Installing and Running pytest

Writing `test_*` functions is a habit; **pytest** is the tool that finds and runs them for you. It's the standard testing tool in the Python world, and installing it takes one command:

```bash
pip install pytest
```

Now, from your project folder, run:

```bash
pytest
```

That's it. You don't tell pytest *which* files or *which* functions to run — it figures that out by looking at names. This is called **test discovery**, and the rules are simple:

- pytest looks for files named `test_*.py` (or `*_test.py`).
- inside those files, it collects functions named `test_*`.
- it runs each one and records whether it raised.

So if you save the code above in a file called `test_total.py` and run `pytest`, it finds `test_total_adds_numbers`, `test_total_of_empty_list_is_zero`, and the rest — automatically, because they follow the naming convention. The same `test_*` naming that our labs expect is the naming pytest expects. What you practice here is the real thing.

:::example
A tiny but complete project is just two things in one folder: your code and your tests. You could even put them in the same file. Save this as `test_total.py`, run `pytest`, and you have a working test suite:

```python
def total(items):
    return sum(items)


def test_total_adds_numbers():
    assert total([10, 20, 5]) == 35


def test_total_of_empty_list_is_zero():
    assert total([]) == 0
```
:::

## Reading pytest's Output

When the tests pass, pytest is refreshingly quiet. Each passing test prints a single green dot:

```text
$ pytest
========================= test session starts =========================
collected 2 items

test_total.py ..                                                [100%]

========================== 2 passed in 0.01s ==========================
```

Those two dots after `test_total.py` are your two passing tests — one dot each. A screen full of dots and a green `passed` line at the bottom is the feeling you're chasing.

Now watch what happens when something is wrong. Suppose `total` had a bug and returned the wrong sum. pytest doesn't just say "failed" — it shows you a red `F` for the failing test, then a detailed report:

```text
$ pytest
========================= test session starts =========================
collected 2 items

test_total.py F.                                                [100%]

============================== FAILURES ===============================
_______________________ test_total_adds_numbers _______________________

    def test_total_adds_numbers():
>       assert total([10, 20, 5]) == 35
E       assert 34 == 35
E        +  where 34 = total([10, 20, 5])

test_total.py:6: AssertionError
======================= 1 passed, 1 failed in 0.02s ===================
```

Read that failure block carefully — it's telling you everything. The `>` marks the exact line that failed. The `E` lines are the diff: `assert 34 == 35` means your code produced `34` but the test expected `35`. It even shows `34 = total([10, 20, 5])` so you can see which call produced the wrong value. You don't have to guess — pytest points straight at the defect.

## Why You Don't Need assertEqual

If you've seen other testing tools, you may have met methods like `assertEqual(a, b)`, `assertTrue(x)`, or `assertGreater(a, b)`. pytest doesn't need any of them. You just write plain `assert`:

```python
assert answer == 35          # not assertEqual(answer, 35)
assert result > 0            # not assertGreater(result, 0)
assert name in names         # not assertIn(name, names)
```

How does pytest produce that helpful `assert 34 == 35` diff from a plain `assert` statement? It rewrites your `assert` lines behind the scenes to inspect both sides of the comparison when the test fails. This is called **assert introspection**. The payoff is huge: one keyword, `assert`, covers every kind of check, and pytest still shows you both the expected and the actual value when things go wrong.

:::warning
`assert` is for *tests*, not for validating user input in real programs. Python can be run with the `-O` (optimize) flag, which strips out every `assert` statement entirely. In test code that never matters — you never run pytest with `-O`. But don't rely on `assert` to guard important logic in production code; use a real `if` and raise an exception there.
:::

## Practice

:::quiz
Q: pytest runs a `test_*` function. All of its `assert` statements are true and the function reaches the end. What does pytest record?
- A failure, because the test didn't return a value
- A pass, because no `AssertionError` was raised *
- Nothing, because you must call a report function
- An error, unless you also import `assertEqual`
E: A passing test simply raises nothing. If every `assert` is true, the function runs to the end quietly and pytest counts it as a pass — no return value or special reporting needed.
:::

:::predict
Q: You run `pytest` and see `test_math.py .F.` in the output. How many tests failed?
- Zero — dots and letters both mean pass
- One — each `F` is one failed test *
- Two — the dots are failures
- Three — every character is a failure
E: A dot is a passing test and an `F` is a failing one. `.F.` is three tests: pass, fail, pass — so exactly one failed.
:::

## Recap

- Automated tests catch regressions, make refactoring safe, and document what your code is supposed to do.
- Every test follows the same shape: **Arrange** the inputs, **Act** by calling the function once, **Assert** the result.
- A test is a `test_*` function full of plain `assert` statements; a passing test raises nothing, and a false `assert` raises `AssertionError`.
- Install with `pip install pytest` and run with `pytest`. Discovery finds files named `test_*.py` and functions named `test_*` automatically — the same names our labs use.
- Green dots mean pass; an `F` and a failure block mean fail. The `E` diff shows the actual value vs. the expected one, pointing straight at the bug.
- Thanks to **assert introspection**, plain `assert` replaces `assertEqual` and friends while still giving you a clear expected-vs-actual diff.

**Next up:** Writing Good Tests — choosing which cases to test, testing that code raises errors, and covering many inputs without repeating yourself.
