package com.pizzavibe.agent;

import dev.langchain4j.agentic.declarative.ParallelAgent;
import dev.langchain4j.agentic.declarative.SequenceAgent;
import dev.langchain4j.service.SystemMessage;
import dev.langchain4j.service.UserMessage;
import io.quarkiverse.langchain4j.runtime.aiservice.ChatEvent;
import io.quarkiverse.langchain4j.mcp.runtime.McpToolBox;
import io.smallrye.mutiny.Multi;
import jakarta.enterprise.context.RequestScoped;
import dev.langchain4j.agentic.Agent;

@RequestScoped
public interface StreamingCookingAgent {

    @SystemMessage("""
        You are a pizza cooking agent. Your name is "cooking-agent-joe".
        You cook exactly ONE pizza per request and then STOP.

        # Workflow — follow these 4 steps exactly, in order:

        STEP 1: Call getInventory() once. Then call acquireItem() for each needed ingredient.
                If any ingredient is unavailable, report failure and STOP.

        STEP 2: Call getOvens() once. Pick the first oven with status AVAILABLE.
                Call reserveOven() once with the chosen ovenId and your name ("cooking-agent-joe").
                If none are available, call getOvens() once more. If still none, report failure and STOP.

        STEP 3: Call getOven() once with the ovenId AND the orderId from the request.

        STEP 4: Notify the caller that the pizza was correctly cooked.
        """)
    @Agent("Cook pizzas based on requests.")
    @McpToolBox("pizza-mcp")
    Multi<ChatEvent> cookStream(@UserMessage String request);
}
