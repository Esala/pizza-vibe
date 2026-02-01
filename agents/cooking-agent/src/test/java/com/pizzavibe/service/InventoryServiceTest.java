package com.pizzavibe.service;

import com.pizzavibe.model.Ingredient;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

@QuarkusTest
class InventoryServiceTest {

    @Inject
    InventoryService inventoryService;

    @Test
    void shouldHaveInitialInventoryWithMockData() {
        Map<Ingredient, Integer> inventory = inventoryService.getInventory();

        assertNotNull(inventory);
        assertFalse(inventory.isEmpty());
        assertTrue(inventory.containsKey(Ingredient.DOUGH));
        assertTrue(inventory.containsKey(Ingredient.TOMATO_SAUCE));
        assertTrue(inventory.containsKey(Ingredient.MOZZARELLA));
    }

    @Test
    void shouldCheckIfIngredientIsAvailable() {
        assertTrue(inventoryService.hasIngredient(Ingredient.DOUGH, 1));
    }

    @Test
    void shouldConsumeIngredients() {
        int initialDough = inventoryService.getQuantity(Ingredient.DOUGH);

        inventoryService.consumeIngredient(Ingredient.DOUGH, 1);

        assertEquals(initialDough - 1, inventoryService.getQuantity(Ingredient.DOUGH));
    }

    @Test
    void shouldReturnFalseWhenInsufficientIngredients() {
        assertFalse(inventoryService.hasIngredient(Ingredient.DOUGH, 1000));
    }
}
