import { useNavigate } from "react-router-dom";

declare global {
  interface Window {
    navigate: ReturnType<typeof useNavigate> | undefined;
  }

  // __BINDPLANE_VERSION__ is set by _globals.js at load time
  var __BINDPLANE_VERSION__: string;
}
