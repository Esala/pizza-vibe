# Figma Page Designs — PizzaVibe App

**File Key**: `ODWoi5IypS0nsrCS99CowJ`
**Base URL**: `https://www.figma.com/design/ODWoi5IypS0nsrCS99CowJ/PizzaVibe-App`

## Pages

### Orders (page: `0:1`)
- Page URL: https://www.figma.com/design/ODWoi5IypS0nsrCS99CowJ/PizzaVibe-App?node-id=0-1
- Contains all Store page artboards

#### Store / New Order / Empty Cart (artboard: `1:2`)
- URL: https://www.figma.com/design/ODWoi5IypS0nsrCS99CowJ/PizzaVibe-App?node-id=1-2
- Layout: Header → Tabs → Main Content (Order Section | divider | Cart Section) → Footer
- Cart state: empty — disabled "PLACE ORDER" button + EmptyBlock
- Key sections:
  - Header (`50:575`): Logo + Navigation (Store active)
  - Tabs (`58:692`): "NEW ORDER" (active) + "YOUR ORDERS"
  - Main Content (`12:60`): padding `--space-padding` vertical
    - Content Wrapper (`12:63`): flex row, gap `--space-gap`
    - Order Section (`49:537`): flex column, gap 48px
      - Order Info (`12:69`): H2 "Order Here" + Body Default description, gap `--space-spacing-m`
      - Pizza Item Container (`49:532`): flex wrap, gap `--space-gap`, centered
        - 4x PizzaItem: Margherita $10, Pepperoni $15, Hawaiian $15, Vegan $12
    - Vertical divider (`12:72`): line between sections
    - Cart Section (`12:70`): flex column, gap `--space-gap`, stretch height
      - Cart Wrapper (`13:148`): flex row, gap 48px
        - Cart Text Container (`12:128`): H2 "Cart" + Body Default "Your cart is empty", gap `--space-spacing-s`
        - Button: disabled state, text "Place Order"
      - EmptyBlock: flex-1, full width
  - Footer (`50:637`)

#### Store / New Order / Cart with Items (artboard: `65:1441`)
- URL: https://www.figma.com/design/ODWoi5IypS0nsrCS99CowJ/PizzaVibe-App?node-id=65-1441
- Layout: same as Empty Cart, but cart has items
- Cart state: populated — active "PLACE ORDER - $35" button + CartItem list
- Differences from Empty Cart:
  - Cart description: "3 pizzas in the cart"
  - Button: active state, text "Place Order - $35"
  - Cart Items (`65:1718`): flex column, replaces EmptyBlock
    - CartItem: Margherita, $20, QuantitySelector (qty 2, default variant)
    - CartItem: Pepperoni, $15, QuantitySelector (qty 1, delete variant)
