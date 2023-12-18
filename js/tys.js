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
    <h1 class="text-3xl font-bold mb-4">Test Yourself</h1>
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

      // const testArr = await DB.tests.where("name").startsWith("tys").toArray();

      await createStorage("tests", {
        id: randomID(),
        name: "tys-" + randomID(),
        is_completed: false,
        language: language,
        numQuestions,
      });

      await updateStorage("states", state);

      const page = document.querySelector("main");
      TestYourselfSection(page);
    });
  } else if (state.current === "tys-quiz") {
    const page = document.querySelector("main");

    TestYourselfSection(page);
  }
}

function tysRandomizeOptions(options, answer = null) {
  const randomOptions = options.sort(() => Math.random() - 0.5);
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
  questions.sort(() => Math.random() - 0.5);
  const questionsLimit = questions.slice(0, limit);
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
      // disabled all input element
      // document.querySelectorAll("input[type=radio]").forEach((el) => {
      //   if (el.checked) {
      //     storeAnswer(el.value);
      //   } else {
      //     storeAnswer(null);
      //   }
      //   el.classList.add("cursor-not-allowed");
      //   el.setAttribute("disabled", true);
      // });
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

  let testQuestions = await DB.table("test_questions")
    .where("test_id")
    .equals(test.id)
    .toArray();

  console.log({ yep: testQuestions });

  if (testQuestions.length <= 0) {
    const questions = await DB.questions.toArray();

    testQuestions = await createStorage("test_questions", {
      test_id: test.id,
      id: randomID(),
      questions: getQuestionsByLimit(
        questions.filter((question) => question.category === test.language),
        Number(test.numQuestions)
      ),
    });

    console.log({ yah: testQuestions})
  }

  console.log({ testQuestions });

  testQuestions = Array.isArray(testQuestions)
    ? testQuestions[0]
    : testQuestions;

  console.log({ testQuestions });

  if (state.current === "tys-quiz") {
    const elHTML = [].concat(testQuestions.questions).map((question, index) => {
      return `
      <p class="text-lg font-bold mb-4">${index + 1}. ${escapeHTMLToEntities(
        question.details.question
      )}</p>
      <form>
        ${tysRandomizeOptions(question.details.options).join("")}
      </form>
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
        <button id="submit" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Submit</button>
      </div>
      </div><br/><br/><br/>
    `;

    tysCountDown(sumDuration);

    const SUBMIT_BUTTON = document.querySelector("#submit");
    const CANCEL_BUTTON = document.querySelector("#tys-cancel");

    CANCEL_BUTTON.addEventListener("click", () => {
      state.current = "tys";
      state.previous = "tys-quiz";
      state.next = null;
      updateStorage("states", state);

      // reload page and remove test_questions
      window.location.reload();
      removeStorage("test_questions", testQuestions.id);
    });

    SUBMIT_BUTTON.addEventListener("click", () => {
      state.currentState = "quiz";
      const page = document.querySelector("main");
      createStorage("state", state);
      QuizPage(page);
    });
  }
}

async function TestResultPage(htmlEl) {
  const state = (await DB.states.where("name").equals("tys").toArray())[0];
  const test = await DB.tests
    .where("name")
    .startsWith("tys")
    .and((test) => !test.is_completed)
    .last();

  if (state.current === "tys-result") {
    htmlEl.innerHTML = `
      <div class="max-w-6xl mx-auto bg-white rounded-lg shadow p-8">
        <h1>Results of ${test.language.toUpperCase()} test.</h1>
        <p></p>
        <div class="flex justify-between gap-4">
          <button id="home" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Go Home</button>
          <button id="tys-another" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Take Another</button>
        </div>
      </div>
      `;
  }
}

// export default TestYourselfPage;
