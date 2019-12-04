import decode from "jwt-decode";

let saved = JSON.parse(localStorage.getItem("auth") || "null");

export function auth() {
  return (
    saved || {
      accessToken: null,
      refreshToken: null,
      __typename: "AuthState"
    }
  );
}

export function user() {
  const { accessToken } = auth();

  if (accessToken) {
    return {
      ...decode(accessToken),
      __typename: "LocalUser"
    };
  }

  return null;
}

export function setAuth(
  _,
  { accessToken, refreshToken, stayLoggedIn },
  { cache }
) {
  const newAuth = { accessToken, refreshToken, __typename: "AuthState" };

  if (stayLoggedIn) {
    localStorage.setItem("auth", JSON.stringify(newAuth));
  } else {
    localStorage.removeItem("auth");
  }

  saved = newAuth;
  cache.writeData({ data: { auth: newAuth } });

  return true;
}

export function logout(_, __, { cache }) {
  localStorage.removeItem("auth");
  saved = null;

  cache.writeData({ data: { auth: auth() } });

  return true;
}
