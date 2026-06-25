# Objects and Classes

This is the big one. Everything you've learned so far, variables, control flow, methods, comes together here. Java is an **object-oriented** language, and objects are the idea at its very core. Once this clicks, you'll understand *why* Java looks the way it does and you'll start designing programs the way professionals do. Take your time. This lesson rewards patience.

## Blueprints and Instances

Here's the single most important idea in this lesson, so let it land fully.

A **class** is a blueprint. An **object** is a thing built from that blueprint.

Think of an architect's blueprint for a house. The blueprint isn't a house, you can't live in it. But from that one blueprint you can build many actual houses, each with its own address, its own color, its own family inside. The blueprint defines the *structure* (every house has rooms, a door, windows); each real house is a specific *instance* with its own details.

![The Dog class is a blueprint defining name and bark(); from it you create individual Dog objects Rex, Bella, and Max, each with its own name](/images/lessons/java-class-object.svg)

:::analogy
A class is a cookie cutter; objects are the cookies. The cutter (class) defines the shape every cookie shares. Each cookie (object) is a separate, real thing you can decorate differently, one with sprinkles, one with icing, but all the same fundamental shape.
:::

In code, you define the blueprint once with the `class` keyword, then stamp out as many objects as you want. Let's build a `Dog` blueprint and see this in action.

## Fields: What an Object Knows

A class describes two things: what its objects **know** (data) and what they can **do** (behavior). The data lives in variables called **fields** (sometimes called attributes).

```java
public class Dog {
    String name;
    int age;
}
```

This `Dog` blueprint says: every dog has a `name` and an `age`. The fields are the blanks on the blueprint, each individual dog will fill them in with its own values.

To build an actual dog from the blueprint, you use the **`new`** keyword:

```java
Dog rex = new Dog();
rex.name = "Rex";
rex.age = 3;

System.out.println(rex.name + " is " + rex.age);   // Rex is 3
```

`new Dog()` constructs a fresh object in memory and hands it back, which we store in `rex`. The dot (`.`) reaches inside the object to read or set its fields. Crucially, you can make *many* dogs, each independent:

```java
Dog luna = new Dog();
luna.name = "Luna";
luna.age = 5;
// rex is still "Rex", age 3. luna is its own separate dog.
```

:::key
A **class** is the blueprint; an **object** is a specific instance built from it with `new`. Fields hold each object's own data, and every object is independent of the others.
:::

## Methods: What an Object Can Do

Objects don't just hold data, they have behavior too, through **methods** that live inside the class. Because these methods belong to individual objects, they're the **instance methods** we hinted at earlier.

```java
public class Dog {
    String name;
    int age;

    void bark() {
        System.out.println(name + " says Woof!");
    }
}
```

Notice `bark()` can use `name` directly, it has access to the fields of whatever dog it's called on. You call an instance method through an object with the dot:

```java
Dog rex = new Dog();
rex.name = "Rex";
rex.bark();   // Rex says Woof!
```

When you write `rex.bark()`, the method runs *for Rex specifically*, so `name` refers to Rex's name. Call `luna.bark()` and it'll say Luna's name instead. The behavior is shared (defined once), but it operates on each object's own data.

## Constructors: Building Objects Properly

Setting every field by hand after `new` is clumsy and easy to forget. A **constructor** fixes this. It's a special method that runs automatically when you create an object, letting you set everything up in one step.

A constructor has the **same name as the class** and no return type:

```java
public class Dog {
    String name;
    int age;

    // Constructor
    Dog(String name, int age) {
        this.name = name;
        this.age = age;
    }

    void bark() {
        System.out.println(name + " says Woof!");
    }
}
```

Now creating a dog is one clean line, with the values passed right in:

```java
Dog rex = new Dog("Rex", 3);
rex.bark();   // Rex says Woof!
```

Notice the keyword **`this`**. Inside the constructor, the parameter is also called `name`, which would be ambiguous. `this.name` means "the *field* belonging to this object," while plain `name` refers to the parameter. So `this.name = name;` reads as "store the incoming value into this object's name field."

:::tip
`this` always means "the object I'm currently working with." It's how an object refers to itself. You'll use it constantly to tell a field apart from a parameter that shares its name.
:::

:::warning
A constructor must be named *exactly* like the class and must have **no return type**, not even `void`. If you accidentally write `void Dog(...)`, Java treats it as an ordinary method, not a constructor, and it won't run when you say `new Dog(...)`. This is a sneaky bug, so watch for it.
:::

## Encapsulation: Protecting Your Data

There's one more professional habit to learn: **encapsulation**. Right now, anyone can reach in and set `rex.age = -100;`, a negative age, which is nonsense. Good design protects an object's data so it can't be put into an invalid state.

The technique is to make fields **`private`** (off-limits from outside) and provide controlled access through methods, often called **getters** and **setters**.

```java
public class Player {
    private String name;
    private int health;

    public Player(String name) {
        this.name = name;
        this.health = 100;   // every new player starts at full health
    }

    // Getter: lets outsiders read the value
    public String getName() {
        return name;
    }

    public int getHealth() {
        return health;
    }

    // Setter with a guard rule
    public void takeDamage(int amount) {
        health = health - amount;
        if (health < 0) {
            health = 0;   // never let health go below zero
        }
    }
}
```

Now the outside world can't corrupt a player's health directly, it *has* to go through `takeDamage`, which enforces the rule that health never drops below zero. The object guards its own integrity.

```java
public class Main {
    public static void main(String[] args) {
        Player hero = new Player("Pixel");
        hero.takeDamage(30);
        System.out.println(hero.getName() + " has " + hero.getHealth() + " HP");
        // Pixel has 70 HP
    }
}
```

:::analogy
Encapsulation is like an ATM. You can't reach into the vault and grab cash directly, that would be chaos. You interact through a controlled slot (the methods), which enforces rules: the right PIN, sufficient balance. The vault (private data) stays protected while still being usable.
:::

## Check Your Understanding

:::quiz
Q: In object-oriented terms, what is the relationship between a class and an object?
- They are the same thing
- A class is a blueprint; an object is an instance built from it *
- An object is a blueprint; a class is built from it
- A class can only ever make one object
E: A class is the blueprint, and you stamp out as many objects (instances) from it as you want, each with its own data.
:::

:::fill
Q: Complete the keyword used to create a new object from the Dog class.
`Dog rex = ___ Dog("Rex", 3);`
- new *
- class
- this
E: The `new` keyword constructs a fresh object in memory from the class blueprint.
:::

:::predict
```java
public class Cat {
    String name;
    Cat(String name) {
        this.name = name;
    }
    void meow() {
        System.out.println(name + " says Meow");
    }
}
// in main:
Cat c = new Cat("Milo");
c.meow();
```
- Milo says Meow *
- name says Meow
- says Meow
E: The constructor stores "Milo" into the object's `name` field via `this.name`, so `meow()` prints "Milo says Meow".
:::

## Talk about it

> Encapsulation means an object protects its own data and only allows changes through controlled methods. Why might letting any part of a program change any value directly become dangerous as software grows? Think of a real-world system (a bank, a vending machine) and describe the "rules" it enforces to keep its data valid.

## What's next

You now understand the beating heart of Java: classes as blueprints, objects as instances, fields and methods, constructors, `this`, and encapsulation. This is the mental model professional Java developers use every single day. In the final lesson, **Next Steps with Java**, we'll look out at the wider landscape, inheritance, collections, frameworks like Spring and Android, and the project ideas that will turn your new knowledge into real, working software. Pixel is so proud of how far you've come.
