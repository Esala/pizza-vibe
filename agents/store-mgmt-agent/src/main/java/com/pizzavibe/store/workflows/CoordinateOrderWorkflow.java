package com.pizzavibe.store.workflows;

import com.pizzavibe.store.agent.CookingRemoteAgent;
import com.pizzavibe.store.agent.DrinksAgent;
import dev.langchain4j.agentic.declarative.ParallelAgent;

public interface CoordinateOrderWorkflow {

  @ParallelAgent(outputKey = "orderKitchenResult",
      subAgents = { CookingRemoteAgent.class, DrinksAgent.class })
  String coordinateKitchenAndDrinks(String request);
}
