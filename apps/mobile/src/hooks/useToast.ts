import { useCallback, useEffect, useRef, useState } from 'react';

// Toast แบบ auto-dismiss — ใช้ร่วมกับ <Toast/> component
export function useToast(duration = 2000) {
  const [message, setMessage] = useState<string | null>(null);
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const show = useCallback(
    (msg: string) => {
      setMessage(msg);
      if (timer.current) {
        clearTimeout(timer.current);
      }
      timer.current = setTimeout(() => setMessage(null), duration);
    },
    [duration],
  );

  useEffect(
    () => () => {
      if (timer.current) {
        clearTimeout(timer.current);
      }
    },
    [],
  );

  return { message: message ?? '', visible: message !== null, show };
}
