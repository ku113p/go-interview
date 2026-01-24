### Channel-Based Transport Plan

1. **Transport primitives** *(implemented: internal/common/transport)*
   - Create `internal/common/transport` with:
     - `Envelope` struct (fields: `TaskID`, `Ctx`, `Module`, `Action`, `Payload any`, `Metadata map[string]string`).
     - `Result` struct (fields: `TaskID`, `Module`, `Action`, `Value any`, `Err error`, `FinishedAt time.Time`).
     - `Dispatcher` struct holding `requests chan Envelope`, shared `results chan<- Result`, worker count, `sync.WaitGroup`, and stop signal.
     - Constructor `NewDispatcher(handler interface{}, queueSize int, workers int, results chan<- Result)` that spins up workers immediately.
     - Methods `Send(ctx, env)` (enqueue or return error on cancellation/overflow) and `Shutdown(ctx)` (close queue, wait for workers).

2. **Per-module worker pools**
   - Each domain module (life_area, criterion, history, interview, memory, message, user, etc.) owns its dedicated dispatcher wired to its command/query handlers.
   - Put worker and queue sizes in `pkg/config/worker_pools.go`.
   - Modules expose thin helpers to wrap commands into `Envelope`s and call the router.

3. **Shared results channel**
   - In bootstrap (`cmd/app/main.go`) create one buffered channel `results := make(chan transport.Result, cfg.Router.ResultsBuffer)`.
   - Pass it to every dispatcher so all results converge to a single sink.

4. **Router module**
   - Build `internal/router` containing:
     - `Router` struct with:
       - `dispatchers map[moduleAction]*transport.Dispatcher` (module/action → dispatcher lookup).
       - Reference to the shared results channel (receive-only).
       - `sync.Map` acting as in-memory task storage (`taskState`).
       - Goroutine `run()` that continuously reads results and updates task states.
     - API:
       - `Register(module, action string, dispatcher *transport.Dispatcher)`.
       - `Dispatch(ctx, module, action string, payload any, metadata map[string]string) (taskID string, error)` — generates a TaskID, creates a `taskState`, enqueues the `Envelope` into the target dispatcher.
       - `Wait(ctx, taskID string) (*transport.Result, error)` — blocks until the result arrives or the context is canceled.
       - `Shutdown(ctx)` — stops the reader and all dispatchers.
     - `taskState` struct: `result *transport.Result`, `done chan struct{}` for synchronous waiting.

5. **Task flow**
   - Entry layer (currently main/CLI) calls `router.Dispatch`, receives a `TaskID`, and optionally calls `router.Wait` to retrieve the result synchronously.
   - Module workers process tasks and publish `Result` into the shared channel.
   - Router updates each `taskState`, acting as temporary storage until the future pipeline storage lands.

6. **Bootstrap order**
   1. Load config (defaults are fine for now).
   2. Create the shared results channel.
   3. Initialize infra dependencies and handlers for each module.
   4. Build dispatchers and register them in the router (`module/action → dispatcher`).
   5. Expose `router.Dispatch` and `router.Wait` as the main API for submitting and reading tasks.
   6. On shutdown, call `router.Shutdown`.

7. **Current-stage limits**
   - Logging and automated tests are not required right now.
   - Cleanup/TTL for `taskState` can be added later; for now results remain in memory until explicit removal or process restart.
