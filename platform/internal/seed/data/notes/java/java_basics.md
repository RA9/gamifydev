# Java Basics

Now that you've met Java and run your first program, it's time to learn the raw materials every program is built from: **data** and the **variables** that hold it. Almost everything you'll ever write boils down to storing information, transforming it, and showing it to someone. This lesson gives you that foundation. Take it slowly, type the examples out yourself, and let it sink in.

## Variables: Labeled Boxes for Your Data

A **variable** is a named container that holds a value. You give it a name, you put something inside, and later you can read it back or change it.

```java
int score = 100;
```

Read that left to right. `int` says what *kind* of value this box can hold (a whole number). `score` is the **name** we chose. The `=` means "put this value into the box." And `100` is the value. The line ends with a semicolon, as all Java statements do.

:::analogy
A variable is like a labeled box on a shelf. The label (`score`) tells you what's inside, and you can open the box anytime to read the value or swap it for a new one. The *type* (`int`) is like a rule stamped on the box saying "whole numbers only."
:::

Once you've declared a variable, you can change its value without repeating the type:

```java
int score = 100;
score = 250;   // the box now holds 250 instead of 100
```

## Static Typing: Java Wants to Know in Advance

Notice that you had to say `int` before `score`. Java is a **statically typed** language, which means every variable's type is decided up front and cannot change. A box stamped "whole numbers only" will never accept text.

```java
int age = 30;
age = "thirty";   // ERROR: you can't put text into an int box
```

This feels strict, but it's a feature. Because Java knows the type of everything ahead of time, it can catch a whole category of mistakes *before* your program even runs. Many bugs simply become impossible.

:::warning
A common beginner frustration is trying to assign the wrong type to a variable, like putting a decimal into an `int`. Java will refuse and report an error. When that happens, check that your value matches the box's stamped type.
:::

## The Primitive Types

Java has a small set of built-in **primitive types**, the most basic kinds of data. These four are the ones you'll use constantly:

```java
int count = 42;          // whole numbers: -7, 0, 1000
double price = 19.99;    // decimal numbers: 3.14, -0.5
boolean isReady = true;  // a yes/no value: only true or false
char grade = 'A';        // a single character, in single quotes
```

Let's give each one a clear mental picture:

- **`int`** holds whole numbers, with no decimal point. Counting things, scores, ages.
- **`double`** holds numbers *with* decimals. Prices, measurements, averages.
- **`boolean`** holds exactly one of two values: `true` or `false`. Perfect for questions like "is the door open?"
- **`char`** holds a single character, written in single quotes like `'A'` or `'?'`.

:::tip
Notice `char` uses **single quotes** (`'A'`) while text uses **double quotes** (`"Hello"`). Mixing these up is one of the most common early mistakes. Single quote, single character.
:::

## Strings: Holding Text

A single character is rarely enough, you usually want whole words and sentences. For that, Java gives you the `String` type. Unlike the primitives above, `String` starts with a capital letter, and its values go in **double quotes**.

```java
String name = "Pixel";
String greeting = "Welcome to Java!";
```

One of the most useful things you can do with strings is **concatenation**, which is a fancy word for gluing pieces together using the `+` operator:

```java
String firstName = "Ada";
String lastName = "Lovelace";
String fullName = firstName + " " + lastName;
System.out.println(fullName);   // prints: Ada Lovelace
```

You can even glue numbers onto strings, and Java will convert them to text for you:

```java
int level = 7;
System.out.println("You are on level " + level);   // You are on level 7
```

:::analogy
String concatenation is like building a train. Each `+` couples another car onto the line. `"You are on level " + 7` couples the text car and the number car into one longer train of text.
:::

## Printing Output

You've already seen `System.out.println`. Let's understand it properly, because you'll use it in nearly every program.

```java
System.out.println("This prints with a new line after it.");
System.out.print("This prints ");
System.out.print("on the same line.");
```

- `println` prints your text and then moves to a **new line** ("print line").
- `print` prints your text and stays on the **same line**.

Here's a fuller example bringing it together:

```java
public class Main {
    public static void main(String[] args) {
        String player = "Pixel";
        int score = 1500;
        System.out.println("Player: " + player);
        System.out.println("Score: " + score);
    }
}
```

This prints two lines:
```
Player: Pixel
Score: 1500
```

## Constants, Comments, and Arithmetic

A few finishing touches that make your code clearer and more powerful.

**Constants** are variables that should never change. Mark them with the keyword `final`, and by convention name them in ALL_CAPS:

```java
final double PI = 3.14159;
final int MAX_PLAYERS = 4;
```

If you ever try to reassign a `final` variable, Java stops you. This protects values that are meant to stay fixed.

**Comments** are notes for humans that Java ignores completely. Use them to explain *why* your code does something:

```java
// This is a single-line comment
int health = 100;   // the player starts at full health

/* This is a
   multi-line comment */
```

**Arithmetic** works the way you'd expect, with one twist worth knowing:

```java
int a = 10;
int b = 3;
System.out.println(a + b);   // 13
System.out.println(a - b);   // 7
System.out.println(a * b);   // 30
System.out.println(a / b);   // 3   <- not 3.33!
System.out.println(a % b);   // 1   <- the remainder (modulo)
```

:::warning
When you divide two `int` values, Java throws away the decimal part. `10 / 3` gives `3`, not `3.33`. If you want the decimal, use `double` values instead: `10.0 / 3.0` gives `3.333...`. The `%` operator gives the **remainder**, which is surprisingly handy.
:::

:::key
Variables are typed boxes. Use `int`, `double`, `boolean`, and `char` for primitive data, `String` for text, `final` for constants, and remember that integer division drops the decimal.
:::

## Check Your Understanding

:::fill
Q: Complete the type that stores a whole number.
`___ lives = 3;`
- int *
- String
- boolean
E: `int` stores whole numbers like 3. `String` is for text and `boolean` is for true/false.
:::

:::match
Q: Match each Java type to the value it stores.
- `int` | Whole numbers like 42
- `double` | Decimal numbers like 19.99
- `boolean` | true or false
- `String` | Text like "hello"
E: Each type is a box stamped for a specific kind of value.
:::

:::predict
```java
int x = 7;
int y = 2;
System.out.println(x / y);
```
- 3 *
- 3.5
- 4
E: Both values are `int`, so Java performs integer division and drops the decimal: 7 / 2 = 3, not 3.5.
:::

## Talk about it

> Java forces you to declare a type for every variable before you use it. Some languages don't require this. What do you think are the advantages of being made to say "this box only holds whole numbers"? Can you imagine a situation where that strictness would have saved you from a bug?

## What's next

You can now store data, label it with the right type, do arithmetic, and print results. That's a real toolkit. But so far our programs run straight from top to bottom. Next, in **Control Flow in Java**, you'll teach your programs to make *decisions* and *repeat* actions, the moment your code starts to feel genuinely alive. Pixel can't wait.
