# Writing Good Tests

Knowing *how* to write a test is only half the skill. The other half is knowing *what* to test. A test that only ever checks `total([10, 20, 5]) == 35` proves your code works on one happy little list — and tells you nothing about what happens when the list is empty, or full of negatives, or has a single item. This lesson is about writing tests that actually catch bugs, without drowning in repetition.

## Test the Edges, Not Just the Middle

The comfortable, obvious case — a normal list of a few positive numbers — is where bugs *don't* hide. Bugs live at the **edges**: the boundaries where your loop starts, stops, or gets nothing to chew on. When you sit down to test a function, deliberately ask: what are the extremes?

For a function like `average(numbers)`, the edges to worry about are:

- **Empty** — an empty list. Does it crash with a `ZeroDivisionError`? Return `0`? You have to decide, then test it.
- **One item** — the smallest non-empty case. Off-by-one loop bugs love this one.
- **Zero** — a value of `0` in the data. Easy to accidentally treat as "no value."
- **Negative** — negative numbers. Code that assumed everything was positive breaks here.

```python
def average(numbers):
    if not numbers:
        return 0
    return sum(numbers) / len(numbers)


def test_average_of_several_numbers():
    assert average([10, 20, 30]) == 20


def test_average_of_empty_list_is_zero():
    assert average([]) == 0          # the empty edge


def test_average_of_single_item():
    assert average([7]) == 7         # the one-item edge


def test_average_handles_negatives():
    assert average([-4, 4]) == 0     # the negative edge
```

:::key
The happy path proves your code *can* work. The edge cases — empty, zero, one, negative — prove it *keeps* working when reality gets weird. Bugs cluster at boundaries, so aim your tests there.
:::

## One Behavior Per Test, with a Name That Explains

Notice that each test above checks exactly one thing and says so in its name. That's not an accident — it's a rule worth following.

When a test checks **one behavior**, a failure tells you precisely what broke. If `test_average_of_empty_list_is_zero` fails, you know the empty case is the problem, full stop. But if you crammed all four checks into one giant `test_average` function, a failure would only tell you "something about averaging is wrong" — and because the first failed `assert` stops the function, you might not even discover the other broken cases until you fix that one.

The name is documentation. `test_average_of_empty_list_is_zero` reads like a sentence: *the average of an empty list is zero*. Compare that to `test_avg2`. When this test fails six months from now, which name would you rather see in the output?

:::tip
A good test name finishes the sentence "test that…". If you can't name a test clearly, it's often a sign the test is trying to check too many things at once. Split it.
:::

## Testing That Code Raises Errors

Sometimes the *correct* behavior is to fail. If someone asks for the biggest number in an empty list, there is no right answer — the function should refuse. In Python, refusing means raising an exception:

```python
def biggest(numbers):
    if not numbers:
        raise ValueError("biggest() needs at least one number")
    result = numbers[0]
    for n in numbers[1:]:
        if n > result:
            result = n
    return result
```

How do you write a test that passes when the code *raises*? A plain `assert` won't help — the exception would blow up your test before any assertion runs. pytest gives you a tool for exactly this: `pytest.raises`, used as a context manager (`with`).

```python
import pytest

def test_biggest_of_empty_list_raises():
    with pytest.raises(ValueError):
        biggest([])
```

Read it as: "inside this `with` block, I expect a `ValueError` to happen." If `biggest([])` raises `ValueError`, the test **passes**. If it raises nothing — or raises a *different* exception — the test **fails**. You've now pinned down not just what your code returns, but how it behaves when it must reject bad input.

:::warning
Keep the `with pytest.raises(...)` block tiny — ideally just the one call you expect to raise. If you wrap several lines, an unexpected exception from the *wrong* line would still make the test pass, hiding a real bug. One block, one expected failure.
:::

## Many Cases Without Repetition: parametrize

Look back at the `average` tests. Four functions, four nearly identical bodies — only the inputs and expected answers changed. That repetition is a smell. pytest's `@pytest.mark.parametrize` lets you write the test body once and feed it many input/output pairs:

```python
import pytest

@pytest.mark.parametrize("numbers, expected", [
    ([10, 20, 30], 20),
    ([], 0),
    ([7], 7),
    ([-4, 4], 0),
])
def test_average(numbers, expected):
    assert average(numbers) == expected
```

That single function now runs **four times**, once per row in the list. The best part: pytest treats each row as its own test. If the empty-list row fails, you see `test_average[[]-0]` fail specifically, while the other three still pass. You get the clarity of separate tests with the brevity of one.

Parametrize shines for functions like `is_even`, where you want to sweep a range of inputs quickly:

```python
def is_even(n):
    return n % 2 == 0

@pytest.mark.parametrize("n, expected", [
    (0, True),    # zero is even — an easy case to get wrong
    (2, True),
    (3, False),
    (-4, True),   # negatives too
])
def test_is_even(n, expected):
    assert is_even(n) == expected
```

:::example
Without parametrize you'd write four `test_is_even_*` functions. With it, one function and four rows cover `0`, a positive even, an odd, and a negative even — including the tricky "is zero even?" boundary — in six lines. Adding a fifth case is a one-line edit.
:::

## A First Look at Fixtures

As tests grow, several of them often need the *same* setup — the same sample data, the same configured object. Copy-pasting that setup into every test is the same repetition problem, one level up. pytest's answer is a **fixture**: a function that builds the shared thing, which tests request just by naming it as a parameter.

```python
import pytest

@pytest.fixture
def sample_scores():
    return [90, 80, 70, 100]


def test_average_of_scores(sample_scores):
    assert average(sample_scores) == 85


def test_biggest_score(sample_scores):
    assert biggest(sample_scores) == 100
```

Both tests ask for `sample_scores` by putting it in their parameter list, and pytest calls the fixture and hands each test a fresh result. The setup lives in exactly one place. If the sample data ever needs to change, you edit the fixture, not a dozen tests. For now, think of a fixture as "shared Arrange" — reusable setup that keeps your tests DRY. There's much more they can do, but this is the core idea.

## Red, Green, Refactor — and Finding the Bug

There's a rhythm many developers follow called **test-driven development (TDD)**, and it flips the usual order around:

1. **Red** — write a test for the behavior you want, and watch it fail. It *should* fail: you haven't written the code yet.
2. **Green** — write just enough code to make the test pass. Nothing fancy.
3. **Refactor** — now that the test protects you, clean up the code. If it stays green, your cleanup didn't break anything.

Red first feels backwards, but it's the point: a test that has never failed might be testing nothing at all. Seeing it go red proves the test can actually detect a problem; seeing it go green proves you fixed it.

This is also exactly how the **"find the bug"** labs work. You're given code that's already broken and a test that's already **red**. The failing test is a signpost — its `E` diff shows the wrong value your code produced next to the value it *should* produce. Your job is to change the code until the test goes **green**. You're not guessing in the dark; the failing assertion points straight at the defect, and green is the unambiguous signal that you nailed it.

:::analogy
Red-green-refactor is like climbing with a rope. You clip in *first* (write the failing test), then you make the move (write the code) — so if you slip, the rope catches you. Refactoring without tests is free-climbing: exhilarating right up until the moment it isn't.
:::

## Practice

:::predict
Q: You want a test to pass when `biggest([])` raises a `ValueError`. Which approach works?
- `assert biggest([]) == ValueError`
- `with pytest.raises(ValueError): biggest([])` *
- `assert biggest([]) raises ValueError`
- `try: biggest([]) — then nothing`
E: `pytest.raises` used as a `with` block passes only if the code inside raises the expected exception. A plain `assert` can't catch a raise — the exception would abort the test before the assertion runs.
:::

:::quiz
Q: You write a `@pytest.mark.parametrize` test with four rows of inputs. How many tests does pytest run and report?
- One test that stops at the first failing row
- Four separate tests, one per row *
- One test that averages the four results
- Zero — parametrize only documents inputs
E: parametrize runs the test body once per row and reports each as its own test. If one row fails you see exactly which, while the other three still pass independently.
:::

## Recap

- Aim tests at the **edges** — empty, zero, one item, negative — because that's where bugs actually hide, not on the happy path.
- Keep **one behavior per test** and give it a **descriptive name**, so a failure points at exactly one thing and reads like documentation.
- Use `with pytest.raises(SomeError):` to assert that code correctly *raises*; keep the block down to the single call you expect to fail.
- `@pytest.mark.parametrize` runs one test body over many input/output rows — each row is reported as its own test, killing repetition without losing clarity.
- A **fixture** (`@pytest.fixture`) is reusable "shared Arrange": tests request it by name and pytest hands each one a fresh copy.
- **Red → Green → Refactor** is the TDD rhythm; the "find the bug" labs hand you a red test whose failure diff points straight at the defect, and green means you fixed it.

**Next up:** you'll put all of this to work in the testing labs — writing `test_*` functions, driving buggy code from red to green, and reading pytest's output like a map to the fix.
