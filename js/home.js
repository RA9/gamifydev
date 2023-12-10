import QuizPage from "./quiz.js";
import { createStorage, getStorage } from "./storage.js";

function HomePage(htmlEl) {
  let state = getStorage("state");

  if (!state) {
    const newState = {
      currentState: "home",
      test: {}
    };
    state = createStorage("state", newState);
  }

  if (state.currentState === "home") {
    htmlEl.innerHTML = `
  <div class="max-w-7xl">
    <div class="bg-white shadow p-8 rounded-lg">
      <div class="card-content">
            <h1 class="text-2xl font-bold mb-4">Welcome to GamifyDev v1.0</h1>
            <p class="text-xl font-sm p-2">
              Are you a beginner who struggles to learn how to code? Do you want
              to learn how to code but don't know where to start? Do you want to
              learn how to code but don't have access to the internet? Do you
              want to learn how to code but don't have a laptop? If you answered
              yes to any of these questions, then you are in the right place.
            </p>
            <p class="text-xl font-sm p-2">
              Yeah! You heard me right, you're in the right place. 
              Actually, you do need a laptop and a good internet connection to get started. 
              But, you don't need to know how to code to get started. 
              You will learn how to code by building real-world projects. So, go and sort
              out your internet issues, and get a laptop if you don't have one, before moving on.
            </p>
            <p class="text-xl font-sm p-2">
              <span style="color: orange">Warning:</span> This is not a get rich quick scheme. You will not become a 
              software developer overnight. You will not become a software developer 
              in a week. You will not become a software developer in a month. You will
              not become a software developer in a year. We are not here to teach you how to become a software developer.
              We are here to provide you with the fundamentals of programming and help you build real-world projects.
               It is up to you to decide what you want to do with the knowledge you will gain from this platform.
            </p><br/><br />

            <h2 class="text-2xl font-bold">A message from the creator</h2>
            <p class="text-xl font-sm p-2">
              Hi there, my name is <a class="underline text-blue-400" href="https://twitter.com/rademejs" target="_blank">Carlos S. Nah</a>, I am a software engineer and I am
              excited to have you here. I have been a software engineer for over 5 years and I have been teaching people how to code for over 3+ years.<br /><br /> And on this journey; it is 
              my hope that you learn enough to be dangerous. If at any point you need to dive deeper; I recommend you give 
              <a class="underline text-blue-400" href="https://twitter.com/rademejs" target="_blank">Kwagei Innovators Training</a> a try.<br/>
              <b>To get started, click the button below.</b>
            </p>
            <button id="get-started" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Get Started</button>
      </div>
    </div>
  </div> <br><br>
    `;

    const GET_STARTED_BUTTON = document.querySelector("#get-started");

    // Show the preference card
    GET_STARTED_BUTTON.addEventListener("click", () => {
      state.currentState = "user_info";
      const page = document.querySelector("main");
      createStorage("state", state);
      UserInfoSection(page);
    });
  } else if (state.currentState === "user_info") {
    const page = document.querySelector("main");
    UserInfoSection(page);
  } else if (state.currentState === "preference") {
    const page = document.querySelector("main");
    PreferenceSection(page);
  } else if (
    state.currentState === "quiz" ||
    state.currentState === "quiz-started" ||
    state.currentState === "quiz-completed"
  ) {
    console.log("Yeah");
    const page = document.querySelector("main");
    QuizPage(page);
  }
}

function UserInfoSection(htmlEl) {
  let state = getStorage("state");

  if (!state) {
    const newState = {
      currentState: "user_info",
      test: {}
    };
    state = createStorage("state", newState);
  }

  if (state.currentState === "user_info") {
    htmlEl.innerHTML = `
  <div class="card" style="background: #fff; margin: auto; box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2); transition: 0.3s; border-radius: 5px;">
  <div class="card-content" style="padding: 2em;">
      <p class="card-text" style="margin-bottom: 1em;">Enter your name and select your preferred field below to get started.</p>
      <div class="card-action" style="display: flex; flex-direction: column; gap: 1em;">
      <input type="text" id="name" placeholder="Enter your name" style="padding: 0.5em; font-size: 1em; border: 1px solid #757195; border-radius: 5px;"/>
      <select id="preference" style="padding: 0.5em; font-size: 1em; border: 1px solid #757195; border-radius: 5px;">
      <option value="">Choose your field</option>
  <option value="fullstack">Fullstack</option>
  <option value="frontend">Frontend</option>
  <option value="backend">Backend</option>
  <option value="notsure">Not sure</option>
</select>
          <button id="submit-name" class="submit-button" style="padding: 0.5em; font-size: 1em; background-color: #757195; color: white; border: none; border-radius: 5px; cursor: pointer;">Submit</button>
      </div>
  </div>
</div>
  `;
  }

  const SUBMIT_NAME_BUTTON = document.querySelector("#submit-name");

  SUBMIT_NAME_BUTTON.addEventListener("click", () => {
    const name = document.querySelector("#name").value;
    const preference = document.querySelector("#preference").value;
    let user = getStorage("user");

    if (!user) {
      user = createStorage("user", {
        name: "",
        preference: "",
        test: {}
      });
    }

    state.currentState = "preference";
    user.name = name;
    user.preference = preference;
    createStorage("state", state);
    createStorage("user", user);
    const page = document.querySelector("main");
    PreferenceSection(page);
  });
}

function PreferenceSection(htmlEl) {
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
      currentState: "home",
    };
    createStorage("state", newState);
  }

  if (state.currentState === "preference") {
    htmlEl.innerHTML = `
  <div class="flex justify-center items-center">
  <div class="bg-white shadow-md rounded p-8">
  <h1 class="text-xl font-bold mb-4">Current Path Selected: ${user.preference.toUpperCase()}</h1>
  <p class="mb-4">Hi ${
    user.name
  }, welcome to GamifyDev. We are excited to have you here. 
  You have selected <b>${user.preference.toUpperCase()}</b> as your preferred path, click any of the buttons below to get started.
  </p>
  <div class="flex flex-col md:flex-row justify-between items-center gap-4">
    <button id="get-started" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Take a quiz</button>
    <button id="get-started" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Start from scratch</button>
  </div>
</div>
</div><br/>
  `;

    const TAKE_QUIZ_BUTTON = document.querySelector("#get-started");
    const START_FROM_SCRATCH_BUTTON = document.querySelector("#get-started");

    TAKE_QUIZ_BUTTON.addEventListener("click", () => {
      state.currentState = "quiz";
      const page = document.querySelector("main");
      createStorage("state", state);
      QuizPage(page);
    });

    START_FROM_SCRATCH_BUTTON.addEventListener("click", () => {
      state.currentState = "scratch";
      // window.location.href = "scratch.html";
    });
  }
}

export default HomePage;
