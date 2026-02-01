package com.pizzavibe;

import com.pizzavibe.service.InventoryService;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.hasItem;
import static org.hamcrest.Matchers.hasSize;

@QuarkusTest
class CookingResourceTest {

    @Inject
    InventoryService inventoryService;

    @BeforeEach
    void setUp() {
        inventoryService.resetInventory();
    }

    @Test
    void testHelloEndpoint() {
        given()
          .when().get("/cook")
          .then()
             .statusCode(200)
             .body(is("Hello from Cooking Agent"));
    }

    @Test
    void testCookPizzaEndpoint() {
        given()
            .contentType("application/json")
            .body("{\"pizzas\": [\"Margherita\"]}")
          .when().post("/cook")
          .then()
             .statusCode(200)
             .body("cookedPizzas", hasSize(1))
             .body("cookedPizzas", hasItem("Margherita"))
             .body("failedPizzas", hasSize(0));
    }

    @Test
    void testCookMultiplePizzasEndpoint() {
        given()
            .contentType("application/json")
            .body("{\"pizzas\": [\"Margherita\", \"Pepperoni\", \"Veggie\"]}")
          .when().post("/cook")
          .then()
             .statusCode(200)
             .body("cookedPizzas", hasSize(3))
             .body("failedPizzas", hasSize(0));
    }

    @Test
    void testCookUnknownPizzaEndpoint() {
        given()
            .contentType("application/json")
            .body("{\"pizzas\": [\"SuperSpecial\"]}")
          .when().post("/cook")
          .then()
             .statusCode(200)
             .body("cookedPizzas", hasSize(0))
             .body("failedPizzas", hasSize(1))
             .body("failedPizzas", hasItem("SuperSpecial"));
    }
}