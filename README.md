## Project Overview

This Go server is a lightweight test and diagnostic tool built using the Gin web framework.
It is designed to help developers test various HTTP behaviors, system responses, and failure scenarios in a controlled environment.

The server provides multiple endpoints that cover a wide range of real-world situations such as:

Reading and logging incoming request bodies and headers.

Simulating high CPU and memory load.

Testing delayed or long-running responses.

Triggering intentional server crashes and exits.

Handling graceful shutdowns when receiving termination signals (e.g., Ctrl+C, SIGTERM).

Each endpoint demonstrates a different server behavior, making this tool useful for:

Load and stress testing.

Client timeout and retry logic testing.

Debugging request/response handling.

Observing resource usage and stability under heavy or delayed conditions.

Experimenting with graceful shutdown and panic recovery.

The server runs on port 8080 by default and logs detailed information about incoming requests, including headers, query parameters, and body content.

It can be safely used in development or staging environments for experimenting with client-server interaction patterns, testing API resilience, or verifying monitoring and alerting setups.

---

## Api Spec
| **Endpoint**           | **Method** | **Description**                                           | **Request Body**           | **Query / Params**           | **Response (200)**              | **Error / Behavior**                                                            |
| ---------------------- | ---------- | --------------------------------------------------------- | -------------------------- | ---------------------------- | ------------------------------- | ------------------------------------------------------------------------------- |
| `/post`                | `POST`     | Reads and prints request body and headers.                | Any raw JSON or text.      | None                         | `got it`                        | `400` if body cannot be read. Logs headers and body to console.                 |
| `/post-heavy`          | `POST`     | Simulates CPU and memory load by allocating 10 MB blocks. | `{ "block_count": <int> }` | None                         | `got it`                        | `400` if invalid JSON or `block_count ≤ 0`. High memory and CPU usage possible. |
| `/get`                 | `GET`      | Handles GET requests with query parameters.               | None                       | `query=<string>`             | `got it`                        | None — logs received query and headers.                                         |
| `/post-panic`          | `POST`     | Triggers a panic intentionally (nil pointer dereference). | None                       | None                         | *(no response, process panics)* | Server crashes due to unhandled panic.                                          |
| `/post-exit`           | `POST`     | Forces server to exit after responding.                   | None                       | None                         | `server shutting down...`       | Server terminates immediately via `os.Exit(1)`.                                 |
| `/post-long-time`      | `POST`     | Simulates a long-running request (30 seconds delay).      | None                       | None                         | `server Response after 30s`     | None — useful for timeout testing.                                              |
| `/post/delay/:seconds` | `POST`     | Delays response by given number of seconds.               | None                       | `:seconds` (integer in path) | `Response after {seconds}`      | `400` if invalid duration.                                                      |
| *(Signal Handler)*     | *(N/A)*    | Handles `SIGINT` and `SIGTERM` for graceful shutdown.     | None                       | None                         | *(Server shuts down cleanly)*   | Ensures current requests complete before stopping.                              |
---

## Additional Details

Server Port: 8080

Content-Type: Most endpoints accept application/json or raw text.

Logging: All requests log headers, query parameters, and/or body.

Graceful Shutdown: Handled automatically on Ctrl+C or OS termination signal.

---
## Example Commands
| **Scenario**           | **Command Example**                                                                                        |
| ---------------------- | ---------------------------------------------------------------------------------------------------------- |
| Basic POST             | `curl -X POST http://localhost:8080/post -d '{"hello":"world"}'`                                           |
| Heavy Load             | `curl -X POST http://localhost:8080/post-heavy -H "Content-Type: application/json" -d '{"block_count":3}'` |
| GET Query              | `curl "http://localhost:8080/get?query=test"`                                                              |
| Delayed Response (10s) | `curl -X POST http://localhost:8080/post/delay/10`                                                         |
| Long Response (30s)    | `curl -X POST http://localhost:8080/post-long-time`                                                        |
| Simulate Panic         | `curl -X POST http://localhost:8080/post-panic`                                                            |
| Force Exit             | `curl -X POST http://localhost:8080/post-exit`                                                             |
