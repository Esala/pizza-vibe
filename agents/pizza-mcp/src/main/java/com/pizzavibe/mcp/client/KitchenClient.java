package com.pizzavibe.mcp.client;

import com.pizzavibe.mcp.model.OvenProgressEvent;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.core.MediaType;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

@RegisterRestClient(configKey = "kitchen-api")
@Path("/oven-progress")
public interface KitchenClient {

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    void sendProgress(OvenProgressEvent event);
}
