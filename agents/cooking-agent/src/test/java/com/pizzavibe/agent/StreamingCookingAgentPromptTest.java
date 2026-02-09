package com.pizzavibe.agent;

import dev.langchain4j.service.SystemMessage;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Method;

import static org.junit.jupiter.api.Assertions.*;

class StreamingCookingAgentPromptTest {

    @Test
    void agentNameShouldBeCookingAgentJoe() throws NoSuchMethodException {
        Method method = StreamingCookingAgent.class.getMethod("cookStream", String.class);
        SystemMessage annotation = method.getAnnotation(SystemMessage.class);
        assertNotNull(annotation, "cookStream should have @SystemMessage");

        String systemMessage = String.join("\n", annotation.value());
        assertTrue(systemMessage.contains("cooking-agent-joe"),
            "System message should reference agent name 'cooking-agent-joe'");
        assertFalse(systemMessage.contains("cooking-agent-joe-stream"),
            "System message should NOT reference old name 'cooking-agent-joe-stream'");
    }

    @Test
    void systemMessageShouldInstructToReserveOvenWithAgentName() throws NoSuchMethodException {
        Method method = StreamingCookingAgent.class.getMethod("cookStream", String.class);
        SystemMessage annotation = method.getAnnotation(SystemMessage.class);
        String systemMessage = String.join("\n", annotation.value());

        // The prompt should instruct the agent to use its name when reserving an oven
        assertTrue(systemMessage.contains("reserveOven"),
            "System message should mention reserveOven tool");
        // The name used for reservation should be cooking-agent-joe
        assertTrue(systemMessage.contains("\"cooking-agent-joe\""),
            "System message should instruct using name 'cooking-agent-joe' for oven reservation");
    }
}
