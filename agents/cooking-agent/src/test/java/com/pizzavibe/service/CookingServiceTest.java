package com.pizzavibe.service;

import com.pizzavibe.model.Ingredient;
import com.pizzavibe.model.Pizza;
import com.pizzavibe.model.CookingResult;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

@QuarkusTest
class CookingServiceTest {

    @Inject
    CookingService cookingService;

    @Inject
    InventoryService inventoryService;

    @BeforeEach
    void setUp() {
        inventoryService.resetInventory();
    }

    @Test
    void shouldCookPizzaWhenIngredientsAvailable() {
        CookingResult result = cookingService.cookPizzas(List.of("Margherita"));

        assertNotNull(result);
        assertFalse(result.cookedPizzas().isEmpty());
        assertEquals(1, result.cookedPizzas().size());
        assertEquals("Margherita", result.cookedPizzas().get(0));
    }

    @Test
    void shouldReturnEmptyListWhenIngredientsNotAvailable() {
        List<String> manyPizzas = List.of(
            "Margherita", "Margherita", "Margherita", "Margherita", "Margherita",
            "Margherita", "Margherita", "Margherita", "Margherita", "Margherita",
            "Margherita", "Margherita", "Margherita", "Margherita", "Margherita",
            "Margherita", "Margherita", "Margherita", "Margherita", "Margherita",
            "Margherita"
        );

        CookingResult result = cookingService.cookPizzas(manyPizzas);

        assertNotNull(result);
        assertTrue(result.cookedPizzas().size() < manyPizzas.size());
        assertFalse(result.failedPizzas().isEmpty());
    }

    @Test
    void shouldConsumeIngredientsAfterCooking() {
        int initialDough = inventoryService.getQuantity(Ingredient.DOUGH);

        cookingService.cookPizzas(List.of("Margherita"));

        assertTrue(inventoryService.getQuantity(Ingredient.DOUGH) < initialDough);
    }

    @Test
    void shouldCookMultiplePizzasOfDifferentTypes() {
        CookingResult result = cookingService.cookPizzas(List.of("Margherita", "Pepperoni"));

        assertNotNull(result);
        assertEquals(2, result.cookedPizzas().size());
        assertTrue(result.cookedPizzas().contains("Margherita"));
        assertTrue(result.cookedPizzas().contains("Pepperoni"));
    }

    @Test
    void shouldReportFailedPizzasWhenSomeCannotBeMade() {
        CookingResult result = cookingService.cookPizzas(List.of("UnknownPizza"));

        assertNotNull(result);
        assertTrue(result.cookedPizzas().isEmpty());
        assertFalse(result.failedPizzas().isEmpty());
        assertEquals("UnknownPizza", result.failedPizzas().get(0));
    }
}
