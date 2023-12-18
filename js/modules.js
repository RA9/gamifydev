// import { DB, createStorage, getStorage, updateStorage } from "./storage.js";
// import { randomID } from "./utils.js";

class Note {
  constructor(title, description, current = false, nextNote = null) {
    this.title = title;
    this.description = description;
    this.current = current;
    this.nextNote = nextNote;
  }
}

const createLinkedList = (notes) => {
  const noteMap = new Map();
  notes.forEach((note, index) => {
    noteMap.set(
      note.title,
      new Note(note.title, note.description, note.current)
    );
  });

  notes.forEach((note, index) => {
    if (note.next) {
      noteMap.get(note.title).nextNote = noteMap.get(note.next);
    }
  });

  let headOfLinkedList = null;
  notes.forEach((note, index) => {
    if (!note.previous) {
      headOfLinkedList = noteMap.get(note.title);
    }
  });

  return headOfLinkedList; // Return the head of the linked list
};

async function handleNotePage() {
  const state = await DB.states.where("name").equals("general").last();
  state.current = "note";
  state.previous = "scratch";
  state.next = "note-quiz";
  updateStorage("states", state);
  notePage(document.querySelector("main"));
}

async function fetchNotes(preference, noteFileName) {
  const response = await fetch(
    `notes/${(preference, toLowerCase())}/${noteFileName}`
  );
  return await response.text();
}

// Function to create interactive exercise HTML
function createInteractiveExercise(exercise) {
  return `
    <div class="p-4">
      <p class="text-xl font-semibold mb-4">${exercise.content}</p>
      ${createInputFields(exercise.input)}
      <button onclick="runTests()">Run Tests</button>
      <div id="results" class="mt-4"></div>
    </div>
  `;
}

// Function to create input fields for the exercise
function createInputFields(inputs) {
  return inputs
    .map(
      (input) => `
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700">${
        input.label
      }</label>
      <textarea
        id="${input.label.toLowerCase().replace(/\s+/g, "-")}"
        class="mt-1 p-2 block w-full rounded-md border border-gray-300"
        rows="4"
        placeholder="${input.value}"
      ></textarea>
      <p class="mt-2 text-sm text-gray-500">${input.hint}</p>
    </div>
  `
    )
    .join("");
}

// Function to run tests on user input
function runTests() {
  const userInput = collectUserInput();
  const results = runTestsOnInput(
    userInput,
    exercises[currentExerciseIndex].tests
  );
  displayResults(results);
}

// Function to collect user input from input fields
function collectUserInput() {
  const userInputs = exercises[currentExerciseIndex].input.map((input) => {
    const inputId = input.label.toLowerCase().replace(/\s+/g, "-");
    return document.getElementById(inputId).value;
  });

  return userInputs.join("\n");
}

// Function to run tests on user input and compare with expected output
function runTestsOnInput(userInput, tests) {
  const results = tests.map((test) => {
    try {
      const userOutput = eval(userInput);
      const expectedOutput = eval(test.input);

      return {
        input: test.input,
        passed: userOutput === expectedOutput,
      };
    } catch (error) {
      return {
        input: test.input,
        passed: false,
        error: error.message,
      };
    }
  });

  return results;
}

// Function to display results on the UI
function displayResults(results) {
  const resultsContainer = document.getElementById("results");
  resultsContainer.innerHTML = ""; // Clear previous results

  results.forEach((result) => {
    const resultDiv = document.createElement("div");
    resultDiv.textContent = `Input: ${result.input}, Passed: ${
      result.passed ? "Yes" : "No"
    }`;
    resultDiv.style.color = result.passed ? "green" : "red";

    if (!result.passed && result.error) {
      const errorDiv = document.createElement("div");
      errorDiv.textContent = `Error: ${result.error}`;
      resultDiv.appendChild(errorDiv);
    }

    resultsContainer.appendChild(resultDiv);
  });
}



// Function to fetch notes, merge exercises, and display them
function mergeNotesAndExercises(currentExerciseIndex = 0) {
  fetchNotes(exercises[currentExerciseIndex].noteFileName)
    .then((noteContent) => {
      const exerciseContainer = document.getElementById("exercise-container");
      exerciseContainer.innerHTML = createInteractiveExercise(
        exercises[currentExerciseIndex]
      );
    })
    .catch((error) =>
      console.error(
        `Error merging notes for ${exercises[currentExerciseIndex].name}:`,
        error
      )
    );
}

async function scratchPage(htmlEl) {
  const state = await DB.states.where("name").equals("general").last();
  const user = (await DB.users.toArray())[0];
  const notes = await DB.paths
    .where("path_name")
    .equals(user.preference)
    .toArray();

  const scratchNotes = createLinkedList(notes);

  let currentNote = scratchNotes;
  const scratchNotesArray = [];

  while (currentNote !== null) {
    // console.log(currentNote.title);
    // disable all buttons except the currentNote
    if (currentNote.current) {
      scratchNotesArray.push(`
    <div class="bg-white rounded-lg shadow p-8">
    <h2 class="text-2xl font-bold mb-4">${currentNote.title}</h2>
    <p class="text-gray-700 mb-4">${currentNote.description}</p>
    <button onclick="handleNotePage()" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded inline-block">Enroll Now</button>
    </div>
    `);
    } else {
      scratchNotesArray.push(`
    <div class="bg-white rounded-lg shadow p-8">
    <h2 class="text-2xl font-bold mb-4">${currentNote.title}</h2>
    <p class="text-gray-700 mb-4">${currentNote.description}</p>
    <button onclick="handleNotePage()" class="bg-blue-500 hover:cursor-not-allowed hover:bg-blue-600 text-white font-bold py-2 px-4 rounded inline-block" disabled>Enroll Now</a>    
  </div>
    `);
    }

    currentNote = currentNote.nextNote;
  }

  if (state.current === "scratch") {
    htmlEl.innerHTML = `
        <div class="grid grid-cols-1 sm:grid-cols-2  lg:grid-cols-3 gap-4">
        ${scratchNotesArray.join("")}
      </div><br /> <br /><br />
        `;
  } else if (state.current === "note") {
    return notePage(document.querySelector("main"));
  }
}

async function notePage(htmlEl) {
  const state = await DB.states.where("name").equals("general").last();
  const user = (await DB.users.toArray())[0];
  const note = await DB.paths
    .where("path_name")
    .equals(user.preference)
    .and((c) => c.current)
    .last();

  console.log({ note });

  if (state.current === "note") {
    htmlEl.innerHTML = `
    <div id="exercise-container"></div>
    `;

    mergeNotesAndExercises();
  }
}

// export { scratchPage };
