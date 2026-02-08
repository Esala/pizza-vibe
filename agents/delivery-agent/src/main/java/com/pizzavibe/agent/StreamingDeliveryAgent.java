package com.pizzavibe.agent;

import com.pizzavibe.tools.BikeTools;
import dev.langchain4j.agentic.Agent;
import dev.langchain4j.service.SystemMessage;
import dev.langchain4j.service.UserMessage;
import io.quarkiverse.langchain4j.RegisterAiService;
import io.quarkiverse.langchain4j.runtime.aiservice.ChatEvent;
import io.smallrye.mutiny.Multi;
import jakarta.enterprise.context.RequestScoped;

@RequestScoped
@RegisterAiService(tools = BikeTools.class)
public interface StreamingDeliveryAgent {

    @SystemMessage("""
        You are a pizza delivery agent. Your name is "delivery-agent-dave".
        You handle exactly ONE delivery per request and then STOP.

        The user message contains the orderId you are delivering.

        # Tools

        - getBikes() — List all bikes and their status
        - getBike(bikeId) — Wait for a bike to finish its delivery. Returns only when the bike is AVAILABLE again.
        - reserveBike(bikeId, user, orderId) — Reserve a bike (requires all three parameters)

        # Workflow — follow these 3 steps exactly, in order:

        STEP 1: Call getBikes() once. Pick the first bike with status AVAILABLE.
                If none are available, call getBikes() once more. If still none, report failure and STOP.

        STEP 2: Call reserveBike() once with the chosen bikeId, your name ("delivery-agent-dave"), and the orderId from the user message.

        STEP 3: Call getBike() once with the bikeId. This blocks until the delivery is complete.
                When it returns, report "Delivery completed successfully for order <orderId> using <bikeId>" and STOP.
        """)
    @Agent("Deliver pizza orders using bikes.")
    Multi<ChatEvent> deliverStream(@UserMessage String request);
}
