export function cart() {
  const saved = JSON.parse(localStorage.getItem("cart") || "null");

  return (
    saved || {
      variants: [],
      __typename: "CartState"
    }
  );
}

export function clearCart(_, __, { cache }) {
  const existing = cart();

  const newCart = {
    ...existing,
    variants: []
  };

  localStorage.setItem("cart", JSON.stringify(newCart));
  cache.writeData({ data: { cart: newCart } });

  return true;
}

export function removeFromCart(_, { variantId }, { cache }) {
  const existing = cart();

  const newCart = {
    ...existing,
    variants: existing.variants.filter(
      variant => variant.variantId !== variantId
    )
  };

  localStorage.setItem("cart", JSON.stringify(newCart));
  cache.writeData({ data: { cart: newCart } });

  return true;
}

export function changeCartQuantity(_, { variantId, quantity }, { cache }) {
  const existing = cart();

  let found = false;
  const variants = existing.variants.map(variant => {
    if (variant.variantId === variantId) {
      found = true;

      const newQuantity = variant.quantity + quantity;

      return {
        ...variant,
        quantity: newQuantity < 1 ? 1 : newQuantity
      };
    }

    return variant;
  });

  if (!found) {
    variants.push({
      variantId,
      quantity: quantity < 1 ? 1 : quantity,
      __typename: "CartVariant"
    });
  }

  const newCart = {
    ...existing,
    variants
  };

  localStorage.setItem("cart", JSON.stringify(newCart));
  cache.writeData({ data: { cart: newCart } });

  return true;
}
