## Store Service

/ralph-loop:ralph-loop "Implement store service in Go using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements:
1. Create a new directory called store and place all store service files inside
2. The store service should expose the following endpoints:
    1. POST endpoint /order to place a pizza order
    2. POST endpoint /events to receive updates from the kitchen and delivery services
    3. Create a websocket connection to the frontend application to send order updates
3. The order data model should include orderId(UUID), OrderItems, orderData and orderStatus
    1. OrderItems must container the pizzaType and the number of pizzas requested for that type
4. Use Go Chi for REST endpoints
5. Document all code and progress

Output <promise>DONE</promise> when all tests green." --max-iterations 25 --completion-promise "DONE"


## Kitchen Service

/ralph-loop:ralph-loop "Implement kitchen service in Go using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements:
1. Create a new directory called kitchen and place all kitchen service files inside
2. The kitchen service should expose the following endpoints:
    1. POST endpoint /cook to cook the OrderItems from the store order
    2. The payload should only be the orderId and orderItems
    3. For each orderitem it should take a random time from 1 to 10 seconds to cook the item. Each item cooked should be printed in the terminal with the amount that it took to be printed
3. Use Go Chi for REST endpoints
4. Document all code and progress
5. Create a docker-compose file to run all services of the application in the root directory
6. Add instruction on how to run all the services in the README.md file at the root directory

Output <promise>DONE</promise> when all tests green." --max-iterations 25 --completion-promise "DONE"


## Implement Kitchen Interaction Feature in Store

/ralph-loop:ralph-loop "Implement calling kitchen and exchange events using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements: 
- For the store service: 
  - Return the Order ID so it can be used to track the order in the frontend application
  - Call the kitchen service to cook the order when the order is placed passing the orderId and orderItems
  - Accept update and Done events from the kitchen service to track the order status. Keep track of events per orderId
  - When a done event is received, update the order status to COOKED
- For the kitchen service: 
  - Print the amount of time it took to cook each order item
  - Send update events to the store service every second while the order is cooking. Events are sent using HTTP to the the store service /events endpoint
  - Send a DONE event when the order is done cooking

Output <promise>DONE</promise> when all tests green." --max-iterations 20 --completion-promise "DONE"


## Implement Event and Order List in Management Page Feature in Store

/ralph-loop:ralph-loop "Implement list events and list orders in management page using TDD.

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

Output <promise>DONE</promise> when all tests green." --max-iterations 20 --completion-promise "DONE"


## Implement websockets event exchange between store and frontend

/ralph-loop:ralph-loop "Implement websocket events between store and frontend using TDD.

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

Output <promise>DONE</promise> when all tests green." --max-iterations 20 --completion-promise "DONE"