# History of the Web

Before you build for the web, it helps to know what the web actually *is*, where it came from, and how it works under the hood. Knowing the story makes a lot of "why does it work this way?" questions suddenly click into place.

By the end of this lesson you'll be able to explain the web's history, trace exactly what happens when you type a URL and hit Enter, and describe the three languages every frontend developer uses every single day.

## First, the web is not the internet

People use these words interchangeably, but they're different things:

- The **internet** is the *infrastructure* — the global network of cables, routers, and computers that lets machines talk to each other. It came first, in the 1970s and 80s.
- The **web** (the "World Wide Web") is *one thing you can do* on the internet — a way of linking and viewing documents in a browser. It was invented later, on top of the internet.

:::analogy
The internet is the road network; the web is one kind of vehicle that drives on it. Email, video calls, file transfers, and online games are other "vehicles" using the same roads.
:::

## Where the internet came from: ARPANET

In 1969, a US research project called **ARPANET** connected a handful of university computers so they could share data even if part of the network failed. Its big idea was **packet switching** — chopping messages into small "packets" that each find their own route to the destination and get reassembled on arrival. That idea still powers the internet today. Over the 1970s and 80s, a common set of rules called **TCP/IP** let many separate networks join into one global "inter-network" — the internet.

## The idea of the web (1989–1990)

At **CERN**, a European physics lab, a scientist named **Tim Berners-Lee** had a problem: researchers' information was scattered across incompatible computers. His solution was elegantly simple — link documents together using **hypertext**, so you could jump from one to another with a click.

By 1990 he'd built the three pillars the web still runs on today:

- **HTML** — a language to write the documents.
- **HTTP** — a set of rules for computers to request and send those documents.
- The first **browser** and **web server** — software to view and to host pages.

:::key
The web was invented to **share and link information**. Documents connected by clickable links is still the beating heart of every website you use.
:::

## The first website (1991)

In 1991 the first-ever website went live at **info.cern.ch**. It was just text and links explaining what the web was. But it proved the idea worked. Crucially, CERN gave the web away **for free**, with no patent — a big reason it spread across the whole planet.

## Browsers, CSS, and JavaScript (1990s)

The web exploded once browsers became visual and friendly. **Mosaic** (1993) and then **Netscape Navigator** showed **images alongside text** and were easy enough for ordinary people. Two more inventions made the web what it is now:

- **CSS** (1996) split *style* away from *content*, so designers could control colors, fonts, and layout without messing up the HTML.
- **JavaScript** (1995) added *behaviour* — pages could now respond to clicks, validate forms, and update without reloading.

:::quiz
Q: What's the difference between the internet and the web?
- They're the same thing
- The internet is the global network; the web is a way of linking documents on it *
- The web came first, then the internet
E: The internet is the underlying network (1970s–80s). The web is one service built on top of it (1989+) — documents linked by hypertext.
:::

## Mobile and today

The **iPhone** (2007) put the web in everyone's pocket and forced **responsive design** — layouts that adapt to any screen size. Today the web runs apps as powerful as desktop software, all inside a browser, on any device — yet still built on the same three pillars from 1990.

## How the web works today

Let's go deeper, because this is what you'll actually be building on.

### Client vs server

Every web interaction has two sides:

- The **client** is the program asking for something — usually your **browser** (Chrome, Safari, Firefox).
- The **server** is a computer somewhere that *stores* the website and *sends* it when asked.

:::analogy
Think of a restaurant. You (the **client**) order from a menu. The kitchen (the **server**) prepares the dish and sends it out. You don't see the kitchen — you just get the result.
:::

### What a browser does

A browser is a surprisingly sophisticated piece of software. When it receives a page it:

1. **Reads the HTML** and builds a tree of elements (the "DOM").
2. **Applies the CSS** to decide how everything looks.
3. **Runs the JavaScript** to make the page interactive.
4. **Paints** the final pixels on your screen.

### The request/response cycle

Here's what happens when you visit a page, step by step:

```bash
1. You type a URL and press Enter.
2. The browser finds the server's address (DNS lookup).
3. The browser sends an HTTP REQUEST to that server.
4. The server sends back an HTTP RESPONSE (the HTML).
5. The browser requests extra files it needs (CSS, JS, images).
6. The browser renders the finished page.
```

This back-and-forth is the **request/response cycle**, and it happens for every page and every file.

### Parts of a URL

A **URL** (Uniform Resource Locator) is the address of a resource on the web. Break this one apart:

```bash
https://shop.example.com:443/products/shoes?color=red#reviews
```

- `https` — the **protocol** (how to talk to the server)
- `shop` — the **subdomain**
- `example.com` — the **domain name**
- `:443` — the **port** (usually hidden)
- `/products/shoes` — the **path** (which resource)
- `?color=red` — the **query string** (extra parameters)
- `#reviews` — the **fragment** (jump to a section on the page)

### DNS: turning names into numbers

Computers find each other using numeric **IP addresses** like `93.184.216.34`. But humans remember *names*. **DNS** (Domain Name System) is the system that translates a domain name into an IP address.

:::analogy
DNS is the web's phone book. You know your friend's *name*, but the phone network needs their *number*. DNS looks up the name and hands back the number so the call can connect.
:::

### HTTP and HTTPS

**HTTP** (HyperText Transfer Protocol) is the language clients and servers use to make requests and send responses. **HTTPS** is the same thing but **encrypted**, so nobody between you and the server can read or tamper with the data. The little padlock in your browser bar means you're on HTTPS.

:::warning
Never type passwords or card numbers into a site that shows `http://` without the "s". Without HTTPS, that data travels in plain text that others on the network can read.
:::

### The three frontend languages

Everything you see in a browser is built from three languages, each with one job:

- **HTML — structure.** The content and its meaning: headings, paragraphs, lists, buttons.
- **CSS — style.** How it looks: colors, fonts, spacing, layout.
- **JavaScript — behaviour.** What it does: respond to clicks, fetch data, update the page.

:::analogy
A web page is like a house. **HTML** is the framing and walls (structure). **CSS** is the paint, furniture, and decor (style). **JavaScript** is the electricity and plumbing that make things *work* (behaviour).
:::

:::quiz
Q: Which language is responsible for how a page *looks* (colors, fonts, spacing)?
- HTML
- CSS *
- JavaScript
E: CSS controls presentation. HTML provides structure and meaning; JavaScript adds behaviour and interactivity.
:::

## What a frontend developer actually does

Day to day, a frontend developer:

- Turns designs into real, working pages using HTML and CSS.
- Makes pages **responsive** so they work on phones, tablets, and desktops.
- Adds interactivity with JavaScript — menus, forms, sliders, live updates.
- Connects the page to a server to fetch and display data.
- Tests across browsers and fixes **bugs**.
- Cares about **accessibility** (usable by everyone) and **performance** (fast loading).

:::quiz
Q: In the restaurant analogy, what does the *server* represent?
- Your browser
- The computer that stores the website and sends it when asked *
- The internet cables
E: The server is the kitchen — it stores the site and sends it back in response to the client's (browser's) request.
:::

## Recap

- The **internet** is the global network; the **web** is documents linked by hypertext built on top of it.
- ARPANET (1969) and TCP/IP created the internet; Tim Berners-Lee invented the web at CERN (1989–91).
- HTML, HTTP, and the browser were the original pillars; CSS and JavaScript made the web visual and interactive.
- A **client** (browser) sends an **HTTP request** to a **server**, which sends back a **response** — the request/response cycle.
- A **URL** has a protocol, domain, path, query, and fragment; **DNS** turns the domain name into an IP address.
- The three frontend languages: **HTML** = structure, **CSS** = style, **JavaScript** = behaviour.

**Next up:** Intro to Programming — the core ideas behind *every* coding language, and the foundation for everything you'll build.
