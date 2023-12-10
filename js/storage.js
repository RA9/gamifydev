/*
* @description - Creates a storage
* @param {string} name - The name of the storage
* @param {object|array} value - The value to store
* @returns {object} - The value stored
*/
export function createStorage(name, value) {
  localStorage.setItem(name, JSON.stringify(value));

  return value;
}

/*
* @description - Get the value stored in the storage
* @param {string} name - The name of the storage
* @returns {object} - The value stored
*/
export function getStorage(name) {
  return JSON.parse(localStorage.getItem(name));
}

/*
* @description - Remove the storage
* @param {string} name - The name of the storage
*/
export function removeStorage(name) {
  localStorage.removeItem(name);
}

/*
* @description - Update the storage
* @param {string} name - The name of the storage
* @param {object|array} value - The value to store
* @returns {object} - The value stored
*/
export function updateStorage(name, value) {
    // get the value from the storage
    const storageValue = getStorage(name);
    // check if the value exists
    if (!storageValue) {
      // if it doesn't exist, create it
      return createStorage(name, value);
    }

    // if it exists, update it
    const updatedValue = { ...storageValue, ...value };
    
  return createStorage(name, updatedValue);;
}
