import { useState } from 'react';
import {
  createChatApi,
  fetchChatsApi,
  getChatApi,
  joinChatApi,
  type ChatInformation,
  type Message,
} from '../api/chat';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export function useChats() {
  const [currentChatInfo, setCurrentChatInfo] = useState<
    ChatInformation | undefined
  >(undefined);

  const [chatConnection, setChatConnection] = useState<WebSocket | undefined>(
    undefined
  );

  const queryClient = useQueryClient();
  const { mutateAsync: createChat } = useMutation({
    mutationKey: ['chats'],
    mutationFn: createChatApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chats'] });
    },
  });

  const { data: chats } = useQuery({
    queryKey: ['chats'],
    queryFn: fetchChatsApi,
    refetchInterval: 5000,
  });

  const joinChat = async (chatName: string, clientName: string) => {
    if (chatConnection) {
      chatConnection.close(1000);
    }
    const chatInfo = await getChatApi(chatName);
    console.log(chatInfo.chatName);
    setCurrentChatInfo(chatInfo);

    const ws = await joinChatApi(chatName, clientName);
    setChatConnection(ws);

    ws.addEventListener('message', function (event) {
      console.log('message received:', event.data);
      const m = JSON.parse(event.data) as Message;
      console.log(m);
      setCurrentChatInfo((prevChatInfo) => {
        if (prevChatInfo) {
          const newMessages = [...prevChatInfo?.messages, m];
          return {
            chatName: prevChatInfo.chatName,
            participants: prevChatInfo?.participants,
            messages: newMessages,
          };
        }
        return prevChatInfo;
      });
    });
  };

  const sendMessage = async (message: string) => {
    if (chatConnection) {
      chatConnection.send(message);
    }
  };

  return {
    chats,
    createChat,
    joinChat,
    currentChatInfo,
    sendMessage,
  };
}
