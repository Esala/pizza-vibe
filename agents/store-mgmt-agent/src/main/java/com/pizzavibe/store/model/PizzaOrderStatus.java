package com.pizzavibe.store.model;

public record PizzaOrderStatus(OrderFinalStatus status, String kitchenReport, String deliveryReport) {
}

