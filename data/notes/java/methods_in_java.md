# Methods in Java

As your programs grow, you'll find yourself writing the same steps over and over: greeting a user, calculating a total, checking whether a move is valid. Copying that logic everywhere is tedious and error-prone, change it in one place and you have to remember to change it in ten others. **Methods** are the cure. A method is a named, reusable block of behavior you can call whenever you need it. This lesson is where your code starts to feel truly *organized*.

## What a Method Is

A method bundles a set of instructions under a single name. Once defined, you can run all those instructions just by calling the name, as many times as you like.

```java
public class Main {
    static void greet() {
        System.out.println("Hello there!");
    }

    public static void main(String[] args) {
        greet();   // runs the greet method
        greet();   // runs it again
    }
}
```

That `greet();` is a **method call**. When Java reaches it, it jumps into the `greet` method, runs everything inside, and then returns to where it left off. This prints "Hello there!" twice without us repeating the print statement.

:::analogy
A method is like a button on a microwave. Someone wired up all the steps once (set power, run timer, beep). You don't redo that work, you just press "Popcorn" whenever you want it. Defining a method is wiring the button; calling it is pressing it.
:::

You've actually been using a method all along: `main` is itself a method, the special one Java runs first.

## Parameters: Giving a Method Input

A method that always does exactly the same thing is limited. Often you want to hand it some information to work with. Those inputs are called **parameters**.

```java
static void greet(String name) {
    System.out.println("Hello, " + name + "!");
}
```

Now `greet` expects a `String` called `name`. When you call it, you pass in an **argument**, the actual value:

```java
greet("Pixel");   // Hello, Pixel!
greet("Ada");     // Hello, Ada!
```

The same button now adapts to whatever you feed it. You can have several parameters, separated by commas:

```java
static void describe(String item, int quantity) {
    System.out.println("You have " + quantity + " " + item + "(s).");
}

// called like:
describe("apple", 3);   // You have 3 apple(s).
```

:::tip
**Parameter** vs **argument** trips people up. The *parameter* is the placeholder in the method definition (`String name`). The *argument* is the real value you supply when calling (`"Pixel"`). Same idea, different moments.
:::

## Return Types: Getting a Result Back

So far our methods just *did* something (printed text) and handed nothing back. Often you want a method to *compute* a value and give it to you. That's where **return types** come in.

The word right before the method name declares what type it returns. Use `return` to send the value back:

```java
static int add(int a, int b) {
    return a + b;
}

public static void main(String[] args) {
    int sum = add(4, 9);
    System.out.println(sum);   // 13
}
```

Here `add` has a return type of `int`. It takes two numbers, adds them, and `return`s the result. The caller captures that result in the variable `sum`.

When a method returns *nothing*, its return type is the special word **`void`**. Our `greet` method was `void` because it only printed, it didn't hand back a value.

```java
static void sayBye() {   // void: returns nothing
    System.out.println("Goodbye!");
}
```

:::analogy
A `void` method is like a vending machine that just plays a jingle, you get an effect but nothing in your hand. A method with a return type is a vending machine that drops a snack: you walk away holding something you can use.
:::

:::warning
If a method declares a return type like `int`, every path through it **must** return an `int`. Forgetting to return a value, or returning the wrong type, is an error Java catches at compile time. A `void` method, by contrast, simply ends when it reaches the bottom.
:::

## Method Signatures, `static`, and Overloading

A method's **signature** is the combination of its name and its parameter list. This is how Java tells methods apart. Together with the return type and body, the signature defines the method completely.

You've seen the word `static` on every method here. A gentle introduction: a **static** method belongs to the class itself, so you can call it without creating an object first. That's why our examples work right inside `main`. Later, when you learn about objects, you'll meet **instance** methods, which belong to individual objects and often need one to exist before you call them. For now, `static` is the simpler world, and it's perfectly fine to stay here while you're learning.

:::analogy
A `static` method is like a public service hotline, anyone can dial it directly, no membership needed. An instance method is like a phone *inside* someone's house: you need that specific house (object) to exist before you can ring it. We'll build houses in the next lesson.
:::

Because Java identifies methods by their signature, you can have several methods with the **same name** as long as their parameters differ. This is called **overloading**, and it's surprisingly natural:

```java
static int multiply(int a, int b) {
    return a * b;
}

static double multiply(double a, double b) {
    return a * b;
}
```

Both are called `multiply`, but Java picks the right one based on the arguments you pass. `multiply(2, 3)` uses the first; `multiply(2.5, 4.0)` uses the second. The name stays meaningful while the method handles different types.

## Why Decomposition Matters

Here's the deeper reason methods matter, beyond avoiding copy-paste. Breaking a big problem into small, well-named methods is called **decomposition**, and it's one of the most important skills in all of programming.

Compare these two ways of writing the same program. First, everything crammed into `main`:

```java
public static void main(String[] args) {
    // 40 lines of tangled logic doing five different things...
}
```

Now the same work, decomposed:

```java
public static void main(String[] args) {
    String name = askForName();
    int score = playRound();
    saveResult(name, score);
    printSummary(name, score);
}
```

Even without seeing inside those methods, you can *read* the second version like a sentence and understand what the program does. Each method has one clear job. When a bug appears in scoring, you know exactly where to look: `playRound`. This clarity is what separates code that's a joy to maintain from code that becomes a nightmare.

:::key
Methods are named, reusable blocks of behavior. Give them inputs through **parameters**, get results through a **return type** (or use **void** for none), and use **decomposition** to break big problems into small, clearly-named pieces.
:::

## Check Your Understanding

:::predict
```java
static int triple(int n) {
    return n * 3;
}

public static void main(String[] args) {
    System.out.println(triple(4));
}
```
- 12 *
- 7
- 43
E: `triple(4)` runs the method with n = 4 and returns 4 * 3 = 12, which is then printed.
:::

:::quiz
Q: What does a return type of `void` mean for a method?
- The method returns a whole number
- The method returns no value *
- The method cannot be called
- The method must take no parameters
E: `void` means the method does its work but hands back no value. Methods that compute a result declare a real type like `int` instead.
:::

:::fill
Q: Complete the call that passes the argument "Sam" to the greet method.
`greet(___);`
- "Sam" *
- String name
- void
E: When calling a method you pass the actual value (the *argument*), here the text `"Sam"`. `String name` is the parameter in the definition, not what you pass at the call.
:::

## Talk about it

> Good programmers obsess over breaking problems into small, well-named methods. Think about a complex task you know well, cooking a meal, planning a trip. How would you split it into named "methods," each doing one job? Why might that make the whole task easier to explain to someone else?

## What's next

You can now package behavior into tidy, reusable methods, and you've had a first taste of `static` versus instance. That last idea points straight at the big one. Next, in **Objects and Classes**, you'll learn to bundle data *and* methods together into objects, the very heart of object-oriented programming and the reason Java is built the way it is. This is the lesson everything has been building toward. Pixel is thrilled for you.
