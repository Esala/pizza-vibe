---
active: true
iteration: 5
max_iterations: 25
completion_promise: "DONE"
started_at: "2026-01-22T10:44:59Z"
---

Implement Frontend using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements: 
1. Create a new directory called front-end
2. Create a front-end project using Next.js
3. Configure the front-end project with the default configuration but without Tailwind, only base CSS modules.
4. Clean all styles and pages to start with a blank page.
5. Create the main page connected to the store service that allow me to place an order. To place an order send a post request to the store service /order.
6. Configure Next.js to proxy requests to the store service and configure crossOrigin.  
7. Use Jest for testing

Output <promise>DONE</promise> when all tests green.
