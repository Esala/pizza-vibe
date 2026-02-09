## Prompt to create Quarkus agent with Langchain4j

/ralph-loop:ralph-loop "Implement a new agent using Quarkus and Langchain4j using TDD.

Process:
1. Write failing test for next requirement
2. Implement minimal code to pass
3. Run tests
4. If failing, fix and retry
5. Refactor if needed
6. Repeat for all requirements

Requirements:
- Create a new Quarkus Agent using Langchain4j by following the documentation located in the following places:
  - https://quarkus.io/quarkus-workshop-langchain4j/
  - https://quarkus.io/guides/langchain4j
- The agent should be called cooking-agent inside the agents/ folder. 
- Use Maven to create the project.
- The agent should have the goal to cook pizzas based on ingredients available in the inventory.
- The agent should have an internal inventory of ingredients with mock data.
- If the agent has enough ingredients to cook pizzas, it should return a list of pizzas that were cooked.
- The agent should expose a REST endpoint to cook pizza orders.

Output <promise>DONE</promise> when all tests green." --max-iterations 20 --completion-promise "DONE"

