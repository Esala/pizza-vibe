package com.pizzavibe.service;

import com.pizzavibe.model.CookingResult;
import com.pizzavibe.model.Ingredient;
import com.pizzavibe.model.Pizza;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@ApplicationScoped
public class CookingService {

    @Inject
    InventoryService inventoryService;

    public CookingResult cookPizzas(List<String> pizzaNames) {
        List<String> cookedPizzas = new ArrayList<>();
        List<String> failedPizzas = new ArrayList<>();

        for (String pizzaName : pizzaNames) {
            Optional<Pizza> pizzaType = findPizzaByName(pizzaName);

            if (pizzaType.isEmpty()) {
                failedPizzas.add(pizzaName);
                continue;
            }

            Pizza pizza = pizzaType.get();
            if (canCookPizza(pizza)) {
                cookPizza(pizza);
                cookedPizzas.add(pizzaName);
            } else {
                failedPizzas.add(pizzaName);
            }
        }

        if (failedPizzas.isEmpty()) {
            return CookingResult.success(cookedPizzas);
        } else if (cookedPizzas.isEmpty()) {
            return CookingResult.failure(failedPizzas, "Could not cook any pizzas due to insufficient ingredients or unknown pizza types");
        } else {
            return CookingResult.partial(cookedPizzas, failedPizzas);
        }
    }

    private Optional<Pizza> findPizzaByName(String name) {
        return Pizza.getAllPizzaTypes().stream()
            .filter(p -> p.name().equalsIgnoreCase(name))
            .findFirst();
    }

    private boolean canCookPizza(Pizza pizza) {
        for (Map.Entry<Ingredient, Integer> entry : pizza.requiredIngredients().entrySet()) {
            if (!inventoryService.hasIngredient(entry.getKey(), entry.getValue())) {
                return false;
            }
        }
        return true;
    }

    private void cookPizza(Pizza pizza) {
        for (Map.Entry<Ingredient, Integer> entry : pizza.requiredIngredients().entrySet()) {
            inventoryService.consumeIngredient(entry.getKey(), entry.getValue());
        }
    }
}
