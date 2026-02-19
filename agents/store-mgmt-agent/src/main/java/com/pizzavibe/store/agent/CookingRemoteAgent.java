package com.pizzavibe.store.agent;

import dev.langchain4j.agentic.declarative.A2AClientAgent;

public interface CookingRemoteAgent {

  @A2AClientAgent(a2aServerUrl = "http://localhost:8087",
      outputKey = "kitchenReport",
      description = "Agent that coordinate the cooking of an order.")
  String cook(String request);
}
