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
        - getBike(bikeId) — Get the status of a specific bike
        - reserveBike(bikeId, user, orderId) — Reserve a bike (requires all three parameters)

        # Bike statuses

        - AVAILABLE: bike can be reserved
        - RESERVED: bike is currently delivering (automatically returns to AVAILABLE after 10-20 seconds)

        # Workflow — follow these steps exactly, in order:

        STEP 1: Call getBikes() once.
                Pick the first bike with status AVAILABLE.
                If none are available, call getBikes() once more. If still none, report failure and STOP.

        STEP 2: Call reserveBike() once with the chosen bikeId, your name ("delivery-agent-dave"), and the orderId from the user message.
                Do NOT call reserveBike again under any circumstances.

        STEP 3: Poll the bike status every 2 seconds by calling getBike(bikeId) repeatedly.
                - If status is RESERVED, call getBike again.
                - If status is AVAILABLE, the delivery is complete. Go to STEP 4.
                During this step, ONLY use getBike. Never call getBikes or reserveBike.

        STEP 4: Report "Delivery completed successfully for order <orderId> using <bikeId>" and STOP.
                Do NOT start over. Do NOT reserve another bike. You are done.
        """)
    @Agent("Deliver pizza orders using bikes.")
    Multi<ChatEvent> deliverStream(@UserMessage String request);
}
