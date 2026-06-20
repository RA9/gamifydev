// import { DB, createStorage, getStorage, updateStorage } from "./storage.js";
// import questions from "../data/questions.json" assert { type: "json" };
// import { randomID, htmlencode } from "./utils.js";

async function TestPage(htmlEl) {
  const stateArr = await DB.states.where("name").equals("tys").toArray();
  const state =
    stateArr.length > 0
      ? stateArr[stateArr.length - 1]
      : await createStorage("states", {
          id: randomID(),
          name: "tys",
          previous: null,
          current: "tys",
          next: "tys-quiz",
        });

  if (state.current === "tys") {
    const languages = [
      { value: "html", label: "HTML" },
      { value: "css", label: "CSS" },
      { value: "javascript", label: "JavaScript" },
      { value: "c", label: "C" },
      { value: "python", label: "Python" },
      { value: "java", label: "Java" },
      { value: "sql", label: "SQL" },
    ];

    htmlEl.innerHTML = `
      <div class="max-w-3xl mx-auto animate-fade-up">
        <div class="text-center mb-6">
          <span class="gd-chip gd-chip-brand mb-3">${icon("zap", "w-3.5 h-3.5")} Test Yourself</span>
          <h1 class="text-2xl sm:text-3xl font-extrabold">Choose your challenge</h1>
          <p class="text-slate-500 mt-1">Pick a language and how many questions you want.</p>
        </div>
        <div class="gd-card space-y-6">
          <div>
            <span class="gd-label">Language</span>
            <input type="hidden" id="language" value="html" />
            <div id="lang-grid" class="grid grid-cols-2 sm:grid-cols-3 gap-3">
              ${languages
                .map(
                  (l, i) => `
                <button type="button" data-lang="${l.value}"
                  class="lang-tile gd-option justify-center flex-col gap-2 py-4 ${
                    i === 0 ? "gd-option-selected" : ""
                  }">
                  ${langBadge(l.value, "h-11 w-11 text-sm")}
                  <span class="text-sm font-extrabold">${l.label}</span>
                </button>`
                )
                .join("")}
            </div>
          </div>

          <div>
            <label for="numQuestions" class="gd-label">Number of questions</label>
            <div class="flex items-center gap-3">
              <input id="numQuestions" name="numQuestions" type="number" min="1" max="50" value="20" class="gd-input w-28 text-center" />
              <div class="flex gap-2">
                ${[10, 20, 30]
                  .map(
                    (n) =>
                      `<button type="button" class="num-preset gd-btn gd-btn-secondary !py-2 !px-4 !text-sm" data-num="${n}">${n}</button>`
                  )
                  .join("")}
              </div>
            </div>
          </div>

          <button id="start-tys" class="gd-btn gd-btn-primary gd-btn-block">Start Quiz</button>
        </div>
      </div>
      `;

    // Language tile selection -> hidden input.
    htmlEl.querySelectorAll(".lang-tile").forEach((tile) => {
      tile.addEventListener("click", () => {
        htmlEl
          .querySelectorAll(".lang-tile")
          .forEach((t) => t.classList.remove("gd-option-selected"));
        tile.classList.add("gd-option-selected");
        document.querySelector("#language").value = tile.dataset.lang;
      });
    });

    // Quick presets for question count.
    htmlEl.querySelectorAll(".num-preset").forEach((btn) => {
      btn.addEventListener("click", () => {
        document.querySelector("#numQuestions").value = btn.dataset.num;
      });
    });

    const START_TYS_BUTTON = document.querySelector("#start-tys");

    START_TYS_BUTTON.addEventListener("click", async () => {
      const language = document.querySelector("#language").value;
      const numQuestions = document.querySelector("#numQuestions").value;
      state.current = "tys-quiz";
      state.previous = "tys";
      state.next = null;

      await createStorage("tests", {
        id: randomID(),
        name: "tys-" + randomID(),
        is_completed: false,
        language: language,
        numQuestions,
        created_at: new Date(),
      });

      await updateStorage("states", state);

      // window.location.reload()
      const page = document.querySelector("main");
      TestYourselfSection(page);
    });
  } else if (state.current === "tys-quiz") {
    const page = document.querySelector("main");
    TestYourselfSection(page);
  } else if (state.current === "tys-quiz-result") {
    TestResultPage(document.querySelector("main"));
  }
}

function tysRandomizeOptions(options, answer = null) {
  const randomOptions = shuffle(options);
  return randomOptions.map((option) => {
    const id = randomID();
    return `
    <label for="${id}" class="gd-option mb-2">
      <input type="radio" id="${id}" name="option" value="${escapeHTMLToEntities(
      option
    )}" class="h-5 w-5 shrink-0 accent-brand-500">
      <span>${escapeHTMLToEntities(option)}</span>
    </label>
  `;
  });
}

function getRandomItem(limit) {
  return Math.floor(Math.random() * limit + 1);
}

function randomizedQuestions(questions, limit) {
  const questionsLimit = shuffle(questions).slice(0, limit);
  return questionsLimit;
}

function getQuestionsByLimit(questions, limit) {
  return randomizedQuestions(questions, limit);
}

function tysCountDown(duration) {
  let timer = Number(duration),
    minutes,
    seconds;
  // document.querySelector("#tys-duration").textContent = "Duration: " + timer;
  const interval = setInterval(function () {
    const durationEl = document.querySelector("#tys-duration");
    // Quiz DOM is gone (submitted/cancelled/navigated away) — stop ticking.
    if (!durationEl) {
      clearInterval(interval);
      return;
    }
    minutes = parseInt(timer / 60, 10);
    seconds = parseInt(timer % 60, 10);
    // console.log({ timer, minutes, seconds });
    minutes = minutes < 10 ? "0" + minutes : minutes;
    seconds = seconds < 10 ? "0" + seconds : seconds;

    // console.log({ minutes, seconds });

    if (timer === 0) {
      clearInterval(interval);
      document.querySelector(
        "#tys-duration"
      ).innerHTML = `<span class="flex items-center gap-1.5 text-rose-500">${icon(
        "clock",
        "w-4 h-4"
      )} Time's up!</span>`;
      // Lock in the current answers and auto-submit when time runs out.
      document.querySelectorAll("input[type=radio]").forEach((el) => {
        el.classList.add("cursor-not-allowed");
        el.setAttribute("disabled", true);
      });
      document.querySelector("#tys-submit")?.click();
    } else {
      const low = timer <= 10;
      document.querySelector(
        "#tys-duration"
      ).innerHTML = `<span class="flex items-center gap-1.5 ${
        low ? "text-rose-500" : "text-grass-600"
      }">${icon("clock", "w-4 h-4")} ${minutes}:${seconds}</span>`;
      timer--;
    }
  }, 1000);
}

async function TestYourselfSection(htmlEl) {
  const state = (await DB.states.where("name").equals("tys").toArray())[0];
  const test = await DB.tests
    .where("name")
    .startsWith("tys")
    .and((test) => !test.is_completed)
    .last();

  let testQuestions = (await DB.table("test_questions").toArray()) || [];

  if (testQuestions.length <= 0) {
    const questions = await DB.questions.toArray();

    await createStorage("test_questions", {
      test_id: test.id,
      id: randomID(),
      questions: getQuestionsByLimit(
        questions.filter((question) => question.category === test.language),
        Number(test.numQuestions)
      ),
    });

    testQuestions = (await DB.table("test_questions").toArray())[0];
  }

  testQuestions = Array.isArray(testQuestions)
    ? testQuestions[0]
    : testQuestions;

  // console.log({ testQuestions });

  if (state.current === "tys-quiz") {
    const total = [].concat(testQuestions.questions).length;
    const elHTML = [].concat(testQuestions.questions).map((question, index) => {
      return `
      <div class="question gd-card-sm">
        <div class="flex items-start gap-3 mb-4">
          <span class="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-brand-100 text-brand-700 font-extrabold text-sm">${
            index + 1
          }</span>
          <p class="text-lg font-bold pt-0.5">${escapeHTMLToEntities(
            question.details.question
          )}</p>
        </div>
        <form class="space-y-2">
          ${tysRandomizeOptions(question.details.options).join("")}
        </form>
      </div>
      `;
    });

    const sumDuration = testQuestions.questions.reduce((acc, question) => {
      // console.log({ question });
      const duration =
        typeof question.details.duration == "string"
          ? question.details.duration.split(" ")
          : question.details.duration;
      if (typeof duration !== "number") {
        if (duration.includes("minutes")) {
          acc += Number(duration[0]) * 60;
        } else if (duration.includes("seconds")) {
          acc += Number(duration[0]);
        }
      } else {
        acc += duration;
      }
      return acc;
    }, 0);

    htmlEl.innerHTML = `
      <div class="max-w-3xl mx-auto animate-fade-up">
        <div class="sticky top-16 z-30 -mx-4 px-4 py-3 mb-4 bg-white/80 dark:bg-slate-900/70 backdrop-blur rounded-b-2xl border-b border-slate-200/70">
          <div class="flex items-center justify-between gap-4">
            <span class="gd-chip gd-chip-brand">${test.language.toUpperCase()}</span>
            <span class="text-sm font-bold text-slate-500">${total} question${
      total === 1 ? "" : "s"
    }</span>
            <span id="tys-duration" class="flex items-center gap-1.5 rounded-full bg-slate-100 px-3 py-1.5 font-extrabold text-sm"></span>
          </div>
        </div>

        <div class="space-y-4">
          ${elHTML.join("")}
        </div>

        <div class="flex gap-3 mt-6">
          <button id="tys-cancel" class="gd-btn gd-btn-secondary flex-1">Cancel</button>
          <button id="tys-submit" class="gd-btn gd-btn-primary flex-[2]">Submit Quiz</button>
        </div>
      </div>
    `;

    tysCountDown(sumDuration);

    const SUBMIT_BUTTON = document.querySelector("#tys-submit");
    const CANCEL_BUTTON = document.querySelector("#tys-cancel");

    CANCEL_BUTTON.addEventListener("click", () => {
      state.current = "tys";
      state.previous = "tys-quiz";
      state.next = null;
      updateStorage("states", state);

      // reload page and remove test_questions
      window.location.reload();
      removeStorage("test_questions", testQuestions.id);
      updateStorage("tests", { ...test, is_completed: true });
    });

    // intercept the location.reload() function and cancel the test
    window.addEventListener("beforeunload", (e) => {
      if (state.current === "tys-quiz") {
        e.preventDefault();
        e.returnValue = "";
        state.current = "tys";
        state.previous = "tys-quiz";
        state.next = null;
        updateStorage("states", state);

        // reload page and remove test_questions
        removeStorage("test_questions", testQuestions.id);
        updateStorage("tests", { ...test, is_completed: true });
      }
    });

    SUBMIT_BUTTON.addEventListener("click", async () => {
      state.current = "tys-quiz-result";
      state.previous = "tys-quiz";
      state.next = null;
      // const page = document.querySelector("main");

      // Query all elements with the .question class
      const questionElements = document.querySelectorAll(".question");

      // Create an array to store the selected options for each question
      const selectedOptions = [];

      // Loop through each question element and get the value of the checked input
      questionElements.forEach((questionElement) => {
        const selectedOption =
          questionElement.querySelector('input[name="option"]:checked')
            ?.value || null;
        selectedOptions.push(selectedOption);
      });

      const scoreDetails = {
        id: randomID(),
        test_id: test.id,
        score: 0,
        numWrong: 0,
        numCorrect: 0,
        details: {},
      };

      // Loop through each test question and compare it with the selected options
      const quizDetails = calculateQuizScore(testQuestions, selectedOptions);
      scoreDetails.score = quizDetails.score;
      scoreDetails.numCorrect = quizDetails.numCorrect;
      scoreDetails.numWrong = quizDetails.numWrong;

      // add details
      scoreDetails.details = {
        questions: testQuestions.questions,
        selectedOptions,
      };

      scoreDetails.created_at = new Date();

      await updateStorage("states", state);
      await updateStorage("tests", { ...test, is_completed: true });
      await createStorage("scores", scoreDetails);

      // QuizPage(page);
      TestResultPage(document.querySelector("main"));
    });
  }
}

function calculateQuizScore(testQuestions, selectedOptions) {
  let score = 0;
  let numCorrect = 0;
  let numWrong = 0;

  testQuestions.questions.forEach((question, index) => {
    const selectedOption = selectedOptions[index];
    const correctAnswer = question.details.answer;

    if (selectedOption === correctAnswer) {
      numCorrect++;
    } else {
      numWrong++;
    }
  });

  score = Math.round((numCorrect / testQuestions.questions.length) * 100);

  return { score, numCorrect, numWrong };
}

// Builds the per-question review: highlights the correct option, flags the
// user's wrong pick, and shows the explanation for each question.
function buildTysReview(testDetails) {
  const { questions, selectedOptions } = testDetails.details;

  return questions
    .map((q, i) => {
      const correct = q.details.answer;
      const chosen = selectedOptions[i];
      const isRight = chosen === correct;

      const optionsHTML = q.details.options
        .map((opt) => {
          const isCorrectOpt = opt === correct;
          const isChosenWrong = opt === chosen && !isCorrectOpt;
          let cls = "border-slate-200 text-slate-600";
          let tag = "";
          if (isCorrectOpt) {
            cls = "border-green-500 bg-green-50 text-green-800";
            tag = `<span class="inline-flex items-center gap-1 font-extrabold shrink-0">${icon(
              "checkCircle",
              "w-4 h-4"
            )} correct answer</span>`;
          } else if (isChosenWrong) {
            cls = "border-red-500 bg-red-50 text-red-800";
            tag = `<span class="inline-flex items-center gap-1 font-extrabold shrink-0">${icon(
              "xCircle",
              "w-4 h-4"
            )} your answer</span>`;
          }
          return `<li class="flex items-center justify-between gap-2 border-2 ${cls} rounded-xl px-3 py-2 mb-1.5 list-none font-semibold"><span>${escapeHTMLToEntities(
            opt
          )}</span>${tag}</li>`;
        })
        .join("");

      return `
      <div class="border-b border-gray-200 py-4">
        <p class="font-bold mb-2">${i + 1}. ${escapeHTMLToEntities(
        q.details.question
      )}
          <span class="text-sm font-normal ${
            isRight ? "text-green-600" : "text-red-600"
          }">(${isRight ? "Correct" : "Incorrect"})</span>
        </p>
        <ul class="mb-2">${optionsHTML}</ul>
        ${
          chosen === null
            ? '<p class="text-sm text-gray-500 mb-1">You did not answer this question.</p>'
            : ""
        }
        ${
          q.details.explanation
            ? `<p class="text-sm text-gray-600"><span class="font-semibold">Explanation:</span> ${escapeHTMLToEntities(
                q.details.explanation
              )}</p>`
            : ""
        }
      </div>`;
    })
    .join("");
}

async function TestResultPage(htmlEl) {
  const state = (await DB.states.where("name").equals("tys").toArray())[0];
  const test = await DB.tests
    .where("name")
    .startsWith("tys")
    .and((test) => test.is_completed)
    .last();

  const testDetails = await DB.scores.where("test_id").equals(test.id).last();
  console.log({ testDetails }, "What is this");

  if (state.current === "tys-quiz-result") {
    const score = testDetails.score;
    const tone = score >= 85 ? "grass" : score >= 70 ? "amber" : "rose";
    const ringColor =
      tone === "grass" ? "#46a302" : tone === "amber" ? "#f59e0b" : "#f43f5e";
    const headline =
      score >= 85 ? "Outstanding!" : score >= 70 ? "Nice work!" : "Keep practicing!";
    const resultIcon = score >= 85 ? "trophy" : score >= 70 ? "medal" : "target";
    const xpEarned = testDetails.numCorrect * 10;

    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <div class="grid h-16 w-16 mx-auto mb-4 place-items-center rounded-2xl bg-${tone}-50 text-${tone}-500 animate-pop-in">${icon(
      resultIcon,
      "w-8 h-8"
    )}</div>
          <p class="gd-chip gd-chip-brand mx-auto mb-3">${test.language.toUpperCase()} Results</p>
          <h1 class="text-2xl sm:text-3xl font-extrabold mb-6">${headline}</h1>

          <div class="relative mx-auto h-40 w-40 mb-6">
            <div class="absolute inset-0 rounded-full"
              style="background: conic-gradient(${ringColor} ${score}%, rgba(148,163,184,0.2) 0);"></div>
            <div class="absolute inset-[10px] rounded-full bg-white dark:bg-[#161c30] grid place-items-center">
              <div>
                <p class="text-4xl font-extrabold">${score}<span class="text-xl">%</span></p>
                <p class="text-xs font-bold uppercase tracking-wide text-slate-400">Score</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-3 gap-3 mb-6">
            <div class="rounded-2xl bg-grass-50 border border-grass-200 py-3">
              <p class="text-2xl font-extrabold text-grass-600">${testDetails.numCorrect}</p>
              <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Correct</p>
            </div>
            <div class="rounded-2xl bg-rose-50 border border-rose-200 py-3">
              <p class="text-2xl font-extrabold text-rose-500">${testDetails.numWrong}</p>
              <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Wrong</p>
            </div>
            <div class="rounded-2xl bg-amber-50 border border-amber-200 py-3">
              <p class="text-2xl font-extrabold text-amber-500">+${xpEarned}</p>
              <p class="flex items-center justify-center gap-1 text-xs font-bold uppercase tracking-wide text-slate-500">${icon(
                "gem",
                "w-3.5 h-3.5"
              )} XP</p>
            </div>
          </div>

          <button id="tys-review" class="gd-btn gd-btn-secondary gd-btn-block mb-3">Review Answers</button>
          <div id="tys-review-container" class="hidden mb-4 text-left"></div>
          <div class="flex gap-3">
            <button id="home" class="gd-btn gd-btn-secondary flex-1">Home</button>
            <button id="tys-another" class="gd-btn gd-btn-primary flex-1">Take Another</button>
          </div>
        </div>
      </div>
      `;

    const REVIEW_BUTTON = document.querySelector("#tys-review");
    const REVIEW_CONTAINER = document.querySelector("#tys-review-container");
    REVIEW_BUTTON.addEventListener("click", () => {
      if (REVIEW_CONTAINER.classList.contains("hidden")) {
        if (!REVIEW_CONTAINER.dataset.rendered) {
          REVIEW_CONTAINER.innerHTML = buildTysReview(testDetails);
          REVIEW_CONTAINER.dataset.rendered = "true";
        }
        REVIEW_CONTAINER.classList.remove("hidden");
        REVIEW_BUTTON.textContent = "Hide Review";
      } else {
        REVIEW_CONTAINER.classList.add("hidden");
        REVIEW_BUTTON.textContent = "Review Answers";
      }
    });

    document.querySelector("#home").addEventListener("click", () => {
      window.location.hash = "";
    });

    document
      .querySelector("#tys-another")
      .addEventListener("click", async () => {
        // Reset to the language/question-count selection screen for a fresh test.
        await DB.test_questions.clear();
        state.current = "tys";
        state.previous = "tys-quiz-result";
        state.next = "tys-quiz";
        await updateStorage("states", state);
        TestPage(document.querySelector("main"));
      });
  }
}

// export default TestYourselfPage;
