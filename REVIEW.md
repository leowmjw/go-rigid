# REVIEW.md — TDD Test Case Analysis

## High-Level Objective
Implement Rama-in-Go (RIG) in Go 1.25, mirroring Rama REST APIs (`/rest/{module}/depot/*/append`, `/rest/{module}/pstate/$$/select[One]`, `/rest/{module}/query/{name}/invoke`) using stdlib HTTP + slog, and Temporal for workflows/activities. Enforce TDD via tutorials 1–6, anonymous function decoupling for side effects, deterministic workflows, and non-determinism in activities.

## Critical Gaps and Feedback (Checklist Ordered by Priority)

### 1. **Missing Assertions in Tests**
   - **Issue**: Most tests (T2, T4, T6) perform HTTP calls but lack assertions on response bodies or status codes. E.g., T2 unmarshals count but doesn't verify it (expected 3 for "to"); T4 and T6 don't check selectOne results.
   - **Impact**: Tests pass vacuously; no verification of Rama logic.
   - **Recommendation**: Add explicit assertions for expected values post-selectOne.

### 2. **Incomplete Temporal Integration**
   - **Issue**: No tests exercise Temporal workflows/activities. T5 uses `testutil.NewEnv()` but doesn't run workflows or activities.
   - **Impact**: Core Temporal usage (deterministic workflows, activity retries) untested.
   - **Recommendation**: Implement workflow tests using `testsuite.TestWorkflowEnvironment` for determinism.

### 3. **No Determinism and Non-Determinism Tests**
   - **Issue**: No tests verify workflow determinism or activity non-determinism.
   - **Impact**: Violates Rama principles.
   - **Recommendation**: Add tests replaying workflows with same inputs for identical outcomes; test activity retries/failures.

### 4. **Anonymous Function Decoupling Not Tested**
   - **Issue**: No tests demonstrate replacing side-effect functions (e.g., HTTP calls) with mocks.
   - **Impact**: Side-effect coupling risks non-determinism.
   - **Recommendation**: Add dependency injection tests for activities.

### 5. **Microbatch Exactly-Once Not Verified**
   - **Issue**: T5 Microbatch test is empty; no simulation of exactly-once processing.
   - **Impact**: ETL guarantees untested.
   - **Recommendation**: Add scenario with duplicate appends, verify no double-processing.

### 6. **Stream Retries and Error Handling Missing**
   - **Issue**: T5 Stream test is empty; no retry mode tests.
   - **Impact**: Fault tolerance unverified.
   - **Recommendation**: Add scenarios for activity failures with retry policies.

### 7. **Query Invocation Untested**
   - **Issue**: No tests for `/rest/{module}/query/{name}/invoke` endpoint.
   - **Impact**: Query functionality incomplete in TDD.
   - **Recommendation**: Add TDD suite for query invocation.

### 8. **Partitioning Edge Cases**
   - **Issue**: T3 only tests same-key partitioning; no tests for different keys, hash collisions, or custom partitioners.
   - **Impact**: Partitioning reliability unproven.
   - **Recommendation**: Add tests for key distribution and load balancing.

### 9. **No Integration Tests Across Modules**
   - **Issue**: Tests are siloed per tutorial; no end-to-end flows.
   - **Impact**: Inter-module interactions unverified.
   - **Recommendation**: Add cross-tutorial integration tests.

### 10. **Logging and Error Scenarios**
    - **Issue**: No tests for slog logging or error responses (e.g., invalid paths, depot not found).
    - **Impact**: Observability and robustness untested.
    - **Recommendation**: Add error injection tests.

## Additional Test Scenarios (Not Implemented)
- **T7_QueryInvoke**: Test invoking named queries with parameters, asserting computed results.
- **T8_WorkflowReplay**: Run workflow twice with same inputs, assert identical outputs despite activity mocks.
- **T9_ActivityRetries**: Simulate activity failure, verify retry logic and eventual success/failure.
- **T10_MicrobatchDeduplication**: Append duplicate data, confirm exactly-once via idempotent keys.

## Questions for Clarification
- Should tests include Temporal cluster setup or mock only?
- Any specific Rama tutorial implementations to mirror exactly?
- Priorities for phases 0–6 in TODO.md?
