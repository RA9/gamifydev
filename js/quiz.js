import { createStorage, getStorage } from "./storage.js";
import data from "../data/questions.json" assert { type: "json" };

function countDown(prop, element) {
  const timer = document.querySelector("#timer");
  let time = prop.duration;
  const interval = setInterval(() => {
    if (time === 0) {
      clearInterval(interval);
      timer.innerHTML = `<span class="text-red-400">Time is up!</span>`;
      // disabled all input element
      element.querySelectorAll("input").forEach((el) => {
        if (el.checked) {
          storeAnswer(el.value);
        } else {
          storeAnswer(null);
        }
        el.classList.add("cursor-not-allowed");
        el.setAttribute("disabled", true);
      });
    } else {
      timer.innerHTML = `<span class="text-green-400">${time} seconds</span>`;
      time--;
    }
  }, 1000);
}

function getQuestions(limit) {
  const questions = data.frontend;
  // randomize questions without repeating them
  questions.sort(() => Math.random() - 0.5);
  const questionsLimit = questions.slice(0, limit);
  return questionsLimit;
}

function storeAnswer(answer) {
  const user = getStorage("users");
  const questions = getStorage("questions");
  const state = getStorage("states");

  if (!user.currentQuiz) {
    user.currentQuiz = {};
  }

  user.currentQuiz[questions[user.index].question] = {
    answer,
  };

  updateStorage("user", user);
  updateStorage("state", state);
}

function QuizPage(htmlEl) {
  const user = getStorage("user");
  const state = getStorage("state");

  if (!user) {
    const newUser = {
      name: "",
      preference: "",
    };
    createStorage("user", newUser);
  }

  if (!state) {
    const newState = {
      currentState: "quiz",
    };
    createStorage("state", newState);
  }

  if (user.name && state.currentState === "quiz") {
    htmlEl.innerHTML = `
    <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <div class="flex flex-col justify-center items-center">
        <h1 class="text-3xl font-bold">Technical Test: ${user.preference.toUpperCase()}</h1>
        <p class="leading text-xl my-4">
        Hi ${
          user.name
        }, I am glad you've made it here. We have put together a list of questions that will help us determine the preferred section within your selected path.
        It will only take about 15 minutes. Click the button below to get started.
        </p>
        <button id="start-quiz" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Start Quiz</button>
    </div>
    </div>
    `;

    const START_QUIZ_BUTTON = document.querySelector("#start-quiz");

    START_QUIZ_BUTTON.addEventListener("click", () => {
      state.currentState = "quiz-started";
      const page = document.querySelector("main");

      updateStorage("state", state);
      QuizStartedPage(page);
    });
  } else if (!user.name) {
    htmlEl.innerHTML = ``;
  } else if (user.name && state.currentState === "quiz-started") {
    const page = document.querySelector("main");

    QuizStartedPage(page);
  } else if (state.currentState === "quiz-completed") {
    // console.log("quiz-completed");
    const page = document.querySelector("main");

    showQuizResultsPage(page);
  }
}

function getRandomItem(limit) {
  return Math.floor(Math.random() * limit + 1);
}

// https://stackoverflow.com/questions/2613582/convert-tags-to-html-entities
function htmlencode(str) {
  return str.replace(/[&<>"']/g, function ($0) {
    return (
      "&" +
      { "&": "amp", "<": "lt", ">": "gt", '"': "quot", "'": "#39" }[$0] +
      ";"
    );
  });
}

function randomizeOptions(options, answer = null) {
  const randomOptions = options.sort(() => Math.random() - 0.5);
  return randomOptions.map((option) => {
    // console.log({ answer, option });
    return `<input type="radio" id="${getRandomItem(options.length)}" ${
      answer && htmlencode(answer) === option ? "checked" : ""
    } name="favoriteLanguage" value="${option}">
    <label for="${option}">${option}</label><br>`;
  });
}

function disableOptions() {
  const options = document
    .querySelector("#quiz-options")
    .querySelectorAll("input");
  options.forEach((el) => {
    el.classList.add("cursor-not-allowed");
    el.setAttribute("disabled", true);
  });
}

function QuizStartedPage(htmlEl) {
  const user = getStorage("user");
  const state = getStorage("state");

  let questions = getStorage("questions");
  if (!questions) {
    questions = createStorage("questions", getQuestions(30));
  }

  if (!user) {
    const newUser = {
      name: "",
      preference: "",
    };
    createStorage("user", newUser);
  }

  if (!state) {
    const newState = {
      currentState: "quiz-started",
    };
    createStorage("state", newState);
  }

  if (!state.quiz) {
    state.quiz = {
      current: "current",
      next: "next",
    };
  }

  if (user.name && state.currentState === "quiz-started") {
    if (!user.index) {
      user.index = 0;
      updateStorage("user", user);
    }

    if (!user.currentQuiz) {
      user.currentQuiz = {};
    }
    const previousAnswer =
      user?.currentQuiz[questions[user.index].question]?.answer || null;
    htmlEl.innerHTML = `
    <div id="quiz-container" class="flex justify-center items-center">
    <div class="w-full bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <div class="flex justify-between">
    <span class="text-xl font-bold mb-4">Question ${user.index + 1} of ${
      questions.length
    }</span>
    <span id="timer" class="text-xl font-bold mb-4"></span>
    </div>
    <p class="mb-4">${questions[user.index].question}</p>
    <div id="quiz-options" class="mb-4">
      ${randomizeOptions(questions[user.index].options, previousAnswer).join(
        ""
      )}
    </div>
    <div class="flex justify-between gap-4">
      ${
        user.index > 0
          ? '<button id="previous" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Previous</button>'
          : ""
      }
      ${
        user.index < questions.length - 1
          ? '<button id="next" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Next</button>'
          : '<button id="submit" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Submit</button>'
      }
    </div>
  </div>
    </div>
    `;

    if (
      !Object.keys(user.currentQuiz).includes(questions[user.index].question)
    ) {
      countDown(questions[user.index], document.querySelector("#quiz-options"));
    }

    if (
      Object.keys(user.currentQuiz).includes(questions[user.index].question)
    ) {
      disableOptions();
    }

    const NEXT_BUTTON = document.querySelector("#next");
    const PREVIOUS_BUTTON = document.querySelector("#previous");
    const SUBMIT_BUTTON = document.querySelector("#submit");

    if (user.index < questions.length - 1) {
      NEXT_BUTTON.addEventListener("click", () => {
        state.currentState = "quiz-started";
        const page = document.querySelector("main");
        const options = document
          .querySelector("#quiz-options")
          .querySelectorAll("input");

        options.forEach((curr) => {
          if (curr.checked) {
            const answer = user.currentQuiz;
            answer[questions[user.index].question] = {
              answer: curr.value,
            };
          }
        });

        if (user.index >= 0 && user.index < questions.length) {
          user.index += 1;
        }

        state.quiz = {
          current: "current",
          next: "next",
        };
        updateStorage("user", user);
        updateStorage("state", state);

        QuizStartedPage(page);
      });
    } else if (user.index === questions.length - 1) {
      SUBMIT_BUTTON.addEventListener("click", () => {
        state.currentState = "quiz-completed";
        const page = document.querySelector("main");

        const options = document
          .querySelector("#quiz-options")
          .querySelectorAll("input");
        options.forEach((curr) => {
          if (curr.checked) {
            const answer = user.currentQuiz;
            answer[questions[user.index].question] = {
              answer: curr.value,
            };
          }
        });

        updateStorage("user", user);
        updateStorage("state", state);

        showQuizResultsPage(page);
      });
    }

    if (user.index > 0) {
      PREVIOUS_BUTTON.addEventListener("click", () => {
        state.currentState = "quiz-started";
        const page = document.querySelector("main");

        if (user.index <= 0) {
          user.index = 0;
        } else {
          user.index -= 1;
        }

        state.quiz = {
          current: "previous",
          next: "next",
        };
        updateStorage("user", user);
        updateStorage("state", state);

        QuizStartedPage(page);
      });
    }
  } else if (!user.name) {
    htmlEl.innerHTML = ``;
  }
}

function showQuizResultsPage(htmlEl) {
  const user = getStorage("user");
  const state = getStorage("state");

  if (!user) {
    const newUser = {
      name: "",
      preference: "",
    };
    createStorage("user", newUser);
  }

  if (!state) {
    const newState = {
      currentState: "quiz-results",
    };
    createStorage("state", newState);
  }

  if (user.name && state.currentState === "quiz-completed") {
    // const page = document.querySelector("main");
    const questions = getStorage("questions");
    const userAnswers = Object.values(user.currentQuiz);
    const correctAnswers = questions.map((question) => {
      return question.answer;
    });

    
    const correctAnswersCount = userAnswers.filter((answer, index) => {
      return answer.answer === correctAnswers[index];
    }).length;

    const percentage = Math.round(
      (correctAnswersCount / questions.length) * 100
    );

    if (!user.score) {
      user.score = 0;
      updateStorage("user", user);
    } 
    
    if(user.score < percentage){
      user.score = percentage;
      updateStorage("user", user);
    }



    htmlEl.innerHTML = `
    <div class="flex justify-center items-center">
    <div class="w-full bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <div class="flex flex-col justify-center items-center">
        <h1 class="text-3xl font-bold">Technical Test for  ${user.preference.toUpperCase()} Results</h1>
        <p class="leading text-xl my-4">
        Hi ${
          user.name
        }, you have completed the technical test for ${user.preference.toUpperCase()} path.
        Out of ${
          questions.length
        } questions, you got ${correctAnswersCount} questions correctly. And your score is ${
      percentage > 80
        ? '<span class="text-green-400">' + percentage + "</span>"
        : '<span class="text-red-400">' + percentage + "</span>"
    }%.
        You can now proceed to the next path.
        </p>
        <button id="my-path" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">My Path</button>
    </div>
    </div>
    </div>
    `;

    const MY_PATH_BUTTON = document.querySelector("#my-path");

    MY_PATH_BUTTON.addEventListener("click", () => {
      state.currentState = "home";
      const page = document.querySelector("main");

      updateStorage("state", state);
      HomePage(page);
    });
  }
}

export default QuizPage;
