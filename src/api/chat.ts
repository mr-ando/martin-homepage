export async function createChatApi(chatName: string) {
  const res = await fetch(`http://localhost:8080/api/chat/create/${chatName}`, {
    method: 'POST',
  });
}

interface ChatsResponse {
  chats: string[];
}

export async function fetchChatsApi() {
  const res = await fetch('http://localhost:8080/api/chats', {
    method: 'GET',
  })
    .then((res) => res.json())
    .then((res) => res as ChatsResponse);
  return res;
}

export async function joinChatApi(
  chatName: string,
  clientName: string
): Promise<WebSocket> {
  const ws = new WebSocket(
    `ws://localhost:8080/api/chat/join/${chatName}?clientName=${clientName}`
  );

  return ws;
}

export async function getChatApi(chatName: string): Promise<ChatInformation> {
  const chatInfo = await fetch(`http://localhost:8080/api/chat/${chatName}`, {
    method: 'GET',
  })
    .then((res) => res.json())
    .then((res) => res as ChatInformation);

  return chatInfo;
}

export interface Message {
  sender: string;
  message: string;
  timestamp: string;
}

export interface ChatInformation {
  chatName: string;
  messages: Message[];
  participants: string[];
}
