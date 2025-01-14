import type {
  Peer,
  PeerInit,
  UserAccount,
  APIKey,
  LoginBody,
  ServerInfo,
  UserAccountWithPass,
  APIKeyInit,
  APIKeyWithToken
} from "@/types/shared";

export async function POST_PreLogin(): Promise<Response> {
  const response = await fetch("/api/v1/prelogin", {
    method: "POST"
  });
  if (!response.ok) {
    const err = await response.json();
    throw err;
  }

  return response;
}

export async function POST_Login(body: LoginBody): Promise<Response> {
  const response = await fetch("/api/v1/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(body)
  });
  if (!response.ok) {
    const err = await response.json();
    throw err;
  }

  return response;
}

export async function POST_Logout(): Promise<void> {
  const response = await fetch("/api/v1/logout", {
    method: "POST"
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function GET_Peers(): Promise<Peer[]> {
  const response = await fetch("/api/v1/peers");
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}

export async function GET_Peer(uuid: string): Promise<Peer> {
  const response = await fetch("/api/v1/peers/" + uuid);
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}

export async function PUT_Peer(peer: Peer): Promise<void> {
  const response = await fetch("/api/v1/peers/" + peer.uuid, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(peer)
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function PATCH_Peer(peer: Peer): Promise<void> {
  const response = await fetch("/api/v1/peers/" + peer.uuid, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(peer)
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function DELETE_Peer(uuid: string): Promise<void> {
  const response = await fetch("/api/v1/peers/" + uuid, {
    method: "DELETE"
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function GET_PeerInit(): Promise<PeerInit> {
  const response = await fetch("/api/v1/peers/init");
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}

export async function GET_Accounts(): Promise<UserAccount[]> {
  const response = await fetch("/api/v1/accounts");
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}

export async function PUT_Account(account: UserAccountWithPass): Promise<void> {
  const response = await fetch("/api/v1/accounts/" + account.email, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(account)
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function PATCH_Account(account: UserAccount): Promise<void> {
  const response = await fetch("/api/v1/accounts/" + account.email, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(account)
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function PATCH_AccountPassword(email: string, password: string): Promise<void> {
  const response = await fetch("/api/v1/accounts/" + email + "/password", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ password })
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function DELETE_Account(email: string): Promise<void> {
  const response = await fetch("/api/v1/accounts/" + email, {
    method: "DELETE"
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function GET_APIKeys(): Promise<APIKey[]> {
  const response = await fetch("/api/v1/apikeys");
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}

export async function PUT_APIKey(key: APIKeyWithToken): Promise<void> {
  const response = await fetch("/api/v1/apikeys/" + key.uuid, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(key)
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function PATCH_APIKey(key: APIKey): Promise<void> {
  const response = await fetch("/api/v1/apikeys/" + key.uuid, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(key)
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function DELETE_APIKey(uuid: string): Promise<void> {
  const response = await fetch("/api/v1/apikeys/" + uuid, {
    method: "DELETE"
  });
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return;
}

export async function GET_APIKeyInit(): Promise<APIKeyInit> {
  const response = await fetch("/api/v1/apikeys/init");
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}

export async function GET_ServerInfo(): Promise<ServerInfo> {
  const response = await fetch("/api/v1/serverinfo");
  if (!response.ok) {
    const err = await response.text();
    throw new Error(err);
  }

  return response.json();
}
