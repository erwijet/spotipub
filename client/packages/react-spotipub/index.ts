import { useEffect, useRef, useState } from "react";
import { payloadSchema, Payload } from "zod-spotipub";

function useInterval(timeout: number, fn: () => unknown) {
  const ref = useRef<number>();
  useEffect(() => {
    ref.current = setInterval(fn, timeout);
    return () => clearInterval(ref.current);
  });
}

export type UseSpotipubOptions = {
  src: string;
};

export function useSpotipub(options: UseSpotipubOptions) {
  const [state, setState] = useState<Payload>();

  useInterval(500, () => {
    setState((s) => {
      if (s) {
        return { ...s, progress_ms: s.progress_ms + 500 };
      }

      return s;
    });
  });

  useEffect(() => {
    const evtsrc = new EventSource(options.src);

    function handleUpdate(evt: MessageEvent) {
      const payload = payloadSchema.parse(JSON.parse(evt.data));
      setState(payload);
    }

    evtsrc.addEventListener("initial", handleUpdate);
    evtsrc.addEventListener("update", handleUpdate);

    return () => {
      evtsrc.removeEventListener("initial", handleUpdate);
      evtsrc.removeEventListener("update", handleUpdate);
    };
  }, []);

  return state;
}
