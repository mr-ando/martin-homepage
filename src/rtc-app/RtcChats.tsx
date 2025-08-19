import {
  Button,
  Stack,
  Container,
  Flex,
  Input,
  Table,
  Text,
  chakra,
  Box,
  Field,
} from '@chakra-ui/react';
import { useEffect, useRef, useState } from 'react';
import ResizeTextarea from 'react-textarea-autosize';

const AutoResizeTextArea = chakra(ResizeTextarea);

import { useChats } from '../hooks/useChats';

export default function RtcChats({ username }: { username: string }) {
  const { chats, createChat, joinChat, currentChatInfo, sendMessage } =
    useChats();
  const [chatName, setChatName] = useState('');
  const [chatNameHidden, setChatNameHidden] = useState(true);
  const chatNameRef = useRef<HTMLInputElement>(null);
  const chatBoxRef = useRef<HTMLTextAreaElement>(null);

  const handleCreateChat = async () => {
    await createChat(chatName);
    setChatName('');
    setChatNameHidden(true);
  };

  const handleSendMessage = () => {
    if (chatBoxRef.current) {
      sendMessage(chatBoxRef.current?.value);
      chatBoxRef.current.value = '';
    }
  };

  useEffect(() => {
    if (!chatNameHidden) {
      console.log('is activated, current:' + chatNameRef.current);
      chatNameRef.current?.focus();
    }
  }, [chatNameHidden]);

  return (
    <Container>
      <Flex height="full" flexDir="column" gap="4">
        <Text textAlign="center" fontSize="4xl">
          Real-time chat application
        </Text>
        <Flex height="full" flexDir="row" gap="5">
          <Flex flexDir="column">
            <Box height="md" maxHeight="md" overflowY="auto" border="md">
              <Table.Root minW="md" showColumnBorder>
                <Table.Caption />
                <Table.Header>
                  <Table.Row>
                    <Table.ColumnHeader>Name</Table.ColumnHeader>
                    <Table.ColumnHeader width="5" textAlign="center">
                      People
                    </Table.ColumnHeader>
                  </Table.Row>
                </Table.Header>
                <Table.Body gap="0">
                  {chats?.chats.map((c) => {
                    return (
                      <Table.Row
                        key={c}
                        onClick={async () => {
                          joinChat(c, username);
                        }}
                        _hover={{ border: 'md' }}
                        border={c === currentChatInfo?.chatName ? 'md' : ''}
                        bgColor={
                          c == currentChatInfo?.chatName ? 'green.950' : ''
                        }
                        flex="flex"
                        role="button"
                        cursor="pointer"
                      >
                        <Table.Cell>{c}</Table.Cell>
                        <Table.Cell textAlign="center">0</Table.Cell>
                      </Table.Row>
                    );
                  })}
                </Table.Body>
              </Table.Root>
            </Box>
            <Flex gap="1">
              <Button
                width="fit-content"
                onClick={() => {
                  setChatNameHidden(!chatNameHidden);
                }}
              >
                {chatNameHidden ? '+' : '-'}
              </Button>
              <Field.Root
                rounded="md"
                position="relative"
                backgroundColor="white"
                width="fit"
                flexDir="row"
                border="sm"
                hidden={chatNameHidden}
              >
                <Input
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      handleCreateChat();
                    }
                  }}
                  border="none"
                  ref={chatNameRef}
                  color="black"
                  width="2xs"
                  value={chatName}
                  outline="none"
                  onChange={(e) => setChatName(e.target.value)}
                />
                <Button onClick={handleCreateChat} colorPalette="blue">
                  Add Chat
                </Button>
              </Field.Root>
            </Flex>
          </Flex>
          <Flex
            border="sm"
            gap="4"
            flexDir="column"
            padding="4"
            rounded="md"
            flexGrow="1"
            minHeight="md"
          >
            <Text>Chat</Text>
            <Stack overflowY="auto" maxHeight="sm" minH="sm">
              {currentChatInfo?.messages.map((m, i) => {
                return (
                  <Flex key={i} gap="2">
                    <Text color="red.500">{m.sender}:</Text>
                    <Text>{m.message}</Text>
                    <Text marginLeft="auto" color="red.600">
                      {m.timestamp}
                    </Text>
                  </Flex>
                );
              })}
            </Stack>
            <Flex
              bgColor="gray.800"
              onClick={() => {}}
              position="relative"
              marginTop="auto"
              border="sm"
              rounded="sm"
              alignItems="center"
              gap="2"
              padding="2"
            >
              <AutoResizeTextArea
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleSendMessage();
                  }
                }}
                w="full"
                rows={1}
                maxRows={12}
                ref={chatBoxRef}
                bg="none"
                fontSize="lg"
                border="none"
                resize="none"
                outline="none"
                scrollbar="hidden"
              />
              <Button onClick={handleSendMessage} textAlign="right">
                Send
              </Button>
            </Flex>
          </Flex>
        </Flex>
      </Flex>
    </Container>
  );
}
