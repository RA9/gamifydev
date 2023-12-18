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