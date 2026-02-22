package com.pizzavibe.delivery;

import com.pizzavibe.delivery.agent.DeliveryAgent;
import com.pizzavibe.delivery.model.DeliveryRequest;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/deliver")
public class DeliveryResource {

    @Inject
    DeliveryAgent deliveryAgent;

    @GET
    @Produces(MediaType.TEXT_PLAIN)
    public String hello() {
        return "Hello from Delivery Agent";
    }

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    public String deliverOrderStream(DeliveryRequest request) {
        return deliveryAgent.deliverOrder(request.orderId());
    }

}
