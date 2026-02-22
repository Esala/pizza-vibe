package com.pizzavibe.cooking;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.pizzavibe.cooking.agent.CookingAgent;
import com.pizzavibe.cooking.model.CookRequest;
import com.pizzavibe.cooking.model.Pizza;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

import java.util.Arrays;

@Path("/cook")
public class CookingResource {

  @Inject
  CookingAgent cookingAgent;

  @Inject
  ObjectMapper objectMapper;

  @GET
  @Produces(MediaType.TEXT_PLAIN)
  public String hello() {
    return "Hello from Cooking Agent";
  }


  @POST
  @Consumes(MediaType.APPLICATION_JSON)
  public String cookPizza(CookRequest request) {
    return cookingAgent.cook(request.orderId(), Arrays.toString(request.orderItems().toArray()));
  }

}
