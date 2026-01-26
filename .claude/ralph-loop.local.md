---
active: true
iteration: 2
max_iterations: 20
completion_promise: "DONE"
started_at: "2026-01-26T10:11:56Z"
---

Implement websocket events between store and frontend using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements:
- For the store service:
    - Everytime that an event is received to the /events endpoint, send it to the frontend using websockets
    - Create the websocket event format as a go type and send it using websockets
    - Create a websocket connection using the client id
- For front-end
    - In the main page, when the order is placed by the user subscribe to the websocket events and display the incoming websocket events related to the order in the UI
    - When the order is placed display the returning order id
    - Add a websocket connection indicator in the UI
    - When connecting to the websocket create a unique client id, this client ide will be used to subscribe to the websocket events
    - Use a table format to display the events related to the order

Output <promise>DONE</promise> when all tests green.
