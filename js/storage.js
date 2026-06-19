// import { randomID } from "./utils.js";

const db = new Dexie("gamifydev");

db.version(8).stores({
  questions: "id, title, category, details, created_at, updated_at",
  app: "id, theme,mode,language",
  paths:
    "id, path_name, title, description, resources, author_name, current, is_completed, next, previous, created_at, updated_at",
  // author: "id, name, email, created_at, updated_at",
  users: "id, name, preference",
  states: "id, name, previous, next, current",
  tests: "id, name, language, numOfQuestions, is_completed, created_at, updated_at",
  test_questions: "id, test_id, questions",
  scores: "id, test_id, score, numCorrect, numWrong, details, created_at, updated_at",
});

async function createQuestions() {
  try {
    const questions = await fetch("./data/questions.json");
    const data = await questions.json();

    // Seed per category and add only questions not already present (matched by
    // question text). This imports newly added categories (e.g. Java) and
    // expands existing ones (e.g. more SQL) for existing users, without
    // creating duplicates.
    for (const key of Object.keys(data)) {
      const existingRows = await db.questions.where("category").equals(key).toArray();
      const existingQuestions = new Set(existingRows.map((r) => r.details?.question));

      for (const question of data[key]) {
        if (existingQuestions.has(question.question)) continue;
        await db.questions.add({
          id: randomID(),
          title: question.title || question.question,
          category: key,
          details: { ...question, category: key },
          created_at: new Date(),
        });
      }
    }
  } catch (error) {
    console.log(error);
    return null;
  }
}

async function createPaths() {
  try {
    const path = await db.paths.toArray();
    const paths = await fetch("./data/app.json");
    const data = (await paths.json()).config.paths;

    // console.log({ data });

    if (path.length > 0 && data.length <= path.length) return;

    // join all properties
    const pathsData = data;

    pathsData.forEach(async (path) => {
      path.modules.forEach(async (mod) => {
        // console.log({ mod });
        await db.paths.add({
          id: randomID(),
          path_name: path.name.toLowerCase(),
          title: mod.title,
          description: mod.description,
          resources: mod.resources,
          author_name: mod.author,
          created_at: new Date(),
          current: mod.current,
          previous: mod.previous,
          next: mod.next,
          is_completed: mod.is_completed,
        });
      });
    });
  } catch (error) {
    console.log(error);
    return null;
  }
}

async function getQuestions(questionID) {
  try {
    const question = await db.questions.get(questionID);
    return question;
  } catch (error) {
    console.log(error);
    return null;
  }
}

async function createStorage(table, data) {
  try {
    const storage = await db.table(table).add(data);
    return await getStorage("states", storage);
  } catch (error) {
    console.log(error);
    return null;
  }
}

async function updateStorage(table, data) {
  try {
    const storage = await db.table(table).update(data.id, data);
    return storage;
  } catch (error) {
    console.log(error);
    return null;
  }
}

async function getStorage(table, id) {
  try {
    const storage = await db.table(table).get(id);
    return storage;
  } catch (error) {
    console.log(error);
    return null;
  }
}

function removeStorage(table, id) {
  try {
    db.table(table).delete(id);
  } catch (error) {
    console.log(error);
    return null;
  }
}

createQuestions();
createPaths();

const DB = db;

// export { db as DB, getQuestions, createStorage, getStorage, updateStorage };
