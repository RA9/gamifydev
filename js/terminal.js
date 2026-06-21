// Terminal Trainer (Phase 5) — a safe, simulated Bash shell.
// A fully in-memory virtual filesystem + command interpreter (no real shell),
// wrapped in a mission sequence so learners practise real commands by doing.

// --- The simulated shell ----------------------------------------------------
function createShell() {
  const root = { type: "dir", children: {} };
  const ensureDir = (node, name) => {
    if (!node.children[name]) node.children[name] = { type: "dir", children: {} };
    return node.children[name];
  };

  // Starting filesystem: /home/dev with a readme and a notes folder.
  const home = ensureDir(ensureDir(root, "home"), "dev");
  home.children["readme.txt"] = {
    type: "file",
    content: "Welcome to the terminal!\nTry: ls, cat, cd, mkdir, touch, echo.\n",
  };
  const notes = ensureDir(home, "notes");
  notes.children["todo.txt"] = { type: "file", content: "- learn bash\n- build cool stuff\n" };

  let cwd = ["home", "dev"];
  let lastCmd = "";

  const nodeAt = (parts) => {
    let n = root;
    for (const p of parts) {
      if (n.type !== "dir" || !n.children[p]) return null;
      n = n.children[p];
    }
    return n;
  };
  const cwdNode = () => nodeAt(cwd);

  const resolve = (path) => {
    let parts;
    if (!path || path === "~") return ["home", "dev"];
    if (path.startsWith("/")) parts = path.split("/").filter(Boolean);
    else if (path.startsWith("~/")) parts = ["home", "dev", ...path.slice(2).split("/").filter(Boolean)];
    else parts = [...cwd, ...path.split("/").filter(Boolean)];
    const out = [];
    for (const p of parts) {
      if (p === ".") continue;
      else if (p === "..") out.pop();
      else out.push(p);
    }
    return out;
  };

  const pathStr = (parts) => {
    if (parts[0] === "home" && parts[1] === "dev") {
      const rest = parts.slice(2);
      return "~" + (rest.length ? "/" + rest.join("/") : "");
    }
    return "/" + parts.join("/");
  };

  const tokenize = (line) => {
    const toks = [];
    let cur = "";
    let q = false;
    for (const ch of line) {
      if (ch === '"' || ch === "'") {
        q = !q;
        continue;
      }
      if (ch === " " && !q) {
        if (cur) {
          toks.push(cur);
          cur = "";
        }
        continue;
      }
      cur += ch;
    }
    if (cur) toks.push(cur);
    return toks;
  };

  function run(line) {
    line = (line || "").trim();
    lastCmd = line;
    if (!line) return { output: "" };

    let redirect = null;
    let body = line;
    const gt = line.indexOf(">");
    if (gt !== -1) {
      redirect = line.slice(gt + 1).trim().replace(/^["']|["']$/g, "");
      body = line.slice(0, gt).trim();
    }
    const toks = tokenize(body);
    const cmd = toks[0];
    const args = toks.slice(1);
    const here = cwdNode();

    switch (cmd) {
      case "help":
        return { output: "Commands: ls [-a], cd, pwd, cat, mkdir, touch, echo, rm [-r], clear, whoami, help" };
      case "clear":
        return { cleared: true, output: "" };
      case "whoami":
        return { output: "dev" };
      case "pwd":
        return { output: pathStr(cwd) };
      case "ls": {
        const target = args.find((a) => !a.startsWith("-"));
        const showHidden = args.includes("-a");
        const node = target ? nodeAt(resolve(target)) : here;
        if (!node) return { output: `ls: ${target}: No such file or directory` };
        if (node.type === "file") return { output: target };
        const names = Object.keys(node.children)
          .filter((n) => showHidden || !n.startsWith("."))
          .sort();
        return { output: names.map((n) => (node.children[n].type === "dir" ? n + "/" : n)).join("  ") };
      }
      case "cd": {
        const target = args[0] || "~";
        const parts = resolve(target);
        const node = nodeAt(parts);
        if (!node) return { output: `cd: ${target}: No such file or directory` };
        if (node.type !== "dir") return { output: `cd: ${target}: Not a directory` };
        cwd = parts;
        return { output: "" };
      }
      case "mkdir": {
        if (!args[0]) return { output: "mkdir: missing operand" };
        if (here.children[args[0]]) return { output: `mkdir: cannot create directory '${args[0]}': File exists` };
        here.children[args[0]] = { type: "dir", children: {} };
        return { output: "" };
      }
      case "touch": {
        if (!args[0]) return { output: "touch: missing file operand" };
        if (!here.children[args[0]]) here.children[args[0]] = { type: "file", content: "" };
        return { output: "" };
      }
      case "cat": {
        if (!args[0]) return { output: "cat: missing operand" };
        const node = nodeAt(resolve(args[0]));
        if (!node) return { output: `cat: ${args[0]}: No such file or directory` };
        if (node.type === "dir") return { output: `cat: ${args[0]}: Is a directory` };
        return { output: node.content.replace(/\n$/, "") };
      }
      case "echo": {
        const text = args.join(" ");
        if (redirect) {
          here.children[redirect] = { type: "file", content: text + "\n" };
          return { output: "" };
        }
        return { output: text };
      }
      case "rm": {
        const recursive = args.some((a) => /^-[rf]+$/.test(a));
        const target = args.find((a) => !a.startsWith("-"));
        if (!target) return { output: "rm: missing operand" };
        if (!here.children[target]) return { output: `rm: cannot remove '${target}': No such file or directory` };
        if (here.children[target].type === "dir" && !recursive)
          return { output: `rm: cannot remove '${target}': Is a directory` };
        delete here.children[target];
        return { output: "" };
      }
      default:
        return { output: `${cmd}: command not found (try 'help')` };
    }
  }

  return {
    run,
    prompt: () => `dev@gamifydev:${pathStr(cwd)}$`,
    cwdName: () => cwd[cwd.length - 1],
    get lastCmd() {
      return lastCmd;
    },
    childType: (name) => {
      const h = cwdNode();
      return h && h.children[name] ? h.children[name].type : null;
    },
    fileContent: (name) => {
      const h = cwdNode();
      const c = h && h.children[name];
      return c && c.type === "file" ? c.content : null;
    },
  };
}

// --- Missions ---------------------------------------------------------------
const TERMINAL_CHALLENGES = [
  { goal: "Print the folder you're currently in.", check: { type: "ran", cmd: "pwd" }, hint: "Type `pwd` and press Enter." },
  { goal: "List what's in this folder.", check: { type: "ran", cmd: "ls" }, hint: "Type `ls`." },
  { goal: "Read the contents of readme.txt.", check: { type: "ran", cmd: "cat readme.txt" }, hint: "Use `cat readme.txt`." },
  { goal: "Create a new folder called projects.", check: { type: "dir", name: "projects" }, hint: "`mkdir projects`" },
  { goal: "Move into the projects folder.", check: { type: "cwdName", name: "projects" }, hint: "`cd projects`" },
  { goal: "Create an empty file called app.js.", check: { type: "file", name: "app.js" }, hint: "`touch app.js`" },
  { goal: "Go back up to the parent folder.", check: { type: "cwdName", name: "dev" }, hint: "`cd ..`" },
  { goal: "Make a file notes.md that contains the word hello.", check: { type: "fileContains", name: "notes.md", text: "hello" }, hint: '`echo hello > notes.md`' },
  { goal: "Delete the notes.md file you just made.", check: { type: "noFile", name: "notes.md" }, hint: "`rm notes.md`" },
];

function termNormalize(s) {
  return (s || "").trim().replace(/\s+/g, " ");
}

function challengePassed(check, shell) {
  switch (check.type) {
    case "ran":
      return termNormalize(shell.lastCmd) === termNormalize(check.cmd);
    case "dir":
      return shell.childType(check.name) === "dir";
    case "file":
      return shell.childType(check.name) === "file";
    case "noFile":
      return shell.childType(check.name) === null;
    case "cwdName":
      return shell.cwdName() === check.name;
    case "fileContains": {
      const c = shell.fileContent(check.name);
      return c != null && c.includes(check.text);
    }
    default:
      return false;
  }
}

// --- The Terminal Trainer page ---------------------------------------------
async function TerminalPage(htmlEl) {
  const challenges = TERMINAL_CHALLENGES;
  let current = 0;
  if (typeof getMeta === "function") {
    const saved = await getMeta("terminal-progress", 0);
    current = Math.min(Math.max(0, saved | 0), challenges.length);
  }
  let shell = createShell();

  const completionCard = () => `
    <div class="gd-card text-center">
      <div class="text-4xl mb-2">🏆</div>
      <h2 class="text-2xl font-extrabold mb-1">Bash Bootcamp complete!</h2>
      <p class="text-slate-600 mb-5">You ran real commands and cleared every mission. You can find your way around a terminal now — a genuine dev superpower.</p>
      <div class="flex justify-center gap-3">
        <button class="term-restart gd-btn gd-btn-secondary">Start over</button>
        <a href="#today" class="gd-btn gd-btn-primary">Back to Today</a>
      </div>
    </div>`;

  function render() {
    const allDone = current >= challenges.length;
    const pct = Math.round((Math.min(current, challenges.length) / challenges.length) * 100);
    const ch = challenges[current];

    htmlEl.innerHTML = `
      <div class="max-w-3xl mx-auto animate-fade-up space-y-4">
        <div class="gd-card">
          <span class="gd-chip gd-chip-brand mb-2">⌨️ Terminal Trainer</span>
          <h1 class="text-2xl font-extrabold">Bash Bootcamp</h1>
          <p class="text-slate-500 text-sm">Type real commands to clear each mission. Stuck? Type <code class="gd-code">help</code>.</p>
          <div class="gd-progress h-2.5 mt-3"><div class="gd-progress-fill" style="width: ${pct}%"></div></div>
        </div>
        ${
          allDone
            ? completionCard()
            : `
        <div class="gd-card">
          <div class="flex items-center justify-between mb-2 gap-3">
            <span class="gd-chip gd-chip-brand">Mission ${current + 1} / ${challenges.length}</span>
            <button class="term-hint text-sm font-bold text-brand-600">Show hint</button>
          </div>
          <p class="font-extrabold mb-1">🎯 ${ch.goal}</p>
          <p class="term-hint-text hidden text-sm text-slate-500 mb-1">${ch.hint}</p>
          <div class="gd-terminal mt-3">
            <div class="gd-terminal-bar"><span></span><span></span><span></span></div>
            <div class="gd-terminal-body">
              <div class="term-output"></div>
              <div class="term-line"><span class="term-prompt">${shell.prompt()}</span><input class="term-input" autocomplete="off" autocapitalize="off" autocorrect="off" spellcheck="false" aria-label="Terminal input" /></div>
            </div>
          </div>
          <div class="term-success hidden mt-3"></div>
        </div>`
        }
      </div>`;

    if (allDone) {
      const restart = htmlEl.querySelector(".term-restart");
      if (restart)
        restart.addEventListener("click", async () => {
          current = 0;
          shell = createShell();
          if (typeof setMeta === "function") await setMeta("terminal-progress", 0);
          render();
        });
      return;
    }
    wire();
  }

  function wire() {
    const input = htmlEl.querySelector(".term-input");
    const output = htmlEl.querySelector(".term-output");
    const promptEl = htmlEl.querySelector(".term-prompt");
    const successEl = htmlEl.querySelector(".term-success");
    const body = htmlEl.querySelector(".gd-terminal-body");
    const hintBtn = htmlEl.querySelector(".term-hint");
    const hintText = htmlEl.querySelector(".term-hint-text");
    if (hintBtn) hintBtn.addEventListener("click", () => hintText.classList.toggle("hidden"));
    if (body) body.addEventListener("click", () => input && input.focus());
    if (!input) return;
    input.focus();

    const append = (text, cls) => {
      const div = document.createElement("div");
      if (cls) div.className = cls;
      div.textContent = text;
      output.appendChild(div);
    };
    let solved = false;

    input.addEventListener("keydown", (e) => {
      if (e.key !== "Enter") return;
      const line = input.value;
      append(`${shell.prompt()} ${line}`, "text-slate-400");
      const res = shell.run(line);
      if (res.cleared) output.innerHTML = "";
      else if (res.output) append(res.output);
      input.value = "";
      promptEl.textContent = shell.prompt();
      if (body) body.scrollTop = body.scrollHeight;

      if (!solved && challengePassed(challenges[current].check, shell)) {
        solved = true;
        input.disabled = true;
        const last = current + 1 >= challenges.length;
        successEl.classList.remove("hidden");
        successEl.innerHTML = `<div class="flex items-center justify-between gap-3 rounded-xl bg-grass-50 border border-grass-200 text-grass-800 px-3 py-2 text-sm font-bold">
          <span>✅ Mission complete!</span>
          <button class="term-next gd-btn gd-btn-primary !py-1.5 !px-4 !text-xs">${last ? "Finish 🏆" : "Next mission"}</button>
        </div>`;
        if (typeof setMeta === "function") setMeta("terminal-progress", current + 1);
        successEl.querySelector(".term-next").addEventListener("click", () => {
          current++;
          render();
        });
      }
    });
  }

  render();
}
