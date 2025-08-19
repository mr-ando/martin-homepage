'use client';

import { Button, Stack, Input, Text, Field } from '@chakra-ui/react';
import { useEffect, useRef, useState } from 'react';
import RtcChats from './RtcChats';

export default function RTCApp() {
  const [userName, setUserName] = useState<string | null>(null);
  const usernameRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    setUserName(localStorage.getItem('username'));
  }, []);

  const handleSetUsername = () => {
    if (usernameRef.current?.value) {
      localStorage.setItem('username', usernameRef.current.value);
      setUserName(usernameRef.current?.value);
    }
  };
  return (
    <Stack height="full">
      {userName == undefined ? (
        <Stack>
          <Text textAlign="center" fontSize="2xl">
            Real-time chat application
          </Text>
          <Field.Root>
            <Field.Label>username</Field.Label>
            <Field.Root flexDir="row">
              <Input ref={usernameRef} border="sm" maxW="sm" />
              <Button onClick={handleSetUsername}>Confirm</Button>
            </Field.Root>
          </Field.Root>
        </Stack>
      ) : (
        <RtcChats username={userName} />
      )}
    </Stack>
  );
}
