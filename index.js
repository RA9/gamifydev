function displayContent() {
    const page = document.querySelector('main');
    // Clear any lingering celebration/briefing overlay when navigating.
    document.querySelectorAll('.gd-overlay').forEach((o) => o.remove());
    const currentURL = (window.location.href).split('#')[1];
    // Parameterised route: #project/<id> opens the build-along player.
    if (currentURL && currentURL.indexOf('project/') === 0) {
        ProjectPlayerPage(page, currentURL.slice('project/'.length));
        updateActiveNav('projects');
        updateHud();
        return;
    }
    // Parameterised route: #code/<id> opens a coding challenge.
    if (currentURL && currentURL.indexOf('code/') === 0) {
        ChallengePage(page, currentURL.slice('code/'.length));
        updateActiveNav('code');
        updateHud();
        return;
    }
    // Parameterised route: #assess/<id> runs a final assessment.
    if (currentURL && currentURL.indexOf('assess/') === 0) {
        AssessmentPage(page, currentURL.slice('assess/'.length));
        updateActiveNav('certify');
        updateHud();
        return;
    }
    // Parameterised route: #certificate/<id> shows an earned certificate.
    if (currentURL && currentURL.indexOf('certificate/') === 0) {
        CertificatePage(page, currentURL.slice('certificate/'.length));
        updateActiveNav('certify');
        updateHud();
        return;
    }
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
        case 'worlds':
            WorldsPage(page);
            break;
        case 'projects':
            ProjectsPage(page);
            break;
        case 'code':
            CodeLabPage(page);
            break;
        case 'certify':
            CertifyPage(page);
            break;
        case 'review':
            ReviewPage(page);
            break;
        case 'terminal':
            TerminalPage(page);
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
    // Celebrate any newly-earned cross-world achievements (guards against
    // stacking on another overlay; silently baselines on first ever run).
    if (typeof maybeCelebrateAchievements === 'function') {
        setTimeout(() => maybeCelebrateAchievements(), 400);
    }
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
