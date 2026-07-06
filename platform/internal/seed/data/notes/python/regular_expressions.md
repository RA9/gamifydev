# Regular Expressions

A regular expression — regex for short — is a tiny pattern language for describing *shapes* of text. Instead of asking "does this string equal `hp 30`?", regex lets you ask "does this string contain a run of digits after the word `hp`?" That power to match patterns rather than exact text makes regex the go-to tool for parsing chat commands, validating player tags, and pulling numbers out of messy log lines.

You already know how to slice and search strings with methods like `.startswith()` and `.split()`. Regex picks up where those leave off: when the thing you're looking for varies but follows a *rule*.

## What Regex Is and When to Use It

A regex is a small string of symbols that describes a pattern. `\d+` means "one or more digits." `[A-Z]` means "any single capital letter." You hand this pattern to Python's `re` module along with some text, and it tells you where — and whether — the pattern matches.

Regex earns its keep when text is *semi-structured*: it follows a loose rule but not a rigid format. Damage in a game log might appear as `dealt 42 damage` or `dealt 7 damage` — the number varies, but the shape is constant. That's a perfect regex job.

But regex is famously easy to overuse. If a plain string method does the job, use it.

:::warning
Don't reach for regex when a simple method works. To check a prefix, `text.startswith("/")` is clearer than a regex. To split on commas, `text.split(",")` beats a pattern. And never try to parse deeply nested formats like HTML or JSON with regex — use a real parser. Regex is for flat, pattern-shaped text.
:::

:::analogy
Think of regex as a search filter in your game's inventory. Typing an exact item name finds one thing; a filter like "all rare swords under 50 gold" finds a whole category by its *shape*. Regex is that category filter, but for text.
:::

## The `re` Module

Everything lives in Python's built-in `re` module. You `import re`, then call one of a handful of functions. The five you'll use constantly:

**`re.search(pattern, text)`** scans the *whole* string and returns a match object at the first place the pattern appears, or `None` if it never does.

```python
import re

log = "player Zed dealt 42 damage"
match = re.search(r"\d+", log)
print(match.group())  # 42
```

**`re.match` vs `re.fullmatch`.** `re.match` only checks whether the pattern matches at the *start* of the string. `re.fullmatch` requires the pattern to match the *entire* string. This distinction matters enormously for validation.

```python
print(re.match(r"\d+", "42 gold"))      # matches (starts with digits)
print(re.match(r"\d+", "gold 42"))      # None (doesn't start with digits)
print(re.fullmatch(r"\d+", "42 gold"))  # None (whole string isn't digits)
print(re.fullmatch(r"\d+", "42"))       # matches (entire string is digits)
```

**`re.findall(pattern, text)`** returns a *list* of every non-overlapping match, as plain strings. Perfect for grabbing all the numbers in a line.

```python
print(re.findall(r"\d+", "hp 30 mp 5 xp 100"))  # ['30', '5', '100']
```

**`re.sub(pattern, replacement, text)`** replaces every match with something else — great for censoring or cleaning text.

```python
print(re.sub(r"\d+", "###", "code 1234 secret 99"))
# code ### secret ###
```

:::key
Use `search` to find a pattern anywhere, `fullmatch` to validate a whole string, `findall` to collect every occurrence, and `sub` to replace. `match` (start-only) is easy to confuse with `search` — most of the time you actually want `search` or `fullmatch`.
:::

## Core Syntax

Regex patterns are built from a small vocabulary. Here are the pieces you'll use most.

**Character shorthands** stand for whole categories:

- `\d` — any digit (`0`–`9`)
- `\w` — any "word" character: letters, digits, or underscore
- `\s` — any whitespace (space, tab, newline)
- `.` — *any* single character except a newline

```python
print(re.findall(r"\w", "a1_!"))   # ['a', '1', '_']  (! is not a word char)
print(re.findall(r"\s", "a b\tc")) # [' ', '\t']
```

**Character classes** `[...]` let you spell out your own set. `[aeiou]` matches any one vowel; `[A-Z]` matches any capital letter; `[0-9a-f]` matches a hex digit. A `^` at the *start* of the class negates it: `[^0-9]` means "anything that is not a digit."

```python
print(re.findall(r"[A-Z]", "Zed vs Kai"))   # ['Z', 'K']
print(re.findall(r"[^aeiou ]", "goblin"))   # ['g', 'b', 'l', 'n']
```

**Anchors** pin a pattern to a position rather than matching a character. `^` means "start of the string" and `$` means "end of the string." They match nothing themselves — they assert *where* you are.

```python
print(bool(re.search(r"^/", "/attack")))   # True  (starts with a slash)
print(bool(re.search(r"gold$", "500 gold"))) # True (ends with "gold")
```

**Quantifiers** say *how many* of the preceding thing to match:

- `*` — zero or more
- `+` — one or more
- `?` — zero or one (optional)
- `{n}` — exactly `n`
- `{n,m}` — between `n` and `m`

```python
print(re.findall(r"\d+", "hp30 mp5"))    # ['30', '5']  (runs of digits)
print(re.fullmatch(r"[A-Z]{3}", "ZED"))  # matches (exactly 3 capitals)
print(re.fullmatch(r"colou?r", "color")) # matches (u is optional)
```

:::predict
Q: What does this return?
```python
import re
print(re.findall(r"\d+", "hp 30 mp 5"))
```
- ['30', '5'] *
- ['3', '0', '5']
- [30, 5]
- '305'
E: `\d+` matches runs of one-or-more digits, so it captures "30" and "5" as separate strings. Without the `+`, `\d` alone would return each digit individually.
:::

## Capture Groups

Matching tells you a pattern is present; often you also want to *pull a value out* of it. Parentheses `( )` create a **capture group** — a sub-part of the pattern whose matched text you can retrieve separately.

After a successful match, `.group(0)` (or just `.group()`) is the whole match, `.group(1)` is the first parenthesized group, `.group(2)` the second, and so on.

```python
import re

line = "Zed dealt 42 damage"
m = re.search(r"(\w+) dealt (\d+) damage", line)
print(m.group(0))  # Zed dealt 42 damage  (the whole match)
print(m.group(1))  # Zed                  (first group: the name)
print(m.group(2))  # 42                   (second group: the amount)
```

This is how you parse a chat command into its parts. Say players type `/give Kai 100` and you want the action, target, and amount:

```python
cmd = "/give Kai 100"
m = re.match(r"/(\w+) (\w+) (\d+)", cmd)
if m:
    action, target, amount = m.group(1), m.group(2), m.group(3)
    print(f"{action} -> {target}: {int(amount)}")  # give -> Kai: 100
```

`findall` cooperates with groups too: if your pattern has groups, `findall` returns tuples of the captured pieces instead of whole matches.

```python
pairs = re.findall(r"(\w+):(\d+)", "hp:30 mp:5")
print(pairs)  # [('hp', '30'), ('mp', '5')]
```

:::tip
Wrap exactly the sub-parts you want to extract in parentheses, and nothing more. Group the *number* in `dealt (\d+) damage`, not the whole phrase — then `.group(1)` hands you `"42"` directly, ready to convert with `int()`.
:::

## Raw Strings: Why `r"..."` Matters

You've seen every pattern above written as `r"..."` — a **raw string**. This isn't decoration; it's essential for regex, and skipping it causes baffling bugs.

The problem is that both Python *and* regex use backslashes. In a normal Python string, `\n` means a newline, `\t` a tab — the backslash is an escape character. But regex *also* uses backslashes for its shorthands like `\d` and `\w`. When you write a plain `"\d"`, Python first tries to interpret `\d` as an escape sequence before regex ever sees it.

A raw string, marked with the `r` prefix, tells Python to leave backslashes completely alone and pass them straight through to the regex engine.

```python
import re

# Raw string: the regex engine sees exactly \d
print(re.findall(r"\d+", "abc 123"))   # ['123']  correct

# Without r: Python mangles the backslash before regex sees it
print(re.findall("\d+", "abc 123"))    # ['123'] but raises a warning in 3.12
```

In modern Python, forgetting the `r` on `\d` triggers a `SyntaxWarning` (and someday an error), because `\d` isn't a valid Python escape. The truly dangerous cases are sequences that *are* valid escapes — `\b` means a backspace character in a normal string but "word boundary" in regex, so `"\bword"` and `r"\bword"` mean entirely different things.

:::warning
Always write regex patterns as raw strings: `r"\d+"`, never `"\d+"`. Without the `r`, Python may quietly transform your backslashes before the regex engine sees them, and `\b` in particular changes meaning silently. Make `r"..."` a reflex for every pattern.
:::

Now let's put the whole toolkit together to validate a player tag. Say a valid tag is a `#` followed by exactly three uppercase letters and three digits, like `#ZED042`. `fullmatch` plus anchoring gives a strict yes/no:

```python
def valid_tag(tag):
    return re.fullmatch(r"#[A-Z]{3}\d{3}", tag) is not None

print(valid_tag("#ZED042"))  # True
print(valid_tag("#zed042"))  # False (lowercase letters)
print(valid_tag("#ZED42"))   # False (only two digits)
print(valid_tag("#ZED042!"))  # False (trailing character)
```

Because `fullmatch` demands the *entire* string fit the pattern, extra characters like the trailing `!` correctly fail — exactly what you want for validation.

:::quiz
Q: Which function should you use to check that an *entire* string is a valid player tag and nothing more?
- re.fullmatch *
- re.search
- re.findall
- re.sub
E: `re.fullmatch` requires the whole string to match the pattern, so trailing junk fails validation. `re.search` would match even if the tag were buried inside other text.
:::

## Recap

- A regex is a pattern that describes the *shape* of text, ideal for semi-structured strings like log lines and chat commands.
- Reach for plain string methods first; use regex only when the target varies but follows a rule — and never to parse HTML or JSON.
- The `re` module gives you `search` (find anywhere), `match` (match at the start), `fullmatch` (match the whole string), `findall` (collect all matches), and `sub` (replace).
- Core syntax: `\d \w \s` and `.` for character categories, `[...]` for custom classes (`[^...]` to negate), `^ $` anchors, and quantifiers `* + ? {n} {n,m}`.
- Capture groups `( )` let you pull values out with `.group(1)`, `.group(2)`, and so on; `findall` returns tuples when groups are present.
- Always write patterns as raw strings (`r"\d+"`) so Python passes backslashes through untouched — `\b` especially changes meaning without the `r`.

**Next up:** compiling patterns with `re.compile` for reuse, plus named groups and lookahead for the trickier parsing jobs.
