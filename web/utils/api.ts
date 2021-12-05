export const fetcher = (path: string) =>
  fetch(`http://localhost:1317${path}`).then((res) => res.json());
