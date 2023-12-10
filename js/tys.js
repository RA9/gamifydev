import { createStorage, getStorage } from "./storage.js";
import questions from "../data/questions.json" assert { type: "json" };

function TestYourselfPage(htmlEl) {
  const state = getStorage("state");
  const user = getStorage("user");

  if (!user) {
    const newUser = {
      name: "",
      preference: "",
      test: {},
    };
    createStorage("user", newUser);
  }

  if (!state) {
    const newState = {
      currentState: "home",
      test: {},
    };
    createStorage("state", newState);
  }

  if (state.test.currentState !== "tys-quiz") {
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

    START_TYS_BUTTON.addEventListener("click", () => {
      const language = document.querySelector("#language").value;
      const numQuestions = document.querySelector("#numQuestions").value;
      const state = getStorage("state");
      const user = getStorage("user");

      if (!user) {
        const newUser = {
          name: "",
          preference: "",
          test: {},
        };
        createStorage("user", newUser);
      }

      if (!state) {
        const newState = {
          currentState: "home",
          test: {},
        };
        createStorage("state", newState);
      }

      state.test.currentState = "tys-quiz";
      user.test.language = language;
      user.test.numQuestions = numQuestions;
      createStorage("state", state);
      createStorage("user", user);
      const page = document.querySelector("main");
      TestYourselfSection(page);
    });
  } else if (state.test.currentState === "tys-quiz") {
    const page = document.querySelector("main");

    TestYourselfSection(page);
  }
}

function randomizedQuestions(questions, limit) {
  questions.sort(() => Math.random() - 0.5);
  const questionsLimit = questions.slice(0, limit);
  return questionsLimit;
}

function getQuestions(questions,  limit) {
  return randomizedQuestions(questions, limit)
}

function TestYourselfSection(htmlEl) {
  const user = getStorage("user");
  const state = getStorage("state");

  if (!user) {
    const newUser = {
      name: "",
      preference: "",
      test: {},
    };
    createStorage("user", newUser);
  }

  if(!user.test.questions) {
    user.test.questions = getQuestions(questions[user.test.language], Number(user.test.numQuestions))
  }

  if (!state) {
    const newState = {
      currentState: "home",
      test: {},
    };
    createStorage("state", newState);
  }

  if (state.test.currentState === "tys-quiz") {
    htmlEl.innerHTML = `
        <div class="max-w-6xl mx-auto bg-white rounded-lg shadow p-8">
        <div class="flex justify-between items-center mb-4">
          <p class="text-gray-700 font-bold">Question 1 of 20</p>
          <p class="text-gray-700">Duration: 30 seconds</p>
        </div>
        <p class="text-lg font-bold mb-4">What is the purpose of the 'def' keyword in Python?</p>
        <form>
          <div class="mb-4">
            <input type="radio" id="option1" name="option" value="To define a class">
            <label for="option1">To define a class</label>
          </div>
          <div class="mb-4">
            <input type="radio" id="option2" name="option" value="To define a function">
            <label for="option2">To define a function</label>
          </div>
          <div class="mb-4">
            <input type="radio" id="option3" name="option" value="To declare a variable">
            <label for="option3">To declare a variable</label>
          </div>
          <div class="mb-4">
            <input type="radio" id="option4" name="option" value="To import a module">
            <label for="option4">To import a module</label>
          </div>
        </form>
        <div class="flex justify-between gap-4">
         <button id="previous" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Previous</button>
         <button id="next" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Next</button>
         <button id="submit" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Submit</button>
        </div>
      </div><br/>
    `;

    const PREVIOUS_BUTTON = document.querySelector("#previous");
    const NEXT_BUTTON = document.querySelector("#next");
    const SUBMIT_BUTTON = document.querySelector("#submit")

    PREVIOUS_BUTTON.addEventListener("click", () => {
      state.currentState = "quiz";
      const page = document.querySelector("main");
      createStorage("state", state);
      QuizPage(page);
    });

    NEXT_BUTTON.addEventListener("click", () => {
      state.currentState = "scratch";
      const page = document.querySelector("main");
      createStorage("state", state);
      ScratchPage(page);
    });
  }
}

export default TestYourselfPage;
