import { gql } from "apollo-boost";

import { auth, logout, setAuth, user } from "./auth";
import { cart, changeCartQuantity, clearCart, removeFromCart } from "./cart";

export const typeDefs = gql`
  type LocalUser {
    id: Int
    email: String
    role: String
  }

  type AuthState {
    accessToken: String
    refreshToken: String
    user: LocalUser
  }

  type CartVariant {
    variantId: Int!
    quantity: Int!
  }

  type CartState {
    variants: [CartVariant!]!
  }

  extend type Query {
    auth: AuthState!
    cart: CartState!
  }

  extend type Mutation {
    logout: Boolean
    setAuth(
      accessToken: String
      refreshToken: String
      stayLoggedIn: Boolean!
    ): Boolean
    changeCartQuantity(variantId: Int!, quantity: Int!): Boolean
    clearCart: Boolean
    removeFromCart(variantId: Int!): Boolean
  }
`;

export const resolvers = {
  AuthState: {
    user
  },
  Query: {
    auth,
    cart
  },
  Mutation: {
    logout,
    setAuth,
    changeCartQuantity,
    clearCart,
    removeFromCart
  }
};
