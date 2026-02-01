package com.pizzavibe.service;

import com.pizzavibe.model.Ingredient;
import jakarta.annotation.PostConstruct;
import jakarta.enterprise.context.ApplicationScoped;

import java.util.EnumMap;
import java.util.Map;

@ApplicationScoped
public class InventoryService {

    private final Map<Ingredient, Integer> inventory = new EnumMap<>(Ingredient.class);

    @PostConstruct
    void init() {
        resetInventory();
    }

    public void resetInventory() {
        inventory.clear();
        inventory.put(Ingredient.DOUGH, 20);
        inventory.put(Ingredient.TOMATO_SAUCE, 15);
        inventory.put(Ingredient.MOZZARELLA, 25);
        inventory.put(Ingredient.PEPPERONI, 10);
        inventory.put(Ingredient.MUSHROOMS, 12);
        inventory.put(Ingredient.OLIVES, 8);
        inventory.put(Ingredient.BELL_PEPPER, 10);
        inventory.put(Ingredient.ONION, 10);
        inventory.put(Ingredient.HAM, 8);
        inventory.put(Ingredient.PINEAPPLE, 6);
        inventory.put(Ingredient.BACON, 10);
        inventory.put(Ingredient.BASIL, 15);
    }

    public Map<Ingredient, Integer> getInventory() {
        return new EnumMap<>(inventory);
    }

    public int getQuantity(Ingredient ingredient) {
        return inventory.getOrDefault(ingredient, 0);
    }

    public boolean hasIngredient(Ingredient ingredient, int quantity) {
        return getQuantity(ingredient) >= quantity;
    }

    public void consumeIngredient(Ingredient ingredient, int quantity) {
        int current = getQuantity(ingredient);
        if (current >= quantity) {
            inventory.put(ingredient, current - quantity);
        } else {
            throw new IllegalStateException("Not enough " + ingredient + " in inventory");
        }
    }
}
