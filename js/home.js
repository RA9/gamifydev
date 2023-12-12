import { scratchPage } from "./modules.js";
import QuizPage from "./quiz.js";
import { DB, createStorage, getStorage, updateStorage } from "./storage.js";
import { randomID } from "./utils.js";

async function HomePage(htmlEl) {
  let state = await DB.states.where("name").equals("general").last();

  if (!state) {
    state = await createStorage("states", {
      id: randomID(),
      previous: null,
      current: "home",
      next: "user_info",
      name: "general",
    });
  }

  if (state.current === "home") {
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
      state.current = "user_info";
      state.previous = "home";
      state.next = "preference";
      const page = document.querySelector("main");
      updateStorage("states", state);
      UserInfoSection(page);
    });
  } else if (state.current === "user_info") {
    const page = document.querySelector("main");
    UserInfoSection(page);
  } else if (state.current === "preference") {
    const page = document.querySelector("main");
    PreferenceSection(page);
  } else if (
    state.current === "quiz" ||
    state.current === "quiz-started" ||
    state.current === "quiz-completed"
  ) {
    const page = document.querySelector("main");
    QuizPage(page);
  } else if(state.current === "scratch") {
    const page = document.querySelector("main");

    scratchPage(page);
  }
}

async function UserInfoSection(htmlEl) {
  let state = await DB.states.where("name").equals("general").last();

  if (state.current === "user_info") {
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

  SUBMIT_NAME_BUTTON.addEventListener("click", async () => {
    const name = document.querySelector("#name").value;
    const preference = document.querySelector("#preference").value;
    let user = await DB.users.where("name").equals(name).last();

    if (name.length < 4) {
      return alert("Your name should be at least 4 characters!");
    }

    if (!user) {
      user = await createStorage("users", {
        id: randomID(),
        name,
        preference,
        // test: {},
      });
    }

    state.current = "preference";
    // user.name = name;
    // user.preference = preference;
    await updateStorage("states", state);
    // await createStorage("users", user);
    const page = document.querySelector("main");
    PreferenceSection(page);
  });
}

async function PreferenceSection(htmlEl) {
  const user = (await DB.users.toArray())[0];
  const state = await DB.states.where("name").equals("general").last();

  if (state.current === "preference") {
    htmlEl.innerHTML = `
  <div class="flex justify-center items-center">
  <div class="bg-white shadow-md rounded p-4">
  <h1 class="text-xl font-bold mb-4">Current Path Selected: ${user.preference.toUpperCase()}</h1>
  <p class="mb-4">Hi ${
    user.name
  }, welcome to GamifyDev. We are excited to have you here. 
  You have selected <b>${user.preference.toUpperCase()}</b> as your preferred path, click any of the buttons below to get started.
  </p>
  <div class="flex flex-col sm:flex-row justify-between items-center gap-4">
    <button id="get-started" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Take a quiz</button>
    <button id="scratch" class="w-full bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Start from scratch</button>
  </div>
</div>
</div><br/><br /> <br />
  `;

    const TAKE_QUIZ_BUTTON = document.querySelector("#get-started");
    const START_FROM_SCRATCH_BUTTON = document.querySelector("#scratch");

    TAKE_QUIZ_BUTTON.addEventListener("click", () => {
      state.current = "quiz";
      state.previous = "user_info";
      state.next = "preference";
      const page = document.querySelector("main");
      updateStorage("states", state);
      QuizPage(page);
    });

    START_FROM_SCRATCH_BUTTON.addEventListener("click", () => {
      state.current = "scratch";
      state.previous = "preference";
      state.next = null;

      console.log({ state });

      // window.location.href = "scratch.html";
      updateStorage("states", state);
      const page = document.querySelector("main");
      scratchPage(page);
    });
  } else if (state.current === "scratch") {
    const page = document.querySelector("main");
    scratchPage(page);
  }
}

export default HomePage;
