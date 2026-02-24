# TODO

## Specification Tasks

## TASK-001
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Update `semantic-validation.md` to clarify that fragment input validation happens during semantic validation phase, not during resolution. See ISSUES.md Issue #2 for context and decision rationale.

## TASK-002
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Update `fragment-input-validation.md` to clarify it runs during semantic validation phase, not as separate step during resolution. See ISSUES.md Issue #2 for context and decision rationale.

## TASK-003
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Update `error-handling.md` to document fail-fast validation pipeline: Schema errors → return, Semantic errors (including input validation) → return, Success. See ISSUES.md Issue #4 for context and decision rationale.

## TASK-004
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Update `core-types.md` to replace `Body interface{}` with explicit union type. See ISSUES.md Issue #3 and `docs/body-type-design.md` for design details and rationale.

## TASK-005
- Priority: 1
- Status: TODO
- Dependencies: [TASK-004]
- Description: Update `fragment-resolution.md` to reflect new Body type structure with String/Array fields instead of interface{}. See ISSUES.md Issue #3 for context.

## TASK-006
- Priority: 2
- Status: TODO
- Dependencies: [TASK-001, TASK-002, TASK-003, TASK-004, TASK-005]
- Description: Review all specs for consistency and completeness after changes. Verify all cross-references are updated and design decisions from ISSUES.md are properly reflected.

## TASK-007
- Priority: 3
- Status: TODO
- Dependencies: [TASK-006]
- Description: Remove working documents: `docs/body-type-design.md` and `ISSUES.md` after all specification tasks complete. These are temporary design documents per ISSUES.md decision log.
