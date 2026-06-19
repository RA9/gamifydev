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
    htmlEl.innerHTML = `
      <div class="max-w-4xl mx-auto bg-white rounded-lg shadow p-8">
    <div class="mb-4">
      <label for="language" class="block text-gray-700 text-sm font-bold mb-2">Choose a language:</label>
      <select id="language" name="language" class="block w-full bg-white border border-gray-400 rounded py-2 px-3">
        <option value="html">HTML</option>
        <option value="css">CSS</option>
        <option value="javascript">JavaScript</option>
        <option value="c">C</option>
        <option value="python">Python</option>
        <option value="java">Java</option>
        <option value="sql">SQL</option>
      </select>
    </div>
    <div class="mb-4">
      <label for="numQuestions" class="block text-gray-700 text-sm font-bold mb-2">Number of questions:</label>
      <input id="numQuestions" name="numQuestions" type="number" value="20" class="block w-full bg-white border border-gray-400 rounded py-2 px-3">
    </div>
    <button id="start-tys" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Start</button>
  </div>
      `;

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
    <div class="mb-4 px-2">
            <input type="radio" id="${id}" name="option" value="${escapeHTMLToEntities(
      option
    )}">
            <label for="${id}">${escapeHTMLToEntities(option)}</label>
    </div>
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
      ).innerHTML = `<span class="text-red-400">Time is up!</span>`;
      // Lock in the current answers and auto-submit when time runs out.
      document.querySelectorAll("input[type=radio]").forEach((el) => {
        el.classList.add("cursor-not-allowed");
        el.setAttribute("disabled", true);
      });
      document.querySelector("#tys-submit")?.click();
    } else {
      document.querySelector(
        "#tys-duration"
      ).innerHTML = `<span class="text-green-400">${minutes}:${seconds} seconds</span>`;
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
    const elHTML = [].concat(testQuestions.questions).map((question, index) => {
      return `
      <div class="question mb-2">
      <p class="text-lg font-bold mb-4">${index + 1}. ${escapeHTMLToEntities(
        question.details.question
      )}</p>
      <form>
        ${tysRandomizeOptions(question.details.options).join("")}
      </form>
      </div>
      <br/>
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
        <div class="max-w-6xl mx-auto bg-white rounded-lg shadow p-8">
        <div class="flex justify-between gap-4">
         <p class="text-lg font-bold mb-4">Language: ${test.language.toUpperCase()}</p>
          <p id="tys-duration" class="text-lg font-bold mb-4"></p>
        </div>
        ${elHTML.join("")} <br/>
        <div class="flex justify-between gap-4">
        <button id="tys-cancel" class="w-full bg-red-500 hover:bg-reds-700 text-white font-bold py-2 px-4 rounded">Cancel</button>
        <button id="tys-submit" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Submit</button>
      </div>
      </div><br/><br/><br/>
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
          let cls = "border-gray-200 text-gray-700";
          let tag = "";
          if (isCorrectOpt) {
            cls = "border-green-500 bg-green-50 text-green-800";
            tag = ` <span class="font-semibold">✓ correct answer</span>`;
          } else if (isChosenWrong) {
            cls = "border-red-500 bg-red-50 text-red-800";
            tag = ` <span class="font-semibold">✗ your answer</span>`;
          }
          return `<li class="border ${cls} rounded px-3 py-2 mb-1 list-none">${escapeHTMLToEntities(
            opt
          )}${tag}</li>`;
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
    htmlEl.innerHTML = `
      <div class="max-w-6xl mx-auto bg-white rounded-lg shadow p-8">
        <h1 class="font-medium text-2xl">Results of ${test.language.toUpperCase()} test.</h1>
        <p class="py-4">You got ${testDetails.numCorrect} correct and ${
      testDetails.numWrong
    } wrong. You total score is <span class="text-${
      testDetails.score > 70
        ? testDetails.score > 85
          ? "green"
          : "yellow"
        : "red"
    }-500">${testDetails.score}%</span>.

     </p>
        <button id="tys-review" class="w-full mb-4 bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded">Review Answers</button>
        <div id="tys-review-container" class="hidden mb-4"></div>
        <div class="flex justify-between gap-4">
          <button id="home" class="w-full bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded">Back to Home</button>
          <button id="tys-another" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Take Another</button>
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
