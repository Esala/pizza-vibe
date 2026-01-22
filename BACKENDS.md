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
