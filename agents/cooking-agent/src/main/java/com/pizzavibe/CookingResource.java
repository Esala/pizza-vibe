package com.pizzavibe;

import com.pizzavibe.model.CookRequest;
import com.pizzavibe.model.CookingResult;
import com.pizzavibe.service.CookingService;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/cook")
public class CookingResource {

    @Inject
    CookingService cookingService;

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String hello() {
        return "Hello from Cooking Agent";
    }

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public CookingResult cookPizzas(CookRequest request) {
        return cookingService.cookPizzas(request.pizzas());
    }
}
