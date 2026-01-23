---
active: true
iteration: 1
max_iterations: 20
completion_promise: "DONE"
started_at: "2026-01-23T14:25:37Z"
---

Implement list events and list orders in management page using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements:
- For the store service:
    - Make sure that there is an endpoint that returns all the orders (GET /orders)
    - Make sure that there is an endpoint that returns all events per order (GET /events)
- For front-end
    - in the management page consume both endpoints to list all the orders and their status
    - all the events per order, when the order is selected.

Output <promise>DONE</promise> when all tests green.
