import { randomID } from "./utils.js";

const db = new Dexie("gamifydev");

db.version(3).stores({
  questions: "id, title, category, details, created_at, updated_at",
  app: "id, theme,mode,language",
  users: "id, name, preference",
  states: "id, name, previous, next, current",
  tests: "id, name, language, numOfQuestions, is_completed, index",
  test_questions: "id, test_id, questions",
  scores: "id, test_id, score, numCorrect, numWrong, created_at, updated_at"
});

async function createQuestions() {
  try {
    const ques = await db.questions.toArray();

    if (ques.length > 0) return;

    const questions = await fetch("./data/questions.json");
    const data = await questions.json();

    console.log({ data });

    // join all properties
    const questionsData = Object.keys(data).map((key) => {
      return [
        ...data[key].map((question) => {
          return {
            ...question,
            category: key,
          };
        }),
      ];
    });

    questionsData.forEach(async (ques) => {
      ques.forEach(async (question) => {
        await db.questions.add({
          id: randomID(),
          title: question.title,
          category: question.category,
          details: question,
          created_at: new Date(),
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
    return await getStorage('states', storage);
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

createQuestions();

export { db as DB, getQuestions, createStorage, getStorage, updateStorage };
