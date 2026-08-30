# How the Internet Works

You use the internet every day, but what actually happens when you open a website? Understanding the plumbing beneath the web — how machines find each other, talk to each other, and send data back and forth — makes you a fundamentally better developer. Bugs that seem mysterious suddenly make sense when you know what's happening under the hood.

## The client-server model

Every interaction on the web follows a simple pattern: one computer **asks** and another computer **answers**.

- The **client** is the program making the request — usually your web browser (Chrome, Firefox, Safari).
- The **server** is a computer somewhere on the internet that *stores* the website's files and *sends* them when asked.

```
Client (your browser)          Server (stores the site)
       |                               |
       |---- "Give me index.html" ---->|
       |                               |
       |<--- here's the HTML file -----|
       |                               |
       |---- "Give me style.css" ----->|
       |                               |
       |<--- here's the CSS file ------|
```

This back-and-forth is called the **request-response cycle**, and it happens for every resource — HTML files, CSS files, JavaScript files, images, fonts.

:::key
The web runs on a simple conversation: the client *requests*, the server *responds*. Every page load, every API call, every image — it's always this same pattern. Client asks, server answers.
:::

## IP addresses: finding computers on the network

Every device on the internet has a unique address called an **IP address** — like a street address for computers:

```
IPv4:  93.184.216.34         (four numbers, 0–255 each)
IPv6:  2606:2800:0220:0001:0248:1893:25c8:1946   (longer, for more addresses)
```

When your browser needs to reach a server, it ultimately needs the server's IP address. But you don't type IP addresses — you type domain names like `google.com`. That's where DNS comes in.

## DNS: the internet's phone book

**DNS** (Domain Name System) translates human-readable domain names into IP addresses.

When you type `example.com`, here's what happens:

1. Your browser checks its local **cache** — has it looked this up recently?
2. If not, it asks your operating system, which asks your **ISP's DNS resolver**.
3. The resolver queries the DNS hierarchy: **root servers** → **.com nameservers** → **example.com's nameserver**.
4. The answer comes back: `example.com` = `93.184.216.34`.
5. Your browser can now connect to that IP address.

This takes milliseconds and is cached so repeat visits are instant.

```
You type: example.com
         ↓
DNS lookup: "What IP is example.com?"
         ↓
Answer: 93.184.216.34
         ↓
Browser connects to 93.184.216.34
```

:::tip
You can see DNS in action. Open a terminal and run `nslookup google.com` — it shows you the IP address your computer resolves that domain to. Every domain you visit goes through this translation step.
:::

## URLs: the anatomy of a web address

A **URL** (Uniform Resource Locator) is the full address of a resource on the web. Let's break one apart:

```
https://shop.example.com:443/products/shoes?color=red&size=10#reviews
```

| Part            | Value               | Purpose                          |
|-----------------|---------------------|----------------------------------|
| **Protocol**    | `https`             | How to communicate (HTTP/HTTPS)  |
| **Subdomain**   | `shop`              | A subdivision of the domain      |
| **Domain**      | `example.com`       | The server's name                |
| **Port**        | `:443`              | The network door (usually hidden)|
| **Path**        | `/products/shoes`   | Which resource on the server     |
| **Query string**| `?color=red&size=10`| Extra parameters (key=value)     |
| **Fragment**    | `#reviews`          | Jump to a section on the page    |

The protocol, domain, and path are the essentials. Query strings pass data to the server. Fragments stay in the browser — the server never sees them.

## The HTTP request-response cycle

HTTP (HyperText Transfer Protocol) is the language clients and servers use. Every HTTP interaction has two parts:

### The request

When your browser asks for a page, it sends an HTTP request:

```
GET /products/shoes HTTP/1.1
Host: shop.example.com
Accept: text/html
User-Agent: Chrome/120
```

- **Method** — `GET` (retrieve data), `POST` (send data), `PUT`, `DELETE`
- **Path** — which resource
- **Headers** — metadata about the request (browser type, accepted formats, cookies)

### The response

The server answers with an HTTP response:

```
HTTP/1.1 200 OK
Content-Type: text/html
Content-Length: 5120

<!DOCTYPE html>
<html>...the page content...</html>
```

- **Status code** — a number indicating what happened
- **Headers** — metadata about the response (content type, caching rules)
- **Body** — the actual content (HTML, JSON, an image, etc.)

## HTTP status codes

Status codes tell you what happened with the request. You don't need to memorize all of them, but these are the ones you'll see constantly:

| Code | Meaning            | When you see it                        |
|------|--------------------|----------------------------------------|
| 200  | OK                 | Everything worked                      |
| 201  | Created            | A new resource was created (POST)      |
| 301  | Moved Permanently  | The URL changed — follow the redirect  |
| 304  | Not Modified       | Use your cached version                |
| 400  | Bad Request        | The server can't understand your request|
| 401  | Unauthorized       | You need to log in                     |
| 403  | Forbidden          | Logged in, but not allowed             |
| 404  | Not Found          | The resource doesn't exist             |
| 500  | Internal Server Error | Something broke on the server       |
| 503  | Service Unavailable | Server is overloaded or down          |

:::key
Status codes follow a pattern: **2xx** = success, **3xx** = redirect, **4xx** = client error (you did something wrong), **5xx** = server error (the server broke). Knowing this pattern means you can diagnose problems even with codes you've never seen.
:::

## HTTPS and TLS: security basics

**HTTPS** is HTTP with encryption. The "S" stands for Secure. It uses a protocol called **TLS** (Transport Layer Security) to:

1. **Encrypt** the data so no one between you and the server can read it.
2. **Authenticate** the server so you know you're talking to the real `bank.com`, not an impostor.
3. **Protect integrity** so the data can't be tampered with in transit.

When your browser shows a **padlock icon** in the address bar, you're on HTTPS.

```
HTTP:  Data travels as plain text — anyone on the network can read it
HTTPS: Data is encrypted — only your browser and the server can read it
```

:::warning
Any site that handles passwords, payment information, or personal data *must* use HTTPS. Modern browsers flag HTTP sites as "Not Secure." As a developer, always use HTTPS — free certificates are available from services like Let's Encrypt.
:::

## What happens when you type a URL

Putting it all together — the full journey when you type `https://www.example.com` and press Enter:

```
1. PARSE THE URL
   → Protocol: HTTPS, Domain: www.example.com, Path: /

2. DNS LOOKUP
   → Browser asks DNS: "What's the IP for www.example.com?"
   → DNS responds: 93.184.216.34

3. TCP CONNECTION
   → Browser opens a connection to 93.184.216.34 on port 443

4. TLS HANDSHAKE
   → Browser and server exchange encryption keys
   → A secure tunnel is established

5. HTTP REQUEST
   → Browser sends: GET / HTTP/1.1

6. SERVER PROCESSING
   → Server finds the resource and prepares the response

7. HTTP RESPONSE
   → Server sends back: 200 OK + HTML content

8. BROWSER RENDERING
   → Browser parses HTML, discovers CSS and JS files
   → Sends additional requests for those files
   → Builds the page and paints it on screen
```

Each of these steps happens in milliseconds. A typical page load involves 20–100 individual requests for HTML, CSS, JS, images, and fonts — all following this same cycle.

:::tip
Open your browser's DevTools (F12) → **Network tab** → reload any page. You'll see every single request, its status code, size, and timing. This is one of the most important debugging tools you'll ever use.
:::

:::quiz
Q: What does DNS do?
- Encrypts data between client and server
- Translates domain names (like google.com) into IP addresses *
- Compresses HTML files for faster delivery
E: DNS (Domain Name System) is the internet's phone book — it converts human-readable domain names into the numeric IP addresses that computers use to find each other.
:::

:::quiz
Q: A server responds with status code 404. What does that mean?
- The server crashed
- The requested resource was not found *
- The request was successful
E: 404 means "Not Found" — the server understood the request but the resource at that URL doesn't exist. It's a 4xx code, meaning the error is on the client side (wrong URL).
:::

:::quiz
Q: What is the main benefit of HTTPS over HTTP?
- HTTPS is faster
- HTTPS encrypts data so it can't be read or tampered with in transit *
- HTTPS allows larger file transfers
E: HTTPS uses TLS encryption to protect data in transit. Without it, data like passwords and credit card numbers travel as plain text that anyone on the network can intercept.
:::

## Recap

- The web uses a **client-server model**: the browser (client) requests, the server responds.
- **IP addresses** identify computers; **DNS** translates domain names into IP addresses.
- A **URL** has a protocol, domain, path, query string, and fragment.
- **HTTP** is the request-response protocol. Requests have a method and headers; responses have a status code and body.
- Status codes: **2xx** = success, **3xx** = redirect, **4xx** = client error, **5xx** = server error.
- **HTTPS** encrypts communication with TLS — always use it.
- When you type a URL: DNS lookup → TCP connection → TLS handshake → HTTP request → server processing → response → browser rendering.

**Next up:** How Browsers Render Pages — what happens inside the browser after it receives the HTML.
