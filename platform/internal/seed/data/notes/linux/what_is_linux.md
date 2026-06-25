# What is Linux?

Welcome to the very beginning of your Linux adventure. Before you type a single command, it helps to know what Linux actually *is* — and why it quietly powers most of the digital world around you. By the end of this lesson you'll be able to explain Linux to a friend without memorizing a thing. Pixel is right here cheering you on, so let's set off.

## First, what is an operating system?

Every computer is a pile of hardware: a processor that does the thinking, memory that holds what it's working on, a disk that stores your files, plus a screen, keyboard, and network card. On their own, those parts don't know how to cooperate. Something has to sit in the middle and coordinate them.

That something is the **operating system**, or **OS**. It's the master program that boots up first and stays running the whole time your machine is on. It decides which program gets to use the processor next, where your files live on disk, and how a keypress travels from your keyboard into the app you're using.

:::analogy
Think of your computer as a busy restaurant. The hardware is the kitchen, the pantry, and the dining room. The operating system is the manager who seats guests, hands tickets to the cooks, and makes sure two waiters don't grab the same plate. Without the manager, you just have a room full of equipment and a lot of chaos.
:::

Windows and macOS are operating systems you've probably met already. **Linux** is another one — and it's the one developers reach for again and again.

## The kernel at the heart of it

When people say "Linux," they're really pointing at the **kernel** — the core piece that talks directly to the hardware. The kernel is the part Linus Torvalds first released in 1991 as a student project, and it has been refined by thousands of contributors ever since.

The kernel handles the deep, unglamorous work: scheduling which task runs when, handing out slices of memory, and shuttling data to and from disks and the network. Everything else you interact with — the menus, the apps, the terminal — sits *on top of* the kernel and asks it for help.

:::key
The **kernel** is the heart of Linux. It's the layer that controls the hardware so your programs don't have to. When you run a command later, it's ultimately the kernel doing the heavy lifting underneath.
:::

## Open source: built in the open, by everyone

Here's what makes Linux genuinely special: it's **open source**. The complete source code — the human-readable instructions that make it work — is published for anyone to read, study, change, and share. You don't need permission, and you don't pay a license fee.

That openness changed everything. Instead of one company quietly building an OS behind closed doors, a worldwide community improves Linux together. A fix written by a developer in Nairobi can land alongside one from Berlin and one from São Paulo, all in the same week. Because anyone can inspect it, bugs and security holes get spotted fast.

:::tip
Open source isn't only "free as in no cost" — it's "free as in freedom." You're free to see how things work under the hood. That transparency is exactly why curious developers fall in love with Linux: nothing is hidden from you.
:::

## Distributions: many flavors, one family

Because Linux is open, anyone can bundle the kernel together with useful tools, a way to install software, and a friendly setup, then share the whole package. Each of these packages is called a **distribution**, or **distro** for short.

You may have heard some names already:

- **Ubuntu** — beginner-friendly and hugely popular, great for desktops and servers.
- **Debian** — rock-solid and stable; Ubuntu is actually built on top of it.
- **Fedora** — fast-moving and modern, often showing off the newest features.
- **Alpine** — tiny and lightweight, beloved for cloud containers.

They look a little different, but they all share the same Linux heart. Learn the skills in this path on one distro, and they carry over to the rest.

:::quiz
Q: What do we mean by a Linux "distribution"?
- A bundle of the Linux kernel plus tools and software, packaged together *
- A single program that replaces the kernel
- A company that sells Windows licenses
- A type of computer keyboard
E: A distribution (distro) packages the Linux kernel together with tools and software so it's ready to use. Ubuntu, Debian, and Fedora are all distros built around the same kernel.
:::

## Linux is everywhere

You might think Linux is rare. The opposite is true — it's likely the most widely used operating system on Earth, even if it's often invisible.

- **Servers and the cloud.** The vast majority of websites and cloud services run on Linux. When you load almost any app, a Linux server is answering.
- **Android phones.** Android is built on the Linux kernel, so billions of pockets carry Linux every day.
- **Supercomputers.** Effectively all of the world's fastest supercomputers run Linux.
- **Tiny devices.** Raspberry Pi boards, smart TVs, routers, and cars often run Linux too.

:::example
Open a new browser tab and visit your favorite website. Behind the scenes, your request almost certainly traveled to a Linux server, which assembled the page and sent it back — all in a fraction of a second. You've been using Linux for years without realizing it.
:::

## Why developers learn Linux

If Linux runs the servers and the cloud, then learning Linux means learning the environment your code will actually *live* in. Developers use it to deploy websites, run databases, automate repetitive work, and control machines they may never physically touch. It's a foundational skill — the kind that quietly makes every other skill easier.

## GUI vs the command line

Most people use computers through a **GUI** — a graphical user interface, with windows, icons, and a mouse. Linux has beautiful GUIs too. But its real superpower is the **command line**: a place where you *type* instructions instead of clicking them.

The command line feels old-fashioned at first, yet it's faster, more precise, and far easier to automate. Once you can describe a task in words, you can save those words and run them again a thousand times. That's the door we're about to open.

:::fill
Q: Complete the sentence about how Linux earned its reputation.
Linux is `___` source, meaning anyone can read and improve its code.
- open *
- closed
- hidden
E: Linux is *open* source — its code is published for everyone to study, change, and share.
:::

## Talk about it

> You've now seen that Linux is an operating system, that its kernel talks to the hardware, and that it runs nearly everything from phones to supercomputers. In your own words, how would you explain to a friend why so many developers choose Linux — and what the difference is between clicking in a GUI and typing on the command line?

## What's next

Now that you know *what* Linux is, it's time to actually talk to it. In the next lesson you'll meet the **command line** — the typed conversation between you and your computer — and you'll run your very first commands. Take a breath; the fun part starts now.
