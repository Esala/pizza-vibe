package com.pizzavibe.mcp.tool;

import com.pizzavibe.mcp.client.OvenClient;
import com.pizzavibe.mcp.model.Oven;
import io.quarkiverse.mcp.server.Tool;
import io.quarkiverse.mcp.server.ToolArg;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.rest.client.inject.RestClient;

import java.util.List;
import java.util.stream.Collectors;

@ApplicationScoped
public class OvenTool {

    @Inject
    @RestClient
    OvenClient ovenClient;

    @Tool(description = "Get all pizza ovens with their current status (AVAILABLE or RESERVED)")
    public String getOvens() {
        List<Oven> ovens = ovenClient.getAll();
        return ovens.stream()
            .map(o -> "Oven: " + o.id() + ", Status: " + o.status() + (o.user() != null ? ", User: " + o.user() : ""))
            .collect(Collectors.joining("\n"));
    }

    @Tool(description = "Wait for a pizza oven to finish cooking. Polls every 2 seconds and returns only when the oven is AVAILABLE again. Call this once after reserving an oven.")
    public String getOven(@ToolArg(description = "The oven ID (e.g., oven-1, oven-2, oven-3, oven-4)") String ovenId) {
        int maxAttempts = 30;
        for (int i = 0; i < maxAttempts; i++) {
            Oven oven = ovenClient.getById(ovenId);
            if (Oven.STATUS_AVAILABLE.equals(oven.status())) {
                return "Oven: " + oven.id() + ", Status: " + oven.status();
            }
            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                return "Error: polling interrupted for oven " + ovenId;
            }
        }
        return "Error: Oven " + ovenId + " did not become AVAILABLE after " + (maxAttempts * 2) + " seconds";
    }

    @Tool(description = "Reserve a pizza oven for cooking. The oven will be automatically released after 5-20 seconds. Returns the reserved oven or an error if already reserved.")
    public String reserveOven(
            @ToolArg(description = "The oven ID to reserve (e.g., oven-1)") String ovenId,
            @ToolArg(description = "The user/cook name reserving the oven") String user) {
        Oven oven = ovenClient.reserve(ovenId, user);
        return "Oven: " + oven.id() + ", Status: " + oven.status() + ", User: " + oven.user();
    }
}
