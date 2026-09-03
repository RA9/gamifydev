// Timezone band picker on the placement result page.
//
// The band decides which cohort a learner can join and when their standup
// window opens, so getting it roughly right matters. We pre-select from the
// browser's own timezone offset — but it stays a visible, editable choice
// rather than a silent guess, because travellers and VPNs exist.
//
// The server treats an unknown value as the default band, so a failure here
// degrades to "Europe & Africa" rather than blocking enrollment.
(function () {
  var select = document.getElementById("bandPick");
  if (!select) return;

  var fields = document.querySelectorAll(".pl-band-field");

  function bandForOffset(minutesWestOfUTC) {
    // getTimezoneOffset is positive west of UTC, so invert for a UTC offset.
    var hours = -minutesWestOfUTC / 60;
    if (hours <= -3) return "americas";
    if (hours >= 4) return "asia_pacific";
    return "europe_africa";
  }

  try {
    var guess = bandForOffset(new Date().getTimezoneOffset());
    if ([].some.call(select.options, function (o) { return o.value === guess; })) {
      select.value = guess;
    }
  } catch (e) {
    /* keep the server default */
  }

  function sync() {
    for (var i = 0; i < fields.length; i++) fields[i].value = select.value;
  }
  select.addEventListener("change", sync);
  sync();
})();
