package com.pizzavibe.store.agent;


import dev.langchain4j.agentic.declarative.A2AClientAgent;

public interface DeliveryRemoteAgent {
  @A2AClientAgent(a2aServerUrl = "http://localhost:8089",
      outputKey = "deliveryReport",
      description = "Agent that delivers an order.")
  String deliverOrder(String request);

}
