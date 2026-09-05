-- Bind each generated enrollment to the exact passing placement decision that
-- created it. The copied exemptions are the immutable scheduling snapshot;
-- later diagnostic results must not rewrite work already assigned to a cohort.

ALTER TABLE enrollments ADD COLUMN placement_result_id INTEGER REFERENCES placement_results(id) ON DELETE RESTRICT;
ALTER TABLE enrollments ADD COLUMN placement_exemptions TEXT NOT NULL DEFAULT '[]';
ALTER TABLE sanctions ADD COLUMN enrollment_id INTEGER REFERENCES enrollments(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_enrollments_placement_result
  ON enrollments(placement_result_id);

-- Best-effort backfill for enrollments created before this invariant existed.
-- Match a passing result for the same user and path from no later than placement.
UPDATE enrollments
SET placement_result_id = (
      SELECT pr.id
      FROM placement_results pr
      WHERE pr.user_id = enrollments.user_id
        AND pr.passed = 1
        AND pr.recommended_path_id = enrollments.path_id
        AND (enrollments.placed_at IS NULL OR datetime(pr.created_at) <= datetime(enrollments.placed_at))
      ORDER BY datetime(pr.created_at) DESC, pr.id DESC
      LIMIT 1
    )
WHERE path_id IS NOT NULL
  AND placement_result_id IS NULL;

UPDATE enrollments
SET placement_exemptions = COALESCE((
      SELECT pr.exemptions
      FROM placement_results pr
      WHERE pr.id = enrollments.placement_result_id
    ), '[]')
WHERE placement_result_id IS NOT NULL;

-- Some legacy/operator enrollments were not created from a result matching their
-- current path. Preserve the schedule overlay they had immediately before this
-- migration without claiming an exact result binding that cannot be proven.
UPDATE enrollments
SET placement_exemptions = COALESCE((
      SELECT pr.exemptions
      FROM placement_results pr
      WHERE pr.user_id = enrollments.user_id
      ORDER BY datetime(pr.created_at) DESC, pr.id DESC
      LIMIT 1
    ), '[]')
WHERE placement_result_id IS NULL;

-- Bind historical attendance drops to the enrollment they most likely ended so
-- an appeal cannot accidentally restore a later enrollment in the same cohort.
UPDATE sanctions
SET enrollment_id = (
      SELECT e.id
      FROM enrollments e
      WHERE e.user_id = sanctions.user_id
        AND e.state = 'dropped'
        AND e.ended_at IS NOT NULL
        AND (sanctions.cohort_id IS NULL OR e.cohort_id = sanctions.cohort_id)
      ORDER BY ABS(strftime('%s', e.ended_at) - strftime('%s', sanctions.applied_at)),
               e.id DESC
      LIMIT 1
    )
WHERE kind = 'drop' AND shadow = 0 AND enrollment_id IS NULL;
