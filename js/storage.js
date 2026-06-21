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

// v9: daily-loop state.
//   daily — per-day record, used to mark days bridged by a streak freeze.
//   meta  — small key/value app state (streak freezes, milestone tracking).
// Dexie carries the v8 stores forward; existing users gain the new stores
// without losing any data.
db.version(9).stores({
  daily: "date",
  meta: "key",
});

// v10: spaced-repetition review items. One row per answered question, keyed by
// "<category>::<question text>", with a Leitner box and a due day-number.
db.version(10).stores({
  reviews: "key, category, due, box",
});

async function createQuestions() {
  try {
    const questions = await fetch("./data/questions.json");
    const data = await questions.json();

    // Sync each category from the source, keyed by question text:
    //  - de-duplicate: keep one row per question text, delete extra rows that
    //    earlier seeding may have created,
    //  - update existing questions whose content changed (e.g. a corrected
    //    answer), so content fixes reach existing users,
    //  - add questions not present yet (new categories, expanded banks).
    for (const key of Object.keys(data)) {
      const existingRows = await db.questions.where("category").equals(key).toArray();
      const byText = {};
      for (const r of existingRows) {
        const qt = r.details && r.details.question;
        if (!qt) continue;
        if (byText[qt]) {
          await db.questions.delete(r.id); // duplicate row — remove it
        } else {
          byText[qt] = r;
        }
      }

      for (const question of data[key]) {
        const details = { ...question, category: key };
        const existing = byText[question.question];
        if (existing) {
          if (JSON.stringify(existing.details) !== JSON.stringify(details)) {
            await db.questions.update(existing.id, {
              details,
              title: question.title || question.question,
            });
          }
        } else {
          const row = {
            id: randomID(),
            title: question.title || question.question,
            category: key,
            details,
            created_at: new Date(),
          };
          await db.questions.add(row);
          byText[question.question] = row; // guard against re-adding within run
        }
      }
    }
  } catch (error) {
    console.log(error);
    return null;
  }
}

async function createPaths() {
  try {
    const paths = await fetch("./data/app.json");
    const data = (await paths.json()).config.paths;

    // Upsert modules keyed by (path_name, title). app.json is the source of
    // truth for STRUCTURE (description, resources, next/previous) so the path
    // can grow over time, but a learner's PROGRESS (is_completed, current) is
    // preserved on modules they already have. This adds new paths (Backend,
    // Fullstack) and new modules to existing users without duplicating or
    // resetting anything.
    for (const path of data) {
      const pathName = path.name.toLowerCase();
      for (const mod of path.modules) {
        const existing = await db.paths
          .where("path_name")
          .equals(pathName)
          .and((p) => p.title === mod.title)
          .first();

        if (existing) {
          await db.paths.update(existing.id, {
            description: mod.description,
            resources: mod.resources,
            author_name: mod.author,
            previous: mod.previous,
            next: mod.next,
          });
        } else {
          await db.paths.add({
            id: randomID(),
            path_name: pathName,
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
        }
      }
    }
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
