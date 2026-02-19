package com.pizzavibe.store.workflows;

import com.pizzavibe.store.agent.CookingRemoteAgent;
import com.pizzavibe.store.agent.DeliveryRemoteAgent;
import com.pizzavibe.store.model.OrderFinalStatus;
import com.pizzavibe.store.model.PizzaOrderStatus;
import dev.langchain4j.agentic.declarative.Output;
import dev.langchain4j.agentic.declarative.SequenceAgent;
import dev.langchain4j.service.UserMessage;

public interface PizzaOrderWorkflow {

    @SequenceAgent(outputKey = "pizzaOrderAgentResult",
      subAgents = { CookingRemoteAgent.class, DeliveryRemoteAgent.class})
    PizzaOrderStatus processPizzaOrder(@UserMessage String request);

  @Output
  static PizzaOrderStatus output(String kitchenReport, String deliveryReport) {
    boolean kitchenFailed = false;
    boolean deliveryFailed = false;
    OrderFinalStatus status = OrderFinalStatus.SUCCESS;
    if(kitchenReport == null || kitchenReport.contains("ERROR") || kitchenReport.contains("FAILED")) {
      kitchenFailed = true;
    }
    if(deliveryReport == null || deliveryReport.contains("ERROR") || deliveryReport.contains("FAILED")) {
      deliveryFailed = true;
    }
    if(kitchenFailed || deliveryFailed) {
      status = OrderFinalStatus.FAILED;
    }
    return new PizzaOrderStatus(status, kitchenReport, deliveryReport);
  }
}
