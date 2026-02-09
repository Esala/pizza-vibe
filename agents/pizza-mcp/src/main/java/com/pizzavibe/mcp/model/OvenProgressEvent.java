package com.pizzavibe.mcp.model;

public record OvenProgressEvent(String orderId, String ovenId, int progress, String status) {
}
