import { DB, createStorage, getStorage, updateStorage } from "./storage.js";
import { randomID } from "./utils.js";

async function scratchPage(htmlEl) {
  const state = await DB.states.where("name").equals("general").last();

  if (state.current === "scratch") {
    htmlEl.innerHTML = `
        <div class="grid grid-cols-1 sm:grid-cols-2  lg:grid-cols-3 gap-4">
        <div class="bg-white rounded-lg shadow p-8">
          <h2 class="text-2xl font-bold mb-4">Course Title 1</h2>
          <p class="text-gray-700 mb-4">Course Description goes here. This is a brief overview of the course content.</p>
          <a href="#" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded inline-block">Enroll Now</a>
        </div>
        <div class="bg-white rounded-lg shadow p-8">
          <h2 class="text-2xl font-bold mb-4">Course Title 2</h2>
          <p class="text-gray-700 mb-4">Course Description goes here. This is a brief overview of the course content.</p>
          <a href="#" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded inline-block">Enroll Now</a>
        </div>
        <div class="bg-white rounded-lg shadow p-8">
          <h2 class="text-2xl font-bold mb-4">Course Title 3</h2>
          <p class="text-gray-700 mb-4">Course Description goes here. This is a brief overview of the course content.</p>
          <a href="#" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded inline-block">Enroll Now</a>
        </div>
        <!-- Repeat for the remaining courses -->
      </div><br /> <br /><br />
        `;
  }
}

export { scratchPage };
