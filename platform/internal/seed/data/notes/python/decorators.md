# Decorators

You've written the same three lines of logging at the top of a dozen functions. You've pasted the same "is the player still alive?" check into every combat action. Copy-paste like that is exactly what decorators erase. A decorator wraps a function to add behavior — logging, validation, timing, caching — *without* editing the function itself. This lesson builds up from functions-as-values to writing your own robust, argument-preserving decorators.

## Functions Are First-Class Objects

The whole idea rests on one fact: in Python, functions are ordinary values. You can store one in a variable, put it in a list, pass it as an argument, and return it from another function — just like an integer or a string.

```python
def attack():
    return "slash!"

action = attack        # store the function itself (no parentheses!)
print(action)          # → <function attack at 0x...>
print(action())        # → slash!  (now we call it)
```

Notice `attack` without parentheses is the function *object*; `attack()` with parentheses *calls* it. That distinction is everything in this lesson. Because functions are values, you can pass them around:

```python
def run_twice(fn):
    fn()
    fn()

def spawn():
    print("enemy spawned")

run_twice(spawn)
# → enemy spawned
# → enemy spawned
```

`run_twice` receives a function and calls it. A function that takes or returns other functions is called a **higher-order function** — and a decorator is just a particular kind of higher-order function.

:::key
`attack` is the function object; `attack()` runs it and gives you the result. Decorators pass functions around as objects, so watch the parentheses closely — a missing or extra pair is the most common decorator bug.
:::

## Closures, Briefly

To wrap a function you usually define an inner function that "remembers" something from the outer one. That memory is called a **closure**: an inner function keeps access to the variables of the outer function even after the outer function has returned.

```python
def make_multiplier(factor):
    def multiply(x):
        return x * factor      # `factor` is remembered from the enclosing scope
    return multiply

double = make_multiplier(2)
triple = make_multiplier(3)

print(double(10))  # → 20
print(triple(10))  # → 30
```

`double` still knows `factor` is 2, and `triple` still knows it's 3, long after `make_multiplier` finished. That captured-variable trick is the machinery underneath every decorator.

:::analogy
A closure is a backpack. When the inner function leaves the outer function, it packs the variables it needs into a backpack and carries them along. `double` walks off with a backpack containing `factor = 2`; `triple` carries `factor = 3`. Same code, different backpack.
:::

## What a Decorator Is

A **decorator** is a function that takes a function and returns a new function — usually one that calls the original but adds something before or after. Here's a decorator that logs every attack:

```python
def log_calls(fn):
    def wrapper():
        print(f"[LOG] calling {fn.__name__}")
        result = fn()
        print(f"[LOG] {fn.__name__} finished")
        return result
    return wrapper

def attack():
    print("dealt 10 damage")

attack = log_calls(attack)   # replace attack with the wrapped version
attack()
# → [LOG] calling attack
# → dealt 10 damage
# → [LOG] attack finished
```

Read the key line carefully: `attack = log_calls(attack)`. We pass the original `attack` into `log_calls`, get back `wrapper`, and rebind the name `attack` to it. Now every call to `attack()` actually runs `wrapper`, which logs, calls the real function via its backpack (`fn`), and logs again.

## The @ Syntax Is Just Sugar

Writing `attack = log_calls(attack)` by hand is tedious and easy to forget. Python gives you `@` syntax that does exactly the same thing, placed right above the function:

```python
@log_calls
def attack():
    print("dealt 10 damage")

# The line above is 100% equivalent to:
# attack = log_calls(attack)
```

That's the entire secret. `@log_calls` on the line before `def attack` means "after defining `attack`, immediately run `attack = log_calls(attack)`." Nothing more mysterious than that.

:::tip
Whenever a decorator confuses you, mentally rewrite `@deco` above `def f` as `f = deco(f)` right after the definition. Every decorator, no matter how fancy, is just that reassignment.
:::

:::predict
Q: What does this print?
```python
def shout(fn):
    def wrapper():
        return fn().upper()
    return wrapper

@shout
def greet():
    return "hello hero"

print(greet())
```
- HELLO HERO*
- hello hero
- <function wrapper>
- error
E: `@shout` makes `greet = shout(greet)`, so calling `greet()` runs `wrapper`, which calls the real greet ("hello hero") and uppercases it → HELLO HERO.
:::

## Wrapping Functions with *args and **kwargs

The `log_calls` above only works for functions that take no arguments. A real decorator has to wrap functions with *any* signature — one argument, five, keyword arguments, whatever. The fix is to make `wrapper` accept `*args` and `**kwargs` and pass them straight through to the original function.

```python
def log_calls(fn):
    def wrapper(*args, **kwargs):
        print(f"[LOG] {fn.__name__} called with {args} {kwargs}")
        result = fn(*args, **kwargs)      # forward everything, unchanged
        print(f"[LOG] {fn.__name__} returned {result}")
        return result
    return wrapper

@log_calls
def deal_damage(target, amount, *, critical=False):
    dmg = amount * (2 if critical else 1)
    return f"{target} takes {dmg}"

print(deal_damage("orc", 15, critical=True))
# → [LOG] deal_damage called with ('orc', 15) {'critical': True}
# → [LOG] deal_damage returned orc takes 30
# → orc takes 30
```

`*args` scoops up any positional arguments into a tuple; `**kwargs` scoops up any keyword arguments into a dict. Forwarding them with `fn(*args, **kwargs)` means the wrapper is signature-agnostic — it works for `deal_damage`, `greet`, or anything else. This `*args, **kwargs` wrapper is the standard template you'll reuse in nearly every decorator.

## Preserving Identity with functools.wraps

There's a subtle cost to wrapping: the returned `wrapper` replaces the original, so the function's name and docstring get clobbered. Introspection lies:

```python
@log_calls
def cast_spell():
    """Cast the equipped spell."""
    ...

print(cast_spell.__name__)  # → wrapper   😱  (we wanted cast_spell)
print(cast_spell.__doc__)   # → None       (lost the docstring)
```

The fix is `functools.wraps`, itself a decorator you apply to your inner `wrapper`. It copies the original function's name, docstring, and other metadata onto the wrapper.

```python
import functools

def log_calls(fn):
    @functools.wraps(fn)              # copy fn's identity onto wrapper
    def wrapper(*args, **kwargs):
        print(f"[LOG] calling {fn.__name__}")
        return fn(*args, **kwargs)
    return wrapper

@log_calls
def cast_spell():
    """Cast the equipped spell."""
    ...

print(cast_spell.__name__)  # → cast_spell  ✅
print(cast_spell.__doc__)   # → Cast the equipped spell.
```

:::warning
Always add `@functools.wraps(fn)` to your wrapper. Without it, debuggers, log output, and documentation tools all report `wrapper` instead of the real function name — a maddening bug when a stack trace shows ten functions all named `wrapper`.
:::

## A Call-Counter / Logging Decorator

Because the wrapper is a closure, it can carry *state* between calls in its backpack. Here's a decorator that counts how many times a function has been called — useful for tracking how often the player attacks:

```python
import functools

def count_calls(fn):
    @functools.wraps(fn)
    def wrapper(*args, **kwargs):
        wrapper.calls += 1
        print(f"[{fn.__name__}] call #{wrapper.calls}")
        return fn(*args, **kwargs)
    wrapper.calls = 0                 # attribute lives on the wrapper itself
    return wrapper

@count_calls
def attack():
    return "slash!"

attack()   # → [attack] call #1
attack()   # → [attack] call #2
attack()   # → [attack] call #3
print(attack.calls)  # → 3
```

We hang a `calls` attribute on `wrapper` and bump it each time. Because there's one `wrapper` object per decorated function, the count is per-function and persists across calls.

## A Decorator That Guards Arguments

Decorators are perfect for **validation**: check a condition before letting the real function run, and reject the call if something's wrong. A classic in games — refuse any action if the player isn't alive.

```python
import functools

def require_alive(fn):
    @functools.wraps(fn)
    def wrapper(player, *args, **kwargs):
        if player["hp"] <= 0:
            raise ValueError(f"{player['name']} is dead and cannot act")
        return fn(player, *args, **kwargs)
    return wrapper

@require_alive
def cast_fireball(player, target):
    return f"{player['name']} torches {target}"

hero = {"name": "Ada", "hp": 30}
print(cast_fireball(hero, "goblin"))   # → Ada torches goblin

ghost = {"name": "Bram", "hp": 0}
print(cast_fireball(ghost, "goblin"))  # → raises ValueError: Bram is dead and cannot act
```

The wrapper inspects the first argument (`player`), enforces the rule, and only then forwards to the real function. Stick `@require_alive` on every combat action and the "must be alive" check lives in exactly one place — no more copy-pasted guards scattered across your codebase.

:::key
A guard decorator validates *before* calling the wrapped function and either forwards the call or raises. Put the rule in the decorator once, apply it everywhere with `@require_alive`, and every action stays protected without repetition.
:::

:::quiz
Q: Why do we add `@functools.wraps(fn)` to a decorator's inner wrapper?
- It makes the decorator run faster
- It copies the original function's name and docstring onto the wrapper so introspection stays honest*
- It's required syntax or the decorator won't work
- It automatically adds `*args` and `**kwargs`
E: `functools.wraps` copies metadata (name, docstring, module) from the original onto the wrapper, so tools and tracebacks report the real function name instead of `wrapper`. It's strongly recommended, not strictly required.
:::

## Recap

- Functions are **first-class objects**: store them, pass them, return them. `f` is the object; `f()` calls it.
- A **closure** is an inner function carrying variables from its enclosing scope in a "backpack" — the mechanism behind decorators.
- A **decorator** takes a function and returns a new (usually wrapping) function.
- `@deco` above `def f` is pure sugar for `f = deco(f)` — rewrite it that way whenever you're stuck.
- Give your wrapper `(*args, **kwargs)` and forward with `fn(*args, **kwargs)` so it works for any signature.
- Always decorate the wrapper with `@functools.wraps(fn)` to preserve the original's name and docstring.
- Wrappers can hold **state** (a call counter) and enforce **rules** (a `require_alive` guard) — one definition, applied everywhere.

**Next up:** Context Managers — using `with` to guarantee setup and cleanup always run, even when things go wrong.
