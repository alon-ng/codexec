import { useCallback, useEffect, useRef, useState } from 'react';
import type { ExercisesExerciseCodeData } from '~/api/generated/model';

interface ExecuteResponse {
  job_id: string;
  stdout: string;
  stderr: string;
  exit_code: number;
  time: number;
  memory: number;
  cpu: number;
}

export const useWebSocket = () => {
  const [lastResult, setLastResult] = useState<ExecuteResponse | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    // Determine WS protocol based on current protocol
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    // Use window.location.host to connect to the same host/port (proxied by Vite or Nginx)
    // Assuming /ws is proxied to backend /ws or /api/v1/ws.
    // The backend route is /ws but attached to gin router which usually has prefix /api/v1?
    // In `pkg/api/v1/routes.go`: `v1 := r.Group("/api/v1")`, but `/ws` was registered on `r` directly?
    // Let's check `routes.go` update.
    // `r.GET("/ws", ...)` was added on `r`. So it is at root `/ws`.
    // Wait, in `main.go`, `router := api.NewRouter(...)` returns `*gin.Engine`.
    // So `/ws` is at root.
    // However, if we run locally, frontend is 5173, backend is 8080.
    // Vite proxy usually proxies /api. Does it proxy /ws?

    // Let's assume we connect to `/ws`.
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    // In development with Vite, we might need to point to backend port if proxy isn't set up for /ws.
    // But typically we configure proxy for /ws too.

    const socket = new WebSocket(wsUrl);
    socketRef.current = socket;

    socket.onopen = () => {
      setIsConnected(true);
    };

    socket.onmessage = (event) => {
      try {
        const response = JSON.parse(event.data);
        // Check if it looks like ExecuteResponse
        if (response.job_id) {
          setLastResult(response);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    socket.onclose = () => {
      setIsConnected(false);
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      socket.close();
    };
  }, []);

  const submitCode = useCallback((exerciseUuid: string, submission: ExercisesExerciseCodeData) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ "exercise_uuid": exerciseUuid, "submission": submission }));
    } else {
      console.error('WebSocket is not connected');
    }
  }, []);

  return { submitCode, lastResult, isConnected };
};
