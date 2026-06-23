# To-Do List (CLI)

Every developer needs a place to track their tasks — so let's build one from scratch. In this capstone you'll create a real command-line To-Do List manager in Python, the kind of tool you could actually use every day.

:::project
By the end you'll have a working terminal app that lets you add tasks, list them, mark them complete, and delete them through a menu — with everything saved to a file so your tasks are still there the next time you run it.
:::

You'll need Python 3 installed and a terminal to run your program in. Create a new file called `todo.py` and build it up step by step. If you ever want to test a small snippet on its own before dropping it in, you can prototype it in the in-app [Code Lab](#code).

## Step 1 — Store your tasks and say hello

Every program needs a starting point. We'll keep all the tasks in a single list called `tasks`, and print a friendly welcome banner so the app feels alive the moment it runs. Right now the list is empty — but it's the home where every task will live.

```python
tasks = []

print("=" * 30)
print("   📋  MY TO-DO LIST")
print("=" * 30)
print(f"You have {len(tasks)} task(s) right now.")
```

## Step 2 — Show the task list

Before we can add tasks, we need a way to *see* them. Let's write a function that prints every task with a number next to it so the user can refer to tasks by their position. If the list is empty, we'll say so clearly instead of printing nothing — a blank screen always confuses people.

```python
def show_tasks(tasks):
    print("\n--- Your Tasks ---")
    if not tasks:
        print("(no tasks yet — add one!)")
        return
    for index, task in enumerate(tasks, start=1):
        print(f"{index}. {task}")
```

## Step 3 — Add a task

Now the fun part: putting tasks *into* the list. This function asks the user what they want to do, then appends their answer to the `tasks` list. We use `.strip()` to trim accidental spaces, and we politely refuse to add an empty task.

```python
def add_task(tasks):
    text = input("What do you need to do? ").strip()
    if text == "":
        print("⚠️  Empty task ignored.")
        return
    tasks.append(text)
    print(f"✅ Added: {text}")
```

## Step 4 — Build the menu loop

Time to tie it together. A menu loop keeps the program running, shows the options, reads the user's choice, and calls the right function. We use a `while True` loop so the app keeps going until the user chooses to quit with `break`.

```python
def menu():
    tasks = []
    while True:
        print("\n--- MENU ---")
        print("1. Add task")
        print("2. List tasks")
        print("3. Quit")
        choice = input("Choose an option: ").strip()

        if choice == "1":
            add_task(tasks)
        elif choice == "2":
            show_tasks(tasks)
        elif choice == "3":
            print("👋 Goodbye!")
            break
        else:
            print("❓ Unknown option, try again.")

menu()
```

## Step 5 — Mark a task complete

A to-do list isn't useful unless you can check things off. To track completion we'll upgrade each task from a plain string into a small dictionary like `{"text": "...", "done": False}`. Update `add_task` to store dictionaries, update `show_tasks` to display a `[x]` or `[ ]` checkbox, and add a `complete_task` function. Add a menu option `4` that calls it.

:::tip
Switching from strings to dictionaries lets each task carry extra information beyond its name. This is a pattern you'll reuse constantly in real programs.
:::

```python
def add_task(tasks):
    text = input("What do you need to do? ").strip()
    if text == "":
        print("⚠️  Empty task ignored.")
        return
    tasks.append({"text": text, "done": False})
    print(f"✅ Added: {text}")

def show_tasks(tasks):
    print("\n--- Your Tasks ---")
    if not tasks:
        print("(no tasks yet — add one!)")
        return
    for index, task in enumerate(tasks, start=1):
        box = "[x]" if task["done"] else "[ ]"
        print(f"{index}. {box} {task['text']}")

def complete_task(tasks):
    show_tasks(tasks)
    if not tasks:
        return
    choice = input("Number to mark complete: ").strip()
    if choice.isdigit() and 1 <= int(choice) <= len(tasks):
        tasks[int(choice) - 1]["done"] = True
        print("🎉 Marked complete!")
    else:
        print("⚠️  That isn't a valid task number.")
```

## Step 6 — Delete a task

Sometimes a task no longer matters and you just want it gone. This function shows the list, asks for a number, and removes that task with `pop()`. We reuse the same number-validation pattern from Step 5 so bad input never crashes the program. We'll wire this to menu option `5`.

```python
def delete_task(tasks):
    show_tasks(tasks)
    if not tasks:
        return
    choice = input("Number to delete: ").strip()
    if choice.isdigit() and 1 <= int(choice) <= len(tasks):
        removed = tasks.pop(int(choice) - 1)
        print(f"🗑️  Deleted: {removed['text']}")
    else:
        print("⚠️  That isn't a valid task number.")
```

## Step 7 — Save and load so tasks persist

The final upgrade makes your app *remember*. We'll use Python's built-in `json` module to write tasks to `tasks.json` whenever they change, and to load them back when the app starts. Save after every add, complete, and delete. Here is the complete, runnable program — copy it into `todo.py`, run `python todo.py`, and you've built your capstone.

:::warning
Run the program from the same folder each time. `tasks.json` is created next to it, so launching from a different directory will look like your tasks vanished.
:::

```python
import json

FILENAME = "tasks.json"

def load_tasks():
    try:
        with open(FILENAME, "r") as f:
            return json.load(f)
    except (FileNotFoundError, json.JSONDecodeError):
        return []

def save_tasks(tasks):
    with open(FILENAME, "w") as f:
        json.dump(tasks, f, indent=2)

def show_tasks(tasks):
    print("\n--- Your Tasks ---")
    if not tasks:
        print("(no tasks yet — add one!)")
        return
    for index, task in enumerate(tasks, start=1):
        box = "[x]" if task["done"] else "[ ]"
        print(f"{index}. {box} {task['text']}")

def add_task(tasks):
    text = input("What do you need to do? ").strip()
    if text == "":
        print("⚠️  Empty task ignored.")
        return
    tasks.append({"text": text, "done": False})
    save_tasks(tasks)
    print(f"✅ Added: {text}")

def complete_task(tasks):
    show_tasks(tasks)
    if not tasks:
        return
    choice = input("Number to mark complete: ").strip()
    if choice.isdigit() and 1 <= int(choice) <= len(tasks):
        tasks[int(choice) - 1]["done"] = True
        save_tasks(tasks)
        print("🎉 Marked complete!")
    else:
        print("⚠️  That isn't a valid task number.")

def delete_task(tasks):
    show_tasks(tasks)
    if not tasks:
        return
    choice = input("Number to delete: ").strip()
    if choice.isdigit() and 1 <= int(choice) <= len(tasks):
        removed = tasks.pop(int(choice) - 1)
        save_tasks(tasks)
        print(f"🗑️  Deleted: {removed['text']}")
    else:
        print("⚠️  That isn't a valid task number.")

def menu():
    tasks = load_tasks()
    print("=" * 30)
    print("   📋  MY TO-DO LIST")
    print("=" * 30)
    while True:
        print("\n--- MENU ---")
        print("1. Add task")
        print("2. List tasks")
        print("3. Mark complete")
        print("4. Delete task")
        print("5. Quit")
        choice = input("Choose an option: ").strip()

        if choice == "1":
            add_task(tasks)
        elif choice == "2":
            show_tasks(tasks)
        elif choice == "3":
            complete_task(tasks)
        elif choice == "4":
            delete_task(tasks)
        elif choice == "5":
            print("👋 Goodbye! Your tasks are saved.")
            break
        else:
            print("❓ Unknown option, try again.")

menu()
```
