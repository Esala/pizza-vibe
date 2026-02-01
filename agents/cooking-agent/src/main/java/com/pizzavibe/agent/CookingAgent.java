package com.pizzavibe.agent;

import dev.langchain4j.service.SystemMessage;
import dev.langchain4j.service.UserMessage;
import io.quarkiverse.langchain4j.RegisterAiService;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
@RegisterAiService(tools = InventoryTool.class)
public interface CookingAgent {

    @SystemMessage("""
        You are a pizza cooking agent. Your job is to cook pizzas based on the available ingredients in the inventory.
        You have access to tools to check the inventory and cook pizzas.
        When asked to cook pizzas, first check if you have enough ingredients, then cook the pizzas that can be made.
        Report back which pizzas were successfully cooked and which ones failed due to insufficient ingredients.
        """)
    String cook(@UserMessage String request);
}
