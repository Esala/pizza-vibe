package com.pizzavibe.agent;

import com.pizzavibe.model.CookingResult;
import com.pizzavibe.model.Ingredient;
import com.pizzavibe.service.CookingService;
import com.pizzavibe.service.InventoryService;
import dev.langchain4j.agent.tool.Tool;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@ApplicationScoped
public class InventoryTool {

    @Inject
    InventoryService inventoryService;

    @Inject
    CookingService cookingService;

    @Tool("Get the current inventory of ingredients with their quantities")
    public String getInventory() {
        Map<Ingredient, Integer> inventory = inventoryService.getInventory();
        return inventory.entrySet().stream()
            .map(e -> e.getKey().name() + ": " + e.getValue())
            .collect(Collectors.joining(", "));
    }

    @Tool("Check if a specific ingredient is available in the required quantity")
    public boolean hasIngredient(String ingredientName, int quantity) {
        try {
            Ingredient ingredient = Ingredient.valueOf(ingredientName.toUpperCase());
            return inventoryService.hasIngredient(ingredient, quantity);
        } catch (IllegalArgumentException e) {
            return false;
        }
    }

    @Tool("Cook the specified pizzas. Returns a result with cooked and failed pizzas.")
    public String cookPizzas(List<String> pizzaNames) {
        CookingResult result = cookingService.cookPizzas(pizzaNames);
        return result.message() + ". Cooked: " + result.cookedPizzas() + ". Failed: " + result.failedPizzas();
    }
}
