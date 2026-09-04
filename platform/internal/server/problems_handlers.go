package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/content"
	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// judgeWait is how long a submit request will sit waiting for a sandbox slot
// and a verdict before handing the submission to the background judge.
//
// Long enough that an ordinary problem is judged in the same request the
// visitor submitted it in, short enough that a queue behind a heavy run turns
// into "we're on it" rather than a spinning browser tab.
const judgeWait = 12 * time.Second

// --- who is solving ---------------------------------------------------------

// solver identifies the current visitor without creating anything.
//
// Reads never mint a guest session: /problems is a public page, and giving
// every crawler a database row would fill the table with visitors who will
// never submit anything.
func (s *Server) solver(r *http.Request) store.Solver {
	if u := auth.CurrentUser(r.Context()); u != nil {
		return store.Solver{UserID: u.ID}
	}
	c, err := r.Cookie(auth.GuestCookie)
	if err != nil || c.Value == "" {
		return store.Solver{}
	}
	id, err := s.st.GuestByToken(r.Context(), c.Value)
	if err != nil {
		// A cookie for a session that was pruned or already claimed. Treat the
		// visitor as new rather than as an error; a fresh one is minted when
		// they next submit.
		return store.Solver{}
	}
	return store.Solver{GuestID: id}
}

// solverForWrite is solver, but mints a guest session when there is nobody yet
// — the first time an anonymous visitor actually submits something.
func (s *Server) solverForWrite(w http.ResponseWriter, r *http.Request) (store.Solver, error) {
	if by := s.solver(r); by != (store.Solver{}) {
		return by, nil
	}
	token, err := auth.NewToken()
	if err != nil {
		return store.Solver{}, err
	}
	id, err := s.st.CreateGuestSession(r.Context(), token)
	if err != nil {
		return store.Solver{}, err
	}
	auth.SetGuestCookie(w, token, s.secure)
	return store.Solver{GuestID: id}, nil
}

// rateKey is the run limiter's key for a solver.
//
// Guest ids are negated so they cannot collide with user ids: the two are
// separate sequences, and without this, guest 7 would share a throttle with
// user 7.
func rateKey(by store.Solver) int64 {
	if by.UserID != 0 {
		return by.UserID
	}
	return -by.GuestID
}

// claimGuestWork moves anything the visitor solved before signing in onto their
// account. Called on the way through sign-in and registration alike — someone
// who practised as a guest and then logged into an account they already had
// should keep their work just as much as someone who registered.
//
// A failure here is logged and swallowed: it costs the learner some history,
// which is bad, but blocking sign-in over it would be worse.
func (s *Server) claimGuestWork(w http.ResponseWriter, r *http.Request, userID int64) {
	c, err := r.Cookie(auth.GuestCookie)
	if err != nil || c.Value == "" {
		return
	}
	guestID, err := s.st.GuestByToken(r.Context(), c.Value)
	if err != nil {
		// Already claimed, or long gone. Clear the spent cookie either way.
		auth.ClearGuestCookie(w, s.secure)
		return
	}
	if n, err := s.st.ClaimGuestWork(r.Context(), guestID, userID); err != nil {
		log.Printf("claim guest work for user %d: %v", userID, err)
		return
	} else if n > 0 {
		log.Printf("user %d claimed %d guest submission(s)", userID, n)
	}
	auth.ClearGuestCookie(w, s.secure)
}

// --- pages ------------------------------------------------------------------

func (s *Server) handleProblems(w http.ResponseWriter, r *http.Request) {
	by := s.solver(r)
	list, err := s.st.ListProblems(r.Context(), by)
	if err != nil {
		http.Error(w, "could not load the problems", http.StatusInternalServerError)
		return
	}
	solved := 0
	for _, p := range list {
		if p.Solved {
			solved++
		}
	}
	s.render(w, r, "problems.html", ViewData{
		Title: "Practice problems",
		Data: map[string]any{
			"problems": list,
			"solved":   solved,
			// A guest with work to lose is the only one who should be nudged to
			// sign up, and they should be told exactly what they'd be keeping.
			"guestSolved": by.GuestID != 0 && solved > 0,
			"judgeable":   s.exec != nil && s.exec.Enabled(),
		},
	})
}

func (s *Server) handleProblem(w http.ResponseWriter, r *http.Request) {
	p, err := s.st.GetProblemBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	starters, err := s.st.Starters(r.Context(), p.ID)
	if err != nil {
		http.Error(w, "could not load the problem", http.StatusInternalServerError)
		return
	}
	tests, err := s.st.ProblemTests(r.Context(), p.ID)
	if err != nil {
		http.Error(w, "could not load the problem", http.StatusInternalServerError)
		return
	}
	// Only the visible tests. The hidden ones are what stop a solution that
	// pattern-matches the examples from passing as an understanding of the
	// problem, and listing them here would defeat them entirely.
	var samples []store.Check
	for _, t := range tests {
		if !t.Hidden {
			samples = append(samples, t)
		}
	}

	by := s.solver(r)
	data := map[string]any{
		"problem":   p,
		"statement": content.Render(p.Statement),
		"samples":   samples,
		"languages": languageChoices(starters),
		"starters":  starters,
		"judgeable": s.exec != nil && s.exec.Enabled(),
	}
	if sub, err := s.st.LatestProblemSubmission(r.Context(), p.ID, by); err == nil {
		data["submission"] = sub
		// Rendered through the same partial the submit and poll requests swap
		// in, so a reloaded page and a live one can't disagree about a verdict.
		data["result"] = resultData(sub, "")
	} else if !errors.Is(err, store.ErrNotFound) {
		log.Printf("problem %s: latest submission: %v", p.Slug, err)
	}
	s.render(w, r, "problem.html", ViewData{Title: p.Title, Data: data})
}

// languageChoices orders the languages a problem accepts, so the picker doesn't
// reshuffle between page loads the way a map iteration would.
func languageChoices(starters map[string]string) []string {
	// Author-intent order: the languages the bank teaches in, first.
	preferred := []string{"python", "c", "shell"}
	var out []string
	for _, l := range preferred {
		if _, ok := starters[l]; ok {
			out = append(out, l)
		}
	}
	// Anything else the author added, in a stable order of its own.
	var extra []string
	for l := range starters {
		if !slices.Contains(out, l) {
			extra = append(extra, l)
		}
	}
	slices.Sort(extra)
	return append(out, extra...)
}

// --- submitting -------------------------------------------------------------

func (s *Server) handleProblemSubmit(w http.ResponseWriter, r *http.Request) {
	p, err := s.st.GetProblemBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	if s.exec == nil || !s.exec.Enabled() {
		s.problemResult(w, r, nil, "Code execution is switched off on this server, so nothing can be judged right now.")
		return
	}

	code := strings.TrimSpace(r.FormValue("code"))
	if code == "" {
		s.problemResult(w, r, nil, "Write some code first.")
		return
	}
	lang, err := s.acceptedLanguage(r.Context(), p.ID, r.FormValue("language"))
	if err != nil {
		s.problemResult(w, r, nil, "That language isn't accepted for this problem.")
		return
	}

	by, err := s.solverForWrite(w, r)
	if err != nil {
		http.Error(w, "could not start a session", http.StatusInternalServerError)
		return
	}
	if wait, ok := s.runlim.allow(rateKey(by)); !ok {
		s.problemResult(w, r, nil, "Easy — wait "+strconv.Itoa(int(wait.Seconds()+1))+"s before running again.")
		return
	}

	subID, err := s.st.CreateProblemSubmission(r.Context(), p.ID, by, lang, code)
	if err != nil {
		http.Error(w, "could not save your submission", http.StatusInternalServerError)
		return
	}

	// Judge it right now if the sandbox has room. The row is already saved, so
	// anything that doesn't finish here is simply left for judge:run — which is
	// also what happens if this process dies mid-judge.
	ctx, cancel := context.WithTimeout(r.Context(), judgeWait)
	defer cancel()
	if rel, ok := s.runlim.acquire(ctx); ok {
		_, err := jobs.Judge(ctx, s.st, s.exec, subID)
		rel()
		if err != nil {
			log.Printf("problem %s: inline judge of submission %d: %v", p.Slug, subID, err)
		}
	}

	sub, err := s.st.GetProblemSubmission(r.Context(), subID)
	if err != nil {
		http.Error(w, "could not read your submission back", http.StatusInternalServerError)
		return
	}
	s.problemResult(w, r, sub, "")
}

// handleProblemResult is what the page polls while a submission is still with
// the background judge.
func (s *Server) handleProblemResult(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	sub, err := s.st.GetProblemSubmission(r.Context(), id)
	if err != nil {
		s.notFound(w, r)
		return
	}
	// A submission is someone's code. Knowing its id is not permission to read
	// it, so the poll only answers the solver who wrote it.
	if !s.solver(r).Owns(sub) {
		http.Error(w, "404 — no such submission.", http.StatusNotFound)
		return
	}
	s.problemResult(w, r, sub, "")
}

// problemResult renders the verdict panel, the one fragment both submitting and
// polling swap into the page.
func (s *Server) problemResult(w http.ResponseWriter, r *http.Request, sub *store.ProblemSubmission, msg string) {
	s.renderPartial(w, r, "problem.html", "result", resultData(sub, msg))
}

// resultData is what the verdict partial reads. Kept out of the template so the
// page render and the htmx swap are fed by the same code.
func resultData(sub *store.ProblemSubmission, msg string) map[string]any {
	data := map[string]any{"message": msg}
	if sub != nil {
		data["sub"] = sub
		data["pending"] = sub.Verdict == store.VerdictPending
		data["headline"] = verdictHeadline(sub.Verdict)
		data["ok"] = sub.Verdict == store.VerdictAccepted
	}
	return data
}

// acceptedLanguage refuses a language the problem has no starter for.
//
// The picker only offers what the problem accepts, but the form is a POST body
// and anyone can put anything in one. Without this, a submission could name a
// language the author never wrote tests for and be judged against tests that
// assume a different one.
func (s *Server) acceptedLanguage(ctx context.Context, problemID int64, want string) (string, error) {
	starters, err := s.st.Starters(ctx, problemID)
	if err != nil {
		return "", err
	}
	if _, ok := starters[want]; ok {
		return want, nil
	}
	return "", errors.New("language not accepted")
}

// verdictHeadline is what the learner reads first.
func verdictHeadline(v string) string {
	switch v {
	case store.VerdictAccepted:
		return "Accepted"
	case store.VerdictWrongAnswer:
		return "Wrong answer"
	case store.VerdictTimeLimit:
		return "Too slow"
	case store.VerdictCompileError:
		return "Didn't compile"
	case store.VerdictRuntimeError:
		return "Crashed"
	case store.VerdictPending:
		return "Judging…"
	default:
		return "Couldn't judge this"
	}
}
