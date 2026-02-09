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

