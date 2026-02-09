package com.pizzavibe.agent;

import com.pizzavibe.service.InventoryService;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

@QuarkusTest
class InventoryToolTest {

    @Inject
    InventoryTool inventoryTool;

    @Inject
    InventoryService inventoryService;

    @BeforeEach
    void setUp() {
        inventoryService.resetInventory();
    }

    @Test
    void shouldGetInventory() {
        String inventory = inventoryTool.getInventory();

        assertNotNull(inventory);
        assertTrue(inventory.contains("DOUGH"));
        assertTrue(inventory.contains("MOZZARELLA"));
        assertTrue(inventory.contains("TOMATO_SAUCE"));
    }

    @Test
    void shouldCheckIngredientAvailability() {
        assertTrue(inventoryTool.hasIngredient("DOUGH", 1));
        assertTrue(inventoryTool.hasIngredient("dough", 1));
        assertFalse(inventoryTool.hasIngredient("DOUGH", 1000));
        assertFalse(inventoryTool.hasIngredient("UNKNOWN_INGREDIENT", 1));
    }

    @Test
    void shouldCookPizzas() {
        String result = inventoryTool.cookPizzas(List.of("Margherita"));

        assertNotNull(result);
        assertTrue(result.contains("Margherita"));
        assertTrue(result.contains("Cooked"));
    }

    @Test
    void shouldReportFailedPizzas() {
        String result = inventoryTool.cookPizzas(List.of("UnknownPizza"));

        assertNotNull(result);
        assertTrue(result.contains("UnknownPizza"));
        assertTrue(result.contains("Failed"));
    }
}
