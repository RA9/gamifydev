# Bank Account

A bank account is the perfect way to learn object-oriented programming: it has data it protects (your balance) and rules for changing it (deposits and withdrawals). You'll build a `BankAccount` class and put it to work.

:::project
You'll have a runnable Java program with a `BankAccount` class that supports deposits, guarded withdrawals, and a balance getter, plus a `main` method that demonstrates the whole thing.
:::

To follow along you need the JDK installed. Put everything in a file called `Main.java`, then compile and run it with `javac Main.java` followed by `java Main`.

## Step 1 — Create the Main skeleton

Start with a `Main` class and a `main` method — the entry point Java runs first. Print something so you can confirm it compiles and runs.

```java
public class Main {
    public static void main(String[] args) {
        System.out.println("Bank is open!");
    }
}
```

## Step 2 — Define the BankAccount class

Below `Main`, add a `BankAccount` class with a `private` balance field so it can't be changed from outside, and a constructor that sets the opening balance. `private` is what makes the data protected.

```java
class BankAccount {
    private double balance;

    public BankAccount(double startingBalance) {
        balance = startingBalance;
    }
}
```

## Step 3 — Add a deposit method

Add a `deposit` method that increases the balance. Guarding against zero or negative amounts keeps the account honest.

```java
public void deposit(double amount) {
    if (amount <= 0) {
        System.out.println("Deposit must be positive.");
        return;
    }
    balance += amount;
    System.out.println("Deposited " + amount);
}
```

## Step 4 — Add a guarded withdraw method

Withdrawals must never overdraw the account. Check that there's enough balance first; if not, print a warning and leave the balance untouched.

```java
public void withdraw(double amount) {
    if (amount > balance) {
        System.out.println("Insufficient funds! Withdrawal denied.");
        return;
    }
    balance -= amount;
    System.out.println("Withdrew " + amount);
}

public double getBalance() {
    return balance;
}
```

## Step 5 — Demonstrate it in main (complete program)

Finally, use the class from `main`: open an account, deposit, try a withdrawal that's too big, do a valid one, and print the final balance with the getter. Here's the full, runnable program.

```java
public class Main {
    public static void main(String[] args) {
        BankAccount account = new BankAccount(100.0);

        account.deposit(50.0);
        account.withdraw(200.0);
        account.withdraw(30.0);

        System.out.println("Final balance: " + account.getBalance());
    }
}

class BankAccount {
    private double balance;

    public BankAccount(double startingBalance) {
        balance = startingBalance;
    }

    public void deposit(double amount) {
        if (amount <= 0) {
            System.out.println("Deposit must be positive.");
            return;
        }
        balance += amount;
        System.out.println("Deposited " + amount);
    }

    public void withdraw(double amount) {
        if (amount > balance) {
            System.out.println("Insufficient funds! Withdrawal denied.");
            return;
        }
        balance -= amount;
        System.out.println("Withdrew " + amount);
    }

    public double getBalance() {
        return balance;
    }
}
```
