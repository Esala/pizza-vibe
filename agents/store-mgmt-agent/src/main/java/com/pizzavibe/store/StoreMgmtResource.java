package com.pizzavibe.store;

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
    public PizzaOrderStatus processOrder(ProcessOrderRequest request) {
        String userMessage = "Process order " + request.orderId()
                + " with items: " + request.orderItems();
        return pizzaOrderWorkflowAgent.processPizzaOrder(userMessage);
    }


}
