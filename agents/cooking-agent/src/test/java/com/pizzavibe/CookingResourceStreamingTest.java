package com.pizzavibe;

import com.pizzavibe.agent.StreamingCookingAgent;
import dev.langchain4j.agent.tool.ToolExecutionRequest;
import dev.langchain4j.data.message.AiMessage;
import dev.langchain4j.model.chat.response.ChatResponse;
import dev.langchain4j.service.tool.ToolExecution;
import io.quarkiverse.langchain4j.runtime.aiservice.ChatEvent;
import io.quarkus.test.InjectMock;
import io.quarkus.test.junit.QuarkusTest;
import io.restassured.http.ContentType;
import io.smallrye.mutiny.Multi;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static io.restassured.RestAssured.given;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests that the /cook/stream endpoint produces properly formatted SSE events
 * with cooking action updates (e.g., checking inventory, reserving oven).
 */
@QuarkusTest
class CookingResourceStreamingTest {

    @InjectMock
    StreamingCookingAgent streamingCookingAgent;

    @Test
    void shouldStreamToolExecutionEventsAsSse() {
        // Given: agent emits tool execution events for a typical cooking flow
        ChatEvent checkInventory = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("getInventory")
                .arguments("{}")
                .build()
        );
        ChatEvent acquireItem = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("acquireItem")
                .arguments("{\"itemName\": \"mozzarella\", \"quantity\": 2}")
                .build()
        );
        ChatEvent reserveOven = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("reserveOven")
                .arguments("{\"ovenId\": \"oven-1\", \"reservedBy\": \"cooking-agent-joe\"}")
                .build()
        );
        ChatEvent completed = new ChatEvent.ChatCompletedEvent(
            ChatResponse.builder()
                .aiMessage(AiMessage.from("Margherita pizza cooked successfully!"))
                .build()
        );

        Mockito.when(streamingCookingAgent.cookStream(Mockito.anyString()))
            .thenReturn(Multi.createFrom().items(checkInventory, acquireItem, reserveOven, completed));

        // When
        String body = given()
            .contentType(ContentType.JSON)
            .body("{\"pizzas\": [\"Margherita\"]}")
            .when().post("/cook/stream")
            .then()
            .statusCode(200)
            .contentType("text/event-stream")
            .extract().body().asString();

        // Then: SSE body should contain data lines with JSON CookingUpdate objects
        assertNotNull(body);
        assertFalse(body.isEmpty(), "SSE body should not be empty");

        // Verify each action event is present in the SSE stream
        assertTrue(body.contains("\"action\":\"checking_inventory\""),
            "Should contain checking_inventory action, got: " + body);
        assertTrue(body.contains("\"action\":\"acquiring_ingredients\""),
            "Should contain acquiring_ingredients action, got: " + body);
        assertTrue(body.contains("\"action\":\"reserving_oven\""),
            "Should contain reserving_oven action, got: " + body);
        assertTrue(body.contains("\"type\":\"result\""),
            "Should contain result event, got: " + body);
    }

    @Test
    void shouldStreamOvenPollingEventsAsSse() {
        // Given: agent polls oven status (checking_oven_status action)
        ChatEvent checkOvens = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("getOvens")
                .arguments("{}")
                .build()
        );
        ChatEvent checkOvenStatus = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("getOven")
                .arguments("oven-1")
                .build()
        );
        ChatEvent completed = new ChatEvent.ChatCompletedEvent(
            ChatResponse.builder()
                .aiMessage(AiMessage.from("Done"))
                .build()
        );

        Mockito.when(streamingCookingAgent.cookStream(Mockito.anyString()))
            .thenReturn(Multi.createFrom().items(checkOvens, checkOvenStatus, completed));

        // When
        String body = given()
            .contentType(ContentType.JSON)
            .body("{\"pizzas\": [\"Margherita\"]}")
            .when().post("/cook/stream")
            .then()
            .statusCode(200)
            .extract().body().asString();

        // Then
        assertTrue(body.contains("\"action\":\"checking_ovens\""),
            "Should contain checking_ovens action");
        assertTrue(body.contains("\"action\":\"checking_oven_status\""),
            "Should contain checking_oven_status action");
    }

    @Test
    void shouldIncludeToolNameAndInputInSseEvents() {
        // Given
        ChatEvent event = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("reserveOven")
                .arguments("{\"ovenId\": \"oven-2\"}")
                .build()
        );
        ChatEvent completed = new ChatEvent.ChatCompletedEvent(
            ChatResponse.builder()
                .aiMessage(AiMessage.from("Done"))
                .build()
        );

        Mockito.when(streamingCookingAgent.cookStream(Mockito.anyString()))
            .thenReturn(Multi.createFrom().items(event, completed));

        // When
        String body = given()
            .contentType(ContentType.JSON)
            .body("{\"pizzas\": [\"Margherita\"]}")
            .when().post("/cook/stream")
            .then()
            .statusCode(200)
            .extract().body().asString();

        // Then: SSE events should include toolName and toolInput fields
        assertTrue(body.contains("\"toolName\":\"reserveOven\""),
            "Should contain toolName field");
        assertTrue(body.contains("\"toolInput\":\"{\\\"ovenId\\\": \\\"oven-2\\\"}\""),
            "Should contain toolInput field");
    }

    @Test
    void shouldStreamSseWithDataPrefix() {
        // Given: the SSE format must use "data:" prefix for each event
        ChatEvent event = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("getInventory")
                .arguments("{}")
                .build()
        );
        ChatEvent completed = new ChatEvent.ChatCompletedEvent(
            ChatResponse.builder()
                .aiMessage(AiMessage.from("Done"))
                .build()
        );

        Mockito.when(streamingCookingAgent.cookStream(Mockito.anyString()))
            .thenReturn(Multi.createFrom().items(event, completed));

        // When
        String body = given()
            .contentType(ContentType.JSON)
            .body("{\"pizzas\": [\"Margherita\"]}")
            .when().post("/cook/stream")
            .then()
            .statusCode(200)
            .contentType("text/event-stream")
            .extract().body().asString();

        // Then: each SSE event line should start with "data:"
        String[] lines = body.split("\n");
        boolean foundDataLine = false;
        for (String line : lines) {
            if (line.startsWith("data:")) {
                foundDataLine = true;
                // Extract JSON after "data:" (with optional space)
                String json = line.startsWith("data: ") ? line.substring(6) : line.substring(5);
                assertTrue(json.startsWith("{"), "Data line should contain JSON: " + json);
                assertTrue(json.contains("\"type\":"), "JSON should have type field: " + json);
            }
        }
        assertTrue(foundDataLine, "SSE body should contain at least one data: line, got: " + body);
    }
}
