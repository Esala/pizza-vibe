package com.pizzavibe.model;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public record DeliveryUpdate(
    String type,
    String action,
    String message,
    String toolName,
    String toolInput
) {
    public static DeliveryUpdate action(String action, String message) {
        return new DeliveryUpdate("action", action, message, null, null);
    }

    public static DeliveryUpdate toolExecution(String toolName, String toolInput) {
        String action = mapToolToAction(toolName);
        String message = generateToolMessage(toolName, toolInput);
        return new DeliveryUpdate("action", action, message, toolName, toolInput);
    }

    public static DeliveryUpdate partial(String message) {
        return new DeliveryUpdate("partial", null, message, null, null);
    }

    public static DeliveryUpdate result(String message) {
        return new DeliveryUpdate("result", "completed", message, null, null);
    }

    private static String mapToolToAction(String toolName) {
        if (toolName == null) {
            return "unknown";
        }
        return switch (toolName.toLowerCase()) {
            case "getbikes" -> "checking_bikes";
            case "getbike" -> "checking_bike_status";
            case "reservebike" -> "reserving_bike";
            default -> "processing";
        };
    }

    private static String extractBikeId(String toolInput) {
        if (toolInput == null) return "";
        // Try to extract bikeId from JSON input like {"bikeId":"bike-1","user":"dave","orderId":"123"}
        try {
            int idx = toolInput.indexOf("\"bikeId\"");
            if (idx >= 0) {
                int start = toolInput.indexOf("\"", idx + 8) + 1;
                int end = toolInput.indexOf("\"", start);
                if (start > 0 && end > start) {
                    return toolInput.substring(start, end);
                }
            }
        } catch (Exception ignored) {}
        return toolInput;
    }

    private static String generateToolMessage(String toolName, String toolInput) {
        if (toolName == null) {
            return "Processing...";
        }
        return switch (toolName.toLowerCase()) {
            case "getbikes" -> "Checking available bikes for delivery";
            case "getbike" -> "Checking bike status: " + extractBikeId(toolInput);
            case "reservebike" -> "Reserving bike for delivery: " + extractBikeId(toolInput);
            default -> "Processing: " + toolName;
        };
    }
}
