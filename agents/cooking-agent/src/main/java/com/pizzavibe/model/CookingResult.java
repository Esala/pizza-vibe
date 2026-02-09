package com.pizzavibe.model;

import java.util.List;

public record CookingResult(
    List<String> cookedPizzas,
    List<String> failedPizzas,
    String message
) {
    public static CookingResult success(List<String> cookedPizzas) {
        return new CookingResult(cookedPizzas, List.of(), "Successfully cooked " + cookedPizzas.size() + " pizza(s)");
    }

    public static CookingResult partial(List<String> cookedPizzas, List<String> failedPizzas) {
        return new CookingResult(cookedPizzas, failedPizzas,
            "Cooked " + cookedPizzas.size() + " pizza(s), " + failedPizzas.size() + " failed due to insufficient ingredients");
    }

    public static CookingResult failure(List<String> failedPizzas, String reason) {
        return new CookingResult(List.of(), failedPizzas, reason);
    }
}
