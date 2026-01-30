import { useCallback, useEffect, useRef, useState } from 'react';
import type { ModelsExerciseCodeData } from '~/api/generated/model';
import type { ExecuteResponse, UserExerciseQuizData } from '~/api/types';


export const useWebSocket = (onSubmissionResponse?: (result: ExecuteResponse) => void) => {
  const [lastResult, setLastResult] = useState<ExecuteResponse | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    let timeoutId: NodeJS.Timeout;
    let socket: WebSocket | null = null;
    let attempts = 0;
    let isMounted = true;

    const connect = () => {
      if (!isMounted) return;

      // Determine WS protocol based on current protocol
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/api/v1/ws`;

      socket = new WebSocket(wsUrl);
      socketRef.current = socket;

      socket.onopen = () => {
        if (!isMounted) return;
        setIsConnected(true);
        attempts = 0;
      };

      socket.onmessage = (event) => {
        if (!isMounted) return;
        try {
          const response = JSON.parse(event.data);
          // Check if it looks like ExecuteResponse
          if (response.job_id) {
            setLastResult(response);
            onSubmissionResponse?.(response);
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      socket.onclose = () => {
        if (!isMounted) return;
        setIsConnected(false);

        // Calculate backoff delay with a cap (e.g., max 30 seconds)
        const delay = Math.min(1000 * Math.pow(2, attempts), 30000);
        attempts++;

        console.log(`WebSocket disconnected. Retrying in ${delay}ms...`);
        timeoutId = setTimeout(connect, delay);
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    };

    connect();

    return () => {
      isMounted = false;
      clearTimeout(timeoutId);
      if (socket) {
        socket.onclose = null;
        socket.onerror = null;
        socket.close();
      }
    };
  }, []);

  const submit = useCallback((exerciseUuid: string, submission: ModelsExerciseCodeData | UserExerciseQuizData) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ "exercise_uuid": exerciseUuid, "submission": submission }));
    } else {
      console.error('WebSocket is not connected');
    }
  }, []);

  return { submit, lastResult, isConnected };
};
