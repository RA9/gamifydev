function randomID() {
    // generate random id
    const str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
    const limit = 8;
  
    let id = "";
    for (let i = 0; i < limit; i++) {
      id += str[Math.floor(Math.random() * str.length)];
    }
  
    return id;
  }

// Fisher–Yates shuffle. Returns a NEW array; does not mutate the input.
// Replaces `arr.sort(() => Math.random() - 0.5)`, which is statistically biased.
function shuffle(array) {
  const result = [...array];
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [result[i], result[j]] = [result[j], result[i]];
  }
  return result;
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

function escapeHTML(htmlStr) {
  return htmlStr.replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#39;");        

}

function escapeHTMLToEntities(str) {
  return str.replace(/[\u00A0-\u9999<>\&]/gim, function(i) {
    return '&#'+i.charCodeAt(0)+';';
  });
}

  // export { randomID, htmlencode };