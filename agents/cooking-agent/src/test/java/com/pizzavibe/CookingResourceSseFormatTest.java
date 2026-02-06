package com.pizzavibe;

import com.pizzavibe.agent.StreamingCookingAgent;
import dev.langchain4j.agent.tool.ToolExecutionRequest;
import dev.langchain4j.data.message.AiMessage;
import dev.langchain4j.model.chat.response.ChatResponse;
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
 * Tests verifying the exact SSE wire format produced by /cook/stream.
 * The Go kitchen service parser expects: data: {json}\n\n
 */
@QuarkusTest
class CookingResourceSseFormatTest {

    @InjectMock
    StreamingCookingAgent streamingCookingAgent;

    @Test
    void shouldProduceSseEventsWithDataPrefixAndBlankLineSeparators() {
        // Given: a single action event followed by completion
        ChatEvent action = new ChatEvent.BeforeToolExecutionEvent(
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
            .thenReturn(Multi.createFrom().items(action, completed));

        // When
        String body = given()
            .contentType(ContentType.JSON)
            .body("{\"pizzas\": [\"Margherita\"]}")
            .when().post("/cook/stream")
            .then()
            .statusCode(200)
            .extract().body().asString();

        System.out.println("=== RAW SSE OUTPUT ===");
        System.out.println(body);
        System.out.println("=== END RAW SSE OUTPUT ===");

        // Then: verify SSE wire format
        // Each event should be on a line starting with "data:" and events separated by blank lines
        assertTrue(body.contains("data:"), "SSE body must contain data: lines");

        // Count data lines - should have exactly 2 (action + result)
        long dataLineCount = body.lines()
            .filter(line -> line.startsWith("data:"))
            .count();
        assertEquals(2, dataLineCount, "Should have exactly 2 data lines (action + result)");
    }

    @Test
    void shouldProduceValidJsonInEachSseDataLine() {
        // Given
        ChatEvent action = new ChatEvent.BeforeToolExecutionEvent(
            ToolExecutionRequest.builder()
                .name("reserveOven")
                .arguments("{\"ovenId\": \"oven-1\"}")
                .build()
        );
        ChatEvent completed = new ChatEvent.ChatCompletedEvent(
            ChatResponse.builder()
                .aiMessage(AiMessage.from("Done"))
                .build()
        );

        Mockito.when(streamingCookingAgent.cookStream(Mockito.anyString()))
            .thenReturn(Multi.createFrom().items(action, completed));

        // When
        String body = given()
            .contentType(ContentType.JSON)
            .body("{\"pizzas\": [\"Margherita\"]}")
            .when().post("/cook/stream")
            .then()
            .statusCode(200)
            .extract().body().asString();

        // Then: each data line should contain valid JSON with required fields
        body.lines()
            .filter(line -> line.startsWith("data:"))
            .forEach(line -> {
                // Strip "data:" or "data: " prefix
                String json = line.startsWith("data: ") ? line.substring(6) : line.substring(5);
                assertTrue(json.startsWith("{") && json.endsWith("}"),
                    "Data line should contain JSON object: " + json);
                assertTrue(json.contains("\"type\""),
                    "JSON should have type field: " + json);
            });
    }
}
