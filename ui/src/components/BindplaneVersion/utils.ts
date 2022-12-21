export function version(): string {
  try {
    return __BINDPLANE_VERSION__;
  } catch (err) {
    return "unknown";
  }
}
