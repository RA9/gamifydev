function displayContent() {
    const page = document.querySelector('main');
    const currentURL = (window.location.href).split('#')[1];
    switch (currentURL) {
        case 'about':
            AboutPage(page);
            break;
        case 'contact':
            ContactPage(page);
            break;
        case 'test':
            TestPage(page);
            break;
        case 'progress':
            ProgressPage(page);
            break;
        case 'today':
            TodayPage(page);
            break;
        case 'journey':
            JourneyPage(page);
            break;
        case 'launch':
            LaunchWeekPage(page);
            break;
        default:
            // Returning learners land on their Daily Standup; TodayPage falls
            // back to the marketing home for first-time visitors.
            TodayPage(page);
    }
    updateActiveNav(currentURL);
    updateHud();
}

// Highlight the nav link matching the current route.
function updateActiveNav(route) {
    document.querySelectorAll('[data-route]').forEach((el) => {
        el.classList.toggle('active', el.dataset.route === route);
    });
}

// Populate the header HUD (streak + XP) from past quiz attempts. Reuses the
// helpers defined in progress.js (computeStreak, XP_PER_CORRECT). Stays hidden
// until the learner has at least one recorded attempt.
async function updateHud() {
    const hud = document.getElementById('gd-hud');
    if (!hud || typeof DB === 'undefined') return;
    try {
        const scores = await DB.scores.toArray();
        if (!scores.length) {
            hud.classList.add('hidden');
            hud.classList.remove('flex');
            return;
        }
        const totalCorrect = scores.reduce((a, s) => a + (s.numCorrect || 0), 0);
        const xp = totalCorrect * (typeof XP_PER_CORRECT !== 'undefined' ? XP_PER_CORRECT : 10);
        const streak = typeof getStreakInfo === 'function'
            ? (await getStreakInfo()).streak
            : (typeof computeStreak === 'function' ? computeStreak(scores.map((s) => s.created_at)) : 0);
        const streakEl = document.getElementById('gd-hud-streak');
        const xpEl = document.getElementById('gd-hud-xp');
        if (streakEl) streakEl.textContent = streak;
        if (xpEl) xpEl.textContent = xp;
        hud.classList.remove('hidden');
        hud.classList.add('flex');
    } catch (e) {
        // No data yet / DB not ready — leave the HUD hidden.
    }
}

// Maintain the streak (auto-spend a freeze for a missed day, grant milestone
// freezes) before the first render, then route.
(async () => {
    try {
        if (typeof maintainStreak === 'function') await maintainStreak();
    } catch (e) {
        /* DB not ready / no data — render anyway */
    }
    displayContent();
})();

window.addEventListener('hashchange', displayContent);
