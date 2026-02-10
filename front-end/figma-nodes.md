# Figma Node Tree — PizzaVibe UI Kit

**File Key**: `Iia6bIqfQwSvXxTnfedTXj`
**Base URL**: `https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit`



### Cover (node: `7:2`)
- URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=7-2
- Brand showcase page

## Design Tokens

### Typography (page: `0:1`, frame: `1:2`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=0-1
- Frame URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=1-2
- Tokens: H1 (Knewave), H2, H3, Body Default, Body Small (Geist)

### Color (page: `97:17`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=97-17
- **Colors frame** (`97:18`): Raw palette — Neutrals, Green, Red, Yellow, Blue, Pink (9 shades each)
  - URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=97-18
- **Palette frame** (`97:189`): Semantic token assignments (source of truth for token values)
  - URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=97-189
  - Backgrounds (`97:190`): default, subtle, primary, secondary, tertiary, inverted, inverted/subtle
  - Border (`98:13`): default, subtle
  - Text (`98:3`): default, subtle, primary, secondary, tertiary, inverted

### Sizing (page: `98:22`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=98-22
- Parent page containing: Padding, Margin, Gap, Sizing Scale, Border Widths, Corners
- Main frame: `98:23` ("Space and Size")

### Border Widths (node: `98:65`)
- URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=98-65
- Tokens: thin, default, thick, thicker

### Spacing (nodes: `98:24`, `98:44`, `98:51`)
- URL (padding): https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=98-24
- URL (margin): https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=98-44
- URL (gap): https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=98-51
- Tokens: padding, margin, gap

### Spacing Scale (node: `102:100`)
- URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=102-100
- Tokens: s, m, l, xl, xxl

### Corners (node: `127:20`)
- URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=127-20
- Tokens: s, m, l, xl

## Components

### Logo (page: `107:147`, component: `107:201`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=107-147
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=107-201
- Pizza icon + "PIZZAVIBE" text (Knewave font, bold italic), 375×62
- No variants

### Header (page: `111:303`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=111-303
- **Header** (component: `111:390`)
  - Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=111-390
  - Full-width bar: Logo (left) + Navigation (right)
  - Padding: `--space-padding` (48px) vertical
  - Nav gap: `--space-gap` (48px)
  - Layout: flex, justify-between, align-end
  - No variants
- **HeaderNavItem** (component set: `112:53`)
  - Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=112-53
  - Font: H3 (Geist, 20px, extrabold/800, uppercase)
  - Padding: `--space-spacing-l` (24px) vertical
  - Variants:
    - `State=Default` (`112:52`): text color `--color-text-subtle`, no border
    - `State=Hover` (`112:54`): text color `--color-text-default`, no border
    - `State=Active` (`112:57`): text color `--color-text-primary-default`, border-bottom 4px `--color-text-primary-default`

### Footer (page: `112:78`, component: `112:177`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=112-78
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=112-177
- Top border: `--border-width-thick` solid `--color-border-default`
- Layout: flex, justify-between, align-center
- Padding: `--space-padding` (48px) vertical
- Left: Logo (small size, ~50% of default)
- Right ("Footer Info", `112:173`): gap `--space-gap` (48px)
  - "PizzaVibe V1.0" (`112:174`): Body Small typography
  - "GitHub" (`112:175`): Body Small typography, underlined link
- No variants

### Button (page: `114:154`, component set: `114:184`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=114-154
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=114-184
- Font: H3 (Geist, 20px, extrabold/800, uppercase)
- Padding: `--space-spacing-l` (24px)
- Border radius: 90px (pill)
- Variants:
  - `State=Default` (`114:183`): bg `--color-background-primary-default`, text `--color-text-default`, no border
  - `State=Hover` (`114:185`): bg `--color-background-primary-default`, text `--color-text-default`, border 4px `--color-border-primary-subtle`
  - `State=Active` (`114:187`): bg `--color-background-primary-default`, text `--color-text-primary-default`, border 4px `--color-border-primary-subtle`
  - `State=Disabled` (`114:190`): bg `--color-background-default`, text `--color-text-disabled`, border 4px `--color-border-disabled`, opacity 0.5
- New tokens discovered: `--color-border-primary-subtle` (#068f52), `--color-border-disabled` (#9faaa5), `--color-text-disabled` (#9faaa5)

### Tabs (page: `113:46`, component: `114:142`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=113-46
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=114-142
- Container: bg `--color-background-subtle`, border `--border-width-thick` solid `--color-border-default`, border-radius `--corner-xl`
- Padding top: `--space-spacing-xxl` (40px)
- Inner "Button Wrapper": bg `--color-background-default`, border `--border-width-thick` `--color-border-default` (top/left/right only), border-radius `--corner-xl` top corners, gap `--space-spacing-s`, padding `--space-spacing-l`
- Layout: flex, justify-center
- No variants
- **TabItem** (component set: `114:132`)
  - Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=114-132
  - Font: H3 (Geist, 20px, extrabold/800, uppercase)
  - Padding: `--space-spacing-l` (24px)
  - Border radius: 90px (pill)
  - Variants:
    - `State=Default` (`114:131`): border `--border-width-thick` `--color-border-default`, text `--color-text-default`
    - `State=Hover` (`114:136`): bg `--color-background-default`, border `--border-width-thick` `--color-border-primary-default`, text `--color-text-primary-default`
    - `State=Active` (`114:130`): bg `--color-background-inverted-default`, text `--color-text-inverted-default`, no border
- New tokens discovered: `--color-background-subtle` (#9faaa5), `--color-border-primary-default` (#045934)

### Icons (page: `116:281`, frame: `117:470`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=116-281
- Frame URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=117-470
- Container size: `--space-spacing-xxl` (40px)
- Fill color: `--color-text-default` (#0c0d0d)
- Icons:
  - `Icon/minus` (`117:477`): horizontal bar, viewBox 24×4
  - `Icon/add` (`117:476`): plus sign, viewBox 24×24
  - `Icon/delete` (`117:557`): trash can, viewBox 20×20
- No variants
- No new tokens

### Quantity Selector (page: `116:256`, component set: `117:613`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=116-256
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=117-613
- Container: border `--border-width-thick` `--color-border-default`, border-radius `--corner-xl`, gap `--space-spacing-s`, padding `--space-spacing-s`
- Quantity text: H3 typography, centered, min-width `--space-spacing-l`
- Layout: inline-flex, align-center
- Variants:
  - `Property 1=Default` (`117:612`): minus button (default type) + quantity + add button
  - `Property 1=Delete` (`117:614`): delete button (delete type) + quantity + add button
- **SelectorButton** (component set: `117:496`)
  - Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=117-496
  - Size: `--space-spacing-xxl` (40px), border-radius `--corner-xl`
  - Type=Default variants:
    - `State=Default` (`117:495`): border `--border-width-thick` `--color-border-default`, icon `--color-text-default`
    - `State=Hover` (`117:497`): bg `--color-background-subtle`, border `--color-border-default`
    - `State=Active` (`117:501`): bg `--color-background-inverted-default`, icon `--color-text-inverted-default`
    - `State=Disabled` (`117:505`): border `--color-border-disabled`, icon `--color-text-disabled`
  - Type=Delete variants:
    - `State=Default` (`117:530`): bg `--color-background-secondary-default`, border `--color-border-secondary-subtle`, icon `--color-text-inverted-default`
    - `State=Hover` (`117:532`): bg `--color-background-secondary-default`, border `--color-border-secondary-default`
    - `State=Active` (`117:534`): bg `--color-background-secondary-default`, border `--color-border-secondary-subtle`
    - `State=Disabled` (`117:536`): bg `--color-background-subtle`, border `--color-border-disabled`
- New tokens discovered: `--color-border-secondary-default` (#591304), `--color-border-secondary-subtle` (#b82508)

### CartItem (page: `116:17`, component: `118:143`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=116-17
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=118-143
- Layout: flex row, align-center, gap `--space-spacing-l`, padding vertical `--space-spacing-l`
- Bottom border: dashed `--border-width-thick` `--color-border-default`
- Width: 100%
- Children:
  - Name (flex: 1 0 0): Body Default typography, `--color-text-default`, overflow ellipsis
  - Price (shrink-0): Body Default typography, `--color-text-subtle`
  - QuantitySelector (shrink-0): existing component with `deleteAtMin`
- No variants
- No new tokens

### EmptyBlock (page: `116:251`, component: `117:433`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=116-251
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=117-433
- Border: dotted `--border-width-thick` `--color-border-default`
- Border radius: `--corner-xl`
- Padding: `--space-spacing-l`
- Width: 100%, height: content-driven
- No variants
- No new tokens

### PizzaItem (page: `115:9`, component set: `115:124`)
- Page URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=115-9
- Component URL: https://www.figma.com/design/Iia6bIqfQwSvXxTnfedTXj/PizzaVibe-UI-Kit?node-id=115-124
- Size: 300×394 fixed width
- Layout: flex column, gap `--space-spacing-l`
- **Illustration** (300×300 circle): 5 concentric rings built with CSS
  - Outer ring: bg `--color-background-subtle`, border `--border-width-thick` `--color-border-default`, round, padding `--space-spacing-m`
  - Middle ring: same as outer
  - Plate ring: bg `--color-background-default`, same border, `box-shadow: 0 7px 0 rgba(0,0,0,0.5)`
  - Sauce ring: bg `--color-background-secondary-default`
  - Cheese/Inner: bg `--color-background-tertiary-default` — accepts topping image
- **Add Button** (hover only): positioned top-right of illustration
  - Outer container: bg `--color-background-default`, border `--border-width-thick` `--color-border-subtle`, round pill, padding `--space-spacing-m`
  - Inner: uses existing Button component with "Add" text
- **Info Container**: gap `--space-spacing-s`
  - Header: flex row, H3 typography uppercase, name (flex-1 left) + price (flex-1 right), `--color-text-default`
  - Description: Body Small typography, `--color-text-subtle`
- Variants (2 states × 4 types):
  - `State=Default`: no add button visible
  - `State=Hover`: add button visible top-right
  - Types: Margherita (`115:123`), Pepperoni (`115:237`), Hawaiian (`115:397`), Vegan (`115:603`)
- New tokens discovered: `--color-background-tertiary-default` (#f2d96b)

