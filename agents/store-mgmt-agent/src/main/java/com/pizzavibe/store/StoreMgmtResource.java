package com.pizzavibe.store;

import com.pizzavibe.store.model.DrinkItem;
import com.pizzavibe.store.model.OrderItem;
import com.pizzavibe.store.model.PizzaOrderStatus;
import com.pizzavibe.store.model.ProcessOrderRequest;
import com.pizzavibe.store.workflows.PizzaOrderWorkflow;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

import java.util.Arrays;
import java.util.List;

@Path("/mgmt")
public class StoreMgmtResource {

    @Inject
    PizzaOrderWorkflow pizzaOrderWorkflowAgent;

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String hello() {
        return "Hello from Store Management Agent";
    }

    @POST
    @Path("/processOrder")
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public PizzaOrderStatus processOrder(ProcessOrderRequest request) {
      System.out.println(request.toString());
        return pizzaOrderWorkflowAgent.processPizzaOrder(request.orderId(),
            Arrays.toString(request.orderItems().toArray()),
            Arrays.toString(request.drinkItems().toArray()));
    }


}
