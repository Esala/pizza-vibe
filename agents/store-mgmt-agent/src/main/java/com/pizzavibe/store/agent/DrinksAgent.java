package com.pizzavibe.store.agent;

import dev.langchain4j.agentic.Agent;
import dev.langchain4j.service.SystemMessage;
import dev.langchain4j.service.UserMessage;
import io.quarkiverse.langchain4j.mcp.runtime.McpToolBox;
import jakarta.enterprise.context.RequestScoped;

@RequestScoped
public interface DrinksAgent {
  @SystemMessage("""
        You are an agent in charge of preparing drinks to be delivered.
        
        Get all the drinks from the request and fetch them from the inventory.
        
        Call getInventory() once. Then call acquireItem() for each drink.
                If a drink is unavailable, report failure and STOP.
        """)
  @Agent("Prepare drinks for delivery.")
  @McpToolBox("pizza-mcp")
  String prepareDrinksForDelivery(@UserMessage String request);
}
